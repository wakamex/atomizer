#!/bin/bash

# Script to monitor open orders on Derive exchange

cd "$(dirname "$0")"

# Load environment variables
export $(grep -v '^#' .env | grep -v '^$' | xargs)

# Create a monitoring script
cat > /tmp/monitor_derive.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "os"
    "time"
    ccxt "github.com/ccxt/ccxt/go/v4"
)

func main() {
    privateKey := os.Getenv("PRIVATE_KEY")
    if privateKey == "" {
        log.Fatal("PRIVATE_KEY not set")
    }

    log.Println("Connecting to Derive...")
    
    exchange := ccxt.NewDerive(map[string]interface{}{
        "rateLimit":       10,
        "enableRateLimit": true,
        "privateKey":      privateKey,
        "options": map[string]interface{}{
            "defaultType": "option",
        },
    })

    // Load markets first
    log.Println("Loading markets...")
    marketsChan := exchange.LoadMarkets()
    marketsRaw := <-marketsChan
    
    if err, ok := marketsRaw.(error); ok {
        log.Fatalf("Error loading markets: %v", err)
    }

    for {
        fmt.Println("\n========================================")
        fmt.Printf("Checking orders at %s\n", time.Now().Format("15:04:05"))
        fmt.Println("========================================")

        // Fetch open orders
        openOrders, err := exchange.FetchOpenOrders()
        if err != nil {
            log.Printf("Error fetching open orders: %v", err)
        } else {
            if len(openOrders) == 0 {
                fmt.Println("No open orders")
            } else {
                fmt.Printf("Found %d open orders:\n", len(openOrders))
                for _, order := range openOrders {
                    fmt.Printf("\nOrder ID: %s\n", order.Id)
                    fmt.Printf("  Symbol: %s\n", order.Symbol)
                    fmt.Printf("  Type: %s\n", order.Type)
                    fmt.Printf("  Side: %s\n", order.Side)
                    fmt.Printf("  Price: %f\n", order.Price)
                    fmt.Printf("  Amount: %f\n", order.Amount)
                    fmt.Printf("  Status: %s\n", order.Status)
                    if order.Datetime != nil {
                        fmt.Printf("  Created: %s\n", *order.Datetime)
                    }
                }
            }
        }

        // Also check recent trades
        fmt.Println("\n--- Recent Trades ---")
        trades, err := exchange.FetchMyTrades()
        if err != nil {
            log.Printf("Error fetching trades: %v", err)
        } else {
            if len(trades) == 0 {
                fmt.Println("No recent trades")
            } else {
                // Show last 5 trades
                start := 0
                if len(trades) > 5 {
                    start = len(trades) - 5
                }
                for i := start; i < len(trades); i++ {
                    trade := trades[i]
                    fmt.Printf("\nTrade ID: %s\n", trade.Id)
                    fmt.Printf("  Symbol: %s\n", trade.Symbol)
                    fmt.Printf("  Side: %s\n", trade.Side)
                    fmt.Printf("  Price: %f\n", trade.Price)
                    fmt.Printf("  Amount: %f\n", trade.Amount)
                    if trade.Datetime != nil {
                        fmt.Printf("  Time: %s\n", *trade.Datetime)
                    }
                }
            }
        }

        // Wait 30 seconds before next check
        fmt.Println("\nWaiting 30 seconds before next check...")
        time.Sleep(30 * time.Second)
    }
}
EOF

# Run the monitor
echo "Starting Derive order monitor..."
echo "This will check for open orders every 30 seconds"
echo "Press Ctrl+C to stop"
echo ""
go run /tmp/monitor_derive.go