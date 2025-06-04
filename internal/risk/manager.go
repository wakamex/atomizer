package risk

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wakamex/atomizer/internal/config"
	"github.com/wakamex/atomizer/internal/types"
	"github.com/shopspring/decimal"
)

// Manager implements risk management for trading operations
type Manager struct {
	config           *config.Config
	positions        map[string]*types.Position
	maxPositionSize  decimal.Decimal
	maxDeltaExposure decimal.Decimal
	maxGammaExposure decimal.Decimal
	stopLossThreshold decimal.Decimal
	mu               sync.RWMutex
}

// NewManager creates a new risk manager
func NewManager(cfg *config.Config) *Manager {
	// Default limits - should be configurable
	maxPosition := decimal.NewFromFloat(1000)
	maxDelta := decimal.NewFromFloat(500)
	maxGamma := decimal.NewFromFloat(100)
	stopLoss := decimal.NewFromFloat(0.1) // 10% stop loss
	
	// Override from config if available
	if cfg.MaxPositionDelta > 0 {
		maxDelta = decimal.NewFromFloat(cfg.MaxPositionDelta)
	}
	
	return &Manager{
		config:            cfg,
		positions:         make(map[string]*types.Position),
		maxPositionSize:   maxPosition,
		maxDeltaExposure:  maxDelta,
		maxGammaExposure:  maxGamma,
		stopLossThreshold: stopLoss,
	}
}

// ValidateTrade checks if a trade is within risk limits
func (m *Manager) ValidateTrade(trade *types.TradeEvent) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Construct instrument name
	instrumentName := m.constructInstrumentName(trade)
	
	// Check position size limit
	currentPosition, exists := m.positions[instrumentName]
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
	
	// Check absolute position size
	if newPositionSize.Abs().GreaterThan(m.maxPositionSize) {
		return fmt.Errorf("trade would exceed maximum position size of %s", m.maxPositionSize.String())
	}
	
	// Estimate Greeks for the trade
	tradeDelta := m.estimateDelta(trade)
	tradeGamma := m.estimateGamma(trade)
	
	// Calculate portfolio Greeks after trade
	totalDelta, totalGamma := m.calculatePortfolioGreeks()
	
	if trade.IsTakerBuy {
		// We're selling, so our delta decreases
		totalDelta = totalDelta.Sub(tradeDelta.Mul(trade.Quantity))
		totalGamma = totalGamma.Sub(tradeGamma.Mul(trade.Quantity))
	} else {
		// We're buying, so our delta increases
		totalDelta = totalDelta.Add(tradeDelta.Mul(trade.Quantity))
		totalGamma = totalGamma.Add(tradeGamma.Mul(trade.Quantity))
	}
	
	// Check delta exposure
	if totalDelta.Abs().GreaterThan(m.maxDeltaExposure) {
		return fmt.Errorf("trade would exceed maximum delta exposure of %s", m.maxDeltaExposure.String())
	}
	
	// Check gamma exposure
	if totalGamma.Abs().GreaterThan(m.maxGammaExposure) {
		return fmt.Errorf("trade would exceed maximum gamma exposure of %s", m.maxGammaExposure.String())
	}
	
	return nil
}

// UpdatePosition updates position tracking after a trade
func (m *Manager) UpdatePosition(trade *types.TradeEvent) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	instrumentName := m.constructInstrumentName(trade)
	
	position, exists := m.positions[instrumentName]
	if !exists {
		position = &types.Position{
			Instrument:  instrumentName,
			Quantity:    decimal.Zero,
			AvgPrice:    decimal.Zero,
			Delta:       decimal.Zero,
			Gamma:       decimal.Zero,
			LastUpdated: time.Now(),
		}
		m.positions[instrumentName] = position
	}
	
	// Update position quantity and average price
	if trade.IsTakerBuy {
		// We're selling
		position.Quantity = position.Quantity.Sub(trade.Quantity)
	} else {
		// We're buying
		if position.Quantity.IsZero() {
			position.AvgPrice = trade.Price
		} else {
			// Calculate weighted average price
			totalValue := position.Quantity.Mul(position.AvgPrice).Add(trade.Quantity.Mul(trade.Price))
			totalQuantity := position.Quantity.Add(trade.Quantity)
			if !totalQuantity.IsZero() {
				position.AvgPrice = totalValue.Div(totalQuantity)
			}
		}
		position.Quantity = position.Quantity.Add(trade.Quantity)
	}
	
	// Update Greeks
	position.Delta = m.estimateDelta(trade)
	position.Gamma = m.estimateGamma(trade)
	position.LastUpdated = time.Now()
	
	// Log position update
	log.Printf("Updated position %s: Qty=%s, AvgPrice=%s, Delta=%s, Gamma=%s",
		instrumentName, position.Quantity.String(), position.AvgPrice.String(),
		position.Delta.String(), position.Gamma.String())
		
	// Check for stop loss
	if m.shouldTriggerStopLoss(position) {
		log.Printf("WARNING: Position %s has exceeded stop loss threshold!", instrumentName)
		// In production, this would trigger automated unwinding
	}
}

