package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shopspring/decimal"
)

// DeriveMarketMakerExchange implements MarketMakerExchange for Derive/Lyra
type DeriveMarketMakerExchange struct {
	wsClient    *DeriveWSClient
	subaccountID uint64
	
	// Ticker subscriptions
	tickerConn   *websocket.Conn
	tickerMu     sync.Mutex
	subscriptions map[string]bool
}

// NewDeriveMarketMakerExchange creates a new Derive exchange adapter
func NewDeriveMarketMakerExchange(privateKey, walletAddress string) (*DeriveMarketMakerExchange, error) {
	// Create WebSocket client
	wsClient, err := NewDeriveWSClient(privateKey, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create WebSocket client: %w", err)
	}
	
	// Get default subaccount
	subaccountID := wsClient.GetDefaultSubaccount()
	
	return &DeriveMarketMakerExchange{
		wsClient:      wsClient,
		subaccountID:  subaccountID,
		subscriptions: make(map[string]bool),
	}, nil
}

// SubscribeTickers subscribes to real-time ticker updates
func (d *DeriveMarketMakerExchange) SubscribeTickers(ctx context.Context, instruments []string) (<-chan TickerUpdate, error) {
	tickerChan := make(chan TickerUpdate, 100)
	
	// Connect to ticker WebSocket
	wsURL := "wss://api.lyra.finance/ws"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ticker WebSocket: %w", err)
	}
	
	d.tickerMu.Lock()
	d.tickerConn = conn
	d.tickerMu.Unlock()
	
	// Subscribe to each instrument
	for _, instrument := range instruments {
		subscribeMsg := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "public/subscribe",
			"params": map[string]interface{}{
				"instrument_name": instrument,
				"channels": []string{
					"ticker",
					"orderbook", // For best bid/ask updates
				},
			},
			"id": fmt.Sprintf("sub_%s_%d", instrument, time.Now().UnixNano()),
		}
		
		if err := conn.WriteJSON(subscribeMsg); err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to subscribe to %s: %w", instrument, err)
		}
		
		d.subscriptions[instrument] = true
		log.Printf("Subscribed to ticker updates for %s", instrument)
	}
	
	// Start processing messages
	go d.processTickerMessages(ctx, conn, tickerChan)
	
	return tickerChan, nil
}

// processTickerMessages handles incoming WebSocket messages
func (d *DeriveMarketMakerExchange) processTickerMessages(ctx context.Context, conn *websocket.Conn, tickerChan chan<- TickerUpdate) {
	defer close(tickerChan)
	defer conn.Close()
	
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Set read deadline
			conn.SetReadDeadline(time.Now().Add(60 * time.Second))
			
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading ticker message: %v", err)
				return
			}
			
			// Parse message
			var msg struct {
				Method string          `json:"method"`
				Params json.RawMessage `json:"params"`
			}
			
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("Failed to parse ticker message: %v", err)
				continue
			}
			
			// Handle ticker updates
			if msg.Method == "ticker" || msg.Method == "orderbook" {
				d.handleTickerUpdate(msg.Params, tickerChan)
			}
		}
	}
}

// handleTickerUpdate processes ticker/orderbook updates
func (d *DeriveMarketMakerExchange) handleTickerUpdate(params json.RawMessage, tickerChan chan<- TickerUpdate) {
	var data struct {
		InstrumentName string          `json:"instrument_name"`
		BestBid        decimal.Decimal `json:"best_bid_price"`
		BestBidAmount  decimal.Decimal `json:"best_bid_amount"`
		BestAsk        decimal.Decimal `json:"best_ask_price"`
		BestAskAmount  decimal.Decimal `json:"best_ask_amount"`
		LastPrice      decimal.Decimal `json:"last_price"`
		MarkPrice      decimal.Decimal `json:"mark_price"`
	}
	
	if err := json.Unmarshal(params, &data); err != nil {
		log.Printf("Failed to parse ticker data: %v", err)
		return
	}
	
	update := TickerUpdate{
		Instrument:  data.InstrumentName,
		BestBid:     data.BestBid,
		BestBidSize: data.BestBidAmount,
		BestAsk:     data.BestAsk,
		BestAskSize: data.BestAskAmount,
		LastPrice:   data.LastPrice,
		MarkPrice:   data.MarkPrice,
		Timestamp:   time.Now(),
	}
	
	select {
	case tickerChan <- update:
	default:
		log.Printf("Ticker channel full, dropping update for %s", data.InstrumentName)
	}
}

// PlaceLimitOrder places a limit order on Derive
func (d *DeriveMarketMakerExchange) PlaceLimitOrder(instrument string, side string, price, amount decimal.Decimal) (string, error) {
	// Convert to Derive order format
	order := map[string]interface{}{
		"subaccount_id":   d.subaccountID,
		"instrument_name": instrument,
		"direction":       side, // "buy" or "sell"
		"order_type":      "limit",
		"limit_price":     price.String(),
		"amount":          amount.String(),
		"time_in_force":   "gtc", // Good till cancelled
		"mmp":             true,  // Market maker protection
	}
	
	// Sign and submit order
	response, err := d.wsClient.SubmitOrder(order)
	if err != nil {
		return "", fmt.Errorf("failed to submit order: %w", err)
	}
	
	return response.Result.OrderID, nil
}

// CancelOrder cancels an order on Derive
func (d *DeriveMarketMakerExchange) CancelOrder(orderID string) error {
	id := fmt.Sprintf("%d", time.Now().UnixMilli())
	
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "private/cancel_order",
		"params": map[string]interface{}{
			"order_id": orderID,
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
func (d *DeriveMarketMakerExchange) GetOpenOrders() ([]MarketMakerOrder, error) {
	rawOrders, err := d.wsClient.GetOpenOrders(d.subaccountID)
	if err != nil {
		return nil, err
	}
	
	orders := make([]MarketMakerOrder, 0, len(rawOrders))
	for _, raw := range rawOrders {
		order := MarketMakerOrder{
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
func (d *DeriveMarketMakerExchange) GetPositions() ([]ExchangePosition, error) {
	rawPositions, err := d.wsClient.GetPositions(d.subaccountID)
	if err != nil {
		return nil, err
	}
	
	positions := make([]ExchangePosition, 0, len(rawPositions))
	for _, raw := range rawPositions {
		position := ExchangePosition{
			InstrumentName: getString(raw, "instrument_name"),
			Amount:         getFloat64(raw, "amount"),
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
	if v, ok := m[key].(float64); ok {
		return v
	}
	return 0
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