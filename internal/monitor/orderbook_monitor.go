package monitor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// OrderBookMonitor handles order book collection from multiple exchanges
type OrderBookMonitor struct {
	config            *Config
	deribitCollector  *DeribitOrderBookCollector
	deriveWSCollector *DeriveWSOrderBookCollector
	spotCollector     *DeriveSpotCollector
	storage           *OrderBookStorage
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
}

// NewOrderBookMonitor creates a new order book monitor
func NewOrderBookMonitor(config *Config, depth int) (*OrderBookMonitor, error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	m := &OrderBookMonitor{
		config:  config,
		storage: NewOrderBookStorage(config.VictoriaMetricsURL),
		ctx:     ctx,
		cancel:  cancel,
	}
	
	// Initialize collectors based on configured exchanges
	for _, exchange := range config.Exchanges {
		switch exchange {
		case "deribit":
			m.deribitCollector = NewDeribitOrderBookCollector(depth)
		case "derive":
			wsCollector, err := NewDeriveWSOrderBookCollector(depth)
			if err != nil {
				cancel()
				return nil, fmt.Errorf("failed to create Derive WebSocket collector: %w", err)
			}
			m.deriveWSCollector = wsCollector
		}
	}
	
	// Always create spot collector for ETH and BTC prices
	spotCollector, err := NewDeriveSpotCollector()
	if err != nil {
		log.Printf("Warning: Failed to create spot collector: %v", err)
		// Don't fail completely, just log the error
	} else {
		m.spotCollector = spotCollector
		// Subscribe to ETH and BTC spot feeds
		if err := spotCollector.Subscribe([]string{"ETH", "BTC"}); err != nil {
			log.Printf("Warning: Failed to subscribe to spot feeds: %v", err)
		}
		
		// Set spot collector on Deribit collector for USD display
		if m.deribitCollector != nil {
			m.deribitCollector.SetSpotCollector(spotCollector)
		}
	}
	
	return m, nil
}

func (m *OrderBookMonitor) Start() error {
	// Subscribe to instruments on Derive WebSocket
	if m.deriveWSCollector != nil {
		if err := m.deriveWSCollector.Subscribe(m.config.InstrumentPatterns); err != nil {
			return fmt.Errorf("failed to subscribe to Derive instruments: %w", err)
		}
	}
	
	// Start collection loops
	if m.deribitCollector != nil {
		m.wg.Add(1)
		go m.deribitCollectionLoop()
	}
	
	if m.deriveWSCollector != nil {
		m.wg.Add(1)
		go m.deriveWSCollectionLoop()
	}
	
	// Start spot price collection loop
	if m.spotCollector != nil {
		log.Println("Starting spot price collection loop")
		m.wg.Add(1)
		go m.spotCollectionLoop()
	} else {
		log.Println("WARNING: Spot collector is nil, not starting spot collection loop")
	}
	
	return nil
}

func (m *OrderBookMonitor) Stop() error {
	m.cancel()
	m.wg.Wait()
	
	// Close WebSocket connections
	if m.deriveWSCollector != nil {
		if err := m.deriveWSCollector.Close(); err != nil {
			log.Printf("Error closing Derive WebSocket: %v", err)
		}
	}
	
	if m.spotCollector != nil {
		if err := m.spotCollector.Close(); err != nil {
			log.Printf("Error closing spot collector: %v", err)
		}
	}
	
	return nil
}

func (m *OrderBookMonitor) deribitCollectionLoop() {
	defer m.wg.Done()
	
	ticker := time.NewTicker(m.config.Interval)
	defer ticker.Stop()
	
	log.Printf("Starting Deribit order book collection every %v", m.config.Interval)
	
	// Initial collection
	m.collectDeribitOrderBooks()
	
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.collectDeribitOrderBooks()
		}
	}
}

