package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// SpotPrice represents a spot price update
type SpotPrice struct {
	Currency  string
	Price     float64
	Timestamp time.Time
}

// DeriveSpotCollector collects spot prices via WebSocket
type DeriveSpotCollector struct {
	conn        *websocket.Conn
	spotPrices  map[string]*SpotPrice
	mu          sync.RWMutex
	subscribers map[string]bool
	ctx         context.Context
	cancel      context.CancelFunc
	reconnectCh chan struct{}
}

// DeriveSpotUpdate represents spot price update from WebSocket
type DeriveSpotUpdate struct {
	Channel string `json:"channel"`
	Data    struct {
		Timestamp int64                  `json:"timestamp"`
		Feeds     map[string]interface{} `json:"feeds"`
	} `json:"data"`
}

func NewDeriveSpotCollector() (*DeriveSpotCollector, error) {
	ctx, cancel := context.WithCancel(context.Background())

	collector := &DeriveSpotCollector{
		spotPrices:  make(map[string]*SpotPrice),
		subscribers: make(map[string]bool),
		ctx:         ctx,
		cancel:      cancel,
		reconnectCh: make(chan struct{}, 1),
	}

	// Connect to WebSocket
	if err := collector.connect(); err != nil {
		cancel()
		return nil, err
	}

	// Start message handler
	go collector.handleMessages()

	// Start reconnection handler
	go collector.handleReconnection()

	return collector, nil
}

func (d *DeriveSpotCollector) connect() error {
	wsURL := "wss://api.lyra.finance/ws"
	log.Printf("[Derive Spot] Connecting to %s", wsURL)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to Derive WebSocket: %w", err)
	}

	d.conn = conn

	// Set up ping/pong handlers
	conn.SetPingHandler(func(appData string) error {
		log.Printf("[Derive Spot] Ping received, sending Pong")
		return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(5*time.Second))
	})

	// Send heartbeat periodically
	go d.sendHeartbeat()

	// Re-subscribe to all currencies
	d.mu.RLock()
	currencies := make([]string, 0, len(d.subscribers))
	for currency := range d.subscribers {
		currencies = append(currencies, currency)
	}
	d.mu.RUnlock()

	for _, currency := range currencies {
		if err := d.subscribeCurrency(currency); err != nil {
			log.Printf("[Derive Spot] Failed to resubscribe to %s: %v", currency, err)
		}
	}

	return nil
}

func (d *DeriveSpotCollector) handleReconnection() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case <-d.reconnectCh:
			log.Println("[Derive Spot] Attempting to reconnect...")

			// Wait a bit before reconnecting
			time.Sleep(5 * time.Second)

			if err := d.connect(); err != nil {
				log.Printf("[Derive Spot] Reconnection failed: %v", err)
				// Trigger another reconnection attempt
				select {
				case d.reconnectCh <- struct{}{}:
				default:
				}
			} else {
				log.Println("[Derive Spot] Reconnected successfully")
			}
		}
	}
}

func (d *DeriveSpotCollector) sendHeartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-d.ctx.Done():
			return
		case <-ticker.C:
			msg := map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  "public/heartbeat",
				"params":  map[string]interface{}{},
			}

			if err := d.conn.WriteJSON(msg); err != nil {
				log.Printf("[Derive Spot] Failed to send heartbeat: %v", err)
				d.triggerReconnect()
				return
			}
		}
	}
}

func (d *DeriveSpotCollector) triggerReconnect() {
	select {
	case d.reconnectCh <- struct{}{}:
	default:
	}
}

func (d *DeriveSpotCollector) handleMessages() {
	for {
		select {
		case <-d.ctx.Done():
			return
		default:
			var msg json.RawMessage
			if err := d.conn.ReadJSON(&msg); err != nil {
				log.Printf("[Derive Spot] Read error: %v", err)
				d.triggerReconnect()
				return
			}

			// Debug: Log raw message
			msgStr := string(msg)
			if len(msgStr) > 200 {
				Debugf("[DEBUG] Derive Spot message (truncated): %s...", msgStr[:200])
			} else {
				Debugf("[DEBUG] Derive Spot message: %s", msgStr)
			}

			// Try to parse as subscription update
			var update struct {
				Method string           `json:"method"`
				Params DeriveSpotUpdate `json:"params"`
			}

			if err := json.Unmarshal(msg, &update); err == nil && update.Method == "subscription" {
				d.processSpotUpdate(update.Params)
			}
		}
	}
}

