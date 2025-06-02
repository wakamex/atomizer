package main

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