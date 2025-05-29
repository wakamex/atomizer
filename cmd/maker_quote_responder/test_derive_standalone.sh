#!/bin/bash

# Test script to verify Derive exchange functionality

cd "$(dirname "$0")"

# Load environment variables
export $(grep -v '^#' .env | grep -v '^$' | xargs)

# Create a minimal Go program to test Derive
cat > /tmp/test_derive.go << 'EOF'
package main

import (
    "log"
    "os"
    "strings"
    ccxt "github.com/ccxt/ccxt/go/v4"
)

func getMapKeys(m map[string]interface{}) []string {
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}

func main() {
    privateKey := os.Getenv("PRIVATE_KEY")
    if privateKey == "" {
        log.Fatal("PRIVATE_KEY not set")
    }

    log.Println("Creating Derive exchange...")
    
    exchange := ccxt.NewDerive(map[string]interface{}{
        "rateLimit":       10,
        "enableRateLimit": true,
        "privateKey":      privateKey,
        "options": map[string]interface{}{
            "defaultType": "option",
        },
    })

    // Test 1: Load markets
    log.Println("Loading markets...")
    marketsChan := exchange.LoadMarkets()
    marketsRaw := <-marketsChan
    
    if err, ok := marketsRaw.(error); ok {
        log.Printf("Error loading markets: %v", err)
        return
    }
    
    log.Printf("Markets loaded, type: %T", marketsRaw)
    
    // Test 2: Try to fetch an order book for a test symbol
    symbols := []string{
        "ETH-20250131-3000-C",
        "ETH-PERP",
        "ETH/USD",
        "ETH:ETH-20250131-3000-C",
    }
    
    for _, symbol := range symbols {
        log.Printf("\nTrying symbol: %s", symbol)
        orderBookChan := exchange.FetchOrderBook(symbol)
        orderBookRaw := <-orderBookChan
        
        if err, ok := orderBookRaw.(error); ok {
            log.Printf("  Error: %v", err)
        } else {
            log.Printf("  Success! Type: %T", orderBookRaw)
            log.Printf("  Value: %v", orderBookRaw)
            
            // Try different type assertions
            switch v := orderBookRaw.(type) {
            case ccxt.OrderBook:
                log.Printf("  OrderBook: Bids=%d, Asks=%d", len(v.Bids), len(v.Asks))
            case *ccxt.OrderBook:
                log.Printf("  *OrderBook: Bids=%d, Asks=%d", len(v.Bids), len(v.Asks))
            case map[string]interface{}:
                log.Printf("  Map with keys: %v", getMapKeys(v))
            case string:
                log.Printf("  String value: %s", v)
            default:
                log.Printf("  Unknown type")
            }
        }
    }
    
    // Test 3: Try FetchTicker instead
    log.Printf("\n\nTrying FetchTicker...")
    tickerSymbols := []string{
        "ETH/USDC:USDC-25-12-26-3400-C",
        "ETH/USDC",
    }
    
    for _, symbol := range tickerSymbols {
        log.Printf("\nFetchTicker for: %s", symbol)
        ticker, err := exchange.FetchTicker(symbol)
        if err != nil {
            log.Printf("  Error: %v", err)
        } else {
            log.Printf("  Success!")
            if ticker.Bid != nil {
                log.Printf("  Bid: %f", *ticker.Bid)
            }
            if ticker.Ask != nil {
                log.Printf("  Ask: %f", *ticker.Ask)
            }
            if ticker.Last != nil {
                log.Printf("  Last: %f", *ticker.Last)
            }
        }
    }
    
    // Test 4: Check the Markets field
    log.Printf("\nChecking exchange.Markets field...")
    if exchange.Markets != nil {
        log.Printf("Markets map has %d entries", len(exchange.Markets))
        
        // Show option markets specifically
        optionCount := 0
        var ethOptions []string
        for symbol, marketRaw := range exchange.Markets {
            if marketMap, ok := marketRaw.(map[string]interface{}); ok {
                if typeStr, exists := marketMap["type"]; exists && typeStr == "option" {
                    optionCount++
                    if base, exists := marketMap["base"]; exists {
                        if base == "ETH" {
                            ethOptions = append(ethOptions, symbol)
                        }
                    }
                }
            }
        }
        log.Printf("Total option markets: %d", optionCount)
        
        // Show all ETH options
        log.Printf("\nAll ETH options (%d):", len(ethOptions))
        
        // Look for June 2025 options specifically
        june2025Options := []string{}
        for _, opt := range ethOptions {
            if strings.Contains(opt, "25-06") {
                june2025Options = append(june2025Options, opt)
            }
        }
        
        if len(june2025Options) > 0 {
            log.Printf("\nJune 2025 ETH options (%d):", len(june2025Options))
            for _, opt := range june2025Options {
                log.Printf("  %s", opt)
            }
        } else {
            log.Printf("\nNo June 2025 options found")
        }
        
        // Show sample of all options
        log.Printf("\nSample of all ETH options:")
        for i, opt := range ethOptions {
            if i < 20 { // Show first 20
                log.Printf("  %s", opt)
            }
        }
    } else {
        log.Printf("Markets field is nil")
    }
}
EOF

# Run the test
echo "Running Derive connection test..."
go run /tmp/test_derive.go

# Clean up
rm /tmp/test_derive.go