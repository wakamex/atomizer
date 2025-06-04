package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// MessageHandler processes incoming WebSocket messages
type MessageHandler func(messageType int, data []byte)

// AuthProvider handles authentication for the WebSocket connection
type AuthProvider interface {
	// Authenticate performs authentication after connection
	Authenticate(conn *websocket.Conn) error
	// RequiresAuth returns true if authentication is needed
	RequiresAuth() bool
}

// ClientConfig contains configuration for the WebSocket client
type ClientConfig struct {
	URL                string
	Name               string
	AuthProvider       AuthProvider
	MessageHandler     MessageHandler
	ReconnectDelay     time.Duration
	MaxReconnectDelay  time.Duration
	PingInterval       time.Duration
	PongTimeout        time.Duration
	MaxMessageSize     int64
}

// Client is a generic WebSocket client with reconnection and auth support
type Client struct {
	config         ClientConfig
	conn           *websocket.Conn
	mu             sync.Mutex
	writeMu        sync.Mutex
	
	// Connection state
	isConnected    bool
	lastActivity   time.Time
	reconnectChan  chan struct{}
	shutdownChan   chan struct{}
	
	// Request tracking for RPC-style communication
	requests       map[string]chan json.RawMessage
	requestsMu     sync.Mutex
	
	// Subscription management
	subscriptions  map[string]bool
	subMu          sync.RWMutex
}

// NewClient creates a new WebSocket client
func NewClient(config ClientConfig) *Client {
	// Set defaults
	if config.ReconnectDelay == 0 {
		config.ReconnectDelay = 1 * time.Second
	}
	if config.MaxReconnectDelay == 0 {
		config.MaxReconnectDelay = 30 * time.Second
	}
	if config.PingInterval == 0 {
		config.PingInterval = 30 * time.Second
	}
	if config.PongTimeout == 0 {
		config.PongTimeout = 10 * time.Second
	}
	if config.MaxMessageSize == 0 {
		config.MaxMessageSize = 1024 * 1024 // 1MB
	}
	
	return &Client{
		config:        config,
		reconnectChan: make(chan struct{}, 1),
		shutdownChan:  make(chan struct{}),
		requests:      make(map[string]chan json.RawMessage),
		subscriptions: make(map[string]bool),
	}
}

// Connect establishes the WebSocket connection
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	if c.isConnected {
		return nil
	}
	
	return c.connect()
}

// connect performs the actual connection (must be called with mutex held)
func (c *Client) connect() error {
	log.Printf("[%s] Connecting to %s", c.config.Name, c.config.URL)
	
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second
	
	conn, _, err := dialer.Dial(c.config.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}
	
	c.conn = conn
	c.conn.SetReadLimit(c.config.MaxMessageSize)
	
	// Set up ping/pong handlers
	c.conn.SetPongHandler(func(string) error {
		c.updateActivity()
		return nil
	})
	
	// Authenticate if required
	if c.config.AuthProvider != nil && c.config.AuthProvider.RequiresAuth() {
		if err := c.config.AuthProvider.Authenticate(conn); err != nil {
			conn.Close()
			return fmt.Errorf("authentication failed: %w", err)
		}
	}
	
	c.isConnected = true
	c.lastActivity = time.Now()
	
	// Start read pump
	go c.readPump()
	
	// Start ping ticker
	go c.pingPump()
	
	log.Printf("[%s] Successfully connected", c.config.Name)
	
	return nil
}

// Start begins the connection management loop
func (c *Client) Start(ctx context.Context) error {
	// Initial connection
	if err := c.Connect(); err != nil {
		return err
	}
	
	// Start connection monitor
	go c.connectionMonitor(ctx)
	
	// Wait for context cancellation
	<-ctx.Done()
	
	c.Close()
	return ctx.Err()
}

// Close gracefully shuts down the connection
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	close(c.shutdownChan)
	
	if c.conn != nil {
		c.conn.Close()
	}
	
	c.isConnected = false
}

// Send sends a message over the WebSocket connection
func (c *Client) Send(data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	
	c.mu.Lock()
	if !c.isConnected || c.conn == nil {
		c.mu.Unlock()
		return fmt.Errorf("not connected")
	}
	c.mu.Unlock()
	
	c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// SendJSON sends a JSON message
func (c *Client) SendJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.Send(data)
}

