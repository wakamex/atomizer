package hedging

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wakamex/atomizer/internal/exchange/derive"
	"github.com/wakamex/atomizer/internal/types"
)

// PureGammaHedger implements aggressive gamma hedging for options positions
// Unlike regular gamma hedging, this focuses on risk reduction, not market making
type PureGammaHedger struct {
	exchange         types.MarketMakerExchange
	debugMode        bool  // Enable debug logging
	
	// Greek tracking
	positions        map[string]*OptionPosition  // instrument -> position with greeks
	netDelta         decimal.Decimal
	netGamma         decimal.Decimal
	
	// Hedging parameters
	deltaThreshold   decimal.Decimal  // Max delta before hedging
	minHedgeSize     decimal.Decimal  // Minimum hedge size (exchange minimum)
	hedgeInterval    time.Duration    // How often to check/rehedge
	aggressiveness   decimal.Decimal  // How far through the spread to place orders (0-1)
	
	// Hedge tracking
	currentHedge     *HedgePosition   // Current perp/spot hedge position
	lastHedgeTime    time.Time
	consecutiveFails int              // Track consecutive hedge failures
	
	// Control
	ctx              context.Context
	cancel           context.CancelFunc
	running          bool
	mu               sync.RWMutex
}

// OptionPosition represents an options position with Greeks
type OptionPosition struct {
	Instrument   string
	Quantity     decimal.Decimal
	AvgPrice     decimal.Decimal
	Delta        decimal.Decimal
	Gamma        decimal.Decimal
	LastUpdated  time.Time
}

// HedgePosition represents a hedge position in perps/spot
type HedgePosition struct {
	Instrument string
	Quantity   decimal.Decimal
	AvgPrice   decimal.Decimal
	UpdatedAt  time.Time
}

// NewPureGammaHedger creates a new pure gamma hedger
func NewPureGammaHedger(exchange types.MarketMakerExchange) *PureGammaHedger {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &PureGammaHedger{
		exchange:        exchange,
		positions:       make(map[string]*OptionPosition),
		deltaThreshold:  decimal.NewFromFloat(0.1),    // 0.1 ETH delta threshold
		minHedgeSize:    decimal.NewFromFloat(0.1),    // Exchange minimum
		hedgeInterval:   30 * time.Second,
		aggressiveness:  decimal.NewFromFloat(0.5),    // Cross 50% of the spread
		ctx:             ctx,
		cancel:          cancel,
	}
}

// SetParameters updates hedger parameters
func (gh *PureGammaHedger) SetParameters(deltaThreshold, minHedgeSize, aggressiveness decimal.Decimal, hedgeInterval time.Duration) {
	gh.mu.Lock()
	defer gh.mu.Unlock()
	
	gh.deltaThreshold = deltaThreshold
	gh.minHedgeSize = minHedgeSize
	gh.aggressiveness = aggressiveness
	gh.hedgeInterval = hedgeInterval
}

// SetDebugMode enables/disables debug logging
func (gh *PureGammaHedger) SetDebugMode(debug bool) {
	gh.debugMode = debug
}

// Start begins the pure gamma hedging loop
func (gh *PureGammaHedger) Start() error {
	gh.mu.Lock()
	if gh.running {
		gh.mu.Unlock()
		return fmt.Errorf("pure gamma hedger already running")
	}
	gh.running = true
	gh.mu.Unlock()
	
	log.Printf("Starting Pure Gamma Hedger")
	log.Printf("Configuration:")
	log.Printf("  Delta Threshold: %s ETH", gh.deltaThreshold.String())
	log.Printf("  Min Hedge Size: %s ETH", gh.minHedgeSize.String()) 
	log.Printf("  Hedge Interval: %v", gh.hedgeInterval)
	log.Printf("  Aggressiveness: %s (0=passive, 1=aggressive)", gh.aggressiveness.String())
	
	// Subscribe to ETH-PERP orderbook
	log.Printf("Subscribing to ETH-PERP orderbook...")
	// Try to cast to derive exchange to access SubscribeOrderBook
	if deriveExchange, ok := gh.exchange.(*derive.DeriveMarketMakerExchange); ok {
		if err := deriveExchange.SubscribeOrderBook("ETH-PERP"); err != nil {
			log.Printf("Warning: Failed to subscribe to ETH-PERP orderbook: %v", err)
		}
	} else {
		log.Printf("Warning: Exchange does not support orderbook subscription")
	}
	
	// Start hedging loop
	go gh.hedgingLoop()
	
	return nil
}

