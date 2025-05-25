# Send Quote Tool

A WebSocket-based tool for sending RFQs (Request for Quotes) to the Rysk taker endpoint and measuring response times.

## Prerequisites

- `markets.json` file containing asset addresses (automatically fetched if not present)
- Go 1.24.2 or later

## Usage

### Quick Start

```bash
./send_rfq.sh
```

### Manual Usage

```bash
./send_quote \
  -url wss://rip-testnet.rysk.finance/taker \
  -chainId 84532 \
  -quantity 1000000000000000000 \
  -taker 0x0000000000000000000000000000000000000000 \
  -markets markets.json
```

## Flags

- `-url` - WebSocket URL (default: "wss://rip-testnet.rysk.finance/taker")
- `-chainId` - Chain ID (default: 84532)
- `-quantity` - Quantity in wei (default: "1000000000000000000" = 1 token)
- `-taker` - Taker address (default: zero address for anonymous RFQs)
- `-markets` - Path to markets.json file (default: "markets.json")

## How It Works

1. **Load Markets Data**: Reads asset addresses from `markets.json` file
2. **Fetch Inventory**: Gets current available strikes and expiries from WebSocket
3. **Send RFQs**: Sends RFQs to all active markets simultaneously
4. **Measure Response Times**: Tracks time between sending RFQ and receiving quote
5. **Display Summary**: Shows response time statistics grouped by asset

## Output Example

```
📊 Loading markets from markets.json...
✅ Loaded 8 markets for chain 84532

📊 Fetching current inventory for strike/expiry data...
✅ Inventory received

🚀 Sending RFQs to all available markets simultaneously...

📤 Sent 16 RFQs, waiting for responses...

📈 RESPONSE TIME SUMMARY:
================================================================================

ETH Markets (4 responses):
  Min: 45.2ms
  Max: 52.1ms
  Avg: 48.3ms

BTC Markets (4 responses):
  Min: 43.8ms
  Max: 51.7ms
  Avg: 47.2ms

OVERALL (8/16 responded):
  Min: 43.8ms
  Max: 52.1ms
  Avg: 47.7ms
================================================================================
```

## Markets File Format

The tool expects a `markets.json` file with the following structure:

```json
{
  "84532": [
    {
      "symbol": "WETH",
      "address": "0xb67bfa7b488df4f2efa874f4e59242e9130ae61f",
      "decimals": 18,
      "chainId": 84532,
      "active": true,
      "underlying": "ETH"
    }
  ]
}
```

If the file doesn't exist, it will automatically be fetched from the Rysk API.