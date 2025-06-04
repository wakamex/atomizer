package marketmaker

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/wakamex/atomizer/internal/types"
)

func TestCalculateQuotes(t *testing.T) {
	tests := []struct {
		name           string
		aggression     decimal.Decimal
		improvement    decimal.Decimal
		bestBid        decimal.Decimal
		bestAsk        decimal.Decimal
		expectedBidMin decimal.Decimal
		expectedBidMax decimal.Decimal
		expectedAskMin decimal.Decimal
		expectedAskMax decimal.Decimal
		description    string
	}{
		{
			name:           "Conservative_0.0_JoinBest",
			aggression:     decimal.NewFromFloat(0.0),
			improvement:    decimal.NewFromFloat(0.1),
			bestBid:        decimal.NewFromFloat(100),
			bestAsk:        decimal.NewFromFloat(101),
			expectedBidMin: decimal.NewFromFloat(100),
			expectedBidMax: decimal.NewFromFloat(100),
			expectedAskMin: decimal.NewFromFloat(101),
			expectedAskMax: decimal.NewFromFloat(101),
			description:    "Aggression 0.0 should join best bid/ask",
		},
		{
			name:           "Conservative_0.5_Halfway",
			aggression:     decimal.NewFromFloat(0.5),
			improvement:    decimal.NewFromFloat(0.1),
			bestBid:        decimal.NewFromFloat(100),
			bestAsk:        decimal.NewFromFloat(102),
			expectedBidMin: decimal.NewFromFloat(100.5),
			expectedBidMax: decimal.NewFromFloat(100.5),
			expectedAskMin: decimal.NewFromFloat(101.5),
			expectedAskMax: decimal.NewFromFloat(101.5),
			description:    "Aggression 0.5 should place orders halfway to mid",
		},
		{
			name:           "Conservative_0.9_NearMid",
			aggression:     decimal.NewFromFloat(0.9),
			improvement:    decimal.NewFromFloat(0.1),
			bestBid:        decimal.NewFromFloat(100),
			bestAsk:        decimal.NewFromFloat(102),
			expectedBidMin: decimal.NewFromFloat(100.9),
			expectedBidMax: decimal.NewFromFloat(100.9),
			expectedAskMin: decimal.NewFromFloat(101.1),
			expectedAskMax: decimal.NewFromFloat(101.1),
			description:    "Aggression 0.9 should place orders very close to mid",
		},
		{
			name:           "Aggressive_1.0_CrossSpread",
			aggression:     decimal.NewFromFloat(1.0),
			improvement:    decimal.NewFromFloat(0.1),
			bestBid:        decimal.NewFromFloat(100),
			bestAsk:        decimal.NewFromFloat(101),
			expectedBidMin: decimal.NewFromFloat(100.1),
			expectedBidMax: decimal.NewFromFloat(100.1),
			expectedAskMin: decimal.NewFromFloat(100.9),
			expectedAskMax: decimal.NewFromFloat(100.9),
			description:    "Aggression 1.0 should cross spread with improvement",
		},
		{
			name:           "Aggressive_2.0_CrossSpread",
			aggression:     decimal.NewFromFloat(2.0),
			improvement:    decimal.NewFromFloat(0.2),
			bestBid:        decimal.NewFromFloat(100),
			bestAsk:        decimal.NewFromFloat(101),
			expectedBidMin: decimal.NewFromFloat(100.2),
			expectedBidMax: decimal.NewFromFloat(100.2),
			expectedAskMin: decimal.NewFromFloat(100.8),
			expectedAskMax: decimal.NewFromFloat(100.8),
			description:    "Aggression > 1.0 should still use improvement logic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create market maker with test config
			config := &types.MarketMakerConfig{
				Aggression:   tt.aggression,
				Improvement:  tt.improvement,
				MinSpreadBps: 1, // 0.01%
			}
			
			mm := &MarketMaker{
				config: config,
			}

			// Create ticker with test data
			ticker := &types.TickerUpdate{
				BestBid: tt.bestBid,
				BestAsk: tt.bestAsk,
			}

			// Calculate quotes
			bidPrice, askPrice := mm.calculateQuotes(ticker, nil)

			// Assert results
			assert.True(t, bidPrice.GreaterThanOrEqual(tt.expectedBidMin), 
				"%s: Bid price %s should be >= %s", tt.description, bidPrice, tt.expectedBidMin)
			assert.True(t, bidPrice.LessThanOrEqual(tt.expectedBidMax), 
				"%s: Bid price %s should be <= %s", tt.description, bidPrice, tt.expectedBidMax)
			assert.True(t, askPrice.GreaterThanOrEqual(tt.expectedAskMin), 
				"%s: Ask price %s should be >= %s", tt.description, askPrice, tt.expectedAskMin)
			assert.True(t, askPrice.LessThanOrEqual(tt.expectedAskMax), 
				"%s: Ask price %s should be <= %s", tt.description, askPrice, tt.expectedAskMax)
			
			// Ensure minimum spread is maintained
			spread := askPrice.Sub(bidPrice)
			assert.True(t, spread.GreaterThan(decimal.Zero), 
				"%s: Spread should be positive", tt.description)
		})
	}
}

