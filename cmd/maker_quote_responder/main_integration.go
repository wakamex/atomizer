package main

import (
	"log"
)

// InitializeArbitrageSystem sets up the arbitrage components
func InitializeArbitrageSystem(cfg *AppConfig, exchange Exchange) (*ArbitrageOrchestrator, error) {
	// Create the arbitrage orchestrator
	orchestrator := NewArbitrageOrchestrator(cfg, exchange, []ExchangePosition{})
	
	// Start the orchestrator
	if err := orchestrator.Start(); err != nil {
		return nil, err
	}
	
	// Start HTTP API if enabled
	if cfg.EnableManualTrades {
		// HTTP server is now started in main.go
		log.Printf("Manual trades enabled via HTTP API")
	}
	
	return orchestrator, nil
}

// IntegrateWithRFQFlow modifies the RFQ processing to use the orchestrator
func IntegrateWithRFQFlow(orchestrator *ArbitrageOrchestrator) func(RFQResult) {
	return func(rfq RFQResult) {
		// Submit RFQ to orchestrator for tracking
		trade, err := orchestrator.SubmitRFQTrade(rfq)
		if err != nil {
			log.Printf("Failed to submit RFQ to orchestrator: %v", err)
			return
		}
		
		log.Printf("RFQ %s submitted to orchestrator as trade %s", rfq.ID, trade.ID)
	}
}

// IntegrateWithConfirmationFlow handles trade confirmations through the orchestrator
func IntegrateWithConfirmationFlow(orchestrator *ArbitrageOrchestrator) func(RFQConfirmation) {
	return func(conf RFQConfirmation) {
		orchestrator.OnRFQConfirmation(conf)
	}
}

// CreateCacheWithConfig creates the appropriate cache based on configuration
func CreateCacheWithConfig(cfg *AppConfig) (MarketCache, error) {
	switch cfg.CacheBackend {
	case "valkey":
		log.Printf("Using Valkey cache at %s", cfg.ValkeyAddr)
		return NewValkeyMarketCache(cfg.ValkeyAddr)
	default:
		log.Printf("Using file cache in ./cache directory")
		return NewFileMarketCache("./cache")
	}
}

// Example of how to modify main() to use the arbitrage system:
/*
func main() {
	// ... existing initialization code ...
	
	// Create cache with configuration
	cache, err := CreateCacheWithConfig(cfg)
	if err != nil {
		log.Fatalf("Failed to create cache: %v", err)
	}
	
	// Create exchange with cache
	// (modify exchange creation to use cache)
	
	// Initialize arbitrage system
	orchestrator, err := InitializeArbitrageSystem(cfg, exchange)
	if err != nil {
		log.Fatalf("Failed to initialize arbitrage system: %v", err)
	}
	defer orchestrator.Stop()
	
	// ... rest of main function ...
	
	// In the RFQ handling section, integrate with orchestrator:
	// When processing RFQ requests:
	onRFQReceived := IntegrateWithRFQFlow(orchestrator)
	
	// When processing confirmations:
	onConfirmationReceived := IntegrateWithConfirmationFlow(orchestrator)
}
*/