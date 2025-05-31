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
	
	// Order tracking - maintains exactly 1 buy and 1 sell order per instrument
	activeOrders map[string]*MarketMakerOrder // orderID -> Order
	ordersByInstrument map[string]map[string]*MarketMakerOrder // instrument -> side -> order (only one per side)
	
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
	
	// Error suppression
	orderbookErrorLogged map[string]bool
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
		orderbookErrorLogged: make(map[string]bool),
	}
}

// Start begins market making
func (mm *MarketMaker) Start() error {
	log.Printf("Starting market maker: %d instruments, 1 buy + 1 sell per instrument", len(mm.config.Instruments))
	
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
	
	// Fetch orderbook if reference size is configured
	var orderBook *MarketMakerOrderBook
	if mm.config.ImprovementReferenceSize.GreaterThan(decimal.Zero) {
		var err error
		orderBook, err = mm.exchange.GetOrderBook(instrument)
		if err != nil {
			mm.mu.Lock()
			if !mm.orderbookErrorLogged[instrument] {
				log.Printf("Failed to fetch orderbook for %s: %v, using ticker data", instrument, err)
				mm.orderbookErrorLogged[instrument] = true
			}
			mm.mu.Unlock()
			// Continue with ticker data only
		} else {
			// Clear error flag on success
			mm.mu.Lock()
			delete(mm.orderbookErrorLogged, instrument)
			mm.mu.Unlock()
		}
	}
	
	// Calculate our quotes
	bidPrice, askPrice := mm.calculateQuotes(ticker, orderBook)
	
	// Check risk limits
	if !mm.checkRiskLimits(instrument, mm.config.QuoteSize) {
		log.Printf("Risk limits exceeded for %s, skipping quote update", instrument)
		return nil
	}
	
	// Get existing orders for this instrument
	mm.mu.RLock()
	existingOrders := mm.ordersByInstrument[instrument]
	mm.mu.RUnlock()
	
	// If we have existing orders, try to replace them
	if existingOrders != nil && len(existingOrders) > 0 {
		var bidOrder, askOrder *MarketMakerOrder
		
		// Find existing bid and ask orders
		for side, order := range existingOrders {
			if side == "buy" {
				bidOrder = order
			} else if side == "sell" {
				askOrder = order
			}
		}
		
		// Replace existing orders if they exist
		var replacedBid, replacedAsk bool
		
		if bidOrder != nil {
			newOrderID, err := mm.exchange.ReplaceOrder(bidOrder.OrderID, instrument, "buy", bidPrice, mm.config.QuoteSize)
			if err != nil {
				log.Printf("Failed to replace bid order %s: %v, will cancel and recreate", bidOrder.OrderID, err)
			} else {
				// Update our tracking
				mm.updateOrderTracking(bidOrder.OrderID, newOrderID, instrument, "buy", bidPrice, mm.config.QuoteSize)
				replacedBid = true
				debugLog("Replaced bid order %s with %s for %s @ %s", bidOrder.OrderID, newOrderID, instrument, bidPrice)
			}
		}
		
		if askOrder != nil {
			newOrderID, err := mm.exchange.ReplaceOrder(askOrder.OrderID, instrument, "sell", askPrice, mm.config.QuoteSize)
			if err != nil {
				log.Printf("Failed to replace ask order %s: %v, will cancel and recreate", askOrder.OrderID, err)
			} else {
				// Update our tracking
				mm.updateOrderTracking(askOrder.OrderID, newOrderID, instrument, "sell", askPrice, mm.config.QuoteSize)
				replacedAsk = true
				debugLog("Replaced ask order %s with %s for %s @ %s", askOrder.OrderID, newOrderID, instrument, askPrice)
			}
		}
		
		// If we successfully replaced both orders, we're done
		if replacedBid && replacedAsk {
			return nil
		}
		
		// Otherwise, fall back to cancel and recreate for failed replacements
		if !replacedBid && bidOrder != nil {
			mm.cancelOrder(bidOrder.OrderID)
		}
		if !replacedAsk && askOrder != nil {
			mm.cancelOrder(askOrder.OrderID)
		}
		
		// Place new orders for any that weren't replaced
		if !replacedBid {
			if err := mm.placeSingleQuote(instrument, "buy", bidPrice); err != nil {
				log.Printf("Failed to place bid order: %v", err)
			}
		}
		if !replacedAsk {
			if err := mm.placeSingleQuote(instrument, "sell", askPrice); err != nil {
				log.Printf("Failed to place ask order: %v", err)
			}
		}
	} else {
		// No existing orders, place new ones
		if err := mm.placeQuotes(instrument, bidPrice, askPrice); err != nil {
			return fmt.Errorf("failed to place quotes: %w", err)
		}
	}
	
	return nil
}