func (m *OrderBookMonitor) collectDeribitOrderBooks() {
	orderBooks, err := m.deribitCollector.CollectOrderBooks(m.ctx, m.config.InstrumentPatterns)
	if err != nil {
		log.Printf("Deribit order book collection error: %v", err)
		return
	}
	
	if len(orderBooks) > 0 {
		if err := m.storage.WriteOrderBooks(orderBooks); err != nil {
			log.Printf("Storage error: %v", err)
			return
		}
		
		log.Printf("Collected %d order books from Deribit", len(orderBooks))
		
		// Log sample order book info
		if len(orderBooks) > 0 {
			ob := orderBooks[0]
			if len(ob.Bids) > 0 && len(ob.Asks) > 0 {
				spread := ob.Asks[0].Price - ob.Bids[0].Price
				
				// For Deribit ETH options, convert prices to USD for display
				displayBidPrice := ob.Bids[0].Price
				displayAskPrice := ob.Asks[0].Price
				displaySpread := spread
				priceUnit := "ETH"
				
				if m.spotCollector != nil && strings.Contains(ob.Instrument, "ETH") && ob.Exchange == "deribit" {
					if ethSpot, ok := m.spotCollector.GetSpotPrice("ETH"); ok {
						// Convert ETH prices to USD for display only
						displayBidPrice = ob.Bids[0].Price * ethSpot
						displayAskPrice = ob.Asks[0].Price * ethSpot
						displaySpread = spread * ethSpot
						priceUnit = fmt.Sprintf("USD (ETH=$%.2f)", ethSpot)
					}
				}
				
				log.Printf("  %s: Bid %.2f x %.2f, Ask %.2f x %.2f, Spread %.2f (%.3f%%) %s",
					ob.Instrument,
					displayBidPrice, ob.Bids[0].Size,
					displayAskPrice, ob.Asks[0].Size,
					displaySpread, (spread/ob.Bids[0].Price)*100,
					priceUnit)
			}
		}
	}
}

func (m *OrderBookMonitor) deriveWSCollectionLoop() {
	defer m.wg.Done()
	
	ticker := time.NewTicker(m.config.Interval)
	defer ticker.Stop()
	
	log.Printf("Starting Derive WebSocket order book collection every %v", m.config.Interval)
	
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.collectDeriveOrderBooks()
		}
	}
}

func (m *OrderBookMonitor) spotCollectionLoop() {
	defer m.wg.Done()
	
	// Use a faster interval for spot prices as they change frequently
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	log.Printf("Starting spot price collection every 10s")
	
	// Initial collection
	m.collectSpotPrices()
	
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.collectSpotPrices()
		}
	}
}

func (m *OrderBookMonitor) collectSpotPrices() {
	spotPrices := m.spotCollector.GetAllSpotPrices()
	
	if len(spotPrices) == 0 {
		log.Println("No spot prices available yet")
		return
	}
	
	// Convert spot prices to metrics
	metrics := make([]Metric, 0, len(spotPrices))
	for currency, spot := range spotPrices {
		// Create a metric for the spot price
		metric := Metric{
			Exchange:   "derive",
			Instrument: fmt.Sprintf("%s-SPOT", currency),
			Timestamp:  spot.Timestamp,
			BidPrice:   spot.Price,
			AskPrice:   spot.Price,
			LastPrice:  spot.Price,
		}
		metrics = append(metrics, metric)
	}
	
	// Write to storage using the regular VMStorage
	storage := NewVMStorage(m.config.VictoriaMetricsURL)
	if err := storage.Write(metrics); err != nil {
		log.Printf("Failed to write spot prices: %v", err)
		return
	}
	
	log.Printf("Collected spot prices: %v", spotPrices)
}

func (m *OrderBookMonitor) collectDeriveOrderBooks() {
	// Get current order books from WebSocket collector
	orderBooks := m.deriveWSCollector.GetOrderBooks()
	
	if len(orderBooks) > 0 {
		if err := m.storage.WriteOrderBooks(orderBooks); err != nil {
			log.Printf("Storage error: %v", err)
			return
		}
		
		log.Printf("Collected %d order books from Derive WebSocket", len(orderBooks))
		
		// Log sample order book info
		if len(orderBooks) > 0 {
			ob := orderBooks[0]
			if len(ob.Bids) > 0 && len(ob.Asks) > 0 {
				spread := ob.Asks[0].Price - ob.Bids[0].Price
				// Check if this is an ETH option and show spot price
				spotInfo := ""
				if m.spotCollector != nil && strings.Contains(ob.Instrument, "ETH") {
					if ethSpot, ok := m.spotCollector.GetSpotPrice("ETH"); ok {
						spotInfo = fmt.Sprintf(" (ETH Spot: $%.2f)", ethSpot)
					}
				}
				
				log.Printf("  %s: Bid %.4f x %.2f, Ask %.4f x %.2f, Spread %.4f (%.3f%%)%s",
					ob.Instrument,
					ob.Bids[0].Price, ob.Bids[0].Size,
					ob.Asks[0].Price, ob.Asks[0].Size,
					spread, (spread/ob.Bids[0].Price)*100,
					spotInfo)
			}
		}
	}
}