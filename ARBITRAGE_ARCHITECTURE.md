# Cross-Protocol Arbitrage Bot Architecture

## Overview

This document describes the modular architecture for a cross-protocol arbitrage bot that:
1. Buys call options on Rysk Finance
2. Hedges positions on a single exchange (Derive or Deribit)
3. Implements dynamic delta hedging using gamma hedging strategies

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     Arbitrage Orchestrator                       │
├─────────────────────────────────────────────────────────────────┤
│  - Coordinates all modules                                       │
│  - Manages async trade execution                                 │
│  - Handles manual trade initiation                               │
│  - Monitors P&L and risk metrics                                 │
└─────────────────────────────────────────────────────────────────┘
                    ▲                    ▲                    ▲
                    │                    │                    │
    ┌───────────────┴────────┐  ┌───────┴────────┐  ┌───────┴────────┐
    │  Trade Source Module   │  │ Hedge Manager   │  │ Risk Manager    │
    ├────────────────────────┤  ├─────────────────┤  ├─────────────────┤
    │ - Rysk RFQ listener    │  │ - Single exchange│  │ - Position limits│
    │ - Manual trade input   │  │ - Order placement│  │ - Greeks calc   │
    │ - Trade validation     │  │ - Execution algo│  │ - Stop loss     │
    └────────────────────────┘  └─────────────────┘  └─────────────────┘
                    ▲                    ▲                    ▲
                    │                    │                    │
    ┌───────────────┴────────┐  ┌───────┴────────┐  ┌───────┴────────┐
    │  Rysk Connector        │  │ Exchange Client │  │ Gamma DDH Algo  │
    ├────────────────────────┤  ├─────────────────┤  ├─────────────────┤
    │ - WebSocket client     │  │ - Deribit OR    │  │ - Delta calc    │
    │ - Quote responder      │  │ - Derive        │  │ - Gamma calc    │
    │ - Trade executor       │  │ - Single conn   │  │ - Hedge orders  │
    └────────────────────────┘  └─────────────────┘  └─────────────────┘
```

## Core Components

### 1. Arbitrage Orchestrator

The central coordinator that manages the entire arbitrage flow.

```go
type ArbitrageOrchestrator struct {
    tradeSource    TradeSource
    hedgeManager   *HedgeManager
    riskManager    *RiskManager
    gammaDDH       *GammaDDHAlgo
    tradeQueue     chan TradeEvent
    config         *ArbitrageConfig
}

type ArbitrageConfig struct {
    MaxPositionSize      decimal.Decimal
    MaxDeltaExposure     decimal.Decimal
    HedgeDelayMs         int
    EnableManualTrades   bool
    EnableAutoHedging    bool
    TargetSpreadBps      int
    ExchangeName         string  // "derive" or "deribit"
}
```

### 2. Trade Source Module

Handles incoming trades from multiple sources:

```go
type TradeSource interface {
    StartListening(ctx context.Context) error
    OnTradeReceived(handler TradeHandler)
    SubmitManualTrade(trade ManualTrade) error
}

type TradeEvent struct {
    ID              string
    Source          TradeSourceType // RYSK_RFQ or MANUAL
    Instrument      string
    Strike          decimal.Decimal
    Expiry          time.Time
    IsPut           bool
    Quantity        decimal.Decimal
    Price           decimal.Decimal
    Timestamp       time.Time
}

type TradeHandler func(trade TradeEvent) error
```

### 3. Hedge Manager

Manages hedge execution on a single configured exchange:

```go
type HedgeManager struct {
    exchange       Exchange
    executionAlgo  ExecutionAlgorithm
    config         *HedgeConfig
}

type HedgeConfig struct {
    ExchangeName   string
    MaxSpreadBps   int
    OrderType      OrderType
    TimeInForce    TimeInForce
}

type HedgeOrder struct {
    TradeID        string
    Instrument     string
    Direction      Direction
    Quantity       decimal.Decimal
    PriceLimit     decimal.Decimal
    TimeInForce    TimeInForce
}

