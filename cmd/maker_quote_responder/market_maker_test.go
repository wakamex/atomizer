package main

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// Initialize debug mode for tests
func init() {
	debugMode = false
}

// Mock debugLog for tests if not defined
func debugLogTest(format string, args ...interface{}) {
	if debugMode {
		log.Printf(format, args...)
	}
}

// MockExchange for testing
type MockExchange struct {
	mu                    sync.Mutex
	placedOrders          []PlacedOrder
	shouldFailBid         bool
	shouldFailAsk         bool
	bidDelay              time.Duration
	askDelay              time.Duration
	cancelledOrders       map[string]bool
	openOrders            []MarketMakerOrder
	orderIDCounter        int
	shouldFailCancel      bool
	cancelFailureOrderIDs map[string]bool
}

type PlacedOrder struct {
	OrderID    string
	Instrument string
	Side       string
	Price      decimal.Decimal
	Amount     decimal.Decimal
	Timestamp  time.Time
}

func NewMockExchange() *MockExchange {
	return &MockExchange{
		placedOrders:          make([]PlacedOrder, 0),
		cancelledOrders:       make(map[string]bool),
		openOrders:            make([]MarketMakerOrder, 0),
		cancelFailureOrderIDs: make(map[string]bool),
	}
}

func (m *MockExchange) PlaceLimitOrder(instrument, side string, price, amount decimal.Decimal) (string, error) {
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

	m.orderIDCounter++
	orderID := "order-" + side + "-" + string(rune(m.orderIDCounter))
	
	m.placedOrders = append(m.placedOrders, PlacedOrder{
		OrderID:    orderID,
		Instrument: instrument,
		Side:       side,
		Price:      price,
		Amount:     amount,
		Timestamp:  time.Now(),
	})

	return orderID, nil
}

func (m *MockExchange) CancelOrder(orderID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.shouldFailCancel || m.cancelFailureOrderIDs[orderID] {
		return errors.New("cancel failed")
	}

	m.cancelledOrders[orderID] = true
	return nil
}

func (m *MockExchange) GetOpenOrders() ([]MarketMakerOrder, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.openOrders, nil
}

func (m *MockExchange) GetOrderBook(instrument string) (*MarketMakerOrderBook, error) {
	return &MarketMakerOrderBook{
		Bids: []OrderBookLevel{
			{Price: decimal.NewFromFloat(100), Size: decimal.NewFromFloat(10)},
		},
		Asks: []OrderBookLevel{
			{Price: decimal.NewFromFloat(101), Size: decimal.NewFromFloat(10)},
		},
	}, nil
}

func (m *MockExchange) GetPositions() ([]Position, error) {
	return []Position{}, nil
}

func (m *MockExchange) SubscribeTickers(ctx context.Context, instruments []string) (<-chan TickerUpdate, error) {
	ch := make(chan TickerUpdate)
	go func() {
		<-ctx.Done()
		close(ch)
	}()
	return ch, nil
}

func (m *MockExchange) ReplaceOrder(orderID, instrument, side string, price, amount decimal.Decimal) (string, error) {
	return "", errors.New("not implemented")
}

// Test concurrent order placement - both orders succeed
func TestPlaceQuotesConcurrent_BothSucceed(t *testing.T) {
	mock := NewMockExchange()
	config := &MarketMakerConfig{
		Instruments:     []string{"BTC-PERPETUAL"},
		QuoteSize:       decimal.NewFromFloat(0.1),
		SpreadBps:       100,
		MinSpreadBps:    50,
		Improvement:     decimal.NewFromFloat(0.01),
		RefreshInterval: 5 * time.Second,
	}

	mm := NewMarketMaker(config, mock)
	
	bidPrice := decimal.NewFromFloat(100)
	askPrice := decimal.NewFromFloat(101)

	err := mm.placeQuotes("BTC-PERPETUAL", bidPrice, askPrice)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check that both orders were placed
	mock.mu.Lock()
	defer mock.mu.Unlock()

	if len(mock.placedOrders) != 2 {
		t.Fatalf("Expected 2 orders, got %d", len(mock.placedOrders))
	}

	// Verify bid order
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

	// Check internal tracking
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if len(mm.activeOrders) != 2 {
		t.Errorf("Expected 2 tracked orders, got %d", len(mm.activeOrders))
	}
}

