package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

// HedgeOrder receives a filled order from RyskV12 and hedges it on Deribit
func HedgeOrder(conf RFQConfirmation, underlying string, cfg *AppConfig) error {
	// Initialize Deribit exchange with API credentials
	exchange := ccxt.NewDeribit(map[string]interface{}{
		"rateLimit":       10,
		"enableRateLimit": true,
		"apiKey":          cfg.DeribitApiKey,    // This is your Deribit Client ID
		"secret":          cfg.DeribitApiSecret, // This is your Deribit Client Secret
		"urls": map[string]interface{}{
			"api": map[string]interface{}{
				"rest": "https://test.deribit.com", // Using testnet
			},
		},
		"options": map[string]interface{}{
			"defaultType":             "option", // Required for options trading
			"adjustForTimeDifference": true,
			"recvWindow":              5000,
		},
	})

	// Convert option details to Deribit instrument name
	instrument, err := convertOptionDetailsToInstrument(underlying, conf.Strike, int64(conf.Expiry), conf.IsPut)
	if err != nil {
		return fmt.Errorf("failed to convert option details: %w", err)
	}

	// Convert quantity from wei to decimal
	quantityFloat, err := strconv.ParseFloat(conf.Quantity, 64)
	if err != nil {
		return fmt.Errorf("failed to parse quantity: %w", err)
	}
	quantityETH := quantityFloat / math.Pow(10, 18) // Convert from wei to ETH

	// Calculate price with slippage
	priceBigInt, ok := new(big.Int).SetString(conf.Price, 10)
	if !ok {
		return fmt.Errorf("invalid price")
	}
	price := priceBigInt.Div(priceBigInt, new(big.Int).SetUint64(1e13)).String()
	// convert the price string to a float
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return fmt.Errorf("failed to parse price: %w", err)
	}
	// convert the priceFloat to the correct units
	priceFloat = priceFloat / 1e5
	// Determine order side based on whether maker needs to buy or sell to hedge
	side := "sell" 
	if conf.IsTakerBuy { // If taker is selling, maker needs to buy to hedge
		side = "buy"
		priceFloat = (priceFloat * 90) / 100 // 10% lower for buy orders for slippage
	} else {
		priceFloat = (priceFloat * 110) / 100 // 10% higher for sell orders for slippage
	}

	// Place hedge order on Deribit with required parameters
	params := map[string]interface{}{
		"instrument_name": instrument,
		"amount":          quantityETH,
		"type":            "limit",
		"price":           priceFloat,
		"time_in_force":   "fill_or_kill",
		"advanced":        "usd",
	}
	// log the params
	order, err := exchange.CreateOrder(instrument, "limit", side, quantityETH, ccxt.WithCreateOrderParams(params))
	if err != nil {
		return fmt.Errorf("failed to place hedge order: %w", err)
	}
	log.Printf("[Hedge %s] Placed %s order for %f contracts of %s at price %f", conf.QuoteNonce, side, quantityETH, instrument, priceFloat)
	log.Printf("[Hedge %s] Order: %+v", conf.QuoteNonce, order)
	// check if the order is filled by looking at the response for the order_state on the order
	if order.Remaining == nil {
		return fmt.Errorf("order remaining is nil")
	}
	if *order.Remaining != 0 {
		log.Printf("[Hedge %s] Order not filled", conf.QuoteNonce)
		return fmt.Errorf("order not filled")
	}
	log.Printf("[Hedge %s] Order filled", conf.QuoteNonce)
	return nil
}
