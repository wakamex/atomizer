# Atomizer - Unified Options Trading Toolkit

Atomizer is a high-performance command-line toolkit for automated options trading, providing market making, RFQ response, position management, and real-time market analysis across multiple exchanges.

## Quick Start

```bash
# Clone and install
git clone https://github.com/wakamex/atomizer.git
cd atomizer
go install ./cmd/atomizer

# Set up environment
export DERIVE_PRIVATE_KEY=your_private_key
export DERIVE_WALLET_ADDRESS=your_wallet_address

# Run RFQ responder
atomizer rfq-responder

# Run market maker (aggressive mode by default)
atomizer market-maker --expiry 20250530 --strikes 3000 --size 0.1

# Run market maker (conservative mode)
atomizer market-maker --expiry 20250530 --strikes 3000 --size 0.1 --aggression 0.5

# Run pure gamma hedger (closes perp positions when no options exist)
atomizer pure-gamma-hedger --aggressiveness 1.0  # Cross spread for immediate fills
```

## Features

- **Automated Market Making**: Continuous two-sided quoting with smart order management
- **RFQ Response System**: Real-time quote generation with automatic hedging
- **Multi-Exchange Support**: Unified interface for Derive/Lyra, Deribit, and CCXT-compatible exchanges
- **Risk Management**: Position limits, Greeks calculation, and automated hedging
- **Market Analysis**: Liquidity scoring, volatility analysis, and spread metrics
- **Real-Time Monitoring**: Market data collection with time-series storage

## Documentation

- **[Architecture Guide](./ARCHITECTURE.md)**: System design, components, and technical implementation
- **[Operations Guide](./OPERATIONS.md)**: Deployment, configuration, trading strategies, and troubleshooting

## Commands

### Core Trading Commands

```bash
# RFQ Responder - Automated quote generation
atomizer rfq-responder [options]
  --exchange string     Exchange to use (derive/deribit)
  --enable-hedging      Enable automatic hedging
  --dummy-price string  Fallback price for testing

# Market Maker - Continuous quoting
atomizer market-maker [options]
  --expiry string       Option expiry (YYYYMMDD)
  --strikes string      Comma-separated strike prices
  --size float          Quote size
  --spread int          Spread in basis points
  --aggression float    Aggression: 0=join best, 0.9=near mid, 1.0+=cross spread

# Manual Orders - Direct order placement
atomizer manual-order [options]
  --instrument string   Instrument to trade
  --side string         Order side (buy/sell)
  --price float         Order price
  --amount float        Order amount
```

### Analysis & Monitoring

```bash
# Market Analysis - Liquidity and pricing analysis
atomizer analyze --underlying ETH --expiry 0

# Position Inventory - Current positions and P&L
atomizer inventory

# Market Monitor - Real-time data collection
atomizer market-monitor start --instruments ETH-*
```

### Utility Commands

```bash
# List available markets
atomizer markets --underlying ETH

# Send test quote
atomizer send-quote --instrument ETH-20250530-3000-C --side buy --price 100

# System health check
atomizer health --verbose
```

## Installation

### Prerequisites
- Go 1.21 or later
- Git
- Linux/macOS (Windows via WSL)

### Building from Source

```bash
# Clone repository
git clone https://github.com/wakamex/atomizer.git
cd atomizer

# Install the CLI
go install ./cmd/atomizer

# Or build locally
go build -o atomizer ./cmd/atomizer
```

### Docker

```bash
# Build image
docker build -t atomizer .

# Run container
docker run -e DERIVE_PRIVATE_KEY=$DERIVE_PRIVATE_KEY atomizer rfq-responder
```

## Configuration

### Environment Variables

```bash
# Exchange Configuration
EXCHANGE_NAME=derive              # derive, deribit, or ccxt
EXCHANGE_TEST_MODE=false          # Use testnet

# Derive Authentication
DERIVE_PRIVATE_KEY=0x...          # Private key (without 0x prefix)
DERIVE_WALLET_ADDRESS=0x...       # Wallet address

# Deribit Authentication  
DERIBIT_API_KEY=your_key
DERIBIT_API_SECRET=your_secret

# RFQ Settings
RFQ_ASSET_ADDRESSES=0x7b79...     # Contract addresses
QUOTE_VALID_DURATION=30           # Quote validity in seconds
```

### Risk Parameters

```bash
# Position Limits
MAX_POSITION_DELTA=10.0           # Maximum delta exposure
MAX_POSITION_SIZE=100.0           # Per-instrument limit

# Hedging
ENABLE_GAMMA_HEDGING=true         # Enable gamma hedging
GAMMA_THRESHOLD=0.5               # Gamma hedge threshold
```

## Project Structure

```
atomizer/
├── cmd/
│   └── atomizer/          # Main CLI application
├── internal/              # Core business logic
│   ├── quoter/            # Quote generation engine
│   ├── arbitrage/         # Trade orchestration
│   ├── marketmaker/       # Market making logic
│   ├── hedging/           # Hedging strategies
│   ├── risk/              # Risk management
│   └── exchange/          # Exchange adapters
├── scripts/               # Utility scripts
└── sdk/                   # Rysk SDK submodule
```

## Development

```bash
# Run tests
go test ./...

# Run with debug logging
atomizer rfq-responder --debug --log-level trace

# Profile performance
atomizer debug profile --cpu --duration 30s

# Format code
go fmt ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- **Issues**: [GitHub Issues](https://github.com/wakamex/atomizer/issues)
- **Documentation**: [Architecture](./ARCHITECTURE.md) | [Operations](./OPERATIONS.md)
- **Examples**: See the `examples/` directory for sample configurations