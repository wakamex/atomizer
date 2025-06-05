package derive

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
	"github.com/wakamex/atomizer/internal/exchange/shared"
	"github.com/wakamex/atomizer/internal/types"
)

// DeriveWSClient handles WebSocket connection to Derive
type DeriveWSClient struct {
	conn        *websocket.Conn
	auth        *DeriveAuth
	wallet      string
	subaccounts []int
	mu          sync.Mutex
	writeMu     sync.Mutex // Protects WebSocket writes
	requests    map[string]chan json.RawMessage

	// Orderbook data
	orderbooks    map[string]*OrderBookData
	orderbookMu   sync.RWMutex
	orderbookSubs map[string]bool

	// Connection state management
	isConnected       bool
	lastActivity      time.Time
	reconnectChan     chan struct{}
	shutdownChan      chan struct{}
	wsURL             string
	reconnectDelay    time.Duration
	maxReconnectDelay time.Duration
	pingTicker        *time.Ticker
	heartbeatChan     chan struct{}
}

// OrderBookData represents orderbook state
type OrderBookData struct {
	Bids      []types.OrderBookLevel
	Asks      []types.OrderBookLevel
	Timestamp time.Time
	ChangeID  int64
}

// NewDeriveWSClient creates a new Derive WebSocket client
func NewDeriveWSClient(privateKey string, deriveWallet string) (*DeriveWSClient, error) {
	auth, err := NewDeriveAuth(privateKey)
	if err != nil {
		return nil, err
	}

	client := &DeriveWSClient{
		auth:              auth,
		wallet:            deriveWallet,
		requests:          make(map[string]chan json.RawMessage),
		orderbooks:        make(map[string]*OrderBookData),
		orderbookSubs:     make(map[string]bool),
		wsURL:             "wss://api.lyra.finance/ws",
		reconnectDelay:    1 * time.Second,
		maxReconnectDelay: 30 * time.Second,
		reconnectChan:     make(chan struct{}, 1),
		shutdownChan:      make(chan struct{}),
		heartbeatChan:     make(chan struct{}, 1),
	}

	// Establish initial connection
	if err := client.connect(); err != nil {
		return nil, fmt.Errorf("failed to establish initial connection: %w", err)
	}

	// Start connection monitor
	go client.connectionMonitor()

	// Start heartbeat to keep connection alive
	go client.heartbeat()

	return client, nil
}

// connect establishes WebSocket connection and performs login
func (c *DeriveWSClient) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close existing connection if any
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	log.Printf("[Derive WS] Connecting to %s", c.wsURL)
	shared.DeriveDebugLog("[Derive WS] Connecting to %s", c.wsURL)

	// Set connection timeout
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second

	conn, _, err := dialer.Dial(c.wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to Derive WebSocket: %w", err)
	}

	c.conn = conn
	c.lastActivity = time.Now()

	// Set up ping/pong handlers
	conn.SetPingHandler(func(appData string) error {
		shared.DeriveDebugLog("[Derive WS] Ping received, sending Pong")
		c.updateActivity()
		err := conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(5*time.Second))
		if err != nil {
			shared.DeriveDebugLog("[Derive WS] Error sending pong: %v", err)
		}
		return nil
	})

	conn.SetPongHandler(func(appData string) error {
		shared.DeriveDebugLog("[Derive WS] Pong received")
		c.updateActivity()
		return nil
	})

	// Set read deadline for initial messages
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// Start message handler
	go c.handleMessages()

	// Unlock mutex temporarily for login (which needs to send/receive messages)
	c.mu.Unlock()

	// Login
	if err := c.login(); err != nil {
		c.mu.Lock()
		conn.Close()
		c.conn = nil
		c.isConnected = false
		return fmt.Errorf("failed to login: %w", err)
	}

	c.mu.Lock()
	c.isConnected = true

	// Clear read deadline after successful login
	c.conn.SetReadDeadline(time.Time{})

	// Resubscribe to orderbooks
	c.resubscribeOrderbooks()

	log.Printf("[Derive WS] Successfully connected and authenticated")

	return nil
}

