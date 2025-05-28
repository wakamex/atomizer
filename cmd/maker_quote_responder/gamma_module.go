package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// GammaModule wraps the gamma hedging algorithm for integration
type GammaModule struct {
	config       *AppConfig
	algo         interface{} // Will be *algo.GammaDDHAlgo when integrated
	marketData   *MarketDataProvider
	wsClient     interface{} // Will be algo.WsClient when integrated
	exchange     Exchange
	positions    map[string]*PositionInfo
	mu           sync.RWMutex
	running      bool
}

// PositionInfo tracks position details for gamma hedging
type PositionInfo struct {
	Instrument string
	Amount     decimal.Decimal
	Delta      decimal.Decimal
	Gamma      decimal.Decimal
	Expiry     int64
}

// MarketDataProvider implements the algo.MarketData interface
type MarketDataProvider struct {
	exchange  Exchange
	positions map[string]*PositionInfo
	mu        sync.RWMutex
}

// NewGammaModule creates a new gamma hedging module
func NewGammaModule(config *AppConfig) *GammaModule {
	return &GammaModule{
		config:    config,
		positions: make(map[string]*PositionInfo),
		// algo will be initialized when Start is called with proper dependencies
	}
}

// Start begins the gamma hedging loop
func (gm *GammaModule) Start(ctx context.Context) error {
	gm.mu.Lock()
	if gm.running {
		gm.mu.Unlock()
		return fmt.Errorf("gamma module already running")
	}
	gm.running = true
	gm.mu.Unlock()
	
	log.Println("Starting gamma hedging module...")
	
	// Initialize market data provider
	gm.marketData = &MarketDataProvider{
		exchange:  gm.exchange,
		positions: gm.positions,
	}
	
	// Initialize WebSocket client wrapper
	// This would wrap the existing exchange WebSocket connection
	gm.wsClient = &ExchangeWSClientWrapper{
		exchange: gm.exchange,
		config:   gm.config,
	}
	
	// Run the gamma hedging algorithm
	// TODO: Integrate with algo.StartHedger when gamma.go is properly packaged
	log.Println("Gamma hedging loop would start here")
	<-ctx.Done()
	err := ctx.Err()
	
	gm.mu.Lock()
	gm.running = false
	gm.mu.Unlock()
	
	return err
}

// OnNewPosition notifies the gamma module of a new position
func (gm *GammaModule) OnNewPosition(trade *TradeEvent) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	
	// Update or create position info
	pos, exists := gm.positions[trade.Instrument]
	if !exists {
		pos = &PositionInfo{
			Instrument: trade.Instrument,
			Amount:     decimal.Zero,
			Expiry:     trade.Expiry,
		}
		gm.positions[trade.Instrument] = pos
	}
	
	// Update position amount
	if trade.IsTakerBuy {
		// We sold
		pos.Amount = pos.Amount.Sub(trade.Quantity)
	} else {
		// We bought
		pos.Amount = pos.Amount.Add(trade.Quantity)
	}
	
	// Greeks would be fetched from pricing service
	// For now, use simplified estimates
	pos.Delta = gm.estimateDelta(trade)
	pos.Gamma = gm.estimateGamma(trade)
	
	log.Printf("Gamma module updated position %s: amount=%s, delta=%s, gamma=%s",
		trade.Instrument,
		pos.Amount.String(),
		pos.Delta.String(),
		pos.Gamma.String())
	
	// Remove if position is closed
	if pos.Amount.IsZero() {
		delete(gm.positions, trade.Instrument)
	}
}

// GetNetGreeks returns current net delta and gamma
func (gm *GammaModule) GetNetGreeks() (decimal.Decimal, decimal.Decimal) {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	
	netDelta := decimal.Zero
	netGamma := decimal.Zero
	
	for _, pos := range gm.positions {
		netDelta = netDelta.Add(pos.Delta.Mul(pos.Amount))
		netGamma = netGamma.Add(pos.Gamma.Mul(pos.Amount))
	}
	
	return netDelta, netGamma
}

// Helper methods

func (gm *GammaModule) estimateDelta(trade *TradeEvent) decimal.Decimal {
	// Simplified delta estimation
	// Real implementation would use Black-Scholes or fetch from pricing service
	
	delta := decimal.NewFromFloat(0.5) // ATM assumption
	
	if trade.IsPut {
		delta = delta.Sub(decimal.NewFromInt(1)) // Put delta = Call delta - 1
	}
	
	return delta
}

func (gm *GammaModule) estimateGamma(trade *TradeEvent) decimal.Decimal {
	// Simplified gamma estimation
	// Gamma is highest for ATM options
	return decimal.NewFromFloat(0.01)
}

// MarketDataProvider methods would implement algo.MarketData interface when integrated

func (mdp *MarketDataProvider) GetNetGreeks() (decimal.Decimal, decimal.Decimal) {
	mdp.mu.RLock()
	defer mdp.mu.RUnlock()
	
	netDelta := decimal.Zero
	netGamma := decimal.Zero
	
	for _, pos := range mdp.positions {
		netDelta = netDelta.Add(pos.Delta.Mul(pos.Amount))
		netGamma = netGamma.Add(pos.Gamma.Mul(pos.Amount))
	}
	
	return netDelta, netGamma
}

// ExchangeWSClientWrapper implements algo.WsClient interface
type ExchangeWSClientWrapper struct {
	exchange Exchange
	config   *AppConfig
}

func (w *ExchangeWSClientWrapper) Login() error {
	// Already logged in via main connection
	return nil
}

func (w *ExchangeWSClientWrapper) EnableCancelOnDisconnect() error {
	// Exchange-specific implementation
	return nil
}

func (w *ExchangeWSClientWrapper) SendOrder(instrumentName string, direction string, amount, limitPrice decimal.Decimal) error {
	// Convert and send order via exchange
	log.Printf("Gamma hedge order: %s %s @ %s", 
		direction,
		amount.String(),
		limitPrice.String())
	
	// This would use the exchange's order placement method
	return nil
}

func (w *ExchangeWSClientWrapper) SendReplace(cancelID uuid.UUID, instrumentName string, direction string, amount, limitPrice decimal.Decimal) error {
	// Cancel and replace order
	log.Printf("Gamma hedge replace order %s", cancelID.String())
	return nil
}

func (w *ExchangeWSClientWrapper) CancelAll(subaccountID int64) error {
	log.Printf("Cancelling all gamma hedge orders")
	return nil
}

func (w *ExchangeWSClientWrapper) CancelByLabel(subaccountID int64, label string) error {
	log.Printf("Cancelling gamma hedge orders with label %s", label)
	return nil
}

func (w *ExchangeWSClientWrapper) Ping() error {
	// Keep connection alive
	return nil
}