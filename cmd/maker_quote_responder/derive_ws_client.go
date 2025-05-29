package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// DeriveWSClient handles WebSocket connection to Derive
type DeriveWSClient struct {
	conn         *websocket.Conn
	auth         *DeriveAuth
	wallet       string
	subaccounts  []int
	mu           sync.Mutex
	requests     map[string]chan json.RawMessage
}

// NewDeriveWSClient creates a new Derive WebSocket client
func NewDeriveWSClient(privateKey string, deriveWallet string) (*DeriveWSClient, error) {
	auth, err := NewDeriveAuth(privateKey)
	if err != nil {
		return nil, err
	}

	client := &DeriveWSClient{
		auth:     auth,
		wallet:   deriveWallet,
		requests: make(map[string]chan json.RawMessage),
	}

	// Connect to Derive WebSocket
	wsURL := "wss://api.lyra.finance/ws"
	log.Printf("[Derive WS] Connecting to %s", wsURL)
	
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Derive WebSocket: %w", err)
	}
	client.conn = conn
	
	// Set up ping/pong handlers with label
	conn.SetPingHandler(func(appData string) error {
		log.Printf("[Derive WS] Ping received, sending Pong")
		err := conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(5*time.Second))
		if err != nil {
			log.Printf("[Derive WS] Error sending pong: %v", err)
		}
		return nil
	})
	
	conn.SetPongHandler(func(appData string) error {
		log.Printf("[Derive WS] Pong received")
		return nil
	})

	// Start message handler
	go client.handleMessages()

	// Login
	if err := client.login(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to login: %w", err)
	}
	
	// Start heartbeat to keep connection alive
	go client.heartbeat()

	return client, nil
}

// heartbeat sends periodic pings to keep the connection alive
func (c *DeriveWSClient) heartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			if c.conn == nil {
				c.mu.Unlock()
				return
			}
			
			// Send ping
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[Derive WS] Error sending ping: %v", err)
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()
			
		case <-time.After(60 * time.Second):
			// If we haven't received anything in 60 seconds, consider the connection dead
			log.Printf("[Derive WS] No activity for 60 seconds, connection may be dead")
			return
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
	
	log.Printf("[Derive WS] Login - EOA: %s", ownerEOA)
	log.Printf("[Derive WS] Login - Derive Wallet: %s", c.wallet)
	log.Printf("[Derive WS] Login - Timestamp: %d (%s)", timestamp, timestampStr)
	log.Printf("[Derive WS] Login - Signature: %s", signature)

	// Use JSON-RPC format - server expects it
	loginReq := map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "public/login",
		"params": map[string]interface{}{
			"wallet":    c.wallet, // Keep original checksummed format
			"timestamp": timestampStr, // String format
			"signature": signature,
		},
		"id": fmt.Sprintf("%d", time.Now().UnixMilli()),
	}

	log.Printf("[Derive WS] Sending login request: %+v", loginReq)

	// Send login request
	respChan := c.sendRequest(loginReq)
	
	select {
	case resp := <-respChan:
		log.Printf("[Derive WS] Login response: %s", string(resp))
		
		var result struct {
			Result []int `json:"result"` // Array of subaccount IDs
			Error  *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(resp, &result); err != nil {
			return fmt.Errorf("failed to parse login response: %w", err)
		}
		if result.Error != nil {
			return fmt.Errorf("login error: %s", result.Error.Message)
		}
		
		// Store subaccount IDs
		c.subaccounts = result.Result
		
		log.Printf("[Derive WS] Login successful. Subaccounts: %v", result.Result)
		return nil
	case <-time.After(10 * time.Second):
		return fmt.Errorf("login timeout")
	}
}

