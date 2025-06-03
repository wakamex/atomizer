package hedging

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DeriveMarketDataAdapter adapts Derive API to the MarketData interface
type DeriveMarketDataAdapter struct {
	wsClient     *DeriveWSClient
	subaccountID uint64
}

func (d *DeriveMarketDataAdapter) GetTicker(instrumentName string) (*TickerData, error) {
	// For ETH-PERP
	if instrumentName == "ETH-PERP" {
		// Fetch actual ticker data
		perpTicker, err := d.fetchTickerFromAPI(instrumentName)
		if err != nil {
			// Fallback to default values
			return &TickerData{
				InstrumentName: instrumentName,
				MarkPrice:      decimal.NewFromFloat(3000),
				MinimumAmount:  decimal.NewFromFloat(0.01),
				AmountStep:     decimal.NewFromFloat(0.01),
				TickSize:       decimal.NewFromFloat(0.01), // 2 decimal places for perps
			}, nil
		}
		return perpTicker, nil
	}

	// For options, fetch ticker data which includes Greeks
	log.Printf("Fetching ticker for option: %s", instrumentName)
	ticker, err := d.fetchTickerFromAPI(instrumentName)
	if err != nil {
		log.Printf("Failed to fetch ticker for %s: %v", instrumentName, err)
		// Use default values for testing
		return &TickerData{
			InstrumentName: instrumentName,
			MarkPrice:      decimal.NewFromFloat(0.1),
			MinimumAmount:  decimal.NewFromFloat(0.1),
			AmountStep:     decimal.NewFromFloat(0.1),
			TickSize:       decimal.NewFromFloat(0.0001),
			OptionPricing: &OptionPricing{
				Delta: decimal.NewFromFloat(0.5),
				Gamma: decimal.NewFromFloat(0.01),
			},
			OptionDetails: &OptionDetails{
				Expiry: time.Now().Add(30 * 24 * time.Hour).Unix(),
			},
		}, nil
	}

	log.Printf("Successfully fetched ticker for %s", instrumentName)
	return ticker, nil
}

// fetchTickerFromAPI calls the public/get_ticker endpoint
func (d *DeriveMarketDataAdapter) fetchTickerFromAPI(instrumentName string) (*TickerData, error) {
	// Use the exchange's fetchTicker method which now includes Greeks
	exchange := &DeriveMarketMakerExchange{
		wsClient:     d.wsClient,
		subaccountID: d.subaccountID,
	}

	// Call the enhanced fetchTicker method
	ticker, err := exchange.fetchTicker(instrumentName)
	if err != nil {
		return nil, err
	}

	// Set proper tick size and amount step based on instrument type
	var tickSize, amountStep, minAmount decimal.Decimal
	if instrumentName == "ETH-PERP" {
		tickSize = decimal.NewFromFloat(0.01)   // 2 decimal places for perps
		amountStep = decimal.NewFromFloat(0.01) // 2 decimal places for amount
		minAmount = decimal.NewFromFloat(0.01)
	} else {
		tickSize = decimal.NewFromFloat(0.0001) // 4 decimal places for options
		amountStep = decimal.NewFromFloat(0.1)  // 1 decimal place for amount
		minAmount = decimal.NewFromFloat(0.1)
	}

	// Convert TickerUpdate to TickerData
	tickerData := &TickerData{
		InstrumentName: ticker.Instrument,
		MarkPrice:      ticker.MarkPrice,
		MinimumAmount:  minAmount,
		AmountStep:     amountStep,
		TickSize:       tickSize,
	}

	// Copy Greeks if available
	if ticker.Delta != nil && ticker.Gamma != nil {
		tickerData.OptionPricing = &OptionPricing{
			Delta: *ticker.Delta,
			Gamma: *ticker.Gamma,
		}
		if ticker.Vega != nil {
			tickerData.OptionPricing.Vega = *ticker.Vega
		}
		if ticker.Theta != nil {
			tickerData.OptionPricing.Theta = *ticker.Theta
		}

		log.Printf("Fetched Greeks for %s - Delta: %s, Gamma: %s",
			instrumentName, ticker.Delta.String(), ticker.Gamma.String())
	}

	// Copy option details
	if ticker.Expiry != nil {
		tickerData.OptionDetails = &OptionDetails{
			Expiry: ticker.Expiry.Unix(),
		}
	}

	return tickerData, nil
}

