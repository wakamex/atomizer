#!/bin/bash

# Script to place a test order on Derive exchange

cd "$(dirname "$0")"

# Load environment variables
export $(grep -v '^#' .env | grep -v '^$' | xargs)

# Default values
DEFAULT_SYMBOL="ETH/USDC:USDC-25-12-26-3400-C"
DEFAULT_QUANTITY="0.1"
DEFAULT_MULTIPLIER="2.0"

# Parse command line arguments
SYMBOL="${1:-$DEFAULT_SYMBOL}"
QUANTITY="${2:-$DEFAULT_QUANTITY}"
MULTIPLIER="${3:-$DEFAULT_MULTIPLIER}"

cat > /tmp/test_order.go << EOF
package main

import (
    "fmt"
    "log"
    "os"
    "strconv"
    ccxt "github.com/ccxt/ccxt/go/v4"
)

func main() {
    privateKey := os.Getenv("PRIVATE_KEY")
    if privateKey == "" {
        log.Fatal("PRIVATE_KEY not set")
    }

    symbol := "$SYMBOL"
    quantity, _ := strconv.ParseFloat("$QUANTITY", 64)
    multiplier, _ := strconv.ParseFloat("$MULTIPLIER", 64)

    log.Printf("Test order parameters:")
    log.Printf("  Symbol: %s", symbol)
    log.Printf("  Quantity: %f", quantity)
    log.Printf("  Price multiplier: %fx", multiplier)

    exchange := ccxt.NewDerive(map[string]interface{}{
        "rateLimit":       10,
        "enableRateLimit": true,
        "privateKey":      privateKey,
        "options": map[string]interface{}{
            "defaultType": "option",
        },
    })

    // Load markets
    log.Println("Loading markets...")
    marketsChan := exchange.LoadMarkets()
    marketsRaw := <-marketsChan
    
    if err, ok := marketsRaw.(error); ok {
        log.Fatalf("Error loading markets: %v", err)
    }

    // Fetch ticker to get current price
    log.Printf("Fetching ticker for %s...", symbol)
    ticker, err := exchange.FetchTicker(symbol)
    if err != nil {
        log.Fatalf("Error fetching ticker: %v", err)
    }

    if ticker.Ask == nil || *ticker.Ask <= 0 {
        log.Fatal("No ask price available")
    }

    bestAsk := *ticker.Ask
    orderPrice := bestAsk * multiplier

    log.Printf("Current best ask: %f", bestAsk)
    log.Printf("Order price (%fx): %f", multiplier, orderPrice)

    // Confirm before placing order
    fmt.Printf("\nReady to place SELL order:\n")
    fmt.Printf("  Symbol: %s\n", symbol)
    fmt.Printf("  Quantity: %f\n", quantity)
    fmt.Printf("  Price: %f USDC\n", orderPrice)
    fmt.Printf("\nPress Enter to continue or Ctrl+C to cancel...")
    fmt.Scanln()

    // Place the order
    log.Println("Placing order...")
    order, err := exchange.CreateOrder(
        symbol,
        "limit",
        "sell",
        quantity,
        ccxt.WithCreateOrderPrice(orderPrice),
    )

    if err != nil {
        log.Fatalf("Error placing order: %v", err)
    }

    fmt.Println("\n✅ Order placed successfully!")
    fmt.Printf("Order ID: %s\n", order.Id)
    fmt.Printf("Symbol: %s\n", order.Symbol)
    fmt.Printf("Side: %s\n", order.Side)
    fmt.Printf("Type: %s\n", order.Type)
    fmt.Printf("Quantity: %f\n", order.Amount)
    fmt.Printf("Price: %f\n", order.Price)
    fmt.Printf("Status: %s\n", order.Status)
}
EOF

echo "Test Order Placement Script"
echo "=========================="
echo "Symbol: $SYMBOL"
echo "Quantity: $QUANTITY"
echo "Price multiplier: ${MULTIPLIER}x"
echo ""
echo "Usage: $0 [symbol] [quantity] [multiplier]"
echo "Example: $0 'ETH/USDC:USDC-25-12-26-3400-C' 0.1 2.0"
echo ""

go run /tmp/test_order.go
rm /tmp/test_order.go