// handleMessages processes incoming WebSocket messages
func (c *DeriveWSClient) handleMessages() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("[Derive WS] Read error: %v", err)
			return
		}

		var msg struct {
			ID     string          `json:"id"`
			Result json.RawMessage `json:"result"`
			Error  json.RawMessage `json:"error"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("[Derive WS] Failed to parse message: %v", err)
			continue
		}

		// Handle response
		c.mu.Lock()
		if ch, ok := c.requests[msg.ID]; ok {
			ch <- message
			delete(c.requests, msg.ID)
		}
		c.mu.Unlock()
	}
}

// sendRequest sends a request and returns a channel for the response
func (c *DeriveWSClient) sendRequest(req map[string]interface{}) <-chan json.RawMessage {
	respChan := make(chan json.RawMessage, 1)
	
	id, _ := req["id"].(string)
	c.mu.Lock()
	c.requests[id] = respChan
	c.mu.Unlock()

	// Marshal to JSON to log exact format
	jsonBytes, _ := json.Marshal(req)
	log.Printf("[Derive WS] Sending JSON: %s", string(jsonBytes))
	
	if err := c.conn.WriteJSON(req); err != nil {
		log.Printf("[Derive WS] Failed to send request: %v", err)
		close(respChan)
		return respChan
	}

	return respChan
}

// SubmitOrder submits an order via WebSocket
func (c *DeriveWSClient) SubmitOrder(order map[string]interface{}) (*DeriveOrderResponse, error) {
	id := fmt.Sprintf("%d", time.Now().UnixMilli())
	
	// Use JSON-RPC format
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "private/order",
		"params": order,
		"id":     id,
	}

	log.Printf("[Derive WS] Submitting order: %+v", order)
	
	respChan := c.sendRequest(req)
	
	select {
	case resp := <-respChan:
		log.Printf("[Derive WS] Order response: %s", string(resp))
		
		// First check if there's an error
		var errorCheck struct {
			Error *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
				Data    *struct {
					Limit     string `json:"limit"`
					Bandwidth string `json:"bandwidth"`
				} `json:"data"`
			} `json:"error"`
		}
		if err := json.Unmarshal(resp, &errorCheck); err == nil && errorCheck.Error != nil {
			// If we get price band error, include bandwidth info
			if errorCheck.Error.Code == 11013 && errorCheck.Error.Data != nil {
				return nil, fmt.Errorf("order error: %s (limit: %s, bandwidth: %s)", 
					errorCheck.Error.Message, errorCheck.Error.Data.Limit, errorCheck.Error.Data.Bandwidth)
			}
			return nil, fmt.Errorf("order error: %s", errorCheck.Error.Message)
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
		
		log.Printf("[Derive WS] Order placed - ID: %s, Status: %s", result.Result.OrderID, result.Result.Status)
		
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
		"method": "private/get_open_orders",
		"params": map[string]interface{}{
			"subaccount_id": subaccountID,
		},
		"id": id,
	}
	
	log.Printf("[Derive WS] Querying open orders for subaccount %d", subaccountID)
	
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

// GetPositions fetches all positions for the subaccount
func (c *DeriveWSClient) GetPositions(subaccountID uint64) ([]map[string]interface{}, error) {
	id := fmt.Sprintf("%d", time.Now().UnixMilli())
	
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "private/get_positions",
		"params": map[string]interface{}{
			"subaccount_id": subaccountID,
		},
		"id": id,
	}
	
	log.Printf("[Derive WS] Querying positions for subaccount %d", subaccountID)
	
	respChan := c.sendRequest(req)
	
	select {
	case resp := <-respChan:
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
			return nil, fmt.Errorf("failed to parse positions response: %w", err)
		}
		
		if result.Error != nil {
			return nil, fmt.Errorf("get positions error: %s", result.Error.Message)
		}
		
		log.Printf("[Derive WS] Found %d positions", len(result.Result.Positions))
		
		return result.Result.Positions, nil
		
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("get positions timeout")
	}
}

// Close closes the WebSocket connection
func (c *DeriveWSClient) Close() error {
	return c.conn.Close()
}