package derive

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "math/big"
    "net/http"
    "os"
    "strconv"
    "sync"
    "time"
    
    "github.com/gorilla/websocket"
    "github.com/shopspring/decimal"
    "github.com/wakamex/atomizer/internal/types"
)

// DeriveMarketMakerExchange implements types.MarketMakerExchange for Derive/Lyra
type DeriveMarketMakerExchange struct {
	wsClient     *DeriveWSClient
	subaccountID uint64
	privateKey   string // Store private key for order signing
	
	// Ticker subscriptions
	tickerConn   *websocket.Conn
	tickerMu     sync.Mutex
	subscriptions map[string]bool
	
	// Cache instrument details to avoid repeated fetches
	instrumentCache map[string]*DeriveInstrumentDetails
	cacheMu         sync.RWMutex
}

// NewDeriveMarketMakerExchange creates a new Derive exchange adapter
func NewDeriveMarketMakerExchange(privateKey, walletAddress string) (*DeriveMarketMakerExchange, error) {
	// Create WebSocket client
	wsClient, err := NewDeriveWSClient(privateKey, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create WebSocket client: %w", err)
	}
	
	// Get subaccount ID from environment or use default
	var subaccountID uint64
	if subaccountIDStr := os.Getenv("DERIVE_SUBACCOUNT_ID"); subaccountIDStr != "" {
		parsed, err := strconv.ParseUint(subaccountIDStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid DERIVE_SUBACCOUNT_ID: %w", err)
		}
		subaccountID = parsed
		log.Printf("Using subaccount ID from environment: %d", subaccountID)
	} else {
		// Fall back to default subaccount
		subaccountID = wsClient.GetDefaultSubaccount()
		log.Printf("Using default subaccount ID: %d", subaccountID)
	}
	
	return &DeriveMarketMakerExchange{
		wsClient:        wsClient,
		subaccountID:    subaccountID,
		privateKey:      privateKey,
		subscriptions:   make(map[string]bool),
		instrumentCache: make(map[string]*DeriveInstrumentDetails),
	}, nil
}

// getInstrumentDetails gets instrument details from cache or fetches if needed
func (d *DeriveMarketMakerExchange) getInstrumentDetails(instrument string) (*DeriveInstrumentDetails, error) {
	// Check cache first
	d.cacheMu.RLock()
	details, exists := d.instrumentCache[instrument]
	d.cacheMu.RUnlock()
	
	if exists {
		return details, nil
	}
	
	// Fetch if not cached
	details, err := FetchDeriveInstrumentDetails(instrument)
	if err != nil {
		return nil, err
	}
	
	// Store in cache
	d.cacheMu.Lock()
	d.instrumentCache[instrument] = details
	d.cacheMu.Unlock()
	
	return details, nil
}

// SubscribeTickers polls ticker updates for given instruments
func (d *DeriveMarketMakerExchange) SubscribeTickers(ctx context.Context, instruments []string) (<-chan types.TickerUpdate, error) {
	tickerChan := make(chan types.TickerUpdate, 100)
	
	// Store subscriptions
	for _, instrument := range instruments {
		d.subscriptions[instrument] = true
		log.Printf("Starting ticker polling for %s", instrument)
	}
	
	// Start polling for each instrument
	go d.pollTickers(ctx, instruments, tickerChan)
	
	return tickerChan, nil
}

// pollTickers polls for ticker updates
func (d *DeriveMarketMakerExchange) pollTickers(ctx context.Context, instruments []string, tickerChan chan<- types.TickerUpdate) {
	defer close(tickerChan)
	
	ticker := time.NewTicker(500 * time.Millisecond) // Poll every 500ms
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Poll each instrument
			for _, instrument := range instruments {
				go func(inst string) {
					ticker, err := d.fetchTicker(inst)
					if err != nil {
						// Only log periodically to avoid spam
						if time.Now().Unix() % 10 == 0 {
							log.Printf("Failed to fetch ticker for %s: %v", inst, err)
						}
						return
					}
					
					// Send update
					select {
					case tickerChan <- *ticker:
					default:
						log.Printf("Ticker channel full, dropping update for %s", inst)
					}
				}(instrument)
			}
		}
	}
}

