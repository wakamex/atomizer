package types

import (
	"time"

	"github.com/shopspring/decimal"
)

// ExchangePosition represents an open position on the exchange
type ExchangePosition struct {
	InstrumentName string
	Amount         float64
	Direction      string // "buy" or "sell"
	AveragePrice   float64
	MarkPrice      float64
	IndexPrice     float64
	PnL            float64
}

// Position represents a current position with Greeks
// This is used internally by the market maker for risk management
type Position struct {
	Instrument  string
	Quantity    decimal.Decimal
	AvgPrice    decimal.Decimal
	Delta       decimal.Decimal
	Gamma       decimal.Decimal
	LastUpdated time.Time
}

// RiskMetrics contains current risk measurements
type RiskMetrics struct {
	TotalDelta      decimal.Decimal
	TotalGamma      decimal.Decimal
	TotalPositions  int
	MaxPositionSize decimal.Decimal
	UpdatedAt       time.Time
}

// Exchange defines the general exchange interface for arbitrage
type Exchange interface {
	// Get order book for RFQ pricing
	GetOrderBook(req RFQResult, asset string) (CCXTOrderBook, error)
	
	// Place an order based on RFQ confirmation
	PlaceOrder(conf RFQConfirmation, instrument string, cfg interface{}) error
	
	// Convert option parameters to exchange-specific instrument name
	ConvertToInstrument(asset string, strike string, expiry int64, isPut bool) (string, error)
	
	// Get current positions
	GetPositions() ([]ExchangePosition, error)
}

// CCXTOrderBook represents an order book in CCXT format
type CCXTOrderBook struct {
	Symbol string      `json:"symbol"`
	Bids   [][]float64 `json:"bids"`
	Asks   [][]float64 `json:"asks"`
	Index  float64     `json:"index"` // Index price of the underlying asset
}
