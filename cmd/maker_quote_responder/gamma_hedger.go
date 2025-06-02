package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// GammaHedger implements dynamic gamma hedging for options positions
type GammaHedger struct {
	exchange         Exchange
	config           *AppConfig
	hedgeManager     *HedgeManager
	
	// Greek tracking
	positions        map[string]*OptionPosition  // instrument -> position with greeks
	netDelta         decimal.Decimal
	netGamma         decimal.Decimal
	
	// Hedging parameters
	deltaThreshold   decimal.Decimal  // Max delta before hedging
	gammaThreshold   decimal.Decimal  // Min gamma to trigger dynamic hedging
	hedgeInterval    time.Duration    // How often to check/rehedge
	minHedgeSize     decimal.Decimal  // Minimum hedge size
	
	// Hedge tracking
	currentHedge     *HedgePosition   // Current perp/spot hedge position
	
	// Control
	ctx              context.Context
	cancel           context.CancelFunc
	running          bool
	mu               sync.RWMutex
}

// OptionPosition tracks an option position with its Greeks
type OptionPosition struct {
	Instrument       string
	Amount           decimal.Decimal  // Positive for long, negative for short
	Strike           decimal.Decimal
	Expiry           int64
	IsPut            bool
	
	// Greeks (per unit)
	Delta            decimal.Decimal
	Gamma            decimal.Decimal
	Vega             decimal.Decimal
	Theta            decimal.Decimal
	
	// Market data
	UnderlyingPrice  decimal.Decimal
	ImpliedVol       decimal.Decimal
	TimeToExpiry     decimal.Decimal  // In years
	
	LastUpdate       time.Time
}

// HedgePosition tracks the hedge position
type HedgePosition struct {
	Instrument       string           // ETH-PERP or spot ETH
	Amount           decimal.Decimal  // Current hedge amount
	AvgPrice         decimal.Decimal
	LastHedgeTime    time.Time
	LastHedgeDelta   decimal.Decimal  // Delta at last hedge
}

// NewGammaHedger creates a new gamma hedger
func NewGammaHedger(exchange Exchange, config *AppConfig, hedgeManager *HedgeManager) *GammaHedger {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &GammaHedger{
		exchange:        exchange,
		config:          config,
		hedgeManager:    hedgeManager,
		positions:       make(map[string]*OptionPosition),
		deltaThreshold:  decimal.NewFromFloat(0.1),   // 0.1 ETH delta threshold
		gammaThreshold:  decimal.NewFromFloat(0.01),  // 0.01 gamma threshold
		hedgeInterval:   30 * time.Second,
		minHedgeSize:    decimal.NewFromFloat(0.1),   // Min 0.1 ETH hedge (exchange minimum)
		ctx:             ctx,
		cancel:          cancel,
	}
}

// Start begins the gamma hedging loop
func (gh *GammaHedger) Start() error {
	gh.mu.Lock()
	if gh.running {
		gh.mu.Unlock()
		return fmt.Errorf("gamma hedger already running")
	}
	gh.running = true
	gh.mu.Unlock()
	
	log.Println("Starting gamma hedger...")
	
	// Start hedging loop
	go gh.hedgingLoop()
	
	return nil
}

// Stop stops the gamma hedger
func (gh *GammaHedger) Stop() {
	log.Println("Stopping gamma hedger...")
	gh.cancel()
	
	gh.mu.Lock()
	gh.running = false
	gh.mu.Unlock()
}

