package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
	
	// Import internal packages
	"github.com/wakamex/atomizer/internal/api"
	"github.com/wakamex/atomizer/internal/arbitrage"
	"github.com/wakamex/atomizer/internal/config"
	"github.com/wakamex/atomizer/internal/exchange"
	"github.com/wakamex/atomizer/internal/hedging"
	"github.com/wakamex/atomizer/internal/hedging/gamma"
	"github.com/wakamex/atomizer/internal/manual"
	"github.com/wakamex/atomizer/internal/marketmaker"
	"github.com/wakamex/atomizer/internal/rfq"
	"github.com/wakamex/atomizer/internal/risk"
	"github.com/wakamex/atomizer/internal/types"
	"github.com/wakamex/atomizer/internal/websocket"
)

func main() {
	// Define subcommands
	if len(os.Args) < 2 {
		fmt.Println("Usage: atomizer <command> [options]")
		fmt.Println("Commands:")
		fmt.Println("  market-maker    Run the market maker")
		fmt.Println("  rfq-responder   Run the RFQ responder")
		fmt.Println("  manual-order    Place a manual order")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "market-maker":
		runMarketMaker(os.Args[2:])
	case "rfq-responder":
		runRFQResponder(os.Args[2:])
	case "manual-order":
		runManualOrder(os.Args[2:])
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
	
	// Aggression parameter
	aggression := fs.Float64("aggression", 1.0, "Aggression level: 0=join best, 0.9=near mid, 1.0+=cross spread (default: 1.0)")
	
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
	
	// Validate aggression parameter
	if *aggression < 0 {
		log.Fatal("Aggression must be non-negative")
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
		Aggression:       decimal.NewFromFloat(*aggression),
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
	// Parse flags
	fs := flag.NewFlagSet("rfq-responder", flag.ExitOnError)
	
	// WebSocket configuration
	wsURL := fs.String("ws-url", "wss://rip-testnet.rysk.finance/maker", "WebSocket URL for RFQ stream")
	rfqAssets := fs.String("rfq-assets", "", "Comma-separated list of asset addresses for RFQ streams")
	
	// Exchange configuration
	exchangeName := fs.String("exchange", "derive", "Exchange to use (derive, deribit)")
	testMode := fs.Bool("test", false, "Use exchange testnet")
	
	// Trading configuration
	dummyPrice := fs.String("dummy-price", "1000000", "Fallback price for quotes")
	quoteDuration := fs.Int64("quote-duration", 30, "Quote validity duration in seconds")
	
	// Risk configuration
	maxDelta := fs.Float64("max-delta", 10.0, "Maximum position delta exposure")
	enableGamma := fs.Bool("enable-gamma", false, "Enable gamma hedging")
	gammaThreshold := fs.Float64("gamma-threshold", 0.1, "Gamma threshold for hedging")
	
	// API configuration
	httpPort := fs.Int("http-port", 8080, "Port for HTTP API server")
	enableManual := fs.Bool("enable-manual", true, "Enable manual trade API")
	
	fs.Parse(args)
	
	// Validate required parameters
	if *rfqAssets == "" {
		log.Fatal("RFQ asset addresses are required (--rfq-assets)")
	}
	
	// Load environment variables
	makerAddress := os.Getenv("MAKER_ADDRESS")
	privateKey := os.Getenv("PRIVATE_KEY")
	
	if makerAddress == "" || privateKey == "" {
		log.Fatal("MAKER_ADDRESS and PRIVATE_KEY environment variables are required")
	}
	
	// Create configuration
	cfg := &config.Config{
		WebSocketURL:              *wsURL,
		RFQAssetAddressesCSV:      *rfqAssets,
		MakerAddress:              makerAddress,
		PrivateKey:                privateKey,
		ExchangeName:              *exchangeName,
		ExchangeTestMode:          *testMode,
		DummyPrice:                *dummyPrice,
		QuoteValidDurationSeconds: *quoteDuration,
		MaxPositionDelta:          *maxDelta,
		EnableGammaHedging:        *enableGamma,
		GammaThreshold:            *gammaThreshold,
		HTTPPort:                  fmt.Sprintf("%d", *httpPort),
		EnableManualTrades:        *enableManual,
		AssetMapping:              config.DefaultAssetMapping, // Use default mappings
	}
	
	// Parse private key
	if err := parsePrivateKey(cfg); err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	
	// Create exchange
	exchange, err := createExchange(cfg)
	if err != nil {
		log.Fatalf("Failed to create exchange: %v", err)
	}
	
	// Create components
	riskManager := risk.NewManager(cfg)
	hedgeManager := hedging.NewManager(exchange, cfg)
	gammaModule := gamma.NewModule(cfg.GammaThreshold)
	gammaHedger := gamma.NewHedger(exchange, cfg, hedgeManager)
	
	// Create arbitrage orchestrator
	orchestrator := arbitrage.NewOrchestrator(
		cfg, exchange, hedgeManager, riskManager, gammaModule, gammaHedger,
	)
	
	// Start orchestrator
	orchestrator.Start()
	defer orchestrator.Stop()
	
	// Create and start HTTP server if enabled
	if cfg.EnableManualTrades {
		httpServer := api.NewServer(orchestrator, riskManager, *httpPort)
		go func() {
			log.Printf("Starting HTTP API server on port %d", *httpPort)
			if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
				log.Printf("HTTP server error: %v", err)
			}
		}()
		defer httpServer.Stop()
	}
	
	// Create RFQ processor
	rfqProcessor := rfq.NewProcessor(cfg, exchange)
	
	// Create WebSocket client
	wsClient := websocket.NewSimpleRFQClient(cfg, orchestrator, rfqProcessor)
	
	// Start WebSocket client
	if err := wsClient.Start(); err != nil {
		log.Fatalf("Failed to start WebSocket client: %v", err)
	}
	defer wsClient.Stop()
	
	log.Println("RFQ responder started. Press Ctrl+C to stop.")
	
	// Wait for interrupt signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	
	log.Println("Shutting down RFQ responder...")
}

// parsePrivateKey parses and validates the private key from the config
func parsePrivateKey(cfg *config.Config) error {
	// Validate format
	if err := validatePrivateKey(cfg.PrivateKey, cfg.MakerAddress); err != nil {
		return err
	}
	
	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}
	cfg.ParsedPrivateKey = privateKey
	return nil
}

