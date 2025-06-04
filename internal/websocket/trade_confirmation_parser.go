package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wakamex/atomizer/internal/types"
)

// ParseTradeConfirmation parses trade confirmations from various message formats
func ParseTradeConfirmation(notification types.RFQNotification) (*types.RFQConfirmation, error) {
	// Case 1: Method is "rfq_confirmation" with params
	if notification.Method == "rfq_confirmation" {
		// Parse from Params which should be JSON
		var conf types.RFQConfirmation
		
		// Params is already unmarshaled as RFQResult, but confirmation has more fields
		// We need to re-parse from the raw JSON if available
		// For now, construct from available fields
		conf = types.RFQConfirmation{
			ID:           notification.ID,
			AssetAddress: notification.Params.Asset,
			Strike:       notification.Params.Strike,
			Quantity:     notification.Params.Quantity,
			IsPut:        notification.Params.IsPut,
			IsTakerBuy:   notification.Params.IsTakerBuy,
			Expiry:       int(notification.Params.Expiry),
			// Note: Price is missing from RFQResult - this is a problem
		}
		
		log.Printf("Parsed rfq_confirmation (partial): Quote ID %s", conf.ID)
		return &conf, nil
	}
	
	// Case 2: ID is "trade" with result containing confirmation
	if notification.ID == "trade" && notification.Result.Asset != "" {
		// Result contains the confirmation details
		conf := types.RFQConfirmation{
			ID:           notification.ID,
			AssetAddress: notification.Result.Asset,
			Strike:       notification.Result.Strike,
			Quantity:     notification.Result.Quantity,
			IsPut:        notification.Result.IsPut,
			IsTakerBuy:   notification.Result.IsTakerBuy,
			Expiry:       int(notification.Result.Expiry),
		}
		
		log.Printf("Parsed trade result confirmation: Quote ID %s", conf.ID)
		return &conf, nil
	}
	
	return nil, fmt.Errorf("unknown confirmation format")
}

// CreateTradeEvent creates a TradeEvent from RFQ result and confirmation
func CreateTradeEvent(rfqResult types.RFQResult, conf *types.RFQConfirmation) *types.TradeEvent {
	// Parse decimals safely
	quantity, _ := decimal.NewFromString(conf.Quantity)
	if quantity.IsZero() && rfqResult.Quantity != "" {
		quantity, _ = decimal.NewFromString(rfqResult.Quantity)
	}
	
	price, _ := decimal.NewFromString(conf.Price)
	strike, _ := decimal.NewFromString(conf.Strike)
	if strike.IsZero() && rfqResult.Strike != "" {
		strike, _ = decimal.NewFromString(rfqResult.Strike)
	}
	
	// Use confirmation data with fallback to RFQ data
	expiry := int64(conf.Expiry)
	if expiry == 0 {
		expiry = rfqResult.Expiry
	}
	
	return &types.TradeEvent{
		ID:         conf.QuoteNonce,
		Source:     types.TradeSourceRysk,
		Status:     types.TradeStatusExecuted,
		RFQId:      rfqResult.RFQId,
		Instrument: conf.AssetAddress,
		Strike:     strike,
		Expiry:     expiry,
		IsPut:      conf.IsPut,
		Quantity:   quantity,
		Price:      price,
		IsTakerBuy: conf.IsTakerBuy,
		Timestamp:  time.Now(),
	}
}

// Enhanced notification structure to handle raw JSON
type RawNotification struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *types.JSONRPCError `json:"error,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// ParseRawConfirmation attempts to parse confirmation from raw JSON
func ParseRawConfirmation(data []byte) (*types.RFQConfirmation, *types.RFQResult, error) {
	var raw RawNotification
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, err
	}
	
	// Handle trade confirmation in result
	if raw.ID == "trade" && raw.Result != nil {
		var conf types.RFQConfirmation
		if err := json.Unmarshal(raw.Result, &conf); err != nil {
			return nil, nil, fmt.Errorf("failed to parse trade confirmation: %w", err)
		}
		return &conf, nil, nil
	}
	
	// Handle rfq_confirmation method
	if raw.Method == "rfq_confirmation" && raw.Params != nil {
		var conf types.RFQConfirmation
		if err := json.Unmarshal(raw.Params, &conf); err != nil {
			return nil, nil, fmt.Errorf("failed to parse rfq_confirmation: %w", err)
		}
		return &conf, nil, nil
	}
	
	// Handle RFQ request
	if raw.Method == "rfq" || (raw.Result != nil && raw.Method == "") {
		var rfq types.RFQResult
		if raw.Params != nil {
			if err := json.Unmarshal(raw.Params, &rfq); err != nil {
				return nil, nil, fmt.Errorf("failed to parse rfq params: %w", err)
			}
		} else if raw.Result != nil {
			if err := json.Unmarshal(raw.Result, &rfq); err != nil {
				return nil, nil, fmt.Errorf("failed to parse rfq result: %w", err)
			}
		}
		rfq.RFQId = raw.ID
		return nil, &rfq, nil
	}
	
	return nil, nil, fmt.Errorf("unknown message format")
}