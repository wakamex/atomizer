package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	
	"github.com/wakamex/atomizer/internal/types"
)

// RyskRFQClient wraps the generic WebSocket client for Rysk RFQ
type RyskRFQClient struct {
	client         *Client
	messageHandler func(types.RFQNotification)
}

// NewRyskRFQClient creates a new Rysk RFQ WebSocket client
func NewRyskRFQClient(url string, handler func(types.RFQNotification)) *RyskRFQClient {
	r := &RyskRFQClient{
		messageHandler: handler,
	}
	
	config := ClientConfig{
		URL:            url,
		Name:           "RyskRFQ",
		AuthProvider:   nil, // Rysk doesn't require auth for RFQ streams
		MessageHandler: r.handleMessage,
	}
	
	r.client = NewClient(config)
	return r
}

// Start begins the WebSocket connection
func (r *RyskRFQClient) Start(ctx context.Context) error {
	return r.client.Start(ctx)
}

// Send sends data to the WebSocket
func (r *RyskRFQClient) Send(data []byte) {
	if err := r.client.Send(data); err != nil {
		log.Printf("[RyskRFQ] Failed to send: %v", err)
	}
}

// Close closes the WebSocket connection
func (r *RyskRFQClient) Close() {
	r.client.Close()
}

// handleMessage processes incoming WebSocket messages
func (r *RyskRFQClient) handleMessage(messageType int, data []byte) {
	// Parse as RFQ notification
	var notification types.RFQNotification
	if err := json.Unmarshal(data, &notification); err != nil {
		log.Printf("[RyskRFQ] Failed to parse message: %v", err)
		return
	}
	
	// Call the handler
	if r.messageHandler != nil {
		r.messageHandler(notification)
	}
}

// RyskRFQManager manages multiple RFQ connections for different assets
type RyskRFQManager struct {
	baseURL    string
	mainClient *RyskRFQClient
	rfqClients map[string]*RyskRFQClient
	handler    func(asset string, notification types.RFQNotification)
}

// NewRyskRFQManager creates a new RFQ manager
func NewRyskRFQManager(baseURL string, handler func(string, types.RFQNotification)) *RyskRFQManager {
	return &RyskRFQManager{
		baseURL:    baseURL,
		rfqClients: make(map[string]*RyskRFQClient),
		handler:    handler,
	}
}

// Connect establishes connections for the given assets
func (m *RyskRFQManager) Connect(ctx context.Context, assets []string) error {
	// Connect to main endpoint
	mainURL := m.baseURL
	m.mainClient = NewRyskRFQClient(mainURL, func(n types.RFQNotification) {
		m.handler("main", n)
	})
	
	if err := m.mainClient.Start(ctx); err != nil {
		return fmt.Errorf("failed to start main client: %w", err)
	}
	
	// Connect to asset-specific RFQ streams
	baseURL := strings.TrimSuffix(m.baseURL, "/maker")
	
	for _, asset := range assets {
		asset = strings.TrimSpace(asset)
		if asset == "" {
			continue
		}
		
		rfqURL := fmt.Sprintf("%s/rfq?asset=%s", baseURL, asset)
		
		client := NewRyskRFQClient(rfqURL, func(n types.RFQNotification) {
			m.handler(asset, n)
		})
		
		m.rfqClients[asset] = client
		
		// Start in background
		go func(a string) {
			if err := client.Start(ctx); err != nil {
				log.Printf("[RyskRFQ] Failed to start client for asset %s: %v", a, err)
			}
		}(asset)
	}
	
	return nil
}

// SendQuote sends a quote response
func (m *RyskRFQManager) SendQuote(quote interface{}) error {
	data, err := json.Marshal(quote)
	if err != nil {
		return err
	}
	
	if m.mainClient != nil {
		m.mainClient.Send(data)
	}
	
	return nil
}

// Close closes all connections
func (m *RyskRFQManager) Close() {
	if m.mainClient != nil {
		m.mainClient.Close()
	}
	
	for _, client := range m.rfqClients {
		client.Close()
	}
}