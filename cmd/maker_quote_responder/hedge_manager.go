package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// HedgeResult contains the result of a hedge execution
type HedgeResult struct {
	OrderID      string
	Exchange     string
	Instrument   string
	Direction    string
	Quantity     decimal.Decimal
	Price        decimal.Decimal
	Status       string
	ExecutedAt   time.Time
}

// HedgeManager manages hedge execution on exchanges
type HedgeManager struct {
	exchange      Exchange
	config        *AppConfig
	maxRetries    int
	retryDelayMs  int
}

// NewHedgeManager creates a new hedge manager
func NewHedgeManager(exchange Exchange, config *AppConfig) *HedgeManager {
	return &HedgeManager{
		exchange:     exchange,
		config:       config,
		maxRetries:   3,
		retryDelayMs: 1000,
	}
}

// ExecuteHedge places a hedge order for the given trade
func (hm *HedgeManager) ExecuteHedge(trade *TradeEvent) (*HedgeResult, error) {
	log.Printf("Executing hedge for trade %s on %s", trade.ID, hm.config.ExchangeName)
	
	// Convert trade to hedge parameters
	hedgeParams, err := hm.buildHedgeParams(trade)
	if err != nil {
		return nil, fmt.Errorf("failed to build hedge params: %w", err)
	}
	
	// Get current order book
	orderBook, err := hm.getOrderBookWithRetry(trade)
	if err != nil {
		return nil, fmt.Errorf("failed to get order book: %w", err)
	}
	
	// Calculate hedge price within spread
	hedgePrice := hm.calculateHedgePrice(orderBook, hedgeParams.isBuy)
	
	// Execute hedge with retries
	var lastErr error
	for attempt := 1; attempt <= hm.maxRetries; attempt++ {
		result, err := hm.executeSingleHedge(hedgeParams, hedgePrice)
		if err == nil {
			log.Printf("Hedge successful on attempt %d: %+v", attempt, result)
			return result, nil
		}
		
		lastErr = err
		log.Printf("Hedge attempt %d failed: %v", attempt, err)
		
		if attempt < hm.maxRetries {
			time.Sleep(time.Duration(hm.retryDelayMs) * time.Millisecond)
		}
	}
	
	return nil, fmt.Errorf("hedge failed after %d attempts: %w", hm.maxRetries, lastErr)
}

// buildHedgeParams converts trade event to hedge parameters
func (hm *HedgeManager) buildHedgeParams(trade *TradeEvent) (*hedgeParams, error) {
	// For Rysk trades, we're the maker
	// If taker buys from us, we sell to them, so we need to buy on exchange to hedge
	hedgeDirection := !trade.IsTakerBuy
	
	// Convert instrument name for exchange
	instrument, err := hm.convertInstrumentName(trade)
	if err != nil {
		return nil, err
	}
	
	return &hedgeParams{
		instrument: instrument,
		quantity:   trade.Quantity,
		isBuy:      hedgeDirection,
		tradeID:    trade.ID,
	}, nil
}

// convertInstrumentName converts from internal format to exchange format
func (hm *HedgeManager) convertInstrumentName(trade *TradeEvent) (string, error) {
	// Extract components from trade
	underlying := hm.extractUnderlying(trade.Instrument)
	
	// Log the strike value being passed
	log.Printf("[HedgeManager] Converting instrument with strike: %s (decimal: %s)", 
		trade.Strike.String(), trade.Strike.String())
	
	// Use exchange's conversion method
	instrument, err := hm.exchange.ConvertToInstrument(
		underlying,
		trade.Strike.String(),
		trade.Expiry,
		trade.IsPut,
	)
	if err != nil {
		return "", fmt.Errorf("failed to convert instrument: %w", err)
	}
	
	return instrument, nil
}

// extractUnderlying gets the underlying asset from instrument name
func (hm *HedgeManager) extractUnderlying(instrument string) string {
	// Simple extraction - assumes format like "ETH-20231225-3000-C"
	parts := strings.Split(instrument, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return "ETH" // Default
}

// getOrderBookWithRetry fetches order book with retry logic
func (hm *HedgeManager) getOrderBookWithRetry(trade *TradeEvent) (CCXTOrderBook, error) {
	// Extract underlying asset from instrument name (e.g., "ETH" from "ETH-20250529-2550-C")
	underlying := hm.extractUnderlying(trade.Instrument)
	
	// Create RFQ with actual trade values for the order book request
	rfq := RFQResult{
		Strike: trade.Strike.String(), // Pass actual strike in wei
		Expiry: trade.Expiry,
		IsPut:  trade.IsPut,
	}
	
	var lastErr error
	for i := 0; i < 3; i++ {
		orderBook, err := hm.exchange.GetOrderBook(rfq, underlying)
		if err == nil {
			return orderBook, nil
		}
		lastErr = err
		time.Sleep(500 * time.Millisecond)
	}
	
	return CCXTOrderBook{}, lastErr
}

// calculateHedgePrice determines the price for hedging within the spread
func (hm *HedgeManager) calculateHedgePrice(orderBook CCXTOrderBook, isBuy bool) decimal.Decimal {
	if isBuy {
		// Buying - use best ask or slightly below
		if len(orderBook.Asks) > 0 && len(orderBook.Asks[0]) >= 2 {
			bestAsk := decimal.NewFromFloat(orderBook.Asks[0][0])
			// Place order slightly below best ask to ensure fill
			return bestAsk.Mul(decimal.NewFromFloat(0.999))
		}
	} else {
		// Selling - use defensive 2x ask strategy during testing
		if len(orderBook.Asks) > 0 && len(orderBook.Asks[0]) >= 2 {
			bestAsk := decimal.NewFromFloat(orderBook.Asks[0][0])
			// Place our ask at 2x the best ask for safety (far from top of book)
			hedgePrice := bestAsk.Mul(decimal.NewFromFloat(2.0))
			log.Printf("[HedgeManager] Defensive sell: best ask %s, placing at %s (2x)", 
				bestAsk.String(), hedgePrice.String())
			return hedgePrice
		}
	}
	
	// Fallback to a safe price if no order book
	log.Printf("Warning: No order book data, using fallback price")
	return decimal.NewFromFloat(0.1) // Conservative fallback
}

// executeSingleHedge executes a single hedge attempt
func (hm *HedgeManager) executeSingleHedge(params *hedgeParams, price decimal.Decimal) (*HedgeResult, error) {
	// Build confirmation for exchange
	confirmation := RFQConfirmation{
		Nonce:      params.tradeID,
		Quantity:   params.quantity.Mul(decimal.New(1, 18)).String(), // Convert to wei
		IsTakerBuy: !params.isBuy, // Invert because we're hedging
		// Other fields would be populated from trade details
	}
	
	// Place hedge order
	err := hm.exchange.PlaceHedgeOrder(confirmation, params.instrument, hm.config)
	if err != nil {
		return nil, err
	}
	
	// Build result
	direction := "sell"
	if params.isBuy {
		direction = "buy"
	}
	
	return &HedgeResult{
		OrderID:    fmt.Sprintf("hedge_%s_%d", params.tradeID, time.Now().Unix()),
		Exchange:   hm.config.ExchangeName,
		Instrument: params.instrument,
		Direction:  direction,
		Quantity:   params.quantity,
		Price:      price,
		Status:     "placed",
		ExecutedAt: time.Now(),
	}, nil
}

// hedgeParams contains parameters for hedge execution
type hedgeParams struct {
	instrument string
	quantity   decimal.Decimal
	isBuy      bool
	tradeID    string
}