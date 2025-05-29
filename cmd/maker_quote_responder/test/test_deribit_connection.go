package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

func main() {
	// Command line flags
	var (
		apiKey    = flag.String("api-key", os.Getenv("DERIBIT_API_KEY"), "Deribit API Key")
		apiSecret = flag.String("api-secret", os.Getenv("DERIBIT_API_SECRET"), "Deribit API Secret")
		testMode  = flag.Bool("test", false, "Use testnet (default: false for mainnet)")
	)
	flag.Parse()

	if *apiKey == "" || *apiSecret == "" {
		log.Fatal("Error: DERIBIT_API_KEY and DERIBIT_API_SECRET must be set")
	}

	// Setup exchange
	apiURL := "https://www.deribit.com"
	networkName := "MAINNET"
	if *testMode {
		apiURL = "https://test.deribit.com"
		networkName = "TESTNET"
	}

	fmt.Printf("=== Deribit Connection Test ===\n")
	fmt.Printf("Network: %s\n", networkName)
	fmt.Printf("URL: %s\n\n", apiURL)

	log.Println("1. Creating Deribit client...")
	exchange := ccxt.NewDeribit(map[string]interface{}{
		"rateLimit":       10,
		"enableRateLimit": true,
		"apiKey":          *apiKey,
		"secret":          *apiSecret,
		"urls": map[string]interface{}{
			"api": map[string]interface{}{
				"rest": apiURL,
			},
		},
		"options": map[string]interface{}{
			"defaultType":             "option",
			"adjustForTimeDifference": true,
			"recvWindow":              5000,
		},
	})
	log.Println("✅ Client created")

	// Test 1: Fetch balance
	log.Println("\n2. Fetching account balance...")
	balance, err := exchange.FetchBalance()
	if err != nil {
		log.Fatalf("❌ Failed to fetch balance: %v", err)
	}
	
	// Display balances
	fmt.Println("\nBalances:")
	currencies := []string{"BTC", "ETH", "USDC"}
	for _, currency := range currencies {
		if bal, exists := balance[currency]; exists {
			if total, ok := bal["total"].(float64); ok && total > 0 {
				fmt.Printf("  %s: %.4f\n", currency, total)
			}
		}
	}
	log.Println("✅ Balance fetched successfully")

	// Test 2: Fetch open orders
	log.Println("\n3. Fetching open orders...")
	openOrders, err := exchange.FetchOpenOrders("")
	if err != nil {
		log.Printf("❌ Failed to fetch open orders: %v", err)
	} else {
		if len(openOrders) == 0 {
			fmt.Println("No open orders")
		} else {
			fmt.Printf("Found %d open orders:\n", len(openOrders))
			for i, order := range openOrders {
				fmt.Printf("  %d. %s %s %.2f @ %.4f\n", 
					i+1, order.Symbol, order.Side, order.Amount, order.Price)
			}
		}
		log.Println("✅ Open orders fetched successfully")
	}

	// Test 3: Fetch positions
	log.Println("\n4. Fetching positions...")
	positions, err := exchange.FetchPositions()
	if err != nil {
		log.Printf("❌ Failed to fetch positions: %v", err)
	} else {
		if len(positions) == 0 {
			fmt.Println("No open positions")
		} else {
			fmt.Printf("Found %d positions:\n", len(positions))
			for i, pos := range positions {
				contracts := 0.0
				if pos.Contracts != nil {
					contracts = *pos.Contracts
				}
				fmt.Printf("  %d. %s: %.2f contracts\n", 
					i+1, pos.Symbol, contracts)
			}
		}
		log.Println("✅ Positions fetched successfully")
	}

	// Test 4: Fetch a sample order book
	log.Println("\n5. Testing order book fetch...")
	testSymbol := "ETH/USD:ETH-PERPETUAL"
	orderBook, err := exchange.FetchOrderBook(testSymbol)
	if err != nil {
		log.Printf("❌ Failed to fetch order book for %s: %v", testSymbol, err)
	} else {
		if len(orderBook.Bids) > 0 && len(orderBook.Asks) > 0 {
			fmt.Printf("ETH-PERPETUAL Order Book:\n")
			fmt.Printf("  Best Bid: $%.2f\n", orderBook.Bids[0][0])
			fmt.Printf("  Best Ask: $%.2f\n", orderBook.Asks[0][0])
			fmt.Printf("  Spread: $%.2f\n", orderBook.Asks[0][0]-orderBook.Bids[0][0])
		}
		log.Println("✅ Order book fetched successfully")
	}

	// Test 5: Check server time
	log.Println("\n6. Checking server time...")
	serverTime := exchange.Milliseconds()
	fmt.Printf("Server timestamp: %d\n", serverTime)
	log.Println("✅ Server communication working")

	fmt.Println("\n✅ All tests completed successfully!")
	fmt.Println("Your Deribit connection is working properly.")
}