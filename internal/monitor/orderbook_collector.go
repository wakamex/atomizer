package monitor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// OrderBookLevel represents a single price level
type OrderBookLevel struct {
	Price  float64
	Size   float64
	Orders int // Number of orders at this level (if available)
}

// OrderBookMetric represents order book data for storage
type OrderBookMetric struct {
	Exchange   string
	Instrument string
	Timestamp  time.Time
	Bids       []OrderBookLevel // Sorted by price descending (best bid first)
	Asks       []OrderBookLevel // Sorted by price ascending (best ask first)
}

// DeribitOrderBookCollector collects full order book data
type DeribitOrderBookCollector struct {
	client        *http.Client
	baseURL       string
	converter     *InstrumentConverter
	depth         int
	spotCollector *DeriveSpotCollector
}

func NewDeribitOrderBookCollector(depth int) *DeribitOrderBookCollector {
	if depth <= 0 {
		depth = 10
	}
	return &DeribitOrderBookCollector{
		client:    NewHTTPClient(10 * time.Second),
		baseURL:   "https://www.deribit.com/api/v2",
		converter: NewInstrumentConverter(),
		depth:     depth,
	}
}

func (d *DeribitOrderBookCollector) Name() string {
	return "deribit-orderbook"
}

func (d *DeribitOrderBookCollector) SetSpotCollector(spotCollector *DeriveSpotCollector) {
	d.spotCollector = spotCollector
}

func (d *DeribitOrderBookCollector) CollectOrderBooks(ctx context.Context, instruments []string) ([]OrderBookMetric, error) {
	// Convert input patterns to Deribit format
	deribitInstruments := d.converter.ConvertInstrumentList(instruments, "deribit")

	metrics := []OrderBookMetric{}
	for _, instrument := range deribitInstruments {
		orderBook, err := d.getOrderBook(ctx, instrument)
		if err != nil {
			continue // Skip failed instruments
		}
		metrics = append(metrics, orderBook)
	}

	return metrics, nil
}