// GetGreeks returns current portfolio Greeks
func (m *Manager) GetGreeks() (delta, gamma decimal.Decimal) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	return m.calculatePortfolioGreeks()
}

// GetPositions returns current positions
func (m *Manager) GetPositions() map[string]types.Position {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Return a copy to prevent external modification
	positions := make(map[string]types.Position)
	for k, v := range m.positions {
		if v != nil {
			positions[k] = *v
		}
	}
	
	return positions
}

// GetMetrics returns current risk metrics
func (m *Manager) GetMetrics() types.RiskMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	delta, gamma := m.calculatePortfolioGreeks()
	
	maxPosition := decimal.Zero
	for _, pos := range m.positions {
		if pos != nil && pos.Quantity.Abs().GreaterThan(maxPosition) {
			maxPosition = pos.Quantity.Abs()
		}
	}
	
	return types.RiskMetrics{
		TotalDelta:      delta,
		TotalGamma:      gamma,
		TotalPositions:  len(m.positions),
		MaxPositionSize: maxPosition,
		UpdatedAt:       time.Now(),
	}
}

// constructInstrumentName creates instrument name from trade details
func (m *Manager) constructInstrumentName(trade *types.TradeEvent) string {
	// If we already have an instrument name, use it
	if trade.Instrument != "" && trade.Instrument != trade.Instrument {
		return trade.Instrument
	}
	
	// Otherwise construct from components
	optionType := "C"
	if trade.IsPut {
		optionType = "P"
	}
	
	// Format: ASSET-EXPIRY-STRIKE-TYPE
	expiryTime := time.Unix(trade.Expiry, 0)
	expiryStr := expiryTime.Format("20060102")
	
	return fmt.Sprintf("%s-%s-%s-%s", 
		"ETH", // Default to ETH, should map from asset address
		expiryStr,
		trade.Strike.String(),
		optionType)
}

// estimateDelta estimates option delta based on simple heuristics
func (m *Manager) estimateDelta(trade *types.TradeEvent) decimal.Decimal {
	// Simple estimation based on moneyness and time to expiry
	// In production, use proper Black-Scholes or get from exchange
	
	timeToExpiry := time.Until(time.Unix(trade.Expiry, 0))
	daysToExpiry := timeToExpiry.Hours() / 24
	
	// Very simple approximation
	if trade.IsPut {
		if daysToExpiry < 7 {
			return decimal.NewFromFloat(-0.3)
		}
		return decimal.NewFromFloat(-0.5)
	} else {
		if daysToExpiry < 7 {
			return decimal.NewFromFloat(0.7)
		}
		return decimal.NewFromFloat(0.5)
	}
}

// estimateGamma estimates option gamma
func (m *Manager) estimateGamma(trade *types.TradeEvent) decimal.Decimal {
	// Simple estimation - gamma is highest ATM and near expiry
	timeToExpiry := time.Until(time.Unix(trade.Expiry, 0))
	daysToExpiry := timeToExpiry.Hours() / 24
	
	if daysToExpiry < 7 {
		return decimal.NewFromFloat(0.1)
	}
	return decimal.NewFromFloat(0.05)
}

// calculatePortfolioGreeks calculates total portfolio Greeks
func (m *Manager) calculatePortfolioGreeks() (delta, gamma decimal.Decimal) {
	delta = decimal.Zero
	gamma = decimal.Zero
	
	for _, pos := range m.positions {
		if pos != nil && !pos.Quantity.IsZero() {
			delta = delta.Add(pos.Delta.Mul(pos.Quantity))
			gamma = gamma.Add(pos.Gamma.Mul(pos.Quantity))
		}
	}
	
	return delta, gamma
}

// shouldTriggerStopLoss checks if position has exceeded loss threshold
func (m *Manager) shouldTriggerStopLoss(position *types.Position) bool {
	if position.Quantity.IsZero() || position.AvgPrice.IsZero() {
		return false
	}
	
	// Calculate P&L percentage
	// This is simplified - in production would use mark price
	currentPrice := position.AvgPrice // Should use current market price
	pnlPercent := currentPrice.Sub(position.AvgPrice).Div(position.AvgPrice)
	
	// Check if loss exceeds threshold
	return pnlPercent.LessThan(m.stopLossThreshold.Neg())
}

// Ensure Manager implements the RiskManager interface
var _ types.RiskManager = (*Manager)(nil)