// Test concurrent order placement - bid fails, ask should be cancelled
func TestPlaceQuotesConcurrent_BidFailsAskCancelled(t *testing.T) {
	mock := NewMockExchange()
	mock.shouldFailBid = true
	
	config := &MarketMakerConfig{
		Instruments:     []string{"BTC-PERPETUAL"},
		QuoteSize:       decimal.NewFromFloat(0.1),
		SpreadBps:       100,
		MinSpreadBps:    50,
		Improvement:     decimal.NewFromFloat(0.01),
		RefreshInterval: 5 * time.Second,
	}

	mm := NewMarketMaker(config, mock)
	
	bidPrice := decimal.NewFromFloat(100)
	askPrice := decimal.NewFromFloat(101)

	err := mm.placeQuotes("BTC-PERPETUAL", bidPrice, askPrice)
	if err == nil {
		t.Fatal("Expected error when bid fails")
	}

	// Wait a bit for cancellation
	time.Sleep(100 * time.Millisecond)

	mock.mu.Lock()
	defer mock.mu.Unlock()

	// Check that ask was placed (1 order)
	if len(mock.placedOrders) != 1 {
		t.Fatalf("Expected 1 order (ask), got %d", len(mock.placedOrders))
	}

	// Verify it was the ask order
	if mock.placedOrders[0].Side != "sell" {
		t.Error("Expected sell order")
	}

	// Check that ask was cancelled
	askOrderID := mock.placedOrders[0].OrderID
	if !mock.cancelledOrders[askOrderID] {
		t.Error("Ask order should have been cancelled when bid failed")
	}

	// Check internal tracking is empty
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if len(mm.activeOrders) != 0 {
		t.Errorf("Expected 0 tracked orders after failure, got %d", len(mm.activeOrders))
	}
}

// Test concurrent order placement - ask fails, bid should be cancelled
func TestPlaceQuotesConcurrent_AskFailsBidCancelled(t *testing.T) {
	mock := NewMockExchange()
	mock.shouldFailAsk = true
	
	config := &MarketMakerConfig{
		Instruments:     []string{"BTC-PERPETUAL"},
		QuoteSize:       decimal.NewFromFloat(0.1),
		SpreadBps:       100,
		MinSpreadBps:    50,
		Improvement:     decimal.NewFromFloat(0.01),
		RefreshInterval: 5 * time.Second,
	}

	mm := NewMarketMaker(config, mock)
	
	bidPrice := decimal.NewFromFloat(100)
	askPrice := decimal.NewFromFloat(101)

	err := mm.placeQuotes("BTC-PERPETUAL", bidPrice, askPrice)
	if err == nil {
		t.Fatal("Expected error when ask fails")
	}

	// Wait a bit for cancellation
	time.Sleep(100 * time.Millisecond)

	mock.mu.Lock()
	defer mock.mu.Unlock()

	// Check that bid was placed (1 order)
	if len(mock.placedOrders) != 1 {
		t.Fatalf("Expected 1 order (bid), got %d", len(mock.placedOrders))
	}

	// Verify it was the bid order
	if mock.placedOrders[0].Side != "buy" {
		t.Error("Expected buy order")
	}

	// Check that bid was cancelled
	bidOrderID := mock.placedOrders[0].OrderID
	if !mock.cancelledOrders[bidOrderID] {
		t.Error("Bid order should have been cancelled when ask failed")
	}

	// Check internal tracking is empty
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if len(mm.activeOrders) != 0 {
		t.Errorf("Expected 0 tracked orders after failure, got %d", len(mm.activeOrders))
	}
}