func TestConservativeModeSpreadConstraints(t *testing.T) {
	// Test that conservative mode never crosses the spread
	bestBid := decimal.NewFromFloat(100)
	bestAsk := decimal.NewFromFloat(101)
	mid := bestBid.Add(bestAsk).Div(decimal.NewFromFloat(2)) // 100.5

	aggressionLevels := []float64{0.0, 0.1, 0.3, 0.5, 0.7, 0.9}
	
	for _, aggression := range aggressionLevels {
		t.Run(fmt.Sprintf("Aggression_%.1f", aggression), func(t *testing.T) {
			config := &types.MarketMakerConfig{
				Aggression:   decimal.NewFromFloat(aggression),
				Improvement:  decimal.NewFromFloat(0.1), // Should be ignored in conservative mode
				MinSpreadBps: 1,
			}
			
			mm := &MarketMaker{
				config: config,
			}

			ticker := &types.TickerUpdate{
				BestBid: bestBid,
				BestAsk: bestAsk,
			}

			bidPrice, askPrice := mm.calculateQuotes(ticker, nil)

			// In conservative mode:
			// - Bid should never be higher than mid
			// - Ask should never be lower than mid
			assert.True(t, bidPrice.LessThanOrEqual(mid), 
				"Conservative bid should not exceed mid (bid: %s, mid: %s)", bidPrice, mid)
			assert.True(t, askPrice.GreaterThanOrEqual(mid), 
				"Conservative ask should not go below mid (ask: %s, mid: %s)", askPrice, mid)
			
			// Bids should be between bestBid and mid
			assert.True(t, bidPrice.GreaterThanOrEqual(bestBid), 
				"Bid should be >= best bid")
			assert.True(t, bidPrice.LessThanOrEqual(mid), 
				"Bid should be <= mid")
			
			// Asks should be between mid and bestAsk
			assert.True(t, askPrice.GreaterThanOrEqual(mid), 
				"Ask should be >= mid")
			assert.True(t, askPrice.LessThanOrEqual(bestAsk), 
				"Ask should be <= best ask")
		})
	}
}

