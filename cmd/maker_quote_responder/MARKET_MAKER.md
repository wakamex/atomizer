# Market Maker Documentation

The market maker is a subcommand of the `maker_quote_responder` binary that provides automated market making for options on Derive and Deribit exchanges.

## Overview

The market maker continuously posts bid and ask orders for specified options, improving market liquidity by tightening spreads. It monitors real-time market data via WebSocket connections and updates quotes based on current market conditions.

## Usage

```bash
./maker_quote_responder market-maker [flags]
```

## Architecture

```
MarketMaker
    ├── MarketMakerInterface       (Generic exchange interface)
    ├── DeriveMarketMakerExchange  (Derive implementation)
    ├── MarketMakerConfig          (Configuration parameters)
    ├── Order Tracking             (Active order management)
    └── Risk Management            (Position and exposure limits)
```

## Command Line Flags

### Required Flags

- **`-expiry string`**: Expiry date to make markets on (format: YYYYMMDD, e.g., "20250606")
  - Must specify an exact expiry date
  - The market maker will only trade options with this expiry

- **Strike Selection** (one required):
  - **`-strikes string`**: Comma-separated list of strikes to trade (e.g., "2800,3000,3200")
  - **`-all-strikes`**: Make markets on all available strikes for the expiry
  - Must specify either `-strikes` or `-all-strikes`

### Exchange Configuration

- **`-exchange string`** (default: "derive"): Exchange to use ("derive" or "deribit")
- **`-test`**: Use exchange testnet instead of mainnet
- **`-underlying string`** (default: "ETH"): Underlying asset ("ETH" or "BTC")

### Market Making Parameters

- **`-spread int`** (default: 10): Target spread in basis points (100 = 1%)
  - Used when reference size filtering falls back to mid price
  - Calculates: `spread_amount = mid_price * spread_bps / 10000`

- **`-size float`** (default: 0.1): Quote size in contracts
  - Number of contracts to quote on each side
  - Same size used for both bid and ask

- **`-refresh int`** (default: 1): Quote refresh interval in seconds
  - How often to cancel and replace quotes
  - Lower values = more responsive to market changes

- **`-min-spread int`** (default: 1000): Minimum spread in basis points (10%)
  - Enforces minimum distance between bid and ask
  - Prevents quotes from being too tight

- **`-improvement float`** (default: 0.1): Amount to improve quotes by
  - How much better than the reference price to quote
  - Bid: reference bid + improvement
  - Ask: reference ask - improvement

- **`-improvement-reference-size float`** (default: 0): Minimum size for best bid/ask selection
  - Only improve against orders with at least this size
  - 0 = use any size (default behavior)
  - When set, fetches full order book to find qualifying levels
  - Falls back to mid price ± spread if no qualifying orders found

### Risk Management

- **`-max-position float`** (default: 1.0): Maximum position per instrument
  - Stops quoting when position exceeds this limit
  - Helps prevent excessive exposure to single strikes

- **`-max-exposure float`** (default: 10.0): Maximum total exposure across all positions
  - Global risk limit across all instruments
  - Calculated as sum of absolute position values

### Authentication

- **`-private-key string`**: Private key for signing orders
  - Overrides environment variables
  - If not provided, uses `DERIVE_PRIVATE_KEY` or `DERIBIT_PRIVATE_KEY`

- **`-wallet string`**: Wallet address (Derive only)
  - Required for Derive exchange
  - If not provided, uses `DERIVE_WALLET_ADDRESS` environment variable

### Derive Specific Configuration

For Derive, the following environment variables are required:
- `DERIVE_PRIVATE_KEY`: Ethereum private key for signing orders
- `DERIVE_WALLET_ADDRESS`: Your Derive wallet address
- `DERIVE_SUBACCOUNT_ID`: The subaccount ID to use (e.g., 50400)

### Utility

- **`-dry-run`**: Print configuration without starting the market maker
  - Useful for verifying parameters
  - Shows all instruments that would be traded

## Pricing Strategy

The market maker uses a configurable improvement strategy with optional size filtering:

### Basic Improvement Strategy
- **Bid Price**: Reference bid + improvement amount
- **Ask Price**: Reference ask - improvement amount