// fetchTicker fetches ticker data for a single instrument
func (d *DeriveMarketMakerExchange) fetchTicker(instrument string) (*types.TickerUpdate, error) {
	url := "https://api.lyra.finance/public/get_ticker"
	payload := map[string]string{"instrument_name": instrument}
	jsonData, _ := json.Marshal(payload)
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	var result struct {
		Result struct {
			InstrumentName string `json:"instrument_name"`
			BestBidPrice   string `json:"best_bid_price"`
			BestAskPrice   string `json:"best_ask_price"`
			BestBidAmount  string `json:"best_bid_amount"`
			BestAskAmount  string `json:"best_ask_amount"`
			LastPrice      string `json:"last_price"`
			MarkPrice      string `json:"mark_price"`
			MinimumAmount  string `json:"minimum_amount"`
			AmountStep     string `json:"amount_step"`
			TickSize       string `json:"tick_size"`
			OptionPricing  *struct {
				Delta string `json:"delta"`
				Gamma string `json:"gamma"`
				Vega  string `json:"vega"`
				Theta string `json:"theta"`
			} `json:"option_pricing"`
			OptionDetails *struct {
				Expiry int64 `json:"expiry"`
				Strike string `json:"strike"`
				OptionType string `json:"option_type"`
			} `json:"option_details"`
		} `json:"result"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	if result.Error != nil {
		return nil, fmt.Errorf("API error: %s", result.Error.Message)
	}
	
	// Convert strings to decimals
	bestBid, _ := decimal.NewFromString(result.Result.BestBidPrice)
	bestAsk, _ := decimal.NewFromString(result.Result.BestAskPrice)
	bestBidAmount, _ := decimal.NewFromString(result.Result.BestBidAmount)
	bestAskAmount, _ := decimal.NewFromString(result.Result.BestAskAmount)
	lastPrice, _ := decimal.NewFromString(result.Result.LastPrice)
	markPrice, _ := decimal.NewFromString(result.Result.MarkPrice)
	
	ticker := &types.TickerUpdate{
		Instrument:  result.Result.InstrumentName,
		BestBid:     bestBid,
		BestBidSize: bestBidAmount,
		BestAsk:     bestAsk,
		BestAskSize: bestAskAmount,
		LastPrice:   lastPrice,
		MarkPrice:   markPrice,
		Timestamp:   time.Now(),
	}
	
	// Add Greeks if available (for options)
	if result.Result.OptionPricing != nil {
		delta, _ := decimal.NewFromString(result.Result.OptionPricing.Delta)
		gamma, _ := decimal.NewFromString(result.Result.OptionPricing.Gamma)
		vega, _ := decimal.NewFromString(result.Result.OptionPricing.Vega)
		theta, _ := decimal.NewFromString(result.Result.OptionPricing.Theta)
		
		ticker.Delta = &delta
		ticker.Gamma = &gamma
		ticker.Vega = &vega
		ticker.Theta = &theta
	}
	
	// Add option details if available
	if result.Result.OptionDetails != nil {
		expiryTime := time.Unix(result.Result.OptionDetails.Expiry, 0)
		ticker.Expiry = &expiryTime
		
		if result.Result.OptionDetails.Strike != "" {
			strike, _ := decimal.NewFromString(result.Result.OptionDetails.Strike)
			ticker.Strike = &strike
		}
		
		if result.Result.OptionDetails.OptionType != "" {
			ticker.OptionType = &result.Result.OptionDetails.OptionType
		}
	}
	
	return ticker, nil
}

// PlaceLimitOrder places a limit order on Derive
func (d *DeriveMarketMakerExchange) PlaceLimitOrder(instrument string, side string, price, amount decimal.Decimal) (string, error) {
	// Create order directly using our existing WebSocket connection
	// This avoids creating a new connection for each order
	
	if debugMode {
		log.Printf("DEBUG PlaceLimitOrder: Starting order placement for %s %s %.6f @ %.6f", 
			instrument, side, amount.InexactFloat64(), price.InexactFloat64())
	}
	
	// Get auth from our existing client
	auth := d.wsClient.GetAuth()
	if debugMode {
		log.Printf("DEBUG PlaceLimitOrder: Auth address: %s", auth.GetAddress())
		log.Printf("DEBUG PlaceLimitOrder: Wallet address: %s", d.wsClient.GetWallet())
		log.Printf("DEBUG PlaceLimitOrder: SubaccountID: %d", d.subaccountID)
	}
	
	// Get instrument details
	instrumentDetails, err := d.getInstrumentDetails(instrument)
	if err != nil {
		return "", fmt.Errorf("failed to fetch instrument: %w", err)
	}
	
	if debugMode {
		log.Printf("DEBUG PlaceLimitOrder: Instrument details fetched:")
		log.Printf("  BaseAssetAddress: %s", instrumentDetails.BaseAssetAddress)
		log.Printf("  BaseAssetSubID: %s", instrumentDetails.BaseAssetSubID)
	}
	
	// Calculate values for action - use exact decimal conversion to avoid precision issues
	limitPriceBigInt := func() *big.Int { 
		// Convert decimal to wei (multiply by 1e18) using exact arithmetic
		wei := price.Mul(decimal.New(1, 18)) // 1e18
		v := new(big.Int)
		v.SetString(wei.String(), 10)
		return v
	}()
	amountBigInt := func() *big.Int { 
		// Convert decimal to wei (multiply by 1e18) using exact arithmetic
		wei := amount.Mul(decimal.New(1, 18)) // 1e18
		v := new(big.Int)
		v.SetString(wei.String(), 10)
		return v
	}()
	maxFeeBigInt := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))
	
	if debugMode {
		log.Printf("DEBUG PlaceLimitOrder: Calculated values:")
		log.Printf("  LimitPrice (wei): %s", limitPriceBigInt.String())
		log.Printf("  Amount (wei): %s", amountBigInt.String())
		log.Printf("  MaxFee (wei): %s", maxFeeBigInt.String())
	}
	
	// Create signed action
	action := &DeriveAction{
		SubaccountID:       d.subaccountID,
		Owner:              d.wsClient.GetWallet(),
		Signer:             auth.GetAddress(),
		SignatureExpirySec: time.Now().Unix() + 3600,
		Nonce:              uint64(time.Now().UnixMilli())*1000 + uint64(time.Now().Nanosecond()%1000),
		ModuleAddress:      "0xB8D20c2B7a1Ad2EE33Bc50eF10876eD3035b5e7b",
		AssetAddress:       instrumentDetails.BaseAssetAddress,
		SubID:              instrumentDetails.BaseAssetSubID,
		LimitPrice:         limitPriceBigInt,
		Amount:             amountBigInt,
		MaxFee:             maxFeeBigInt,
		RecipientID:        d.subaccountID,
		IsBid:              side == "buy",
	}
	
	if debugMode {
		log.Printf("DEBUG PlaceLimitOrder: DeriveAction struct before signing:")
		log.Printf("  SubaccountID: %d", action.SubaccountID)
		log.Printf("  Owner: %s", action.Owner)
		log.Printf("  Signer: %s", action.Signer)
		log.Printf("  SignatureExpirySec: %d", action.SignatureExpirySec)
		log.Printf("  Nonce: %d", action.Nonce)
		log.Printf("  ModuleAddress: %s", action.ModuleAddress)
		log.Printf("  AssetAddress: %s", action.AssetAddress)
		log.Printf("  SubID: %s", action.SubID)
		log.Printf("  LimitPrice: %s", action.LimitPrice.String())
		log.Printf("  Amount: %s", action.Amount.String())
		log.Printf("  MaxFee: %s", action.MaxFee.String())
		log.Printf("  RecipientID: %d", action.RecipientID)
		log.Printf("  IsBid: %v", action.IsBid)
	}
	
	// Sign the action
	if err := action.Sign(auth.GetPrivateKey()); err != nil {
		return "", fmt.Errorf("failed to sign action: %w", err)
	}
	
	if debugMode {
		log.Printf("DEBUG PlaceLimitOrder: Action signed successfully")
		log.Printf("  Signature: %s", action.Signature)
	}
	
	// Prepare order request - ensure correct types match the working test
	orderReq := map[string]interface{}{
		"instrument_name":      instrument,
		"direction":           side,
		"order_type":         "limit",
		"time_in_force":      "gtc",
		"mmp":                true, // Market maker protection
		"subaccount_id":      d.subaccountID,      // int64
		"nonce":              action.Nonce,         // uint64
		"owner":              action.Owner,
		"signer":             action.Signer,
		"signature_expiry_sec": action.SignatureExpirySec, // int64
		"signature":          action.Signature,
		"limit_price":        fmt.Sprintf("%.6f", price.InexactFloat64()),
		"amount":             fmt.Sprintf("%.6f", amount.InexactFloat64()),
		"max_fee":            "100", // string of integer, not decimal
	}
	
	if debugMode {
		log.Printf("DEBUG PlaceLimitOrder: Order request payload:")
		orderReqJSON, _ := json.MarshalIndent(orderReq, "  ", "  ")
		log.Printf("%s", string(orderReqJSON))
	}
	
	// Debug the order first to see what the server expects
	if debugMode {
		debugResp, err := d.wsClient.DebugOrder(orderReq)
		if err != nil {
			if debugMode {
				log.Printf("DEBUG PlaceLimitOrder: Failed to debug order: %v", err)
			}
		} else {
			// Log the debug response
			if result, ok := debugResp["result"].(map[string]interface{}); ok {
				if debugMode {
					log.Printf("DEBUG PlaceLimitOrder: Order debug info from server:")
					log.Printf("  action_hash: %v", result["action_hash"])
					log.Printf("  typed_data_hash: %v", result["typed_data_hash"])
					log.Printf("  encoded_data_hashed: %v", result["encoded_data_hashed"])
					if rawData, ok := result["raw_data"].(map[string]interface{}); ok {
						rawJSON, _ := json.MarshalIndent(rawData, "    ", "  ")
						log.Printf("  raw_data: %s", string(rawJSON))
					}
				}
			}
		}
	}
	
	// Submit order via our existing WebSocket client
	orderResp, err := d.wsClient.SubmitOrder(orderReq)
	if err != nil {
		if debugMode {
			log.Printf("DEBUG PlaceLimitOrder: Failed to submit order: %v", err)
		}
		return "", fmt.Errorf("failed to submit order: %w", err)
	}
	
	// Debug: Log the full response
	if debugMode {
		respJSON, _ := json.MarshalIndent(orderResp, "  ", "  ")
		log.Printf("DEBUG PlaceLimitOrder: Full order response:\n%s", string(respJSON))
	}
	
	if orderResp.Error != nil {
		if debugMode {
			log.Printf("DEBUG PlaceLimitOrder: Order error response: %s", orderResp.Error.Message)
		}
		return "", fmt.Errorf("order error: %s", orderResp.Error.Message)
	}
	
	if orderResp.Result.OrderID == "" {
		if debugMode {
			log.Printf("DEBUG PlaceLimitOrder: Order response has empty order ID")
		}
		return "", fmt.Errorf("order response has empty order ID")
	}
	
	if debugMode {
		log.Printf("DEBUG PlaceLimitOrder: Order placed successfully with ID: %s", orderResp.Result.OrderID)
	}
	return orderResp.Result.OrderID, nil
}

// ReplaceOrder replaces an existing order with new parameters
func (d *DeriveMarketMakerExchange) ReplaceOrder(orderID string, instrument string, side string, price, amount decimal.Decimal) (string, error) {
	if debugMode {
		log.Printf("DEBUG ReplaceOrder: Starting order replacement for order %s -> %s %s %.6f @ %.6f", 
			orderID, instrument, side, amount.InexactFloat64(), price.InexactFloat64())
	}
	
	// Get auth from our existing client
	auth := d.wsClient.GetAuth()
	if debugMode {
		log.Printf("DEBUG ReplaceOrder: Auth address: %s", auth.GetAddress())
	}
	
	// Get instrument details
	instrumentDetails, err := d.getInstrumentDetails(instrument)
	if err != nil {
		return "", fmt.Errorf("failed to fetch instrument: %w", err)
	}
	
	if debugMode {
		log.Printf("DEBUG ReplaceOrder: Instrument details:")
		log.Printf("  BaseAssetAddress: %s", instrumentDetails.BaseAssetAddress)
		log.Printf("  BaseAssetSubID: %s", instrumentDetails.BaseAssetSubID)
	}
	
	// Calculate values for action
	limitPriceBigInt := func() *big.Int { 
		v, _ := new(big.Float).Mul(big.NewFloat(price.InexactFloat64()), big.NewFloat(1e18)).Int(nil)
		return v
	}()
	amountBigInt := func() *big.Int { 
		v, _ := new(big.Float).Mul(big.NewFloat(amount.InexactFloat64()), big.NewFloat(1e18)).Int(nil)
		return v
	}()
	maxFeeBigInt := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18))
	
	// Create signed action for the new order
	action := &DeriveAction{
		SubaccountID:       d.subaccountID,
		Owner:              d.wsClient.GetWallet(),
		Signer:             auth.GetAddress(),
		SignatureExpirySec: time.Now().Unix() + 3600,
		Nonce:              uint64(time.Now().UnixMilli())*1000 + uint64(time.Now().Nanosecond()%1000),
		ModuleAddress:      "0xB8D20c2B7a1Ad2EE33Bc50eF10876eD3035b5e7b",
		AssetAddress:       instrumentDetails.BaseAssetAddress,
		SubID:              instrumentDetails.BaseAssetSubID,
		LimitPrice:         limitPriceBigInt,
		Amount:             amountBigInt,
		MaxFee:             maxFeeBigInt,
		RecipientID:        d.subaccountID,
		IsBid:              side == "buy",
	}
	
	if debugMode {
		log.Printf("DEBUG ReplaceOrder: DeriveAction struct before signing:")
		log.Printf("  SubaccountID: %d", action.SubaccountID)
		log.Printf("  Owner: %s", action.Owner)
		log.Printf("  Signer: %s", action.Signer)
		log.Printf("  SignatureExpirySec: %d", action.SignatureExpirySec)
		log.Printf("  Nonce: %d", action.Nonce)
		log.Printf("  ModuleAddress: %s", action.ModuleAddress)
		log.Printf("  AssetAddress: %s", action.AssetAddress)
		log.Printf("  SubID: %s", action.SubID)
		log.Printf("  LimitPrice: %s", action.LimitPrice.String())
		log.Printf("  Amount: %s", action.Amount.String())
		log.Printf("  MaxFee: %s", action.MaxFee.String())
		log.Printf("  RecipientID: %d", action.RecipientID)
		log.Printf("  IsBid: %v", action.IsBid)
	}
	
	// Sign the action
	if err := action.Sign(auth.GetPrivateKey()); err != nil {
		return "", fmt.Errorf("failed to sign action: %w", err)
	}
	
	if debugMode {
		log.Printf("DEBUG ReplaceOrder: Action signed successfully, signature: %s", action.Signature)
	}
	
	// Prepare replace request
	replaceReq := map[string]interface{}{
		// Cancel parameters - use the correct field name
		"order_id_to_cancel": orderID,
		
		// New order parameters (all required fields)
		"instrument_name":      instrument,
		"direction":           side,
		"order_type":         "limit",
		"time_in_force":      "gtc",
		"amount":             fmt.Sprintf("%.6f", amount.InexactFloat64()),
		"limit_price":        fmt.Sprintf("%.6f", price.InexactFloat64()),
		"max_fee":            "100",
		"subaccount_id":      d.subaccountID,
		"nonce":              action.Nonce,
		"signature_expiry_sec": action.SignatureExpirySec,
		"owner":              action.Owner,
		"signer":             action.Signer,
		"signature":          action.Signature,
		"mmp":                true, // Market maker protection
	}
	
	// Send replace request
	id := fmt.Sprintf("%d", time.Now().UnixMilli())
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "private/replace",
		"params":  replaceReq,
		"id":      id,
	}
	
	log.Printf("Sending replace order request for order %s", orderID)
	respChan := d.wsClient.sendRequest(req)
	
	select {
	case resp := <-respChan:
		var result struct {
			Result *struct {
				Order *struct {
					OrderID string `json:"order_id"`
				} `json:"order"`
				CancelledOrder *struct {
					OrderID string `json:"order_id"`
				} `json:"cancelled_order"`
				CreateOrderError *struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
					Data    string `json:"data"`
				} `json:"create_order_error"`
			} `json:"result"`
			Error *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(resp, &result); err != nil {
			return "", fmt.Errorf("failed to parse response: %w", err)
		}
		if result.Error != nil {
			return "", fmt.Errorf("replace error: %s", result.Error.Message)
		}
		if result.Result != nil {
			// Check for partial failure: cancelled but not created
			if result.Result.CancelledOrder != nil && result.Result.Order == nil {
				// Log the partial failure
				log.Printf("WARNING: Replace order partially failed - cancelled order %s but failed to create new order", 
					result.Result.CancelledOrder.OrderID)
				
				if result.Result.CreateOrderError != nil {
					log.Printf("Create order error: Code=%d, Message=%s, Data=%s", 
						result.Result.CreateOrderError.Code, 
						result.Result.CreateOrderError.Message,
						result.Result.CreateOrderError.Data)
				}
				
				// Attempt recovery by placing a new order
				log.Printf("Attempting to recover by placing a new order: %s %s %.6f @ %.6f", 
					instrument, side, amount.InexactFloat64(), price.InexactFloat64())
				
				return d.PlaceLimitOrder(instrument, side, price, amount)
			}
			
			// Success case: both cancelled and created
			if result.Result.Order != nil {
				if result.Result.CancelledOrder != nil {
					log.Printf("Successfully replaced order %s with new order %s", 
						result.Result.CancelledOrder.OrderID, result.Result.Order.OrderID)
				}
				return result.Result.Order.OrderID, nil
			}
		}
		return "", fmt.Errorf("no order ID in response")
	case <-time.After(5 * time.Second):
		return "", fmt.Errorf("replace order timeout")
	}
}

// CancelOrder cancels an order on Derive
func (d *DeriveMarketMakerExchange) CancelOrder(orderID string) error {
	// First, we need to find the instrument name for this order
	orders, err := d.GetOpenOrders()
	if err != nil {
		return fmt.Errorf("failed to get open orders: %w", err)
	}
	
	var instrumentName string
	for _, order := range orders {
		if order.OrderID == orderID {
			instrumentName = order.Instrument
			break
		}
	}
	
	if instrumentName == "" {
		// Order not found in open orders, might already be cancelled
		return nil
	}
	
	id := fmt.Sprintf("%d", time.Now().UnixMilli())
	
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "private/cancel",
		"params": map[string]interface{}{
			"order_id": orderID,
			"subaccount_id": d.subaccountID,
			"instrument_name": instrumentName,
		},
		"id": id,
	}
	
	respChan := d.wsClient.sendRequest(req)
	
	select {
	case resp := <-respChan:
		var result struct {
			Error *struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		
		if err := json.Unmarshal(resp, &result); err != nil {
			return fmt.Errorf("failed to parse cancel response: %w", err)
		}
		
		if result.Error != nil {
			return fmt.Errorf("cancel error: %s", result.Error.Message)
		}
		
		return nil
		
	case <-time.After(10 * time.Second):
		return fmt.Errorf("cancel order timeout")
	}
}

// GetOpenOrders gets all open orders
func (d *DeriveMarketMakerExchange) GetOpenOrders() ([]types.MarketMakerOrder, error) {
	rawOrders, err := d.wsClient.GetOpenOrders(d.subaccountID)
	if err != nil {
		return nil, err
	}
	
	orders := make([]types.MarketMakerOrder, 0, len(rawOrders))
	for _, raw := range rawOrders {
		order := types.MarketMakerOrder{
			OrderID:    getString(raw, "order_id"),
			Instrument: getString(raw, "instrument_name"),
			Side:       getString(raw, "direction"),
			Price:      getDecimal(raw, "limit_price"),
			Amount:     getDecimal(raw, "amount"),
			FilledAmount: getDecimal(raw, "filled_amount"),
			Status:     getString(raw, "status"),
			CreatedAt:  time.Unix(getInt64(raw, "created_at")/1000, 0),
			UpdatedAt:  time.Unix(getInt64(raw, "updated_at")/1000, 0),
		}
		orders = append(orders, order)
	}
	
	return orders, nil
}

// GetPositions gets current positions
func (d *DeriveMarketMakerExchange) GetPositions() ([]types.ExchangePosition, error) {
	rawPositions, err := d.wsClient.GetPositions(d.subaccountID)
	if err != nil {
		return nil, err
	}
	
	positions := make([]types.ExchangePosition, 0, len(rawPositions))
	for i, raw := range rawPositions {
		// Debug log first position to see field names (only if debug mode is enabled)
		if i == 0 && debugMode {
			log.Printf("DEBUG: First position raw data:")
			for key, value := range raw {
				log.Printf("  %s: %v (type: %T)", key, value, value)
			}
		}
		
		// Try both "amount" and "size" fields
		amount := getFloat64(raw, "amount")
		if amount == 0 {
			amount = getFloat64(raw, "size")
		}
		
		position := types.ExchangePosition{
			InstrumentName: getString(raw, "instrument_name"),
			Amount:         amount,
			Direction:      getString(raw, "direction"),
			AveragePrice:   getFloat64(raw, "average_price"),
			MarkPrice:      getFloat64(raw, "mark_price"),
			IndexPrice:     getFloat64(raw, "index_price"),
			PnL:            getFloat64(raw, "pnl"),
		}
		positions = append(positions, position)
	}
	
	return positions, nil
}

// GetOrderBook returns the cached orderbook from WebSocket subscription
func (d *DeriveMarketMakerExchange) GetOrderBook(instrument string) (*types.MarketMakerOrderBook, error) {
	// Get cached orderbook from WebSocket client
	orderbook := d.wsClient.GetOrderBook(instrument)
	if orderbook == nil {
		return nil, fmt.Errorf("no orderbook data available for %s", instrument)
	}
	
	// Convert to types.MarketMakerOrderBook format
	mmOrderBook := &types.MarketMakerOrderBook{
		Bids:      orderbook.Bids,
		Asks:      orderbook.Asks,
		Timestamp: orderbook.Timestamp,
	}
	
	return mmOrderBook, nil
}

// SubscribeOrderBook subscribes to orderbook updates for an instrument
func (d *DeriveMarketMakerExchange) SubscribeOrderBook(instrument string) error {
	return d.wsClient.SubscribeOrderBook(instrument, 20) // Subscribe with depth 20
}

// Helper functions to extract values from map
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getDecimal(m map[string]interface{}, key string) decimal.Decimal {
	switch v := m[key].(type) {
	case string:
		d, _ := decimal.NewFromString(v)
		return d
	case float64:
		return decimal.NewFromFloat(v)
	default:
		return decimal.Zero
	}
}

func getFloat64(m map[string]interface{}, key string) float64 {
	switch v := m[key].(type) {
	case float64:
		return v
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	default:
		return 0
	}
}

func getInt64(m map[string]interface{}, key string) int64 {
	if v, ok := m[key].(float64); ok {
		return int64(v)
	}
	return 0
}

// Close closes the exchange connections
func (d *DeriveMarketMakerExchange) Close() error {
	d.tickerMu.Lock()
	if d.tickerConn != nil {
		d.tickerConn.Close()
	}
	d.tickerMu.Unlock()
	
	if d.wsClient != nil {
		return d.wsClient.Close()
	}
	
	return nil
}
