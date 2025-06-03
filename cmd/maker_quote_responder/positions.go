package main

import (
    "log"
    "github.com/shopspring/decimal"
)

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

// updatePosition updates position after a fill
func (mm *MarketMaker) updatePosition(instrument string, side string, amount decimal.Decimal) {
    mm.mu.Lock()
    defer mm.mu.Unlock()
    
    if side == "sell" {
        amount = amount.Neg()
    }
    
    currentPosition := mm.positions[instrument]
    mm.positions[instrument] = currentPosition.Add(amount)
    
    log.Printf("Updated position for %s: %s -> %s", 
        instrument, 
        currentPosition.String(), 
        mm.positions[instrument].String())
}

// getNetPosition returns the net position for an instrument
func (mm *MarketMaker) getNetPosition(instrument string) decimal.Decimal {
    mm.mu.RLock()
    defer mm.mu.RUnlock()
    
    return mm.positions[instrument]
}

// getTotalExposure returns total absolute exposure across all instruments
func (mm *MarketMaker) getTotalExposure() decimal.Decimal {
    mm.mu.RLock()
    defer mm.mu.RUnlock()
    
    totalExposure := decimal.Zero
    for _, pos := range mm.positions {
        totalExposure = totalExposure.Add(pos.Abs())
    }
    
    return totalExposure
}