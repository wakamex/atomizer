package config

import (
	"crypto/ecdsa"
)

// Config holds all configuration for the application
type Config struct {
	// WebSocket configuration
	WebSocketURL              string
	RFQAssetAddressesCSV      string
	MakerAddress              string
	PrivateKey                string
	ParsedPrivateKey          *ecdsa.PrivateKey
	
	// Exchange configuration
	ExchangeName              string
	ExchangeTestMode          bool
	DeribitApiKey             string
	DeribitApiSecret          string
	
	// Trading configuration
	DummyPrice                string
	QuoteValidDurationSeconds int64
	AssetMapping              map[string]string
	
	// Risk and hedging configuration
	MaxPositionDelta          float64
	MinLiquidityScore         float64
	EnableGammaHedging        bool
	GammaThreshold            float64
	EnableManualTrades        bool
	
	// Infrastructure configuration
	HTTPPort                  string
	CacheBackend              string
	ValkeyAddr                string
	
	// Market maker specific configuration
	Underlying                string
	Expiry                    string
	Strikes                   []string
	AllStrikes                bool
	Spread                    int
	MinSpread                 int
	Size                      float64
	Improvement               float64
	ImprovementRefSize        float64
	MaxPosition               float64
	MaxExposure               float64
	RefreshInterval           int
	DryRun                    bool
	BidOnly                   bool
	AskOnly                   bool
}

// DefaultAssetMapping provides default asset mappings for known tokens
var DefaultAssetMapping = map[string]string{
	// Mainnet
	"0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2": "ETH",  // WETH
	
	// Sepolia Testnet
	"0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9": "ETH",  // Sepolia WETH
	"0xb67bfa7b488df4f2efa874f4e59242e9130ae61f": "ETH",  // Another testnet ETH
	
	// Add more mappings as needed
}