// Stop stops the hedger
func (gh *PureGammaHedger) Stop() {
	gh.mu.Lock()
	gh.running = false
	gh.mu.Unlock()
	
	gh.cancel()
	log.Printf("Pure Gamma Hedger stopped")
}

// hedgingLoop is the main hedging loop
func (gh *PureGammaHedger) hedgingLoop() {
	// Wait for WebSocket and orderbook to be ready
	log.Printf("Waiting for orderbook data...")
	time.Sleep(3 * time.Second)
	
	// Initial position load and hedge
	log.Printf("Starting initial hedge check...")
	gh.loadPositions()
	gh.checkAndHedge()
	
	ticker := time.NewTicker(gh.hedgeInterval)
	defer ticker.Stop()
	
	log.Printf("Entering main loop with interval %v", gh.hedgeInterval)
	for {
		select {
		case <-gh.ctx.Done():
			log.Printf("Hedging loop cancelled")
			return
		case <-ticker.C:
			log.Printf("Hedge interval tick - checking positions...")
			gh.loadPositions()
			gh.checkAndHedge()
		}
	}
}

// loadPositions loads current options positions from exchange
func (gh *PureGammaHedger) loadPositions() {
	positions, err := gh.exchange.GetPositions()
	if err != nil {
		log.Printf("ERROR: Failed to load positions: %v", err)
		return
	}
	
	gh.mu.Lock()
	defer gh.mu.Unlock()
	
	// Clear old positions
	gh.positions = make(map[string]*OptionPosition)
	gh.netDelta = decimal.Zero
	gh.netGamma = decimal.Zero
	
	// Process each position
	for _, pos := range positions {
		// Create position
		optPos := &OptionPosition{
			Instrument:  pos.InstrumentName,
			Quantity:    decimal.NewFromFloat(pos.Amount),
			AvgPrice:    decimal.NewFromFloat(pos.AveragePrice),
			Delta:       decimal.Zero, // Will be updated from ticker
			Gamma:       decimal.Zero, // Will be updated from ticker
			LastUpdated: time.Now(),
		}
		
		// Log every position we find
		log.Printf("Found position: %s, amount=%.4f, avgPrice=%.2f", 
			pos.InstrumentName, pos.Amount, pos.AveragePrice)
		
		// Check if it's ETH-PERP (our hedge instrument)
		if pos.InstrumentName == "ETH-PERP" {
			// Track current hedge position
			gh.currentHedge = &HedgePosition{
				Instrument: pos.InstrumentName,
				Quantity:   optPos.Quantity,
				AvgPrice:   optPos.AvgPrice,
				UpdatedAt:  time.Now(),
			}
			// ETH-PERP has delta of 1
			optPos.Delta = decimal.NewFromFloat(1.0)
			optPos.Gamma = decimal.Zero
		}
		
		// Store position
		gh.positions[pos.InstrumentName] = optPos
		
		if gh.debugMode {
			log.Printf("[DEBUG] Loaded position: %s, qty=%s", pos.InstrumentName, optPos.Quantity.String())
		}
	}
	
	log.Printf("Loaded %d option positions", len(gh.positions))
}

