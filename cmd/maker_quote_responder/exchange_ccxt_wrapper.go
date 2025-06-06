package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

// CCXTDeriveExchange wraps the CCXT Derive implementation
type CCXTDeriveExchange struct {
	exchange       *ccxt.Derive
	config         ExchangeConfig
	deriveMarkets  map[string]DeriveInstrument // Custom market data from Derive API
	wsClient       *DeriveWSClient            // Persistent WebSocket connection
	wsClientMu     sync.Mutex                 // Mutex to protect wsClient access
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
	
	// Load markets to ensure connection is working
	log.Printf("[Derive] Loading markets...")
	marketsChan := exchange.LoadMarkets()
	marketsRaw := <-marketsChan
	
	if err, ok := marketsRaw.(error); ok {
		return nil, fmt.Errorf("failed to load Derive markets: %w", err)
	}
	
	// After LoadMarkets, the markets should be available in exchange.Markets
	log.Printf("[Derive] Markets loaded successfully")
	
	// Check if we need to handle pagination
	// Many exchanges limit to 100-1000 markets per request
	if len(exchange.Markets) == 100 {
		log.Printf("[Derive] Exactly 100 markets loaded - this might be paginated. Attempting to load more...")
		
		// Some exchanges support reloading with parameters
		// Try common pagination approaches
		reloadChan := exchange.LoadMarkets(true) // Force reload
		reloadRaw := <-reloadChan
		if err, ok := reloadRaw.(error); ok {
			log.Printf("[Derive] Warning: Could not reload markets: %v", err)
		} else {
			log.Printf("[Derive] After reload: %d markets available", len(exchange.Markets))
		}
	}
	
	// Check the Markets field
	if exchange.Markets != nil {
		optionCount := 0
		var sampleOptions []string
		
		for symbol, marketRaw := range exchange.Markets {
			// Try to convert to MarketInterface
			if market, ok := marketRaw.(ccxt.MarketInterface); ok {
				if market.Type != nil && *market.Type == "option" {
					optionCount++
					// Log first few option symbols as examples
					if len(sampleOptions) < 10 {
						sampleOptions = append(sampleOptions, symbol)
					}
				}
			}
		}
		
		log.Printf("[Derive] Found %d markets (%d options)", len(exchange.Markets), optionCount)
		if len(sampleOptions) > 0 {
			log.Printf("[Derive] Sample option symbols: %v", sampleOptions)
		}
	} else {
		log.Printf("[Derive] Warning: Markets field is nil after loading")
	}
	
	// Load all markets using cached loader
	cache, err := NewFileMarketCache("./cache")
	if err != nil {
		log.Printf("[Derive] Warning: Failed to create cache: %v", err)
	}
	
	// Cache markets for 1 hour
	cachedLoader := NewCachedMarketLoader(cache, 1*time.Hour)
	
	log.Printf("[Derive] Loading all markets (with caching)...")
	deriveMarkets, err := cachedLoader.LoadDeriveMarkets()
	if err != nil {
		log.Printf("[Derive] Warning: Failed to load markets: %v", err)
		deriveMarkets = make(map[string]DeriveInstrument)
	} else {
		log.Printf("[Derive] Successfully loaded %d total markets", len(deriveMarkets))
		
		// Count ETH options
		ethOptionCount := 0
		for _, market := range deriveMarkets {
			if market.BaseCurrency == "ETH" && market.InstrumentType == "option" {
				ethOptionCount++
			}
		}
		log.Printf("[Derive] Found %d ETH options", ethOptionCount)
	}
	
	deriveExchange := &CCXTDeriveExchange{
		exchange:      &exchange,
		config:        config,
		deriveMarkets: deriveMarkets,
	}
	
	// Initialize persistent WebSocket connection if we have the necessary config
	if err := deriveExchange.initializeWebSocket(); err != nil {
		log.Printf("[Derive] Warning: Failed to initialize persistent WebSocket: %v", err)
		// Continue without persistent connection - will fall back to creating new connections
	}
	
	return deriveExchange, nil
}

