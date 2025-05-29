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
	
	// Check if we need to update quotes
	// Only update if market has moved more than 0.5% or we don't have orders
	mm.mu.RLock()
	existingOrders := mm.ordersByInstrument[instrument]
	mm.mu.RUnlock()
	
	needsUpdate := false
	if existingOrders == nil || len(existingOrders) == 0 {
		needsUpdate = true
	} else {
		// Check if market has moved significantly
		threshold := decimal.NewFromFloat(0.005) // 0.5% threshold
		
		if bidOrder, exists := existingOrders["buy"]; exists {
			priceDiff := bidPrice.Sub(bidOrder.Price).Abs()
			if priceDiff.Div(bidOrder.Price).GreaterThan(threshold) {
				needsUpdate = true
			}
		}
		
		if askOrder, exists := existingOrders["sell"]; exists {
			priceDiff := askPrice.Sub(askOrder.Price).Abs()
			if priceDiff.Div(askOrder.Price).GreaterThan(threshold) {
				needsUpdate = true
			}
		}
	}
	
	if !needsUpdate {
		// Market hasn't moved enough, skip update
		log.Printf("Market hasn't moved significantly for %s, keeping existing orders", instrument)
		return nil
	}
	
	// Check risk limits
	if !mm.checkRiskLimits(instrument, mm.config.QuoteSize) {
		log.Printf("Risk limits exceeded for %s, skipping quote update", instrument)
		return nil
	}
	
	// Update quotes - use replace if we have existing orders
	if err := mm.updateQuotes(instrument, bidPrice, askPrice); err != nil {
		return fmt.Errorf("failed to update quotes: %w", err)
	}
	
	return nil
}

// calculateQuotes calculates bid and ask prices based on current market
func (mm *MarketMaker) calculateQuotes(ticker *TickerUpdate) (bidPrice, askPrice decimal.Decimal) {
	// Calculate mid price
	midPrice := ticker.BestBid.Add(ticker.BestAsk).Div(decimal.NewFromInt(2))
	
	// Our quotes: improve the market
	// Bid: best bid + improvement amount
	// Ask: best ask - improvement amount
	improvementAmount := mm.config.PriceImprovement
	
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

// updateQuotes updates quotes for an instrument, using replace when possible
func (mm *MarketMaker) updateQuotes(instrument string, bidPrice, askPrice decimal.Decimal) error {
	mm.mu.RLock()
	existingOrders := mm.ordersByInstrument[instrument]
	mm.mu.RUnlock()
	
	// If we have existing orders, try to replace them
	if existingOrders != nil && len(existingOrders) > 0 {
		var bidOrderID, askOrderID string
		var replacedBid, replacedAsk bool
		
		// Replace bid order if it exists
		if bidOrder, exists := existingOrders["buy"]; exists {
			newBidOrderID, err := mm.exchange.ReplaceOrder(bidOrder.OrderID, instrument, "buy", bidPrice, mm.config.QuoteSize)
			if err != nil {
				log.Printf("Failed to replace bid order %s: %v, will cancel and recreate", bidOrder.OrderID, err)
			} else {
				bidOrderID = newBidOrderID
				replacedBid = true
				log.Printf("Replaced bid order %s with %s at %.4f", bidOrder.OrderID, newBidOrderID, bidPrice)
			}
		}
		
		// Replace ask order if it exists
		if askOrder, exists := existingOrders["sell"]; exists {
			newAskOrderID, err := mm.exchange.ReplaceOrder(askOrder.OrderID, instrument, "sell", askPrice, mm.config.QuoteSize)
			if err != nil {
				log.Printf("Failed to replace ask order %s: %v, will cancel and recreate", askOrder.OrderID, err)
			} else {
				askOrderID = newAskOrderID
				replacedAsk = true
				log.Printf("Replaced ask order %s with %s at %.4f", askOrder.OrderID, newAskOrderID, askPrice)
			}
		}
		
		// Update tracking for replaced orders
		mm.mu.Lock()
		if replacedBid && bidOrderID != "" {
			// Remove old order
			if oldBid, exists := existingOrders["buy"]; exists {
				delete(mm.activeOrders, oldBid.OrderID)
			}
			// Track new order
			mm.trackOrder(bidOrderID, instrument, "buy", bidPrice, mm.config.QuoteSize)
			mm.stats.OrdersReplaced++
		}
		if replacedAsk && askOrderID != "" {
			// Remove old order
			if oldAsk, exists := existingOrders["sell"]; exists {
				delete(mm.activeOrders, oldAsk.OrderID)
			}
			// Track new order
			mm.trackOrder(askOrderID, instrument, "sell", askPrice, mm.config.QuoteSize)
			mm.stats.OrdersReplaced++
		}
		mm.mu.Unlock()
		
		// If both were replaced successfully, we're done
		if replacedBid && replacedAsk {
			log.Printf("Updated quotes for %s: Bid %.4f, Ask %.4f", instrument, bidPrice, askPrice)
			return nil
		}
		
		// Otherwise, cancel remaining orders and place new ones
		if !replacedBid || !replacedAsk {
			mm.cancelOrdersForInstrument(instrument)
			return mm.placeQuotes(instrument, bidPrice, askPrice)
		}
	}
	
	// No existing orders, place new ones
	return mm.placeQuotes(instrument, bidPrice, askPrice)
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