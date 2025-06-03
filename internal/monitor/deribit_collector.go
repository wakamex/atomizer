package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type DeribitCollector struct {
	client    *http.Client
	baseURL   string
	converter *InstrumentConverter
}

type DeribitResponse struct {
	Result []DeribitTicker `json:"result"`
}

type DeribitTicker struct {
	InstrumentName string  `json:"instrument_name"`
	BestBidPrice   float64 `json:"best_bid_price"`
	BestAskPrice   float64 `json:"best_ask_price"`
	BestBidAmount  float64 `json:"best_bid_amount"`
	BestAskAmount  float64 `json:"best_ask_amount"`
	LastPrice      float64 `json:"last_price"`
	Stats          struct {
		Volume      float64 `json:"volume"`
		High        float64 `json:"high"`
		Low         float64 `json:"low"`
		PriceChange float64 `json:"price_change"`
	} `json:"stats"`
	MarkPrice float64 `json:"mark_price"`
}

func NewDeribitCollector() *DeribitCollector {
	return &DeribitCollector{
		client:    NewHTTPClient(10 * time.Second),
		baseURL:   "https://www.deribit.com/api/v2",
		converter: NewInstrumentConverter(),
	}
}

func (d *DeribitCollector) Name() string {
	return "deribit"
}

func (d *DeribitCollector) Collect(ctx context.Context, instruments []string) ([]Metric, error) {
	// Get all active instruments first
	instrumentsData, err := d.getInstruments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get instruments: %w", err)
	}

	// Convert input patterns to Deribit format
	deribitPatterns := d.converter.ConvertInstrumentList(instruments, "deribit")

	// Filter instruments based on converted patterns
	filteredInstruments := FilterInstruments(instrumentsData, deribitPatterns)

	if len(filteredInstruments) == 0 {
		return []Metric{}, nil
	}

	// Batch fetch ticker data
	metrics := []Metric{}
	for i := 0; i < len(filteredInstruments); i += 100 {
		end := i + 100
		if end > len(filteredInstruments) {
			end = len(filteredInstruments)
		}

		batch := filteredInstruments[i:end]
		batchMetrics, err := d.fetchTickerBatch(ctx, batch)
		if err != nil {
			// Log error but continue with other batches
			fmt.Printf("Error fetching batch %d-%d: %v\n", i, end, err)
			continue
		}

		metrics = append(metrics, batchMetrics...)
	}

	return metrics, nil
}

func (d *DeribitCollector) getInstruments(ctx context.Context) ([]string, error) {
	// Get both BTC and ETH instruments
	allInstruments := []string{}

	for _, currency := range []string{"BTC", "ETH"} {
		url := fmt.Sprintf("%s/public/get_instruments?currency=%s&expired=false", d.baseURL, currency)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := d.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var result struct {
			Result []struct {
				InstrumentName string `json:"instrument_name"`
			} `json:"result"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, err
		}

		for _, inst := range result.Result {
			allInstruments = append(allInstruments, inst.InstrumentName)
		}
	}

	return allInstruments, nil
}

func (d *DeribitCollector) fetchTickerBatch(ctx context.Context, instruments []string) ([]Metric, error) {
	url := fmt.Sprintf("%s/public/ticker", d.baseURL)

	// Use concurrent fetching with worker pool
	const maxWorkers = 5
	jobs := make(chan string, len(instruments))
	results := make(chan *Metric, len(instruments))

	// Start workers
	var wg sync.WaitGroup
	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for instrument := range jobs {
				metric := d.fetchSingleTicker(ctx, url, instrument)
				if metric != nil {
					results <- metric
				}
			}
		}()
	}

	// Send jobs
	for _, instrument := range instruments {
		jobs <- instrument
	}
	close(jobs)

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	metrics := []Metric{}
	for metric := range results {
		if metric != nil {
			metrics = append(metrics, *metric)
		}
	}

	return metrics, nil
}

func (d *DeribitCollector) fetchSingleTicker(ctx context.Context, baseURL, instrument string) *Metric {
	tickerURL := fmt.Sprintf("%s?instrument_name=%s", baseURL, instrument)

	req, err := http.NewRequestWithContext(ctx, "GET", tickerURL, nil)
	if err != nil {
		return nil
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var tickerResp struct {
		Result DeribitTicker `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tickerResp); err != nil {
		return nil
	}

	ticker := tickerResp.Result

	// Calculate open price from price change
	openPrice := ticker.LastPrice
	if ticker.Stats.PriceChange != 0 {
		openPrice = ticker.LastPrice / (1 + ticker.Stats.PriceChange/100)
	}

	return &Metric{
		Exchange:   "deribit",
		Instrument: ticker.InstrumentName,
		Timestamp:  time.Now(),
		BidPrice:   ticker.BestBidPrice,
		AskPrice:   ticker.BestAskPrice,
		BidSize:    ticker.BestBidAmount,
		AskSize:    ticker.BestAskAmount,
		LastPrice:  ticker.LastPrice,
		Volume24h:  ticker.Stats.Volume,
		OpenPrice:  openPrice,
		HighPrice:  ticker.Stats.High,
		LowPrice:   ticker.Stats.Low,
	}
}
