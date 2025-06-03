package main

import (
    "context"
    "time"
    "github.com/shopspring/decimal"
)

// MarketMakerExchange interface - your exchange abstraction
type MarketMakerExchange interface {
    PlaceLimitOrder(instrument, side string, price, amount decimal.Decimal) (string, error)
    CancelOrder(orderID string) error
    ReplaceOrder(orderID, instrument, side string, price, amount decimal.Decimal) (string, error)
    GetOpenOrders() ([]MarketMakerOrder, error)
    GetOrderBook(instrument string) (*MarketMakerOrderBook, error)
    GetPositions() ([]Position, error)
    SubscribeTickers(ctx context.Context, instruments []string) (<-chan TickerUpdate, error)
}

// MarketMakerOrder represents an order
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

// MarketMakerOrderBook represents order book data
type MarketMakerOrderBook struct {
    Bids []OrderBookLevel
    Asks []OrderBookLevel
}

type OrderBookLevel struct {
    Price decimal.Decimal
    Size  decimal.Decimal
}

// TickerUpdate represents market data update
type TickerUpdate struct {
    Instrument string
    BestBid    decimal.Decimal
    BestAsk    decimal.Decimal
    MarkPrice  decimal.Decimal
    Timestamp  time.Time
}

// Position represents a position
type Position struct {
    InstrumentName string
    Direction      string
    Amount         float64
}

// MarketMakerConfig holds configuration
type MarketMakerConfig struct {
    Instruments              []string
    QuoteSize               decimal.Decimal
    SpreadBps               int
    MinSpreadBps            int
    Improvement             decimal.Decimal
    ImprovementReferenceSize decimal.Decimal
    RefreshInterval         time.Duration
    CancelThreshold         decimal.Decimal
    MaxPositionSize         decimal.Decimal
    MaxTotalExposure        decimal.Decimal
    BidOnly                 bool
    AskOnly                 bool
}

// MarketMakerStats tracks statistics
type MarketMakerStats struct {
    OrdersPlaced    int64
    OrdersCancelled int64
    OrdersFilled    int64
    BidAskSpread    map[string]decimal.Decimal
    UptimeSeconds   int64
    LastUpdate      time.Time
}