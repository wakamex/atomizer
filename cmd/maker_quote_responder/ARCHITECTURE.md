# Exchange Integration Architecture

## Overview

This document describes the architecture for integrating multiple cryptocurrency exchanges into the market maker system. The design supports both CCXT-compatible exchanges and custom implementations while minimizing code duplication.

## Design Principles

1. **DRY (Don't Repeat Yourself)**: Share common code between CCXT exchanges
2. **Open/Closed**: Easy to add new exchanges without modifying existing code
3. **Interface Segregation**: Clean interfaces that exchanges can implement
4. **Composition over Inheritance**: Use composition since Go doesn't support inheritance

## Architecture Diagram

```
┌─────────────────────────────────────────────────────┐
│                 Exchange Interface                   │
├─────────────────────────────────────────────────────┤
│  - GetOrderBook()                                   │
│  - PlaceHedgeOrder()                                │
│  - ConvertToInstrument()                            │
└─────────────────────────────────────────────────────┘
                    ▲              ▲
                    │              │
    ┌───────────────┴───┐    ┌────┴────────────────┐
    │  CCXTExchange     │    │  CustomExchange     │
    │  (Generic)        │    │  (Base)             │
    ├───────────────────┤    ├────────────────────┤
    │  - ccxtClient     │    │  - httpClient      │
    │  - symbolFormat   │    │  - baseURL         │
    │  - commonMethods  │    │  - authenticate()  │
    └───────────────────┘    └────────────────────┘
              ▲                        ▲
              │                        │
    ┌─────────┴─────────┐    ┌────────┴──────────┐
    │ Per-Exchange      │    │ DeribitAsymmetric │
    │ Configuration     │    │ - Ed25519 auth    │
    │ - Deribit         │    │ - Custom signing  │
    │ - Bybit           │    └───────────────────┘
    │ - OKX             │
    │ - Derive          │
    └───────────────────┘
```

## Components

### 1. Exchange Interface

The core interface that all exchanges must implement:

```go
type Exchange interface {
    GetOrderBook(req RFQResult, asset string) (CCXTOrderBook, error)
    PlaceHedgeOrder(conf RFQConfirmation, underlying string, cfg *AppConfig) error
    ConvertToInstrument(asset, strike string, expiry int64, isPut bool) (string, error)
}
```

### 2. Generic CCXT Adapter

Handles any CCXT-supported exchange with configuration:

```go
type CCXTGenericExchange struct {
    exchange ccxt.IExchange
    config   CCXTConfig
}

type CCXTConfig struct {
    ExchangeName      string
    APIKey           string
    APISecret        string
    TestMode         bool
    
    // Exchange-specific formatters
    SymbolFormatter  func(asset, instrument string) string
    InstrumentParser func(asset, strike string, expiry int64, isPut bool) string
    
    // Optional overrides
    GetOrderBook     func(symbol string) (*OrderBook, error)
    PlaceOrder       func(symbol, side string, amount, price float64) (*Order, error)
}
```

### 3. Exchange Registry

Configuration-driven exchange definitions:

```go
type ExchangeDefinition struct {
    CCXTName         string
    SymbolFormat     string  // e.g., "%s/USD:%s" for Deribit
    DateFormat       string  // e.g., "2Jan06" for Deribit
    StrikeDivisor    int64   // e.g., 1e8 for converting from wei
    SupportsOptions  bool
    RequiresUSDFlag  bool    // Some exchanges need "advanced": "usd"
}

var exchangeRegistry = map[string]ExchangeDefinition{
    "deribit": {
        CCXTName:        "deribit",
        SymbolFormat:    "%s/USD:%s",
        DateFormat:      "2Jan06",
        StrikeDivisor:   1e8,
        SupportsOptions: true,
        RequiresUSDFlag: true,
    },
    "derive": {
        CCXTName:        "derive",
        SymbolFormat:    "%s:%s",
        DateFormat:      "20060102",
        StrikeDivisor:   1e8,
        SupportsOptions: true,
        RequiresUSDFlag: false,
    },
    "bybit": {
        CCXTName:        "bybit",
        SymbolFormat:    "%s-%s",
        DateFormat:      "02JAN06",
        StrikeDivisor:   1,
        SupportsOptions: true,
        RequiresUSDFlag: false,
    },
}
```

### 4. Exchange Factory

Simplified factory that handles both generic and custom exchanges:

```go
func CreateExchange(name string, cfg *AppConfig) (Exchange, error) {
    // Check for custom implementations first
    switch name {
    case "deribit":
        // Special handling for asymmetric keys
        if hasAsymmetricKeys() {
            return NewDeribitAsymmetricExchange(cfg)
        }
    }
    
    // Look up in registry
    def, exists := exchangeRegistry[name]
    if !exists {
        return nil, fmt.Errorf("unsupported exchange: %s", name)
    }
    
    // Create generic CCXT exchange
    return NewCCXTGenericExchange(def, cfg)
}
```

## Adding a New Exchange

### For CCXT-Supported Exchanges

1. Add entry to `exchangeRegistry`:
```go
"newexchange": {
    CCXTName:        "newexchange",
    SymbolFormat:    "%s-%s",
    DateFormat:      "060102",
    StrikeDivisor:   1e8,
    SupportsOptions: true,
},
```

2. Add credentials to environment:
```bash
NEWEXCHANGE_API_KEY=xxx
NEWEXCHANGE_API_SECRET=xxx
```

3. That's it! The generic adapter handles the rest.

### For Custom Exchanges

1. Implement the `Exchange` interface
2. Add case to factory switch statement
3. Handle authentication and API calls

## Benefits

1. **Minimal Code**: New CCXT exchanges need only configuration
2. **Flexibility**: Can override any method for specific needs
3. **Type Safety**: Strongly typed interfaces
4. **Testability**: Easy to mock exchanges for testing
5. **Maintainability**: Changes to common logic affect all exchanges

## Migration Strategy

1. **Phase 1**: Keep existing implementations working
2. **Phase 2**: Create generic CCXT adapter
3. **Phase 3**: Migrate existing exchanges to use generic adapter
4. **Phase 4**: Add new exchanges using configuration only

## Example Usage

```go
// Automatically selects the right implementation
exchange, err := CreateExchange("deribit", appConfig)

// Works the same for all exchanges
orderBook, err := exchange.GetOrderBook(rfq, "ETH")
err = exchange.PlaceHedgeOrder(confirmation, "ETH", appConfig)
```

## Future Enhancements

1. **Exchange Capabilities**: Query what each exchange supports
2. **Failover**: Automatic fallback to secondary exchange
3. **Smart Routing**: Choose best exchange based on liquidity/price
4. **Unified Error Handling**: Standard error types across exchanges
5. **Metrics Collection**: Track performance per exchange

## Testing Strategy

1. **Interface Tests**: Verify all exchanges implement interface correctly
2. **Mock Exchange**: For testing business logic without API calls
3. **Integration Tests**: Real API calls with small amounts
4. **Configuration Tests**: Verify registry entries are valid