// checkAndHedge checks if hedging is needed and executes it
func (gh *PureGammaHedger) checkAndHedge() {
	// Load latest positions
	gh.loadPositions()
	
	// Update Greeks from market data
	gh.updateGreeks()
	
	// Calculate net exposure
	gh.mu.RLock()
	netDelta := gh.netDelta
	gh.mu.RUnlock()
	
	// Check if we have any options positions
	hasOptions := false
	for _, pos := range gh.positions {
		if isOptionInstrument(pos.Instrument) && !pos.Quantity.IsZero() {
			hasOptions = true
			break
		}
	}
	
	// Special case: no options but have perp position - always close it
	if !hasOptions && !netDelta.IsZero() {
		log.Printf("No options positions found, closing perp hedge of %s ETH", netDelta.Neg().StringFixed(4))
		// Force hedge regardless of threshold
	} else if netDelta.Abs().LessThan(gh.deltaThreshold) {
		// Normal threshold check
		if gh.debugMode {
			log.Printf("[DEBUG] No hedge needed. Net delta: %s (threshold: %s)", 
				netDelta.String(), gh.deltaThreshold.String())
		}
		return
	}
	
	// Calculate hedge size
	hedgeSize := netDelta.Neg() // Hedge in opposite direction
	
	// Check minimum size (skip for closing positions when no options)
	if hedgeSize.Abs().LessThan(gh.minHedgeSize) {
		// Always allow closing positions when no options
		if !hasOptions && !netDelta.IsZero() {
			log.Printf("Closing position of %s ETH (below min size %s, but closing allowed)", 
				hedgeSize.Abs().StringFixed(4), gh.minHedgeSize.StringFixed(4))
		} else {
			log.Printf("Hedge size %s below minimum %s, skipping", 
				hedgeSize.Abs().StringFixed(4), gh.minHedgeSize.StringFixed(4))
			return
		}
	}
	
	// Execute hedge
	log.Printf("Executing hedge: Net delta=%s, Hedge size=%s", 
		netDelta.String(), hedgeSize.String())
	
	// Debug: show what we're about to do
	action := "BUY"
	if hedgeSize.IsNegative() {
		action = "SELL"
	}
	log.Printf("Action: Will %s %s ETH to hedge", action, hedgeSize.Abs().StringFixed(4))
	
	if err := gh.executeHedge(hedgeSize); err != nil {
		// Check if it's a minimum size error
		if strings.Contains(err.Error(), "Order amount must be >") {
			log.Printf("NOTICE: Cannot close position of %s ETH - exchange minimum is 0.1 ETH", hedgeSize.Abs().StringFixed(4))
			log.Printf("NOTICE: Options to handle this position:")
			log.Printf("  1. Open a larger opposite position (e.g., buy 0.1 ETH), then close the net position")
			log.Printf("  2. Use the exchange web interface which may allow smaller closes")
			log.Printf("  3. Wait for options positions that require hedging >= 0.1 ETH")
			// Don't increment failure count for minimum size errors
			return
		}
		
		gh.consecutiveFails++
		log.Printf("ERROR: Hedge failed (attempt %d): %v", gh.consecutiveFails, err)
		
		// After 3 failures, try market order
		if gh.consecutiveFails >= 3 {
			log.Printf("WARNING: 3 consecutive hedge failures. Attempting MARKET order...")
			if err := gh.executeMarketHedge(hedgeSize); err != nil {
				// Check again for minimum size error
				if strings.Contains(err.Error(), "Order amount must be >") {
					log.Printf("NOTICE: Market order also below minimum. Position size %s < 0.1 ETH minimum", hedgeSize.Abs().StringFixed(4))
					log.Printf("NOTICE: Manual intervention required - see options above")
				} else {
					log.Printf("ERROR: Market hedge also failed: %v", err)
					log.Printf("CRITICAL: Manual intervention required!")
				}
			} else {
				gh.consecutiveFails = 0
				gh.lastHedgeTime = time.Now()
				log.Printf("Market hedge successful at %s", time.Now().Format("15:04:05"))
			}
		}
	} else {
		gh.consecutiveFails = 0
		gh.lastHedgeTime = time.Now()
		log.Printf("Hedge successful at %s", time.Now().Format("15:04:05"))
	}
}

