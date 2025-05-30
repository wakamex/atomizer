package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// DeriveWSOrderBookCollector collects order book data via WebSocket
type DeriveWSOrderBookCollector struct {
	conn         *websocket.Conn
	converter    *InstrumentConverter
	depth        int
	orderBooks   map[string]*OrderBookMetric
	mu           sync.RWMutex
	subscribers  map[string]bool
	ctx          context.Context
	cancel       context.CancelFunc
	reconnectCh  chan struct{}
}

// DeriveWSMessage represents a WebSocket message
type DeriveWSMessage struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
	ID     interface{}     `json:"id,omitempty"`
}

// DeriveOrderBookUpdate represents order book update from WebSocket
type DeriveOrderBookUpdate struct {
	Channel string `json:"channel"`
	Data    struct {
		Timestamp      int64           `json:"timestamp"`
		InstrumentName string          `json:"instrument_name"`
		Bids           [][]json.Number `json:"bids"` // [price, size]
		Asks           [][]json.Number `json:"asks"` // [price, size]
		ChangeID       int64           `json:"change_id"`
	} `json:"data"`
}

func NewDeriveWSOrderBookCollector(depth int) (*DeriveWSOrderBookCollector, error) {
	if depth <= 0 {
		depth = 10
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	collector := &DeriveWSOrderBookCollector{
		converter:   NewInstrumentConverter(),
		depth:       depth,
		orderBooks:  make(map[string]*OrderBookMetric),
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

func (d *DeriveWSOrderBookCollector) connect() error {
	wsURL := "wss://api.lyra.finance/ws"
	log.Printf("[Derive WS OrderBook] Connecting to %s", wsURL)
	
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to Derive WebSocket: %w", err)
	}
	
	d.conn = conn
	
	// Set up ping/pong handlers
	conn.SetPingHandler(func(appData string) error {
		log.Printf("[Derive WS OrderBook] Ping received, sending Pong")
		return conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(5*time.Second))
	})
	
	// Send heartbeat periodically
	go d.sendHeartbeat()
	
	// Re-subscribe to all instruments
	d.mu.RLock()
	instruments := make([]string, 0, len(d.subscribers))
	for inst := range d.subscribers {
		instruments = append(instruments, inst)
	}
	d.mu.RUnlock()
	
	for _, inst := range instruments {
		if err := d.subscribeInstrument(inst); err != nil {
			log.Printf("[Derive WS OrderBook] Failed to resubscribe to %s: %v", inst, err)
		}
	}
	
	return nil
}

func (d *DeriveWSOrderBookCollector) handleReconnection() {
	for {
		select {
		case <-d.ctx.Done():
			return
		case <-d.reconnectCh:
			log.Println("[Derive WS OrderBook] Attempting to reconnect...")
			
			// Wait a bit before reconnecting
			time.Sleep(5 * time.Second)
			
			if err := d.connect(); err != nil {
				log.Printf("[Derive WS OrderBook] Reconnection failed: %v", err)
				// Trigger another reconnection attempt
				select {
				case d.reconnectCh <- struct{}{}:
				default:
				}
			} else {
				log.Println("[Derive WS OrderBook] Reconnected successfully")
			}
		}
	}
}

func (d *DeriveWSOrderBookCollector) sendHeartbeat() {
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
				log.Printf("[Derive WS OrderBook] Failed to send heartbeat: %v", err)
				d.triggerReconnect()
				return
			}
		}
	}
}

func (d *DeriveWSOrderBookCollector) triggerReconnect() {
	select {
	case d.reconnectCh <- struct{}{}:
	default:
	}
}

func (d *DeriveWSOrderBookCollector) handleMessages() {
	for {
		select {
		case <-d.ctx.Done():
			return
		default:
			var msg json.RawMessage
			if err := d.conn.ReadJSON(&msg); err != nil {
				log.Printf("[Derive WS OrderBook] Read error: %v", err)
				d.triggerReconnect()
				return
			}
			
			// Debug: Log raw message
			msgStr := string(msg)
			if len(msgStr) > 200 {
				Debugf("[DEBUG] Derive WS message (truncated): %s...", msgStr[:200])
			} else {
				Debugf("[DEBUG] Derive WS message: %s", msgStr)
			}
			
			// Try to parse as subscription update
			var update struct {
				Params DeriveOrderBookUpdate `json:"params"`
			}
			
			if err := json.Unmarshal(msg, &update); err == nil && update.Params.Channel != "" {
				Debugf("[DEBUG] Processing order book update for %s", update.Params.Data.InstrumentName)
				d.processOrderBookUpdate(update.Params)
			}
		}
	}
}

