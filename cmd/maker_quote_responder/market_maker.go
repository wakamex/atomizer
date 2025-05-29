package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// MarketMaker manages automated market making
type MarketMaker struct {
	config   *MarketMakerConfig
	exchange MarketMakerExchange
	
	// Order tracking
	activeOrders map[string]*MarketMakerOrder // orderID -> Order
	ordersByInstrument map[string]map[string]*MarketMakerOrder // instrument -> side -> orders
	
	// Market data
	latestTickers map[string]*TickerUpdate
	
	// Position tracking
	positions map[string]decimal.Decimal // instrument -> net position
	
	// Statistics
	stats MarketMakerStats
	
	// Control
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex
}

// NewMarketMaker creates a new market maker instance
func NewMarketMaker(config *MarketMakerConfig, exchange MarketMakerExchange) *MarketMaker {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &MarketMaker{
		config:             config,
		exchange:           exchange,
		activeOrders:       make(map[string]*MarketMakerOrder),
		ordersByInstrument: make(map[string]map[string]*MarketMakerOrder),
		latestTickers:      make(map[string]*TickerUpdate),
		positions:          make(map[string]decimal.Decimal),
		stats:              MarketMakerStats{BidAskSpread: make(map[string]decimal.Decimal)},
		ctx:                ctx,
		cancel:             cancel,
	}
}

// Start begins market making
func (mm *MarketMaker) Start() error {
	log.Printf("Starting market maker for %d instruments", len(mm.config.Instruments))
	
	// Load existing positions
	if err := mm.loadPositions(); err != nil {
		return fmt.Errorf("failed to load positions: %w", err)
	}
	
	// Load existing orders
	if err := mm.loadActiveOrders(); err != nil {
		return fmt.Errorf("failed to load active orders: %w", err)
	}
	
	// Subscribe to ticker updates
	tickerChan, err := mm.exchange.SubscribeTickers(mm.ctx, mm.config.Instruments)
	if err != nil {
		return fmt.Errorf("failed to subscribe to tickers: %w", err)
	}
	
	// Start ticker processor
	mm.wg.Add(1)
	go mm.processTickers(tickerChan)
	
	// Start quote updater
	mm.wg.Add(1)
	go mm.quoteUpdater()
	
	// Start statistics reporter
	mm.wg.Add(1)
	go mm.statsReporter()
	
	mm.stats.UptimeSeconds = 0
	log.Println("Market maker started successfully")
	
	return nil
}

// Stop gracefully shuts down the market maker
func (mm *MarketMaker) Stop() error {
	log.Println("Stopping market maker...")
	
	// Cancel context to stop all goroutines
	mm.cancel()
	
	// Cancel all active orders
	mm.cancelAllOrders()
	
	// Wait for goroutines to finish
	mm.wg.Wait()
	
	log.Println("Market maker stopped")
	return nil
}

// processTickers handles incoming ticker updates
func (mm *MarketMaker) processTickers(tickerChan <-chan TickerUpdate) {
	defer mm.wg.Done()
	
	for {
		select {
		case <-mm.ctx.Done():
			return
		case ticker, ok := <-tickerChan:
			if !ok {
				log.Println("Ticker channel closed")
				return
			}
			
			mm.mu.Lock()
			mm.latestTickers[ticker.Instrument] = &ticker
			mm.mu.Unlock()
			
			// Check if we need to update quotes
			if mm.shouldUpdateQuotes(ticker.Instrument) {
				mm.updateQuotesForInstrument(ticker.Instrument)
			}
		}
	}
}

// quoteUpdater periodically updates all quotes
func (mm *MarketMaker) quoteUpdater() {
	defer mm.wg.Done()
	
	ticker := time.NewTicker(mm.config.RefreshInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-mm.ctx.Done():
			return
		case <-ticker.C:
			mm.updateAllQuotes()
		}
	}
}

// updateAllQuotes updates quotes for all instruments
func (mm *MarketMaker) updateAllQuotes() {
	for _, instrument := range mm.config.Instruments {
		if err := mm.updateQuotesForInstrument(instrument); err != nil {
			log.Printf("Failed to update quotes for %s: %v", instrument, err)
		}
	}
}

