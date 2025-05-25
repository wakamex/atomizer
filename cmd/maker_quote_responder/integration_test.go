package main

import (
	"testing"
	"time"
)

func TestQuoterIntegration(t *testing.T) {
	// Create a test RFQ for a future option
	futureDate := time.Now().AddDate(0, 6, 0) // 6 months from now
	
	rfq := RFQResult{
		Asset:      "0xtest", // This won't have mapping, but we're testing the flow
		ChainID:    1,
		Expiry:     futureDate.Unix(),
		IsPut:      false,
		IsTakerBuy: true,
		Quantity:   "100000000000000000", // 0.1 * 1e18
		Strike:     "400000000000",        // 4000 * 1e8
	}

	t.Run("Test getOrderBook with ETH", func(t *testing.T) {
		book, err := getOrderBook(rfq, "ETH")
		if err != nil {
			// This might fail if the specific option doesn't exist
			t.Logf("getOrderBook failed (expected if option doesn't exist): %v", err)
			
			// But let's check if it's failing on the perpetual or the option
			if contains(err.Error(), "ETH-PERP") {
				t.Errorf("Still trying to fetch ETH-PERP instead of ETH-PERPETUAL")
			}
			return
		}
		
		t.Logf("Successfully fetched order book")
		t.Logf("Index price: %f", book.Index)
		if len(book.Bids) > 0 {
			t.Logf("Best bid: %f", book.Bids[0][0])
		}
		if len(book.Asks) > 0 {
			t.Logf("Best ask: %f", book.Asks[0][0])
		}
	})

	t.Run("Test instrument name conversion", func(t *testing.T) {
		testCases := []struct {
			name   string
			expiry int64
			strike string
			asset  string
		}{
			{
				name:   "Future ETH option",
				expiry: time.Date(2025, 12, 26, 0, 0, 0, 0, time.UTC).Unix(), // Dec 26, 2025
				strike: "500000000000", // 5000
				asset:  "ETH",
			},
			{
				name:   "Future BTC option",  
				expiry: time.Date(2025, 12, 26, 0, 0, 0, 0, time.UTC).Unix(),
				strike: "10000000000000", // 100000
				asset:  "BTC",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				instrument, err := convertOptionDetailsToInstrument(tc.asset, tc.strike, tc.expiry, false)
				if err != nil {
					t.Errorf("Failed to convert: %v", err)
					return
				}
				t.Logf("Converted to instrument: %s", instrument)
				
				// Verify format
				expectedFormat := tc.asset + "-26DEC25-"
				if !contains(instrument, expectedFormat) {
					t.Errorf("Instrument %s doesn't contain expected format %s", instrument, expectedFormat)
				}
			})
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr)))
}