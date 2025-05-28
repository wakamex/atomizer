package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// Position represents a current position
type Position struct {
	Instrument   string
	Quantity     decimal.Decimal
	AvgPrice     decimal.Decimal
	Delta        decimal.Decimal
	Gamma        decimal.Decimal
	LastUpdated  time.Time
}

// RiskMetrics contains current risk measurements
type RiskMetrics struct {
	TotalDelta       decimal.Decimal
	TotalGamma       decimal.Decimal
	TotalPositions   int
	MaxPositionSize  decimal.Decimal
	UpdatedAt        time.Time
}

// RiskManager manages position and risk limits
type RiskManager struct {
	config           *AppConfig
	positions        map[string]*Position
	maxPositionSize  decimal.Decimal
	maxDeltaExposure decimal.Decimal
	maxGammaExposure decimal.Decimal
	mu               sync.RWMutex
}

// NewRiskManager creates a new risk manager
func NewRiskManager(config *AppConfig) *RiskManager {
	// Default limits - should be configurable
	maxPosition := decimal.NewFromFloat(1000)
	maxDelta := decimal.NewFromFloat(500)
	maxGamma := decimal.NewFromFloat(100)
	
	// Override from config if available
	if config.MaxPositionSize != "" {
		if val, err := decimal.NewFromString(config.MaxPositionSize); err == nil {
			maxPosition = val
		}
	}
	if config.MaxDeltaExposure != "" {
		if val, err := decimal.NewFromString(config.MaxDeltaExposure); err == nil {
			maxDelta = val
		}
	}
	
	return &RiskManager{
		config:           config,
		positions:        make(map[string]*Position),
		maxPositionSize:  maxPosition,
		maxDeltaExposure: maxDelta,
		maxGammaExposure: maxGamma,
	}
}

// ValidateTrade checks if a trade is within risk limits
func (rm *RiskManager) ValidateTrade(trade *TradeEvent) error {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	// Check position size limit
	currentPosition, exists := rm.positions[trade.Instrument]
	newPositionSize := trade.Quantity
	if exists {
		// Calculate net position after trade
		if trade.IsTakerBuy {
			// We're selling, so position decreases
			newPositionSize = currentPosition.Quantity.Sub(trade.Quantity)
		} else {
			// We're buying, so position increases
			newPositionSize = currentPosition.Quantity.Add(trade.Quantity)
		}
	}
	
	if newPositionSize.Abs().GreaterThan(rm.maxPositionSize) {
		return fmt.Errorf("position size %s exceeds limit %s", 
			newPositionSize.Abs().String(), 
			rm.maxPositionSize.String())
	}
	
	// Estimate delta impact (simplified - should use proper Greeks)
	estimatedDelta := rm.estimateDeltaImpact(trade)
	totalDelta := rm.calculateTotalDelta().Add(estimatedDelta)
	
	if totalDelta.Abs().GreaterThan(rm.maxDeltaExposure) {
		return fmt.Errorf("delta exposure %s exceeds limit %s", 
			totalDelta.Abs().String(), 
			rm.maxDeltaExposure.String())
	}
	
	// Check if we have too many positions
	maxPositions := 50 // Configurable
	if !exists && len(rm.positions) >= maxPositions {
		return fmt.Errorf("maximum number of positions (%d) reached", maxPositions)
	}
	
	log.Printf("Trade validated: position=%s, delta impact=%s", 
		newPositionSize.String(), estimatedDelta.String())
	
	return nil
}

