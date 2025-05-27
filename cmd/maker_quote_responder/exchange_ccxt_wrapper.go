package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

// CCXTDeriveExchange wraps the CCXT Derive implementation
type CCXTDeriveExchange struct {
	exchange *ccxt.Derive
	config   ExchangeConfig
}

// NewCCXTDeriveExchange creates a new Derive exchange using CCXT
func NewCCXTDeriveExchange(config ExchangeConfig) (*CCXTDeriveExchange, error) {
	// Initialize Derive with CCXT using private key authentication
	exchangeConfig := map[string]interface{}{
		"rateLimit":       config.RateLimit,
		"enableRateLimit": true,
		"options": map[string]interface{}{
			"defaultType": "option",
		},
	}
	
	// Derive uses privateKey for authentication
	if config.APIKey != "" {
		exchangeConfig["privateKey"] = config.APIKey
	}
	
	exchange := ccxt.NewDerive(exchangeConfig)
	
	return &CCXTDeriveExchange{
		exchange: &exchange,
		config:   config,
	}, nil
}

// GetOrderBook fetches the order book for a given option
func (d *CCXTDeriveExchange) GetOrderBook(req RFQResult, asset string) (CCXTOrderBook, error) {
	// Convert to instrument
	instrument, err := d.ConvertToInstrument(asset, req.Strike, req.Expiry, req.IsPut)
	if err != nil {
		return CCXTOrderBook{}, err
	}
	
	// Format symbol - adjust based on Derive's actual format
	symbol := fmt.Sprintf("%s:%s", asset, instrument)
	
	// Fetch order book
	orderBookChan := d.exchange.FetchOrderBook(symbol)
	orderBookRaw := <-orderBookChan
	
	// Check for errors
	if err, ok := orderBookRaw.(error); ok {
		return CCXTOrderBook{}, err
	}
	
	// Convert to OrderBook type
	orderBook, ok := orderBookRaw.(ccxt.OrderBook)
	if !ok {
		return CCXTOrderBook{}, fmt.Errorf("unexpected order book type: %T", orderBookRaw)
	}
	
	// Get underlying price
	indexPrice := 0.0
	
	// Try option ticker first
	optionTicker, err := d.exchange.FetchTicker(symbol)
	if err == nil && optionTicker.Info != nil {
		if underlyingPrice, exists := optionTicker.Info["underlying_price"]; exists {
			if price, ok := underlyingPrice.(float64); ok {
				indexPrice = price
			}
		}
	}
	
	// If no underlying price, try spot/perpetual
	if indexPrice == 0.0 {
		underlyingSymbols := []string{
			asset + "/USD",
			asset + "/USDT",
			asset + "-PERPETUAL",
		}
		
		for _, sym := range underlyingSymbols {
			ticker, err := d.exchange.FetchTicker(sym)
			if err == nil && ticker.Last != nil {
				indexPrice = *ticker.Last
				break
			}
		}
	}
	
	return CCXTOrderBook{
		Bids:  orderBook.Bids,
		Asks:  orderBook.Asks,
		Index: indexPrice,
	}, nil
}

// PlaceHedgeOrder places a hedge order on Derive
func (d *CCXTDeriveExchange) PlaceHedgeOrder(conf RFQConfirmation, underlying string, cfg *AppConfig) error {
	// Convert to instrument
	instrument, err := d.ConvertToInstrument(underlying, conf.Strike, int64(conf.Expiry), conf.IsPut)
	if err != nil {
		return fmt.Errorf("failed to convert instrument: %w", err)
	}
	
	// Convert quantity from wei
	quantityFloat, err := strconv.ParseFloat(conf.Quantity, 64)
	if err != nil {
		return fmt.Errorf("failed to parse quantity: %w", err)
	}
	quantity := quantityFloat / math.Pow(10, 18)
	
	// Format symbol
	symbol := fmt.Sprintf("%s:%s", underlying, instrument)
	
	// Get current order book
	orderBookChan := d.exchange.FetchOrderBook(symbol)
	orderBookRaw := <-orderBookChan
	
	// Check for errors
	if err, ok := orderBookRaw.(error); ok {
		return fmt.Errorf("failed to fetch order book: %w", err)
	}
	
	// Convert to OrderBook type
	orderBook, ok := orderBookRaw.(ccxt.OrderBook)
	if !ok {
		return fmt.Errorf("unexpected order book type: %T", orderBookRaw)
	}
	
	if len(orderBook.Asks) == 0 {
		return fmt.Errorf("no asks available")
	}
	
	bestAsk := orderBook.Asks[0][0]
	hedgePrice := bestAsk * 2.0
	
	log.Printf("[Hedge] Derive best ask: %f, placing at: %f (2x)", bestAsk, hedgePrice)
	
	// Place order
	order, err := d.exchange.CreateOrder(
		symbol,
		"limit",
		"sell",
		quantity,
		ccxt.WithCreateOrderPrice(hedgePrice),
	)
	
	if err != nil {
		return fmt.Errorf("failed to place order: %w", err)
	}
	
	log.Printf("[Hedge] Order placed on Derive - ID: %s, Quantity: %f, Price: %f",
		order.Id, quantity, hedgePrice)
	
	return nil
}

// ConvertToInstrument converts option details to Derive format
func (d *CCXTDeriveExchange) ConvertToInstrument(asset string, strike string, expiry int64, isPut bool) (string, error) {
	// Convert strike from wei
	strikeBigInt, ok := new(big.Int).SetString(strike, 10)
	if !ok {
		return "", fmt.Errorf("invalid strike")
	}
	strikeNum := strikeBigInt.Div(strikeBigInt, new(big.Int).SetUint64(1e8)).String()
	
	// Format expiry - adjust based on Derive's actual format
	expiryTime := time.Unix(expiry, 0)
	expiryStr := expiryTime.Format("20060102") // YYYYMMDD
	
	// Option type
	optionType := "C"
	if isPut {
		optionType = "P"
		return "", fmt.Errorf("puts not supported")
	}
	
	// Build instrument - adjust based on Derive's actual format
	// Example: "ETH-20250131-3000-C"
	return fmt.Sprintf("%s-%s-%s-%s", asset, expiryStr, strikeNum, optionType), nil
}