// updateActivity updates the last activity timestamp
func (c *DeriveWSClient) updateActivity() {
	c.mu.Lock()
	c.lastActivity = time.Now()
	c.mu.Unlock()
}

// connectionMonitor monitors connection health and triggers reconnection when needed
func (c *DeriveWSClient) connectionMonitor() {
	for {
		select {
		case <-c.reconnectChan:
			log.Printf("[Derive WS] Reconnection requested")
			c.performReconnection()

		case <-c.shutdownChan:
			log.Printf("[Derive WS] Connection monitor shutting down")
			return
		}
	}
}

// performReconnection handles the reconnection logic with exponential backoff
func (c *DeriveWSClient) performReconnection() {
	c.mu.Lock()
	c.isConnected = false
	delay := c.reconnectDelay
	c.mu.Unlock()

	for {
		select {
		case <-c.shutdownChan:
			return
		default:
		}

		log.Printf("[Derive WS] Attempting reconnection in %v", delay)
		time.Sleep(delay)

		if err := c.connect(); err != nil {
			log.Printf("[Derive WS] Reconnection failed: %v", err)

			// Exponential backoff
			delay = delay * 2
			if delay > c.maxReconnectDelay {
				delay = c.maxReconnectDelay
			}
		} else {
			// Reset delay on successful reconnection
			c.mu.Lock()
			c.reconnectDelay = 1 * time.Second
			c.mu.Unlock()
			return
		}
	}
}

// triggerReconnection triggers a reconnection attempt
func (c *DeriveWSClient) triggerReconnection() {
	select {
	case c.reconnectChan <- struct{}{}:
		shared.DeriveDebugLog("[Derive WS] Reconnection triggered")
	default:
		// Channel is full, reconnection already scheduled
	}
}

// resubscribeOrderbooks re-subscribes to all previously subscribed orderbooks
func (c *DeriveWSClient) resubscribeOrderbooks() {
	c.orderbookMu.RLock()
	channels := make([]string, 0, len(c.orderbookSubs))
	for channel := range c.orderbookSubs {
		channels = append(channels, channel)
	}
	c.orderbookMu.RUnlock()

	if len(channels) > 0 {
		msg := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "subscribe",
			"params": map[string]interface{}{
				"channels": channels,
			},
			"id": fmt.Sprintf("resubscribe_%d", time.Now().UnixNano()),
		}

		c.writeMu.Lock()
		err := c.conn.WriteJSON(msg)
		c.writeMu.Unlock()

		if err != nil {
			log.Printf("[Derive WS] Failed to resubscribe to orderbooks: %v", err)
		} else {
			log.Printf("[Derive WS] Resubscribed to %d orderbook channels", len(channels))
		}
	}
}

// heartbeat sends periodic pings to keep the connection alive
func (c *DeriveWSClient) heartbeat() {
	c.pingTicker = time.NewTicker(15 * time.Second)
	defer c.pingTicker.Stop()

	activityCheckTicker := time.NewTicker(45 * time.Second)
	defer activityCheckTicker.Stop()

	for {
		select {
		case <-c.shutdownChan:
			return

		case <-c.pingTicker.C:
			c.mu.Lock()
			if c.conn == nil || !c.isConnected {
				c.mu.Unlock()
				continue
			}

			// Get connection reference while holding lock
			conn := c.conn
			c.mu.Unlock()

			// Send ping with write mutex
			c.writeMu.Lock()
			err := conn.WriteMessage(websocket.PingMessage, nil)
			c.writeMu.Unlock()

			if err != nil {
				log.Printf("[Derive WS] Error sending ping: %v", err)
				c.mu.Lock()
				c.isConnected = false
				c.mu.Unlock()
				c.triggerReconnection()
				continue
			}
			shared.DeriveDebugLog("[Derive WS] Ping sent")

		case <-activityCheckTicker.C:
			// Check for connection health based on last activity
			c.mu.Lock()
			if c.isConnected && time.Since(c.lastActivity) > 90*time.Second {
				log.Printf("[Derive WS] No activity for %v, connection may be dead", time.Since(c.lastActivity))
				c.isConnected = false
				c.mu.Unlock()
				c.triggerReconnection()
			} else {
				c.mu.Unlock()
			}

		case <-c.heartbeatChan:
			// External trigger to check connection
			c.mu.Lock()
			if !c.isConnected {
				c.mu.Unlock()
				c.triggerReconnection()
			} else {
				c.mu.Unlock()
			}
		}
	}
}

