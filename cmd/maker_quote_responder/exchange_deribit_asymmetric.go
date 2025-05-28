package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"
)

// DeribitAsymmetricExchange implements the Exchange interface using Ed25519 authentication
type DeribitAsymmetricExchange struct {
	client *DeribitClient
	config ExchangeConfig
}

// NewDeribitAsymmetricExchange creates a new Deribit exchange with asymmetric key authentication
func NewDeribitAsymmetricExchange(config ExchangeConfig, clientID string, privateKeyPEM string) (*DeribitAsymmetricExchange, error) {
	client, err := NewDeribitClient(clientID, privateKeyPEM, config.TestMode)
	if err != nil {
		return nil, fmt.Errorf("failed to create Deribit client: %w", err)
	}

	return &DeribitAsymmetricExchange{
		client: client,
		config: config,
	}, nil
}

// GetOrderBook fetches the order book for a given option
func (d *DeribitAsymmetricExchange) GetOrderBook(req RFQResult, asset string) (CCXTOrderBook, error) {
	// Convert option details to instrument name
	instrumentName, err := d.ConvertToInstrument(asset, req.Strike, req.Expiry, req.IsPut)
	if err != nil {
		return CCXTOrderBook{}, err
	}

	// Fetch order book using native client
	orderBook, err := d.client.GetOrderBook(instrumentName)
	if err != nil {
		return CCXTOrderBook{}, fmt.Errorf("failed to fetch order book: %v", err)
	}

	// Convert to CCXT format
	book := CCXTOrderBook{
		Bids:  orderBook.Bids,
		Asks:  orderBook.Asks,
		Index: orderBook.UnderlyingPrice,
	}

	// If no underlying price, try to get it from index
	if book.Index == 0 && orderBook.IndexPrice > 0 {
		book.Index = orderBook.IndexPrice
	}

	return book, nil
}

// PlaceHedgeOrder places a hedge order on Deribit
// Since Rysk users are always selling calls (we buy from them), we hedge by selling calls on Deribit
func (d *DeribitAsymmetricExchange) PlaceHedgeOrder(conf RFQConfirmation, instrument string, cfg *AppConfig) error {
	// The instrument parameter is already in the correct format
	// No need to convert again

	// Convert quantity from wei to decimal
	quantityFloat, err := strconv.ParseFloat(conf.Quantity, 64)
	if err != nil {
		return fmt.Errorf("failed to parse quantity: %w", err)
	}
	quantityETH := quantityFloat / math.Pow(10, 18) // Convert from wei to ETH

	// Get current order book to find best ask
	orderBook, err := d.client.GetOrderBook(instrument)
	if err != nil {
		return fmt.Errorf("failed to fetch order book for hedge: %w", err)
	}
	
	// Get best ask price
	if len(orderBook.Asks) == 0 {
		return fmt.Errorf("no asks available in order book")
	}
	bestAsk := orderBook.Asks[0][0] // First ask price
	
	// Since Rysk users are always selling calls to us (we buy from them),
	// we always need to hedge by selling calls on Deribit
	orderSide := "sell"
	
	// Place our ask at 2x the best ask for safety (far from top of book)
	hedgePrice := bestAsk * 2.0
	
	log.Printf("[Hedge] Best ask: %f, placing our ask at: %f (2x best ask)", bestAsk, hedgePrice)

	// Place the hedge order
	order, err := d.client.PlaceOrder(instrument, quantityETH, "limit", orderSide, hedgePrice)
	if err != nil {
		return fmt.Errorf("failed to place hedge order: %w", err)
	}

	log.Printf("[Hedge] Order placed successfully - ID: %s, Side: %s, Quantity: %f ETH, Price: %f",
		order.OrderID, orderSide, quantityETH, hedgePrice)

	return nil
}

// ConvertToInstrument converts option details to Deribit-specific instrument format
func (d *DeribitAsymmetricExchange) ConvertToInstrument(asset string, strike string, expiry int64, isPut bool) (string, error) {
	// This is the same as the standard Deribit implementation
	return convertDeribitInstrument(asset, strike, expiry, isPut)
}

// Helper function shared with standard Deribit implementation
func convertDeribitInstrument(asset string, strike string, expiry int64, isPut bool) (string, error) {
	// Convert the strike from a big.Int string to a normal number
	strikeBigInt, ok := new(big.Int).SetString(strike, 10)
	if !ok {
		return "", fmt.Errorf("invalid strike")
	}
	strike = strikeBigInt.Div(strikeBigInt, new(big.Int).SetUint64(1e8)).String()
	
	// Convert the expiry from a timestamp seconds into a Deribit compatible date time
	expiryTime := time.Unix(expiry, 0)
	deribitExpiry := strings.ToUpper(expiryTime.Format("2Jan06"))
	
	// Convert isPut to "C" or "P"
	optionType := "C"
	if isPut {
		optionType = "P"
		// Note: Remove this restriction if puts are supported
		return "", fmt.Errorf("puts not supported")
	}
	
	// Construct the instrument name in Deribit format: ASSET-EXPIRY-STRIKE-TYPE
	instrumentName := asset + "-" + deribitExpiry + "-" + strike + "-" + optionType
	return instrumentName, nil
}