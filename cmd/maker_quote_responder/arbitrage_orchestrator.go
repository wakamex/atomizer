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

// ArbitrageOrchestrator coordinates the arbitrage flow
type ArbitrageOrchestrator struct {
	config         *AppConfig
	hedgeManager   *HedgeManager
	riskManager    *RiskManager
	gammaModule    *GammaDDHAlgo
	gammaHedger    *GammaHedger  // New integrated gamma hedger
	tradeQueue     chan TradeEvent
	activeTrades   map[string]*TradeEvent
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
}

// NewArbitrageOrchestrator creates a new orchestrator
func NewArbitrageOrchestrator(cfg *AppConfig, exchange Exchange, existingPositions []ExchangePosition) *ArbitrageOrchestrator {
	ctx, cancel := context.WithCancel(context.Background())
	
	// Create risk manager
	riskManager := NewRiskManager(cfg)
	
	// Initialize risk manager with existing positions
	if len(existingPositions) > 0 {
		log.Printf("Initializing risk manager with %d existing positions", len(existingPositions))
		for _, pos := range existingPositions {
			// Convert Derive position to internal position format
			// For now, we'll just log them - actual position tracking would need
			// to parse the instrument name and update the risk manager's state
			log.Printf("  Risk Manager: Adding position %s %s %.4f", 
				pos.Direction, pos.InstrumentName, pos.Amount)
			// TODO: Parse instrument name to extract asset, strike, expiry, etc.
			// and update risk manager's position tracking
		}
	}
	
	hedgeManager := NewHedgeManager(exchange, cfg)
	gammaHedger := NewGammaHedger(exchange, cfg, hedgeManager)
	
	return &ArbitrageOrchestrator{
		config:       cfg,
		hedgeManager: hedgeManager,
		riskManager:  riskManager,
		gammaModule:  NewGammaDDHAlgo(exchange, cfg.GammaThreshold),
		gammaHedger:  gammaHedger,
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
	if o.config.EnableGammaHedging {
		// Start the integrated gamma hedger
		if err := o.gammaHedger.Start(); err != nil {
			log.Printf("Failed to start gamma hedger: %v", err)
		}
		// Also start the legacy gamma module if needed
		go o.gammaModule.Start(o.ctx)
	}
	
	return nil
}

// Stop gracefully shuts down the orchestrator
func (o *ArbitrageOrchestrator) Stop() {
	log.Println("Stopping arbitrage orchestrator...")
	o.cancel()
	
	// Stop gamma hedger if running
	if o.config.EnableGammaHedging {
		o.gammaHedger.Stop()
	}
}

// SubmitRFQTrade converts an RFQ into a trade event
func (o *ArbitrageOrchestrator) SubmitRFQTrade(rfq RFQResult) (*TradeEvent, error) {
	trade := TradeEvent{
		ID:         uuid.New().String(),
		Source:     TradeSourceRysk,
		Status:     TradeStatusPending,
		RFQId:      rfq.ID,
		Instrument: o.buildInstrumentName(rfq),
		Strike:     DecimalFromString(rfq.Strike),
		Expiry:     rfq.Expiry,
		IsPut:      rfq.IsPut,
		Quantity:   DecimalFromString(rfq.Quantity).Div(decimal.New(1, 18)), // Convert from wei
		IsTakerBuy: rfq.IsTakerBuy,
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
	// Convert strike to wei format (multiply by 10^8) for consistency with RFQ trades
	strikeInWei := req.Strike.Mul(decimal.NewFromInt(100000000))
	
	trade := TradeEvent{
		ID:         uuid.New().String(),
		Source:     TradeSourceManual,
		Status:     TradeStatusPending,
		Instrument: req.Instrument,
		Strike:     strikeInWei,
		Expiry:     req.Expiry,
		IsPut:      req.IsPut,
		Quantity:   req.Quantity,
		Price:      req.Price,
		IsTakerBuy: req.Side == "sell", // For manual trades, invert: if we buy, taker (market) is selling to us
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
// HandleManualTrade handles a manually submitted trade
func (o *ArbitrageOrchestrator) HandleManualTrade(trade *TradeEvent) error {
	select {
	case o.tradeQueue <- *trade:
		log.Printf("Manual trade %s queued for processing", trade.ID)
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("failed to queue trade %s: timeout", trade.ID)
	}
}

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
	
	// Handle based on trade source
	if trade.Source == TradeSourceManual {
		// For manual trades, place the actual order on the exchange
		log.Printf("Executing manual trade order on exchange")
		
		// Build order confirmation for exchange
		confirmation := RFQConfirmation{
			Nonce:      trade.ID,
			Quantity:   trade.Quantity.Mul(decimal.New(1, 18)).String(), // Convert to wei
			IsTakerBuy: trade.IsTakerBuy,
			Strike:     trade.Strike.String(),
			Expiry:     int(trade.Expiry),
			IsPut:      trade.IsPut,
			Price:      trade.Price.String(), // Use the manual trade's specified price
		}
		
		// Place the order using hedge manager's exchange
		err := o.hedgeManager.exchange.PlaceOrder(confirmation, trade.Instrument, o.config)
		if err != nil {
			log.Printf("Manual trade order failed for %s: %v", trade.ID, err)
			o.updateTradeStatus(trade.ID, TradeStatusFailed)
			o.updateTradeError(trade.ID, err)
			return
		}
		
		log.Printf("Manual trade order placed successfully for %s", trade.ID)
		o.updateTradeStatus(trade.ID, TradeStatusExecuted)
		
	} else if trade.Status == TradeStatusExecuted {
		// For RFQ trades that are executed, proceed with hedging
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
		if o.gammaModule != nil {
			o.gammaModule.OnNewPosition(trade.Instrument, trade.Quantity, trade.Price)
		}
	}
}

// OnRFQConfirmation handles trade confirmations from Rysk
func (o *ArbitrageOrchestrator) OnRFQConfirmation(conf RFQConfirmation) {
	// Find the trade by RFQ ID
	o.mu.RLock()
	var trade *TradeEvent
	for _, t := range o.activeTrades {
		if t.RFQId == conf.Nonce {
			trade = t
			break
		}
	}
	o.mu.RUnlock()
	
	if trade == nil {
		log.Printf("No trade found for RFQ confirmation %s", conf.Nonce)
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

func (o *ArbitrageOrchestrator) buildInstrumentName(rfq RFQResult) string {
	// Convert to exchange instrument format
	// This is a simplified version - actual implementation would use exchange methods
	strikeStr := DecimalFromString(rfq.Strike).Div(decimal.New(1, 8)).String()
	expiryTime := time.Unix(rfq.Expiry, 0)
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
	Side       string          `json:"side"` // "buy" or "sell" - which side of the order book
}