// updateQuotesForInstrument updates quotes for a specific instrument
func (mm *MarketMaker) updateQuotesForInstrument(instrument string) error {
	mm.mu.RLock()
	ticker, exists := mm.latestTickers[instrument]
	mm.mu.RUnlock()
	
	if !exists || ticker == nil {
		return fmt.Errorf("no ticker data for %s", instrument)
	}
	
	// Calculate our quotes
	bidPrice, askPrice := mm.calculateQuotes(ticker)
	
	// Check risk limits
	if !mm.checkRiskLimits(instrument, mm.config.QuoteSize) {
		log.Printf("Risk limits exceeded for %s, skipping quote update", instrument)
		return nil
	}
	
	// Cancel existing orders
	mm.cancelOrdersForInstrument(instrument)
	
	// Place new orders
	if err := mm.placeQuotes(instrument, bidPrice, askPrice); err != nil {
		return fmt.Errorf("failed to place quotes: %w", err)
	}
	
	return nil
}

// calculateQuotes calculates bid and ask prices based on current market
func (mm *MarketMaker) calculateQuotes(ticker *TickerUpdate) (bidPrice, askPrice decimal.Decimal) {
	// Calculate mid price
	midPrice := ticker.BestBid.Add(ticker.BestAsk).Div(decimal.NewFromInt(2))
	
	// Our quotes: improve the market
	// Bid: best bid + 0.1 (or configured spread)
	// Ask: best ask - 0.1 (or configured spread)
	improvementAmount := decimal.NewFromFloat(0.1)
	
	bidPrice = ticker.BestBid.Add(improvementAmount)
	askPrice = ticker.BestAsk.Sub(improvementAmount)
	
	// Ensure minimum spread
	minSpread := midPrice.Mul(decimal.NewFromInt(int64(mm.config.MinSpreadBps)).Div(decimal.NewFromInt(10000)))
	if askPrice.Sub(bidPrice).LessThan(minSpread) {
		// Widen the spread
		bidPrice = midPrice.Sub(minSpread.Div(decimal.NewFromInt(2)))
		askPrice = midPrice.Add(minSpread.Div(decimal.NewFromInt(2)))
	}
	
	return bidPrice, askPrice
}

// placeQuotes places bid and ask orders
func (mm *MarketMaker) placeQuotes(instrument string, bidPrice, askPrice decimal.Decimal) error {
	// Place bid order
	bidOrderID, err := mm.exchange.PlaceLimitOrder(instrument, "buy", bidPrice, mm.config.QuoteSize)
	if err != nil {
		return fmt.Errorf("failed to place bid order: %w", err)
	}
	
	// Place ask order
	askOrderID, err := mm.exchange.PlaceLimitOrder(instrument, "sell", askPrice, mm.config.QuoteSize)
	if err != nil {
		// Cancel the bid order since we couldn't place the ask
		mm.exchange.CancelOrder(bidOrderID)
		return fmt.Errorf("failed to place ask order: %w", err)
	}
	
	// Track orders
	mm.mu.Lock()
	mm.trackOrder(bidOrderID, instrument, "buy", bidPrice, mm.config.QuoteSize)
	mm.trackOrder(askOrderID, instrument, "sell", askPrice, mm.config.QuoteSize)
	mm.stats.OrdersPlaced += 2
	mm.mu.Unlock()
	
	log.Printf("Placed quotes for %s: Bid %.4f, Ask %.4f", instrument, bidPrice, askPrice)
	
	return nil
}

// trackOrder adds an order to tracking
func (mm *MarketMaker) trackOrder(orderID, instrument, side string, price, amount decimal.Decimal) {
	order := &MarketMakerOrder{
		OrderID:    orderID,
		Instrument: instrument,
		Side:       side,
		Price:      price,
		Amount:     amount,
		Status:     "open",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	mm.activeOrders[orderID] = order
	
	if mm.ordersByInstrument[instrument] == nil {
		mm.ordersByInstrument[instrument] = make(map[string]*MarketMakerOrder)
	}
	mm.ordersByInstrument[instrument][side] = order
}

// cancelOrdersForInstrument cancels all orders for an instrument
func (mm *MarketMaker) cancelOrdersForInstrument(instrument string) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	orders, exists := mm.ordersByInstrument[instrument]
	if !exists {
		return
	}
	
	for _, order := range orders {
		if err := mm.exchange.CancelOrder(order.OrderID); err != nil {
			log.Printf("Failed to cancel order %s: %v", order.OrderID, err)
		} else {
			delete(mm.activeOrders, order.OrderID)
			mm.stats.OrdersCancelled++
		}
	}
	
	delete(mm.ordersByInstrument, instrument)
}