// calculateQuotes calculates bid and ask prices based on current market
func (mm *MarketMaker) calculateQuotes(ticker *TickerUpdate, orderBook *MarketMakerOrderBook) (bidPrice, askPrice decimal.Decimal) {
	// Calculate mid price
	midPrice := ticker.BestBid.Add(ticker.BestAsk).Div(decimal.NewFromInt(2))
	
	// Determine reference prices based on orderbook or ticker
	referenceBid := ticker.BestBid
	referenceAsk := ticker.BestAsk
	
	// If we have orderbook data and reference size is set, find the best bid/ask that meets size requirement
	if orderBook != nil && mm.config.ImprovementReferenceSize.GreaterThan(decimal.Zero) {
		// Find the best bid with sufficient size
		foundBid := false
		for _, bid := range orderBook.Bids {
			if bid.Size.GreaterThanOrEqual(mm.config.ImprovementReferenceSize) {
				referenceBid = bid.Price
				foundBid = true
				break
			}
		}
		
		// Find the best ask with sufficient size
		foundAsk := false
		for _, ask := range orderBook.Asks {
			if ask.Size.GreaterThanOrEqual(mm.config.ImprovementReferenceSize) {
				referenceAsk = ask.Price
				foundAsk = true
				break
			}
		}
		
		// If we couldn't find sufficient size on either side, fall back to mid price with spread
		if !foundBid || !foundAsk {
			spreadAmount := midPrice.Mul(decimal.NewFromInt(int64(mm.config.SpreadBps)).Div(decimal.NewFromInt(10000)))
			if !foundBid {
				referenceBid = midPrice.Sub(spreadAmount.Div(decimal.NewFromInt(2)))
			}
			if !foundAsk {
				referenceAsk = midPrice.Add(spreadAmount.Div(decimal.NewFromInt(2)))
			}
		}
	}
	
	// Our quotes: improve the market
	// Bid: reference bid + improvement (or configured spread)
	// Ask: reference ask - improvement (or configured spread)
	improvementAmount := mm.config.Improvement
	
	bidPrice = referenceBid.Add(improvementAmount)
	askPrice = referenceAsk.Sub(improvementAmount)
	
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
	
	debugLog("Placed quotes for %s: Bid %.4f, Ask %.4f", instrument, bidPrice, askPrice)
	
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

// updateOrderTracking updates our internal tracking when an order is replaced
func (mm *MarketMaker) updateOrderTracking(oldOrderID, newOrderID, instrument, side string, price, amount decimal.Decimal) {
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	// Remove old order
	delete(mm.activeOrders, oldOrderID)
	
	// Add new order
	newOrder := &MarketMakerOrder{
		OrderID:    newOrderID,
		Instrument: instrument,
		Side:       side,
		Price:      price,
		Amount:     amount,
		Status:     "open",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	
	mm.activeOrders[newOrderID] = newOrder
	
	// Update instrument tracking
	if mm.ordersByInstrument[instrument] == nil {
		mm.ordersByInstrument[instrument] = make(map[string]*MarketMakerOrder)
	}
	mm.ordersByInstrument[instrument][side] = newOrder
}

// cancelOrder cancels a single order
func (mm *MarketMaker) cancelOrder(orderID string) {
	if err := mm.exchange.CancelOrder(orderID); err != nil {
		log.Printf("Failed to cancel order %s: %v", orderID, err)
	} else {
		mm.mu.Lock()
		if order, exists := mm.activeOrders[orderID]; exists {
			delete(mm.activeOrders, orderID)
			if instrumentOrders, exists := mm.ordersByInstrument[order.Instrument]; exists {
				delete(instrumentOrders, order.Side)
				if len(instrumentOrders) == 0 {
					delete(mm.ordersByInstrument, order.Instrument)
				}
			}
			mm.stats.OrdersCancelled++
		}
		mm.mu.Unlock()
	}
}

// placeSingleQuote places a single buy or sell order
func (mm *MarketMaker) placeSingleQuote(instrument, side string, price decimal.Decimal) error {
	orderID, err := mm.exchange.PlaceLimitOrder(instrument, side, price, mm.config.QuoteSize)
	if err != nil {
		return err
	}
	
	// Track the order
	mm.mu.Lock()
	order := &MarketMakerOrder{
		OrderID:    orderID,
		Instrument: instrument,
		Side:       side,
		Price:      price,
		Amount:     mm.config.QuoteSize,
		Status:     "open",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	mm.activeOrders[orderID] = order
	
	if mm.ordersByInstrument[instrument] == nil {
		mm.ordersByInstrument[instrument] = make(map[string]*MarketMakerOrder)
	}
	mm.ordersByInstrument[instrument][side] = order
	
	mm.stats.OrdersPlaced++
	mm.mu.Unlock()
	
	debugLog("Placed %s order %s for %s @ %s", side, orderID, instrument, price)
	return nil
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
			
			// Count active orders
			activeCount := 0
			for _, orders := range mm.ordersByInstrument {
				if len(orders) > 0 {
					activeCount++
				}
			}
			
			// Concise stats output
			log.Printf("Stats: Orders=%d/%d/%d (placed/cancelled/filled), Active=%d/%d instruments, Uptime=%ds",
				mm.stats.OrdersPlaced,
				mm.stats.OrdersCancelled,
				mm.stats.OrdersFilled,
				activeCount,
				len(mm.config.Instruments),
				mm.stats.UptimeSeconds)
			
			// Detailed order state in debug mode only
			if debugMode {
				for instrument, orders := range mm.ordersByInstrument {
					if len(orders) > 0 {
						var bidPrice, askPrice string = "none", "none"
						if bidOrder, hasBid := orders["buy"]; hasBid {
							bidPrice = bidOrder.Price.String()
						}
						if askOrder, hasAsk := orders["sell"]; hasAsk {
							askPrice = askOrder.Price.String()
						}
						debugLog("  %s: bid=%s, ask=%s", instrument, bidPrice, askPrice)
					}
				}
			}
			mm.mu.RUnlock()
		}
	}
}