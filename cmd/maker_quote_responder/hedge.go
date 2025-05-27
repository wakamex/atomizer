package main

import (
	"fmt"
	"log"
)

// HedgeOrder receives a filled order from RyskV12 and hedges it on the configured exchange
func HedgeOrder(conf RFQConfirmation, underlying string, cfg *AppConfig, exchange Exchange) error {
	// Use the exchange interface to place the hedge order
	err := exchange.PlaceHedgeOrder(conf, underlying, cfg)
	if err != nil {
		return fmt.Errorf("failed to place hedge order on %s: %w", cfg.ExchangeName, err)
	}
	
	log.Printf("[Hedge %s] Successfully placed hedge order on %s", conf.QuoteNonce, cfg.ExchangeName)
	return nil
}
