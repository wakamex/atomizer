package main

import ()

// Exchange defines the interface for interacting with different cryptocurrency exchanges
type Exchange interface {
	// GetOrderBook fetches the order book for a given option
	GetOrderBook(req RFQResult, asset string) (CCXTOrderBook, error)
	
	// PlaceHedgeOrder places a hedge order on the exchange
	PlaceHedgeOrder(conf RFQConfirmation, underlying string, cfg *AppConfig) error
	
	// ConvertToInstrument converts option details to exchange-specific instrument format
	ConvertToInstrument(asset string, strike string, expiry int64, isPut bool) (string, error)
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