func (hm *HedgeManager) ExecuteHedge(trade TradeEvent) (*HedgeResult, error) {
    // 1. Get current orderbook from exchange
    // 2. Calculate hedge quantity based on trade
    // 3. Place order within best ask
    // 4. Return execution result
}
```

### 4. Exchange Client

Manages connection to a single exchange:

```go
type ExchangeClient struct {
    exchange     Exchange
    healthCheck  *HealthChecker
    cache        MarketCache
    config       ExchangeConfig
}

func (ec *ExchangeClient) GetOrderBook(instrument string) (*OrderBook, error) {
    // Get orderbook from configured exchange
    // Use cache if available
}

func (ec *ExchangeClient) PlaceOrder(order HedgeOrder) (*OrderResult, error) {
    // Place order on exchange
    // Handle retries on failure
}
```

### 5. Market Data Cache

Flexible caching layer supporting both file and Valkey backends:

```go
type CacheConfig struct {
    Backend     string // "file" or "valkey"
    ValkeyAddr  string // "localhost:6379" for Valkey
    FileDir     string // "./cache" for file cache
    DefaultTTL  time.Duration
}

// Switch between cache implementations
func NewMarketCache(config CacheConfig) (MarketCache, error) {
    switch config.Backend {
    case "file":
        return NewFileMarketCache(config.FileDir)
    case "valkey":
        return NewValkeyMarketCache(config.ValkeyAddr)
    default:
        return nil, fmt.Errorf("unknown cache backend: %s", config.Backend)
    }
}

// Usage in ExchangeClient
func NewExchangeClient(exchangeName string, cacheConfig CacheConfig) (*ExchangeClient, error) {
    cache, err := NewMarketCache(cacheConfig)
    if err != nil {
        return nil, err
    }
    
    exchange, err := CreateExchange(exchangeName, appConfig)
    if err != nil {
        return nil, err
    }
    
    return &ExchangeClient{
        exchange: exchange,
        cache:    cache,
        config:   exchangeConfig,
    }, nil
}
```

### 6. Gamma Dynamic Delta Hedge (DDH) Algorithm

Existing implementation from `gamma.go` integrated as a module:

```go
type GammaDDHModule struct {
    algo         *GammaDDHAlgo
    marketData   MarketData
    wsClient     WsClient
    positions    map[string]Position
}

func (g *GammaDDHModule) OnNewPosition(trade TradeEvent) {
    // Update position tracking
    // Trigger hedge recalculation
}

func (g *GammaDDHModule) RunHedgingLoop(ctx context.Context) error {
    // Continuous hedging based on gamma/delta
}
```

## Async Trade Flow

### 1. Trade Reception
```go
func (o *ArbitrageOrchestrator) ProcessTradeAsync(trade TradeEvent) {
    go func() {
        // Validate trade
        if err := o.riskManager.ValidateTrade(trade); err != nil {
            log.Printf("Trade rejected: %v", err)
            return
        }
        
        // Queue for processing
        select {
        case o.tradeQueue <- trade:
        case <-time.After(5 * time.Second):
            log.Printf("Trade queue timeout")
        }
    }()
}
```

### 2. Async Hedge Execution
```go
func (o *ArbitrageOrchestrator) processTradeQueue(ctx context.Context) {
    for {
        select {
        case trade := <-o.tradeQueue:
            go o.executeTrade(trade)
        case <-ctx.Done():
            return
        }
    }
}

func (o *ArbitrageOrchestrator) executeTrade(trade TradeEvent) {
    // Execute on Rysk
    ryskResult, err := o.executeRyskTrade(trade)
    if err != nil {
        log.Printf("Rysk execution failed: %v", err)
        return
    }
    
    // Hedge on configured exchange
    hedgeResult, err := o.hedgeManager.ExecuteHedge(ryskResult)
    if err != nil {
        log.Printf("Hedge failed: %v", err)
        // Could implement retry logic here
        return
    }
    
    log.Printf("Trade %s hedged successfully: %+v", trade.ID, hedgeResult)
    
    // Update gamma hedging
    o.gammaDDH.OnNewPosition(trade)
}
```

## Manual Trade Initiation

```go
type ManualTradeAPI struct {
    orchestrator *ArbitrageOrchestrator
}

