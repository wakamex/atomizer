# Atomizer Project

This repository contains the Atomizer project, a comprehensive market making system for Rysk Finance with automated hedging capabilities.
It utilizes the `ryskV12-cli` as a submodule located in the `sdk` directory.

## Features

- **Automated Market Making**: Connects to Rysk Finance WebSocket API to receive and respond to RFQs
- **Deribit Hedging**: Automatically hedges positions on Deribit when trades are executed
- **Real-time Pricing**: Fetches live option prices from Deribit for competitive quoting
- **Market Analysis Tools**: Command-line utilities for market inventory, RFQ testing, and performance analysis

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


## Prerequisites

- Git
- Go (version as specified in `examples/go.mod` and `sdk/go.mod`)

## Setup

1. Clone the repository:
   ```bash
   git clone <repository_url> atomizer
   cd atomizer
   ```

2. Initialize and update the submodule:
   ```bash
   git submodule update --init --recursive
   ```

## Applications

### 1. Maker Quote Responder (`cmd/maker_quote_responder/`)

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
