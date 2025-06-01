package main

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// MarketMakerExchange defines the interface that exchanges must implement for market making
type MarketMakerExchange interface {
	// Subscribe to real-time ticker updates for given instruments
	SubscribeTickers(ctx context.Context, instruments []string) (<-chan TickerUpdate, error)
	
	// Place a limit order
	PlaceLimitOrder(instrument string, side string, price, amount decimal.Decimal) (string, error)
	
	// Replace an existing order with new parameters (atomic cancel + create)
	ReplaceOrder(orderID string, instrument string, side string, price, amount decimal.Decimal) (string, error)
	
	// Cancel an order
	CancelOrder(orderID string) error
	
	// Get active orders
	GetOpenOrders() ([]MarketMakerOrder, error)
	
	// Get current positions
	GetPositions() ([]ExchangePosition, error)
	
	// Get order book for an instrument
	GetOrderBook(instrument string) (*MarketMakerOrderBook, error)
}

// TickerUpdate represents a real-time ticker update
type TickerUpdate struct {
	Instrument   string
	BestBid      decimal.Decimal
	BestBidSize  decimal.Decimal
	BestAsk      decimal.Decimal
	BestAskSize  decimal.Decimal
	LastPrice    decimal.Decimal
	MarkPrice    decimal.Decimal
	Timestamp    time.Time
	// Greeks for options
	Delta        *decimal.Decimal
	Gamma        *decimal.Decimal
	Vega         *decimal.Decimal
	Theta        *decimal.Decimal
	// Option details
	Expiry       *time.Time
	Strike       *decimal.Decimal
	OptionType   *string
}

// OrderBookLevel represents a price level in the order book
type OrderBookLevel struct {
	Price decimal.Decimal
	Size  decimal.Decimal
}

// MarketMakerOrderBook represents the order book for an instrument
type MarketMakerOrderBook struct {
	Bids []OrderBookLevel // Sorted by price descending (best bid first)
	Asks []OrderBookLevel // Sorted by price ascending (best ask first)
	Timestamp time.Time
}

// MarketMakerOrder represents an active order for market making
type MarketMakerOrder struct {
	OrderID        string
	Instrument     string
	Side           string // "buy" or "sell"
	Price          decimal.Decimal
	Amount         decimal.Decimal
	FilledAmount   decimal.Decimal
	Status         string // "open", "filled", "cancelled"
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// MarketMakerConfig holds configuration for the market maker
type MarketMakerConfig struct {
	// Exchange configuration
	Exchange         string   // "derive" or "deribit"
	ExchangeTestMode bool
	
	// Market making parameters
	Instruments      []string        // List of instruments to make markets on
	SpreadBps        int             // Spread in basis points (100 = 1%)
	QuoteSize        decimal.Decimal // Size of quotes
	RefreshInterval  time.Duration   // How often to update quotes
	
	// Risk parameters
	MaxPositionSize  decimal.Decimal // Maximum position per instrument
	MaxTotalExposure decimal.Decimal // Maximum total exposure across all instruments
	
	// Order management
	CancelThreshold  decimal.Decimal // Price movement threshold to trigger order updates
	MaxOrdersPerSide int             // Maximum orders per side per instrument
	
	// Performance
	MinSpreadBps     int             // Minimum spread to maintain profitability
	TargetFillRate   decimal.Decimal // Target fill rate (0-1)
	
	// Quote improvement
	Improvement      decimal.Decimal // Amount to improve quotes by (tighten spread)
	ImprovementReferenceSize decimal.Decimal // Minimum size for best bid/ask selection
	
	// One-sided quoting
	BidOnly          bool // Only place bid orders (buy side)
	AskOnly          bool // Only place ask orders (sell side)
}

// DefaultMarketMakerConfig returns a default configuration
func DefaultMarketMakerConfig() *MarketMakerConfig {
	return &MarketMakerConfig{
		Exchange:         "derive",
		ExchangeTestMode: false,
		SpreadBps:        10, // 0.1%
		QuoteSize:        decimal.NewFromFloat(1),
		RefreshInterval:  1 * time.Second,
		MaxPositionSize:  decimal.NewFromFloat(10),
		MaxTotalExposure: decimal.NewFromFloat(100),
		CancelThreshold:  decimal.NewFromFloat(0.005), // 0.5% price movement
		MaxOrdersPerSide: 1,
		MinSpreadBps:     5, // 0.05%
		TargetFillRate:   decimal.NewFromFloat(0.1), // 10% fill rate target
		Improvement:      decimal.NewFromFloat(0.1), // Default 0.1 improvement
		ImprovementReferenceSize: decimal.NewFromFloat(0), // Default 0 (use any size)
	}
}

// MarketMakerStats tracks performance statistics
type MarketMakerStats struct {
	OrdersPlaced     int64
	OrdersCancelled  int64
	OrdersFilled     int64
	TotalVolume      decimal.Decimal
	TotalPnL         decimal.Decimal
	BidAskSpread     map[string]decimal.Decimal // Current spreads by instrument
	FillRate         decimal.Decimal
	UptimeSeconds    int64
	LastUpdate       time.Time
}