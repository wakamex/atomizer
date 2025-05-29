package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

func main() {
	// Command line flags
	var (
		apiKey     = flag.String("api-key", os.Getenv("DERIBIT_API_KEY"), "Deribit API Key")
		apiSecret  = flag.String("api-secret", os.Getenv("DERIBIT_API_SECRET"), "Deribit API Secret")
		testMode   = flag.Bool("test", false, "Use testnet (default: false for mainnet)")
		instrument = flag.String("instrument", "ETH-28MAR25-5000-C", "Option instrument (e.g. ETH-28MAR25-5000-C)")
		quantity   = flag.Float64("qty", 0.1, "Quantity in ETH")
		multiplier = flag.Float64("mult", 2.0, "Ask price multiplier (default: 2x best ask)")
	)
	flag.Parse()

	if *apiKey == "" || *apiSecret == "" {
		log.Fatal("Error: DERIBIT_API_KEY and DERIBIT_API_SECRET must be set")
	}

	// Setup exchange
	apiURL := "https://www.deribit.com"
	if *testMode {
		apiURL = "https://test.deribit.com"
		log.Println("Using TESTNET")
	} else {
		log.Println("Using MAINNET - REAL MONEY")
	}

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

	// Construct symbol
	symbol := fmt.Sprintf("ETH/USD:%s", *instrument)
	log.Printf("Testing with instrument: %s", symbol)

	// Fetch order book
	log.Println("Fetching order book...")
	orderBook, err := exchange.FetchOrderBook(symbol)
	if err != nil {
		log.Fatalf("Failed to fetch order book: %v", err)
	}

	// Get best ask
	if len(orderBook.Asks) == 0 {
		log.Fatal("No asks in order book")
	}
	bestAsk := orderBook.Asks[0][0]
	ourPrice := bestAsk * *multiplier

	log.Printf("Best ask: %.4f ETH", bestAsk)
	log.Printf("Our ask: %.4f ETH (%.1fx best ask)", ourPrice, *multiplier)
	log.Printf("Quantity: %.2f ETH", *quantity)

	// Confirm before placing order
	fmt.Printf("\nReady to place SELL order? (y/n): ")
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" {
		log.Println("Order cancelled by user")
		return
	}

	// Place the order
	log.Println("Placing order...")
	startTime := time.Now()
	
	order, err := exchange.CreateOrder(
		symbol,
		"limit",
		"sell",
		*quantity,
		ccxt.WithCreateOrderPrice(ourPrice),
		ccxt.WithCreateOrderParams(map[string]interface{}{
			"advanced": "usd",
		}),
	)
	
	elapsed := time.Since(startTime)
	
	if err != nil {
		log.Fatalf("Failed to place order: %v", err)
	}

	log.Printf("✅ Order placed successfully in %v", elapsed)
	log.Printf("Order ID: %s", order.Id)
	log.Printf("Status: %s", order.Status)
	log.Printf("Price: %.4f", order.Price)
	log.Printf("Amount: %.2f", order.Amount)
	
	// Wait a moment then check order status
	time.Sleep(2 * time.Second)
	
	log.Println("Checking order status...")
	updatedOrder, err := exchange.FetchOrder(order.Id, symbol)
	if err != nil {
		log.Printf("Warning: Could not fetch order status: %v", err)
	} else {
		log.Printf("Current status: %s", updatedOrder.Status)
		log.Printf("Filled: %.2f / %.2f", updatedOrder.Filled, updatedOrder.Amount)
	}

	// Cancel the order
	fmt.Printf("\nCancel the order? (y/n): ")
	fmt.Scanln(&confirm)
	if confirm == "y" || confirm == "Y" {
		log.Println("Cancelling order...")
		_, err := exchange.CancelOrder(order.Id, symbol)
		if err != nil {
			log.Printf("Failed to cancel: %v", err)
		} else {
			log.Println("✅ Order cancelled")
		}
	}
}