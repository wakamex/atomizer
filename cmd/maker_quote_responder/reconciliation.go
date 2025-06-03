package main

import (
    "log"
)

// reconcileOrders finds and cancels any orders not being tracked
func (mm *MarketMaker) reconcileOrders() {
    openOrders, err := mm.exchange.GetOpenOrders()
    if err != nil {
        log.Printf("Failed to get open orders for reconciliation: %v", err)
        return
    }
    
    mm.mu.RLock()
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
        mm.mu.Lock()
        mm.stats.OrdersCancelled += int64(orphanedCount)
        mm.mu.Unlock()
    }
    
    // Verify tracked orders still exist
    mm.verifyTrackedOrders(openOrders)
}

// verifyTrackedOrders removes tracked orders that no longer exist
func (mm *MarketMaker) verifyTrackedOrders(openOrders []MarketMakerOrder) {
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

// reconcileOrdersForInstrument reconciles orders for a specific instrument
func (mm *MarketMaker) reconcileOrdersForInstrument(instrument string) {
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
    
    trackedOrders := mm.ordersByInstrument[instrument]
    if trackedOrders == nil {
        trackedOrders = make(map[string]*MarketMakerOrder)
    }
    
    // Build map of actual order IDs
    actualOrders := make(map[string]bool)
    for _, order := range instrumentOrders {
        actualOrders[order.OrderID] = true
    }
    
    // Remove phantom orders
    for side, order := range trackedOrders {
        if order != nil && !actualOrders[order.OrderID] {
            debugLog("Removing phantom %s order %s for %s", side, order.OrderID, instrument)
            delete(mm.activeOrders, order.OrderID)
            delete(trackedOrders, side)
        }
    }
    
    // Add untracked orders
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