func (d *DeriveSpotCollector) processSpotUpdate(update DeriveSpotUpdate) {
	// Extract currency from channel name (e.g., "spot_feed.ETH" -> "ETH")
	var currency string
	if _, err := fmt.Sscanf(update.Channel, "spot_feed.%s", &currency); err != nil {
		log.Printf("[Derive Spot] Failed to parse channel name: %s", update.Channel)
		return
	}

	// Extract price from feeds
	// The feeds structure contains currency -> {price, confidence, etc}
	var spotPrice float64
	foundPrice := false

	// Check if feeds contains currency data
	if currencyData, ok := update.Data.Feeds[currency]; ok {
		// currencyData is a map with price, confidence, etc.
		if priceData, ok := currencyData.(map[string]interface{}); ok {
			// Look for price field
			if priceVal, ok := priceData["price"]; ok {
				switch v := priceVal.(type) {
				case float64:
					spotPrice = v
					foundPrice = true
				case string:
					if f, err := strconv.ParseFloat(v, 64); err == nil {
						spotPrice = f
						foundPrice = true
					}
				case json.Number:
					if f, err := v.Float64(); err == nil {
						spotPrice = f
						foundPrice = true
					}
				}
			}
		}
	}

	if !foundPrice {
		log.Printf("[Derive Spot] No price found in feeds for %s: %+v", currency, update.Data.Feeds)
		return
	}

	// Store the spot price
	spotUpdate := &SpotPrice{
		Currency:  currency,
		Price:     spotPrice,
		Timestamp: time.Unix(0, update.Data.Timestamp*1000000), // Convert microseconds to nanoseconds
	}

	d.mu.Lock()
	d.spotPrices[currency] = spotUpdate
	d.mu.Unlock()

	log.Printf("[Derive Spot] %s: $%.2f", currency, spotPrice)
}

func (d *DeriveSpotCollector) Subscribe(currencies []string) error {
	for _, currency := range currencies {
		if err := d.subscribeCurrency(currency); err != nil {
			return err
		}

		d.mu.Lock()
		d.subscribers[currency] = true
		d.mu.Unlock()
	}

	return nil
}

func (d *DeriveSpotCollector) subscribeCurrency(currency string) error {
	channel := fmt.Sprintf("spot_feed.%s", currency)

	msg := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"params": map[string]interface{}{
			"channels": []string{channel},
		},
		"id": fmt.Sprintf("subscribe_spot_%s_%d", currency, time.Now().Unix()),
	}

	if err := d.conn.WriteJSON(msg); err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", currency, err)
	}

	log.Printf("[Derive Spot] Subscribed to %s", channel)
	return nil
}

func (d *DeriveSpotCollector) GetSpotPrice(currency string) (float64, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if spot, ok := d.spotPrices[currency]; ok {
		// Only return if price is recent (less than 60 seconds old)
		if time.Since(spot.Timestamp) < 60*time.Second {
			return spot.Price, true
		}
	}

	return 0, false
}

func (d *DeriveSpotCollector) GetAllSpotPrices() map[string]SpotPrice {
	d.mu.RLock()
	defer d.mu.RUnlock()

	prices := make(map[string]SpotPrice)
	for currency, spot := range d.spotPrices {
		if time.Since(spot.Timestamp) < 60*time.Second {
			prices[currency] = *spot
		}
	}

	return prices
}

func (d *DeriveSpotCollector) Close() error {
	d.cancel()

	if d.conn != nil {
		// Unsubscribe from all currencies
		d.mu.RLock()
		currencies := make([]string, 0, len(d.subscribers))
		for currency := range d.subscribers {
			currencies = append(currencies, currency)
		}
		d.mu.RUnlock()

		for _, currency := range currencies {
			channel := fmt.Sprintf("spot_feed.%s", currency)
			msg := map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  "unsubscribe",
				"params": map[string]interface{}{
					"channels": []string{channel},
				},
				"id": fmt.Sprintf("unsubscribe_spot_%s_%d", currency, time.Now().Unix()),
			}
			d.conn.WriteJSON(msg)
		}

		// Close connection
		d.conn.Close()
	}

	return nil
}
