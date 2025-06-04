package websocket

import (
	"context"
	"encoding/json"
	"log"

	"github.com/wakamex/atomizer/internal/types"
)

// EnhancedRFQHandler provides better message parsing for Rysk RFQ WebSocket
type EnhancedRFQHandler struct {
	onRFQ               func(types.RFQResult)
	onTradeConfirmation func(types.RFQConfirmation)
	onError             func(error)
}

// NewEnhancedRFQHandler creates a handler with callbacks
func NewEnhancedRFQHandler(
	onRFQ func(types.RFQResult),
	onConfirmation func(types.RFQConfirmation),
	onError func(error),
) *EnhancedRFQHandler {
	return &EnhancedRFQHandler{
		onRFQ:               onRFQ,
		onTradeConfirmation: onConfirmation,
		onError:             onError,
	}
}

// CreateMessageHandler returns a MessageHandler function for the generic client
func (h *EnhancedRFQHandler) CreateMessageHandler() MessageHandler {
	return func(messageType int, data []byte) {
		// Parse raw message
		conf, rfq, err := ParseRawConfirmation(data)
		if err != nil {
			// Try parsing as standard notification
			var notification types.RFQNotification
			if err := json.Unmarshal(data, &notification); err != nil {
				if h.onError != nil {
					h.onError(err)
				}
				return
			}
			
			// Handle based on method or ID
			switch {
			case notification.Method == "rfq" || (notification.Result.Asset != "" && notification.Method == ""):
				h.handleRFQ(notification)
			case notification.Method == "rfq_confirmation" || notification.ID == "trade":
				h.handleConfirmation(notification)
			default:
				log.Printf("Unknown message type: method=%s, id=%s", notification.Method, notification.ID)
			}
			return
		}
		
		// Handle parsed results
		if conf != nil && h.onTradeConfirmation != nil {
			h.onTradeConfirmation(*conf)
		}
		if rfq != nil && h.onRFQ != nil {
			h.onRFQ(*rfq)
		}
	}
}

// handleRFQ processes RFQ messages
func (h *EnhancedRFQHandler) handleRFQ(notification types.RFQNotification) {
	if h.onRFQ == nil {
		return
	}
	
	rfqResult := notification.Result
	if notification.Params.Asset != "" {
		rfqResult = notification.Params
	}
	
	// Set RFQ ID
	if rfqResult.RFQId == "" {
		rfqResult.RFQId = notification.ID
	}
	
	h.onRFQ(rfqResult)
}

// handleConfirmation processes trade confirmations
func (h *EnhancedRFQHandler) handleConfirmation(notification types.RFQNotification) {
	if h.onTradeConfirmation == nil {
		return
	}
	
	conf, err := ParseTradeConfirmation(notification)
	if err != nil {
		if h.onError != nil {
			h.onError(err)
		}
		return
	}
	
	h.onTradeConfirmation(*conf)
}

// CreateEnhancedRyskClient creates a Rysk client with enhanced message parsing
func CreateEnhancedRyskClient(
	ctx context.Context,
	wsURL string,
	onRFQ func(types.RFQResult),
	onConfirmation func(types.RFQConfirmation),
) (*Client, error) {
	
	handler := NewEnhancedRFQHandler(onRFQ, onConfirmation, func(err error) {
		log.Printf("RFQ Handler Error: %v", err)
	})
	
	config := ClientConfig{
		URL:            wsURL,
		Name:           "RyskRFQ",
		MessageHandler: handler.CreateMessageHandler(),
	}
	
	client := NewClient(config)
	
	// Start in background
	go func() {
		if err := client.Start(ctx); err != nil {
			log.Printf("RFQ client error: %v", err)
		}
	}()
	
	return client, nil
}