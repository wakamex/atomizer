# Market Maker Refactored File Structure

## File Structure
```
market-maker/
├── market_maker.go      # Core struct and main Start/Stop methods
├── orders.go           # Order management (place, cancel, track)
├── quotes.go           # Quote calculation and updates
├── positions.go        # Position and risk management
├── reconciliation.go   # Order reconciliation and cleanup
├── stats.go           # Statistics and reporting
└── types.go           # Shared types and interfaces
```

## 1. types.go - Shared Types and Interfaces

```go
package main

import (
    "context"
    "time"
    "github.com/shopspring/decimal"
)

// MarketMakerExchange interface - your exchange abstraction
type MarketMakerExchange interface {
    PlaceLimitOrder(instrument, side string, price, amount decimal.Decimal) (string, error)
    CancelOrder(orderID string) error
    ReplaceOrder(orderID, instrument, side string, price, amount decimal.Decimal) (string, error)
    GetOpenOrders() ([]MarketMakerOrder, error)
    GetOrderBook(instrument string) (*MarketMakerOrderBook, error)
    GetPositions() ([]Position, error)
    SubscribeTickers(ctx context.Context, instruments []string) (<-chan TickerUpdate, error)
}

// MarketMakerOrder represents an order
type MarketMakerOrder struct {
    OrderID    string
    Instrument string
    Side       string
    Price      decimal.Decimal
    Amount     decimal.Decimal
    Status     string
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

// MarketMakerOrderBook represents order book data
type MarketMakerOrderBook struct {
    Bids []OrderBookLevel
    Asks []OrderBookLevel
}

type OrderBookLevel struct {
    Price decimal.Decimal
    Size  decimal.Decimal
}

// TickerUpdate represents market data update
type TickerUpdate struct {
    Instrument string
    BestBid    decimal.Decimal
    BestAsk    decimal.Decimal
    MarkPrice  decimal.Decimal
    Timestamp  time.Time
}

// Position represents a position
type Position struct {
    InstrumentName string
    Direction      string
    Amount         float64
}

// MarketMakerConfig holds configuration
type MarketMakerConfig struct {
    Instruments              []string
    QuoteSize               decimal.Decimal
    SpreadBps               int
    MinSpreadBps            int
    Improvement             decimal.Decimal
    ImprovementReferenceSize decimal.Decimal
    RefreshInterval         time.Duration
    CancelThreshold         decimal.Decimal
    MaxPositionSize         decimal.Decimal
    MaxTotalExposure        decimal.Decimal
    BidOnly                 bool
    AskOnly                 bool
}

// MarketMakerStats tracks statistics
type MarketMakerStats struct {
    OrdersPlaced    int64
    OrdersCancelled int64
    OrdersFilled    int64
    BidAskSpread    map[string]decimal.Decimal
    UptimeSeconds   int64
    LastUpdate      time.Time
}
```

## 2. market_maker.go - Core Structure

```go
package main

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"
    "github.com/shopspring/decimal"
)

// MarketMaker manages automated market making
type MarketMaker struct {
    config   *MarketMakerConfig
    exchange MarketMakerExchange
    
    // Order tracking
    activeOrders       map[string]*MarketMakerOrder
    ordersByInstrument map[string]map[string]*MarketMakerOrder
    
    // Market data
    latestTickers map[string]*TickerUpdate
    
    // Position tracking
    positions map[string]decimal.Decimal
    
    // Statistics
    stats MarketMakerStats
    
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
func NewMarketMaker(config *MarketMakerConfig, exchange MarketMakerExchange) *MarketMaker {
    ctx, cancel := context.WithCancel(context.Background())
    
    updateLocks := make(map[string]*sync.Mutex)
    for _, instrument := range config.Instruments {
        updateLocks[instrument] = &sync.Mutex{}
    }
    
    return &MarketMaker{
        config:               config,
        exchange:            exchange,
        activeOrders:        make(map[string]*MarketMakerOrder),
        ordersByInstrument:  make(map[string]map[string]*MarketMakerOrder),
        latestTickers:       make(map[string]*TickerUpdate),
        positions:           make(map[string]decimal.Decimal),
        stats:               MarketMakerStats{BidAskSpread: make(map[string]decimal.Decimal)},
        ctx:                 ctx,
        cancel:              cancel,
        orderbookErrorLogged: make(map[string]bool),
        updateLocks:         updateLocks,
        failedCancelAttempts: make(map[string]int),
        lastUpdateTime:      make(map[string]time.Time),
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
    mm.activeOrders = make(map[string]*MarketMakerOrder)
    mm.ordersByInstrument = make(map[string]map[string]*MarketMakerOrder)
    mm.mu.Unlock()
    
    // Load existing positions
    if err := mm.loadPositions(); err != nil {
        return fmt.Errorf("failed to load positions: %w", err)
    }
    
    // Load existing orders
    if err := mm.loadActiveOrders(); err != nil {
        return fmt.Errorf("failed to load active orders: %w", err)
    }
    
    // Cancel all existing orders on startup
    log.Printf("Cancelling all existing orders on startup...")
    mm.cancelAllOrdersOnStartup()
    
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
    mm.reconcileOrders()
    
    mm.stats.UptimeSeconds = 0
    log.Println("Market maker started successfully")
    
    return nil
}

// Stop gracefully shuts down the market maker
func (mm *MarketMaker) Stop() error {
    log.Println("Stopping market maker...")
    
    mm.cancel()
    mm.cancelAllOrders()
    mm.wg.Wait()
    
    log.Println("Market maker stopped")
    return nil
}

// processTickers handles incoming ticker updates
func (mm *MarketMaker) processTickers(tickerChan <-chan TickerUpdate) {
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
                mm.updateQuotesForInstrument(ticker.Instrument)
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
        if err := mm.updateQuotesForInstrument(instrument); err != nil {
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
func (mm *MarketMaker) cancelAllOrdersOnStartup() {
    cancelCount := 0
    mm.mu.RLock()
    ordersCopy := make(map[string]*MarketMakerOrder)
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
```

