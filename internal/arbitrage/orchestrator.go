package arbitrage

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	
	"github.com/wakamex/atomizer/internal/config"
	"github.com/wakamex/atomizer/internal/types"
)

// Ensure Orchestrator implements the api.Orchestrator interface
var _ interface {
	SubmitManualTrade(req types.ManualTradeRequest) (*types.TradeEvent, error)
	GetActiveTrades() []types.TradeEvent
} = (*Orchestrator)(nil)

// Orchestrator coordinates the arbitrage flow between RFQ, manual trades, and hedging
type Orchestrator struct {
	config         *config.Config
	exchange       types.Exchange
	hedgeManager   types.HedgeManager
	riskManager    types.RiskManager
	gammaModule    GammaModule
	gammaHedger    GammaHedger
	tradeQueue     chan types.TradeEvent
	activeTrades   map[string]*types.TradeEvent
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
}


// NewOrchestrator creates a new arbitrage orchestrator
func NewOrchestrator(cfg *config.Config, exchange types.Exchange, hedgeManager types.HedgeManager, 
	riskManager types.RiskManager, gammaModule GammaModule, gammaHedger GammaHedger) *Orchestrator {
	
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Orchestrator{
		config:       cfg,
		exchange:     exchange,
		hedgeManager: hedgeManager,
		riskManager:  riskManager,
		gammaModule:  gammaModule,
		gammaHedger:  gammaHedger,
		tradeQueue:   make(chan types.TradeEvent, 100),
		activeTrades: make(map[string]*types.TradeEvent),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start begins the orchestrator's processing loops
func (o *Orchestrator) Start() {
	log.Println("Starting arbitrage orchestrator...")
	
	// Start gamma hedger if enabled
	if o.config.EnableGammaHedging && o.gammaHedger != nil {
		go o.gammaHedger.Start(o.ctx)
	}
	
	// Start trade processor
	go o.processTrades()
	
	// Start periodic risk check
	go o.monitorRisk()
}

// Stop gracefully shuts down the orchestrator
func (o *Orchestrator) Stop() {
	log.Println("Stopping arbitrage orchestrator...")
	o.cancel()
	
	// Stop gamma hedger
	if o.gammaHedger != nil {
		o.gammaHedger.Stop()
	}
}

// SubmitRFQTrade submits an RFQ trade for processing
func (o *Orchestrator) SubmitRFQTrade(rfqResult types.RFQResult, confirmation *types.RFQConfirmation) error {
	trade := types.TradeEvent{
		ID:         uuid.New().String(),
		Source:     types.TradeSourceRysk,
		Status:     types.TradeStatusPending,
		RFQId:      rfqResult.RFQId,
		Instrument: rfqResult.Asset,
		Strike:     DecimalFromString(rfqResult.Strike),
		Expiry:     rfqResult.Expiry,
		IsPut:      rfqResult.IsPut,
		Quantity:   DecimalFromString(rfqResult.Quantity),
		IsTakerBuy: rfqResult.IsTakerBuy,
		Timestamp:  time.Now(),
	}
	
	// If we have a confirmation, update the trade
	if confirmation != nil {
		trade.Status = types.TradeStatusExecuted
		trade.Price = DecimalFromString(confirmation.Price)
	}
	
	// Add to active trades
	o.mu.Lock()
	o.activeTrades[trade.ID] = &trade
	o.mu.Unlock()
	
	// Queue for processing
	select {
	case o.tradeQueue <- trade:
		return nil
	case <-o.ctx.Done():
		return fmt.Errorf("orchestrator shutting down")
	default:
		return fmt.Errorf("trade queue full")
	}
}

// SubmitManualTrade submits a manual trade for processing
func (o *Orchestrator) SubmitManualTrade(req types.ManualTradeRequest) (*types.TradeEvent, error) {
	// Validate the trade with risk manager
	trade := &types.TradeEvent{
		ID:         uuid.New().String(),
		Source:     types.TradeSourceManual,
		Status:     types.TradeStatusPending,
		Instrument: req.Asset,
		Strike:     DecimalFromString(req.Strike),
		Expiry:     req.Expiry,
		IsPut:      req.IsPut,
		Quantity:   req.Quantity,
		IsTakerBuy: req.IsTakerBuy,
		Timestamp:  time.Now(),
	}
	
	// Validate with risk manager
	if err := o.riskManager.ValidateTrade(trade); err != nil {
		trade.Status = types.TradeStatusFailed
		trade.Error = err
		return trade, err
	}
	
	// Add to active trades
	o.mu.Lock()
	o.activeTrades[trade.ID] = trade
	o.mu.Unlock()
	
	// Queue for processing
	select {
	case o.tradeQueue <- *trade:
		return trade, nil
	case <-o.ctx.Done():
		return nil, fmt.Errorf("orchestrator shutting down")
	default:
		return nil, fmt.Errorf("trade queue full")
	}
}

// GetActiveTrades returns all active trades
func (o *Orchestrator) GetActiveTrades() []types.TradeEvent {
	o.mu.RLock()
	defer o.mu.RUnlock()
	
	trades := make([]types.TradeEvent, 0, len(o.activeTrades))
	for _, trade := range o.activeTrades {
		if trade != nil {
			trades = append(trades, *trade)
		}
	}
	return trades
}

// processTrades handles the main trade processing loop
func (o *Orchestrator) processTrades() {
	for {
		select {
		case trade := <-o.tradeQueue:
			o.processTrade(&trade)
			
		case <-o.ctx.Done():
			return
		}
	}
}

// processTrade processes a single trade through its lifecycle
func (o *Orchestrator) processTrade(trade *types.TradeEvent) {
	log.Printf("Processing trade %s from %s: %s %s %s", 
		trade.ID, trade.Source, trade.Instrument, 
		trade.Quantity.String(), 
		map[bool]string{true: "BUY", false: "SELL"}[trade.IsTakerBuy])
	
	// Update status
	o.updateTradeStatus(trade.ID, types.TradeStatusExecuted)
	
	// Update risk manager
	o.riskManager.UpdatePosition(trade)
	
	// Execute hedge if not already a hedge trade
	if trade.Source != types.TradeSourceHedge {
		if err := o.hedgeManager.ExecuteHedge(o.ctx, trade); err != nil {
			log.Printf("Failed to hedge trade %s: %v", trade.ID, err)
			o.updateTradeStatus(trade.ID, types.TradeStatusFailed)
			trade.Error = err
			return
		}
		o.updateTradeStatus(trade.ID, types.TradeStatusHedged)
	}
	
	// Clean up completed trades after some time
	go func() {
		time.Sleep(5 * time.Minute)
		o.mu.Lock()
		delete(o.activeTrades, trade.ID)
		o.mu.Unlock()
	}()
}

// updateTradeStatus updates the status of a trade
func (o *Orchestrator) updateTradeStatus(tradeID string, status types.TradeStatus) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	if trade, exists := o.activeTrades[tradeID]; exists {
		trade.Status = status
		log.Printf("Trade %s status updated to %s", tradeID, status)
	}
}

// monitorRisk periodically checks risk levels
func (o *Orchestrator) monitorRisk() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			delta, gamma := o.riskManager.GetGreeks()
			log.Printf("Portfolio Greeks - Delta: %s, Gamma: %s", 
				delta.StringFixed(4), gamma.StringFixed(4))
			
			// Check if we need to hedge based on gamma
			if o.config.EnableGammaHedging && o.gammaModule != nil {
				if o.gammaModule.ShouldHedge(gamma) {
					log.Printf("Gamma threshold exceeded, initiating hedge...")
					// Gamma hedger runs independently
				}
			}
			
		case <-o.ctx.Done():
			return
		}
	}
}

// DecimalFromString safely converts string to decimal
func DecimalFromString(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		log.Printf("Error parsing decimal from string %s: %v", s, err)
		return decimal.Zero
	}
	return d
}