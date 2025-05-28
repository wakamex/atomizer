package main

import (
	"fmt"
	"log"
)

// HedgeOrder receives a filled order from RyskV12 and hedges it on the configured exchange
// This is the legacy function - new code should use HedgeManager directly
func HedgeOrder(conf RFQConfirmation, underlying string, cfg *AppConfig, exchange Exchange) error {
	// Use the exchange interface to place the hedge order
	err := exchange.PlaceHedgeOrder(conf, underlying, cfg)
	if err != nil {
		return fmt.Errorf("failed to place hedge order on %s: %w", cfg.ExchangeName, err)
	}
	
	log.Printf("[Hedge %s] Successfully placed hedge order on %s", conf.QuoteNonce, cfg.ExchangeName)
	return nil
}

// HedgeOrderWithManager uses the new HedgeManager for hedging
func HedgeOrderWithManager(conf RFQConfirmation, underlying string, cfg *AppConfig, hedgeManager *HedgeManager) error {
	// Convert confirmation to TradeEvent
	trade := &TradeEvent{
		ID:         conf.QuoteNonce,
		RFQId:      conf.Nonce,
		IsTakerBuy: conf.IsTakerBuy,
		Quantity:   DecimalFromString(conf.Quantity), // Already in correct format
		// Other fields would be populated from the original RFQ
	}
	
	result, err := hedgeManager.ExecuteHedge(trade)
	if err != nil {
		return fmt.Errorf("failed to execute hedge: %w", err)
	}
	
	log.Printf("[Hedge %s] Successfully hedged: %+v", conf.QuoteNonce, result)
	return nil
}
