package hedging

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

// PureGammaHedger implements aggressive gamma hedging for options positions
// Unlike GammaDDHAlgo, this focuses on risk reduction, not market making
type PureGammaHedger struct {
	exchange  MarketMakerExchange
	config    *AppConfig
	debugMode bool // Enable debug logging

	// Greek tracking
	positions map[string]*OptionPosition // instrument -> position with greeks
	netDelta  decimal.Decimal
	netGamma  decimal.Decimal

	// Hedging parameters
	deltaThreshold decimal.Decimal // Max delta before hedging
	minHedgeSize   decimal.Decimal // Minimum hedge size (exchange minimum)
	hedgeInterval  time.Duration   // How often to check/rehedge
	aggressiveness decimal.Decimal // How far through the spread to place orders (0-1)

	// Hedge tracking
	currentHedge     *HedgePosition // Current perp/spot hedge position
	lastHedgeTime    time.Time
	consecutiveFails int // Track consecutive hedge failures

	// Control
	ctx     context.Context
	cancel  context.CancelFunc
	running bool
	mu      sync.RWMutex
}

// NewPureGammaHedger creates a new pure gamma hedger
func NewPureGammaHedger(exchange MarketMakerExchange, config *AppConfig) *PureGammaHedger {
	ctx, cancel := context.WithCancel(context.Background())

	return &PureGammaHedger{
		exchange:       exchange,
		config:         config,
		positions:      make(map[string]*OptionPosition),
		deltaThreshold: decimal.NewFromFloat(0.1), // 0.1 ETH delta threshold
		minHedgeSize:   decimal.NewFromFloat(0.1), // Exchange minimum
		hedgeInterval:  30 * time.Second,
		aggressiveness: decimal.NewFromFloat(0.5), // Cross 50% of the spread
		ctx:            ctx,
		cancel:         cancel,
	}
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

	// Start hedging loop
	go gh.hedgingLoop()

	return nil
}

// hedgingLoop is the main hedging loop
func (gh *PureGammaHedger) hedgingLoop() {
	// Wait a bit for WebSocket to be fully ready
	time.Sleep(2 * time.Second)

	// Initial position load and hedge
	gh.loadPositions()
	gh.checkAndHedge()

	ticker := time.NewTicker(gh.hedgeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-gh.ctx.Done():
			return
		case <-ticker.C:
			gh.checkAndHedge()
		}
	}
}

// checkAndHedge checks if hedging is needed and executes
func (gh *PureGammaHedger) checkAndHedge() {
	// Reload positions and update Greeks
	gh.loadPositions()

	gh.mu.RLock()
	netDelta := gh.netDelta
	positions := len(gh.positions)
	gh.mu.RUnlock()

	if positions == 0 {
		log.Printf("No positions to hedge")
		return
	}

	// Check if we need to hedge
	if netDelta.Abs().LessThan(gh.deltaThreshold) {
		return
	}

	// Calculate hedge amount (negative of net delta to neutralize)
	hedgeAmount := netDelta.Neg()

	// Check minimum size
	if hedgeAmount.Abs().LessThan(gh.minHedgeSize) {
		log.Printf("Hedge amount %s below minimum %s, skipping",
			hedgeAmount.StringFixed(4), gh.minHedgeSize.StringFixed(4))
		return
	}

	// Execute aggressive hedge
	gh.executeAggressiveHedge(hedgeAmount)
}

