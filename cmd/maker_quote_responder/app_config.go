package main

import (
	"flag"
	"log"
	"os"
)

// AppConfig holds all configuration for the application.
type AppConfig struct {
	WebSocketURL              string
	RFQAssetAddressesCSV      string
	MakerAddress              string
	PrivateKey                string
	DummyPrice                string
	QuoteValidDurationSeconds int64
	AssetMapping              map[string]string // Maps asset addresses to underlying symbols (ETH, BTC, SOL)
}

// LoadConfig parses command-line flags and environment variables
// to populate and return an AppConfig struct.
func LoadConfig() *AppConfig {
	cfg := &AppConfig{}

	flag.StringVar(&cfg.WebSocketURL, "websocket_url", "wss://rip-testnet.rysk.finance/maker", "WebSocket URL for RFQ stream and quote submission")
	flag.StringVar(&cfg.RFQAssetAddressesCSV, "rfq_asset_addresses", "", "Comma-separated list of asset addresses for RFQ streams (e.g., 0xAsset1,0xAsset2)")
	flag.StringVar(&cfg.DummyPrice, "dummy_price", "1000000", "Dummy price to quote (ensure format matches Rysk requirements, e.g., units)")
	flag.Int64Var(&cfg.QuoteValidDurationSeconds, "quote_valid_duration_seconds", 30, "How long your quotes will be valid in seconds")
	flag.Parse()

	cfg.MakerAddress = os.Getenv("MAKER_ADDRESS")
	cfg.PrivateKey = os.Getenv("PRIVATE_KEY")

	if cfg.RFQAssetAddressesCSV == "" {
		log.Fatal("Error: --rfq_asset_addresses is required.")
	}
	if cfg.MakerAddress == "" {
		log.Fatal("Error: MAKER_ADDRESS environment variable is not set or empty.")
	}
	if cfg.PrivateKey == "" {
		log.Fatal("Error: PRIVATE_KEY environment variable is not set or empty.")
	}

	// Initialize asset mapping
	// TODO: This should be configurable via environment variables or config file
	cfg.AssetMapping = map[string]string{
		"0xb67bfa7b488df4f2efa874f4e59242e9130ae61f": "ETH", // Example mapping for testnet
		// Add more mappings as needed
	}

	return cfg
}