func (d *DeriveMarketDataAdapter) GetOrderbookExcludeMyOrders(instrumentName string) (*OrderbookData, error) {
	// Fetch real orderbook from Derive WebSocket
	ob := d.wsClient.GetOrderBook(instrumentName)
	if ob == nil {
		return nil, fmt.Errorf("no orderbook data for %s", instrumentName)
	}

	// Convert to OrderbookData format
	orderbook := &OrderbookData{
		Bids: [][]decimal.Decimal{},
		Asks: [][]decimal.Decimal{},
	}

	// Convert bids
	for _, bid := range ob.Bids {
		orderbook.Bids = append(orderbook.Bids, []decimal.Decimal{bid.Price, bid.Size})
	}

	// Convert asks
	for _, ask := range ob.Asks {
		orderbook.Asks = append(orderbook.Asks, []decimal.Decimal{ask.Price, ask.Size})
	}

	// If no orderbook data, subscribe and wait briefly
	if len(orderbook.Bids) == 0 && len(orderbook.Asks) == 0 {
		if err := d.wsClient.SubscribeOrderBook(instrumentName, 10); err != nil {
			return nil, fmt.Errorf("failed to subscribe to orderbook: %w", err)
		}
		// Wait a moment for initial data
		time.Sleep(500 * time.Millisecond)

		// Try again
		ob = d.wsClient.GetOrderBook(instrumentName)
		if ob != nil {
			// Re-convert with new data
			orderbook.Bids = nil
			orderbook.Asks = nil
			for _, bid := range ob.Bids {
				orderbook.Bids = append(orderbook.Bids, []decimal.Decimal{bid.Price, bid.Size})
			}
			for _, ask := range ob.Asks {
				orderbook.Asks = append(orderbook.Asks, []decimal.Decimal{ask.Price, ask.Size})
			}
		}
	}

	return orderbook, nil
}

func (d *DeriveMarketDataAdapter) GetOrders(instrumentName string) []GammaOrder {
	orders, err := d.wsClient.GetOpenOrders(d.subaccountID)
	if err != nil {
		return nil
	}

	var result []GammaOrder
	for _, order := range orders {
		if inst, ok := order["instrument_name"].(string); ok && inst == instrumentName {
			gOrder := parseGammaOrder(order)
			result = append(result, gOrder)
		}
	}
	return result
}

func (d *DeriveMarketDataAdapter) IterPositions() []GammaPosition {
	positions, err := d.wsClient.GetPositions(d.subaccountID)
	if err != nil {
		return nil
	}

	var result []GammaPosition
	for _, pos := range positions {
		if inst, ok := pos["instrument_name"].(string); ok {
			gp := GammaPosition{InstrumentName: inst}
			if amountStr, ok := pos["amount"].(string); ok {
				gp.Amount, _ = decimal.NewFromString(amountStr)
			}
			result = append(result, gp)
		}
	}
	return result
}

// DeriveWsClientAdapter adapts DeriveWSClient to WsClient interface
type DeriveWsClientAdapter struct {
	wsClient     *DeriveWSClient
	subaccountID uint64
}

func (c *DeriveWsClientAdapter) Login() error {
	return nil // Already logged in
}

func (c *DeriveWsClientAdapter) EnableCancelOnDisconnect() error {
	return nil
}

func (c *DeriveWsClientAdapter) SendOrder(ticker *TickerData, subaccountID int64, args OrderArgs) error {
	// Create a temporary exchange instance to place the order
	exchange := &DeriveMarketMakerExchange{
		wsClient:        c.wsClient,
		subaccountID:    uint64(subaccountID),
		privateKey:      "",                                        // Not needed, auth is in wsClient
		instrumentCache: make(map[string]*DeriveInstrumentDetails), // Initialize the map
		subscriptions:   make(map[string]bool),
	}

	_, err := exchange.PlaceLimitOrder(ticker.InstrumentName, args.Direction.String(), args.LimitPrice, args.Amount)
	return err
}

func (c *DeriveWsClientAdapter) SendReplace(ticker *TickerData, subaccountID int64, cancelID uuid.UUID, args OrderArgs) error {
	// For now, just place new order
	// TODO: Implement proper replace when we add CancelOrder method
	return c.SendOrder(ticker, subaccountID, args)
}

func (c *DeriveWsClientAdapter) CancelAll(subaccountID int64) error {
	// TODO: Implement CancelAllOrders
	return nil
}

func (c *DeriveWsClientAdapter) CancelByLabel(subaccountID int64, label string) error {
	orders, err := c.wsClient.GetOpenOrders(uint64(subaccountID))
	if err != nil {
		return err
	}

	for _, order := range orders {
		if l, ok := order["label"].(string); ok && l == label {
			if id, ok := order["order_id"].(string); ok {
				// TODO: Implement CancelOrder
				_ = id
			}
		}
	}
	return nil
}

func (c *DeriveWsClientAdapter) Ping() error {
	return nil
}