## 3. quotes.go - Quote Calculation and Updates

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/shopspring/decimal"
)

// updateQuotesForInstrument updates quotes for a specific instrument
func (mm *MarketMaker) updateQuotesForInstrument(instrument string) error {
    // Prevent concurrent updates
    if lock, exists := mm.updateLocks[instrument]; exists {
        lock.Lock()
        defer lock.Unlock()
    }
    
    // Rate limiting check
    mm.mu.RLock()
    lastUpdate, exists := mm.lastUpdateTime[instrument]
    mm.mu.RUnlock()
    
    if exists && time.Since(lastUpdate) < 2*time.Second {
        return nil
    }
    
    mm.mu.Lock()
    mm.lastUpdateTime[instrument] = time.Now()
    mm.mu.Unlock()
    
    // Get ticker data
    mm.mu.RLock()
    ticker, exists := mm.latestTickers[instrument]
    mm.mu.RUnlock()
    
    if !exists || ticker == nil {
        return fmt.Errorf("no ticker data for %s", instrument)
    }
    
    // Skip if no valid price data
    if ticker.BestBid.IsZero() && ticker.BestAsk.IsZero() && ticker.MarkPrice.IsZero() {
        log.Printf("No valid price data for %s yet, skipping quote update", instrument)
        return nil
    }
    
    // Fetch orderbook if needed
    var orderBook *MarketMakerOrderBook
    if mm.config.ImprovementReferenceSize.GreaterThan(decimal.Zero) {
        var err error
        orderBook, err = mm.exchange.GetOrderBook(instrument)
        if err != nil {
            mm.logOrderbookError(instrument, err)
        } else {
            mm.clearOrderbookError(instrument)
        }
    }
    
    // Calculate quotes
    bidPrice, askPrice := mm.calculateQuotes(ticker, orderBook)
    
    // Check risk limits
    if !mm.checkRiskLimits(instrument, mm.config.QuoteSize) {
        log.Printf("Risk limits exceeded for %s, skipping quote update", instrument)
        return nil
    }
    
    // Update orders
    return mm.updateOrCreateOrders(instrument, bidPrice, askPrice)
}

