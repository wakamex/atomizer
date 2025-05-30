package exchange

import (
	"github.com/shopspring/decimal"
	"github.com/wakamex/atomizer/internal/types"
)

// Exchange defines the interface for interacting with different cryptocurrency exchanges
type Exchange interface {
	// GetOrderBook fetches the order book for a given option
	GetOrderBook(req types.RFQResult, asset string) (OrderBook, error)
	
	// PlaceOrder places an order on the exchange
	PlaceOrder(confirmation types.RFQConfirmation, instrument string) error
	
	// ConvertToInstrument converts option details to exchange-specific instrument format
	ConvertToInstrument(asset string, strike string, expiry int64, isPut bool) (string, error)
	
	// GetPositions fetches all open positions from the exchange
	GetPositions() ([]Position, error)
	
	// GetOpenOrders fetches all open orders
	GetOpenOrders() ([]Order, error)
	
	// CancelOrder cancels an order
	CancelOrder(orderID string) error
}

// Position represents an open position on the exchange
type Position struct {
	InstrumentName string
	Amount         float64
	Direction      string  // "buy" or "sell"
	AveragePrice   float64
	MarkPrice      float64
	IndexPrice     float64
	PnL            float64
}

// Order represents an order on the exchange
type Order struct {
	OrderID      string
	Instrument   string
	Side         string // "buy" or "sell"
	OrderType    string // "limit", "market"
	Price        decimal.Decimal
	Amount       decimal.Decimal
	FilledAmount decimal.Decimal
	Status       string // "open", "filled", "cancelled"
	CreatedAt    int64
	UpdatedAt    int64
}

// OrderBook represents the order book for an instrument
type OrderBook struct {
	InstrumentName string
	Bids           []PriceLevel
	Asks           []PriceLevel
	Timestamp      int64
}

// PriceLevel represents a price and quantity at that level
type PriceLevel struct {
	Price    decimal.Decimal
	Quantity decimal.Decimal
}

// Config holds common configuration for exchanges
type Config struct {
	APIKey    string
	APISecret string
	TestMode  bool
	RateLimit int
}