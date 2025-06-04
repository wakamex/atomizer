package exchange

import (
	"fmt"
	"strings"
	"time"
	
	"github.com/shopspring/decimal"
	"github.com/wakamex/atomizer/internal/types"
)

// Factory creates exchange instances
type Factory struct{}

// NewFactory creates a new exchange factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateExchange creates an exchange instance that implements the general Exchange interface
func (f *Factory) CreateExchange(exchangeName string, config map[string]interface{}) (interface{}, error) {
	// Create a market maker config
	mmConfig := &types.MarketMakerConfig{
		Exchange:         exchangeName,
		ExchangeTestMode: false,
	}
	
	// Check test mode
	if testMode, ok := config["test_mode"].(bool); ok {
		mmConfig.ExchangeTestMode = testMode
	}
	
	// Create the market maker exchange
	mmExchange, err := NewExchange(mmConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create market maker exchange: %w", err)
	}
	
	// Wrap it in the adapter
	adapter := &marketMakerExchangeAdapter{
		mmExchange: mmExchange,
		name:       exchangeName,
	}
	
	return adapter, nil
}

// marketMakerExchangeAdapter adapts a MarketMakerExchange to the general Exchange interface
type marketMakerExchangeAdapter struct {
	mmExchange types.MarketMakerExchange
	name       string
}

// GetOrderBook gets the order book for RFQ pricing
func (a *marketMakerExchangeAdapter) GetOrderBook(req types.RFQResult, asset string) (types.CCXTOrderBook, error) {
	// Convert RFQ to instrument name
	instrument, err := a.ConvertToInstrument(asset, req.Strike, req.Expiry, req.IsPut)
	if err != nil {
		return types.CCXTOrderBook{}, fmt.Errorf("failed to convert to instrument: %w", err)
	}
	
	// Get order book from market maker exchange
	orderBook, err := a.mmExchange.GetOrderBook(instrument)
	if err != nil {
		return types.CCXTOrderBook{}, fmt.Errorf("failed to get order book: %w", err)
	}
	
	// Convert to CCXT format
	ccxtBook := types.CCXTOrderBook{
		Symbol: instrument,
		Bids:   make([][]float64, len(orderBook.Bids)),
		Asks:   make([][]float64, len(orderBook.Asks)),
	}
	
	for i, bid := range orderBook.Bids {
		price, _ := bid.Price.Float64()
		size, _ := bid.Size.Float64()
		ccxtBook.Bids[i] = []float64{price, size}
	}
	
	for i, ask := range orderBook.Asks {
		price, _ := ask.Price.Float64()
		size, _ := ask.Size.Float64()
		ccxtBook.Asks[i] = []float64{price, size}
	}
	
	return ccxtBook, nil
}

// PlaceOrder places an order based on RFQ confirmation
func (a *marketMakerExchangeAdapter) PlaceOrder(conf types.RFQConfirmation, instrument string, cfg interface{}) error {
	// Determine side from confirmation
	side := "buy"
	if conf.IsTakerBuy {
		side = "sell" // We sell if taker buys
	}
	
	// Parse price and quantity
	price, err := decimal.NewFromString(conf.Price)
	if err != nil {
		return fmt.Errorf("invalid price: %w", err)
	}
	
	quantity, err := decimal.NewFromString(conf.Quantity)
	if err != nil {
		return fmt.Errorf("invalid quantity: %w", err)
	}
	
	// Place order via market maker exchange
	_, err = a.mmExchange.PlaceLimitOrder(instrument, side, price, quantity)
	if err != nil {
		return fmt.Errorf("failed to place order: %w", err)
	}
	
	return nil
}

// ConvertToInstrument converts option parameters to exchange-specific instrument name
func (a *marketMakerExchangeAdapter) ConvertToInstrument(asset string, strike string, expiry int64, isPut bool) (string, error) {
	// Convert to standard format: ASSET-YYYYMMDD-STRIKE-TYPE
	expiryTime := time.Unix(expiry, 0).UTC()
	expiryStr := expiryTime.Format("20060102")
	
	optionType := "C"
	if isPut {
		optionType = "P"
	}
	
	// For Derive/Lyra format
	if a.name == "derive" {
		return fmt.Sprintf("%s-%s-%s-%s", asset, expiryStr, strike, optionType), nil
	}
	
	// For Deribit format
	if a.name == "deribit" {
		// Deribit uses format like: ETH-29DEC23-3000-C
		expiryStr = expiryTime.Format("2Jan06")
		expiryStr = strings.ToUpper(expiryStr)
		return fmt.Sprintf("%s-%s-%s-%s", asset, expiryStr, strike, optionType), nil
	}
	
	// Default format
	return fmt.Sprintf("%s-%s-%s-%s", asset, expiryStr, strike, optionType), nil
}

// GetPositions returns current positions
func (a *marketMakerExchangeAdapter) GetPositions() ([]types.ExchangePosition, error) {
	return a.mmExchange.GetPositions()
}