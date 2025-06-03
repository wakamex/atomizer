package marketmaker

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wakamex/atomizer/internal/types"
)

// MarketMaker manages automated market making
type MarketMaker struct {
	config   *types.MarketMakerConfig
	exchange types.MarketMakerExchange

	// Order tracking
	activeOrders       map[string]*types.MarketMakerOrder
	ordersByInstrument map[string]map[string]*types.MarketMakerOrder

	// Market data
	latestTickers map[string]*types.TickerUpdate

	// Position tracking
	positions map[string]decimal.Decimal

	// Statistics
	stats types.MarketMakerStats

	// Control
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	mu     sync.RWMutex

	// Error suppression
	orderbookErrorLogged map[string]bool

	// Per-instrument update locks
	updateLocks map[string]*sync.Mutex

	// Track failed cancel attempts
	failedCancelAttempts map[string]int

	// Track last update time
	lastUpdateTime map[string]time.Time
}

// NewMarketMaker creates a new market maker instance
func NewMarketMaker(config *types.MarketMakerConfig, exchange types.MarketMakerExchange) *MarketMaker {
	ctx, cancel := context.WithCancel(context.Background())

	updateLocks := make(map[string]*sync.Mutex)
	for _, instrument := range config.Instruments {
		updateLocks[instrument] = &sync.Mutex{}
	}

	return &MarketMaker{
		config:               config,
		exchange:             exchange,
		activeOrders:         make(map[string]*types.MarketMakerOrder),
		ordersByInstrument:   make(map[string]map[string]*types.MarketMakerOrder),
		latestTickers:        make(map[string]*types.TickerUpdate),
		positions:            make(map[string]decimal.Decimal),
		stats:                types.MarketMakerStats{BidAskSpread: make(map[string]decimal.Decimal)},
		ctx:                  ctx,
		cancel:               cancel,
		orderbookErrorLogged: make(map[string]bool),
		updateLocks:          updateLocks,
		failedCancelAttempts: make(map[string]int),
		lastUpdateTime:       make(map[string]time.Time),
	}
}

// Start begins market making
func (mm *MarketMaker) Start() error {
	mode := "two-sided (1 buy + 1 sell)"
	if mm.config.BidOnly {
		mode = "bid-only"
	} else if mm.config.AskOnly {
		mode = "ask-only"
	}
	log.Printf("Starting market maker: %d instruments, %s per instrument",
		len(mm.config.Instruments), mode)

	// Clear stale state
	mm.mu.Lock()
	mm.activeOrders = make(map[string]*types.MarketMakerOrder)
	mm.ordersByInstrument = make(map[string]map[string]*types.MarketMakerOrder)
	mm.mu.Unlock()

	// Load existing positions
	if err := mm.LoadPositions(); err != nil {
		return fmt.Errorf("failed to load positions: %w", err)
	}

	// Load existing orders
	if err := mm.LoadActiveOrders(); err != nil {
		return fmt.Errorf("failed to load active orders: %w", err)
	}

	// Cancel all existing orders on startup
	log.Printf("Cancelling all existing orders on startup...")
	mm.CancelAllOrdersOnStartup()

	// Subscribe to ticker updates
	tickerChan, err := mm.exchange.SubscribeTickers(mm.ctx, mm.config.Instruments)
	if err != nil {
		return fmt.Errorf("failed to subscribe to tickers: %w", err)
	}

	// Subscribe to orderbook updates
	mm.subscribeToOrderBooks()

	// Start goroutines
	mm.wg.Add(1)
	go mm.processTickers(tickerChan)

	mm.wg.Add(1)
	go mm.quoteUpdater()

	mm.wg.Add(1)
	go mm.statsReporter()

	// Initial reconciliation
	mm.ReconcileOrders()

	mm.stats.UptimeSeconds = 0
	log.Println("Market maker started successfully")

	return nil
}

// Stop gracefully shuts down the market maker
func (mm *MarketMaker) Stop() error {
	log.Println("Stopping market maker...")

	mm.cancel()
	mm.CancelAllOrders()
	mm.wg.Wait()

	log.Println("Market maker stopped")
	return nil
}

// processTickers handles incoming ticker updates
func (mm *MarketMaker) processTickers(tickerChan <-chan types.TickerUpdate) {
	defer mm.wg.Done()

	for {
		select {
		case <-mm.ctx.Done():
			return
		case ticker, ok := <-tickerChan:
			if !ok {
				log.Println("Ticker channel closed")
				return
			}

			mm.mu.Lock()
			mm.latestTickers[ticker.Instrument] = &ticker
			mm.mu.Unlock()

			if mm.shouldUpdateQuotes(ticker.Instrument) {
				mm.UpdateQuotesForInstrument(ticker.Instrument)
			}
		}
	}
}

// quoteUpdater periodically updates all quotes
func (mm *MarketMaker) quoteUpdater() {
	defer mm.wg.Done()

	ticker := time.NewTicker(mm.config.RefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-mm.ctx.Done():
			return
		case <-ticker.C:
			mm.updateAllQuotes()
		}
	}
}

// updateAllQuotes updates quotes for all instruments
func (mm *MarketMaker) updateAllQuotes() {
	for _, instrument := range mm.config.Instruments {
		if err := mm.UpdateQuotesForInstrument(instrument); err != nil {
			log.Printf("Failed to update quotes for %s: %v", instrument, err)
		}
	}
}

// Helper function for orderbook subscription
func (mm *MarketMaker) subscribeToOrderBooks() {
	if subscriber, ok := mm.exchange.(interface{ SubscribeOrderBook(string) error }); ok {
		for _, instrument := range mm.config.Instruments {
			if err := subscriber.SubscribeOrderBook(instrument); err != nil {
				log.Printf("Failed to subscribe to orderbook for %s: %v", instrument, err)
			} else {
				log.Printf("Subscribed to orderbook for %s", instrument)
			}
		}
		time.Sleep(2 * time.Second)
	}
}

// Helper for startup cancellation
func (mm *MarketMaker) CancelAllOrdersOnStartup() {
	cancelCount := 0
	mm.mu.RLock()
	ordersCopy := make(map[string]*types.MarketMakerOrder)
	for id, order := range mm.activeOrders {
		ordersCopy[id] = order
	}
	mm.mu.RUnlock()

	for _, order := range ordersCopy {
		if err := mm.exchange.CancelOrder(order.OrderID); err != nil {
			log.Printf("Failed to cancel order %s: %v", order.OrderID, err)
		} else {
			cancelCount++
			mm.mu.Lock()
			delete(mm.activeOrders, order.OrderID)
			if mm.ordersByInstrument[order.Instrument] != nil {
				delete(mm.ordersByInstrument[order.Instrument], order.Side)
			}
			mm.mu.Unlock()
		}
	}
	log.Printf("Cancelled %d orders on startup", cancelCount)
}