func (api *ManualTradeAPI) SubmitTrade(req ManualTradeRequest) (*TradeResponse, error) {
    trade := TradeEvent{
        ID:         uuid.New().String(),
        Source:     MANUAL,
        Instrument: req.Instrument,
        Strike:     req.Strike,
        Expiry:     req.Expiry,
        IsPut:      req.IsPut,
        Quantity:   req.Quantity,
        Price:      req.Price,
        Timestamp:  time.Now(),
    }
    
    // Process through same flow as RFQ trades
    o.orchestrator.ProcessTradeAsync(trade)
    
    return &TradeResponse{
        TradeID: trade.ID,
        Status:  "PENDING",
    }, nil
}
```

## Risk Management

```go
type RiskManager struct {
    maxPositionSize   decimal.Decimal
    maxDeltaExposure  decimal.Decimal
    currentPositions  map[string]Position
    mu                sync.RWMutex
}

func (rm *RiskManager) ValidateTrade(trade TradeEvent) error {
    rm.mu.RLock()
    defer rm.mu.RUnlock()
    
    // Check position limits
    currentSize := rm.getPositionSize(trade.Instrument)
    if currentSize.Add(trade.Quantity).GreaterThan(rm.maxPositionSize) {
        return fmt.Errorf("position limit exceeded")
    }
    
    // Check delta exposure
    deltaImpact := rm.calculateDeltaImpact(trade)
    if deltaImpact.GreaterThan(rm.maxDeltaExposure) {
        return fmt.Errorf("delta limit exceeded")
    }
    
    return nil
}
```

## Configuration

```yaml
arbitrage:
  # Trade sources
  sources:
    rysk:
      enabled: true
      websocket_url: "wss://rip-testnet.rysk.finance/maker"
      assets:
        - "0xb67bfa7b488df4f2efa874f4e59242e9130ae61f"
    manual:
      enabled: true
      api_port: 8080
  
  # Hedging configuration
  hedging:
    exchange: "derive"  # or "deribit"
    
    # Execution parameters
    max_spread_bps: 50
    order_type: "post_only"
    retry_attempts: 3
  
  # Cache configuration
  cache:
    backend: "valkey"  # or "file"
    valkey_addr: "localhost:6379"
    file_dir: "./cache"
    default_ttl: 3600  # 1 hour
  
  # Gamma hedging
  gamma_ddh:
    enabled: true
    perp_name: "ETH-PERPETUAL"
    max_abs_delta: "100"
    max_abs_spread: "0.01"
    action_wait_ms: 1000
    
  # Risk limits
  risk:
    max_position_size: "1000"
    max_delta_exposure: "500"
    stop_loss_pct: "10"
```

## Monitoring and Observability

```go
type MetricsCollector struct {
    tradesExecuted   prometheus.Counter
    hedgeLatency     prometheus.Histogram
    deltaExposure    prometheus.Gauge
    gammaExposure    prometheus.Gauge
    pnl              prometheus.Gauge
}

// Prometheus metrics
arbitrage_trades_total{source="rysk|manual", status="success|failed"}
arbitrage_hedge_latency_seconds{exchange="derive|deribit"}
arbitrage_delta_exposure{instrument="ETH"}
arbitrage_gamma_exposure{instrument="ETH"}
arbitrage_pnl_usd{strategy="call_arbitrage"}
```

## Deployment Considerations

### 1. Single Exchange Simplicity
- Configure exchange via environment variable
- Single connection pool per instance
- No routing logic needed

### 2. Latency Optimization
- Colocate with exchange APIs
- Use exchange WebSocket feeds
- Implement pre-hedge calculations

### 3. Error Recovery
- Automatic reconnection
- Position reconciliation
- Dead letter queue for failed trades

### 4. Scaling Strategy
- Run multiple instances with different exchanges
- Each instance handles one exchange only
- Shared monitoring and risk management

## Future Enhancements

1. **Multi-Asset Support**: Extend beyond ETH to BTC, SOL
2. **Advanced Strategies**: Iron condors, butterflies
3. **ML-Based Pricing**: Optimize hedge timing
4. **Cross-Chain Arbitrage**: Include DeFi protocols
5. **Automated Market Making**: Two-way quotes on Rysk