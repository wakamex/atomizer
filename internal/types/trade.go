package types

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// TradeSourceType represents the source of a trade
type TradeSourceType string

const (
	TradeSourceRysk   TradeSourceType = "RYSK_RFQ"
	TradeSourceManual TradeSourceType = "MANUAL"
	TradeSourceHedge  TradeSourceType = "HEDGE"
)

// TradeStatus represents the current status of a trade
type TradeStatus string

const (
	TradeStatusPending   TradeStatus = "PENDING"
	TradeStatusQuoted    TradeStatus = "QUOTED"
	TradeStatusExecuted  TradeStatus = "EXECUTED"
	TradeStatusHedged    TradeStatus = "HEDGED"
	TradeStatusFailed    TradeStatus = "FAILED"
	TradeStatusCancelled TradeStatus = "CANCELLED"
)

// TradeEvent represents a unified trade across all sources
type TradeEvent struct {
	ID              string
	Source          TradeSourceType
	Status          TradeStatus
	RFQId           string // Original RFQ ID if from Rysk
	Instrument      string
	Strike          decimal.Decimal
	Expiry          int64
	IsPut           bool
	Quantity        decimal.Decimal
	Price           decimal.Decimal
	IsTakerBuy      bool
	Timestamp       time.Time
	HedgeOrderID    string
	HedgeExchange   string
	Error           error
}

// ManualTradeRequest represents a request for manual trade execution
type ManualTradeRequest struct {
	Asset      string          `json:"asset"`
	Strike     string          `json:"strike"`
	Expiry     int64           `json:"expiry"`
	IsPut      bool            `json:"isPut"`
	Quantity   decimal.Decimal `json:"quantity"`
	IsTakerBuy bool            `json:"isTakerBuy"`
}

// HedgeManager defines the interface for hedge execution
type HedgeManager interface {
	ExecuteHedge(ctx context.Context, trade *TradeEvent) error
}

// RiskManager defines the interface for risk management
type RiskManager interface {
	ValidateTrade(trade *TradeEvent) error
	UpdatePosition(trade *TradeEvent)
	GetGreeks() (delta, gamma decimal.Decimal)
	GetPositions() map[string]Position
}

