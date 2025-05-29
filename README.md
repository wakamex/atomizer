# Atomizer - Unified Options Trading Toolkit

Atomizer is a comprehensive command-line toolkit for options trading, providing automated market making, RFQ response, market analysis, and position management capabilities.

## Features

- **Unified CLI**: Single `atomizer` command with subcommands for all functionality
- **RFQ Responder**: Automated quote generation for Request-for-Quote systems
- **Market Analysis**: Real-time liquidity analysis with spread metrics and volatility comparisons
- **Position Management**: Track inventory and P&L across multiple exchanges
- **Multi-Exchange Support**: Works with Derive/Lyra and Deribit (more coming soon)
- **Automated Hedging**: Hedge positions automatically when trades are executed

## Progress

sdk actions:
- [x] approve - leave in the cli
- [x] balances - debugging
- [x] connect - working w/ default channel id
- [ ] positions
- [x] quote - maker_quote_response.go
- [ ] transfer
combo:
- [x] maker_quote_response.go (connect and quote)
- [x] hedging integration with Deribit
- [x] market analysis tools (inventory, markets, send_quote)


## Installation

### Prerequisites
- Go 1.21 or later
- Git

### Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/wakamex/atomizer.git
   cd atomizer
   ```

2. Build all components:
   ```bash
   ./build.sh
   ```

3. Add to PATH (optional):
   ```bash
   export PATH=$PATH:$(pwd)/bin
   ```

## Usage

```bash
# Show available commands
atomizer help

# Get help for a specific command
atomizer help <command>

# Run the RFQ responder
atomizer rfq --derive-key $PRIVATE_KEY --derive-wallet $WALLET

# Analyze options markets
atomizer analyze -u ETH -e 0

# Show current positions
atomizer inventory

# List available markets
atomizer markets --underlying ETH

# Send a single quote
atomizer send-quote -i ETH-20250530-2800-C -s buy -p 100 --size 0.1
```

## Commands

### `atomizer rfq`
Runs the RFQ (Request for Quote) responder that connects to exchange WebSocket APIs and automatically provides quotes based on market conditions.

**Key features:**
- Real-time quote generation
- Automatic hedging on execution
- Multi-exchange support (Derive, Deribit)
- Configurable pricing strategies

### `atomizer analyze`
Analyzes options market liquidity and pricing across exchanges.

**Key features:**
- Liquidity scoring based on volume, trades, open interest, and spreads
- Spread analysis (absolute and percentage)
- Implied volatility comparisons
- Greeks and pricing metrics

### `atomizer inventory`
Shows current positions and P&L across configured exchanges.

### `atomizer markets`
Lists available markets and instruments with filtering options.

### `atomizer send-quote`
Utility to send individual quotes/orders for testing.

## Architecture

### Applications

The primary application that connects to Rysk Finance API, listens for RFQs, responds with quotes, and automatically hedges on Deribit.

[**Detailed Instructions**](./cmd/maker_quote_responder/README.md)

Key features:
- Real-time quote generation using Deribit prices
- Automatic hedging when trades are executed
- Configurable price premiums and slippage protection
- Comprehensive error handling and fallback mechanisms

### 2. Market Analysis Tools

#### Inventory Command (`cmd/inventory/`)
Displays current market inventory with bid/ask spreads, delta, and APY.

```bash
cd cmd/inventory && ./inventory.sh
```

#### Markets Command (`cmd/markets/`)
Fetches and saves market asset data from the Rysk API.

```bash
cd cmd/markets && ./markets.sh
```

[**Markets Tool Documentation**](./cmd/markets/README.md)

#### Send Quote Tool (`cmd/send_quote/`)
Sends RFQs to multiple markets simultaneously and measures response times.

```bash
cd cmd/send_quote && ./send_rfq.sh
```

[**Send Quote Documentation**](./cmd/send_quote/README.md)

## Submodule

The `sdk` directory is a submodule pointing to the [ryskV12-cli](https://github.com/wakamex/ryskV12-cli) repository.

To update the submodule to the latest commit from its remote:
```bash
cd sdk
git pull
cd ..
git add sdk
git commit -m "Update sdk submodule"
```
