package main

import ()

// Exchange defines the interface for interacting with different cryptocurrency exchanges
type Exchange interface {
	// GetOrderBook fetches the order book for a given option
	GetOrderBook(req RFQResult, asset string) (CCXTOrderBook, error)
	
	// PlaceOrder places an order on the exchange
	PlaceOrder(conf RFQConfirmation, instrument string, cfg *AppConfig) error
	
	// ConvertToInstrument converts option details to exchange-specific instrument format
	ConvertToInstrument(asset string, strike string, expiry int64, isPut bool) (string, error)
	
	// GetPositions fetches all open positions from the exchange
	GetPositions() ([]ExchangePosition, error)
}

// ExchangePosition represents an open position on the exchange
type ExchangePosition struct {
	InstrumentName string
	Amount         float64
	Direction      string  // "buy" or "sell"
	AveragePrice   float64
	MarkPrice      float64
	IndexPrice     float64
	PnL            float64
}

// ExchangeConfig holds common configuration for exchanges
type ExchangeConfig struct {
	APIKey    string
	APISecret string
	TestMode  bool
	RateLimit int
}

// OrderResult represents the result of placing an order
type OrderResult struct {
	OrderID     string
	FilledPrice float64
	FilledQty   float64
	Status      string
}

// PriceLevel represents a price and quantity at that level
type PriceLevel struct {
	Price    float64
	Quantity float64
}

// OptionQuote represents a quote for an option
type OptionQuote struct {
	Bid           float64
	Ask           float64
	BidSize       float64
	AskSize       float64
	UnderlyingPrice float64
	IV            float64
}