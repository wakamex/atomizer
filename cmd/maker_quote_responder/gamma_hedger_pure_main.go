package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/shopspring/decimal"
)

// PureGammaHedgerCommand runs the pure gamma hedger
func PureGammaHedgerCommand() {
	fs := flag.NewFlagSet("pure-gamma-hedger", flag.ExitOnError)
	
	// Configuration flags
	privateKey := fs.String("private-key", os.Getenv("DERIVE_PRIVATE_KEY"), "Private key")
	wallet := fs.String("wallet", os.Getenv("DERIVE_WALLET"), "Wallet address")
	subaccountID := fs.String("subaccount", os.Getenv("DERIVE_SUBACCOUNT_ID"), "Subaccount ID")
	deltaThreshold := fs.Float64("delta-threshold", 0.1, "Delta threshold before hedging (ETH)")
	minHedgeSize := fs.Float64("min-hedge-size", 0.1, "Minimum hedge size (ETH)")
	hedgeIntervalSec := fs.Int("hedge-interval", 30, "Hedge check interval (seconds)")
	aggressiveness := fs.Float64("aggressiveness", 0.7, "Order aggressiveness (0=passive, 1=aggressive)")
	debug := fs.Bool("debug", false, "Enable debug logging")
	
	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}
	
	// Validate inputs
	if *privateKey == "" || *wallet == "" {
		log.Fatal("Private key and wallet address are required")
	}
	
	if *aggressiveness < 0 || *aggressiveness > 1 {
		log.Fatal("Aggressiveness must be between 0 and 1")
	}
	
	// Parse subaccount ID
	var subID int
	if *subaccountID != "" {
		parsed, err := strconv.Atoi(*subaccountID)
		if err != nil {
			log.Fatalf("Invalid subaccount ID: %v", err)
		}
		subID = parsed
	}
	
	log.Println("========================================")
	log.Printf("Build hash: %s", getBuildHash())
	log.Println("========================================")
	log.Printf("Starting Pure Gamma Hedger")
	log.Printf("Configuration:")
	log.Printf("  Subaccount: %d", subID)
	log.Printf("  Delta Threshold: %.4f ETH", *deltaThreshold)
	log.Printf("  Min Hedge Size: %.4f ETH", *minHedgeSize)
	log.Printf("  Hedge Interval: %d seconds", *hedgeIntervalSec)
	log.Printf("  Aggressiveness: %.2f (%.0f%% through spread)", *aggressiveness, *aggressiveness*100)
	log.Printf("  Debug Mode: %v", *debug)
	
	// Set debug mode
	SetDebugMode(*debug)
	
	// Create derive market maker exchange
	exchange, err := NewDeriveMarketMakerExchange(*privateKey, *wallet)
	if err != nil {
		log.Fatalf("Failed to create exchange: %v", err)
	}
	
	// Subscribe to ETH-PERP orderbook
	log.Printf("Subscribing to ETH-PERP orderbook...")
	if wsClient := exchange.wsClient; wsClient != nil {
		if err := wsClient.SubscribeOrderBook("ETH-PERP", 10); err != nil {
			log.Printf("Warning: Failed to subscribe to orderbook: %v", err)
		}
		defer wsClient.Close()
	}
	
	// Create and configure pure gamma hedger
	hedger := NewPureGammaHedger(exchange, nil)
	hedger.deltaThreshold = decimal.NewFromFloat(*deltaThreshold)
	hedger.minHedgeSize = decimal.NewFromFloat(*minHedgeSize)
	hedger.hedgeInterval = time.Duration(*hedgeIntervalSec) * time.Second
	hedger.aggressiveness = decimal.NewFromFloat(*aggressiveness)
	hedger.debugMode = *debug
	
	// Start hedger
	if err := hedger.Start(); err != nil {
		log.Fatalf("Failed to start hedger: %v", err)
	}
	
	// Handle shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	select {
	case <-sigChan:
		log.Println("Shutting down...")
		hedger.Stop()
	case <-ctx.Done():
		log.Println("Context cancelled")
	}
}