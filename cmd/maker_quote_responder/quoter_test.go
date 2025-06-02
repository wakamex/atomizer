package main

import (
	"fmt"
	"testing"
	"time"

	ccxt "github.com/ccxt/ccxt/go/v4"
	"github.com/stretchr/testify/assert"
)

func TestDeribitOptionPricing(t *testing.T) {
	// Initialize CCXT exchange
	exchange := ccxt.NewDeribit(map[string]interface{}{
		"rateLimit":       10,
		"enableRateLimit": true,
	})

	// Test with a real ETH option
	// Let's use an ETH option with a future expiry
	testCases := []struct {
		name       string
		instrument string
		wantError  bool
	}{
		{
			name:       "ETH Call Option",
			instrument: "ETH-30MAY25-3000-C",
			wantError:  false,
		},
		{
			name:       "ETH Perpetual",
			instrument: "ETH-PERPETUAL",
			wantError:  false,
		},
		{
			name:       "Invalid ETH-PERP",
			instrument: "ETH-PERP",
			wantError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Fetch ticker
			ticker, err := exchange.FetchTicker(tc.instrument)
			if tc.wantError {
				if err == nil {
					t.Errorf("Expected error for %s, but got none", tc.instrument)
				}
				t.Logf("Expected error for %s: %v", tc.instrument, err)
				return
			}

			if err != nil {
				t.Errorf("Failed to fetch ticker for %s: %v", tc.instrument, err)
				return
			}

			// Log ticker details
			t.Logf("Ticker for %s:", tc.instrument)
			if ticker.Bid != nil {
				t.Logf("  Bid: %f", *ticker.Bid)
			}
			if ticker.Ask != nil {
				t.Logf("  Ask: %f", *ticker.Ask)
			}
			if ticker.Last != nil {
				t.Logf("  Last: %f", *ticker.Last)
			}

			// Check if Info contains underlying price for options
			if ticker.Info != nil {
				t.Logf("  Info type: %T", ticker.Info)
				t.Logf("  Info keys: %v", getMapKeys(ticker.Info))
				if underlyingPrice, exists := ticker.Info["underlying_price"]; exists {
					t.Logf("  Underlying price: %v (type: %T)", underlyingPrice, underlyingPrice)
				}
				if indexPrice, exists := ticker.Info["index_price"]; exists {
					t.Logf("  Index price: %v (type: %T)", indexPrice, indexPrice)
				}
			}

			// Fetch order book
			orderBook, err := exchange.FetchOrderBook(tc.instrument)
			if err != nil {
				t.Errorf("Failed to fetch order book for %s: %v", tc.instrument, err)
				return
			}

			t.Logf("Order book for %s:", tc.instrument)
			if len(orderBook.Bids) > 0 {
				t.Logf("  Best bid: Price=%f, Size=%f", orderBook.Bids[0][0], orderBook.Bids[0][1])
			}
			if len(orderBook.Asks) > 0 {
				t.Logf("  Best ask: Price=%f, Size=%f", orderBook.Asks[0][0], orderBook.Asks[0][1])
			}
		})
	}
}

func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestManualDeribitAPI(t *testing.T) {
	// This test manually checks what Deribit returns for options
	t.Skip("Manual test - uncomment to run")

	exchange := ccxt.NewDeribit(map[string]interface{}{
		"rateLimit":       10,
		"enableRateLimit": true,
	})

	// Test a specific option
	instrument := "ETH-28MAR25-5000-C"

	ticker, err := exchange.FetchTicker(instrument)
	if err != nil {
		t.Fatalf("Failed to fetch ticker: %v", err)
	}

	fmt.Printf("Raw ticker response for %s:\n", instrument)
	fmt.Printf("Bid: %v\n", ticker.Bid)
	fmt.Printf("Ask: %v\n", ticker.Ask)
	fmt.Printf("Last: %v\n", ticker.Last)
	fmt.Printf("Info: %+v\n", ticker.Info)
}

func TestMakeQuote(t *testing.T) {
	mockExchange := &MockExchange{}
	a, err := MakeQuote(RFQResult{
		Asset:      "0xb67bfa7b488df4f2efa874f4e59242e9130ae61f",
		Strike:     "260000000000",
		Expiry:     time.Date(2025, 6, 28, 0, 0, 0, 0, time.UTC).Unix(),
		IsPut:      false,
		Quantity:   "1000000000000000000",
		IsTakerBuy: false,
	}, "ETH", "test-nonce", &AppConfig{}, mockExchange)
	assert.NoError(t, err)
	assert.Equal(t, a.Strike, "260000000000")
	assert.Equal(t, a.Expiry, time.Date(2025, 6, 28, 0, 0, 0, 0, time.UTC).Unix())
	assert.Equal(t, a.IsPut, false)
	assert.Equal(t, a.Quantity, "1000000000000000000")
	assert.Equal(t, a.IsTakerBuy, false)
	assert.Equal(t, a.AssetAddress, "0x0000000000000000000000000000000000000000")
	assert.Equal(t, a.ChainID, uint64(1))
}