// login authenticates the WebSocket session
func (c *DeriveWSClient) login() error {
	timestamp := time.Now().UTC().UnixMilli()
	timestampStr := fmt.Sprintf("%d", timestamp)

	signature, err := c.auth.SignMessage(timestampStr)
	if err != nil {
		return fmt.Errorf("failed to sign login message: %w", err)
	}

	ownerEOA := c.auth.GetAddress()

	shared.DeriveDebugLog("[Derive WS] Login - EOA: %s", ownerEOA)
	shared.DeriveDebugLog("[Derive WS] Login - Derive Wallet: %s", c.wallet)
	shared.DeriveDebugLog("[Derive WS] Login - Timestamp: %d (%s)", timestamp, timestampStr)
	shared.DeriveDebugLog("[Derive WS] Login - Signature: %s", signature)

	// Use JSON-RPC format - server expects it
	loginReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "public/login",
		"params": map[string]interface{}{
			"wallet":    c.wallet,     // Keep original checksummed format
			"timestamp": timestampStr, // String format
			"signature": signature,
		},
		"id": fmt.Sprintf("%d", time.Now().UnixMilli()),
	}

	shared.DeriveDebugLog("[Derive WS] Sending login request: %+v", loginReq)

	// Send login request
	respChan := c.sendRequest(loginReq)

	select {
	case resp := <-respChan:
		shared.DeriveDebugLog("[Derive WS] Login response: %s", string(resp))

		var result struct {
			Result []int `json:"result"` // Array of subaccount IDs
			Error  *struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Data    interface{} `json:"data"`
			} `json:"error"`
		}
		if err := json.Unmarshal(resp, &result); err != nil {
			return fmt.Errorf("failed to parse login response: %w", err)
		}
		if result.Error != nil {
			// Always log full error details for debugging
			log.Printf("[Derive WS] Login error - Code: %d, Message: %s, Data: %v",
				result.Error.Code, result.Error.Message, result.Error.Data)
			return fmt.Errorf("login error: %s", result.Error.Message)
		}

		// Store subaccount IDs
		c.subaccounts = result.Result

		shared.DeriveDebugLog("[Derive WS] Login successful. Subaccounts: %v", result.Result)
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("login timeout")
	}
}