// calculateQuotes calculates bid and ask prices based on current market
func (mm *MarketMaker) calculateQuotes(ticker *TickerUpdate, orderBook *MarketMakerOrderBook) (bidPrice, askPrice decimal.Decimal) {
    // Calculate mid price
    var midPrice decimal.Decimal
    if ticker.BestBid.IsZero() || ticker.BestAsk.IsZero() {
        if !ticker.MarkPrice.IsZero() {
            midPrice = ticker.MarkPrice
        } else {
            log.Printf("WARNING: No valid price data for %s", ticker.Instrument)
            midPrice = decimal.NewFromFloat(1.0)
        }
    } else {
        midPrice = ticker.BestBid.Add(ticker.BestAsk).Div(decimal.NewFromInt(2))
    }
    
    // Determine reference prices
    referenceBid := ticker.BestBid
    referenceAsk := ticker.BestAsk
    
    // Handle zero prices
    if referenceBid.IsZero() || referenceAsk.IsZero() {
        spreadAmount := midPrice.Mul(decimal.NewFromInt(int64(mm.config.SpreadBps)).Div(decimal.NewFromInt(10000)))
        if referenceBid.IsZero() {
            referenceBid = midPrice.Sub(spreadAmount.Div(decimal.NewFromInt(2)))
        }
        if referenceAsk.IsZero() {
            referenceAsk = midPrice.Add(spreadAmount.Div(decimal.NewFromInt(2)))
        }
    }
    
    // Use orderbook if reference size is set
    if orderBook != nil && mm.config.ImprovementReferenceSize.GreaterThan(decimal.Zero) {
        mm.adjustPricesForReferenceSize(orderBook, &referenceBid, &referenceAsk, midPrice)
    }
    
    // Calculate our quotes with improvement
    bidPrice = referenceBid.Add(mm.config.Improvement)
    askPrice = referenceAsk.Sub(mm.config.Improvement)
    
    // Ensure minimum spread
    minSpread := midPrice.Mul(decimal.NewFromInt(int64(mm.config.MinSpreadBps)).Div(decimal.NewFromInt(10000)))
    if askPrice.Sub(bidPrice).LessThan(minSpread) {
        bidPrice = midPrice.Sub(minSpread.Div(decimal.NewFromInt(2)))
        askPrice = midPrice.Add(minSpread.Div(decimal.NewFromInt(2)))
    }
    
    return bidPrice, askPrice
}

// adjustPricesForReferenceSize finds best bid/ask with sufficient size
func (mm *MarketMaker) adjustPricesForReferenceSize(orderBook *MarketMakerOrderBook, referenceBid, referenceAsk *decimal.Decimal, midPrice decimal.Decimal) {
    foundBid := false
    for _, bid := range orderBook.Bids {
        if bid.Size.GreaterThanOrEqual(mm.config.ImprovementReferenceSize) {
            *referenceBid = bid.Price
            foundBid = true
            break
        }
    }
    
    foundAsk := false
    for _, ask := range orderBook.Asks {
        if ask.Size.GreaterThanOrEqual(mm.config.ImprovementReferenceSize) {
            *referenceAsk = ask.Price
            foundAsk = true
            break
        }
    }
    
    // Fallback if insufficient size found
    if !foundBid || !foundAsk {
        spreadAmount := midPrice.Mul(decimal.NewFromInt(int64(mm.config.SpreadBps)).Div(decimal.NewFromInt(10000)))
        if !foundBid {
            *referenceBid = midPrice.Sub(spreadAmount.Div(decimal.NewFromInt(2)))
        }
        if !foundAsk {
            *referenceAsk = midPrice.Add(spreadAmount.Div(decimal.NewFromInt(2)))
        }
    }
}

// shouldUpdateQuotes checks if quotes need updating
func (mm *MarketMaker) shouldUpdateQuotes(instrument string) bool {
    mm.mu.RLock()
    defer mm.mu.RUnlock()
    
    orders, exists := mm.ordersByInstrument[instrument]
    if !exists || len(orders) == 0 {
        return true
    }
    
    ticker, exists := mm.latestTickers[instrument]
    if !exists {
        return false
    }
    
    // Check if market has moved significantly
    for side, order := range orders {
        var marketPrice decimal.Decimal
        if side == "buy" {
            marketPrice = ticker.BestBid
        } else {
            marketPrice = ticker.BestAsk
        }
        
        priceDiff := order.Price.Sub(marketPrice).Abs()
        if priceDiff.GreaterThan(order.Price.Mul(mm.config.CancelThreshold)) {
            return true
        }
    }
    
    return false
}

// Helper functions for error logging
func (mm *MarketMaker) logOrderbookError(instrument string, err error) {
    mm.mu.Lock()
    if !mm.orderbookErrorLogged[instrument] {
        log.Printf("Failed to fetch orderbook for %s: %v, using ticker data", instrument, err)
        mm.orderbookErrorLogged[instrument] = true
    }
    mm.mu.Unlock()
}

