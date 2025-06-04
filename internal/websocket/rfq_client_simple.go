package websocket

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/wakamex/atomizer/internal/arbitrage"
	"github.com/wakamex/atomizer/internal/config"
	"github.com/wakamex/atomizer/internal/rfq"
	"github.com/wakamex/atomizer/internal/types"
)

// SimpleRFQClient manages WebSocket connections for RFQ streaming using the generic client
type SimpleRFQClient struct {
	config       *config.Config
	orchestrator *arbitrage.Orchestrator
	processor    *rfq.Processor
	manager      *RyskRFQManager
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewSimpleRFQClient creates a new RFQ WebSocket client
func NewSimpleRFQClient(cfg *config.Config, orchestrator *arbitrage.Orchestrator, processor *rfq.Processor) *SimpleRFQClient {
	ctx, cancel := context.WithCancel(context.Background())
	
	client := &SimpleRFQClient{
		config:       cfg,
		orchestrator: orchestrator,
		processor:    processor,
		ctx:          ctx,
		cancel:       cancel,
	}
	
	// Create RFQ manager with message handler
	client.manager = NewRyskRFQManager(cfg.WebSocketURL, func(asset string, notification types.RFQNotification) {
		client.handleNotification(asset, notification)
	})
	
	return client
}

// Start begins the WebSocket connections
func (c *SimpleRFQClient) Start() error {
	log.Println("Starting RFQ WebSocket client...")
	
	// Parse assets from config
	assets := strings.Split(c.config.RFQAssetAddressesCSV, ",")
	
	// Connect using manager
	return c.manager.Connect(c.ctx, assets)
}

// Stop gracefully shuts down all connections
func (c *SimpleRFQClient) Stop() {
	log.Println("Stopping RFQ WebSocket client...")
	c.cancel()
	c.manager.Close()
}

// handleNotification processes an RFQ notification
func (c *SimpleRFQClient) handleNotification(source string, notification types.RFQNotification) {
	// Handle different message types
	switch notification.Method {
	case "rfq":
		c.handleRFQRequest(notification)
		
	case "trade_confirmation":
		c.handleTradeConfirmation(notification)
		
	default:
		// Check if it's a result message
		if notification.Result.Asset != "" {
			c.handleRFQRequest(notification)
		}
	}
}

// handleRFQRequest processes an RFQ and generates a quote
func (c *SimpleRFQClient) handleRFQRequest(notification types.RFQNotification) {
	rfqResult := notification.Result
	if notification.Params.Asset != "" {
		rfqResult = notification.Params
	}
	
	// Set RFQ ID
	if rfqResult.RFQId == "" {
		rfqResult.RFQId = notification.ID
	}
	
	log.Printf("Received RFQ: ID=%s, Asset=%s, Strike=%s, Expiry=%d, IsPut=%t, Quantity=%s",
		rfqResult.RFQId, rfqResult.Asset, rfqResult.Strike, rfqResult.Expiry, 
		rfqResult.IsPut, rfqResult.Quantity)
	
	// Create a simple client adapter that implements rfq.RyskClient
	client := &ryskClientAdapter{manager: c.manager}
	
	// Generate and send quote
	if err := c.processor.ProcessRFQ(client, rfqResult, rfqResult.RFQId); err != nil {
		log.Printf("Failed to process RFQ: %v", err)
	}
}

// handleTradeConfirmation processes a trade confirmation
func (c *SimpleRFQClient) handleTradeConfirmation(notification types.RFQNotification) {
	log.Printf("Received trade confirmation for RFQ %s", notification.ID)
	
	// Parse confirmation using the proper parser
	conf, err := ParseTradeConfirmation(notification)
	if err != nil {
		log.Printf("Failed to parse trade confirmation: %v", err)
		return
	}
	
	// Check asset mapping
	if c.config.AssetMapping != nil {
		if underlying, hasMapping := c.config.AssetMapping[conf.AssetAddress]; hasMapping {
			log.Printf("Trade confirmation for %s (%s), Quote ID: %s", conf.AssetAddress, underlying, conf.QuoteNonce)
		} else {
			log.Printf("WARNING: No asset mapping for %s. Cannot hedge.", conf.AssetAddress)
		}
	}
	
	// Get RFQ result for additional context
	rfqResult := notification.Params
	if rfqResult.Asset == "" {
		rfqResult = notification.Result
	}
	
	// Submit to orchestrator
	if err := c.orchestrator.SubmitRFQTrade(rfqResult, conf); err != nil {
		log.Printf("Failed to submit trade to orchestrator: %v", err)
	} else {
		log.Printf("Successfully sent trade to orchestrator for Quote ID %s", conf.QuoteNonce)
	}
}

// ryskClientAdapter adapts the manager to the RyskClient interface
type ryskClientAdapter struct {
	manager *RyskRFQManager
}

func (r *ryskClientAdapter) Send(data []byte) {
	// Parse as JSON-RPC request to send properly
	var req map[string]interface{}
	if err := json.Unmarshal(data, &req); err == nil {
		r.manager.SendQuote(req)
	}
}