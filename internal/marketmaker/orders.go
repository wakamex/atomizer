package marketmaker

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wakamex/atomizer/internal/types"
)

// updateOrCreateOrders handles the order update/creation logic
func (mm *MarketMaker) updateOrCreateOrders(instrument string, bidPrice, askPrice decimal.Decimal) error {
	// Get fresh order state
	openOrders, err := mm.exchange.GetOpenOrders()
	if err != nil {
		log.Printf("Failed to get open orders: %v", err)
		return err
	}

	// Update tracking with real orders
	mm.syncOrderTracking(instrument, openOrders)

	// Get existing orders
	mm.mu.RLock()
	existingOrders := mm.ordersByInstrument[instrument]
	mm.mu.RUnlock()

	if existingOrders != nil && len(existingOrders) > 0 {
		log.Printf("Found %d existing orders for %s, will update", len(existingOrders), instrument)
		return mm.replaceExistingOrders(instrument, bidPrice, askPrice, existingOrders)
	}

	// No existing orders, place new ones
	log.Printf("No existing orders for %s, placing new quotes", instrument)
	return mm.placeQuotes(instrument, bidPrice, askPrice)
}

// syncOrderTracking updates our tracking with real orders from exchange
func (mm *MarketMaker) syncOrderTracking(instrument string, openOrders []types.MarketMakerOrder) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	// First, clear ALL tracked orders (they might have been filled/cancelled)
	for _, instrumentOrders := range mm.ordersByInstrument {
		for _, order := range instrumentOrders {
			if order != nil {
				delete(mm.activeOrders, order.OrderID)
			}
		}
	}
	
	// Reset all instrument order maps
	for inst := range mm.ordersByInstrument {
		mm.ordersByInstrument[inst] = make(map[string]*types.MarketMakerOrder)
	}

	// Now add ALL orders back from the exchange response
	for i := range openOrders {
		order := &openOrders[i]
		// Make sure we have a map for this instrument
		if mm.ordersByInstrument[order.Instrument] == nil {
			mm.ordersByInstrument[order.Instrument] = make(map[string]*types.MarketMakerOrder)
		}
		
		log.Printf("Tracking order %s for %s: side=%s, price=%s", order.OrderID, order.Instrument, order.Side, order.Price)
		mm.activeOrders[order.OrderID] = order
		mm.ordersByInstrument[order.Instrument][order.Side] = order
	}
}

// replaceExistingOrders handles order replacement logic
func (mm *MarketMaker) replaceExistingOrders(instrument string, bidPrice, askPrice decimal.Decimal, existingOrders map[string]*types.MarketMakerOrder) error {
	var bidOrder, askOrder *types.MarketMakerOrder

	// Find existing orders
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

	var replacedBid, replacedAsk bool

	// Skip sides based on mode
	if mm.config.AskOnly {
		replacedBid = true
	}
	if mm.config.BidOnly {
		replacedAsk = true
	}

	// Handle bid order
	if bidOrder != nil && !mm.config.AskOnly {
		replacedBid = mm.handleOrderUpdate(bidOrder, instrument, "buy", bidPrice)
	}

	// Handle ask order
	if askOrder != nil && !mm.config.BidOnly {
		replacedAsk = mm.handleOrderUpdate(askOrder, instrument, "sell", askPrice)
	}

	// Place new orders for any that weren't replaced
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

	return nil
}

// handleOrderUpdate processes a single order update
func (mm *MarketMaker) handleOrderUpdate(order *types.MarketMakerOrder, instrument, side string, newPrice decimal.Decimal) bool {
	// Skip if order is already at target price and recent
	if order.Price.Equal(newPrice) && time.Since(order.CreatedAt) < 5*time.Second {
		DebugLog("Skipping %s update - order %s already at target price", side, order.OrderID)
		return true
	}

	// Just cancel - we'll recreate
	if !mm.cancelOrder(order.OrderID) {
		log.Printf("Failed to cancel %s order %s, skipping recreation to avoid duplicates", side, order.OrderID)
		return true
	}

	return false
}

