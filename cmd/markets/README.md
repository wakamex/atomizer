# Markets Command

Fetches and saves market asset data from the Rysk API.

## Usage

### Using the Unified CLI (Recommended)

```bash
# List available markets
atomizer markets

# Filter by underlying asset
atomizer markets --underlying ETH

# Filter by expiry
atomizer markets --expiry 20250530

# Get help
atomizer help markets
```

### Direct Binary Usage

```bash
# Using the script
./markets.sh

# Using the binary directly
./markets -url https://rip-testnet.rysk.finance/api/assets -output markets.json

# With custom output file
OUTPUT_FILE=my-markets.json ./markets.sh
```

## Options

- `-url` - API endpoint URL (default: https://rip-testnet.rysk.finance/api/assets)
- `-output` - Output JSON file (default: markets.json)
- `-pretty` - Pretty print JSON output (default: true)

## Output

The command saves a JSON file with the following structure:

```json
{
  "chainId": [
    {
      "symbol": "WETH",
      "address": "0x...",
      "decimals": 18,
      "chainId": 84532,
      "active": true,
      "price": "...",
      "underlying": "ETH",
      "underlyingAssetAddress": "0x...",
      "minTradeSize": "...",
      "maxTradeSize": "..."
    }
  ]
}
```

## Example Output

```
📊 Fetching markets data from: https://rip-testnet.rysk.finance/api/assets
✅ Markets data saved to: markets.json

📈 MARKETS SUMMARY:
================================================================================

Chain ID 84532: 8 assets
  - WETH     (ETH): 0xb67bfa7b488df4f2efa874f4e59242e9130ae61f [ACTIVE]
  - WBTC     (BTC): 0x1234567890123456789012345678901234567890 [ACTIVE]
  - WHYPE    (HYPE): 0x2345678901234567890123456789012345678901 [ACTIVE]
  ...

================================================================================
Total chains: 3
Output file: markets.json (7.24 KB)
```