package marketmaker

import (
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wakamex/atomizer/internal/types"
)

// updateQuotesForInstrument updates quotes for a specific instrument
func (mm *MarketMaker) UpdateQuotesForInstrument(instrument string) error {
	// Prevent concurrent updates
	if lock, exists := mm.updateLocks[instrument]; exists {
		lock.Lock()
		defer lock.Unlock()
	}

	// Rate limiting check
	mm.mu.RLock()
	lastUpdate, exists := mm.lastUpdateTime[instrument]
	mm.mu.RUnlock()

	if exists && time.Since(lastUpdate) < 2*time.Second {
		return nil
	}

	mm.mu.Lock()
	mm.lastUpdateTime[instrument] = time.Now()
	mm.mu.Unlock()

	// Get ticker data
	mm.mu.RLock()
	ticker, exists := mm.latestTickers[instrument]
	mm.mu.RUnlock()

	if !exists || ticker == nil {
		return fmt.Errorf("no ticker data for %s", instrument)
	}

	// Skip if no valid price data
	if ticker.BestBid.IsZero() && ticker.BestAsk.IsZero() && ticker.MarkPrice.IsZero() {
		log.Printf("No valid price data for %s yet, skipping quote update", instrument)
		return nil
	}

	// Fetch orderbook if needed
	var orderBook *types.MarketMakerOrderBook
	if mm.config.ImprovementReferenceSize.GreaterThan(decimal.Zero) {
		var err error
		orderBook, err = mm.exchange.GetOrderBook(instrument)
		if err != nil {
			mm.logOrderbookError(instrument, err)
		} else {
			mm.clearOrderbookError(instrument)
		}
	}

	// Calculate quotes
	bidPrice, askPrice := mm.calculateQuotes(ticker, orderBook)

	// Check risk limits
	if !mm.checkRiskLimits(instrument, mm.config.QuoteSize) {
		log.Printf("Risk limits exceeded for %s, skipping quote update", instrument)
		return nil
	}

	// Update orders
	return mm.updateOrCreateOrders(instrument, bidPrice, askPrice)
}

// calculateQuotes calculates bid and ask prices based on current market
func (mm *MarketMaker) calculateQuotes(ticker *types.TickerUpdate, orderBook *types.MarketMakerOrderBook) (bidPrice, askPrice decimal.Decimal) {
	// Calculate mid price
	var midPrice decimal.Decimal
	if ticker.BestBid.IsZero() || ticker.BestAsk.IsZero() {
		if !ticker.MarkPrice.IsZero() {
			midPrice = ticker.MarkPrice
		} else {
			log.Printf("WARNING: No valid price data for %s", ticker.Instrument)
			midPrice = decimal.NewFromFloat(1.0)
		}
	} else {
		midPrice = ticker.BestBid.Add(ticker.BestAsk).Div(decimal.NewFromInt(2))
	}

	// Determine reference prices
	referenceBid := ticker.BestBid
	referenceAsk := ticker.BestAsk

	// Handle zero prices
	if referenceBid.IsZero() || referenceAsk.IsZero() {
		spreadAmount := midPrice.Mul(decimal.NewFromInt(int64(mm.config.SpreadBps)).Div(decimal.NewFromInt(10000)))
		if referenceBid.IsZero() {
			referenceBid = midPrice.Sub(spreadAmount.Div(decimal.NewFromInt(2)))
		}
		if referenceAsk.IsZero() {
			referenceAsk = midPrice.Add(spreadAmount.Div(decimal.NewFromInt(2)))
		}
	}

	// Use orderbook if reference size is set
	if orderBook != nil && mm.config.ImprovementReferenceSize.GreaterThan(decimal.Zero) {
		mm.adjustPricesForReferenceSize(orderBook, &referenceBid, &referenceAsk, midPrice)
	}

	// Calculate our quotes with improvement
	bidPrice = referenceBid.Add(mm.config.Improvement)
	askPrice = referenceAsk.Sub(mm.config.Improvement)

	// Ensure minimum spread
	minSpread := midPrice.Mul(decimal.NewFromInt(int64(mm.config.MinSpreadBps)).Div(decimal.NewFromInt(10000)))
	if askPrice.Sub(bidPrice).LessThan(minSpread) {
		bidPrice = midPrice.Sub(minSpread.Div(decimal.NewFromInt(2)))
		askPrice = midPrice.Add(minSpread.Div(decimal.NewFromInt(2)))
	}

	return bidPrice, askPrice
}

// adjustPricesForReferenceSize finds best bid/ask with sufficient size
func (mm *MarketMaker) adjustPricesForReferenceSize(orderBook *types.MarketMakerOrderBook, referenceBid, referenceAsk *decimal.Decimal, midPrice decimal.Decimal) {
	foundBid := false
	for _, bid := range orderBook.Bids {
		if bid.Size.GreaterThanOrEqual(mm.config.ImprovementReferenceSize) {
			*referenceBid = bid.Price
			foundBid = true
			break
		}
	}

	foundAsk := false
	for _, ask := range orderBook.Asks {
		if ask.Size.GreaterThanOrEqual(mm.config.ImprovementReferenceSize) {
			*referenceAsk = ask.Price
			foundAsk = true
			break
		}
	}

	// Fallback if insufficient size found
	if !foundBid || !foundAsk {
		spreadAmount := midPrice.Mul(decimal.NewFromInt(int64(mm.config.SpreadBps)).Div(decimal.NewFromInt(10000)))
		if !foundBid {
			*referenceBid = midPrice.Sub(spreadAmount.Div(decimal.NewFromInt(2)))
		}
		if !foundAsk {
			*referenceAsk = midPrice.Add(spreadAmount.Div(decimal.NewFromInt(2)))
		}
	}
}

// shouldUpdateQuotes checks if quotes need updating
func (mm *MarketMaker) shouldUpdateQuotes(instrument string) bool {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	orders, exists := mm.ordersByInstrument[instrument]
	if !exists || len(orders) == 0 {
		return true
	}

	ticker, exists := mm.latestTickers[instrument]
	if !exists {
		return false
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

// Helper functions for error logging
func (mm *MarketMaker) logOrderbookError(instrument string, err error) {
	mm.mu.Lock()
	if !mm.orderbookErrorLogged[instrument] {
		log.Printf("Failed to fetch orderbook for %s: %v, using ticker data", instrument, err)
		mm.orderbookErrorLogged[instrument] = true
	}
	mm.mu.Unlock()
}

func (mm *MarketMaker) clearOrderbookError(instrument string) {
	mm.mu.Lock()
	delete(mm.orderbookErrorLogged, instrument)
	mm.mu.Unlock()
}