// updateGreeks updates Greeks from market data (ticker updates)
func (gh *PureGammaHedger) updateGreeks() {
	gh.mu.Lock()
	defer gh.mu.Unlock()
	
	// Reset totals
	gh.netDelta = decimal.Zero
	gh.netGamma = decimal.Zero
	
	// Track subtotals for reporting
	optionsDelta := decimal.Zero
	perpsDelta := decimal.Zero
	
	// Process all positions
	for _, pos := range gh.positions {
		instrument := pos.Instrument
		
		// Handle perps (delta = 1:1)
		if instrument == "ETH-PERP" {
			perpsDelta = perpsDelta.Add(pos.Quantity)
			gh.netDelta = gh.netDelta.Add(pos.Quantity)
			continue
		}
		
		// For options, fetch real-time Greeks via ticker
		if ticker, err := gh.fetchTicker(instrument); err == nil {
			// Update Greeks from ticker
			if ticker.Delta != nil {
				pos.Delta = *ticker.Delta
			}
			if ticker.Gamma != nil {
				pos.Gamma = *ticker.Gamma
			}
			pos.LastUpdated = time.Now()
			
			// Calculate position Greeks
			posDelta := pos.Delta.Mul(pos.Quantity)
			posGamma := pos.Gamma.Mul(pos.Quantity)
			
			optionsDelta = optionsDelta.Add(posDelta)
			gh.netDelta = gh.netDelta.Add(posDelta)
			gh.netGamma = gh.netGamma.Add(posGamma)
			
			if gh.debugMode {
				log.Printf("[DEBUG] %s: qty=%s, delta=%s, gamma=%s, posDelta=%s", 
					instrument, pos.Quantity.String(), pos.Delta.String(), 
					pos.Gamma.String(), posDelta.String())
			}
		} else {
			// Fallback to default Greeks if ticker fetch fails
			log.Printf("Warning: Failed to fetch ticker for %s: %v (using defaults)", instrument, err)
			pos.Delta = decimal.NewFromFloat(0.5)
			pos.Gamma = decimal.NewFromFloat(0.01)
			
			posDelta := pos.Delta.Mul(pos.Quantity)
			posGamma := pos.Gamma.Mul(pos.Quantity)
			
			optionsDelta = optionsDelta.Add(posDelta)
			gh.netDelta = gh.netDelta.Add(posDelta)
			gh.netGamma = gh.netGamma.Add(posGamma)
		}
	}
	
	// Check if we have any options positions
	hasOptions := false
	for _, pos := range gh.positions {
		if isOptionInstrument(pos.Instrument) && !pos.Quantity.IsZero() {
			hasOptions = true
			break
		}
	}
	
	// Log portfolio Greeks with threshold status
	if !optionsDelta.IsZero() || !perpsDelta.IsZero() {
		status := "within threshold, no hedge needed"
		if !hasOptions && !perpsDelta.IsZero() {
			status = "no options positions, should close perp hedge"
		} else if gh.netDelta.Abs().GreaterThanOrEqual(gh.deltaThreshold) {
			status = "exceeds threshold, hedging required"
		}
		log.Printf("Portfolio Delta: Options=%s, Perps=%s, Total=%s ETH (%s)", 
			optionsDelta.StringFixed(4), perpsDelta.StringFixed(4), 
			gh.netDelta.StringFixed(4), status)
	}
	
	// If no options positions but we have perp position, keep the actual delta
	// The hedge calculation will invert it to close the position
	if !hasOptions && !perpsDelta.IsZero() {
		// gh.netDelta already contains perpsDelta from the loop above
		// Don't modify it here - the hedge calculation will handle the inversion
	}
}

// fetchTicker fetches real-time ticker data for an instrument
func (gh *PureGammaHedger) fetchTicker(instrument string) (*types.TickerUpdate, error) {
	// Subscribe to ticker if needed
	tickerCh, err := gh.exchange.SubscribeTickers(gh.ctx, []string{instrument})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to ticker: %w", err)
	}
	
	// Get the latest ticker update (with timeout)
	select {
	case ticker := <-tickerCh:
		if ticker.Instrument == instrument {
			return &ticker, nil
		}
	case <-time.After(2 * time.Second):
		return nil, fmt.Errorf("ticker timeout")
	}
	
	return nil, fmt.Errorf("no ticker data received")
}