### Reference Price Selection
1. **Without size filter** (`-improvement-reference-size 0`):
   - Uses current best bid/ask as reference prices
   
2. **With size filter** (`-improvement-reference-size > 0`):
   - Fetches full order book
   - Finds first bid/ask level with size ≥ reference size
   - If no qualifying levels found, falls back to mid price ± spread

### Example Pricing

**Example 1: Basic improvement (no size filter)**
If the current market is:
- Best Bid: $100.00 (size: 0.5)
- Best Ask: $105.00 (size: 0.3)
- Improvement: $0.10

The market maker will quote:
- New Bid: $100.10
- New Ask: $104.90

**Example 2: With size filter**
If `-improvement-reference-size 1.0` and the order book is:
```
Bids:                    Asks:
$100.00 (size: 0.5)     $105.00 (size: 0.3)
$99.50  (size: 1.2) ←   $105.50 (size: 1.5) ←
$99.00  (size: 2.0)     $106.00 (size: 2.0)
```

The market maker will:
- Skip the best bid/ask (sizes < 1.0)
- Use $99.50 as reference bid (first bid ≥ 1.0)
- Use $105.50 as reference ask (first ask ≥ 1.0)
- Quote: Bid $99.60, Ask $105.40

### Minimum Spread Protection

If the calculated quotes violate the minimum spread requirement:
- The market maker widens the quotes symmetrically around the mid price
- Ensures the spread is at least `-min-spread` basis points

## Examples

### Basic Market Making

Make markets on specific ETH call strikes expiring June 1, 2025:

```bash
./maker_quote_responder market-maker \
  -expiry 20250601 \
  -strikes 2600,2800,3000 \
  -size 0.1
```

### All Strikes with Tight Spreads

Make markets on all available strikes with 5 bps minimum spread:

```bash
./maker_quote_responder market-maker \
  -expiry 20250606 \
  -all-strikes \
  -min-spread 5 \
  -size 1.0
```

### Risk-Limited Trading

Trade with position and exposure limits:

```bash
./maker_quote_responder market-maker \
  -expiry 20250606 \
  -strikes 3000,3200,3400 \
  -max-position 5.0 \
  -max-exposure 50.0 \
  -size 0.5
```

### Custom Improvement Strategy

Improve by $0.50 only against orders with size ≥ 5.0:

```bash
./maker_quote_responder market-maker \
  -expiry 20250606 \
  -strikes 3000 \
  -improvement 0.5 \
  -improvement-reference-size 5.0
```

### Avoid Small Orders

Only quote against substantial liquidity:

```bash
./maker_quote_responder market-maker \
  -expiry 20250606 \
  -all-strikes \
  -improvement-reference-size 10.0 \
  -size 2.0
```

### Testnet Trading

Test on Derive testnet:

```bash
./maker_quote_responder market-maker \
  -test \
  -expiry 20250606 \
  -strikes 3000 \
  -dry-run
```

### Using Environment Variables

Set credentials via environment:

```bash
export DERIVE_PRIVATE_KEY="your-private-key"
export DERIVE_WALLET_ADDRESS="your-wallet-address"
export DERIVE_SUBACCOUNT_ID=50400

./maker_quote_responder market-maker \
  -expiry 20250601 \
  -strikes 2600,2800,3000
```

### Using the Shell Script

A convenience script with environment variable configuration:

```bash
# Basic usage with defaults
./run_market_maker.sh

# Custom configuration
EXPIRY=20250613 STRIKES="2800,2900,3000" SIZE=2.0 ./run_market_maker.sh

# All strikes with wider spread
ALL_STRIKES=true SPREAD_BPS=20 ./run_market_maker.sh

# Conservative settings
MAX_POSITION=5.0 MAX_EXPOSURE=50.0 SIZE=0.5 ./run_market_maker.sh

# Custom improvement strategy
IMPROVEMENT=0.5 IMPROVEMENT_REFERENCE_SIZE=10.0 ./run_market_maker.sh

# Dry run to see configuration
DRY_RUN=true ./run_market_maker.sh
```

## Risk Management

The market maker includes several risk controls:

1. **Position Limits**: Stops quoting when position exceeds `-max-position`
2. **Exposure Limits**: Stops all quoting when total exposure exceeds `-max-exposure`
3. **Minimum Spread**: Enforces minimum distance between bid and ask
4. **Order Tracking**: Monitors all open orders and filled positions

## Order Lifecycle

1. **Quote Calculation**: Based on current market ticker data
2. **Risk Check**: Verifies position and exposure limits
3. **Order Cancellation**: Cancels existing orders for the instrument
4. **Order Placement**: Places new bid and ask orders
5. **Order Tracking**: Stores order details for monitoring
6. **Fill Detection**: Updates positions when orders are filled

## Monitoring

The market maker logs:
- Quote updates and order placements
- Fill notifications
- Risk limit warnings
- Connection status
- Error conditions
- Statistics every 30 seconds

Monitor logs to ensure proper operation:

```bash
./maker_quote_responder market-maker -expiry 20250606 -strikes 3000 2>&1 | tee market_maker.log
```

## Safety Features

1. **Graceful Shutdown**: Ctrl+C cancels all orders before exiting
2. **Position Tracking**: Monitors net exposure in real-time
3. **Order Validation**: Checks risk limits before placing orders
4. **Connection Monitoring**: Handles WebSocket disconnections gracefully
5. **Atomic Order Updates**: Ensures bid/ask pairs are updated together
6. **Smart Order Updates**: Only cancels/replaces orders if market moves > 0.5% (CancelThreshold)

## Technical Implementation Details

### Derive Order Signing

The market maker uses EIP-712 signing for Derive orders:

1. **Fetch Instrument Details**: Gets base asset address and sub ID
2. **Create Action**: Builds order action with all parameters
3. **EIP-712 Signing**: Signs with domain separator and action typehash
4. **Submit Order**: Sends signed order via WebSocket

Key parameters for Derive orders:
- `subaccount_id`: From DERIVE_SUBACCOUNT_ID env var
- `max_fee`: Set to 100 USDC (in wei)
- `signature_expiry_sec`: 1 hour from order creation
- `nonce`: Unique value using timestamp + nanoseconds
- `signer`: EOA address that signs the order

### WebSocket Management

- Uses mutex protection for concurrent writes
- Maintains separate connections for orders and ticker data
- Implements heartbeat/ping-pong for connection health
- Automatic reconnection on disconnection

## Future Improvements

1. **Dynamic Spread Calculation**: Use the `-spread` parameter to calculate spread based on:
   - Market volatility
   - Current position
   - Time to expiry
   - Order book depth

2. **Advanced Quoting Strategies**: 
   - Quote around theoretical value instead of market prices
   - Skew quotes based on current position (wider on the side we're long)
   - Adjust spread based on market volatility
   - Time-weighted average price (TWAP) execution

3. **Greeks Management**: 
   - Track delta, gamma, vega, theta exposure
   - Set limits on greek exposures
   - Dynamic hedging strategies

4. **P&L Tracking**: 
   - Real-time profit/loss calculation
   - Mark-to-market positions
   - Fee tracking and reporting

5. **Multi-Expiry Support**: 
   - Trade multiple expiries simultaneously
   - Cross-expiry arbitrage detection
   - Calendar spread strategies

6. **Market Impact Analysis**:
   - Track how our quotes affect the market
   - Measure price improvement provided
   - Optimize quote sizes based on fill rates

## Troubleshooting

### Common Issues

1. **"No ticker data" errors**:
   - Ensure the instrument exists and is active
   - Check WebSocket connection status
   - Verify the expiry date is correct

2. **Orders not placing**:
   - Check authentication credentials
   - Verify wallet has sufficient balance
   - Ensure risk limits aren't exceeded

3. **Quotes not updating**:
   - Check the refresh interval setting
   - Monitor WebSocket ticker updates
   - Verify market data is being received

### Debug Mode

Run with verbose logging:

```bash
RUST_LOG=debug ./maker_quote_responder market-maker -expiry 20250606 -strikes 3000
```

## Performance Considerations

- Each instrument requires ~2 WebSocket subscriptions (ticker + orders)
- Quote updates happen every `-refresh` seconds
- Consider network latency when setting refresh intervals
- Use `-all-strikes` carefully as it may create many subscriptions