// placeQuotes places bid and ask orders concurrently
func (mm *MarketMaker) placeQuotes(instrument string, bidPrice, askPrice decimal.Decimal) error {
	var bidOrderID, askOrderID string
	var bidErr, askErr error
	ordersPlaced := 0

	var wg sync.WaitGroup

	// Place bid order
	if !mm.config.AskOnly {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bidOrderID, bidErr = mm.exchange.PlaceLimitOrder(instrument, "buy", bidPrice, mm.config.QuoteSize)
		}()
	}

	// Place ask order
	if !mm.config.BidOnly {
		wg.Add(1)
		go func() {
			defer wg.Done()
			askOrderID, askErr = mm.exchange.PlaceLimitOrder(instrument, "sell", askPrice, mm.config.QuoteSize)
		}()
	}

	wg.Wait()

	// Handle errors
	if bidErr != nil && askErr != nil {
		return fmt.Errorf("failed to place both orders: bid error: %w, ask error: %v", bidErr, askErr)
	}

	if bidErr != nil && askOrderID != "" {
		mm.exchange.CancelOrder(askOrderID)
		return fmt.Errorf("failed to place bid order (cancelled ask): %w", bidErr)
	}

	if askErr != nil && bidOrderID != "" {
		mm.exchange.CancelOrder(bidOrderID)
		return fmt.Errorf("failed to place ask order (cancelled bid): %w", askErr)
	}

	// Track successfully placed orders
	mm.mu.Lock()
	if bidOrderID != "" && bidErr == nil {
		mm.trackOrder(bidOrderID, instrument, "buy", bidPrice, mm.config.QuoteSize)
		ordersPlaced++
	}
	if askOrderID != "" && askErr == nil {
		mm.trackOrder(askOrderID, instrument, "sell", askPrice, mm.config.QuoteSize)
		ordersPlaced++
	}
	mm.stats.OrdersPlaced += int64(ordersPlaced)
	mm.mu.Unlock()

	// Log placement
	if mm.config.BidOnly {
		log.Printf("Placed bid order for %s @ %s", instrument, bidPrice.String())
	} else if mm.config.AskOnly {
		log.Printf("Placed ask order for %s @ %s", instrument, askPrice.String())
	} else {
		log.Printf("Placed quotes for %s: Bid %s, Ask %s", instrument, bidPrice.String(), askPrice.String())
	}

	return nil
}

// placeSingleQuote places a single buy or sell order
func (mm *MarketMaker) placeSingleQuote(instrument, side string, price decimal.Decimal) error {
	orderID, err := mm.exchange.PlaceLimitOrder(instrument, side, price, mm.config.QuoteSize)
	if err != nil {
		return err
	}

	mm.mu.Lock()
	mm.trackOrder(orderID, instrument, side, price, mm.config.QuoteSize)
	mm.stats.OrdersPlaced++
	mm.mu.Unlock()

	log.Printf("Placed %s order %s for %s @ %s", side, orderID, instrument, price)
	return nil
}

// trackOrder adds an order to tracking
func (mm *MarketMaker) trackOrder(orderID, instrument, side string, price, amount decimal.Decimal) {
	order := &types.MarketMakerOrder{
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
		mm.ordersByInstrument[instrument] = make(map[string]*types.MarketMakerOrder)
	}
	mm.ordersByInstrument[instrument][side] = order
}