// Request sends an RPC-style request and waits for response
func (c *Client) Request(id string, method string, params interface{}) (json.RawMessage, error) {
	// Create response channel
	respChan := make(chan json.RawMessage, 1)
	
	c.requestsMu.Lock()
	c.requests[id] = respChan
	c.requestsMu.Unlock()
	
	// Clean up on return
	defer func() {
		c.requestsMu.Lock()
		delete(c.requests, id)
		c.requestsMu.Unlock()
	}()
	
	// Send request
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"method":  method,
		"params":  params,
	}
	
	if err := c.SendJSON(request); err != nil {
		return nil, err
	}
	
	// Wait for response
	select {
	case resp := <-respChan:
		return resp, nil
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("request timeout")
	}
}

// readPump handles incoming messages
func (c *Client) readPump() {
	defer func() {
		c.triggerReconnect()
	}()
	
	for {
		messageType, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[%s] WebSocket error: %v", c.config.Name, err)
			}
			return
		}
		
		c.updateActivity()
		
		// Handle RPC responses
		var rpcMsg struct {
			ID     string          `json:"id"`
			Result json.RawMessage `json:"result"`
		}
		
		if err := json.Unmarshal(data, &rpcMsg); err == nil && rpcMsg.ID != "" {
			c.requestsMu.Lock()
			if respChan, ok := c.requests[rpcMsg.ID]; ok {
				select {
				case respChan <- rpcMsg.Result:
				default:
				}
			}
			c.requestsMu.Unlock()
		}
		
		// Call message handler
		if c.config.MessageHandler != nil {
			c.config.MessageHandler(messageType, data)
		}
	}
}

// pingPump sends periodic pings
func (c *Client) pingPump() {
	ticker := time.NewTicker(c.config.PingInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			c.writeMu.Lock()
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.writeMu.Unlock()
				return
			}
			c.writeMu.Unlock()
			
		case <-c.shutdownChan:
			return
		}
	}
}

// connectionMonitor monitors connection health and reconnects as needed
func (c *Client) connectionMonitor(ctx context.Context) {
	for {
		select {
		case <-c.reconnectChan:
			c.performReconnection()
			
		case <-ctx.Done():
			return
			
		case <-c.shutdownChan:
			return
		}
	}
}

// performReconnection handles reconnection with exponential backoff
func (c *Client) performReconnection() {
	c.mu.Lock()
	c.isConnected = false
	delay := c.config.ReconnectDelay
	c.mu.Unlock()
	
	for {
		select {
		case <-c.shutdownChan:
			return
		default:
		}
		
		log.Printf("[%s] Attempting reconnection in %v", c.config.Name, delay)
		time.Sleep(delay)
		
		c.mu.Lock()
		err := c.connect()
		c.mu.Unlock()
		
		if err == nil {
			// Resubscribe
			c.resubscribe()
			return
		}
		
		log.Printf("[%s] Reconnection failed: %v", c.config.Name, err)
		
		// Exponential backoff
		delay *= 2
		if delay > c.config.MaxReconnectDelay {
			delay = c.config.MaxReconnectDelay
		}
	}
}

// resubscribe re-establishes subscriptions after reconnection
func (c *Client) resubscribe() {
	c.subMu.RLock()
	subs := make([]string, 0, len(c.subscriptions))
	for sub := range c.subscriptions {
		subs = append(subs, sub)
	}
	c.subMu.RUnlock()
	
	// Resubscribe logic would go here
	// This is application-specific
	log.Printf("[%s] Resubscribing to %d channels", c.config.Name, len(subs))
}

// triggerReconnect triggers a reconnection attempt
func (c *Client) triggerReconnect() {
	select {
	case c.reconnectChan <- struct{}{}:
	default:
	}
}

// updateActivity updates the last activity timestamp
func (c *Client) updateActivity() {
	c.mu.Lock()
	c.lastActivity = time.Now()
	c.mu.Unlock()
}

// IsConnected returns the connection status
func (c *Client) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isConnected
}