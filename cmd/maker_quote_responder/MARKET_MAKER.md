# Market Maker for Derive/Deribit Options

A generic market maker that automatically quotes bid and ask prices on options exchanges.

## Features

- **Exchange Agnostic**: Supports multiple exchanges (currently Derive, Deribit coming soon)
- **Real-time Quotes**: Direct WebSocket connection for low-latency ticker updates
- **Risk Management**: Position limits and exposure controls
- **Automatic Requoting**: Updates quotes when market moves beyond threshold
- **Graceful Shutdown**: Cancels all orders on exit

## Architecture

```
MarketMaker
    ├── MarketMakerInterface    (Generic exchange interface)
    ├── DeriveMarketMakerExchange (Derive implementation)
    ├── MarketMakerConfig       (Configuration)
    └── Order Tracking          (Active order management)
```

## Quick Start

### 1. Set Environment Variables

```bash
# For Derive
export DERIVE_PRIVATE_KEY="your_private_key"
export DERIVE_WALLET_ADDRESS="0x..."

# For Deribit
export DERIBIT_API_KEY="your_api_key"
export DERIBIT_API_SECRET="your_secret"
```

### 2. Run Market Maker

```bash
# Basic usage - ETH options expiring June 6, 2025
./run_market_maker.sh

# Custom configuration
EXPIRY=20250613 STRIKES="2800,2900,3000" ./run_market_maker.sh

# All strikes with wider spread
ALL_STRIKES=true SPREAD_BPS=20 ./run_market_maker.sh

# Dry run to see configuration
DRY_RUN=true ./run_market_maker.sh
```

## Command Line Options

```bash
go run cmd/market_maker/main.go [options]

Options:
  -exchange string     Exchange to use (derive, deribit) (default "derive")
  -test               Use exchange testnet
  -expiry string      Expiry date (e.g., 20250606)
  -underlying string  Underlying asset (default "ETH")
  -strikes string     Comma-separated strikes (e.g., "2800,2900,3000")
  -all-strikes        Make markets on all available strikes
  -spread int         Spread in basis points (default 10)
  -size float         Quote size (default 1.0)
  -refresh int        Quote refresh interval in seconds (default 1)
  -max-position float Maximum position per instrument (default 10.0)
  -max-exposure float Maximum total exposure (default 100.0)
  -min-spread int     Minimum spread in basis points (default 5)
  -dry-run            Print configuration without starting
  -version            Print version and exit
```

## Pricing Strategy

The market maker improves the best bid and ask by 0.1:
- **Our Bid**: Best Bid + 0.1
- **Our Ask**: Best Ask - 0.1

This ensures our quotes are at the top of the order book while maintaining profitability through the spread.

## Risk Management

1. **Position Limits**: Maximum position size per instrument
2. **Total Exposure**: Maximum exposure across all instruments
3. **Minimum Spread**: Ensures profitability
4. **Cancel Threshold**: Requotes when market moves >0.5%

## Examples

### Market Make on Specific Strikes
```bash
./run_market_maker.sh \
  EXPIRY=20250606 \
  STRIKES="2700,2800,2900,3000,3100" \
  SPREAD_BPS=15 \
  SIZE=2.0
```

### Conservative Settings
```bash
./run_market_maker.sh \
  SPREAD_BPS=20 \
  SIZE=0.5 \
  MAX_POSITION=5.0 \
  MAX_EXPOSURE=50.0
```

### Aggressive Settings
```bash
./run_market_maker.sh \
  SPREAD_BPS=5 \
  SIZE=5.0 \
  MAX_POSITION=20.0 \
  REFRESH_SEC=0.5
```

## Monitoring

The market maker logs:
- Order placement and cancellation
- Position updates
- Statistics every 30 seconds
- Error messages

## Safety Features

1. **Graceful Shutdown**: Ctrl+C cancels all orders before exiting
2. **Position Tracking**: Monitors net exposure
3. **Order Validation**: Checks risk limits before placing orders
4. **Connection Monitoring**: Handles WebSocket disconnections

## TODO

- [ ] Filled order tracking and PnL calculation
- [ ] Deribit exchange implementation
- [ ] Dynamic spread adjustment based on volatility
- [ ] Inventory management strategies
- [ ] Web dashboard for monitoring