func (mm *MarketMaker) clearOrderbookError(instrument string) {
    mm.mu.Lock()
    delete(mm.orderbookErrorLogged, instrument)
    mm.mu.Unlock()
}
```

## 4. orders.go - Order Management

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "sync"
    "time"
    "github.com/shopspring/decimal"
)

// updateOrCreateOrders handles the order update/creation logic
func (mm *MarketMaker) updateOrCreateOrders(instrument string, bidPrice, askPrice decimal.Decimal) error {
    // Get fresh order state
    openOrders, err := mm.exchange.GetOpenOrders()
    if err != nil {
        log.Printf("Failed to get open orders: %v", err)
        return err
    }
    
    // Update tracking with real orders
    mm.syncOrderTracking(instrument, openOrders)
    
    // Get existing orders
    mm.mu.RLock()
    existingOrders := mm.ordersByInstrument[instrument]
    mm.mu.RUnlock()
    
    if existingOrders != nil && len(existingOrders) > 0 {
        return mm.replaceExistingOrders(instrument, bidPrice, askPrice, existingOrders)
    }
    
    // No existing orders, place new ones
    return mm.placeQuotes(instrument, bidPrice, askPrice)
}

// syncOrderTracking updates our tracking with real orders from exchange
func (mm *MarketMaker) syncOrderTracking(instrument string, openOrders []MarketMakerOrder) {
    mm.mu.Lock()
    defer mm.mu.Unlock()
    
    // Clear existing tracking for this instrument
    if mm.ordersByInstrument[instrument] != nil {
        for _, order := range mm.ordersByInstrument[instrument] {
            if order != nil {
                delete(mm.activeOrders, order.OrderID)
            }
        }
    }
    mm.ordersByInstrument[instrument] = make(map[string]*MarketMakerOrder)
    
    // Add real orders to tracking
    for i := range openOrders {
        order := &openOrders[i]
        if order.Instrument == instrument {
            mm.activeOrders[order.OrderID] = order
            mm.ordersByInstrument[instrument][order.Side] = order
        }
    }
}

// replaceExistingOrders handles order replacement logic
func (mm *MarketMaker) replaceExistingOrders(instrument string, bidPrice, askPrice decimal.Decimal, existingOrders map[string]*MarketMakerOrder) error {
    var bidOrder, askOrder *MarketMakerOrder
    
    // Find existing orders
    for side, order := range existingOrders {
        if side == "buy" {
            bidOrder = order
        } else if side == "sell" {
            askOrder = order
        }
    }
    
    // Cancel unwanted orders in one-sided mode
    if mm.config.AskOnly && bidOrder != nil {
        log.Printf("Cancelling bid order %s (ask-only mode)", bidOrder.OrderID)
        mm.cancelOrder(bidOrder.OrderID)
        bidOrder = nil
    }
    if mm.config.BidOnly && askOrder != nil {
        log.Printf("Cancelling ask order %s (bid-only mode)", askOrder.OrderID)
        mm.cancelOrder(askOrder.OrderID)
        askOrder = nil
    }
    
    var replacedBid, replacedAsk bool
    
    // Skip sides based on mode
    if mm.config.AskOnly {
        replacedBid = true
    }
    if mm.config.BidOnly {
        replacedAsk = true
    }
    
    // Handle bid order
    if bidOrder != nil && !mm.config.AskOnly {
        replacedBid = mm.handleOrderUpdate(bidOrder, instrument, "buy", bidPrice)
    }
    
    // Handle ask order
    if askOrder != nil && !mm.config.BidOnly {
        replacedAsk = mm.handleOrderUpdate(askOrder, instrument, "sell", askPrice)
    }
    
    // Place new orders for any that weren't replaced
    if !replacedBid && !mm.config.AskOnly {
        if err := mm.placeSingleQuote(instrument, "buy", bidPrice); err != nil {
            log.Printf("Failed to place bid order: %v", err)
        }
    }
    if !replacedAsk && !mm.config.BidOnly {
        if err := mm.placeSingleQuote(instrument, "sell", askPrice); err != nil {
            log.Printf("Failed to place ask order: %v", err)
        }
    }
    
    return nil
}

// handleOrderUpdate processes a single order update
func (mm *MarketMaker) handleOrderUpdate(order *MarketMakerOrder, instrument, side string, newPrice decimal.Decimal) bool {
    // Skip if order is already at target price and recent
    if order.Price.Equal(newPrice) && time.Since(order.CreatedAt) < 5*time.Second {
        debugLog("Skipping %s update - order %s already at target price", side, order.OrderID)
        return true
    }
    
    // Just cancel - we'll recreate
    if !mm.cancelOrder(order.OrderID) {
        log.Printf("Failed to cancel %s order %s, skipping recreation to avoid duplicates", side, order.OrderID)
        return true
    }
    
    return false
}

// placeQuotes places bid and ask orders concurrently
func (mm *MarketMaker) placeQuotes(instrument string, bidPrice, askPrice decimal.Decimal) error {
    var bidOrderID, askOrderID string
    var bidErr, askErr error
    ordersPlaced := 0
    
    var wg sync.WaitGroup
    
    // Place bid order
    if !mm.config.AskOnly {
        wg.Add(1)
        go func() {
            defer wg.Done()
            bidOrderID, bidErr = mm.exchange.PlaceLimitOrder(instrument, "buy", bidPrice, mm.config.QuoteSize)
        }()
    }
    
    // Place ask order
    if !mm.config.BidOnly {
        wg.Add(1)
        go func() {
            defer wg.Done()
            askOrderID, askErr = mm.exchange.PlaceLimitOrder(instrument, "sell", askPrice, mm.config.QuoteSize)
        }()
    }
    
    wg.Wait()
    
    // Handle errors
    if bidErr != nil && askErr != nil {
        return fmt.Errorf("failed to place both orders: bid error: %w, ask error: %v", bidErr, askErr)
    }
    
    if bidErr != nil && askOrderID != "" {
        mm.exchange.CancelOrder(askOrderID)
        return fmt.Errorf("failed to place bid order (cancelled ask): %w", bidErr)
    }
    
    if askErr != nil && bidOrderID != "" {
        mm.exchange.CancelOrder(bidOrderID)
        return fmt.Errorf("failed to place ask order (cancelled bid): %w", askErr)
    }
    
    // Track successfully placed orders
    mm.mu.Lock()
    if bidOrderID != "" && bidErr == nil {
        mm.trackOrder(bidOrderID, instrument, "buy", bidPrice, mm.config.QuoteSize)
        ordersPlaced++
    }
    if askOrderID != "" && askErr == nil {
        mm.trackOrder(askOrderID, instrument, "sell", askPrice, mm.config.QuoteSize)
        ordersPlaced++
    }
    mm.stats.OrdersPlaced += int64(ordersPlaced)
    mm.mu.Unlock()
    
    // Log placement
    if mm.config.BidOnly {
        log.Printf("Placed bid order for %s @ %s", instrument, bidPrice.String())
    } else if mm.config.AskOnly {
        log.Printf("Placed ask order for %s @ %s", instrument, askPrice.String())
    } else {
        log.Printf("Placed quotes for %s: Bid %s, Ask %s", instrument, bidPrice.String(), askPrice.String())
    }
    
    return nil
}

// placeSingleQuote places a single buy or sell order
func (mm *MarketMaker) placeSingleQuote(instrument, side string, price decimal.Decimal) error {
    orderID, err := mm.exchange.PlaceLimitOrder(instrument, side, price, mm.config.QuoteSize)
    if err != nil {
        return err
    }
    
    mm.mu.Lock()
    mm.trackOrder(orderID, instrument, side, price, mm.config.QuoteSize)
    mm.stats.OrdersPlaced++
    mm.mu.Unlock()
    
    log.Printf("Placed %s order %s for %s @ %s", side, orderID, instrument, price)
    return nil
}

// trackOrder adds an order to tracking
func (mm *MarketMaker) trackOrder(orderID, instrument, side string, price, amount decimal.Decimal) {
    order := &MarketMakerOrder{
        OrderID:    orderID,
        Instrument: instrument,
        Side:       side,
        Price:      price,
        Amount:     amount,
        Status:     "open",
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }
    
    mm.activeOrders[orderID] = order
    
    if mm.ordersByInstrument[instrument] == nil {
        mm.ordersByInstrument[instrument] = make(map[string]*MarketMakerOrder)
    }
    mm.ordersByInstrument[instrument][side] = order
}

// cancelOrder cancels a single order
func (mm *MarketMaker) cancelOrder(orderID string) bool {
    // Check failed attempts
    mm.mu.Lock()
    if attempts, exists := mm.failedCancelAttempts[orderID]; exists && attempts >= 3 {
        log.Printf("Order %s has failed to cancel %d times, treating as non-existent", orderID, attempts)
        mm.removeOrderFromTracking(orderID)
        delete(mm.failedCancelAttempts, orderID)
        mm.mu.Unlock()
        return true
    }
    mm.mu.Unlock()
    
    // Try to cancel with retries
    var lastErr error
    for retries := 0; retries < 3; retries++ {
        if retries > 0 {
            time.Sleep(500 * time.Millisecond)
        }
        
        err := mm.exchange.CancelOrder(orderID)
        if err == nil {
            mm.mu.Lock()
            mm.removeOrderFromTracking(orderID)
            mm.stats.OrdersCancelled++
            delete(mm.failedCancelAttempts, orderID)
            mm.mu.Unlock()
            return true
        }
        
        lastErr = err
        
        // Handle "order not found" as success
        if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "Order does not exist") {
            log.Printf("Order %s already cancelled or doesn't exist", orderID)
            mm.mu.Lock()
            mm.removeOrderFromTracking(orderID)
            delete(mm.failedCancelAttempts, orderID)
            mm.mu.Unlock()
            return true
        }
        
        // Retry internal errors
        if strings.Contains(err.Error(), "Internal error") {
            if retries == 0 {
                log.Printf("Internal error cancelling order %s, will retry", orderID)
            }
            continue
        }
        
        break
    }
    
    // Track failed attempts
    mm.mu.Lock()
    mm.failedCancelAttempts[orderID]++
    attempts := mm.failedCancelAttempts[orderID]
    mm.mu.Unlock()
    
    if attempts >= 3 {
        log.Printf("Order %s failed to cancel %d times, treating as cancelled", orderID, attempts)
        mm.mu.Lock()
        mm.removeOrderFromTracking(orderID)
        delete(mm.failedCancelAttempts, orderID)
        mm.mu.Unlock()
        return true
    }
    
    log.Printf("Failed to cancel order %s (attempt %d): %v", orderID, attempts, lastErr)
    return false
}

// removeOrderFromTracking removes an order from internal tracking (must be called with lock held)
func (mm *MarketMaker) removeOrderFromTracking(orderID string) {
    if order, exists := mm.activeOrders[orderID]; exists {
        delete(mm.activeOrders, orderID)
        if instrumentOrders, exists := mm.ordersByInstrument[order.Instrument]; exists {
            delete(instrumentOrders, order.Side)
            if len(instrumentOrders) == 0 {
                delete(mm.ordersByInstrument, order.Instrument)
            }
        }
    }
}

// cancelAllOrders cancels all active orders
func (mm *MarketMaker) cancelAllOrders() {
    mm.mu.Lock()
    defer mm.mu.Unlock()
    
    for orderID := range mm.activeOrders {
        if err := mm.exchange.CancelOrder(orderID); err != nil {
            log.Printf("Failed to cancel order %s: %v", orderID, err)
        } else {
            mm.stats.OrdersCancelled++
        }
    }
    
    mm.activeOrders = make(map[string]*MarketMakerOrder)
    mm.ordersByInstrument = make(map[string]map[string]*MarketMakerOrder)
}

// loadActiveOrders loads current open orders from exchange
func (mm *MarketMaker) loadActiveOrders() error {
    orders, err := mm.exchange.GetOpenOrders()
    if err != nil {
        return err
    }
    
    mm.mu.Lock()
    defer mm.mu.Unlock()
    
    for _, order := range orders {
        mm.activeOrders[order.OrderID] = &order
        
        if mm.ordersByInstrument[order.Instrument] == nil {
            mm.ordersByInstrument[order.Instrument] = make(map[string]*MarketMakerOrder)
        }
        mm.ordersByInstrument[order.Instrument][order.Side] = &order
    }
    
    log.Printf("Loaded %d active orders", len(orders))
    return nil
}

// verifyOrderExists checks if an order actually exists on the exchange
func (mm *MarketMaker) verifyOrderExists(orderID string) bool {
    orders, err := mm.exchange.GetOpenOrders()
    if err != nil {
        debugLog("Failed to verify order %s: %v", orderID, err)
        return true
    }
    
    for _, order := range orders {
        if order.OrderID == orderID {
            return true
        }
    }
    return false
}
```

