package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	
	// Import your internal packages
	"github.com/wakamex/atomizer/internal/marketmaker"
	"github.com/wakamex/atomizer/internal/exchange"
	"github.com/wakamex/atomizer/internal/types"
)

func main() {
	// Define subcommands
	if len(os.Args) < 2 {
		fmt.Println("Usage: atomizer <command> [options]")
		fmt.Println("Commands:")
		fmt.Println("  market-maker    Run the market maker")
		fmt.Println("  rfq-responder   Run the RFQ responder")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "market-maker":
		runMarketMaker(os.Args[2:])
	case "rfq-responder":
		runRFQResponder(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runMarketMaker(args []string) {
	// Parse flags
	fs := flag.NewFlagSet("market-maker", flag.ExitOnError)
	
	// Exchange selection
	exchangeName := fs.String("exchange", "derive", "Exchange to use (derive, deribit)")
	
	// Market maker parameters
	underlying := fs.String("underlying", "ETH", "Underlying asset")
	expiry := fs.String("expiry", "", "Expiry date (YYYYMMDD)")
	strikes := fs.String("strikes", "", "Comma-separated list of strikes")
	allStrikes := fs.Bool("all-strikes", false, "Trade all available strikes")
	
	// Trading parameters
	spread := fs.Int("spread", 50, "Spread in basis points")
	minSpread := fs.Int("min-spread", 10, "Minimum spread in basis points")
	size := fs.Float64("size", 0.1, "Quote size")
	improvement := fs.Float64("improvement", 0.1, "Price improvement")
	improvementRefSize := fs.Float64("improvement-reference-size", 0, "Reference size for improvement")
	
	// Risk parameters
	maxPosition := fs.Float64("max-position", 10.0, "Maximum position per instrument")
	maxExposure := fs.Float64("max-exposure", 100.0, "Maximum total exposure")
	
	// Operational parameters
	refresh := fs.Int("refresh", 5, "Refresh interval in seconds")
	dryRun := fs.Bool("dry-run", false, "Dry run mode (no real orders)")
	test := fs.Bool("test", false, "Use test environment")
	bidOnly := fs.Bool("bid-only", false, "Only place bid orders (buy side)")
	askOnly := fs.Bool("ask-only", false, "Only place ask orders (sell side)")
	
	fs.Parse(args)
	
	// Validate required parameters
	if *expiry == "" {
		log.Fatal("Expiry date is required")
	}
	
	// Validate one-sided flags
	if *bidOnly && *askOnly {
		log.Fatal("Cannot use both --bid-only and --ask-only")
	}
	
	// Build instrument list
	instruments := buildInstrumentList(*underlying, *expiry, *strikes, *allStrikes)
	
	// Create market maker config
	config := &types.MarketMakerConfig{
		Exchange:         *exchangeName,
		ExchangeTestMode: *test,
		Instruments:      instruments,
		SpreadBps:        *spread,
		MinSpreadBps:     *minSpread,
		QuoteSize:        decimal.NewFromFloat(*size),
		RefreshInterval:  time.Duration(*refresh) * time.Second,
		MaxPositionSize:  decimal.NewFromFloat(*maxPosition),
		MaxTotalExposure: decimal.NewFromFloat(*maxExposure),
		Improvement:      decimal.NewFromFloat(*improvement),
		ImprovementReferenceSize: decimal.NewFromFloat(*improvementRefSize),
		CancelThreshold:  decimal.NewFromFloat(0.005), // 0.5% default
		MaxOrdersPerSide: 1,
		BidOnly:          *bidOnly,
		AskOnly:          *askOnly,
		// DryRun: *dryRun, // TODO: Add dry run support
	}
	
	// Show dry run warning if enabled
	if *dryRun {
		log.Println("WARNING: Running in DRY RUN mode - no real orders will be placed")
	}
	
	// Create exchange
	log.Printf("Creating %s exchange (test mode: %v)...", *exchangeName, *test)
	exchangeImpl, err := exchange.NewExchange(config)
	if err != nil {
		log.Fatalf("Failed to create exchange: %v", err)
	}
	
	// Create market maker
	mm := marketmaker.NewMarketMaker(config, exchangeImpl)
	
	// Start market maker
	log.Printf("Starting market maker with %d instruments...", len(instruments))
	log.Printf("Spread: %d bps, Size: %.2f, Refresh: %ds", *spread, *size, *refresh)
	
	if err := mm.Start(); err != nil {
		log.Fatalf("Failed to start market maker: %v", err)
	}
	
	// Wait for interrupt
	log.Println("Market maker running. Press Ctrl+C to stop...")
	select {} // Block forever until killed
}

func buildInstrumentList(underlying, expiry, strikes string, allStrikes bool) []string {
	var instruments []string
	
	if strikes != "" {
		// Parse comma-separated strikes
		strikeList := strings.Split(strikes, ",")
		for _, strike := range strikeList {
			strike = strings.TrimSpace(strike)
			if strike != "" {
				// Add both call and put for each strike
				instruments = append(instruments, 
					fmt.Sprintf("%s-%s-%s-C", underlying, expiry, strike),
					fmt.Sprintf("%s-%s-%s-P", underlying, expiry, strike),
				)
			}
		}
		return instruments
	}
	
	if allStrikes {
		// TODO: Fetch all available strikes from exchange
		log.Println("All strikes mode not yet implemented")
		return []string{}
	}
	
	// No strikes specified
	log.Println("Warning: No strikes specified. Use --strikes or --all-strikes")
	return []string{}
}

func runRFQResponder(args []string) {
	// TODO: Implement RFQ responder command
	log.Println("RFQ responder not yet implemented in new structure")
}