// AddPosition adds or updates an option position
func (gh *GammaHedger) AddPosition(trade *TradeEvent) {
	gh.mu.Lock()
	defer gh.mu.Unlock()
	
	// Get or create position
	pos, exists := gh.positions[trade.Instrument]
	if !exists {
		pos = &OptionPosition{
			Instrument: trade.Instrument,
			Strike:     trade.Strike,
			Expiry:     trade.Expiry,
			IsPut:      trade.IsPut,
			Amount:     decimal.Zero,
		}
		gh.positions[trade.Instrument] = pos
	}
	
	// Update position amount
	if trade.IsTakerBuy {
		// We sold to taker
		pos.Amount = pos.Amount.Sub(trade.Quantity)
	} else {
		// We bought from taker
		pos.Amount = pos.Amount.Add(trade.Quantity)
	}
	
	log.Printf("Gamma hedger: Updated position %s to %s", 
		trade.Instrument, pos.Amount.String())
	
	// Remove if position is closed
	if pos.Amount.IsZero() {
		delete(gh.positions, trade.Instrument)
		log.Printf("Gamma hedger: Position %s closed", trade.Instrument)
	}
	
	// Update Greeks (would be fetched from pricing service)
	gh.updateGreeks(pos)
	
	// Recalculate portfolio Greeks
	gh.calculatePortfolioGreeks()
}

// updateGreeks updates the Greeks for a position
func (gh *GammaHedger) updateGreeks(pos *OptionPosition) {
	// In production, fetch from pricing service or calculate using Black-Scholes
	// For now, use simplified estimates
	
	// Get current market price
	underlyingPrice := gh.getCurrentUnderlyingPrice()
	pos.UnderlyingPrice = underlyingPrice
	
	// Calculate moneyness
	moneyness := underlyingPrice.Div(pos.Strike)
	
	// Time to expiry in years
	timeToExpiry := decimal.NewFromFloat(float64(pos.Expiry-time.Now().Unix()) / (365.25 * 24 * 60 * 60))
	pos.TimeToExpiry = timeToExpiry
	
	// Simplified Greeks estimation
	if pos.IsPut {
		// Put option Greeks
		if moneyness.GreaterThan(decimal.NewFromFloat(1.1)) {
			// OTM put
			pos.Delta = decimal.NewFromFloat(-0.2)
			pos.Gamma = decimal.NewFromFloat(0.005)
		} else if moneyness.LessThan(decimal.NewFromFloat(0.9)) {
			// ITM put
			pos.Delta = decimal.NewFromFloat(-0.8)
			pos.Gamma = decimal.NewFromFloat(0.005)
		} else {
			// ATM put
			pos.Delta = decimal.NewFromFloat(-0.5)
			pos.Gamma = decimal.NewFromFloat(0.02)
		}
	} else {
		// Call option Greeks
		if moneyness.LessThan(decimal.NewFromFloat(0.9)) {
			// OTM call
			pos.Delta = decimal.NewFromFloat(0.2)
			pos.Gamma = decimal.NewFromFloat(0.005)
		} else if moneyness.GreaterThan(decimal.NewFromFloat(1.1)) {
			// ITM call
			pos.Delta = decimal.NewFromFloat(0.8)
			pos.Gamma = decimal.NewFromFloat(0.005)
		} else {
			// ATM call
			pos.Delta = decimal.NewFromFloat(0.5)
			pos.Gamma = decimal.NewFromFloat(0.02)
		}
	}
	
	// Adjust gamma for time decay (gamma increases as expiry approaches for ATM)
	if timeToExpiry.LessThan(decimal.NewFromFloat(0.1)) { // Less than ~36 days
		pos.Gamma = pos.Gamma.Mul(decimal.NewFromFloat(1.5))
	}
	
	pos.LastUpdate = time.Now()
}

// calculatePortfolioGreeks calculates net portfolio Greeks
func (gh *GammaHedger) calculatePortfolioGreeks() {
	netDelta := decimal.Zero
	netGamma := decimal.Zero
	
	for _, pos := range gh.positions {
		// Position Greeks = position size * per-unit Greeks
		netDelta = netDelta.Add(pos.Amount.Mul(pos.Delta))
		netGamma = netGamma.Add(pos.Amount.Mul(pos.Gamma))
	}
	
	// Add hedge position delta (perp/spot is delta 1)
	if gh.currentHedge != nil {
		netDelta = netDelta.Add(gh.currentHedge.Amount)
	}
	
	gh.netDelta = netDelta
	gh.netGamma = netGamma
	
	log.Printf("Portfolio Greeks - Delta: %s, Gamma: %s", 
		netDelta.StringFixed(4), netGamma.StringFixed(4))
}