// executeHedge places the hedge order
func (gh *PureGammaHedger) executeHedge(size decimal.Decimal) error {
	instrument := "ETH-PERP" // Hedge instrument
	minOrderSize := decimal.NewFromFloat(0.1) // Exchange minimum
	
	log.Printf("executeHedge called with size: %s", size.StringFixed(4))
	
	// Debug what we're going to do
	if size.IsPositive() {
		log.Printf("executeHedge: Need to BUY %s ETH", size.StringFixed(4))
	} else {
		log.Printf("executeHedge: Need to SELL %s ETH", size.Abs().StringFixed(4))
	}
	
	// Check if we need to handle minimum size issue
	if size.Abs().LessThan(minOrderSize) {
		// Get current position to understand what we're closing
		positions, err := gh.exchange.GetPositions()
		if err != nil {
			return fmt.Errorf("failed to get positions: %w", err)
		}
		
		var currentPosition decimal.Decimal
		for _, pos := range positions {
			if pos.InstrumentName == instrument {
				currentPosition = decimal.NewFromFloat(pos.Amount)
				break
			}
		}
		
		log.Printf("Position size %s is below minimum %s, using increase-then-close strategy", 
			size.Abs().StringFixed(4), minOrderSize.StringFixed(4))
		log.Printf("Current actual position: %s ETH", currentPosition.StringFixed(4))
		return gh.executeMinSizeClose(size, currentPosition, minOrderSize)
	}
	
	// Get current orderbook
	orderBook, err := gh.exchange.GetOrderBook(instrument)
	if err != nil {
		return fmt.Errorf("failed to get orderbook: %w", err)
	}
	
	log.Printf("Got orderbook - Bids: %d, Asks: %d", len(orderBook.Bids), len(orderBook.Asks))
	
	// Determine side
	side := "buy"
	if size.IsNegative() {
		side = "sell"
		size = size.Abs()
	}
	
	// Calculate hedge price based on aggressiveness
	var hedgePrice decimal.Decimal
	if side == "buy" {
		// For buys, start from best ask and move toward bid
		if len(orderBook.Asks) == 0 {
			return fmt.Errorf("no asks in orderbook")
		}
		
		bestAsk := orderBook.Asks[0].Price
		bestBid := decimal.Zero
		if len(orderBook.Bids) > 0 {
			bestBid = orderBook.Bids[0].Price
		}
		
		// Price = ask - (ask-bid) * aggressiveness
		spread := bestAsk.Sub(bestBid)
		hedgePrice = bestAsk.Sub(spread.Mul(gh.aggressiveness))
	} else {
		// For sells, start from best bid and move toward ask
		if len(orderBook.Bids) == 0 {
			return fmt.Errorf("no bids in orderbook")
		}
		
		bestBid := orderBook.Bids[0].Price
		bestAsk := decimal.NewFromFloat(999999)
		if len(orderBook.Asks) > 0 {
			bestAsk = orderBook.Asks[0].Price
		}
		
		// Price = bid + (ask-bid) * aggressiveness
		spread := bestAsk.Sub(bestBid)
		hedgePrice = bestBid.Add(spread.Mul(gh.aggressiveness))
	}
	
	// Round price to tick size (assume 0.1)
	tickSize := decimal.NewFromFloat(0.1)
	hedgePrice = hedgePrice.Div(tickSize).Round(0).Mul(tickSize)
	
	log.Printf("Placing hedge order: %s %s %s @ %s", 
		side, size.String(), instrument, hedgePrice.String())
	
	// Final debug before placing order
	log.Printf("FINAL CHECK - About to place order: Side=%s, Size=%s", side, size.StringFixed(4))
	
	// Place the order
	orderID, err := gh.exchange.PlaceLimitOrder(instrument, side, hedgePrice, size)
	if err != nil {
		return fmt.Errorf("failed to place order: %w", err)
	}
	
	log.Printf("Hedge order placed: ID=%s", orderID)
	
	// Update current hedge position
	gh.mu.Lock()
	gh.currentHedge = &HedgePosition{
		Instrument: instrument,
		Quantity:   size,
		AvgPrice:   hedgePrice,
		UpdatedAt:  time.Now(),
	}
	gh.mu.Unlock()
	
	return nil
}