// cancelAllOrders cancels all active orders
func (mm *MarketMaker) cancelAllOrders() {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	for orderID := range mm.activeOrders {
		if err := mm.exchange.CancelOrder(orderID); err != nil {
			log.Printf("Failed to cancel order %s: %v", orderID, err)
		} else {
			mm.stats.OrdersCancelled++
		}
	}
	
	mm.activeOrders = make(map[string]*MarketMakerOrder)
	mm.ordersByInstrument = make(map[string]map[string]*MarketMakerOrder)
}

// shouldUpdateQuotes checks if quotes need updating
func (mm *MarketMaker) shouldUpdateQuotes(instrument string) bool {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	orders, exists := mm.ordersByInstrument[instrument]
	if !exists || len(orders) == 0 {
		return true // No orders, need to place them
	}
	
	ticker, exists := mm.latestTickers[instrument]
	if !exists {
		return false // No ticker data
	}
	
	// Check if market has moved significantly
	for side, order := range orders {
		var marketPrice decimal.Decimal
		if side == "buy" {
			marketPrice = ticker.BestBid
		} else {
			marketPrice = ticker.BestAsk
		}
		
		priceDiff := order.Price.Sub(marketPrice).Abs()
		if priceDiff.GreaterThan(order.Price.Mul(mm.config.CancelThreshold)) {
			return true
		}
	}
	
	return false
}

// checkRiskLimits checks if placing an order would exceed risk limits
func (mm *MarketMaker) checkRiskLimits(instrument string, size decimal.Decimal) bool {
	mm.mu.RLock()
	defer mm.mu.RUnlock()
	
	// Check position limit for instrument
	currentPosition := mm.positions[instrument]
	if currentPosition.Add(size).Abs().GreaterThan(mm.config.MaxPositionSize) {
		return false
	}
	
	// Check total exposure
	totalExposure := decimal.Zero
	for _, pos := range mm.positions {
		totalExposure = totalExposure.Add(pos.Abs())
	}
	
	if totalExposure.Add(size).GreaterThan(mm.config.MaxTotalExposure) {
		return false
	}
	
	return true
}

// loadPositions loads current positions from exchange
func (mm *MarketMaker) loadPositions() error {
	positions, err := mm.exchange.GetPositions()
	if err != nil {
		return err
	}
	
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	for _, pos := range positions {
		amount := decimal.NewFromFloat(pos.Amount)
		if pos.Direction == "sell" {
			amount = amount.Neg()
		}
		mm.positions[pos.InstrumentName] = amount
	}
	
	log.Printf("Loaded %d positions", len(positions))
	return nil
}

// loadActiveOrders loads current open orders from exchange
func (mm *MarketMaker) loadActiveOrders() error {
	orders, err := mm.exchange.GetOpenOrders()
	if err != nil {
		return err
	}
	
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	for _, order := range orders {
		mm.activeOrders[order.OrderID] = &order
		
		if mm.ordersByInstrument[order.Instrument] == nil {
			mm.ordersByInstrument[order.Instrument] = make(map[string]*MarketMakerOrder)
		}
		mm.ordersByInstrument[order.Instrument][order.Side] = &order
	}
	
	log.Printf("Loaded %d active orders", len(orders))
	return nil
}

// statsReporter periodically reports statistics
func (mm *MarketMaker) statsReporter() {
	defer mm.wg.Done()
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	startTime := time.Now()
	
	for {
		select {
		case <-mm.ctx.Done():
			return
		case <-ticker.C:
			mm.mu.RLock()
			mm.stats.UptimeSeconds = int64(time.Since(startTime).Seconds())
			mm.stats.LastUpdate = time.Now()
			
			log.Printf("Market Maker Stats: Orders placed: %d, Cancelled: %d, Filled: %d, Uptime: %ds",
				mm.stats.OrdersPlaced,
				mm.stats.OrdersCancelled,
				mm.stats.OrdersFilled,
				mm.stats.UptimeSeconds)
			mm.mu.RUnlock()
		}
	}
}