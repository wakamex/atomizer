# Atomizer Architecture

This document describes the system architecture, design principles, and technical implementation details of the Atomizer options trading toolkit.

## Table of Contents
- [System Overview](#system-overview)
- [Design Principles](#design-principles)
- [Component Architecture](#component-architecture)
- [Package Structure](#package-structure)
- [Exchange Integration](#exchange-integration)
- [Data Flow](#data-flow)
- [Extension Points](#extension-points)

## System Overview

Atomizer is a modular options trading system designed for automated market making, RFQ response, and position management across multiple exchanges. The system follows a clean architecture pattern with clear separation of concerns.

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      CLI Interface                           │
│                   (cmd/atomizer/main.go)                     │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────┴───────────────────────────────────────┐
│                    Core Components                           │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │   Quoter    │  │     RFQ      │  │   Arbitrage     │   │
│  │  (quoter/)  │  │  Processor   │  │  Orchestrator   │   │
│  └─────────────┘  └──────────────┘  └─────────────────┘   │
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │    Risk     │  │   Hedging    │  │  Market Maker   │   │
│  │  Manager    │  │   Manager    │  │   (marketmaker/)│   │
│  └─────────────┘  └──────────────┘  └─────────────────┘   │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────┴───────────────────────────────────────┐
│                 Exchange Abstraction Layer                   │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │   Deribit   │  │    Derive    │  │      CCXT       │   │
│  │  Adapter    │  │   Adapter    │  │    Wrapper      │   │
│  └─────────────┘  └──────────────┘  └─────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Design Principles

### 1. Modularity
- Each component has a single, well-defined responsibility
- Components communicate through interfaces, not concrete implementations
- Easy to add new exchanges or trading strategies

### 2. Asynchronous Processing
- Non-blocking WebSocket communication
- Concurrent order management
- Event-driven architecture for real-time responsiveness

### 3. Fault Tolerance
- Automatic reconnection for WebSocket connections
- Graceful degradation (e.g., fallback to dummy quotes)
- Comprehensive error handling and logging

### 4. Performance
- Minimal latency for quote generation
- Efficient order book processing
- Optimized for high-frequency operations

## Component Architecture

### Core Trading Components

#### Quoter (`internal/quoter/`)
The heart of the pricing engine, responsible for:
- Generating EIP-712 signed quotes
- Calculating prices with slippage protection
- Computing APR for options
- Managing quote validity periods

```go
type Quoter interface {
    MakeQuote(rfq RFQResult, underlying string, rfqID string) (Quote, error)
}
```

#### RFQ Processor (`internal/rfq/`)
Handles incoming RFQ requests:
- Debouncing duplicate requests
- Routing to appropriate pricing sources
- Managing quote responses
- Fallback to dummy quotes when needed

#### Arbitrage Orchestrator (`internal/arbitrage/`)
Coordinates trading activities across different sources:
- RFQ trades from WebSocket
- Manual trades from HTTP API
- Hedging operations
- Position synchronization

### Risk Management

#### Risk Manager (`internal/risk/`)
Monitors and controls risk exposure:
- Position limits per instrument
- Total portfolio exposure limits
- Greeks calculation and limits
- Real-time risk metrics

#### Hedging Manager (`internal/hedging/`)
Executes hedging strategies:
- Delta hedging for options positions
- Gamma hedging for large positions
- Cross-exchange hedging
- Smart order routing

### Market Making

#### Market Maker (`internal/marketmaker/`)
Provides continuous two-sided quotes:
- Order placement and management
- Spread calculation based on volatility
- Position-aware pricing with aggression control
- Inventory management

**Aggression Parameter**:
- `0.0-0.9`: Conservative mode - places orders between best bid/ask and mid
- `1.0+`: Aggressive mode - can cross the spread for better fills

### Infrastructure Components

#### WebSocket Client (`internal/websocket/`)
Reusable WebSocket infrastructure:
- Automatic reconnection with exponential backoff
- Authentication adapters for different exchanges
- Message routing and parsing
- Connection health monitoring

#### Exchange Adapters (`internal/exchange/`)
Standardized interface for multiple exchanges:
- Deribit: Full derivatives support
- Derive/Lyra: EVM-based options protocol
- CCXT: Generic exchange wrapper

## Package Structure

```
internal/
├── api/           # HTTP REST API server
├── arbitrage/     # Trade coordination and arbitrage logic
├── cache/         # Market data caching (in-memory and Redis)
├── config/        # Configuration management
├── exchange/      # Exchange adapters and interfaces
│   ├── ccxt/      # CCXT wrapper for generic exchanges
│   ├── deribit/   # Deribit-specific implementation
│   └── derive/    # Derive/Lyra protocol implementation
├── hedging/       # Hedging strategies and execution
│   └── gamma/     # Gamma hedging module
├── manual/        # Manual order management
├── marketmaker/   # Market making engine
├── monitor/       # Market data collection and monitoring
├── quoter/        # Quote generation and pricing
├── rfq/           # RFQ processing
├── risk/          # Risk management
├── types/         # Shared type definitions
├── utils/         # Common utilities
└── websocket/     # WebSocket client and adapters
```

## Exchange Integration

### Exchange Interface
All exchanges implement a common interface:

```go
type Exchange interface {
    GetOrderBook(rfq RFQResult, asset string) (OrderBook, error)
    PlaceOrder(confirmation RFQConfirmation, instrument string) error
    ConvertToInstrument(asset, strike string, expiry int64, isPut bool) (string, error)
    GetPositions() ([]Position, error)
}
```

### Adding New Exchanges
1. Implement the `Exchange` interface
2. Create an adapter in `internal/exchange/`
3. Register in the exchange factory
4. Add exchange-specific configuration

### Authentication Patterns
- **API Key/Secret**: Traditional REST API auth (Deribit)
- **EIP-712 Signing**: Ethereum-based authentication (Derive)
- **WebSocket Auth**: Custom authentication flows

## Data Flow

### RFQ Flow
1. WebSocket receives RFQ notification
2. RFQ processor validates and deduplicates
3. Quoter generates price using exchange data
4. Quote is signed and sent back
5. Confirmation triggers hedging

### Order Flow
1. Market maker calculates desired orders
2. Compares with existing orders
3. Cancels outdated orders
4. Places new orders
5. Updates internal state

### Risk Flow
1. Position changes trigger risk recalculation
2. Greeks are computed for all positions
3. Risk limits are checked
4. Hedging orders are generated if needed

## Extension Points

### Custom Trading Strategies
Implement the `TradingStrategy` interface:
- Define entry/exit conditions
- Set position sizing rules
- Configure risk parameters

### New Data Sources
Add new market data feeds:
- Implement data adapter
- Configure caching strategy
- Set up monitoring

### Analytics Modules
Extend monitoring capabilities:
- Custom metrics collection
- Performance analytics
- Trade analysis

### Risk Models
Implement custom risk calculations:
- Alternative Greeks models
- Portfolio optimization
- Stress testing

## Performance Considerations

### Latency Optimization
- Pre-compute common calculations
- Cache market data appropriately
- Use connection pooling
- Minimize serialization overhead

### Scalability
- Horizontal scaling for WebSocket connections
- Separate read/write paths
- Async processing for non-critical paths
- Rate limiting and backpressure handling

### Monitoring
- Prometheus metrics for all components
- Structured logging with context
- Distributed tracing support
- Health checks and alerting