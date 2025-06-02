package main

import (
	"context"
	"fmt"
	"log"
	"strings"
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
	
	// Per-instrument update locks to prevent concurrent updates
	updateLocks map[string]*sync.Mutex
	
	// Track orders that consistently fail to cancel (likely don't exist)
	failedCancelAttempts map[string]int
	
	// Track last update time per instrument to prevent duplicate updates
	lastUpdateTime map[string]time.Time
}

// NewMarketMaker creates a new market maker instance
func NewMarketMaker(config *MarketMakerConfig, exchange MarketMakerExchange) *MarketMaker {
	ctx, cancel := context.WithCancel(context.Background())
	
	// Initialize update locks for each instrument
	updateLocks := make(map[string]*sync.Mutex)
	for _, instrument := range config.Instruments {
		updateLocks[instrument] = &sync.Mutex{}
	}
	
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
		updateLocks:        updateLocks,
		failedCancelAttempts: make(map[string]int),
		lastUpdateTime:     make(map[string]time.Time),
	}
}

// Start begins market making
func (mm *MarketMaker) Start() error {
	mode := "two-sided (1 buy + 1 sell)"
	if mm.config.BidOnly {
		mode = "bid-only"
	} else if mm.config.AskOnly {
		mode = "ask-only"
	}
	log.Printf("Starting market maker: %d instruments, %s per instrument", len(mm.config.Instruments), mode)
	
	// Clear any stale state
	mm.mu.Lock()
	mm.activeOrders = make(map[string]*MarketMakerOrder)
	mm.ordersByInstrument = make(map[string]map[string]*MarketMakerOrder)
	mm.mu.Unlock()
	
	// Load existing positions
	if err := mm.loadPositions(); err != nil {
		return fmt.Errorf("failed to load positions: %w", err)
	}
	
	// Load existing orders
	if err := mm.loadActiveOrders(); err != nil {
		return fmt.Errorf("failed to load active orders: %w", err)
	}
	
	// Cancel all existing orders on startup
	log.Printf("Cancelling all existing orders on startup...")
	cancelCount := 0
	mm.mu.RLock()
	ordersCopy := make(map[string]*MarketMakerOrder)
	for id, order := range mm.activeOrders {
		ordersCopy[id] = order
	}
	mm.mu.RUnlock()
	
	for _, order := range ordersCopy {
		if err := mm.exchange.CancelOrder(order.OrderID); err != nil {
			log.Printf("Failed to cancel order %s: %v", order.OrderID, err)
		} else {
			cancelCount++
			mm.mu.Lock()
			delete(mm.activeOrders, order.OrderID)
			if mm.ordersByInstrument[order.Instrument] != nil {
				delete(mm.ordersByInstrument[order.Instrument], order.Side)
			}
			mm.mu.Unlock()
		}
	}
	log.Printf("Cancelled %d orders on startup", cancelCount)
	
	// Subscribe to ticker updates
	tickerChan, err := mm.exchange.SubscribeTickers(mm.ctx, mm.config.Instruments)
	if err != nil {
		return fmt.Errorf("failed to subscribe to tickers: %w", err)
	}
	
	// Subscribe to orderbook updates for each instrument
	if subscriber, ok := mm.exchange.(interface{ SubscribeOrderBook(string) error }); ok {
		for _, instrument := range mm.config.Instruments {
			if err := subscriber.SubscribeOrderBook(instrument); err != nil {
				log.Printf("Failed to subscribe to orderbook for %s: %v", instrument, err)
				// Continue anyway - we can fall back to ticker data
			} else {
				log.Printf("Subscribed to orderbook for %s", instrument)
			}
		}
		// Give orderbook subscriptions time to receive initial data
		time.Sleep(2 * time.Second)
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
	
	// Run an immediate reconciliation to clean up any stale state
	mm.reconcileOrders()
	
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
	// Prevent concurrent updates for the same instrument
	if lock, exists := mm.updateLocks[instrument]; exists {
		lock.Lock()
		defer lock.Unlock()
	}
	
	// Check if we recently updated this instrument (within 2 seconds to avoid rate limits)
	mm.mu.RLock()
	lastUpdate, exists := mm.lastUpdateTime[instrument]
	mm.mu.RUnlock()
	
	if exists && time.Since(lastUpdate) < 2*time.Second {
		// Skip this update to prevent rate limiting
		return nil
	}
	
	// Update the last update time
	mm.mu.Lock()
	mm.lastUpdateTime[instrument] = time.Now()
	mm.mu.Unlock()
	mm.mu.RLock()
	ticker, exists := mm.latestTickers[instrument]
	mm.mu.RUnlock()
	
	if !exists || ticker == nil {
		return fmt.Errorf("no ticker data for %s", instrument)
	}
	
	// Skip if we don't have valid price data yet
	if ticker.BestBid.IsZero() && ticker.BestAsk.IsZero() && ticker.MarkPrice.IsZero() {
		log.Printf("No valid price data for %s yet (BestBid=0, BestAsk=0, MarkPrice=0), skipping quote update", instrument)
		return nil
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
	
	// Get fresh order state from exchange to avoid phantom orders
	openOrders, err := mm.exchange.GetOpenOrders()
	if err != nil {
		log.Printf("Failed to get open orders: %v", err)
		return err
	}
	
	// Update our tracking with real orders
	mm.mu.Lock()
	// Clear existing tracking for this instrument
	if mm.ordersByInstrument[instrument] != nil {
		for _, order := range mm.ordersByInstrument[instrument] {
			if order != nil {
				delete(mm.activeOrders, order.OrderID)
			}
		}
	}
	mm.ordersByInstrument[instrument] = make(map[string]*MarketMakerOrder)
	
	// Add real orders to tracking
	for i := range openOrders {
		order := &openOrders[i]
		if order.Instrument == instrument {
			mm.activeOrders[order.OrderID] = order
			mm.ordersByInstrument[instrument][order.Side] = order
		}
	}
	mm.mu.Unlock()
	
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
		
		// Cancel unwanted orders in one-sided mode
		if mm.config.AskOnly && bidOrder != nil {
			log.Printf("Cancelling bid order %s (ask-only mode)", bidOrder.OrderID)
			mm.cancelOrder(bidOrder.OrderID)
			bidOrder = nil
		}
		if mm.config.BidOnly && askOrder != nil {
			log.Printf("Cancelling ask order %s (bid-only mode)", askOrder.OrderID)
			mm.cancelOrder(askOrder.OrderID)
			askOrder = nil
		}
		
		// Replace existing orders if they exist
		var replacedBid, replacedAsk bool
		
		// Handle one-sided quoting
		if mm.config.AskOnly {
			replacedBid = true // Skip bid side
		}
		if mm.config.BidOnly {
			replacedAsk = true // Skip ask side
		}
		
		// TEMPORARY: Disable ReplaceOrder since it seems broken on Derive
		const useReplaceOrder = false
		
		if bidOrder != nil && bidOrder.OrderID != "" && !mm.config.AskOnly {
			// Skip if order is already at target price and was created recently
			if bidOrder.Price.Equal(bidPrice) && time.Since(bidOrder.CreatedAt) < 5*time.Second {
				debugLog("Skipping bid update - order %s already at target price", bidOrder.OrderID)
				replacedBid = true
			} else if useReplaceOrder {
				newOrderID, err := mm.exchange.ReplaceOrder(bidOrder.OrderID, instrument, "buy", bidPrice, mm.config.QuoteSize)
				if err != nil {
					log.Printf("Failed to replace bid order %s: %v, will cancel and recreate", bidOrder.OrderID, err)
			} else if newOrderID != "" {
				// Update our tracking only if we got a valid order ID
				mm.updateOrderTracking(bidOrder.OrderID, newOrderID, instrument, "buy", bidPrice, mm.config.QuoteSize)
				replacedBid = true
				debugLog("Replaced bid order %s with %s for %s @ %s", bidOrder.OrderID, newOrderID, instrument, bidPrice)
				// Verify the replacement worked
				if !mm.verifyOrderExists(newOrderID) {
					log.Printf("WARNING: Replacement bid order %s not found, will recreate", newOrderID)
					replacedBid = false
					mm.mu.Lock()
					delete(mm.activeOrders, newOrderID)
					delete(mm.ordersByInstrument[instrument], "buy")
					mm.mu.Unlock()
				}
			} else {
				log.Printf("Failed to replace bid order %s: got empty order ID, will cancel and recreate", bidOrder.OrderID)
			}
			} else {
				// Just cancel - we'll recreate below
				if !mm.cancelOrder(bidOrder.OrderID) {
					// Failed to cancel, don't create a new order to avoid duplicates
					log.Printf("Failed to cancel bid order %s, skipping recreation to avoid duplicates", bidOrder.OrderID)
					replacedBid = true // Mark as "replaced" to skip recreation
				}
			}
		}
		
		if askOrder != nil && askOrder.OrderID != "" && !mm.config.BidOnly {
			// Skip if order is already at target price and was created recently
			if askOrder.Price.Equal(askPrice) && time.Since(askOrder.CreatedAt) < 5*time.Second {
				debugLog("Skipping ask update - order %s already at target price", askOrder.OrderID)
				replacedAsk = true
			} else if useReplaceOrder {
				newOrderID, err := mm.exchange.ReplaceOrder(askOrder.OrderID, instrument, "sell", askPrice, mm.config.QuoteSize)
				if err != nil {
					log.Printf("Failed to replace ask order %s: %v, will cancel and recreate", askOrder.OrderID, err)
			} else if newOrderID != "" {
				// Update our tracking only if we got a valid order ID
				mm.updateOrderTracking(askOrder.OrderID, newOrderID, instrument, "sell", askPrice, mm.config.QuoteSize)
				replacedAsk = true
				debugLog("Replaced ask order %s with %s for %s @ %s", askOrder.OrderID, newOrderID, instrument, askPrice)
				// Verify the replacement worked
				if !mm.verifyOrderExists(newOrderID) {
					log.Printf("WARNING: Replacement ask order %s not found, will recreate", newOrderID)
					replacedAsk = false
					mm.mu.Lock()
					delete(mm.activeOrders, newOrderID)
					delete(mm.ordersByInstrument[instrument], "sell")
					mm.mu.Unlock()
				}
			} else {
				log.Printf("Failed to replace ask order %s: got empty order ID, will cancel and recreate", askOrder.OrderID)
			}
			} else {
				// Just cancel - we'll recreate below
				if !mm.cancelOrder(askOrder.OrderID) {
					// Failed to cancel, don't create a new order to avoid duplicates
					log.Printf("Failed to cancel ask order %s, skipping recreation to avoid duplicates", askOrder.OrderID)
					replacedAsk = true // Mark as "replaced" to skip recreation
				}
			}
		}
		
		// If we successfully replaced both orders, we're done
		if replacedBid && replacedAsk {
			return nil
		}
		
		// Otherwise, fall back to cancel and recreate for failed replacements
		if !replacedBid && bidOrder != nil && bidOrder.OrderID != "" {
			if !mm.cancelOrder(bidOrder.OrderID) {
				// Failed to cancel, skip recreation
				log.Printf("Failed to cancel bid order %s in fallback, skipping recreation", bidOrder.OrderID)
				replacedBid = true
			}
		}
		if !replacedAsk && askOrder != nil && askOrder.OrderID != "" {
			if !mm.cancelOrder(askOrder.OrderID) {
				// Failed to cancel, skip recreation
				log.Printf("Failed to cancel ask order %s in fallback, skipping recreation", askOrder.OrderID)
				replacedAsk = true
			}
		}
		
		// If we had any failed cancellations, wait a bit and verify order state
		if (!replacedBid && bidOrder != nil) || (!replacedAsk && askOrder != nil) {
			time.Sleep(1 * time.Second) // Wait for any pending operations to settle
			
			// Re-verify order state from exchange
			openOrders, err := mm.exchange.GetOpenOrders()
			if err == nil {
				// Check if the orders we tried to cancel still exist
				for _, order := range openOrders {
					if order.Instrument == instrument {
						if order.Side == "buy" && bidOrder != nil && order.OrderID == bidOrder.OrderID {
							log.Printf("Bid order %s still exists after cancel attempt, skipping new order", order.OrderID)
							replacedBid = true
						}
						if order.Side == "sell" && askOrder != nil && order.OrderID == askOrder.OrderID {
							log.Printf("Ask order %s still exists after cancel attempt, skipping new order", order.OrderID)
							replacedAsk = true
						}
					}
				}
			}
		}
		
		// Place new orders for any that weren't replaced (respecting one-sided flags)
		if !replacedBid && !mm.config.AskOnly {
			if err := mm.placeSingleQuote(instrument, "buy", bidPrice); err != nil {
				log.Printf("Failed to place bid order: %v", err)
			}
		}
		if !replacedAsk && !mm.config.BidOnly {
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
	// Calculate mid price - use mark price as fallback if orderbook is empty
	var midPrice decimal.Decimal
	if ticker.BestBid.IsZero() || ticker.BestAsk.IsZero() {
		// Use mark price if available, otherwise use a default
		if !ticker.MarkPrice.IsZero() {
			midPrice = ticker.MarkPrice
		} else {
			// This should rarely happen - log it
			log.Printf("WARNING: No valid price data for %s (BestBid=%s, BestAsk=%s, MarkPrice=%s)", 
				ticker.Instrument, ticker.BestBid.String(), ticker.BestAsk.String(), ticker.MarkPrice.String())
			// Use a reasonable default for options (e.g., $1)
			midPrice = decimal.NewFromFloat(1.0)
		}
	} else {
		midPrice = ticker.BestBid.Add(ticker.BestAsk).Div(decimal.NewFromInt(2))
	}
	
	// Determine reference prices based on orderbook or ticker
	referenceBid := ticker.BestBid
	referenceAsk := ticker.BestAsk
	
	// If best bid/ask are zero, create synthetic prices from mid price
	if referenceBid.IsZero() || referenceAsk.IsZero() {
		// Use configured spread to create synthetic bid/ask
		spreadAmount := midPrice.Mul(decimal.NewFromInt(int64(mm.config.SpreadBps)).Div(decimal.NewFromInt(10000)))
		if referenceBid.IsZero() {
			referenceBid = midPrice.Sub(spreadAmount.Div(decimal.NewFromInt(2)))
		}
		if referenceAsk.IsZero() {
			referenceAsk = midPrice.Add(spreadAmount.Div(decimal.NewFromInt(2)))
		}
	}
	
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

// placeQuotes places bid and ask orders (respecting one-sided flags)
func (mm *MarketMaker) placeQuotes(instrument string, bidPrice, askPrice decimal.Decimal) error {
	var bidOrderID, askOrderID string
	var err error
	ordersPlaced := 0
	
	// Place bid order if not ask-only
	if !mm.config.AskOnly {
		bidOrderID, err = mm.exchange.PlaceLimitOrder(instrument, "buy", bidPrice, mm.config.QuoteSize)
		if err != nil {
			return fmt.Errorf("failed to place bid order: %w", err)
		}
		ordersPlaced++
	}
	
	// Place ask order if not bid-only
	if !mm.config.BidOnly {
		askOrderID, err = mm.exchange.PlaceLimitOrder(instrument, "sell", askPrice, mm.config.QuoteSize)
		if err != nil {
			// Cancel the bid order if we placed one but couldn't place the ask
			if bidOrderID != "" {
				mm.exchange.CancelOrder(bidOrderID)
			}
			return fmt.Errorf("failed to place ask order: %w", err)
		}
		ordersPlaced++
	}
	
	// Track orders
	mm.mu.Lock()
	if bidOrderID != "" {
		mm.trackOrder(bidOrderID, instrument, "buy", bidPrice, mm.config.QuoteSize)
	}
	if askOrderID != "" {
		mm.trackOrder(askOrderID, instrument, "sell", askPrice, mm.config.QuoteSize)
	}
	mm.stats.OrdersPlaced += int64(ordersPlaced)
	mm.mu.Unlock()
	
	// Log what was actually placed
	if mm.config.BidOnly {
		log.Printf("Placed bid order for %s @ %s", instrument, bidPrice.String())
	} else if mm.config.AskOnly {
		log.Printf("Placed ask order for %s @ %s", instrument, askPrice.String())
	} else {
		log.Printf("Placed quotes for %s: Bid %s, Ask %s", instrument, bidPrice.String(), askPrice.String())
	}
	
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

// cancelOrder cancels a single order and returns whether it was successfully cancelled
func (mm *MarketMaker) cancelOrder(orderID string) bool {
	// Check if we've already tried to cancel this order many times
	mm.mu.Lock()
	if attempts, exists := mm.failedCancelAttempts[orderID]; exists && attempts >= 3 {
		// This order has consistently failed to cancel, treat as already gone
		log.Printf("Order %s has failed to cancel %d times, treating as non-existent", orderID, attempts)
		// Clean up tracking
		if order, exists := mm.activeOrders[orderID]; exists {
			delete(mm.activeOrders, orderID)
			if instrumentOrders, exists := mm.ordersByInstrument[order.Instrument]; exists {
				delete(instrumentOrders, order.Side)
				if len(instrumentOrders) == 0 {
					delete(mm.ordersByInstrument, order.Instrument)
				}
			}
		}
		delete(mm.failedCancelAttempts, orderID)
		mm.mu.Unlock()
		return true
	}
	mm.mu.Unlock()
	
	// Try to cancel with retries for internal errors
	var lastErr error
	for retries := 0; retries < 3; retries++ {
		if retries > 0 {
			time.Sleep(500 * time.Millisecond) // Wait before retry
		}
		
		err := mm.exchange.CancelOrder(orderID)
		if err == nil {
			// Success - remove from tracking
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
			delete(mm.failedCancelAttempts, orderID)
			mm.mu.Unlock()
			return true
		}
		
		lastErr = err
		// If it's an "order not found" error, treat as success (order already gone)
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "Order does not exist") {
			log.Printf("Order %s already cancelled or doesn't exist", orderID)
			// Remove from tracking
			mm.mu.Lock()
			if order, exists := mm.activeOrders[orderID]; exists {
				delete(mm.activeOrders, orderID)
				if instrumentOrders, exists := mm.ordersByInstrument[order.Instrument]; exists {
					delete(instrumentOrders, order.Side)
					if len(instrumentOrders) == 0 {
						delete(mm.ordersByInstrument, order.Instrument)
					}
				}
			}
			delete(mm.failedCancelAttempts, orderID)
			mm.mu.Unlock()
			return true
		}
		
		// For "Internal error", retry but only log on first retry
		if strings.Contains(err.Error(), "Internal error") {
			if retries == 0 {
				log.Printf("Internal error cancelling order %s, will retry", orderID)
			}
			continue
		}
		
		// For other errors, don't retry
		break
	}
	
	// Track failed attempts
	mm.mu.Lock()
	mm.failedCancelAttempts[orderID]++
	attempts := mm.failedCancelAttempts[orderID]
	mm.mu.Unlock()
	
	// If this is the 3rd+ failed attempt, treat as non-existent
	if attempts >= 3 {
		log.Printf("Order %s failed to cancel %d times (likely doesn't exist), treating as cancelled", orderID, attempts)
		mm.mu.Lock()
		if order, exists := mm.activeOrders[orderID]; exists {
			delete(mm.activeOrders, orderID)
			if instrumentOrders, exists := mm.ordersByInstrument[order.Instrument]; exists {
				delete(instrumentOrders, order.Side)
				if len(instrumentOrders) == 0 {
					delete(mm.ordersByInstrument, order.Instrument)
				}
			}
		}
		delete(mm.failedCancelAttempts, orderID)
		mm.mu.Unlock()
		return true
	}
	
	log.Printf("Failed to cancel order %s (attempt %d): %v", orderID, attempts, lastErr)
	return false
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
	
	log.Printf("Placed %s order %s for %s @ %s", side, orderID, instrument, price)
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
			totalOrders := 0
			for _, orders := range mm.ordersByInstrument {
				if len(orders) > 0 {
					activeCount++
					totalOrders += len(orders)
				}
			}
			
			// Quick consistency check
			if totalOrders != len(mm.activeOrders) {
				log.Printf("WARNING: Order tracking inconsistency detected: %d orders in activeOrders, %d in ordersByInstrument", 
					len(mm.activeOrders), totalOrders)
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

// verifyOrderExists checks if an order actually exists on the exchange
func (mm *MarketMaker) verifyOrderExists(orderID string) bool {
	// Get open orders from exchange
	orders, err := mm.exchange.GetOpenOrders()
	if err != nil {
		debugLog("Failed to verify order %s: %v", orderID, err)
		return true // Assume it exists if we can't check
	}
	
	for _, order := range orders {
		if order.OrderID == orderID {
			return true
		}
	}
	return false
}


// reconcileOrdersForInstrument reconciles orders for a specific instrument
func (mm *MarketMaker) reconcileOrdersForInstrument(instrument string) {
	// Get all open orders from exchange
	openOrders, err := mm.exchange.GetOpenOrders()
	if err != nil {
		debugLog("Failed to get open orders for reconciliation: %v", err)
		return
	}
	
	// Filter to just this instrument
	var instrumentOrders []MarketMakerOrder
	for _, order := range openOrders {
		if order.Instrument == instrument {
			instrumentOrders = append(instrumentOrders, order)
		}
	}
	
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	// Check our tracked orders for this instrument
	trackedOrders := mm.ordersByInstrument[instrument]
	if trackedOrders == nil {
		trackedOrders = make(map[string]*MarketMakerOrder)
	}
	
	// Build a map of actual order IDs
	actualOrders := make(map[string]bool)
	for _, order := range instrumentOrders {
		actualOrders[order.OrderID] = true
	}
	
	// Remove any tracked orders that don't actually exist
	for side, order := range trackedOrders {
		if order != nil && !actualOrders[order.OrderID] {
			debugLog("Removing phantom %s order %s for %s", side, order.OrderID, instrument)
			delete(mm.activeOrders, order.OrderID)
			delete(trackedOrders, side)
		}
	}
	
	// Add any real orders we're not tracking
	for _, order := range instrumentOrders {
		if _, tracked := mm.activeOrders[order.OrderID]; !tracked {
			log.Printf("Found untracked order %s for %s, adding to tracking", order.OrderID, instrument)
			orderCopy := order
			mm.activeOrders[order.OrderID] = &orderCopy
			if mm.ordersByInstrument[instrument] == nil {
				mm.ordersByInstrument[instrument] = make(map[string]*MarketMakerOrder)
			}
			mm.ordersByInstrument[instrument][order.Side] = &orderCopy
		}
	}
}

// reconcileOrders finds and cancels any orders not being tracked
func (mm *MarketMaker) reconcileOrders() {
	// Get all open orders from exchange
	openOrders, err := mm.exchange.GetOpenOrders()
	if err != nil {
		log.Printf("Failed to get open orders for reconciliation: %v", err)
		return
	}
	
	mm.mu.RLock()
	// Create a map of our tracked order IDs
	trackedOrders := make(map[string]bool)
	for orderID := range mm.activeOrders {
		trackedOrders[orderID] = true
	}
	mm.mu.RUnlock()
	
	// Find orphaned orders
	orphanedCount := 0
	for _, order := range openOrders {
		if !trackedOrders[order.OrderID] {
			orphanedCount++
			log.Printf("Found orphaned order %s for %s, cancelling", order.OrderID, order.Instrument)
			if err := mm.exchange.CancelOrder(order.OrderID); err != nil {
				log.Printf("Failed to cancel orphaned order %s: %v", order.OrderID, err)
			}
		}
	}
	
	if orphanedCount > 0 {
		log.Printf("Cancelled %d orphaned orders", orphanedCount)
		mm.stats.OrdersCancelled += int64(orphanedCount)
	}
	
	// Also verify our tracked orders still exist
	mm.mu.Lock()
	defer mm.mu.Unlock()
	
	for orderID, order := range mm.activeOrders {
		found := false
		for _, openOrder := range openOrders {
			if openOrder.OrderID == orderID {
				found = true
				break
			}
		}
		if !found {
			log.Printf("Tracked order %s no longer exists on exchange, removing from tracking", orderID)
			delete(mm.activeOrders, orderID)
			if orders, ok := mm.ordersByInstrument[order.Instrument]; ok {
				delete(orders, order.Side)
				if len(orders) == 0 {
					delete(mm.ordersByInstrument, order.Instrument)
				}
			}
		}
	}
}