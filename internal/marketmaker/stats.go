package marketmaker

import (
	"log"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wakamex/atomizer/internal/types"
)

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
			mm.reportStats(startTime)
		}
	}
}

// reportStats generates and logs statistics
func (mm *MarketMaker) reportStats(startTime time.Time) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

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

	// Consistency check
	if totalOrders != len(mm.activeOrders) {
		log.Printf("WARNING: Order tracking inconsistency: %d in activeOrders, %d in ordersByInstrument",
			len(mm.activeOrders), totalOrders)
	}

	// Log stats
	log.Printf("Stats: Orders=%d/%d/%d (placed/cancelled/filled), Active=%d/%d instruments, Uptime=%ds",
		mm.stats.OrdersPlaced,
		mm.stats.OrdersCancelled,
		mm.stats.OrdersFilled,
		activeCount,
		len(mm.config.Instruments),
		mm.stats.UptimeSeconds)

	// Detailed order state in debug mode
	if debugMode {
		mm.logDetailedOrderState()
	}
}

// logDetailedOrderState logs detailed order information (debug mode only)
func (mm *MarketMaker) logDetailedOrderState() {
	for instrument, orders := range mm.ordersByInstrument {
		if len(orders) > 0 {
			bidPrice, askPrice := "none", "none"
			if bidOrder, hasBid := orders["buy"]; hasBid {
				bidPrice = bidOrder.Price.String()
			}
			if askOrder, hasAsk := orders["sell"]; hasAsk {
				askPrice = askOrder.Price.String()
			}
			DebugLog("  %s: bid=%s, ask=%s", instrument, bidPrice, askPrice)
		}
	}
}

// getStats returns a copy of current statistics
func (mm *MarketMaker) GetStats() types.MarketMakerStats {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	// Create a copy of stats
	statsCopy := mm.stats

	// Copy bid-ask spreads map
	statsCopy.BidAskSpread = make(map[string]decimal.Decimal)
	for k, v := range mm.stats.BidAskSpread {
		statsCopy.BidAskSpread[k] = v
	}

	return statsCopy
}

// updateBidAskSpread updates the recorded spread for an instrument
func (mm *MarketMaker) updateBidAskSpread(instrument string, spread decimal.Decimal) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.stats.BidAskSpread[instrument] = spread
}
