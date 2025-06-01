package main

import (
	"crypto/ecdsa"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

// AppConfig holds all configuration for the application.
type AppConfig struct {
	WebSocketURL              string
	RFQAssetAddressesCSV      string
	MakerAddress              string
	PrivateKey                string
	DeribitApiKey             string
	DeribitApiSecret          string
	DummyPrice                string
	QuoteValidDurationSeconds int64
	AssetMapping              map[string]string // Maps asset addresses to underlying symbols (ETH, BTC, SOL)
	ExchangeName              string            // Name of the exchange to use (e.g., "deribit", "okx", "bybit")
	ExchangeTestMode          bool              // Whether to use the exchange's testnet
	
	// Arbitrage configuration
	HTTPPort                  string            // Port for HTTP API server
	MaxPositionDelta          float64           // Maximum position delta exposure
	MinLiquidityScore         float64           // Minimum liquidity score for trades
	EnableGammaHedging        bool              // Enable gamma hedging
	GammaThreshold            float64           // Gamma threshold for hedging
	EnableManualTrades        bool              // Enable manual trade API
	CacheBackend              string            // Cache backend: "file" or "valkey"
	ValkeyAddr                string            // Valkey server address
}

// LoadConfig parses command-line flags and environment variables
// to populate and return an AppConfig struct.
func LoadConfig() *AppConfig {
	cfg := &AppConfig{}

	// Check if flags have already been parsed (e.g., by a subcommand)
	if flag.Parsed() {
		// If already parsed, just load from environment
		cfg.WebSocketURL = "wss://rip-testnet.rysk.finance/maker"
		cfg.DummyPrice = "1000000"
		cfg.QuoteValidDurationSeconds = 30
		cfg.ExchangeName = "derive"
		cfg.HTTPPort = "8080"
		cfg.MaxPositionDelta = 10.0
		cfg.MinLiquidityScore = 0.001
		cfg.EnableManualTrades = true
		cfg.GammaThreshold = 0.1
		cfg.CacheBackend = "file"
		cfg.ValkeyAddr = "localhost:6379"
	} else {
		flag.StringVar(&cfg.WebSocketURL, "websocket_url", "wss://rip-testnet.rysk.finance/maker", "WebSocket URL for RFQ stream and quote submission")
		flag.StringVar(&cfg.RFQAssetAddressesCSV, "rfq_asset_addresses", "", "Comma-separated list of asset addresses for RFQ streams (e.g., 0xAsset1,0xAsset2)")
		flag.StringVar(&cfg.DummyPrice, "dummy_price", "1000000", "Dummy price to quote (ensure format matches Rysk requirements, e.g., units)")
		flag.Int64Var(&cfg.QuoteValidDurationSeconds, "quote_valid_duration_seconds", 30, "How long your quotes will be valid in seconds")
		flag.StringVar(&cfg.ExchangeName, "exchange", "derive", "Exchange to use for hedging (e.g., derive, deribit, okx, bybit)")
		flag.BoolVar(&cfg.ExchangeTestMode, "exchange_test_mode", false, "Use exchange testnet (true) or mainnet (false)")
		
		// Arbitrage flags
		flag.StringVar(&cfg.HTTPPort, "http_port", "8080", "Port for HTTP API server")
		flag.Float64Var(&cfg.MaxPositionDelta, "max_position_delta", 10.0, "Maximum position delta exposure")
		flag.Float64Var(&cfg.MinLiquidityScore, "min_liquidity_score", 0.001, "Minimum liquidity score for trades")
		flag.BoolVar(&cfg.EnableManualTrades, "enable_manual_trades", true, "Enable manual trade API")
		flag.BoolVar(&cfg.EnableGammaHedging, "enable_gamma_hedging", false, "Enable gamma hedging")
		flag.Float64Var(&cfg.GammaThreshold, "gamma_threshold", 0.1, "Gamma threshold for hedging")
		flag.StringVar(&cfg.CacheBackend, "cache_backend", "file", "Cache backend: file or valkey")
		flag.StringVar(&cfg.ValkeyAddr, "valkey_addr", "localhost:6379", "Valkey server address")
		
		flag.Parse()
	}

	cfg.MakerAddress = os.Getenv("MAKER_ADDRESS")
	cfg.PrivateKey = os.Getenv("PRIVATE_KEY")
	cfg.DeribitApiKey = os.Getenv("DERIBIT_API_KEY")
	cfg.DeribitApiSecret = os.Getenv("DERIBIT_API_SECRET")
	
	// Skip validation if running as a subcommand (flags already parsed)
	if flag.Parsed() && cfg.RFQAssetAddressesCSV == "" {
		// Running as subcommand, skip RFQ validation
		return cfg
	}
	
	if cfg.RFQAssetAddressesCSV == "" {
		log.Fatal("Error: --rfq_asset_addresses is required.")
	}
	if cfg.MakerAddress == "" {
		log.Fatal("Error: MAKER_ADDRESS environment variable is not set or empty.")
	}
	if cfg.PrivateKey == "" {
		log.Fatal("Error: PRIVATE_KEY environment variable is not set or empty.")
	}
	
	// Only check for Deribit credentials if Deribit is the selected exchange
	if strings.ToLower(cfg.ExchangeName) == "deribit" {
		// Check for asymmetric key first
		asymmetricPrivateKey := os.Getenv("ASYMMETRIC_PRIVATE_KEY")
		deribitClientId := os.Getenv("DERIBIT_CLIENT_ID")
		
		// Only require standard API credentials if asymmetric auth is not available
		if asymmetricPrivateKey == "" || deribitClientId == "" {
			if cfg.DeribitApiKey == "" {
				log.Printf("Warning: DERIBIT_API_KEY not set. Deribit exchange will not be available.")
				log.Printf("To use Deribit, set DERIBIT_API_KEY and DERIBIT_API_SECRET, or ASYMMETRIC_PRIVATE_KEY and DERIBIT_CLIENT_ID.")
				// Change to a different exchange if Deribit credentials are not available
				cfg.ExchangeName = "derive"
				log.Printf("Switching to Derive exchange instead.")
			} else if cfg.DeribitApiSecret == "" {
				log.Fatal("Error: DERIBIT_API_SECRET environment variable is not set or empty.")
			}
		}
	}
	
	// Validate private key format
	if len(cfg.PrivateKey) != 64 {
		log.Fatalf("Error: PRIVATE_KEY must be exactly 64 characters long (got %d). Example: 72d4422755956df7a8e225603c24122c97b9650e245af67a40f100f955272064", len(cfg.PrivateKey))
	}
	
	// Check if private key contains only valid hex characters
	validHex := regexp.MustCompile(`^[0-9a-fA-F]+$`)
	if !validHex.MatchString(cfg.PrivateKey) {
		log.Fatalf("Error: PRIVATE_KEY must contain only hexadecimal characters (0-9, a-f). Current value '%s' contains invalid characters.", cfg.PrivateKey)
	}
	
	// Verify that private key matches maker address
	derivedAddress, err := privateKeyToAddress(cfg.PrivateKey)
	if err != nil {
		log.Fatalf("Error: Failed to derive address from PRIVATE_KEY: %v", err)
	}
	
	if !strings.EqualFold(derivedAddress, cfg.MakerAddress) {
		log.Fatalf("Error: PRIVATE_KEY does not match MAKER_ADDRESS.\nDerived address: %s\nMaker address: %s\nPlease ensure the private key corresponds to the maker address.", derivedAddress, cfg.MakerAddress)
	}
	
	log.Printf("✓ Private key validation successful - derived address matches maker address: %s", cfg.MakerAddress)

	// Load arbitrage configuration from environment
	// These can override command-line flags if needed

	// Initialize asset mapping
	// TODO: This should be configurable via environment variables or config file
	cfg.AssetMapping = map[string]string{
		"0xb67bfa7b488df4f2efa874f4e59242e9130ae61f": "ETH", // Example mapping for testnet
		// Add more mappings as needed
	}

	return cfg
}

// privateKeyToAddress derives the Ethereum address from a private key hex string
func privateKeyToAddress(privateKeyHex string) (string, error) {
	// Convert hex string to ECDSA private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	
	// Get the public key
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", err
	}
	
	// Derive the address
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address.Hex(), nil
}