func TestAggressionClampingAndEdgeCases(t *testing.T) {
	tests := []struct {
		name              string
		aggression        decimal.Decimal
		expectedClamped   decimal.Decimal
		shouldBeAggressive bool
	}{
		{
			name:              "Negative_Aggression_Clamped",
			aggression:        decimal.NewFromFloat(-0.5),
			expectedClamped:   decimal.NewFromFloat(0.0),
			shouldBeAggressive: false,
		},
		{
			name:              "Above_0.9_Clamped",
			aggression:        decimal.NewFromFloat(0.95),
			expectedClamped:   decimal.NewFromFloat(0.9),
			shouldBeAggressive: false,
		},
		{
			name:              "Exactly_1.0_Aggressive",
			aggression:        decimal.NewFromFloat(1.0),
			expectedClamped:   decimal.NewFromFloat(1.0),
			shouldBeAggressive: true,
		},
		{
			name:              "Above_1.0_Aggressive",
			aggression:        decimal.NewFromFloat(1.5),
			expectedClamped:   decimal.NewFromFloat(1.5),
			shouldBeAggressive: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &types.MarketMakerConfig{
				Aggression:   tt.aggression,
				Improvement:  decimal.NewFromFloat(0.1),
				MinSpreadBps: 1,
			}
			
			mm := &MarketMaker{
				config: config,
			}

			ticker := &types.TickerUpdate{
				BestBid: decimal.NewFromFloat(100),
				BestAsk: decimal.NewFromFloat(101),
			}

			bidPrice, askPrice := mm.calculateQuotes(ticker, nil)

			if tt.shouldBeAggressive {
				// In aggressive mode, bid can be > bestBid
				expectedBid := ticker.BestBid.Add(config.Improvement)
				assert.Equal(t, expectedBid.String(), bidPrice.String(), 
					"Aggressive mode should add improvement to bid")
			} else {
				// In conservative mode, verify clamping worked
				mid := ticker.BestBid.Add(ticker.BestAsk).Div(decimal.NewFromFloat(2))
				assert.True(t, bidPrice.LessThanOrEqual(mid), 
					"Conservative mode should not exceed mid")
				assert.True(t, askPrice.GreaterThanOrEqual(mid), 
					"Conservative mode should not go below mid")
			}
		})
	}
}

func TestMinimumSpreadEnforcement(t *testing.T) {
	// Test that minimum spread is always maintained
	config := &types.MarketMakerConfig{
		Aggression:   decimal.NewFromFloat(0.9), // Very aggressive conservative
		Improvement:  decimal.NewFromFloat(0.1),
		MinSpreadBps: 100, // 1% minimum spread
		SpreadBps:    50,  // 0.5% target spread
	}
	
	mm := &MarketMaker{
		config: config,
	}

	// Very tight market
	ticker := &types.TickerUpdate{
		BestBid: decimal.NewFromFloat(100),
		BestAsk: decimal.NewFromFloat(100.1), // Only 0.1% spread
	}

	bidPrice, askPrice := mm.calculateQuotes(ticker, nil)
	
	// Calculate actual spread
	spread := askPrice.Sub(bidPrice)
	mid := ticker.BestBid.Add(ticker.BestAsk).Div(decimal.NewFromFloat(2))
	minSpread := mid.Mul(decimal.NewFromInt(int64(config.MinSpreadBps)).Div(decimal.NewFromInt(10000)))
	
	assert.True(t, spread.GreaterThanOrEqual(minSpread), 
		"Spread %s should be >= minimum spread %s", spread, minSpread)
}

func TestZeroPriceHandling(t *testing.T) {
	// Test handling when one side has no price
	config := &types.MarketMakerConfig{
		Aggression:   decimal.NewFromFloat(0.5),
		Improvement:  decimal.NewFromFloat(0.1),
		MinSpreadBps: 10,
		SpreadBps:    50,
	}
	
	mm := &MarketMaker{
		config: config,
	}

	tests := []struct {
		name     string
		bestBid  decimal.Decimal
		bestAsk  decimal.Decimal
		markPrice decimal.Decimal
	}{
		{
			name:      "Zero_Bid",
			bestBid:   decimal.Zero,
			bestAsk:   decimal.NewFromFloat(101),
			markPrice: decimal.NewFromFloat(100),
		},
		{
			name:      "Zero_Ask",
			bestBid:   decimal.NewFromFloat(99),
			bestAsk:   decimal.Zero,
			markPrice: decimal.NewFromFloat(100),
		},
		{
			name:      "Both_Zero",
			bestBid:   decimal.Zero,
			bestAsk:   decimal.Zero,
			markPrice: decimal.NewFromFloat(100),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ticker := &types.TickerUpdate{
				BestBid:   tt.bestBid,
				BestAsk:   tt.bestAsk,
				MarkPrice: tt.markPrice,
			}

			bidPrice, askPrice := mm.calculateQuotes(ticker, nil)
			
			// Should produce valid quotes
			assert.True(t, bidPrice.GreaterThan(decimal.Zero), 
				"Bid price should be positive")
			assert.True(t, askPrice.GreaterThan(decimal.Zero), 
				"Ask price should be positive")
			assert.True(t, askPrice.GreaterThan(bidPrice), 
				"Ask should be greater than bid")
		})
	}
}