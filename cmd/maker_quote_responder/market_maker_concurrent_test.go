package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// Mock types needed for testing
type MarketMakerConfig struct {
	Instruments              []string
	QuoteSize                decimal.Decimal
	SpreadBps                int
	MinSpreadBps             int
	Improvement              decimal.Decimal
	ImprovementReferenceSize decimal.Decimal
	RefreshInterval          time.Duration
	CancelThreshold          decimal.Decimal
	MaxPositionSize          decimal.Decimal
	MaxTotalExposure         decimal.Decimal
	BidOnly                  bool
	AskOnly                  bool
}

type MarketMakerOrder struct {
	OrderID    string
	Instrument string
	Side       string
	Price      decimal.Decimal
	Amount     decimal.Decimal
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MarketMakerStats struct {
	OrdersPlaced    int64
	OrdersCancelled int64
	OrdersFilled    int64
	BidAskSpread    map[string]decimal.Decimal
	UptimeSeconds   int64
	LastUpdate      time.Time
}

type TickerUpdate struct {
	Instrument string
	BestBid    decimal.Decimal
	BestAsk    decimal.Decimal
	MarkPrice  decimal.Decimal
	Timestamp  time.Time
}

// MarketMakerExchange interface needed for testing
type MarketMakerExchange interface {
	PlaceLimitOrder(instrument, side string, price, amount decimal.Decimal) (string, error)
	CancelOrder(orderID string) error
}

// Minimal MarketMaker struct for testing
type MarketMaker struct {
	config             *MarketMakerConfig
	exchange           MarketMakerExchange
	activeOrders       map[string]*MarketMakerOrder
	ordersByInstrument map[string]map[string]*MarketMakerOrder
	latestTickers      map[string]*TickerUpdate
	positions          map[string]decimal.Decimal
	stats              MarketMakerStats
	ctx                context.Context
	cancel             context.CancelFunc
	wg                 sync.WaitGroup
	mu                 sync.RWMutex
}

// placeQuotes places bid and ask orders concurrently (copied from actual implementation)
func (mm *MarketMaker) placeQuotes(instrument string, bidPrice, askPrice decimal.Decimal) error {
	var bidOrderID, askOrderID string
	var bidErr, askErr error
	ordersPlaced := 0
	
	// Use goroutines to place orders concurrently
	var wg sync.WaitGroup
	
	// Place bid order if not ask-only
	if !mm.config.AskOnly {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bidOrderID, bidErr = mm.exchange.PlaceLimitOrder(instrument, "buy", bidPrice, mm.config.QuoteSize)
		}()
	}
	
	// Place ask order if not bid-only
	if !mm.config.BidOnly {
		wg.Add(1)
		go func() {
			defer wg.Done()
			askOrderID, askErr = mm.exchange.PlaceLimitOrder(instrument, "sell", askPrice, mm.config.QuoteSize)
		}()
	}
	
	// Wait for both orders to complete
	wg.Wait()
	
	// Handle errors - if one failed, cancel the other
	if bidErr != nil && askErr != nil {
		return errors.New("failed to place both orders")
	}
	
	if bidErr != nil && askOrderID != "" {
		// Bid failed but ask succeeded, cancel ask
		mm.exchange.CancelOrder(askOrderID)
		return errors.New("failed to place bid order (cancelled ask)")
	}
	
	if askErr != nil && bidOrderID != "" {
		// Ask failed but bid succeeded, cancel bid
		mm.exchange.CancelOrder(bidOrderID)
		return errors.New("failed to place ask order (cancelled bid)")
	}
	
	// Track successfully placed orders
	mm.mu.Lock()
	if bidOrderID != "" && bidErr == nil {
		mm.trackOrder(bidOrderID, instrument, "buy", bidPrice, mm.config.QuoteSize)
		ordersPlaced++
	}
	if askOrderID != "" && askErr == nil {
		mm.trackOrder(askOrderID, instrument, "sell", askPrice, mm.config.QuoteSize)
		ordersPlaced++
	}
	mm.stats.OrdersPlaced += int64(ordersPlaced)
	mm.mu.Unlock()
	
	return nil
}