// handleMessages processes incoming WebSocket messages
func (c *DeriveWSClient) handleMessages() {
	defer func() {
		// Trigger reconnection when message handler exits
		c.mu.Lock()
		c.isConnected = false
		c.mu.Unlock()
		c.triggerReconnection()
	}()

	for {
		c.mu.Lock()
		if c.conn == nil {
			c.mu.Unlock()
			return
		}
		conn := c.conn
		c.mu.Unlock()

		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[Derive WS] Unexpected close error: %v", err)
			} else {
				shared.DeriveDebugLog("[Derive WS] Read error: %v", err)
			}
			return
		}

		// Update activity timestamp
		c.updateActivity()

		// Log raw message for debugging
		if len(message) == 0 {
			log.Printf("[Derive WS] WARNING: Received empty message from WebSocket")
			// Empty messages might indicate connection issues
			c.mu.Lock()
			if time.Since(c.lastActivity) > 60*time.Second {
				c.isConnected = false
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()
			continue
		} else if len(message) < 500 {
			shared.DeriveWSDebugLog("[Derive WS] Raw message: %s", string(message))
		} else {
			shared.DeriveWSDebugLog("[Derive WS] Raw message (truncated): %s...", string(message[:500]))
		}

		// First try to parse as a subscription update
		var subMsg struct {
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
		}
		if err := json.Unmarshal(message, &subMsg); err == nil && subMsg.Method == "subscription" {
			// Handle subscription update
			c.handleSubscriptionUpdate(subMsg.Params)
			continue
		}

		// Otherwise handle as a response
		var msg struct {
			ID     string          `json:"id"`
			Result json.RawMessage `json:"result"`
			Error  json.RawMessage `json:"error"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			shared.DeriveDebugLog("[Derive WS] Failed to parse message: %v", err)
			continue
		}

		// Handle response
		c.mu.Lock()
		if ch, ok := c.requests[msg.ID]; ok {
			// Log the raw message being sent to the channel
			if len(message) == 0 {
				log.Printf("[Derive WS] WARNING: Sending empty message to request %s", msg.ID)
			} else if len(message) < 100 {
				shared.DeriveDebugLog("[Derive WS] Sending response to request %s: %s", msg.ID, string(message))
			} else {
				shared.DeriveDebugLog("[Derive WS] Sending response to request %s (length: %d)", msg.ID, len(message))
			}
			ch <- message
			delete(c.requests, msg.ID)
		}
		c.mu.Unlock()
	}
}

// handleSubscriptionUpdate processes subscription updates (orderbook, trades, etc)
func (c *DeriveWSClient) handleSubscriptionUpdate(params json.RawMessage) {
	var update struct {
		Channel string `json:"channel"`
		Data    struct {
			Timestamp      int64           `json:"timestamp"`
			InstrumentName string          `json:"instrument_name"`
			Bids           [][]json.Number `json:"bids"` // [price, size]
			Asks           [][]json.Number `json:"asks"` // [price, size]
			ChangeID       int64           `json:"change_id"`
		} `json:"data"`
	}

	if err := json.Unmarshal(params, &update); err != nil {
		shared.DeriveDebugLog("[Derive WS] Failed to parse subscription update: %v", err)
		return
	}

	// Check if this is an orderbook update
	if strings.HasPrefix(update.Channel, "orderbook.") {
		c.updateOrderBook(update.Data.InstrumentName, &update.Data)
	}
}

// updateOrderBook updates the cached orderbook for an instrument
func (c *DeriveWSClient) updateOrderBook(instrument string, data *struct {
	Timestamp      int64           `json:"timestamp"`
	InstrumentName string          `json:"instrument_name"`
	Bids           [][]json.Number `json:"bids"`
	Asks           [][]json.Number `json:"asks"`
	ChangeID       int64           `json:"change_id"`
}) {
	// Convert bids and asks to types.OrderBookLevel
	bids := make([]types.OrderBookLevel, 0, len(data.Bids))
	for _, bid := range data.Bids {
		if len(bid) >= 2 {
			price, _ := decimal.NewFromString(bid[0].String())
			amount, _ := decimal.NewFromString(bid[1].String())
			bids = append(bids, types.OrderBookLevel{
				Price: price,
				Size:  amount,
			})
		}
	}

	asks := make([]types.OrderBookLevel, 0, len(data.Asks))
	for _, ask := range data.Asks {
		if len(ask) >= 2 {
			price, _ := decimal.NewFromString(ask[0].String())
			amount, _ := decimal.NewFromString(ask[1].String())
			asks = append(asks, types.OrderBookLevel{
				Price: price,
				Size:  amount,
			})
		}
	}

	// Update cached orderbook
	c.orderbookMu.Lock()
	c.orderbooks[instrument] = &OrderBookData{
		Bids:      bids,
		Asks:      asks,
		Timestamp: time.Unix(0, data.Timestamp*1e6), // Convert microseconds to nanoseconds
		ChangeID:  data.ChangeID,
	}
	c.orderbookMu.Unlock()

	shared.DeriveWSDebugLog("[Derive WS] Updated orderbook for %s: %d bids, %d asks", instrument, len(bids), len(asks))
}

// sendRequest sends a request and returns a channel for the response
func (c *DeriveWSClient) sendRequest(req map[string]interface{}) <-chan json.RawMessage {
	respChan := make(chan json.RawMessage, 1)

	id, _ := req["id"].(string)

	c.mu.Lock()
	// Allow login requests to go through even if not yet connected
	method, _ := req["method"].(string)
	if method != "public/login" && (!c.isConnected || c.conn == nil) {
		c.mu.Unlock()
		log.Printf("[Derive WS] Cannot send request: connection not available")
		close(respChan)
		// Trigger reconnection
		c.triggerReconnection()
		return respChan
	}

	// For login requests, just check if conn is available
	if c.conn == nil {
		c.mu.Unlock()
		log.Printf("[Derive WS] Cannot send request: no connection")
		close(respChan)
		return respChan
	}

	c.requests[id] = respChan
	conn := c.conn
	c.mu.Unlock()

	// Marshal to JSON to log exact format
	jsonBytes, _ := json.Marshal(req)
	shared.DeriveDebugLog("[Derive WS] Sending JSON: %s", string(jsonBytes))

	// Protect WebSocket write with mutex
	c.writeMu.Lock()
	err := conn.WriteJSON(req)
	c.writeMu.Unlock()

	if err != nil {
		log.Printf("[Derive WS] Failed to send request: %v", err)
		c.mu.Lock()
		delete(c.requests, id)
		c.isConnected = false
		c.mu.Unlock()
		close(respChan)
		c.triggerReconnection()
		return respChan
	}

	return respChan
}

// DebugOrder sends order to order_debug endpoint to verify signature
func (c *DeriveWSClient) DebugOrder(order map[string]interface{}) (map[string]interface{}, error) {
	id := fmt.Sprintf("%d", time.Now().UnixMilli())

	// Use JSON-RPC format
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "private/order_debug",
		"params":  order,
		"id":      id,
	}

	shared.DeriveDebugLog("[Derive WS] Sending order_debug request")

	respChan := c.sendRequest(req)

	select {
	case resp := <-respChan:
		shared.DeriveDebugLog("[Derive WS] Order debug response: %s", string(resp))

		var result map[string]interface{}
		if err := json.Unmarshal(resp, &result); err != nil {
			return nil, fmt.Errorf("failed to parse debug response: %w", err)
		}

		return result, nil

	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("order debug timeout")
	}
}

// SubmitOrder submits an order via WebSocket
func (c *DeriveWSClient) SubmitOrder(order map[string]interface{}) (*DeriveOrderResponse, error) {
	id := fmt.Sprintf("%d", time.Now().UnixMilli())

	// Use JSON-RPC format
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "private/order",
		"params":  order,
		"id":      id,
	}

	shared.DeriveDebugLog("[Derive WS] Submitting order: %+v", order)

	respChan := c.sendRequest(req)

	select {
	case resp := <-respChan:
		shared.DeriveDebugLog("[Derive WS] Order response: %s", string(resp))

		// First check if there's an error
		var errorCheck struct {
			Error *struct {
				Code    int         `json:"code"`
				Message string      `json:"message"`
				Data    interface{} `json:"data"`
			} `json:"error"`
		}
		if err := json.Unmarshal(resp, &errorCheck); err == nil && errorCheck.Error != nil {
			// Log the full error for debugging
			shared.DeriveDebugLog("[Derive WS] Order error detected: code=%d, message=%s, data=%v",
				errorCheck.Error.Code, errorCheck.Error.Message, errorCheck.Error.Data)

			// If we get price band error, include bandwidth info
			if errorCheck.Error.Code == 11013 {
				// Try to parse data as bandwidth struct
				if dataMap, ok := errorCheck.Error.Data.(map[string]interface{}); ok {
					limit, _ := dataMap["limit"].(string)
					bandwidth, _ := dataMap["bandwidth"].(string)
					return nil, fmt.Errorf("order error: %s (limit: %s, bandwidth: %s)",
						errorCheck.Error.Message, limit, bandwidth)
				}
			}
			// For signature errors, include the data field if it's a string
			if errorCheck.Error.Code == 14014 {
				if dataStr, ok := errorCheck.Error.Data.(string); ok && dataStr != "" {
					return nil, fmt.Errorf("order error: %s - %s", errorCheck.Error.Message, dataStr)
				}
			}
			return nil, fmt.Errorf("order error: %s (code: %d)", errorCheck.Error.Message, errorCheck.Error.Code)
		}

		// Try to parse the successful response
		var orderResp struct {
			Result struct {
				Order map[string]interface{} `json:"order"`
			} `json:"result"`
		}

		if err := json.Unmarshal(resp, &orderResp); err != nil {
			return nil, fmt.Errorf("failed to parse order response: %w", err)
		}

		// Convert to DeriveOrderResponse
		result := &DeriveOrderResponse{}
		if orderID, ok := orderResp.Result.Order["order_id"].(string); ok {
			result.Result.OrderID = orderID
		}
		if status, ok := orderResp.Result.Order["status"].(string); ok {
			result.Result.Status = status
		}
		if instrumentName, ok := orderResp.Result.Order["instrument_name"].(string); ok {
			result.Result.InstrumentName = instrumentName
		}

		shared.DeriveDebugLog("[Derive WS] Order placed - ID: %s, Status: %s", result.Result.OrderID, result.Result.Status)

		return result, nil

	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("order submission timeout")
	}
}

// GetOpenOrders queries open orders for a subaccount
func (c *DeriveWSClient) GetOpenOrders(subaccountID uint64) ([]map[string]interface{}, error) {
	id := fmt.Sprintf("%d", time.Now().UnixMilli())

	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "private/get_open_orders",
		"params": map[string]interface{}{
			"subaccount_id": subaccountID,
		},
		"id": id,
	}

	shared.DeriveDebugLog("[Derive WS] Querying open orders for subaccount %d", subaccountID)

	respChan := c.sendRequest(req)

	select {
	case resp := <-respChan:
		var result struct {
			Result struct {
				SubaccountID int                      `json:"subaccount_id"`
				Orders       []map[string]interface{} `json:"orders"`
			} `json:"result"`
			Error *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.Unmarshal(resp, &result); err != nil {
			return nil, fmt.Errorf("failed to parse orders response: %w", err)
		}

		if result.Error != nil {
			return nil, fmt.Errorf("get orders error: %s", result.Error.Message)
		}

		return result.Result.Orders, nil

	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("get orders timeout")
	}
}

// GetDefaultSubaccount returns the first subaccount ID
func (c *DeriveWSClient) GetDefaultSubaccount() uint64 {
	if len(c.subaccounts) > 0 {
		return uint64(c.subaccounts[0])
	}
	return 0
}

// GetAuth returns the auth instance for signing
func (c *DeriveWSClient) GetAuth() *DeriveAuth {
	return c.auth
}

// GetWallet returns the wallet address
func (c *DeriveWSClient) GetWallet() string {
	return c.wallet
}

// GetPositions fetches all positions for the subaccount
func (c *DeriveWSClient) GetPositions(subaccountID uint64) ([]map[string]interface{}, error) {
	// Check if WebSocket is still connected
	c.mu.Lock()
	if !c.isConnected || c.conn == nil {
		c.mu.Unlock()
		return nil, fmt.Errorf("WebSocket connection not available")
	}
	c.mu.Unlock()

	id := fmt.Sprintf("%d", time.Now().UnixMilli())

	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "private/get_positions",
		"params": map[string]interface{}{
			"subaccount_id": subaccountID,
		},
		"id": id,
	}

	log.Printf("[Derive WS] Querying positions for subaccount %d with request ID: %s", subaccountID, id)
	shared.DeriveDebugLog("[Derive WS] Full request: %+v", req)

	respChan := c.sendRequest(req)

	select {
	case resp := <-respChan:
		// Check if response is empty
		if len(resp) == 0 {
			log.Printf("[Derive WS] ERROR: Received empty response for positions request")
			return nil, fmt.Errorf("received empty response from WebSocket")
		}

		var result struct {
			Result struct {
				SubaccountID int                      `json:"subaccount_id"`
				Positions    []map[string]interface{} `json:"positions"`
			} `json:"result"`
			Error *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}

		if err := json.Unmarshal(resp, &result); err != nil {
			// Log the response that failed to parse
			if len(resp) < 1000 {
				log.Printf("[Derive WS] Failed to parse positions response. Raw response: %s", string(resp))
			} else {
				log.Printf("[Derive WS] Failed to parse positions response. Raw response length: %d, first 1000 chars: %s", len(resp), string(resp[:1000]))
			}
			return nil, fmt.Errorf("failed to parse positions response: %w", err)
		}

		if result.Error != nil {
			return nil, fmt.Errorf("get positions error: %s", result.Error.Message)
		}

		shared.DeriveDebugLog("[Derive WS] Found %d positions", len(result.Result.Positions))

		return result.Result.Positions, nil

	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("get positions timeout")
	}
}

// Close closes the WebSocket connection
func (c *DeriveWSClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Signal shutdown to all goroutines
	close(c.shutdownChan)

	c.isConnected = false
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

// IsConnected returns the connection status
func (c *DeriveWSClient) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isConnected
}

// GetConnectionInfo returns connection status information
func (c *DeriveWSClient) GetConnectionInfo() (bool, time.Time, time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isConnected, c.lastActivity, time.Since(c.lastActivity)
}

// SubscribeOrderBook subscribes to orderbook updates for an instrument
func (c *DeriveWSClient) SubscribeOrderBook(instrument string, depth int) error {
	// Determine depth parameter
	depthParam := "10"
	if depth > 10 && depth <= 20 {
		depthParam = "20"
	} else if depth > 20 {
		depthParam = "100"
	}

	channel := fmt.Sprintf("orderbook.%s.1.%s", instrument, depthParam)

	// Check if already subscribed
	c.orderbookMu.Lock()
	if c.orderbookSubs[channel] {
		c.orderbookMu.Unlock()
		return nil
	}
	c.orderbookSubs[channel] = true
	c.orderbookMu.Unlock()

	// Subscribe via WebSocket
	msg := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"params": map[string]interface{}{
			"channels": []string{channel},
		},
		"id": fmt.Sprintf("subscribe_%s_%d", instrument, time.Now().UnixNano()),
	}

	c.mu.Lock()
	if !c.isConnected || c.conn == nil {
		c.mu.Unlock()
		c.orderbookMu.Lock()
		delete(c.orderbookSubs, channel)
		c.orderbookMu.Unlock()
		return fmt.Errorf("cannot subscribe: connection not available")
	}
	conn := c.conn
	c.mu.Unlock()

	c.writeMu.Lock()
	err := conn.WriteJSON(msg)
	c.writeMu.Unlock()

	if err != nil {
		c.orderbookMu.Lock()
		delete(c.orderbookSubs, channel)
		c.orderbookMu.Unlock()
		return fmt.Errorf("failed to subscribe to orderbook: %w", err)
	}

	shared.DeriveDebugLog("[Derive WS] Subscribed to orderbook channel: %s", channel)
	return nil
}

// GetOrderBook returns the cached orderbook for an instrument
func (c *DeriveWSClient) GetOrderBook(instrument string) *OrderBookData {
	c.orderbookMu.RLock()
	defer c.orderbookMu.RUnlock()

	if ob, ok := c.orderbooks[instrument]; ok {
		// Return a copy to avoid race conditions
		return &OrderBookData{
			Bids:      append([]types.OrderBookLevel{}, ob.Bids...),
			Asks:      append([]types.OrderBookLevel{}, ob.Asks...),
			Timestamp: ob.Timestamp,
			ChangeID:  ob.ChangeID,
		}
	}
	return nil
}