// GetOrderBook fetches the order book for a given option
func (d *CCXTDeriveExchange) GetOrderBook(req RFQResult, asset string) (CCXTOrderBook, error) {
	log.Printf("[Derive] GetOrderBook called for RFQ: asset=%s, strike=%s, expiry=%d, isPut=%v",
		asset, req.Strike, req.Expiry, req.IsPut)
	
	// Convert to instrument
	instrument, err := d.ConvertToInstrument(asset, req.Strike, req.Expiry, req.IsPut)
	if err != nil {
		return CCXTOrderBook{}, err
	}
	
	// The instrument already contains the full symbol
	symbol := instrument
	
	// Convert to CCXT format for checking markets
	// Derive API format: ETH-20250627-3400-C
	// CCXT format: ETH/USDC:USDC-25-06-27-3400-C
	ccxtSymbol := ""
	if parts := strings.Split(instrument, "-"); len(parts) == 4 {
		asset := parts[0]
		dateStr := parts[1] // YYYYMMDD
		strike := parts[2]
		optType := parts[3]
		
		// Convert date format
		if len(dateStr) == 8 {
			year := dateStr[2:4]  // YY
			month := dateStr[4:6] // MM
			day := dateStr[6:8]   // DD
			ccxtSymbol = fmt.Sprintf("%s/USDC:USDC-%s-%s-%s-%s-%s", asset, year, month, day, strike, optType)
			log.Printf("[Derive] Converted instrument %s to CCXT format: %s", instrument, ccxtSymbol)
		}
	}
	
	// Check if the market exists in our loaded markets
	instrumentName := ""
	ccxtExists := false
	if ccxtSymbol != "" {
		if _, exists := d.exchange.Markets[ccxtSymbol]; exists {
			ccxtExists = true
			log.Printf("[Derive] Found symbol in CCXT markets: %s", ccxtSymbol)
		}
	}
	
	if !ccxtExists {
		log.Printf("[Derive] Symbol %s not found in CCXT loaded markets.", symbol)
		
		// Check our custom loaded markets
		// The instrument variable already contains the Derive format (ETH-20250627-3800-C)
		if market, exists := d.deriveMarkets[instrument]; exists {
			instrumentName = instrument
			log.Printf("[Derive] Found instrument in custom loaded markets: %s", instrumentName)
			log.Printf("[Derive] Instrument details: base=%s, strike=%s, type=%s", 
				market.BaseCurrency, market.OptionDetails.Strike, market.OptionDetails.OptionType)
		} else {
			log.Printf("[Derive] Warning: Instrument %s not found in custom loaded markets. Will try anyway.", instrument)
			
			// Log some similar instruments for debugging
			log.Printf("[Derive] Looking for similar instruments...")
			count := 0
			for name, market := range d.deriveMarkets {
				if market.BaseCurrency == asset && strings.Contains(name, "20250627") {
					log.Printf("[Derive]   Found June 27, 2025 option: %s", name)
					count++
					if count >= 5 {
						break
					}
				}
			}
		}
	}
	
	// For Derive, use CCXT format if available, otherwise fall back to instrument format
	tickerSymbol := instrument
	if ccxtExists && ccxtSymbol != "" {
		tickerSymbol = ccxtSymbol
		log.Printf("[Derive] Using CCXT symbol for ticker: %s", tickerSymbol)
	} else {
		log.Printf("[Derive] Using Derive API format for ticker: %s", tickerSymbol)
	}
	
	// Try CCXT first
	ticker, err := d.exchange.FetchTicker(tickerSymbol)
	if err != nil {
		// If CCXT fails, try direct API call
		log.Printf("[Derive] CCXT FetchTicker failed (%v), trying direct API", err)
		
		// Direct API always uses Derive format, not CCXT format
		deriveTicker, apiErr := FetchDeriveTicker(instrument)
		if apiErr != nil {
			return CCXTOrderBook{}, fmt.Errorf("both CCXT and direct API failed: CCXT=%v, API=%v", err, apiErr)
		}
		
		bidPrice := deriveTicker.GetBidPrice()
		askPrice := deriveTicker.GetAskPrice()
		
		log.Printf("[Derive] Direct API ticker success: bid=%f (size=%s), ask=%f (size=%s)", 
			bidPrice, deriveTicker.BestBidAmount, askPrice, deriveTicker.BestAskAmount)
		
		// Convert to order book format
		orderBook := CCXTOrderBook{
			Bids:  [][]float64{},
			Asks:  [][]float64{},
			Index: deriveTicker.GetIndexPrice(),
		}
		
		if bidPrice > 0 {
			orderBook.Bids = append(orderBook.Bids, []float64{bidPrice, deriveTicker.GetBidSize()})
		}
		
		if askPrice > 0 {
			orderBook.Asks = append(orderBook.Asks, []float64{askPrice, deriveTicker.GetAskSize()})
		}
		
		if len(orderBook.Asks) == 0 {
			return CCXTOrderBook{}, fmt.Errorf("no ask price available from ticker")
		}
		
		return orderBook, nil
	}
	
	// Convert ticker to order book format
	orderBook := CCXTOrderBook{
		Bids: [][]float64{},
		Asks: [][]float64{},
	}
	
	// Add bid/ask from ticker if available
	if ticker.Bid != nil && *ticker.Bid > 0 {
		// Use a dummy size of 1.0 since ticker doesn't provide size
		orderBook.Bids = append(orderBook.Bids, []float64{*ticker.Bid, 1.0})
		log.Printf("[Derive] Ticker bid: %f", *ticker.Bid)
	}
	
	if ticker.Ask != nil && *ticker.Ask > 0 {
		// Use a dummy size of 1.0 since ticker doesn't provide size
		orderBook.Asks = append(orderBook.Asks, []float64{*ticker.Ask, 1.0})
		log.Printf("[Derive] Ticker ask: %f", *ticker.Ask)
	}
	
	if len(orderBook.Asks) == 0 {
		return CCXTOrderBook{}, fmt.Errorf("no ask price available from ticker")
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

// PlaceOrder places an order on Derive
func (d *CCXTDeriveExchange) PlaceOrder(conf RFQConfirmation, instrument string, cfg *AppConfig) error {
	// The instrument parameter is already in the correct format (e.g., "ETH-20250529-2550-C")
	// No need to convert again
	log.Printf("[Derive] PlaceOrder: instrument=%s", instrument)
	
	// Convert quantity from wei
	quantityFloat, err := strconv.ParseFloat(conf.Quantity, 64)
	if err != nil {
		return fmt.Errorf("failed to parse quantity: %w", err)
	}
	quantity := quantityFloat / math.Pow(10, 18)
	
	// Determine order side based on IsTakerBuy
	// IsTakerBuy represents what the taker is doing:
	// - true: taker is buying (we are selling)
	// - false: taker is selling (we are buying)
	orderSide := "buy"
	if conf.IsTakerBuy {
		orderSide = "sell"
	}
	
	log.Printf("[Derive] PlaceOrder: side=%s, quantity=%f", orderSide, quantity)
	
	// Convert to CCXT format for orders
	ccxtSymbol := ""
	if parts := strings.Split(instrument, "-"); len(parts) == 4 {
		asset := parts[0]
		dateStr := parts[1] // YYYYMMDD
		strike := parts[2]
		optType := parts[3]
		
		// Convert date format
		if len(dateStr) == 8 {
			year := dateStr[2:4]  // YY
			month := dateStr[4:6] // MM
			day := dateStr[6:8]   // DD
			ccxtSymbol = fmt.Sprintf("%s/USDC:USDC-%s-%s-%s-%s-%s", asset, year, month, day, strike, optType)
			log.Printf("[Hedge] Converted instrument %s to CCXT format: %s", instrument, ccxtSymbol)
		}
	}
	
	// Use CCXT format if available
	symbol := instrument
	if ccxtSymbol != "" {
		symbol = ccxtSymbol
	}
	
	// Try to get current price
	var bestAsk float64
	
	// Try CCXT first with CCXT symbol
	ticker, err := d.exchange.FetchTicker(symbol)
	if err != nil {
		// If CCXT fails, try direct API
		log.Printf("[Hedge] CCXT FetchTicker failed (%v), trying direct API", err)
		
		deriveTicker, apiErr := FetchDeriveTicker(instrument)
		if apiErr != nil {
			return fmt.Errorf("failed to fetch ticker from both CCXT and API: CCXT=%v, API=%v", err, apiErr)
		}
		
		bestAsk = deriveTicker.GetAskPrice()
		if bestAsk <= 0 {
			return fmt.Errorf("no ask price available from direct API")
		}
	} else {
		if ticker.Ask == nil || *ticker.Ask <= 0 {
			return fmt.Errorf("no ask price available from CCXT")
		}
		bestAsk = *ticker.Ask
	}
	
	// Calculate order price
	var hedgePrice float64
	
	// Check if price is provided in confirmation (for manual trades)
	if conf.Price != "" {
		priceFloat, err := strconv.ParseFloat(conf.Price, 64)
		if err == nil && priceFloat > 0 {
			hedgePrice = priceFloat
			log.Printf("[Order] Using manual price: %f", hedgePrice)
		}
	}
	
	// If no price provided, calculate based on order side
	if hedgePrice == 0 {
		if orderSide == "buy" {
			// For buy orders, place slightly below ask to ensure fill
			hedgePrice = bestAsk * 0.999
			
			// Also check bandwidth limits for buy orders
			// If our price exceeds ~1.8x the best ask, it might be rejected
			maxBuyPrice := bestAsk * 1.8
			if hedgePrice > maxBuyPrice {
				log.Printf("[Hedge] Buy price %f exceeds likely bandwidth, capping at %f", hedgePrice, maxBuyPrice)
				hedgePrice = maxBuyPrice
			}
			
			log.Printf("[Hedge] Buy order: best ask %f, placing at %f", bestAsk, hedgePrice)
		} else {
			// For sell orders, use defensive 2x ask strategy
			hedgePrice = bestAsk * 2.0
			log.Printf("[Hedge] Sell order: best ask %f, placing at %f (2x)", bestAsk, hedgePrice)
		}
	}
	
	// Sanity check - warn if price seems unusual
	if bestAsk < 0.01 {
		log.Printf("[Hedge] ⚠️  WARNING: Best ask price seems very low: %f", bestAsk)
	}
	if hedgePrice > 10000 {
		log.Printf("[Hedge] ⚠️  WARNING: Hedge price seems very high: %f", hedgePrice)
	}
	
	// Extract underlying from instrument (e.g., "ETH" from "ETH-20250529-2550-C")
	underlying := "ETH"
	if parts := strings.Split(instrument, "-"); len(parts) > 0 {
		underlying = parts[0]
	}
	log.Printf("[Order] Details - Symbol: %s, Side: %s, Quantity: %f %s, Price: %f USDC", 
		symbol, strings.ToUpper(orderSide), quantity, underlying, hedgePrice)
	
	// Place order
	order, err := d.exchange.CreateOrder(
		symbol,
		"limit",
		orderSide,
		quantity,
		ccxt.WithCreateOrderPrice(hedgePrice),
	)
	
	if err != nil {
		// If CCXT fails, try direct API
		log.Printf("[Hedge] CCXT CreateOrder failed (%v), trying direct API", err)
		
		// Use the original instrument format for direct API
		deriveWalletAddress := os.Getenv("DERIVE_WALLET_ADDRESS")
		subaccountIDStr := os.Getenv("DERIVE_SUBACCOUNT_ID")
		if deriveWalletAddress == "" || subaccountIDStr == "" {
			log.Printf("[Hedge] ⚠️  DERIVE_WALLET_ADDRESS or DERIVE_SUBACCOUNT_ID not set. Cannot place order.")
			log.Printf("[Hedge] ⚠️  Order details: %s %s %.4f @ %.2f USDC", instrument, strings.ToUpper(orderSide), quantity, hedgePrice)
			return nil
		}
		
		subaccountID, err := strconv.ParseUint(subaccountIDStr, 10, 64)
		if err != nil {
			log.Printf("[Hedge] ⚠️  Invalid DERIVE_SUBACCOUNT_ID: %v", err)
			return nil
		}
		
		orderResp, apiErr := d.placeDeriveOrderWithPersistentWS(instrument, orderSide, "limit", hedgePrice, quantity, subaccountID)
		if apiErr != nil {
			log.Printf("[Hedge] Direct API order failed: %v", apiErr)
			log.Printf("[Hedge] ⚠️  Order details: %s %s %.4f @ %.2f USDC", instrument, strings.ToUpper(orderSide), quantity, hedgePrice)
			return nil
		}
		
		log.Printf("[Hedge] ✅ Order placed via direct API")
		log.Printf("[Hedge] Order ID: %s", orderResp.Result.OrderID)
		log.Printf("[Hedge] Status: %s", orderResp.Result.Status)
		return nil
	}
	
	log.Printf("[Hedge] ✅ Order placed successfully on Derive")
	if order.Id != nil {
		log.Printf("[Hedge] Order ID: %s", *order.Id)
	}
	log.Printf("[Hedge] Symbol: %s", symbol)
	log.Printf("[Hedge] Side: %s", strings.ToUpper(orderSide))
	log.Printf("[Hedge] Quantity: %f", quantity)
	log.Printf("[Hedge] Price: %f USDC", hedgePrice)
	if order.Status != nil {
		log.Printf("[Hedge] Status: %s", *order.Status)
	}
	
	return nil
}

// ConvertToInstrument converts option details to Derive format
func (d *CCXTDeriveExchange) ConvertToInstrument(asset string, strike string, expiry int64, isPut bool) (string, error) {
	// Log the incoming strike for debugging
	log.Printf("[Derive] ConvertToInstrument: strike=%s, asset=%s, expiry=%d", strike, asset, expiry)
	
	// Convert strike from wei
	strikeBigInt, ok := new(big.Int).SetString(strike, 10)
	if !ok {
		return "", fmt.Errorf("invalid strike: %s", strike)
	}
	strikeNum := strikeBigInt.Div(strikeBigInt, new(big.Int).SetUint64(1e8)).String()
	
	// Format expiry - Derive uses YYYYMMDD format
	expiryTime := time.Unix(expiry, 0)
	expiryStr := expiryTime.Format("20060102") // YYYYMMDD
	
	// Option type
	optionType := "C"
	if isPut {
		optionType = "P"
		return "", fmt.Errorf("puts not supported")
	}
	
	// Build instrument - Derive format: ETH-YYYYMMDD-STRIKE-TYPE
	instrument := fmt.Sprintf("%s-%s-%s-%s", asset, expiryStr, strikeNum, optionType)
	
	log.Printf("[Derive] Converting to instrument: asset=%s, strike=%s (wei) -> %s, expiry=%d -> %s, isPut=%v -> instrument=%s",
		asset, strike, strikeNum, expiry, expiryStr, isPut, instrument)
	
	return instrument, nil
}

// initializeWebSocket creates a persistent WebSocket connection for order placement
func (d *CCXTDeriveExchange) initializeWebSocket() error {
	// Check if we have the necessary configuration
	deriveWalletAddress := os.Getenv("DERIVE_WALLET_ADDRESS")
	if deriveWalletAddress == "" {
		return fmt.Errorf("DERIVE_WALLET_ADDRESS not set")
	}
	
	if d.config.APIKey == "" {
		return fmt.Errorf("private key not configured")
	}
	
	// Create WebSocket client
	wsClient, err := NewDeriveWSClient(d.config.APIKey, deriveWalletAddress)
	if err != nil {
		return fmt.Errorf("failed to create WebSocket client: %w", err)
	}
	
	d.wsClientMu.Lock()
	d.wsClient = wsClient
	d.wsClientMu.Unlock()
	
	log.Printf("[Derive] Persistent WebSocket connection established")
	return nil
}

// getOrCreateWebSocket returns the persistent WebSocket client or creates a new one
func (d *CCXTDeriveExchange) getOrCreateWebSocket() (*DeriveWSClient, bool, error) {
	d.wsClientMu.Lock()
	defer d.wsClientMu.Unlock()
	
	// If we have a persistent connection, use it
	if d.wsClient != nil {
		return d.wsClient, false, nil // false means don't close after use
	}
	
	// Otherwise create a new one (will be closed after use)
	deriveWalletAddress := os.Getenv("DERIVE_WALLET_ADDRESS")
	if deriveWalletAddress == "" {
		return nil, false, fmt.Errorf("DERIVE_WALLET_ADDRESS not set")
	}
	
	wsClient, err := NewDeriveWSClient(d.config.APIKey, deriveWalletAddress)
	if err != nil {
		return nil, false, fmt.Errorf("failed to create WebSocket client: %w", err)
	}
	
	return wsClient, true, nil // true means close after use
}

// placeDeriveOrderWithPersistentWS places an order using the persistent WebSocket connection
func (d *CCXTDeriveExchange) placeDeriveOrderWithPersistentWS(instrumentName string, side string, orderType string, price float64, amount float64, subaccountID uint64) (*DeriveOrderResponse, error) {
	// Get WebSocket client (persistent or new)
	wsClient, shouldClose, err := d.getOrCreateWebSocket()
	if err != nil {
		return nil, fmt.Errorf("failed to get WebSocket client: %w", err)
	}
	
	// Only close if we created a new connection
	if shouldClose {
		defer wsClient.Close()
	}
	
	// Get derive wallet address
	deriveWalletAddress := os.Getenv("DERIVE_WALLET_ADDRESS")
	if deriveWalletAddress == "" {
		return nil, fmt.Errorf("DERIVE_WALLET_ADDRESS not set")
	}
	
	// Create auth for signing
	auth, err := NewDeriveAuth(d.config.APIKey)
	if err != nil {
		return nil, err
	}
	
	// Get instrument details
	instrument, err := FetchDeriveInstrumentDetails(instrumentName)
	if err != nil {
		return nil, err
	}
	
	// Create and sign action
	action := &DeriveAction{
		SubaccountID:       subaccountID,
		Owner:              deriveWalletAddress,
		Signer:             auth.address.Hex(),
		SignatureExpirySec: time.Now().Unix() + 3600,
		Nonce:              uint64(time.Now().UnixNano()),
		ModuleAddress:      "0xB8D20c2B7a1Ad2EE33Bc50eF10876eD3035b5e7b", // Trade module
		AssetAddress:       instrument.BaseAssetAddress,
		SubID:              instrument.BaseAssetSubID,
		LimitPrice:         func() *big.Int { v, _ := new(big.Float).Mul(big.NewFloat(price), big.NewFloat(1e18)).Int(nil); return v }(),
		Amount:             func() *big.Int { v, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e18)).Int(nil); return v }(),
		MaxFee:             new(big.Int).Mul(big.NewInt(1000), big.NewInt(1e18)), // 1000 USDC max fee
		RecipientID:        subaccountID,
		IsBid:              side == "buy",
	}
	
	if err := action.Sign(auth.privateKey); err != nil {
		return nil, err
	}
	
	// Prepare order request
	orderReq := map[string]interface{}{
		"instrument_name":      instrumentName,
		"direction":           side,
		"order_type":         orderType,
		"time_in_force":      "gtc",
		"mmp":                false,
		"subaccount_id":      subaccountID,
		"nonce":              action.Nonce,
		"signer":             action.Signer,
		"signature_expiry_sec": action.SignatureExpirySec,
		"signature":          action.Signature,
		"limit_price":        fmt.Sprintf("%.6f", price),
		"amount":             fmt.Sprintf("%.6f", amount),
		"max_fee":            "1000",
	}
	
	// Submit order via WebSocket
	orderResp, err := wsClient.SubmitOrder(orderReq)
	if err != nil {
		return nil, err
	}
	
	// Query open orders to verify
	openOrders, err := wsClient.GetOpenOrders(subaccountID)
	if err != nil {
		log.Printf("[Derive Order] Warning: failed to query open orders: %v", err)
	} else {
		log.Printf("[Derive Order] Open orders after submission: %d orders", len(openOrders))
	}
	
	return orderResp, nil
}