// validatePrivateKey validates private key format and derives address
func validatePrivateKey(privateKeyHex, expectedAddress string) error {
	// Remove 0x prefix if present
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	
	// Check length
	if len(privateKeyHex) != 64 {
		return fmt.Errorf("private key must be 64 hex characters (got %d)", len(privateKeyHex))
	}
	
	// Check if valid hex
	if !regexp.MustCompile(`^[0-9a-fA-F]+$`).MatchString(privateKeyHex) {
		return fmt.Errorf("private key must contain only hexadecimal characters")
	}
	
	// Derive address
	derivedAddr, err := privateKeyToAddress(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to derive address: %w", err)
	}
	
	// Compare addresses
	if !strings.EqualFold(derivedAddr, expectedAddress) {
		return fmt.Errorf("private key doesn't match maker address\nDerived: %s\nExpected: %s", 
			derivedAddr, expectedAddress)
	}
	
	log.Printf("✓ Private key validation successful - derived address matches: %s", expectedAddress)
	return nil
}

// privateKeyToAddress derives the Ethereum address from a private key
func privateKeyToAddress(privateKeyHex string) (string, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("error casting public key to ECDSA")
	}
	
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address.Hex(), nil
}

// createExchange creates an exchange instance based on config
func createExchange(cfg *config.Config) (types.Exchange, error) {
	factory := exchange.NewFactory()
	
	// Map config to exchange config
	exchangeConfig := map[string]interface{}{
		"test_mode": cfg.ExchangeTestMode,
	}
	
	// Add Deribit credentials if needed
	if cfg.ExchangeName == "deribit" {
		exchangeConfig["api_key"] = cfg.DeribitApiKey
		exchangeConfig["api_secret"] = cfg.DeribitApiSecret
	}
	
	// Create the exchange
	ex, err := factory.CreateExchange(cfg.ExchangeName, exchangeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create exchange: %w", err)
	}
	
	// Ensure it implements the Exchange interface
	exchange, ok := ex.(types.Exchange)
	if !ok {
		return nil, fmt.Errorf("exchange %s does not implement Exchange interface", cfg.ExchangeName)
	}
	
	return exchange, nil
}

