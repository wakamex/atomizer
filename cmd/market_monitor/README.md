# Market Monitor

Real-time market data collection and storage system for options trading.

## Overview

The Market Monitor collects order book and spot price data from multiple exchanges and stores it in a time-series database (VictoriaMetrics) for analysis and backtesting.

## Features

- **Multi-Exchange Support**: Collects data from Derive/Lyra and Deribit
- **Order Book Depth**: Configurable depth collection (default: 10 levels)
- **Spot Price Feeds**: Real-time ETH and BTC spot prices for accurate pricing
- **USD Conversion**: Displays Deribit ETH-denominated prices in USD (display only, original values stored)
- **Instrument Conversion**: Automatic conversion between exchange naming formats
  - Derive: `ETH-20250531-2700-C`
  - Deribit: `ETH-31MAY25-2700-C`
- **WebSocket Support**: Real-time data collection via WebSocket connections
- **Debug Mode**: Detailed logging for troubleshooting

## Installation

### Prerequisites

1. Build the market monitor:
```bash
cd /path/to/atomizer
./build.sh
```

2. Set up VictoriaMetrics:
```bash
atomizer market-monitor setup
```

This downloads VictoriaMetrics and creates the necessary directory structure.

## Usage

### Start VictoriaMetrics

```bash
cd market_monitor_data
./victoria-metrics-prod -storageDataPath ./vm-data
```

### Start Market Monitor

Basic ticker collection:
```bash
atomizer market-monitor start
```

Order book depth collection:
```bash
atomizer market-monitor start --orderbook --depth 10
```

Monitor specific instruments:
```bash
atomizer market-monitor start --exchanges derive,deribit --instruments ETH-*-C,BTC-*-C
```

With debug output:
```bash
atomizer market-monitor start --orderbook --debug
```

### Command Line Options

- `--interval DURATION`: Collection interval (default: 5s)
- `--exchanges LIST`: Comma-separated exchanges (default: derive,deribit)
- `--instruments PATTERN`: Instrument patterns or exact names (default: ETH-PERP)
- `--vm-url URL`: VictoriaMetrics URL (default: http://localhost:8428)
- `--workers N`: Number of concurrent workers (default: 10)
- `--orderbook`: Collect order book depth instead of just ticker
- `--depth N`: Order book depth to collect (default: 10)
- `--debug`: Enable debug logging

## Data Storage

Data is stored in VictoriaMetrics in Prometheus format:

### Metrics

**Ticker Data:**
- `options_best_bid{exchange="...", instrument="..."}`
- `options_best_ask{exchange="...", instrument="..."}`
- `options_bid_size{exchange="...", instrument="..."}`
- `options_ask_size{exchange="...", instrument="..."}`

**Order Book Data:**
- `orderbook_bid_price{exchange="...", instrument="...", level="0"}`
- `orderbook_bid_size{exchange="...", instrument="...", level="0"}`
- `orderbook_ask_price{exchange="...", instrument="...", level="0"}`
- `orderbook_ask_size{exchange="...", instrument="...", level="0"}`

**Spot Prices:**
- `market_bid_price{instrument="ETH-SPOT", exchange="derive"}`
- `market_bid_price{instrument="BTC-SPOT", exchange="derive"}`

## Data Analysis

After collecting data, use the analysis tools in `cmd/analyze_correlation/`:

- **analyze_correlation.py**: Analyze price correlation between exchanges
- **inspect_orderbook.py**: View real-time and historical orderbook data

See [analyze_correlation/README.md](../analyze_correlation/README.md) for details.

## Architecture

### Components

1. **Main Monitor** (`cmd/market_monitor/main.go`)
   - CLI interface and configuration
   - Subcommand routing

2. **Monitor Package** (`internal/monitor/`)
   - `monitor.go`: Main monitoring orchestration
   - `orderbook_monitor.go`: Order book collection coordinator
   - `ticker_collector.go`: REST API ticker collection
   - `orderbook_collector.go`: REST API order book collection
   - `derive_ws_orderbook.go`: Derive WebSocket order book collection
   - `spot_collector.go`: Spot price collection via WebSocket
   - `storage.go`: VictoriaMetrics storage interface
   - `instrument_converter.go`: Exchange naming conversion
   - `logger.go`: Debug logging utilities

### Data Flow

1. **Collection**: Data is collected from exchanges via REST APIs or WebSocket
2. **Conversion**: Instrument names are converted between exchange formats
3. **Processing**: Order books are normalized to common format
4. **Storage**: Data is sent to VictoriaMetrics via Prometheus remote write API
5. **Display**: Deribit ETH prices are converted to USD for display (original stored)

## Examples

### Monitor All ETH Options
```bash
atomizer market-monitor start --orderbook --instruments "ETH-*" --interval 10s
```

### High-Frequency Collection
```bash
atomizer market-monitor start --orderbook --interval 1s --workers 20
```

### Debug WebSocket Issues
```bash
atomizer market-monitor start --orderbook --debug --instruments ETH-1JUN25-2600-C
```

## Querying Data

Once data is collected, you can query it using VictoriaMetrics' PromQL API:

```bash
# Get latest ETH spot price
curl 'http://localhost:8428/api/v1/query?query=spot_price{currency="ETH"}'

# Get order book for specific instrument
curl 'http://localhost:8428/api/v1/query?query=orderbook_bid_price{instrument="ETH-1JUN25-2600-C"}'

# Get bid-ask spread over time
curl 'http://localhost:8428/api/v1/query_range?query=options_best_ask-options_best_bid&start=2024-01-01T00:00:00Z&end=2024-01-02T00:00:00Z&step=1m'
```

## Troubleshooting

### No Data Being Collected
- Check VictoriaMetrics is running: `curl http://localhost:8428/health`
- Enable debug mode to see raw messages: `--debug`
- Verify instrument names are valid for the exchange

### WebSocket Connection Issues
- Check network connectivity to exchange APIs
- Enable debug mode to see connection status
- Look for "Ping received, sending Pong" messages

### Instrument Not Found
- Verify the instrument is live on the exchange
- Check date format matches exchange expectations
- Use exact instrument names instead of patterns for testing

## Future Enhancements

- [ ] Historical data export functionality
- [ ] Real-time alerting on market conditions
- [ ] Greeks calculation and storage
- [ ] Volatility surface construction
- [ ] Cross-exchange arbitrage detection