// executeAggressiveHedge places orders that will actually fill
func (gh *PureGammaHedger) executeAggressiveHedge(amount decimal.Decimal) {
	log.Printf("Executing AGGRESSIVE hedge: %s ETH", amount.StringFixed(4))

	// Get current orderbook for ETH-PERP
	orderbook, err := gh.exchange.GetOrderBook("ETH-PERP")
	if err != nil {
		log.Printf("Failed to get orderbook: %v", err)
		gh.consecutiveFails++
		return
	}

	// Ensure we have orderbook data
	if len(orderbook.Bids) == 0 || len(orderbook.Asks) == 0 {
		log.Printf("Empty orderbook, cannot hedge")
		gh.consecutiveFails++
		return
	}

	bestBid := orderbook.Bids[0].Price
	bestAsk := orderbook.Asks[0].Price
	spread := bestAsk.Sub(bestBid)

	log.Printf("ETH-PERP Market: Bid=%s, Ask=%s, Spread=%s",
		bestBid.StringFixed(2), bestAsk.StringFixed(2), spread.StringFixed(2))

	var price decimal.Decimal
	var side string

	if amount.IsNegative() {
		// Need to SELL to hedge positive delta
		side = "sell"
		// Price = bid - (spread * (1 - aggressiveness))
		// aggressiveness=1 means price at bid (most aggressive)
		// aggressiveness=0 means price at ask (least aggressive, won't fill)
		price = bestBid.Sub(spread.Mul(decimal.NewFromFloat(1).Sub(gh.aggressiveness)))
	} else {
		// Need to BUY to hedge negative delta
		side = "buy"
		// Price = ask + (spread * (1 - aggressiveness))
		price = bestAsk.Add(spread.Mul(decimal.NewFromFloat(1).Sub(gh.aggressiveness)))
	}

	// Round to tick size
	tickSize := decimal.NewFromFloat(0.01) // ETH-PERP tick size
	price = price.Div(tickSize).Round(0).Mul(tickSize)

	// Round amount to ETH-PERP's amount step (0.01)
	amountStep := decimal.NewFromFloat(0.01)
	roundedAmount := amount.Abs().Div(amountStep).Round(0).Mul(amountStep)

	log.Printf("Placing %s order: %s ETH @ %s (IOC)",
		side, roundedAmount.StringFixed(2), price.StringFixed(2))

	// Execute via exchange using PlaceLimitOrder
	orderID, err := gh.exchange.PlaceLimitOrder("ETH-PERP", side, price, roundedAmount)

	if err != nil {
		log.Printf("Hedge order failed: %v", err)
		gh.consecutiveFails++

		// If we've failed too many times, try market order
		if gh.consecutiveFails >= 3 {
			log.Printf("Multiple failures, attempting MARKET order")
			gh.executeMarketHedge(amount)
		}
		return
	}

	log.Printf("Hedge order placed successfully: %s", orderID)
	gh.consecutiveFails = 0
	gh.lastHedgeTime = time.Now()

	// Update hedge tracking
	gh.mu.Lock()
	if gh.currentHedge == nil {
		gh.currentHedge = &HedgePosition{
			Instrument: "ETH-PERP",
			Amount:     decimal.Zero,
		}
	}
	if side == "buy" {
		gh.currentHedge.Amount = gh.currentHedge.Amount.Add(amount.Abs())
	} else {
		gh.currentHedge.Amount = gh.currentHedge.Amount.Sub(amount.Abs())
	}
	gh.currentHedge.LastHedgeTime = time.Now()
	gh.currentHedge.LastHedgeDelta = gh.netDelta
	gh.mu.Unlock()
}

// executeMarketHedge uses market orders as last resort
func (gh *PureGammaHedger) executeMarketHedge(amount decimal.Decimal) {
	side := "sell"
	if amount.IsPositive() {
		side = "buy"
	}

	log.Printf("Placing MARKET order: %s %s ETH", side, amount.Abs().StringFixed(4))

	// For market orders, we can use a very aggressive limit order
	// Get orderbook first
	orderbook, err := gh.exchange.GetOrderBook("ETH-PERP")
	if err != nil {
		log.Printf("Failed to get orderbook for market order: %v", err)
		return
	}

	var marketPrice decimal.Decimal
	if side == "sell" {
		// Use a price well below bid to ensure fill
		marketPrice = orderbook.Bids[0].Price.Mul(decimal.NewFromFloat(0.99))
	} else {
		// Use a price well above ask to ensure fill
		marketPrice = orderbook.Asks[0].Price.Mul(decimal.NewFromFloat(1.01))
	}

	// Round amount
	amountStep := decimal.NewFromFloat(0.01)
	roundedAmount := amount.Abs().Div(amountStep).Round(0).Mul(amountStep)

	orderID, err := gh.exchange.PlaceLimitOrder("ETH-PERP", side, marketPrice, roundedAmount)
	if err != nil {
		log.Printf("Market order also failed: %v", err)
		return
	}

	log.Printf("Market hedge executed: %s", orderID)
}