// executeMarketHedge uses market orders as last resort
func (gh *PureGammaHedger) executeMarketHedge(size decimal.Decimal) error {
	instrument := "ETH-PERP"
	
	// Get current orderbook
	orderBook, err := gh.exchange.GetOrderBook(instrument)
	if err != nil {
		return fmt.Errorf("failed to get orderbook for market order: %w", err)
	}
	
	// Determine side
	side := "buy"
	if size.IsNegative() {
		side = "sell"
		size = size.Abs()
	}
	
	// For market orders, use very aggressive pricing
	var marketPrice decimal.Decimal
	if side == "sell" && len(orderBook.Bids) > 0 {
		// Price well below bid to ensure fill
		marketPrice = orderBook.Bids[0].Price.Mul(decimal.NewFromFloat(0.99))
	} else if side == "buy" && len(orderBook.Asks) > 0 {
		// Price well above ask to ensure fill
		marketPrice = orderBook.Asks[0].Price.Mul(decimal.NewFromFloat(1.01))
	} else {
		return fmt.Errorf("no %s liquidity in orderbook", side)
	}
	
	// Round to tick size
	tickSize := decimal.NewFromFloat(0.1)
	marketPrice = marketPrice.Div(tickSize).Round(0).Mul(tickSize)
	
	log.Printf("Placing MARKET order (as aggressive limit): %s %s %s @ %s", 
		side, size.String(), instrument, marketPrice.String())
	
	// Place the order
	orderID, err := gh.exchange.PlaceLimitOrder(instrument, side, marketPrice, size)
	if err != nil {
		return fmt.Errorf("failed to place market order: %w", err)
	}
	
	log.Printf("Market hedge order placed: ID=%s", orderID)
	
	// Update hedge position
	gh.mu.Lock()
	gh.currentHedge = &HedgePosition{
		Instrument: instrument,
		Quantity:   size,
		AvgPrice:   marketPrice,
		UpdatedAt:  time.Now(),
	}
	gh.mu.Unlock()
	
	return nil
}

