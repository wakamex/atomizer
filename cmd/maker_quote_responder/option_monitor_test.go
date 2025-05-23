package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

// TestOptionMonitor provides comprehensive monitoring tools for Deribit options
func TestOptionMonitor(t *testing.T) {
	exchange := ccxt.NewDeribit(map[string]interface{}{
		"rateLimit":       10,
		"enableRateLimit": true,
	})

	t.Run("MarketOverview", func(t *testing.T) {
		// Fetch ETH options
		ethTickers, err := exchange.FetchTickers(ccxt.WithFetchTickersParams(map[string]interface{}{
			"code": "ETH",
		}))
		if err != nil {
			t.Fatalf("Failed to fetch ETH tickers: %v", err)
		}

		// Get ETH index price
		perpTicker, err := exchange.FetchTicker("ETH-PERPETUAL")
		if err != nil {
			t.Logf("Failed to fetch ETH-PERPETUAL: %v", err)
		}
		
		var ethIndexPrice float64
		if perpTicker.Last != nil {
			ethIndexPrice = *perpTicker.Last
			t.Logf("ETH Index Price: $%.2f", ethIndexPrice)
		}

		// Categorize options
		var calls, puts []string
		expiryMap := make(map[string]int)
		
		for symbol := range ethTickers.Tickers {
			cleanSymbol := symbol
			if idx := strings.Index(symbol, ":"); idx >= 0 {
				cleanSymbol = symbol[idx+1:]
			}
			
			if strings.HasSuffix(cleanSymbol, "-C") {
				calls = append(calls, cleanSymbol)
			} else if strings.HasSuffix(cleanSymbol, "-P") {
				puts = append(puts, cleanSymbol)
			}
			
			// Extract expiry
			parts := strings.Split(cleanSymbol, "-")
			if len(parts) >= 2 {
				expiryMap[parts[1]]++
			}
		}

		t.Logf("Market Overview:")
		t.Logf("  Total Options: %d", len(calls)+len(puts))
		t.Logf("  Calls: %d, Puts: %d", len(calls), len(puts))
		t.Logf("  Unique Expiries: %d", len(expiryMap))
		
		// Show nearest expiries
		var expiries []string
		for exp := range expiryMap {
			expiries = append(expiries, exp)
		}
		sort.Strings(expiries)
		
		t.Logf("  Next 5 Expiries:")
		for i := 0; i < 5 && i < len(expiries); i++ {
			t.Logf("    %s (%d options)", expiries[i], expiryMap[expiries[i]])
		}
	})

	t.Run("AnalyzeSpecificOptions", func(t *testing.T) {
		// Analyze options at different strikes and expiries
		testOptions := []struct {
			symbol   string
			expiry   time.Time
			strike   float64
			analysis string
		}{
			{"ETH-250530-2500-C", time.Date(2025, 5, 30, 0, 0, 0, 0, time.UTC), 2500, "Near ATM"},
			{"ETH-250530-3000-C", time.Date(2025, 5, 30, 0, 0, 0, 0, time.UTC), 3000, "OTM"},
			{"ETH-250627-3000-C", time.Date(2025, 6, 27, 0, 0, 0, 0, time.UTC), 3000, "1 month OTM"},
			{"ETH-250926-3000-C", time.Date(2025, 9, 26, 0, 0, 0, 0, time.UTC), 3000, "4 month OTM"},
		}

		// Get ETH price
		perpTicker, _ := exchange.FetchTicker("ETH-PERPETUAL")
		var ethPrice float64
		if perpTicker.Last != nil {
			ethPrice = *perpTicker.Last
		}

		for _, opt := range testOptions {
			fullSymbol := "ETH/USD:" + opt.symbol
			ticker, err := exchange.FetchTicker(fullSymbol)
			if err != nil {
				t.Logf("%s: Failed to fetch - %v", opt.symbol, err)
				continue
			}

			t.Logf("\n%s (%s):", opt.symbol, opt.analysis)
			
			// Calculate days to expiry
			daysToExpiry := opt.expiry.Sub(time.Now()).Hours() / 24
			t.Logf("  Days to Expiry: %.1f", daysToExpiry)
			
			// Price information
			if ticker.Bid != nil && ticker.Ask != nil {
				midPrice := (*ticker.Bid + *ticker.Ask) / 2
				spread := *ticker.Ask - *ticker.Bid
				spreadPct := (spread / midPrice) * 100
				
				t.Logf("  Bid: %.4f ETH ($%.2f)", *ticker.Bid, *ticker.Bid*ethPrice)
				t.Logf("  Ask: %.4f ETH ($%.2f)", *ticker.Ask, *ticker.Ask*ethPrice)
				t.Logf("  Mid: %.4f ETH ($%.2f)", midPrice, midPrice*ethPrice)
				t.Logf("  Spread: %.4f ETH (%.1f%%)", spread, spreadPct)
				
				// Moneyness
				moneyness := ethPrice / opt.strike
				t.Logf("  Moneyness: %.3f", moneyness)
				
				// Price as % of spot
				priceAsPctOfSpot := (midPrice * ethPrice / ethPrice) * 100
				t.Logf("  Price as %% of Spot: %.2f%%", priceAsPctOfSpot)
			}
			
			if ticker.Last != nil {
				t.Logf("  Last: %.4f ETH ($%.2f)", *ticker.Last, *ticker.Last*ethPrice)
			}
		}
	})

	t.Run("FindBestOpportunities", func(t *testing.T) {
		// Find options with specific characteristics
		ethTickers, err := exchange.FetchTickers(ccxt.WithFetchTickersParams(map[string]interface{}{
			"code": "ETH",
		}))
		if err != nil {
			t.Skip("Failed to fetch tickers")
		}

		type opportunity struct {
			symbol    string
			bidAsk    string
			spread    float64
			daysToExp float64
		}

		var tightSpreads []opportunity
		var wideSpreads []opportunity
		var nearExpiry []opportunity

		for symbol, ticker := range ethTickers.Tickers {
			cleanSymbol := symbol
			if idx := strings.Index(symbol, ":"); idx >= 0 {
				cleanSymbol = symbol[idx+1:]
			}

			// Only analyze calls
			if !strings.HasSuffix(cleanSymbol, "-C") {
				continue
			}

			// Calculate spread
			if ticker.Bid != nil && ticker.Ask != nil && *ticker.Bid > 0 {
				mid := (*ticker.Bid + *ticker.Ask) / 2
				spread := (*ticker.Ask - *ticker.Bid) / mid * 100
				
				// Parse expiry
				parts := strings.Split(cleanSymbol, "-")
				if len(parts) >= 2 && len(parts[1]) == 6 {
					year, _ := fmt.Sscanf(parts[1][:2], "%d")
					month, _ := fmt.Sscanf(parts[1][2:4], "%d")
					day, _ := fmt.Sscanf(parts[1][4:6], "%d")
					expiry := time.Date(2000+year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
					daysToExp := expiry.Sub(time.Now()).Hours() / 24
					
					opp := opportunity{
						symbol:    cleanSymbol,
						bidAsk:    fmt.Sprintf("%.4f/%.4f", *ticker.Bid, *ticker.Ask),
						spread:    spread,
						daysToExp: daysToExp,
					}
					
					if spread < 2 && daysToExp > 1 {
						tightSpreads = append(tightSpreads, opp)
					}
					if spread > 10 && daysToExp > 1 {
						wideSpreads = append(wideSpreads, opp)
					}
					if daysToExp > 0 && daysToExp < 7 {
						nearExpiry = append(nearExpiry, opp)
					}
				}
			}
		}

		// Sort by spread
		sort.Slice(tightSpreads, func(i, j int) bool {
			return tightSpreads[i].spread < tightSpreads[j].spread
		})
		sort.Slice(wideSpreads, func(i, j int) bool {
			return wideSpreads[i].spread > wideSpreads[j].spread
		})
		sort.Slice(nearExpiry, func(i, j int) bool {
			return nearExpiry[i].daysToExp < nearExpiry[j].daysToExp
		})

		t.Log("\nTightest Spreads (< 2%):")
		for i := 0; i < 5 && i < len(tightSpreads); i++ {
			o := tightSpreads[i]
			t.Logf("  %s: %s ETH, Spread: %.1f%%, Days: %.1f", 
				o.symbol, o.bidAsk, o.spread, o.daysToExp)
		}

		t.Log("\nWidest Spreads (> 10%):")
		for i := 0; i < 5 && i < len(wideSpreads); i++ {
			o := wideSpreads[i]
			t.Logf("  %s: %s ETH, Spread: %.1f%%, Days: %.1f", 
				o.symbol, o.bidAsk, o.spread, o.daysToExp)
		}

		t.Log("\nNear Expiry (< 7 days):")
		for i := 0; i < 5 && i < len(nearExpiry); i++ {
			o := nearExpiry[i]
			t.Logf("  %s: %s ETH, Spread: %.1f%%, Days: %.1f", 
				o.symbol, o.bidAsk, o.spread, o.daysToExp)
		}
	})

	t.Run("TestQuoteFlow", func(t *testing.T) {
		// Test the full quote flow with a real option
		rfq := RFQResult{
			Asset:      "0xb67bfa7b488df4f2efa874f4e59242e9130ae61f",
			ChainID:    1,
			Expiry:     time.Date(2025, 5, 30, 8, 0, 0, 0, time.UTC).Unix(),
			IsPut:      false,
			IsTakerBuy: true,
			Quantity:   "100000000000000000", // 0.1 ETH
			Strike:     "250000000000",        // 2500
		}

		// Test conversion
		instrument, err := convertOptionDetailsToInstrument("ETH", rfq.Strike, rfq.Expiry, rfq.IsPut)
		if err != nil {
			t.Errorf("Failed to convert: %v", err)
			return
		}
		t.Logf("Converted to: %s", instrument)

		// Test order book fetch
		book, err := getOrderBook(rfq, "ETH")
		if err != nil {
			t.Logf("Failed to get order book: %v", err)
			return
		}

		t.Logf("Order Book:")
		t.Logf("  Index Price: $%.2f", book.Index)
		if len(book.Bids) > 0 {
			t.Logf("  Best Bid: %.4f ETH ($%.2f)", book.Bids[0][0], book.Bids[0][0]*book.Index)
		}
		if len(book.Asks) > 0 {
			t.Logf("  Best Ask: %.4f ETH ($%.2f)", book.Asks[0][0], book.Asks[0][0]*book.Index)
		}

		// Test quote generation
		quote, apr, err := getDeribitQuote(rfq, "ETH")
		if err != nil {
			t.Logf("Failed to get quote: %v", err)
			return
		}
		
		var quoteFloat float64
		fmt.Sscanf(quote, "%f", &quoteFloat)
		t.Logf("Quote: $%.2f (APR: %.2f%%)", quoteFloat/1e8, apr)
	})
}