// cancelOrder cancels a single order
func (mm *MarketMaker) cancelOrder(orderID string) bool {
	// Check failed attempts
	mm.mu.Lock()
	if attempts, exists := mm.failedCancelAttempts[orderID]; exists && attempts >= 3 {
		log.Printf("Order %s has failed to cancel %d times, treating as non-existent", orderID, attempts)
		mm.removeOrderFromTracking(orderID)
		delete(mm.failedCancelAttempts, orderID)
		mm.mu.Unlock()
		return true
	}
	mm.mu.Unlock()

	// Try to cancel with retries
	var lastErr error
	for retries := 0; retries < 3; retries++ {
		if retries > 0 {
			time.Sleep(500 * time.Millisecond)
		}

		err := mm.exchange.CancelOrder(orderID)
		if err == nil {
			mm.mu.Lock()
			mm.removeOrderFromTracking(orderID)
			mm.stats.OrdersCancelled++
			delete(mm.failedCancelAttempts, orderID)
			mm.mu.Unlock()
			return true
		}

		lastErr = err

		// Handle "order not found" as success
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "Order does not exist") {
			log.Printf("Order %s already cancelled or doesn't exist", orderID)
			mm.mu.Lock()
			mm.removeOrderFromTracking(orderID)
			delete(mm.failedCancelAttempts, orderID)
			mm.mu.Unlock()
			return true
		}

		// Retry internal errors
		if strings.Contains(err.Error(), "Internal error") {
			if retries == 0 {
				log.Printf("Internal error cancelling order %s, will retry", orderID)
			}
			continue
		}

		break
	}

	// Track failed attempts
	mm.mu.Lock()
	mm.failedCancelAttempts[orderID]++
	attempts := mm.failedCancelAttempts[orderID]
	mm.mu.Unlock()

	if attempts >= 3 {
		log.Printf("Order %s failed to cancel %d times, treating as cancelled", orderID, attempts)
		mm.mu.Lock()
		mm.removeOrderFromTracking(orderID)
		delete(mm.failedCancelAttempts, orderID)
		mm.mu.Unlock()
		return true
	}

	log.Printf("Failed to cancel order %s (attempt %d): %v", orderID, attempts, lastErr)
	return false
}

// removeOrderFromTracking removes an order from internal tracking (must be called with lock held)
func (mm *MarketMaker) removeOrderFromTracking(orderID string) {
	if order, exists := mm.activeOrders[orderID]; exists {
		delete(mm.activeOrders, orderID)
		if instrumentOrders, exists := mm.ordersByInstrument[order.Instrument]; exists {
			delete(instrumentOrders, order.Side)
			if len(instrumentOrders) == 0 {
				delete(mm.ordersByInstrument, order.Instrument)
			}
		}
	}
}

// cancelAllOrders cancels all active orders
func (mm *MarketMaker) CancelAllOrders() {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	for orderID := range mm.activeOrders {
		if err := mm.exchange.CancelOrder(orderID); err != nil {
			log.Printf("Failed to cancel order %s: %v", orderID, err)
		} else {
			mm.stats.OrdersCancelled++
		}
	}

	mm.activeOrders = make(map[string]*types.MarketMakerOrder)
	mm.ordersByInstrument = make(map[string]map[string]*types.MarketMakerOrder)
}

// loadActiveOrders loads current open orders from exchange
func (mm *MarketMaker) LoadActiveOrders() error {
	orders, err := mm.exchange.GetOpenOrders()
	if err != nil {
		return err
	}

	mm.mu.Lock()
	defer mm.mu.Unlock()

	for _, order := range orders {
		mm.activeOrders[order.OrderID] = &order

		if mm.ordersByInstrument[order.Instrument] == nil {
			mm.ordersByInstrument[order.Instrument] = make(map[string]*types.MarketMakerOrder)
		}
		mm.ordersByInstrument[order.Instrument][order.Side] = &order
	}

	log.Printf("Loaded %d active orders", len(orders))
	return nil
}

// verifyOrderExists checks if an order actually exists on the exchange
func (mm *MarketMaker) verifyOrderExists(orderID string) bool {
	orders, err := mm.exchange.GetOpenOrders()
	if err != nil {
		DebugLog("Failed to verify order %s: %v", orderID, err)
		return true
	}

	for _, order := range orders {
		if order.OrderID == orderID {
			return true
		}
	}
	return false
}