## 5. positions.go - Position and Risk Management

```go
package main

import (
    "log"
    "github.com/shopspring/decimal"
)

// checkRiskLimits checks if placing an order would exceed risk limits
func (mm *MarketMaker) checkRiskLimits(instrument string, size decimal.Decimal) bool {
    mm.mu.RLock()
    defer mm.mu.RUnlock()
    
    // Check position limit for instrument
    currentPosition := mm.positions[instrument]
    if currentPosition.Add(size).Abs().GreaterThan(mm.config.MaxPositionSize) {
        return false
    }
    
    // Check total exposure
    totalExposure := decimal.Zero
    for _, pos := range mm.positions {
        totalExposure = totalExposure.Add(pos.Abs())
    }
    
    if totalExposure.Add(size).GreaterThan(mm.config.MaxTotalExposure) {
        return false
    }
    
    return true
}

// loadPositions loads current positions from exchange
func (mm *MarketMaker) loadPositions() error {
    positions, err := mm.exchange.GetPositions()
    if err != nil {
        return err
    }
    
    mm.mu.Lock()
    defer mm.mu.Unlock()
    
    for _, pos := range positions {
        amount := decimal.NewFromFloat(pos.Amount)
        if pos.Direction == "sell" {
            amount = amount.Neg()
        }
        mm.positions[pos.InstrumentName] = amount
    }
    
    log.Printf("Loaded %d positions", len(positions))
    return nil
}

// updatePosition updates position after a fill
func (mm *MarketMaker) updatePosition(instrument string, side string, amount decimal.Decimal) {
    mm.mu.Lock()
    defer mm.mu.Unlock()
    
    if side == "sell" {
        amount = amount.Neg()
    }
    
    currentPosition := mm.positions[instrument]
    mm.positions[instrument] = currentPosition.Add(amount)
    
    log.Printf("Updated position for %s: %s -> %s", 
        instrument, 
        currentPosition.String(), 
        mm.positions[instrument].String())
}

// getNetPosition returns the net position for an instrument
func (mm *MarketMaker) getNetPosition(instrument string) decimal.Decimal {
    mm.mu.RLock()
    defer mm.mu.RUnlock()
    
    return mm.positions[instrument]
}

// getTotalExposure returns total absolute exposure across all instruments
func (mm *MarketMaker) getTotalExposure() decimal.Decimal {
    mm.mu.RLock()
    defer mm.mu.RUnlock()
    
    totalExposure := decimal.Zero
    for _, pos := range mm.positions {
        totalExposure = totalExposure.Add(pos.Abs())
    }
    
    return totalExposure
}
```