func runManualOrder(args []string) {
	// Parse flags
	fs := flag.NewFlagSet("manual-order", flag.ExitOnError)
	
	// Exchange configuration
	exchangeName := fs.String("exchange", "derive", "Exchange to use (derive, deribit)")
	testMode := fs.Bool("test", false, "Use test environment")
	
	// Order parameters
	instrument := fs.String("instrument", "", "Instrument to trade (required)")
	side := fs.String("side", "", "Order side: buy or sell (required)")
	price := fs.Float64("price", 0, "Order price (required)")
	amount := fs.Float64("amount", 0, "Order amount (required)")
	
	// Deribit specific
	deribitApiKey := fs.String("deribit-api-key", os.Getenv("DERIBIT_API_KEY"), "Deribit API key")
	deribitApiSecret := fs.String("deribit-api-secret", os.Getenv("DERIBIT_API_SECRET"), "Deribit API secret")
	
	// Derive specific
	derivePrivateKey := fs.String("derive-private-key", os.Getenv("DERIVE_PRIVATE_KEY"), "Derive private key")
	deriveWalletAddress := fs.String("derive-wallet-address", os.Getenv("DERIVE_WALLET_ADDRESS"), "Derive wallet address")
	
	fs.Parse(args)
	
	// Validate required parameters
	if *instrument == "" || *side == "" || *price == 0 || *amount == 0 {
		log.Fatal("All order parameters are required: --instrument, --side, --price, --amount")
	}
	
	// Create configuration
	cfg := &config.Config{
		ExchangeName:     *exchangeName,
		ExchangeTestMode: *testMode,
		DeribitApiKey:    *deribitApiKey,
		DeribitApiSecret: *deribitApiSecret,
		PrivateKey:       *derivePrivateKey,
		MakerAddress:     *deriveWalletAddress,
	}
	
	// Parse private key if using Derive
	if *exchangeName == "derive" {
		if *derivePrivateKey == "" || *deriveWalletAddress == "" {
			log.Fatal("Derive requires DERIVE_PRIVATE_KEY and DERIVE_WALLET_ADDRESS")
		}
		if err := parsePrivateKey(cfg); err != nil {
			log.Fatalf("Failed to parse private key: %v", err)
		}
	}
	
	// Create exchange
	exchange, err := createExchange(cfg)
	if err != nil {
		log.Fatalf("Failed to create exchange: %v", err)
	}
	
	// Ensure it's a market maker exchange
	mmExchange, ok := exchange.(types.MarketMakerExchange)
	if !ok {
		log.Fatalf("Exchange %s does not implement MarketMakerExchange interface", *exchangeName)
	}
	
	// Create order configuration
	orderCfg := manual.OrderConfig{
		Instrument: *instrument,
		Side:       *side,
		Price:      *price,
		Amount:     *amount,
	}
	
	// Override from environment if specified
	orderCfg = manual.ParseOrderFromEnv(orderCfg)
	
	// Run manual order
	if err := manual.RunManualOrder(cfg, mmExchange, orderCfg); err != nil {
		log.Fatalf("Failed to place order: %v", err)
	}
}