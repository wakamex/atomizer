package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TradeSourceType represents the source of a trade
type TradeSourceType string

const (
	TradeSourceRysk   TradeSourceType = "RYSK_RFQ"
	TradeSourceManual TradeSourceType = "MANUAL"
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

// ArbitrageOrchestrator coordinates the arbitrage flow
type ArbitrageOrchestrator struct {
	config         *AppConfig
	hedgeManager   *HedgeManager
	riskManager    *RiskManager
	gammaModule    *GammaModule
	tradeQueue     chan TradeEvent
	activeTrades   map[string]*TradeEvent
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewArbitrageOrchestrator creates a new orchestrator
func NewArbitrageOrchestrator(cfg *AppConfig, exchange Exchange) *ArbitrageOrchestrator {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &ArbitrageOrchestrator{
		config:       cfg,
		hedgeManager: NewHedgeManager(exchange, cfg),
		riskManager:  NewRiskManager(cfg),
		gammaModule:  NewGammaModule(cfg),
		tradeQueue:   make(chan TradeEvent, 100),
		activeTrades: make(map[string]*TradeEvent),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start begins the orchestrator's async processing
func (o *ArbitrageOrchestrator) Start() error {
	log.Println("Starting arbitrage orchestrator...")
	
	// Start trade processor
	go o.processTradeQueue()
	
	// Start gamma hedging if enabled
	if o.config.GammaHedgingEnabled {
		go o.gammaModule.Start(o.ctx)
	}
	
	return nil
}

// Stop gracefully shuts down the orchestrator
func (o *ArbitrageOrchestrator) Stop() {
	log.Println("Stopping arbitrage orchestrator...")
	o.cancel()
}

// SubmitRFQTrade converts an RFQ into a trade event
func (o *ArbitrageOrchestrator) SubmitRFQTrade(rfq RFQRequest) (*TradeEvent, error) {
	trade := TradeEvent{
		ID:         uuid.New().String(),
		Source:     TradeSourceRysk,
		Status:     TradeStatusPending,
		RFQId:      rfq.RfqId,
		Instrument: o.buildInstrumentName(rfq),
		Strike:     decimal.NewFromBigInt(rfq.Strike, 0),
		Expiry:     rfq.Expiry.Int64(),
		IsPut:      rfq.IsPut,
		Quantity:   decimal.NewFromBigInt(rfq.Size, -18), // Convert from wei
		IsTakerBuy: rfq.IsBuyOrder,
		Timestamp:  time.Now(),
	}
	
	// Store and queue
	o.mu.Lock()
	o.activeTrades[trade.ID] = &trade
	o.mu.Unlock()
	
	select {
	case o.tradeQueue <- trade:
		log.Printf("Queued RFQ trade %s for processing", trade.ID)
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("trade queue timeout")
	}
	
	return &trade, nil
}

// SubmitManualTrade handles manually initiated trades
func (o *ArbitrageOrchestrator) SubmitManualTrade(req ManualTradeRequest) (*TradeEvent, error) {
	trade := TradeEvent{
		ID:         uuid.New().String(),
		Source:     TradeSourceManual,
		Status:     TradeStatusPending,
		Instrument: req.Instrument,
		Strike:     req.Strike,
		Expiry:     req.Expiry,
		IsPut:      req.IsPut,
		Quantity:   req.Quantity,
		Price:      req.Price,
		IsTakerBuy: false, // Manual trades are maker sells
		Timestamp:  time.Now(),
	}
	
	// Validate with risk manager
	if err := o.riskManager.ValidateTrade(&trade); err != nil {
		trade.Status = TradeStatusCancelled
		trade.Error = err
		return &trade, err
	}
	
	// Store and queue
	o.mu.Lock()
	o.activeTrades[trade.ID] = &trade
	o.mu.Unlock()
	
	select {
	case o.tradeQueue <- trade:
		log.Printf("Queued manual trade %s for processing", trade.ID)
	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("trade queue timeout")
	}
	
	return &trade, nil
}

// processTradeQueue handles trades asynchronously
func (o *ArbitrageOrchestrator) processTradeQueue() {
	for {
		select {
		case trade := <-o.tradeQueue:
			go o.executeTrade(trade)
		case <-o.ctx.Done():
			return
		}
	}
}

// executeTrade processes a single trade through the full flow
func (o *ArbitrageOrchestrator) executeTrade(trade TradeEvent) {
	log.Printf("Executing trade %s from %s", trade.ID, trade.Source)
	
	// Update status
	o.updateTradeStatus(trade.ID, TradeStatusExecuted)
	
	// For manual trades, we skip to hedging
	// For RFQ trades, quote response and execution are handled by existing flow
	
	if trade.Source == TradeSourceManual || trade.Status == TradeStatusExecuted {
		// Execute hedge
		hedgeResult, err := o.hedgeManager.ExecuteHedge(&trade)
		if err != nil {
			log.Printf("Hedge failed for trade %s: %v", trade.ID, err)
			o.updateTradeStatus(trade.ID, TradeStatusFailed)
			o.updateTradeError(trade.ID, err)
			return
		}
		
		// Update with hedge info
		o.mu.Lock()
		if t, exists := o.activeTrades[trade.ID]; exists {
			t.HedgeOrderID = hedgeResult.OrderID
			t.HedgeExchange = hedgeResult.Exchange
		}
		o.mu.Unlock()
		
		o.updateTradeStatus(trade.ID, TradeStatusHedged)
		
		// Notify gamma module of new position
		o.gammaModule.OnNewPosition(&trade)
	}
}

// OnRFQConfirmation handles trade confirmations from Rysk
func (o *ArbitrageOrchestrator) OnRFQConfirmation(conf RFQConfirmation) {
	// Find the trade by RFQ ID
	o.mu.RLock()
	var trade *TradeEvent
	for _, t := range o.activeTrades {
		if t.RFQId == conf.RfqId {
			trade = t
			break
		}
	}
	o.mu.RUnlock()
	
	if trade == nil {
		log.Printf("No trade found for RFQ confirmation %s", conf.RfqId)
		return
	}
	
	// Update status and trigger hedge
	o.updateTradeStatus(trade.ID, TradeStatusExecuted)
	
	// Re-queue for hedge execution
	select {
	case o.tradeQueue <- *trade:
	default:
		log.Printf("Failed to queue trade %s for hedging", trade.ID)
	}
}

// GetActiveTrades returns current active trades
func (o *ArbitrageOrchestrator) GetActiveTrades() []TradeEvent {
	o.mu.RLock()
	defer o.mu.RUnlock()
	
	trades := make([]TradeEvent, 0, len(o.activeTrades))
	for _, trade := range o.activeTrades {
		trades = append(trades, *trade)
	}
	return trades
}

// Helper methods

func (o *ArbitrageOrchestrator) updateTradeStatus(tradeID string, status TradeStatus) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	if trade, exists := o.activeTrades[tradeID]; exists {
		trade.Status = status
		log.Printf("Trade %s status updated to %s", tradeID, status)
	}
}

func (o *ArbitrageOrchestrator) updateTradeError(tradeID string, err error) {
	o.mu.Lock()
	defer o.mu.Unlock()
	
	if trade, exists := o.activeTrades[tradeID]; exists {
		trade.Error = err
	}
}

func (o *ArbitrageOrchestrator) buildInstrumentName(rfq RFQRequest) string {
	// Convert to exchange instrument format
	// This is a simplified version - actual implementation would use exchange methods
	strikeStr := decimal.NewFromBigInt(rfq.Strike, -8).String()
	expiryTime := time.Unix(rfq.Expiry.Int64(), 0)
	optionType := "C"
	if rfq.IsPut {
		optionType = "P"
	}
	
	return fmt.Sprintf("ETH-%s-%s-%s", 
		expiryTime.Format("20060102"),
		strikeStr,
		optionType,
	)
}

// ManualTradeRequest represents a manual trade submission
type ManualTradeRequest struct {
	Instrument string          `json:"instrument"`
	Strike     decimal.Decimal `json:"strike"`
	Expiry     int64           `json:"expiry"`
	IsPut      bool            `json:"is_put"`
	Quantity   decimal.Decimal `json:"quantity"`
	Price      decimal.Decimal `json:"price"`
}