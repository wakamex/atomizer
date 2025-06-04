package gamma

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wakamex/atomizer/internal/config"
	"github.com/wakamex/atomizer/internal/types"
	"github.com/shopspring/decimal"
)

// Hedger implements dynamic gamma hedging for options positions
type Hedger struct {
	exchange         types.Exchange
	config           *config.Config
	hedgeManager     types.HedgeManager
	
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

// NewHedger creates a new gamma hedger
func NewHedger(exchange types.Exchange, cfg *config.Config, hedgeManager types.HedgeManager) *Hedger {
	ctx, cancel := context.WithCancel(context.Background())
	
	gammaThreshold := decimal.NewFromFloat(0.1)
	if cfg.GammaThreshold > 0 {
		gammaThreshold = decimal.NewFromFloat(cfg.GammaThreshold)
	}
	
	return &Hedger{
		exchange:        exchange,
		config:          cfg,
		hedgeManager:    hedgeManager,
		positions:       make(map[string]*OptionPosition),
		deltaThreshold:  decimal.NewFromFloat(0.1),   // 0.1 ETH delta threshold
		gammaThreshold:  gammaThreshold,
		hedgeInterval:   30 * time.Second,
		minHedgeSize:    decimal.NewFromFloat(0.01),  // Min 0.01 ETH hedge
		ctx:             ctx,
		cancel:          cancel,
	}
}

// Start begins the gamma hedging loop
func (h *Hedger) Start(ctx context.Context) {
	h.mu.Lock()
	if h.running {
		h.mu.Unlock()
		return
	}
	h.running = true
	h.mu.Unlock()
	
	log.Println("Starting gamma hedger...")
	
	// Main hedging loop
	ticker := time.NewTicker(h.hedgeInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			h.Stop()
			return
		case <-ticker.C:
			h.performHedgeCheck()
		}
	}
}

// Stop halts the gamma hedger
func (h *Hedger) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if !h.running {
		return
	}
	
	log.Println("Stopping gamma hedger...")
	h.cancel()
	h.running = false
}

// UpdatePosition updates an option position and its Greeks
func (h *Hedger) UpdatePosition(trade *types.TradeEvent, greeks *OptionGreeks) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	instrumentName := h.constructInstrumentName(trade)
	
	position, exists := h.positions[instrumentName]
	if !exists {
		position = &OptionPosition{
			Instrument: instrumentName,
			Strike:     trade.Strike,
			Expiry:     trade.Expiry,
			IsPut:      trade.IsPut,
			Amount:     decimal.Zero,
		}
		h.positions[instrumentName] = position
	}
	
	// Update position amount
	if trade.IsTakerBuy {
		// We sold, so our position decreases
		position.Amount = position.Amount.Sub(trade.Quantity)
	} else {
		// We bought, so our position increases
		position.Amount = position.Amount.Add(trade.Quantity)
	}
	
	// Update Greeks if provided
	if greeks != nil {
		position.Delta = greeks.Delta
		position.Gamma = greeks.Gamma
		position.Vega = greeks.Vega
		position.Theta = greeks.Theta
		position.UnderlyingPrice = greeks.UnderlyingPrice
		position.ImpliedVol = greeks.ImpliedVol
		position.TimeToExpiry = greeks.TimeToExpiry
	}
	
	position.LastUpdate = time.Now()
	
	// Recalculate portfolio Greeks
	h.recalculatePortfolioGreeks()
}

// performHedgeCheck checks if hedging is needed and executes if necessary
func (h *Hedger) performHedgeCheck() {
	h.mu.RLock()
	netDelta := h.netDelta
	netGamma := h.netGamma
	h.mu.RUnlock()
	
	// Check if we need to hedge based on delta
	if netDelta.Abs().GreaterThan(h.deltaThreshold) {
		log.Printf("Delta hedge triggered: Net Delta = %s", netDelta.StringFixed(4))
		h.executeDeltaHedge(netDelta)
	}
	
	// Check if we need dynamic gamma hedging
	if netGamma.Abs().GreaterThan(h.gammaThreshold) {
		// For high gamma positions, we need more frequent rehedging
		h.checkDynamicGammaHedge(netDelta, netGamma)
	}
}