func (d *DeribitOrderBookCollector) getOrderBook(ctx context.Context, instrument string) (OrderBookMetric, error) {
	url := fmt.Sprintf("%s/public/get_order_book", d.baseURL)

	// Parameters for order book request
	params := fmt.Sprintf("?instrument_name=%s&depth=%d", instrument, d.depth)

	req, err := http.NewRequestWithContext(ctx, "GET", url+params, nil)
	if err != nil {
		return OrderBookMetric{}, err
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return OrderBookMetric{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Result struct {
			Bids [][]float64 `json:"bids"` // [price, amount]
			Asks [][]float64 `json:"asks"` // [price, amount]
		} `json:"result"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return OrderBookMetric{}, fmt.Errorf("failed to read response: %w", err)
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return OrderBookMetric{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Debug: Log the response
	if len(result.Result.Bids) > 0 || len(result.Result.Asks) > 0 {
		Debugf("[DEBUG] Deribit order book: %d bids, %d asks", len(result.Result.Bids), len(result.Result.Asks))
		if len(result.Result.Bids) > 0 && len(result.Result.Asks) > 0 {
			bidPrice := result.Result.Bids[0][0]
			askPrice := result.Result.Asks[0][0]
			bidSize := result.Result.Bids[0][1]
			askSize := result.Result.Asks[0][1]

			// Check if we can convert to USD for display
			priceUnit := "ETH"
			displayBidPrice := bidPrice
			displayAskPrice := askPrice

			if d.spotCollector != nil && strings.Contains(instrument, "ETH") {
				if ethSpot, ok := d.spotCollector.GetSpotPrice("ETH"); ok {
					displayBidPrice = bidPrice * ethSpot
					displayAskPrice = askPrice * ethSpot
					priceUnit = fmt.Sprintf("USD (ETH=$%.2f)", ethSpot)
				}
			}

			Debugf("  Best bid: %.2f @ %.2f, Best ask: %.2f @ %.2f %s",
				displayBidPrice, bidSize, displayAskPrice, askSize, priceUnit)
		}
	}

	// Convert to OrderBookMetric
	metric := OrderBookMetric{
		Exchange:   "deribit",
		Instrument: instrument,
		Timestamp:  time.Now(),
		Bids:       make([]OrderBookLevel, 0, len(result.Result.Bids)),
		Asks:       make([]OrderBookLevel, 0, len(result.Result.Asks)),
	}

	// Convert bids
	for _, bid := range result.Result.Bids {
		if len(bid) >= 2 {
			metric.Bids = append(metric.Bids, OrderBookLevel{
				Price: bid[0],
				Size:  bid[1],
			})
		}
	}

	// Convert asks
	for _, ask := range result.Result.Asks {
		if len(ask) >= 2 {
			metric.Asks = append(metric.Asks, OrderBookLevel{
				Price: ask[0],
				Size:  ask[1],
			})
		}
	}

	return metric, nil
}

// DeriveOrderBookCollector collects order book data from Derive
type DeriveOrderBookCollector struct {
	client    *http.Client
	baseURL   string
	converter *InstrumentConverter
	depth     int
}

func NewDeriveOrderBookCollector(depth int) *DeriveOrderBookCollector {
	if depth <= 0 {
		depth = 10
	}
	return &DeriveOrderBookCollector{
		client:    NewHTTPClient(10 * time.Second),
		baseURL:   "https://api.lyra.finance",
		converter: NewInstrumentConverter(),
		depth:     depth,
	}
}

func (d *DeriveOrderBookCollector) Name() string {
	return "derive-orderbook"
}

func (d *DeriveOrderBookCollector) CollectOrderBooks(ctx context.Context, instruments []string) ([]OrderBookMetric, error) {
	// Convert input patterns to Derive format
	deriveInstruments := d.converter.ConvertInstrumentList(instruments, "derive")

	metrics := []OrderBookMetric{}
	for _, instrument := range deriveInstruments {
		orderBook, err := d.getOrderBook(ctx, instrument)
		if err != nil {
			continue // Skip failed instruments
		}
		metrics = append(metrics, orderBook)
	}

	return metrics, nil
}

func (d *DeriveOrderBookCollector) getOrderBook(ctx context.Context, instrument string) (OrderBookMetric, error) {
	// Derive only provides order book data via WebSocket, not REST API
	// For REST polling, we can only get top-of-book from ticker
	return d.getOrderBookFromTicker(ctx, instrument)
}

// Fallback method if Derive doesn't have order book endpoint
func (d *DeriveOrderBookCollector) getOrderBookFromTicker(ctx context.Context, instrument string) (OrderBookMetric, error) {
	// Use the existing ticker endpoint as fallback
	url := fmt.Sprintf("%s/public/get_ticker", d.baseURL)

	payload := fmt.Sprintf(`{"instrument_name": "%s"}`, instrument)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return OrderBookMetric{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return OrderBookMetric{}, err
	}
	defer resp.Body.Close()

	var tickerResp struct {
		Result struct {
			BestBidPrice  string `json:"best_bid_price"`
			BestAskPrice  string `json:"best_ask_price"`
			BestBidAmount string `json:"best_bid_amount"`
			BestAskAmount string `json:"best_ask_amount"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tickerResp); err != nil {
		return OrderBookMetric{}, err
	}

	// Parse string values
	bidPrice, _ := strconv.ParseFloat(tickerResp.Result.BestBidPrice, 64)
	askPrice, _ := strconv.ParseFloat(tickerResp.Result.BestAskPrice, 64)
	bidSize, _ := strconv.ParseFloat(tickerResp.Result.BestBidAmount, 64)
	askSize, _ := strconv.ParseFloat(tickerResp.Result.BestAskAmount, 64)

	// Create single-level order book
	return OrderBookMetric{
		Exchange:   "derive",
		Instrument: instrument,
		Timestamp:  time.Now(),
		Bids: []OrderBookLevel{
			{Price: bidPrice, Size: bidSize},
		},
		Asks: []OrderBookLevel{
			{Price: askPrice, Size: askSize},
		},
	}, nil
}