// Test concurrent order placement - both fail
func TestPlaceQuotesConcurrent_BothFail(t *testing.T) {
	mock := NewMockExchange()
	mock.shouldFailBid = true
	mock.shouldFailAsk = true
	
	config := &MarketMakerConfig{
		Instruments:     []string{"BTC-PERPETUAL"},
		QuoteSize:       decimal.NewFromFloat(0.1),
		SpreadBps:       100,
		MinSpreadBps:    50,
		Improvement:     decimal.NewFromFloat(0.01),
		RefreshInterval: 5 * time.Second,
	}

	mm := NewMarketMaker(config, mock)
	
	bidPrice := decimal.NewFromFloat(100)
	askPrice := decimal.NewFromFloat(101)

	err := mm.placeQuotes("BTC-PERPETUAL", bidPrice, askPrice)
	if err == nil {
		t.Fatal("Expected error when both orders fail")
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()

	// Check that no orders were placed
	if len(mock.placedOrders) != 0 {
		t.Fatalf("Expected 0 orders, got %d", len(mock.placedOrders))
	}

	// Check no cancellations
	if len(mock.cancelledOrders) != 0 {
		t.Error("No orders should have been cancelled")
	}

	// Check internal tracking is empty
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if len(mm.activeOrders) != 0 {
		t.Errorf("Expected 0 tracked orders, got %d", len(mm.activeOrders))
	}
}

// Test one-sided mode - bid only
func TestPlaceQuotesConcurrent_BidOnly(t *testing.T) {
	mock := NewMockExchange()
	
	config := &MarketMakerConfig{
		Instruments:     []string{"BTC-PERPETUAL"},
		QuoteSize:       decimal.NewFromFloat(0.1),
		SpreadBps:       100,
		MinSpreadBps:    50,
		Improvement:     decimal.NewFromFloat(0.01),
		RefreshInterval: 5 * time.Second,
		BidOnly:         true,
	}

	mm := NewMarketMaker(config, mock)
	
	bidPrice := decimal.NewFromFloat(100)
	askPrice := decimal.NewFromFloat(101)

	err := mm.placeQuotes("BTC-PERPETUAL", bidPrice, askPrice)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()

	// Check that only bid was placed
	if len(mock.placedOrders) != 1 {
		t.Fatalf("Expected 1 order (bid), got %d", len(mock.placedOrders))
	}

	if mock.placedOrders[0].Side != "buy" {
		t.Error("Expected buy order in bid-only mode")
	}

	// Check internal tracking
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if len(mm.activeOrders) != 1 {
		t.Errorf("Expected 1 tracked order, got %d", len(mm.activeOrders))
	}
}

// Test one-sided mode - ask only
func TestPlaceQuotesConcurrent_AskOnly(t *testing.T) {
	mock := NewMockExchange()
	
	config := &MarketMakerConfig{
		Instruments:     []string{"BTC-PERPETUAL"},
		QuoteSize:       decimal.NewFromFloat(0.1),
		SpreadBps:       100,
		MinSpreadBps:    50,
		Improvement:     decimal.NewFromFloat(0.01),
		RefreshInterval: 5 * time.Second,
		AskOnly:         true,
	}

	mm := NewMarketMaker(config, mock)
	
	bidPrice := decimal.NewFromFloat(100)
	askPrice := decimal.NewFromFloat(101)

	err := mm.placeQuotes("BTC-PERPETUAL", bidPrice, askPrice)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()

	// Check that only ask was placed
	if len(mock.placedOrders) != 1 {
		t.Fatalf("Expected 1 order (ask), got %d", len(mock.placedOrders))
	}

	if mock.placedOrders[0].Side != "sell" {
		t.Error("Expected sell order in ask-only mode")
	}

	// Check internal tracking
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	if len(mm.activeOrders) != 1 {
		t.Errorf("Expected 1 tracked order, got %d", len(mm.activeOrders))
	}
}

// Test that orders are placed concurrently (timing test)
func TestPlaceQuotesConcurrent_Timing(t *testing.T) {
	mock := NewMockExchange()
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

	mm := NewMarketMaker(config, mock)
	
	bidPrice := decimal.NewFromFloat(100)
	askPrice := decimal.NewFromFloat(101)

	start := time.Now()
	err := mm.placeQuotes("BTC-PERPETUAL", bidPrice, askPrice)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// If orders were placed sequentially, it would take ~200ms
	// With concurrent placement, it should take ~100ms
	if elapsed > 150*time.Millisecond {
		t.Errorf("Orders took too long to place: %v (expected ~100ms for concurrent placement)", elapsed)
	}

	mock.mu.Lock()
	defer mock.mu.Unlock()

	if len(mock.placedOrders) != 2 {
		t.Fatalf("Expected 2 orders, got %d", len(mock.placedOrders))
	}

	// Check that orders were placed close together in time
	timeDiff := mock.placedOrders[1].Timestamp.Sub(mock.placedOrders[0].Timestamp)
	if timeDiff > 20*time.Millisecond {
		t.Errorf("Orders placed too far apart: %v", timeDiff)
	}
}