func (d *DeriveWSOrderBookCollector) processOrderBookUpdate(update DeriveOrderBookUpdate) {
	// Convert to OrderBookMetric
	metric := &OrderBookMetric{
		Exchange:   "derive",
		Instrument: update.Data.InstrumentName,
		Timestamp:  time.Unix(0, update.Data.Timestamp*1000000), // Convert microseconds to nanoseconds
		Bids:       make([]OrderBookLevel, 0, len(update.Data.Bids)),
		Asks:       make([]OrderBookLevel, 0, len(update.Data.Asks)),
	}
	
	// Convert bids
	for _, bid := range update.Data.Bids {
		if len(bid) >= 2 {
			price, _ := bid[0].Float64()
			size, _ := bid[1].Float64()
			metric.Bids = append(metric.Bids, OrderBookLevel{
				Price: price,
				Size:  size,
			})
		}
	}
	
	// Convert asks
	for _, ask := range update.Data.Asks {
		if len(ask) >= 2 {
			price, _ := ask[0].Float64()
			size, _ := ask[1].Float64()
			metric.Asks = append(metric.Asks, OrderBookLevel{
				Price: price,
				Size:  size,
			})
		}
	}
	
	// Store the latest order book
	d.mu.Lock()
	d.orderBooks[update.Data.InstrumentName] = metric
	d.mu.Unlock()
}

func (d *DeriveWSOrderBookCollector) Subscribe(instruments []string) error {
	// Convert instruments to Derive format
	deriveInstruments := d.converter.ConvertInstrumentList(instruments, "derive")
	
	for _, instrument := range deriveInstruments {
		if err := d.subscribeInstrument(instrument); err != nil {
			return err
		}
		
		d.mu.Lock()
		d.subscribers[instrument] = true
		d.mu.Unlock()
	}
	
	return nil
}

func (d *DeriveWSOrderBookCollector) subscribeInstrument(instrument string) error {
	// Determine depth parameter based on configured depth
	depthParam := "10"
	if d.depth > 10 && d.depth <= 20 {
		depthParam = "20"
	} else if d.depth > 20 {
		depthParam = "100"
	}
	
	channel := fmt.Sprintf("orderbook.%s.1.%s", instrument, depthParam)
	
	// Use JSON-RPC format
	msg := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "subscribe",
		"params": map[string]interface{}{
			"channels": []string{channel},
		},
		"id": fmt.Sprintf("subscribe_%s_%d", instrument, time.Now().Unix()),
	}
	
	if err := d.conn.WriteJSON(msg); err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", instrument, err)
	}
	
	log.Printf("[Derive WS OrderBook] Subscribed to %s", channel)
	return nil
}

func (d *DeriveWSOrderBookCollector) Unsubscribe(instruments []string) error {
	// Convert instruments to Derive format
	deriveInstruments := d.converter.ConvertInstrumentList(instruments, "derive")
	
	for _, instrument := range deriveInstruments {
		depthParam := "10"
		if d.depth > 10 && d.depth <= 20 {
			depthParam = "20"
		} else if d.depth > 20 {
			depthParam = "100"
		}
		
		channel := fmt.Sprintf("orderbook.%s.1.%s", instrument, depthParam)
		
		msg := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "unsubscribe",
			"params": map[string]interface{}{
				"channels": []string{channel},
			},
			"id": fmt.Sprintf("unsubscribe_%s_%d", instrument, time.Now().Unix()),
		}
		
		if err := d.conn.WriteJSON(msg); err != nil {
			return fmt.Errorf("failed to unsubscribe from %s: %w", instrument, err)
		}
		
		d.mu.Lock()
		delete(d.subscribers, instrument)
		delete(d.orderBooks, instrument)
		d.mu.Unlock()
	}
	
	return nil
}

func (d *DeriveWSOrderBookCollector) GetOrderBooks() []OrderBookMetric {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	orderBooks := make([]OrderBookMetric, 0, len(d.orderBooks))
	for _, ob := range d.orderBooks {
		// Only return recent order books (less than 10 seconds old)
		if time.Since(ob.Timestamp) < 10*time.Second {
			orderBooks = append(orderBooks, *ob)
		}
	}
	
	return orderBooks
}

func (d *DeriveWSOrderBookCollector) Close() error {
	d.cancel()
	
	if d.conn != nil {
		// Send unsubscribe for all instruments
		d.mu.RLock()
		instruments := make([]string, 0, len(d.subscribers))
		for inst := range d.subscribers {
			instruments = append(instruments, inst)
		}
		d.mu.RUnlock()
		
		d.Unsubscribe(instruments)
		
		// Close connection
		d.conn.Close()
	}
	
	return nil
}

func (d *DeriveWSOrderBookCollector) Name() string {
	return "derive-ws-orderbook"
}