// trackOrder adds an order to tracking
func (mm *MarketMaker) trackOrder(orderID, instrument, side string, price, amount decimal.Decimal) {
	order := &MarketMakerOrder{
		OrderID:    orderID,
		Instrument: instrument,
		Side:       side,
		Price:      price,
		Amount:     amount,
		Status:     "open",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	mm.activeOrders[orderID] = order
	
	if mm.ordersByInstrument[instrument] == nil {
		mm.ordersByInstrument[instrument] = make(map[string]*MarketMakerOrder)
	}
	mm.ordersByInstrument[instrument][side] = order
}

// Simple mock exchange for testing concurrent order placement
type TestMockExchange struct {
	mu                    sync.Mutex
	placedOrders          []TestPlacedOrder
	shouldFailBid         bool
	shouldFailAsk         bool
	bidDelay              time.Duration
	askDelay              time.Duration
	cancelledOrders       map[string]bool
	orderIDCounter        int32
	bidOrderPlacedTime    time.Time
	askOrderPlacedTime    time.Time
}

type TestPlacedOrder struct {
	OrderID    string
	Instrument string
	Side       string
	Price      decimal.Decimal
	Amount     decimal.Decimal
	Timestamp  time.Time
}

func NewTestMockExchange() *TestMockExchange {
	return &TestMockExchange{
		placedOrders:    make([]TestPlacedOrder, 0),
		cancelledOrders: make(map[string]bool),
	}
}

func (m *TestMockExchange) PlaceLimitOrder(instrument, side string, price, amount decimal.Decimal) (string, error) {
	// Simulate delay
	if side == "buy" && m.bidDelay > 0 {
		time.Sleep(m.bidDelay)
	}
	if side == "sell" && m.askDelay > 0 {
		time.Sleep(m.askDelay)
	}

	// Check for failures
	if side == "buy" && m.shouldFailBid {
		return "", errors.New("bid order failed")
	}
	if side == "sell" && m.shouldFailAsk {
		return "", errors.New("ask order failed")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	counter := atomic.AddInt32(&m.orderIDCounter, 1)
	orderID := "test-order-" + side + "-" + string(rune(counter))
	
	order := TestPlacedOrder{
		OrderID:    orderID,
		Instrument: instrument,
		Side:       side,
		Price:      price,
		Amount:     amount,
		Timestamp:  time.Now(),
	}
	
	m.placedOrders = append(m.placedOrders, order)
	
	// Track timing
	if side == "buy" {
		m.bidOrderPlacedTime = order.Timestamp
	} else {
		m.askOrderPlacedTime = order.Timestamp
	}

	return orderID, nil
}

func (m *TestMockExchange) CancelOrder(orderID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cancelledOrders[orderID] = true
	return nil
}

// Test that both orders are placed concurrently and successfully
func TestMarketMakerPlaceQuotesConcurrent_Success(t *testing.T) {
	mock := NewTestMockExchange()
	
	// Create a minimal market maker config
	config := &MarketMakerConfig{
		Instruments:     []string{"BTC-PERPETUAL"},
		QuoteSize:       decimal.NewFromFloat(0.1),
		SpreadBps:       100,
		MinSpreadBps:    50,
		Improvement:     decimal.NewFromFloat(0.01),
		RefreshInterval: 5 * time.Second,
	}

	// Create market maker with mock exchange
	mm := &MarketMaker{
		config:             config,
		exchange:           mock,
		activeOrders:       make(map[string]*MarketMakerOrder),
		ordersByInstrument: make(map[string]map[string]*MarketMakerOrder),
		stats:              MarketMakerStats{},
	}
	
	bidPrice := decimal.NewFromFloat(100)
	askPrice := decimal.NewFromFloat(101)

	err := mm.placeQuotes("BTC-PERPETUAL", bidPrice, askPrice)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify both orders were placed
	if len(mock.placedOrders) != 2 {
		t.Fatalf("Expected 2 orders, got %d", len(mock.placedOrders))
	}

	// Check order details
	bidFound := false
	askFound := false
	for _, order := range mock.placedOrders {
		if order.Side == "buy" && order.Price.Equal(bidPrice) {
			bidFound = true
		}
		if order.Side == "sell" && order.Price.Equal(askPrice) {
			askFound = true
		}
	}

	if !bidFound {
		t.Error("Bid order not found")
	}
	if !askFound {
		t.Error("Ask order not found")
	}
}

// Test that when bid fails, ask is cancelled
func TestMarketMakerPlaceQuotesConcurrent_BidFailsAskCancelled(t *testing.T) {
	mock := NewTestMockExchange()
	mock.shouldFailBid = true
	
	config := &MarketMakerConfig{
		Instruments:     []string{"BTC-PERPETUAL"},
		QuoteSize:       decimal.NewFromFloat(0.1),
		SpreadBps:       100,
		MinSpreadBps:    50,
		Improvement:     decimal.NewFromFloat(0.01),
		RefreshInterval: 5 * time.Second,
	}

	mm := &MarketMaker{
		config:             config,
		exchange:           mock,
		activeOrders:       make(map[string]*MarketMakerOrder),
		ordersByInstrument: make(map[string]map[string]*MarketMakerOrder),
		stats:              MarketMakerStats{},
	}
	
	bidPrice := decimal.NewFromFloat(100)
	askPrice := decimal.NewFromFloat(101)

	err := mm.placeQuotes("BTC-PERPETUAL", bidPrice, askPrice)
	if err == nil {
		t.Fatal("Expected error when bid fails")
	}

	// Wait briefly for cancellation
	time.Sleep(50 * time.Millisecond)

	// Verify ask was placed but then cancelled
	if len(mock.placedOrders) != 1 {
		t.Fatalf("Expected 1 order (ask), got %d", len(mock.placedOrders))
	}

	if mock.placedOrders[0].Side != "sell" {
		t.Error("Expected sell order")
	}

	// Check that ask was cancelled
	askOrderID := mock.placedOrders[0].OrderID
	if !mock.cancelledOrders[askOrderID] {
		t.Error("Ask order should have been cancelled when bid failed")
	}
}

// Test concurrent execution timing
func TestMarketMakerPlaceQuotesConcurrent_ExecutionTiming(t *testing.T) {
	mock := NewTestMockExchange()
	// Set delays to simulate slow order placement
	mock.bidDelay = 100 * time.Millisecond
	mock.askDelay = 100 * time.Millisecond
	
	config := &MarketMakerConfig{
		Instruments:     []string{"BTC-PERPETUAL"},
		QuoteSize:       decimal.NewFromFloat(0.1),
		SpreadBps:       100,
		MinSpreadBps:    50,
		Improvement:     decimal.NewFromFloat(0.01),
		RefreshInterval: 5 * time.Second,
	}

	mm := &MarketMaker{
		config:             config,
		exchange:           mock,
		activeOrders:       make(map[string]*MarketMakerOrder),
		ordersByInstrument: make(map[string]map[string]*MarketMakerOrder),
		stats:              MarketMakerStats{},
	}
	
	bidPrice := decimal.NewFromFloat(100)
	askPrice := decimal.NewFromFloat(101)

	start := time.Now()
	err := mm.placeQuotes("BTC-PERPETUAL", bidPrice, askPrice)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// With concurrent placement, should take ~100ms (not 200ms for sequential)
	if elapsed > 150*time.Millisecond {
		t.Errorf("Orders took too long to place: %v (expected ~100ms for concurrent)", elapsed)
	}

	// Verify orders were placed at nearly the same time
	timeDiff := mock.askOrderPlacedTime.Sub(mock.bidOrderPlacedTime).Abs()
	if timeDiff > 20*time.Millisecond {
		t.Errorf("Orders placed too far apart: %v", timeDiff)
	}
}

// Test one-sided modes
func TestMarketMakerPlaceQuotesConcurrent_OneSided(t *testing.T) {
	testCases := []struct {
		name        string
		bidOnly     bool
		askOnly     bool
		expectSide  string
		expectCount int
	}{
		{
			name:        "BidOnly",
			bidOnly:     true,
			expectSide:  "buy",
			expectCount: 1,
		},
		{
			name:        "AskOnly",
			askOnly:     true,
			expectSide:  "sell",
			expectCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mock := NewTestMockExchange()
			
			config := &MarketMakerConfig{
				Instruments:     []string{"BTC-PERPETUAL"},
				QuoteSize:       decimal.NewFromFloat(0.1),
				SpreadBps:       100,
				MinSpreadBps:    50,
				Improvement:     decimal.NewFromFloat(0.01),
				RefreshInterval: 5 * time.Second,
				BidOnly:         tc.bidOnly,
				AskOnly:         tc.askOnly,
			}

			mm := &MarketMaker{
				config:             config,
				exchange:           mock,
				activeOrders:       make(map[string]*MarketMakerOrder),
				ordersByInstrument: make(map[string]map[string]*MarketMakerOrder),
				stats:              MarketMakerStats{},
			}
			
			bidPrice := decimal.NewFromFloat(100)
			askPrice := decimal.NewFromFloat(101)

			err := mm.placeQuotes("BTC-PERPETUAL", bidPrice, askPrice)
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			// Check only expected order was placed
			if len(mock.placedOrders) != tc.expectCount {
				t.Fatalf("Expected %d order(s), got %d", tc.expectCount, len(mock.placedOrders))
			}

			if mock.placedOrders[0].Side != tc.expectSide {
				t.Errorf("Expected %s order, got %s", tc.expectSide, mock.placedOrders[0].Side)
			}
		})
	}
}