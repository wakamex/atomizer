package main

import (
	"fmt"
	"log"
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
		config := ExchangeConfig{
			APIKey:    cfg.DeribitApiKey,
			APISecret: cfg.DeribitApiSecret,
			TestMode:  cfg.ExchangeTestMode,
			RateLimit: 10,
		}
		return NewDeribitExchange(config), nil
		
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
	return []string{"deribit"} // Add more as they are implemented
}