// executeDeltaHedge executes a delta-neutral hedge
func (h *Hedger) executeDeltaHedge(targetDelta decimal.Decimal) {
	// Calculate hedge size needed
	hedgeSize := targetDelta.Neg() // Opposite of our delta exposure
	
	// Check minimum hedge size
	if hedgeSize.Abs().LessThan(h.minHedgeSize) {
		log.Printf("Hedge size %s below minimum %s, skipping", 
			hedgeSize.StringFixed(4), h.minHedgeSize.StringFixed(4))
		return
	}
	
	// Create a synthetic trade event for the hedge
	hedgeTrade := &types.TradeEvent{
		ID:         fmt.Sprintf("gamma_hedge_%d", time.Now().UnixNano()),
		Source:     types.TradeSourceHedge,
		Status:     types.TradeStatusPending,
		Instrument: "ETH-PERP", // Use perpetual for hedging
		Quantity:   hedgeSize.Abs(),
		IsTakerBuy: hedgeSize.GreaterThan(decimal.Zero), // Buy if we need positive delta
		Timestamp:  time.Now(),
	}
	
	// Execute hedge through hedge manager
	err := h.hedgeManager.ExecuteHedge(h.ctx, hedgeTrade)
	if err != nil {
		log.Printf("Failed to execute delta hedge: %v", err)
		return
	}
	
	// Update hedge tracking
	h.mu.Lock()
	if h.currentHedge == nil {
		h.currentHedge = &HedgePosition{
			Instrument: "ETH-PERP",
			Amount:     decimal.Zero,
		}
	}
	h.currentHedge.Amount = h.currentHedge.Amount.Add(hedgeSize)
	h.currentHedge.LastHedgeTime = time.Now()
	h.currentHedge.LastHedgeDelta = targetDelta
	h.mu.Unlock()
	
	log.Printf("Delta hedge executed: Size = %s, New hedge position = %s", 
		hedgeSize.StringFixed(4), h.currentHedge.Amount.StringFixed(4))
}

// checkDynamicGammaHedge implements dynamic hedging for high gamma positions
func (h *Hedger) checkDynamicGammaHedge(netDelta, netGamma decimal.Decimal) {
	// For high gamma, we rehedge more frequently based on time decay
	h.mu.RLock()
	lastHedgeTime := time.Time{}
	if h.currentHedge != nil {
		lastHedgeTime = h.currentHedge.LastHedgeTime
	}
	h.mu.RUnlock()
	
	// Calculate time-based rehedge threshold
	timeSinceHedge := time.Since(lastHedgeTime)
	
	// For options near expiry with high gamma, hedge more frequently
	var rehedgeThreshold time.Duration
	if h.hasNearExpiryPositions() {
		rehedgeThreshold = 5 * time.Minute  // Hedge every 5 minutes for near expiry
	} else {
		rehedgeThreshold = 30 * time.Minute // Standard rehedge interval
	}
	
	if timeSinceHedge > rehedgeThreshold {
		log.Printf("Dynamic gamma hedge triggered: Gamma = %s, Time since last hedge = %v",
			netGamma.StringFixed(4), timeSinceHedge)
		h.executeDeltaHedge(netDelta)
	}
}

// hasNearExpiryPositions checks if we have positions expiring soon
func (h *Hedger) hasNearExpiryPositions() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	nearExpiryThreshold := 7 * 24 * time.Hour // 7 days
	
	for _, pos := range h.positions {
		if pos.Amount.IsZero() {
			continue
		}
		
		timeToExpiry := time.Until(time.Unix(pos.Expiry, 0))
		if timeToExpiry < nearExpiryThreshold && timeToExpiry > 0 {
			return true
		}
	}
	
	return false
}

// recalculatePortfolioGreeks recalculates net portfolio Greeks
func (h *Hedger) recalculatePortfolioGreeks() {
	h.netDelta = decimal.Zero
	h.netGamma = decimal.Zero
	
	for _, pos := range h.positions {
		if !pos.Amount.IsZero() {
			h.netDelta = h.netDelta.Add(pos.Delta.Mul(pos.Amount))
			h.netGamma = h.netGamma.Add(pos.Gamma.Mul(pos.Amount))
		}
	}
	
	// Include hedge position in delta calculation
	if h.currentHedge != nil && !h.currentHedge.Amount.IsZero() {
		// Perp has delta of 1
		h.netDelta = h.netDelta.Add(h.currentHedge.Amount)
	}
}

// constructInstrumentName creates a standardized instrument name
func (h *Hedger) constructInstrumentName(trade *types.TradeEvent) string {
	optionType := "C"
	if trade.IsPut {
		optionType = "P"
	}
	
	expiryTime := time.Unix(trade.Expiry, 0)
	expiryStr := expiryTime.Format("20060102")
	
	return fmt.Sprintf("ETH-%s-%s-%s", expiryStr, trade.Strike.String(), optionType)
}

// GetMetrics returns current gamma hedging metrics
func (h *Hedger) GetMetrics() GammaMetrics {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	hedgeAmount := decimal.Zero
	if h.currentHedge != nil {
		hedgeAmount = h.currentHedge.Amount
	}
	
	return GammaMetrics{
		NetDelta:       h.netDelta,
		NetGamma:       h.netGamma,
		HedgeAmount:    hedgeAmount,
		PositionCount:  len(h.positions),
		LastUpdateTime: time.Now(),
	}
}

// OptionGreeks contains the Greeks for an option
type OptionGreeks struct {
	Delta           decimal.Decimal
	Gamma           decimal.Decimal
	Vega            decimal.Decimal
	Theta           decimal.Decimal
	UnderlyingPrice decimal.Decimal
	ImpliedVol      decimal.Decimal
	TimeToExpiry    decimal.Decimal
}

// GammaMetrics contains gamma hedging metrics
type GammaMetrics struct {
	NetDelta       decimal.Decimal
	NetGamma       decimal.Decimal
	HedgeAmount    decimal.Decimal
	PositionCount  int
	LastUpdateTime time.Time
}