// Close closes the persistent WebSocket connection if it exists
// GetPositions fetches all open positions from Derive
func (d *CCXTDeriveExchange) GetPositions() ([]ExchangePosition, error) {
	// Get WebSocket client
	wsClient, shouldClose, err := d.getOrCreateWebSocket()
	if err != nil {
		return nil, fmt.Errorf("failed to get WebSocket client: %w", err)
	}
	
	// Only close if we created a new connection
	if shouldClose {
		defer wsClient.Close()
	}
	
	// Use the default subaccount ID from login
	subaccountID := wsClient.GetDefaultSubaccount()
	log.Printf("[Derive] Using subaccount ID: %d", subaccountID)
	
	// Fetch positions from Derive
	rawPositions, err := wsClient.GetPositions(subaccountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}
	
	// Convert raw positions to our ExchangePosition struct
	var positions []ExchangePosition
	for _, rawPos := range rawPositions {
		// Extract fields from raw position
		instrumentName, _ := rawPos["instrument_name"].(string)
		amount, _ := rawPos["amount"].(float64)
		direction, _ := rawPos["direction"].(string)
		averagePrice, _ := rawPos["average_price"].(float64)
		markPrice, _ := rawPos["mark_price"].(float64)
		indexPrice, _ := rawPos["index_price"].(float64)
		pnl, _ := rawPos["pnl"].(float64)
		
		position := ExchangePosition{
			InstrumentName: instrumentName,
			Amount:         amount,
			Direction:      direction,
			AveragePrice:   averagePrice,
			MarkPrice:      markPrice,
			IndexPrice:     indexPrice,
			PnL:            pnl,
		}
		
		positions = append(positions, position)
		
		log.Printf("[Derive] Position: %s %s %.4f @ avg %.2f (mark: %.2f, PnL: %.2f)", 
			direction, instrumentName, amount, averagePrice, markPrice, pnl)
	}
	
	log.Printf("[Derive] Total positions: %d", len(positions))
	
	return positions, nil
}

func (d *CCXTDeriveExchange) Close() error {
	d.wsClientMu.Lock()
	defer d.wsClientMu.Unlock()
	
	if d.wsClient != nil {
		d.wsClient.Close()
		d.wsClient = nil
		log.Printf("[Derive] Persistent WebSocket connection closed")
	}
	
	return nil
}