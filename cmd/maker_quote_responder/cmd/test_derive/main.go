package main

import (
	"log"
	"os"
	"time"
	
	maker "github.com/wakamex/atomizer/cmd/maker_quote_responder"
)

func main() {
	// Set up test configuration
	config := maker.ExchangeConfig{
		APIKey:    os.Getenv("PRIVATE_KEY"),
		APISecret: "",
		TestMode:  true,
		RateLimit: 10,
	}

	// Create Derive exchange
	exchange, err := maker.NewCCXTDeriveExchange(config)
	if err != nil {
		log.Fatalf("Failed to create Derive exchange: %v", err)
	}

	// Create a test RFQ
	rfq := maker.RFQResult{
		Strike: "300000000000", // 3000 in 8 decimal places
		Expiry: time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days from now
		IsPut:  false,
	}

	asset := "ETH"

	log.Printf("Testing Derive with RFQ: asset=%s, strike=%s, expiry=%d", asset, rfq.Strike, rfq.Expiry)

	// Test GetOrderBook
	orderBook, err := exchange.GetOrderBook(rfq, asset)
	if err != nil {
		log.Printf("Failed to get order book: %v", err)
	} else {
		log.Printf("Order book retrieved successfully:")
		log.Printf("  Bids: %d", len(orderBook.Bids))
		log.Printf("  Asks: %d", len(orderBook.Asks))
		log.Printf("  Index price: %f", orderBook.Index)
		
		if len(orderBook.Asks) > 0 {
			log.Printf("  Best ask: %f", orderBook.Asks[0][0])
		}
	}

	// Test instrument conversion
	instrument, err := exchange.ConvertToInstrument(asset, rfq.Strike, rfq.Expiry, rfq.IsPut)
	if err != nil {
		log.Printf("Failed to convert instrument: %v", err)
	} else {
		log.Printf("Converted instrument: %s", instrument)
	}

	// Test a few different strike prices and expiries
	testCases := []struct {
		strike string
		expiry time.Time
		desc   string
	}{
		{"250000000000", time.Now().Add(7 * 24 * time.Hour), "7 days, 2500 strike"},
		{"350000000000", time.Now().Add(30 * 24 * time.Hour), "30 days, 3500 strike"},
		{"400000000000", time.Now().Add(60 * 24 * time.Hour), "60 days, 4000 strike"},
	}

	for _, tc := range testCases {
		rfqTest := maker.RFQResult{
			Strike: tc.strike,
			Expiry: tc.expiry.Unix(),
			IsPut:  false,
		}
		
		instrument, err := exchange.ConvertToInstrument(asset, rfqTest.Strike, rfqTest.Expiry, rfqTest.IsPut)
		if err != nil {
			log.Printf("Test case %s failed: %v", tc.desc, err)
		} else {
			log.Printf("Test case %s: %s", tc.desc, instrument)
		}
	}
}