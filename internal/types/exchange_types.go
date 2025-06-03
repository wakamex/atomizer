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
