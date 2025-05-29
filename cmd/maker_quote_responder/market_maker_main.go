package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shopspring/decimal"
)

// RunMarketMaker runs the market maker with command line arguments
func RunMarketMaker(args []string) {
	// Create a new flag set for the market maker subcommand
	fs := flag.NewFlagSet("market-maker", flag.ExitOnError)
	
	// Command-line flags
	var (
		exchange      = fs.String("exchange", "derive", "Exchange to use (derive, deribit)")
		testMode      = fs.Bool("test", false, "Use exchange testnet")
		expiry        = fs.String("expiry", "", "Expiry date to make markets on (e.g., 20250606)")
		underlying    = fs.String("underlying", "ETH", "Underlying asset (ETH, BTC)")
		strikes       = fs.String("strikes", "", "Comma-separated list of strikes to trade")
		allStrikes    = fs.Bool("all-strikes", false, "Make markets on all available strikes")
		spreadBps     = fs.Int("spread", 10, "Spread in basis points (100 = 1%)")
		size          = fs.Float64("size", 0.1, "Quote size")
		refreshSec    = fs.Int("refresh", 1, "Quote refresh interval in seconds")
		maxPosition   = fs.Float64("max-position", 1.0, "Maximum position per instrument")
		maxExposure   = fs.Float64("max-exposure", 10.0, "Maximum total exposure")
		minSpreadBps  = fs.Int("min-spread", 1000, "Minimum spread in basis points")
		improvement   = fs.Float64("improvement", 0.1, "Price improvement amount")
		privateKey    = fs.String("private-key", "", "Private key (overrides env var)")
		walletAddress = fs.String("wallet", "", "Wallet address (Derive only)")
		dryRun        = fs.Bool("dry-run", false, "Print configuration without starting")
	)
	
	// Parse the arguments
	fs.Parse(args)
	
	// Validate required parameters
	if *expiry == "" {
		log.Fatal("Expiry date is required (use -expiry flag)")
	}
	
	if !*allStrikes && *strikes == "" {
		log.Fatal("Either -strikes or -all-strikes must be specified")
	}
	
	// Get credentials from env if not provided
	if *privateKey == "" {
		if *exchange == "derive" {
			*privateKey = os.Getenv("DERIVE_PRIVATE_KEY")
		} else if *exchange == "deribit" {
			*privateKey = os.Getenv("DERIBIT_PRIVATE_KEY")
		}
	}
	
	if *walletAddress == "" && *exchange == "derive" {
		*walletAddress = os.Getenv("DERIVE_WALLET_ADDRESS")
	}
	
	// Build instrument list
	instruments := buildMarketMakerInstrumentList(*underlying, *expiry, *strikes, *allStrikes)
	
	// Create configuration
	config := &MarketMakerConfig{
		Exchange:         *exchange,
		ExchangeTestMode: *testMode,
		Instruments:      instruments,
		SpreadBps:        *spreadBps,
		QuoteSize:        decimal.NewFromFloat(*size),
		RefreshInterval:  time.Duration(*refreshSec) * time.Second,
		PriceImprovement: decimal.NewFromFloat(*improvement),
		MaxPositionSize:  decimal.NewFromFloat(*maxPosition),
		MaxTotalExposure: decimal.NewFromFloat(*maxExposure),
		CancelThreshold:  decimal.NewFromFloat(0.005), // 0.5% default
		MaxOrdersPerSide: 1,
		MinSpreadBps:     *minSpreadBps,
		TargetFillRate:   decimal.NewFromFloat(0.1), // 10% default
	}
	
	// Print configuration
	log.Println("Market Maker Configuration:")
	log.Printf("  Exchange: %s (test mode: %v)", config.Exchange, config.ExchangeTestMode)
	log.Printf("  Underlying: %s", *underlying)
	log.Printf("  Expiry: %s", *expiry)
	log.Printf("  Instruments: %d", len(config.Instruments))
	if len(config.Instruments) <= 10 {
		for _, inst := range config.Instruments {
			log.Printf("    - %s", inst)
		}
	} else {
		log.Printf("    - %s ... %s", config.Instruments[0], config.Instruments[len(config.Instruments)-1])
	}
	log.Printf("  Spread: %d bps (%.2f%%)", config.SpreadBps, float64(config.SpreadBps)/100)
	log.Printf("  Quote Size: %s", config.QuoteSize)
	log.Printf("  Price Improvement: %s", config.PriceImprovement)
	log.Printf("  Refresh Interval: %s", config.RefreshInterval)
	log.Printf("  Max Position: %s", config.MaxPositionSize)
	log.Printf("  Max Exposure: %s", config.MaxTotalExposure)
	
	if *dryRun {
		log.Println("Dry run mode - exiting without starting market maker")
		return
	}
	
	// Create exchange adapter
	var exchangeAdapter MarketMakerExchange
	var err error
	
	switch config.Exchange {
	case "derive":
		if *privateKey == "" || *walletAddress == "" {
			log.Fatal("Derive requires private key and wallet address")
		}
		exchangeAdapter, err = NewDeriveMarketMakerExchange(*privateKey, *walletAddress)
		if err != nil {
			log.Fatalf("Failed to create Derive exchange: %v", err)
		}
	case "deribit":
		log.Fatal("Deribit market maker not yet implemented")
		// TODO: Implement Deribit adapter
	default:
		log.Fatalf("Unknown exchange: %s", config.Exchange)
	}
	
	// Create market maker
	marketMaker := NewMarketMaker(config, exchangeAdapter)
	
	// Start market maker
	if err := marketMaker.Start(); err != nil {
		log.Fatalf("Failed to start market maker: %v", err)
	}
	
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	log.Println("Market maker running. Press Ctrl+C to stop...")
	
	// Wait for shutdown signal
	<-sigChan
	
	log.Println("Shutting down market maker...")
	if err := marketMaker.Stop(); err != nil {
		log.Printf("Error stopping market maker: %v", err)
	}
	
	// Close exchange connections
	if closer, ok := exchangeAdapter.(interface{ Close() error }); ok {
		closer.Close()
	}
	
	log.Println("Market maker stopped")
}

// buildMarketMakerInstrumentList builds the list of instruments to trade
func buildMarketMakerInstrumentList(underlying, expiry string, strikes string, allStrikes bool) []string {
	var instruments []string
	
	if allStrikes {
		// TODO: Query exchange for all available strikes
		log.Println("All strikes mode - would query exchange for available strikes")
		// For now, use a default set
		defaultStrikes := []string{"2600", "2700", "2800", "2900", "3000", "3100", "3200"}
		for _, strike := range defaultStrikes {
			// Format: ETH-20250606-2800-C
			instruments = append(instruments, 
				fmt.Sprintf("%s-%s-%s-C", underlying, expiry, strike))
		}
	} else {
		// Parse comma-separated strikes
		strikeList := strings.Split(strikes, ",")
		for _, strike := range strikeList {
			strike = strings.TrimSpace(strike)
			if strike != "" {
				instruments = append(instruments, 
					fmt.Sprintf("%s-%s-%s-C", underlying, expiry, strike))
			}
		}
	}
	
	return instruments
}