// loadPositions loads current positions and updates Greeks
func (gh *PureGammaHedger) loadPositions() {
	positions, err := gh.exchange.GetPositions()
	if err != nil {
		log.Printf("Failed to load positions: %v", err)
		return
	}

	if gh.debugMode {
		log.Printf("Loaded %d positions from exchange", len(positions))
	}

	gh.mu.Lock()
	defer gh.mu.Unlock()

	// Clear old positions
	gh.positions = make(map[string]*OptionPosition)
	gh.netDelta = decimal.Zero
	gh.netGamma = decimal.Zero

	// Track subtotals
	optionsDelta := decimal.Zero
	perpsDelta := decimal.Zero

	for _, position := range positions {
		amount := decimal.NewFromFloat(position.Amount)
		instrument := position.InstrumentName

		if gh.debugMode {
			log.Printf("Found position: %s, amount=%f, direction=%s, avgPrice=%f",
				instrument, position.Amount, position.Direction, position.AveragePrice)
		}

		if amount.IsZero() {
			continue
		}

		// Handle non-options (ETH-PERP)
		if !isOption(instrument) {
			if instrument == "ETH-PERP" {
				// Track our perp position
				if gh.currentHedge == nil {
					gh.currentHedge = &HedgePosition{
						Instrument: "ETH-PERP",
						Amount:     decimal.Zero,
					}
				}
				gh.currentHedge.Amount = amount
				// PERP has delta of 1 per contract
				gh.netDelta = gh.netDelta.Add(amount)
				perpsDelta = perpsDelta.Add(amount)
				if gh.debugMode {
					log.Printf("PERP %s: Amount=%s, Delta=%s (1:1)",
						instrument, amount.StringFixed(4), amount.StringFixed(4))
				}
			}
			continue
		}

		// Create position and fetch Greeks
		pos := &OptionPosition{
			Instrument: instrument,
			Amount:     amount,
		}

		// Greeks are already in the position data from Derive
		// We can fetch fresh ones or use what we have
		pos.Delta = decimal.NewFromFloat(0.5)  // Default
		pos.Gamma = decimal.NewFromFloat(0.01) // Default

		// Try to get Greeks from the position data (they were in the debug output)
		// For now, let's fetch fresh Greeks to ensure they're up to date
		if ticker, err := FetchDeriveTicker(instrument); err == nil {
			pos.Delta = decimal.NewFromFloat(ticker.GetDelta())
			pos.Gamma = decimal.NewFromFloat(ticker.GetGamma())
			pos.UnderlyingPrice = decimal.NewFromFloat(ticker.GetIndexPrice())

			if gh.debugMode {
				log.Printf("Option %s: Amount=%s, Delta=%s, Gamma=%s",
					instrument, amount.StringFixed(4),
					pos.Delta.StringFixed(4), pos.Gamma.StringFixed(6))
			}
		} else {
			log.Printf("Failed to fetch Greeks for %s: %v", instrument, err)
		}

		gh.positions[instrument] = pos

		// Update net Greeks
		positionDelta := amount.Mul(pos.Delta)
		gh.netDelta = gh.netDelta.Add(positionDelta)
		gh.netGamma = gh.netGamma.Add(amount.Mul(pos.Gamma))
		optionsDelta = optionsDelta.Add(positionDelta)
	}

	// Hedge position delta is already added above when processing ETH-PERP

	// Log subtotals and total with threshold status
	if !optionsDelta.IsZero() || !perpsDelta.IsZero() {
		status := "within threshold, no hedge needed"
		if gh.netDelta.Abs().GreaterThanOrEqual(gh.deltaThreshold) {
			status = "exceeds threshold, hedging required"
		}
		log.Printf("Portfolio Delta: Options=%s, Perps=%s, Total=%s ETH (%s)",
			optionsDelta.StringFixed(4), perpsDelta.StringFixed(4), gh.netDelta.StringFixed(4), status)
	}
}

// Stop stops the hedger
func (gh *PureGammaHedger) Stop() {
	gh.cancel()
	gh.mu.Lock()
	gh.running = false
	gh.mu.Unlock()
	log.Println("Pure gamma hedger stopped")
}

// Helper function to check if instrument is an option
func isOption(instrument string) bool {
	// ETH-YYYYMMDD-STRIKE-C/P format
	return len(instrument) > 10 && (instrument[len(instrument)-1] == 'C' || instrument[len(instrument)-1] == 'P')
}
