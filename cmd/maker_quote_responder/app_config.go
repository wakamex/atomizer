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
}

// LoadConfig parses command-line flags and environment variables
// to populate and return an AppConfig struct.
func LoadConfig() *AppConfig {
	cfg := &AppConfig{}

	flag.StringVar(&cfg.WebSocketURL, "websocket_url", "wss://rip-testnet.rysk.finance/maker", "WebSocket URL for RFQ stream and quote submission")
	flag.StringVar(&cfg.RFQAssetAddressesCSV, "rfq_asset_addresses", "", "Comma-separated list of asset addresses for RFQ streams (e.g., 0xAsset1,0xAsset2)")
	flag.StringVar(&cfg.DummyPrice, "dummy_price", "1000000", "Dummy price to quote (ensure format matches Rysk requirements, e.g., units)")
	flag.Int64Var(&cfg.QuoteValidDurationSeconds, "quote_valid_duration_seconds", 30, "How long your quotes will be valid in seconds")
	flag.StringVar(&cfg.ExchangeName, "exchange", "derive", "Exchange to use for hedging (e.g., derive, deribit, okx, bybit)")
	flag.BoolVar(&cfg.ExchangeTestMode, "exchange_test_mode", false, "Use exchange testnet (true) or mainnet (false)")
	flag.Parse()

	cfg.MakerAddress = os.Getenv("MAKER_ADDRESS")
	cfg.PrivateKey = os.Getenv("PRIVATE_KEY")
	cfg.DeribitApiKey = os.Getenv("DERIBIT_API_KEY")
	cfg.DeribitApiSecret = os.Getenv("DERIBIT_API_SECRET")
	
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
				log.Fatal("Error: DERIBIT_API_KEY environment variable is not set or empty (or provide ASYMMETRIC_PRIVATE_KEY and DERIBIT_CLIENT_ID).")
			}
			if cfg.DeribitApiSecret == "" {
				log.Fatal("Error: DERIBIT_API_SECRET environment variable is not set or empty (or provide ASYMMETRIC_PRIVATE_KEY and DERIBIT_CLIENT_ID).")
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
