package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// ExchangeFactory creates Exchange instances based on configuration
type ExchangeFactory struct{}

// NewExchangeFactory creates a new exchange factory
func NewExchangeFactory() *ExchangeFactory {
	return &ExchangeFactory{}
}

// CreateExchange creates an exchange instance based on the provided configuration
func (f *ExchangeFactory) CreateExchange(cfg *AppConfig) (Exchange, error) {
	exchangeName := strings.ToLower(cfg.ExchangeName)
	
	log.Printf("Creating exchange: %s (test mode: %v)", exchangeName, cfg.ExchangeTestMode)
	
	switch exchangeName {
	case "deribit":
		// Check if we have asymmetric key credentials
		clientID := os.Getenv("DERIBIT_CLIENT_ID")
		privateKey := os.Getenv("DERIBIT_PRIVATE_KEY")
		privateKeyFile := os.Getenv("DERIBIT_PRIVATE_KEY_FILE")
		
		// If we have client ID and private key, use asymmetric authentication
		if clientID != "" && (privateKey != "" || privateKeyFile != "") {
			log.Printf("Using Deribit with asymmetric key authentication (Client ID: %s)", clientID)
			
			// Read private key from file if specified
			if privateKeyFile != "" && privateKey == "" {
				keyData, err := ioutil.ReadFile(privateKeyFile)
				if err != nil {
					return nil, fmt.Errorf("failed to read private key file: %w", err)
				}
				privateKey = string(keyData)
			}
			
			config := ExchangeConfig{
				TestMode:  cfg.ExchangeTestMode,
				RateLimit: 10,
			}
			return NewDeribitAsymmetricExchange(config, clientID, privateKey)
		}
		
		// Otherwise use standard API key/secret with CCXT
		if cfg.DeribitApiKey == "" || cfg.DeribitApiSecret == "" {
			return nil, fmt.Errorf("Deribit credentials not found. Set either DERIBIT_CLIENT_ID + DERIBIT_PRIVATE_KEY for asymmetric auth, or DERIBIT_API_KEY + DERIBIT_API_SECRET for standard auth")
		}
		
		log.Printf("Using Deribit with standard API key authentication")
		config := ExchangeConfig{
			APIKey:    cfg.DeribitApiKey,
			APISecret: cfg.DeribitApiSecret,
			TestMode:  cfg.ExchangeTestMode,
			RateLimit: 10,
		}
		return NewDeribitExchange(config), nil
		
	case "derive":
		// Use Derive-specific private key
		privateKey := os.Getenv("DERIVE_PRIVATE_KEY")
		if privateKey != "" {
			log.Printf("Using DERIVE_PRIVATE_KEY for Derive exchange")
		} else {
			// Fallback to general private key
			privateKey = os.Getenv("PRIVATE_KEY")
			log.Printf("Using PRIVATE_KEY for Derive exchange (DERIVE_PRIVATE_KEY not set)")
		}
		
		deriveWallet := os.Getenv("DERIVE_WALLET_ADDRESS")
		
		if privateKey == "" || deriveWallet == "" {
			return nil, fmt.Errorf("Derive credentials not found. Set DERIVE_PRIVATE_KEY (or PRIVATE_KEY) and DERIVE_WALLET_ADDRESS")
		}
		
		log.Printf("Using Derive exchange with wallet: %s", deriveWallet)
		config := ExchangeConfig{
			APIKey:    privateKey,
			APISecret: "", // Derive uses private key signing, not API secret
			TestMode:  cfg.ExchangeTestMode,
			RateLimit: 10,
		}
		// Use Derive-specific CCXT wrapper
		return NewCCXTDeriveExchange(config)
		
	case "okx":
		// Check for OKX API credentials
		okxAPIKey := os.Getenv("OKX_API_KEY")
		okxAPISecret := os.Getenv("OKX_API_SECRET")
		okxPassphrase := os.Getenv("OKX_PASSPHRASE")
		
		if okxAPIKey == "" || okxAPISecret == "" || okxPassphrase == "" {
			return nil, fmt.Errorf("OKX credentials not found. Set OKX_API_KEY, OKX_API_SECRET, and OKX_PASSPHRASE")
		}
		
		log.Printf("Using OKX exchange with API key authentication")
		// config := ExchangeConfig{
		// 	APIKey:    okxAPIKey,
		// 	APISecret: okxAPISecret,
		// 	TestMode:  cfg.ExchangeTestMode,
		// 	RateLimit: 10,
		// }
		// TODO: Implement OKX with CCXT
		return nil, fmt.Errorf("OKX implementation pending")
		
	case "bybit":
		// Check for Bybit API credentials
		bybitAPIKey := os.Getenv("BYBIT_API_KEY")
		bybitAPISecret := os.Getenv("BYBIT_API_SECRET")
		
		if bybitAPIKey == "" || bybitAPISecret == "" {
			return nil, fmt.Errorf("Bybit credentials not found. Set BYBIT_API_KEY and BYBIT_API_SECRET")
		}
		
		log.Printf("Using Bybit exchange with API key authentication")
		_ = ExchangeConfig{
			APIKey:    bybitAPIKey,
			APISecret: bybitAPISecret,
			TestMode:  cfg.ExchangeTestMode,
			RateLimit: 10,
		}
		// TODO: Implement Bybit with CCXT
		return nil, fmt.Errorf("Bybit implementation pending")
		
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", exchangeName)
	}
}

// GetSupportedExchanges returns a list of supported exchange names
func (f *ExchangeFactory) GetSupportedExchanges() []string {
	return []string{"deribit", "derive", "okx", "bybit"}
}