// executeMinSizeClose handles closing positions below minimum size
// Strategy: First increase position by minimum order size, then close it entirely
func (gh *PureGammaHedger) executeMinSizeClose(hedgeSize, currentPosition, minOrderSize decimal.Decimal) error {
	instrument := "ETH-PERP"
	
	log.Printf("Executing minimum size close strategy")
	log.Printf("  Current position: %s ETH", currentPosition.StringFixed(4))
	log.Printf("  Target hedge: %s ETH (to close position)", hedgeSize.StringFixed(4))
	
	// Step 1: Increase position by the minimum order size (0.1 ETH)
	// This will make our position larger than minimum so we can close it
	newPosition := currentPosition
	
	// Determine which side increases our position
	var increaseSide string
	if currentPosition.IsNegative() {
		// We're short, need to sell more to make it more negative
		increaseSide = "sell"
		newPosition = currentPosition.Sub(minOrderSize) // More negative
		log.Printf("Step 1: Increasing short position by selling %s ETH (from %s to %s)", 
			minOrderSize.StringFixed(4), currentPosition.StringFixed(4), newPosition.StringFixed(4))
	} else {
		// We're long, need to buy more to make it more positive
		increaseSide = "buy"
		newPosition = currentPosition.Add(minOrderSize) // More positive
		log.Printf("Step 1: Increasing long position by buying %s ETH (from %s to %s)", 
			minOrderSize.StringFixed(4), currentPosition.StringFixed(4), newPosition.StringFixed(4))
	}
	
	// Get orderbook for pricing
	orderBook, err := gh.exchange.GetOrderBook(instrument)
	if err != nil {
		return fmt.Errorf("failed to get orderbook: %w", err)
	}
	
	// Calculate aggressive price for increase order
	var increasePrice decimal.Decimal
	if increaseSide == "buy" {
		if len(orderBook.Asks) == 0 {
			return fmt.Errorf("no asks in orderbook")
		}
		// Price above ask to ensure fill
		increasePrice = orderBook.Asks[0].Price.Mul(decimal.NewFromFloat(1.001))
	} else {
		if len(orderBook.Bids) == 0 {
			return fmt.Errorf("no bids in orderbook")
		}
		// Price below bid to ensure fill
		increasePrice = orderBook.Bids[0].Price.Mul(decimal.NewFromFloat(0.999))
	}
	
	// Round to tick size
	tickSize := decimal.NewFromFloat(0.1)
	increasePrice = increasePrice.Div(tickSize).Round(0).Mul(tickSize)
	
	// Place the increase order (minimum size)
	log.Printf("Placing increase order: %s %s @ %s", increaseSide, minOrderSize.StringFixed(4), increasePrice.StringFixed(2))
	increaseOrderID, err := gh.exchange.PlaceLimitOrder(instrument, increaseSide, increasePrice, minOrderSize)
	if err != nil {
		return fmt.Errorf("failed to place increase order: %w", err)
	}
	
	log.Printf("Increase order placed: %s. Waiting for fill...", increaseOrderID)
	
	// Wait for the order to fill (with timeout)
	filled := false
	for i := 0; i < 10; i++ { // Wait up to 10 seconds
		time.Sleep(1 * time.Second)
		
		// Check if order is filled by looking at open orders
		openOrders, err := gh.exchange.GetOpenOrders()
		if err != nil {
			log.Printf("Error checking order status: %v", err)
			continue
		}
		
		// If order is not in open orders, it's likely filled
		orderFound := false
		for _, order := range openOrders {
			if order.OrderID == increaseOrderID {
				orderFound = true
				break
			}
		}
		
		if !orderFound {
			filled = true
			log.Printf("Increase order filled!")
			break
		}
	}
	
	if !filled {
		// Cancel the increase order if it didn't fill
		log.Printf("Increase order did not fill in time, cancelling...")
		if err := gh.exchange.CancelOrder(increaseOrderID); err != nil {
			log.Printf("Failed to cancel increase order: %v", err)
		}
		return fmt.Errorf("increase order did not fill within timeout")
	}
	
	// Step 2: Now close the full position (which is now larger than minimum)
	closeSide := "buy"
	closeSize := newPosition.Abs() // Size of position to close
	if currentPosition.IsNegative() {
		closeSide = "buy" // Buy to close short
	} else {
		closeSide = "sell" // Sell to close long
	}
	
	log.Printf("Step 2: Closing full position of %s ETH with %s order", closeSize.StringFixed(4), closeSide)
	
	// Refresh orderbook
	orderBook, err = gh.exchange.GetOrderBook(instrument)
	if err != nil {
		return fmt.Errorf("failed to get orderbook for close: %w", err)
	}
	
	// Calculate aggressive price for close order
	var closePrice decimal.Decimal
	if closeSide == "buy" {
		if len(orderBook.Asks) == 0 {
			return fmt.Errorf("no asks in orderbook")
		}
		// Price above ask to ensure fill
		closePrice = orderBook.Asks[0].Price.Mul(decimal.NewFromFloat(1.001))
	} else {
		if len(orderBook.Bids) == 0 {
			return fmt.Errorf("no bids in orderbook")
		}
		// Price below bid to ensure fill
		closePrice = orderBook.Bids[0].Price.Mul(decimal.NewFromFloat(0.999))
	}
	
	closePrice = closePrice.Div(tickSize).Round(0).Mul(tickSize)
	
	// Place the close order
	log.Printf("Placing close order: %s %s @ %s", closeSide, closeSize.StringFixed(4), closePrice.StringFixed(2))
	closeOrderID, err := gh.exchange.PlaceLimitOrder(instrument, closeSide, closePrice, closeSize)
	if err != nil {
		return fmt.Errorf("failed to place close order: %w", err)
	}
	
	log.Printf("Close order placed: %s", closeOrderID)
	log.Printf("Successfully executed minimum size close strategy!")
	
	// Calculate the cost of this operation (spread loss on the extra size we had to trade)
	extraSize := minOrderSize // We traded an extra 0.1 ETH (increased then closed)
	estimatedCost := extraSize.Mul(closePrice.Sub(increasePrice).Abs())
	log.Printf("Estimated cost of this operation: ~$%s (spread loss on %s ETH extra volume)", 
		estimatedCost.StringFixed(2), extraSize.StringFixed(4))
	
	return nil
}

// isOptionInstrument checks if an instrument is an option
func isOptionInstrument(instrument string) bool {
	// Simple check - options usually end with -C or -P
	return len(instrument) > 2 && 
		(instrument[len(instrument)-2:] == "-C" || instrument[len(instrument)-2:] == "-P")
}