// UpdatePosition updates position after trade execution
func (rm *RiskManager) UpdatePosition(trade *TradeEvent, delta, gamma decimal.Decimal) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	
	position, exists := rm.positions[trade.Instrument]
	if !exists {
		position = &Position{
			Instrument: trade.Instrument,
			Quantity:   decimal.Zero,
			AvgPrice:   decimal.Zero,
		}
		rm.positions[trade.Instrument] = position
	}
	
	// Update quantity
	if trade.IsTakerBuy {
		// We sold
		position.Quantity = position.Quantity.Sub(trade.Quantity)
	} else {
		// We bought
		position.Quantity = position.Quantity.Add(trade.Quantity)
	}
	
	// Update average price (simplified - should use weighted average)
	if !position.Quantity.IsZero() {
		position.AvgPrice = trade.Price
	}
	
	// Update Greeks
	position.Delta = delta
	position.Gamma = gamma
	position.LastUpdated = time.Now()
	
	// Remove position if closed
	if position.Quantity.IsZero() {
		delete(rm.positions, trade.Instrument)
		log.Printf("Position closed for %s", trade.Instrument)
	} else {
		log.Printf("Position updated for %s: qty=%s, delta=%s, gamma=%s",
			trade.Instrument,
			position.Quantity.String(),
			position.Delta.String(),
			position.Gamma.String())
	}
}

// GetRiskMetrics returns current risk metrics
func (rm *RiskManager) GetRiskMetrics() RiskMetrics {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	totalDelta := rm.calculateTotalDelta()
	totalGamma := rm.calculateTotalGamma()
	
	var maxPos decimal.Decimal
	for _, pos := range rm.positions {
		if pos.Quantity.Abs().GreaterThan(maxPos) {
			maxPos = pos.Quantity.Abs()
		}
	}
	
	return RiskMetrics{
		TotalDelta:      totalDelta,
		TotalGamma:      totalGamma,
		TotalPositions:  len(rm.positions),
		MaxPositionSize: maxPos,
		UpdatedAt:       time.Now(),
	}
}

// GetPositions returns current positions
func (rm *RiskManager) GetPositions() map[string]Position {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	positions := make(map[string]Position)
	for k, v := range rm.positions {
		positions[k] = *v
	}
	return positions
}

// CheckStopLoss checks if any position needs stop loss
func (rm *RiskManager) CheckStopLoss(currentPrices map[string]decimal.Decimal) []string {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	
	var stopLossPositions []string
	stopLossPct := decimal.NewFromFloat(0.1) // 10% stop loss
	
	for instrument, position := range rm.positions {
		if currentPrice, exists := currentPrices[instrument]; exists {
			loss := position.AvgPrice.Sub(currentPrice).Div(position.AvgPrice)
			if loss.GreaterThan(stopLossPct) {
				stopLossPositions = append(stopLossPositions, instrument)
				log.Printf("Stop loss triggered for %s: loss=%s%%", 
					instrument, 
					loss.Mul(decimal.NewFromInt(100)).String())
			}
		}
	}
	
	return stopLossPositions
}

// Helper methods

func (rm *RiskManager) calculateTotalDelta() decimal.Decimal {
	total := decimal.Zero
	for _, pos := range rm.positions {
		// Delta is per unit, multiply by position size
		total = total.Add(pos.Delta.Mul(pos.Quantity))
	}
	return total
}

func (rm *RiskManager) calculateTotalGamma() decimal.Decimal {
	total := decimal.Zero
	for _, pos := range rm.positions {
		// Gamma is per unit, multiply by position size
		total = total.Add(pos.Gamma.Mul(pos.Quantity))
	}
	return total
}

func (rm *RiskManager) estimateDeltaImpact(trade *TradeEvent) decimal.Decimal {
	// Simplified delta estimation
	// In reality, would fetch Greeks from pricing model
	
	baseDelta := decimal.NewFromFloat(0.5) // ATM assumption
	
	// Adjust for moneyness (very simplified)
	// Put options have negative delta
	if trade.IsPut {
		baseDelta = baseDelta.Neg()
	}
	
	// Impact depends on trade direction
	deltaImpact := baseDelta.Mul(trade.Quantity)
	if trade.IsTakerBuy {
		// We're selling, so our delta decreases
		deltaImpact = deltaImpact.Neg()
	}
	
	return deltaImpact
}