// hedgingLoop is the main hedging loop
func (gh *GammaHedger) hedgingLoop() {
	ticker := time.NewTicker(gh.hedgeInterval)
	defer ticker.Stop()
	
	// Initial hedge on startup
	gh.checkAndHedge()
	
	for {
		select {
		case <-gh.ctx.Done():
			return
		case <-ticker.C:
			gh.checkAndHedge()
		}
	}
}

// checkAndHedge checks if hedging is needed and executes
func (gh *GammaHedger) checkAndHedge() {
	gh.mu.RLock()
	netDelta := gh.netDelta
	netGamma := gh.netGamma
	positions := len(gh.positions)
	gh.mu.RUnlock()
	
	if positions == 0 {
		// No positions to hedge
		return
	}
	
	// Update Greeks for all positions
	gh.mu.Lock()
	for _, pos := range gh.positions {
		gh.updateGreeks(pos)
	}
	gh.calculatePortfolioGreeks()
	
	// Get updated values after recalculation
	netDelta = gh.netDelta
	netGamma = gh.netGamma
	gh.mu.Unlock()
	
	// Check if we need to hedge
	shouldHedge := false
	hedgeReason := ""
	
	// 1. Delta threshold exceeded
	if netDelta.Abs().GreaterThan(gh.deltaThreshold) {
		shouldHedge = true
		hedgeReason = fmt.Sprintf("Delta threshold exceeded: %s > %s", 
			netDelta.StringFixed(4), gh.deltaThreshold.StringFixed(4))
	}
	
	// 2. Gamma-based dynamic hedging
	if netGamma.Abs().GreaterThan(gh.gammaThreshold) && gh.currentHedge != nil {
		// Calculate expected delta change for a 1% move
		underlyingPrice := gh.getCurrentUnderlyingPrice()
		priceMove := underlyingPrice.Mul(decimal.NewFromFloat(0.01))
		expectedDeltaChange := netGamma.Mul(priceMove)
		
		// If expected delta change is significant, hedge proactively
		if expectedDeltaChange.Abs().GreaterThan(gh.deltaThreshold.Mul(decimal.NewFromFloat(0.5))) {
			shouldHedge = true
			hedgeReason = fmt.Sprintf("High gamma risk: %s gamma, %s expected delta change on 1%% move", 
				netGamma.StringFixed(4), expectedDeltaChange.StringFixed(4))
		}
	}
	
	// 3. Time-based rehedge for high gamma positions
	if gh.currentHedge != nil && time.Since(gh.currentHedge.LastHedgeTime) > 5*time.Minute {
		if netGamma.Abs().GreaterThan(gh.gammaThreshold) {
			// Delta might have drifted due to underlying price movement
			deltaChange := netDelta.Sub(gh.currentHedge.LastHedgeDelta).Abs()
			if deltaChange.GreaterThan(gh.deltaThreshold.Mul(decimal.NewFromFloat(0.3))) {
				shouldHedge = true
				hedgeReason = fmt.Sprintf("Time-based rehedge: delta drifted by %s", 
					deltaChange.StringFixed(4))
			}
		}
	}
	
	if !shouldHedge {
		return
	}
	
	log.Printf("Gamma hedger: %s", hedgeReason)
	
	// Calculate hedge amount
	targetHedgeAmount := netDelta.Neg() // Hedge to neutralize delta
	currentHedgeAmount := decimal.Zero
	if gh.currentHedge != nil {
		currentHedgeAmount = gh.currentHedge.Amount
	}
	
	hedgeAdjustment := targetHedgeAmount.Sub(currentHedgeAmount)
	
	// Check minimum hedge size
	if hedgeAdjustment.Abs().LessThan(gh.minHedgeSize) {
		log.Printf("Hedge adjustment %s below minimum %s, skipping", 
			hedgeAdjustment.StringFixed(4), gh.minHedgeSize.StringFixed(4))
		return
	}
	
	// Execute hedge
	gh.executeHedge(hedgeAdjustment, netDelta)
}