## 6. reconciliation.go - Order Reconciliation

```go
package main

import (
    "log"
)

// reconcileOrders finds and cancels any orders not being tracked
func (mm *MarketMaker) reconcileOrders() {
    openOrders, err := mm.exchange.GetOpenOrders()
    if err != nil {
        log.Printf("Failed to get open orders for reconciliation: %v", err)
        return
    }
    
    mm.mu.RLock()
    trackedOrders := make(map[string]bool)
    for orderID := range mm.activeOrders {
        trackedOrders[orderID] = true
    }
    mm.mu.RUnlock()
    
    // Find orphaned orders
    orphanedCount := 0
    for _, order := range openOrders {
        if !trackedOrders[order.OrderID] {
            orphanedCount++
            log.Printf("Found orphaned order %s for %s, cancelling", order.OrderID, order.Instrument)
            if err := mm.exchange.CancelOrder(order.OrderID); err != nil {
                log.Printf("Failed to cancel orphaned order %s: %v", order.OrderID, err)
            }
        }
    }
    
    if orphanedCount > 0 {
        log.Printf("Cancelled %d orphaned orders", orphanedCount)
        mm.mu.Lock()
        mm.stats.OrdersCancelled += int64(orphanedCount)
        mm.mu.Unlock()
    }
    
    // Verify tracked orders still exist
    mm.verifyTrackedOrders(openOrders)
}

// verifyTrackedOrders removes tracked orders that no longer exist
func (mm *MarketMaker) verifyTrackedOrders(openOrders []MarketMakerOrder) {
    mm.mu.Lock()
    defer mm.mu.Unlock()
    
    for orderID, order := range mm.activeOrders {
        found := false
        for _, openOrder := range openOrders {
            if openOrder.OrderID == orderID {
                found = true
                break
            }
        }
        if !found {
            log.Printf("Tracked order %s no longer exists on exchange, removing from tracking", orderID)
            delete(mm.activeOrders, orderID)
            if orders, ok := mm.ordersByInstrument[order.Instrument]; ok {
                delete(orders, order.Side)
                if len(orders) == 0 {
                    delete(mm.ordersByInstrument, order.Instrument)
                }
            }
        }
    }
}

// reconcileOrdersForInstrument reconciles orders for a specific instrument
func (mm *MarketMaker) reconcileOrdersForInstrument(instrument string) {
    openOrders, err := mm.exchange.GetOpenOrders()
    if err != nil {
        debugLog("Failed to get open orders for reconciliation: %v", err)
        return
    }
    
    // Filter to just this instrument
    var instrumentOrders []MarketMakerOrder
    for _, order := range openOrders {
        if order.Instrument == instrument {
            instrumentOrders = append(instrumentOrders, order)
        }
    }
    
    mm.mu.Lock()
    defer mm.mu.Unlock()
    
    trackedOrders := mm.ordersByInstrument[instrument]
    if trackedOrders == nil {
        trackedOrders = make(map[string]*MarketMakerOrder)
    }
    
    // Build map of actual order IDs
    actualOrders := make(map[string]bool)
    for _, order := range instrumentOrders {
        actualOrders[order.OrderID] = true
    }
    
    // Remove phantom orders
    for side, order := range trackedOrders {
        if order != nil && !actualOrders[order.OrderID] {
            debugLog("Removing phantom %s order %s for %s", side, order.OrderID, instrument)
            delete(mm.activeOrders, order.OrderID)
            delete(trackedOrders, side)
        }
    }
    
    // Add untracked orders
    for _, order := range instrumentOrders {
        if _, tracked := mm.activeOrders[order.OrderID]; !tracked {
            log.Printf("Found untracked order %s for %s, adding to tracking", order.OrderID, instrument)
            orderCopy := order
            mm.activeOrders[order.OrderID] = &orderCopy
            if mm.ordersByInstrument[instrument] == nil {
                mm.ordersByInstrument[instrument] = make(map[string]*MarketMakerOrder)
            }
            mm.ordersByInstrument[instrument][order.Side] = &orderCopy
        }
    }
}
```