func parseGammaOrder(order map[string]interface{}) GammaOrder {
	gOrder := GammaOrder{Status: "open"}

	if id, ok := order["order_id"].(string); ok {
		gOrder.OrderID = id
	}
	if dir, ok := order["direction"].(string); ok {
		if dir == "buy" {
			gOrder.Direction = Buy
		} else {
			gOrder.Direction = Sell
		}
	}
	if price, ok := order["price"].(string); ok {
		gOrder.Price, _ = decimal.NewFromString(price)
	}
	if amount, ok := order["amount"].(string); ok {
		gOrder.Amount, _ = decimal.NewFromString(amount)
	}
	if label, ok := order["label"].(string); ok {
		gOrder.Label = label
	}

	return gOrder
}

// RunGammaHedger runs the gamma hedger as a standalone process
func RunGammaHedger(args []string) {
	fs := flag.NewFlagSet("gamma-hedger", flag.ExitOnError)

	// Configuration flags
	privateKey := fs.String("private-key", os.Getenv("DERIVE_PRIVATE_KEY"), "Private key")
	wallet := fs.String("wallet", os.Getenv("DERIVE_WALLET"), "Wallet address")
	subaccountID := fs.String("subaccount", os.Getenv("DERIVE_SUBACCOUNT_ID"), "Subaccount ID")
	gammaThreshold := fs.Float64("gamma-threshold", 0.1, "Max absolute delta before hedging")
	maxSpread := fs.Float64("max-spread", 0.002, "Max spread for hedge orders")
	actionWaitMs := fs.Uint64("wait-ms", 1000, "Milliseconds between hedge checks")
	debug := fs.Bool("debug", false, "Enable debug logging")

	if err := fs.Parse(args); err != nil {
		log.Fatal(err)
	}

	// Set debug mode
	SetDebugMode(*debug)

	// Validate
	if *privateKey == "" {
		log.Fatal("Private key required (--private-key or DERIVE_PRIVATE_KEY)")
	}
	if *wallet == "" {
		log.Fatal("Wallet required (--wallet or DERIVE_WALLET)")
	}

	// Parse subaccount ID
	var subID uint64
	if *subaccountID != "" {
		parsed, err := strconv.ParseUint(*subaccountID, 10, 64)
		if err != nil {
			log.Fatalf("Invalid subaccount ID: %v", err)
		}
		subID = parsed
	}

	log.Printf("Starting Gamma Hedger")
	log.Printf("Configuration:")
	log.Printf("  Subaccount: %d", subID)
	log.Printf("  Gamma Threshold: %.4f", *gammaThreshold)
	log.Printf("  Max Spread: %.4f", *maxSpread)
	log.Printf("  Check Interval: %dms", *actionWaitMs)

	// Create exchange using the standard factory
	exchange, err := NewDeriveMarketMakerExchange(*privateKey, *wallet)
	if err != nil {
		log.Fatalf("Failed to create exchange: %v", err)
	}
	wsClient := exchange.wsClient
	defer wsClient.Close()

	// Get default subaccount if not specified
	if subID == 0 {
		subID = wsClient.GetDefaultSubaccount()
		log.Printf("Using default subaccount: %d", subID)
	}

	// Create adapters
	marketData := &DeriveMarketDataAdapter{
		wsClient:     wsClient,
		subaccountID: subID,
	}

	wsAdapter := &DeriveWsClientAdapter{
		wsClient:     wsClient,
		subaccountID: subID,
	}

	// Create gamma algorithm
	algo := &GammaDDHAlgo{
		SubaccountID: int64(subID),
		PerpName:     "ETH-PERP",
		MaxAbsDelta:  decimal.NewFromFloat(*gammaThreshold),
		MaxAbsSpread: decimal.NewFromFloat(*maxSpread),
		ActionWaitMS: *actionWaitMs,
		PriceTol:     decimal.NewFromFloat(0.0001),
		AmountTol:    decimal.NewFromFloat(0.01),
		MidPriceTol:  decimal.NewFromFloat(0.0001),
	}

	// Subscribe to ETH-PERP orderbook
	log.Printf("Subscribing to ETH-PERP orderbook...")
	if err := wsClient.SubscribeOrderBook("ETH-PERP", 10); err != nil {
		log.Printf("Warning: Failed to subscribe to ETH-PERP orderbook: %v", err)
	}

	// Wait for initial orderbook data
	time.Sleep(1 * time.Second)

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Run hedger
	go func() {
		if err := algo.StartHedger(ctx, marketData, wsAdapter); err != nil {
			log.Printf("Hedger error: %v", err)
			cancel()
		}
	}()

	// Wait for shutdown
	<-sigChan
	log.Println("Shutting down gamma hedger...")
}