// executeHedge executes the hedge trade
func (gh *GammaHedger) executeHedge(amount decimal.Decimal, currentDelta decimal.Decimal) {
	log.Printf("Executing hedge: %s ETH (current delta: %s)", 
		amount.StringFixed(4), currentDelta.StringFixed(4))
	
	// Create a synthetic trade event for the hedge
	hedgeTrade := &TradeEvent{
		ID:         fmt.Sprintf("gamma_hedge_%d", time.Now().UnixNano()),
		Source:     TradeSourceHedge,
		Instrument: "ETH-PERP", // Use perpetual for hedging
		Quantity:   amount.Abs(),
		IsTakerBuy: amount.GreaterThan(decimal.Zero), // Buy if positive hedge needed
		Timestamp:  time.Now(),
	}
	
	// Execute via hedge manager
	result, err := gh.hedgeManager.ExecuteHedge(hedgeTrade)
	if err != nil {
		log.Printf("Gamma hedge failed: %v", err)
		return
	}
	
	// Update hedge tracking
	gh.mu.Lock()
	if gh.currentHedge == nil {
		gh.currentHedge = &HedgePosition{
			Instrument: result.Instrument,
		}
	}
	
	gh.currentHedge.Amount = gh.currentHedge.Amount.Add(amount)
	gh.currentHedge.LastHedgeTime = time.Now()
	gh.currentHedge.LastHedgeDelta = currentDelta
	
	// Recalculate portfolio Greeks after hedge
	gh.calculatePortfolioGreeks()
	gh.mu.Unlock()
	
	log.Printf("Gamma hedge executed: %s @ %s, new portfolio delta: %s", 
		result.Quantity.StringFixed(4), 
		result.Price.StringFixed(2),
		gh.netDelta.StringFixed(4))
}

// getCurrentUnderlyingPrice gets current ETH price
func (gh *GammaHedger) getCurrentUnderlyingPrice() decimal.Decimal {
	// In production, get from market data feed
	// For now, return a default
	return decimal.NewFromFloat(3000)
}

// GetStatus returns current hedging status
func (gh *GammaHedger) GetStatus() map[string]interface{} {
	gh.mu.RLock()
	defer gh.mu.RUnlock()
	
	status := map[string]interface{}{
		"running":        gh.running,
		"net_delta":      gh.netDelta.StringFixed(4),
		"net_gamma":      gh.netGamma.StringFixed(4),
		"positions":      len(gh.positions),
		"delta_threshold": gh.deltaThreshold.StringFixed(4),
		"gamma_threshold": gh.gammaThreshold.StringFixed(4),
	}
	
	if gh.currentHedge != nil {
		status["hedge_amount"] = gh.currentHedge.Amount.StringFixed(4)
		status["last_hedge"] = gh.currentHedge.LastHedgeTime.Format(time.RFC3339)
	}
	
	return status
}

// SetThresholds updates hedging thresholds
func (gh *GammaHedger) SetThresholds(deltaThreshold, gammaThreshold float64) {
	gh.mu.Lock()
	defer gh.mu.Unlock()
	
	gh.deltaThreshold = decimal.NewFromFloat(deltaThreshold)
	gh.gammaThreshold = decimal.NewFromFloat(gammaThreshold)
	
	log.Printf("Gamma hedger thresholds updated - Delta: %s, Gamma: %s",
		gh.deltaThreshold.StringFixed(4),
		gh.gammaThreshold.StringFixed(4))
}