## 7. stats.go - Statistics and Reporting

```go
package main

import (
    "log"
    "time"
)

// statsReporter periodically reports statistics
func (mm *MarketMaker) statsReporter() {
    defer mm.wg.Done()
    
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    startTime := time.Now()
    
    for {
        select {
        case <-mm.ctx.Done():
            return
        case <-ticker.C:
            mm.reportStats(startTime)
        }
    }
}

// reportStats generates and logs statistics
func (mm *MarketMaker) reportStats(startTime time.Time) {
    mm.mu.RLock()
    defer mm.mu.RUnlock()
    
    mm.stats.UptimeSeconds = int64(time.Since(startTime).Seconds())
    mm.stats.LastUpdate = time.Now()
    
    // Count active orders
    activeCount := 0
    totalOrders := 0
    for _, orders := range mm.ordersByInstrument {
        if len(orders) > 0 {
            activeCount++
            totalOrders += len(orders)
        }
    }
    
    // Consistency check
    if totalOrders != len(mm.activeOrders) {
        log.Printf("WARNING: Order tracking inconsistency: %d in activeOrders, %d in ordersByInstrument", 
            len(mm.activeOrders), totalOrders)
    }
    
    // Log stats
    log.Printf("Stats: Orders=%d/%d/%d (placed/cancelled/filled), Active=%d/%d instruments, Uptime=%ds",
        mm.stats.OrdersPlaced,
        mm.stats.OrdersCancelled,
        mm.stats.OrdersFilled,
        activeCount,
        len(mm.config.Instruments),
        mm.stats.UptimeSeconds)
    
    // Detailed order state in debug mode
    if debugMode {
        mm.logDetailedOrderState()
    }
}

// logDetailedOrderState logs detailed order information (debug mode only)
func (mm *MarketMaker) logDetailedOrderState() {
    for instrument, orders := range mm.ordersByInstrument {
        if len(orders) > 0 {
            bidPrice, askPrice := "none", "none"
            if bidOrder, hasBid := orders["buy"]; hasBid {
                bidPrice = bidOrder.Price.String()
            }
            if askOrder, hasAsk := orders["sell"]; hasAsk {
                askPrice = askOrder.Price.String()
            }
            debugLog("  %s: bid=%s, ask=%s", instrument, bidPrice, askPrice)
        }
    }
}

// getStats returns a copy of current statistics
func (mm *MarketMaker) getStats() MarketMakerStats {
    mm.mu.RLock()
    defer mm.mu.RUnlock()
    
    // Create a copy of stats
    statsCopy := mm.stats
    
    // Copy bid-ask spreads map
    statsCopy.BidAskSpread = make(map[string]decimal.Decimal)
    for k, v := range mm.stats.BidAskSpread {
        statsCopy.BidAskSpread[k] = v
    }
    
    return statsCopy
}

// updateBidAskSpread updates the recorded spread for an instrument
func (mm *MarketMaker) updateBidAskSpread(instrument string, spread decimal.Decimal) {
    mm.mu.Lock()
    defer mm.mu.Unlock()
    
    mm.stats.BidAskSpread[instrument] = spread
}

// Utility functions

var debugMode = false

func debugLog(format string, args ...interface{}) {
    if debugMode {
        log.Printf("[DEBUG] "+format, args...)
    }
}

// EnableDebugMode enables debug logging
func EnableDebugMode() {
    debugMode = true
    log.Println("Debug mode enabled")
}

// DisableDebugMode disables debug logging
func DisableDebugMode() {
    debugMode = false
    log.Println("Debug mode disabled")
}
```

## Key Benefits of This Refactoring

1. **Clear Separation of Concerns**
   - Each file has a specific responsibility
   - Easier to find and modify specific functionality

2. **No Unnecessary Abstractions**
   - Still uses the same `MarketMaker` struct
   - Methods are just organized into logical files
   - No new interfaces or layers added

3. **Improved Maintainability**
   - Smaller files are easier to understand
   - Related functions are grouped together
   - Can work on one aspect without touching others

4. **Easy to Navigate**
   - File names clearly indicate their purpose
   - Functions are where you'd expect to find them

5. **Preserved Functionality**
   - All the original logic is maintained
   - No behavior changes, just better organization