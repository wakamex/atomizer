package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DeriveCollector struct {
	client    *http.Client
	baseURL   string
	converter *InstrumentConverter
}

type DeriveTickerResponse struct {
	Result map[string]DeriveTicker `json:"result"`
}

type DeriveTicker struct {
	InstrumentName string `json:"instrument_name"`
	BestBidPrice   string `json:"best_bid_price"`
	BestAskPrice   string `json:"best_ask_price"`
	BestBidAmount  string `json:"best_bid_amount"`
	BestAskAmount  string `json:"best_ask_amount"`
	LastPrice      string `json:"last_price"`
	MarkPrice      string `json:"mark_price"`
	IndexPrice     string `json:"index_price"`
	Stats          struct {
		Volume      string `json:"volume"`
		High        string `json:"high"`
		Low         string `json:"low"`
		PriceChange string `json:"price_change"`
	} `json:"stats"`
}

func NewDeriveCollector() *DeriveCollector {
	return &DeriveCollector{
		client:    NewHTTPClient(10 * time.Second),
		baseURL:   "https://api.lyra.finance",
		converter: NewInstrumentConverter(),
	}
}

func (d *DeriveCollector) Name() string {
	return "derive"
}

func (d *DeriveCollector) Collect(ctx context.Context, instruments []string) ([]Metric, error) {
	// First, get all instruments
	allInstruments, err := d.getAllInstruments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get instruments: %w", err)
	}

	// Convert input patterns to Derive format
	derivePatterns := d.converter.ConvertInstrumentList(instruments, "derive")

	// Filter instruments based on converted patterns
	filteredInstruments := FilterInstruments(allInstruments, derivePatterns)

	if len(filteredInstruments) == 0 {
		return []Metric{}, nil
	}

	// Collect ticker data for each instrument
	metrics := []Metric{}
	for _, instrument := range filteredInstruments {
		ticker, err := d.getTicker(ctx, instrument)
		if err != nil {
			continue // Skip failed tickers
		}
		metrics = append(metrics, ticker)
	}

	return metrics, nil
}

func (d *DeriveCollector) getAllInstruments(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/public/get_all_instruments", d.baseURL)
	instruments := []string{}
	page := 1

	for {
		// Derive uses POST with JSON payload and pagination
		payload := map[string]interface{}{
			"instrument_type": "option", // Get only options
			"page":            page,
			"page_size":       1000,
			"expired":         false,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := d.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Read the response body for debugging
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		var result struct {
			Result struct {
				Instruments []struct {
					InstrumentName string `json:"instrument_name"`
					IsActive       bool   `json:"is_active"`
					InstrumentType string `json:"instrument_type"`
				} `json:"instruments"`
				Pagination struct {
					NumPages    int `json:"num_pages"`
					CurrentPage int `json:"current_page"`
				} `json:"pagination"`
			} `json:"result"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		// Add active instruments
		for _, inst := range result.Result.Instruments {
			if inst.IsActive {
				instruments = append(instruments, inst.InstrumentName)
			}
		}

		// Check if there are more pages
		if page >= result.Result.Pagination.NumPages || result.Result.Pagination.NumPages == 0 {
			break
		}
		page++
	}

	return instruments, nil
}

func (d *DeriveCollector) getTicker(ctx context.Context, instrument string) (Metric, error) {
	url := fmt.Sprintf("%s/public/get_ticker", d.baseURL)

	payload := fmt.Sprintf(`{"instrument_name": "%s"}`, instrument)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return Metric{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return Metric{}, err
	}
	defer resp.Body.Close()

	var tickerResp struct {
		Result DeriveTicker `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tickerResp); err != nil {
		return Metric{}, err
	}

	ticker := tickerResp.Result

	// Parse string values to float64
	parseFloat := func(s string) float64 {
		var f float64
		fmt.Sscanf(s, "%f", &f)
		return f
	}

	bidPrice := parseFloat(ticker.BestBidPrice)
	askPrice := parseFloat(ticker.BestAskPrice)
	bidSize := parseFloat(ticker.BestBidAmount)
	askSize := parseFloat(ticker.BestAskAmount)
	lastPrice := parseFloat(ticker.LastPrice)
	volume := parseFloat(ticker.Stats.Volume)
	high := parseFloat(ticker.Stats.High)
	low := parseFloat(ticker.Stats.Low)
	priceChange := parseFloat(ticker.Stats.PriceChange)

	// Calculate open price from price change
	openPrice := lastPrice
	if priceChange != 0 {
		openPrice = lastPrice / (1 + priceChange/100)
	}

	return Metric{
		Exchange:   "derive",
		Instrument: ticker.InstrumentName,
		Timestamp:  time.Now(),
		BidPrice:   bidPrice,
		AskPrice:   askPrice,
		BidSize:    bidSize,
		AskSize:    askSize,
		LastPrice:  lastPrice,
		Volume24h:  volume,
		OpenPrice:  openPrice,
		HighPrice:  high,
		LowPrice:   low,
	}, nil
}
