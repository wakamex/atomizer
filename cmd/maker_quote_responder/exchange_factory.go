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
		// Check for Derive API credentials
		deriveAPIKey := os.Getenv("DERIVE_API_KEY")
		deriveAPISecret := os.Getenv("DERIVE_API_SECRET")
		
		if deriveAPIKey == "" || deriveAPISecret == "" {
			return nil, fmt.Errorf("Derive credentials not found. Set DERIVE_API_KEY and DERIVE_API_SECRET")
		}
		
		log.Printf("Using Derive exchange with API key authentication")
		config := ExchangeConfig{
			APIKey:    deriveAPIKey,
			APISecret: deriveAPISecret,
			TestMode:  cfg.ExchangeTestMode,
			RateLimit: 10,
		}
		return NewDeriveExchange(config), nil
		
	// Placeholder for future exchanges
	case "okx":
		// TODO: Implement OKX exchange when adding support
		return nil, fmt.Errorf("OKX exchange not yet implemented")
		
	case "bybit":
		// TODO: Implement Bybit exchange when adding support
		return nil, fmt.Errorf("Bybit exchange not yet implemented")
		
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", exchangeName)
	}
}

// GetSupportedExchanges returns a list of supported exchange names
func (f *ExchangeFactory) GetSupportedExchanges() []string {
	return []string{"deribit", "derive"} // Add more as they are implemented
}