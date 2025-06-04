package types

import (
	"context"
	
	"github.com/shopspring/decimal"
)

// UnifiedExchange combines both Exchange and MarketMakerExchange capabilities
type UnifiedExchange interface {
	// RFQ/Arbitrage methods
	GetOrderBookForRFQ(req RFQResult, asset string) (CCXTOrderBook, error)
	PlaceOrder(conf RFQConfirmation, instrument string, cfg interface{}) error
	ConvertToInstrument(asset string, strike string, expiry int64, isPut bool) (string, error)
	
	// Market Maker methods
	GetOrderBook(instrument string) (*MarketMakerOrderBook, error)
	PlaceLimitOrder(instrument string, side string, price, amount decimal.Decimal) (string, error)
	ReplaceOrder(orderID string, instrument string, side string, price, amount decimal.Decimal) (string, error)
	CancelOrder(orderID string) error
	GetOpenOrders() ([]MarketMakerOrder, error)
	
	// Common methods
	GetPositions() ([]ExchangePosition, error)
	SubscribeTickers(ctx context.Context, instruments []string) (<-chan TickerUpdate, error)
}