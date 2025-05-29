# Analyze Options Tool

A command-line tool for analyzing Derive options market data with caching support.

## Features

- **Market Overview**: Display total options count and active percentage
- **Strike Distribution**: Statistical analysis of strike prices by underlying asset
- **Active Options by Expiry**: Visual breakdown of active vs inactive options
- **Export Functionality**: Export filtered options to CSV format
- **Cache Support**: Uses file-based caching for efficient data retrieval

## Usage

```bash
# Build the tool
cd cmd/analyze_options
go build

# Run analysis
./analyze_options

# Examples:
./analyze_options stats        # Show active options stats
./analyze_options export 2     # Export ETH calls expiring in 2 days
./analyze_options query        # Query ETH calls at nearest expiry with liquidity analysis
./analyze_options nearterm 7   # Show options expiring in next 7 days
```

## Commands

1. **Show Strike Distribution** - Displays mean, std dev, and percentiles for each underlying asset
2. **Show Active by Expiry** - Shows active percentage for each expiry date with visual indicators
3. **Export Options** - Export filtered options to CSV with customizable criteria:
   - Underlying asset (BTC/ETH/etc)
   - Option type (call/put)
   - Days to expiry
   - Active status filter
4. **Query ETH Calls** - Fetches real-time ticker data for ETH call options at the nearest expiry:
   - Automatically selects the nearest expiry date with available options
   - Calculates liquidity scores based on volume, trades, open interest, order book depth, and spreads
   - Shows top 10 most liquid options with pricing and Greeks
   - Displays IV comparison (API IV vs Bid/Ask IV)
   - Exports results to CSV with detailed metrics

## Example Output

```
=== Options Market Overview ===
Total Options: 1322
Active Options: 968 (73.2%)

=== Strike Distribution by Underlying ===
BTC: Mean=144998, StdDev=84658, P25=91000, P50=110000, P75=185000
ETH: Mean=4472, StdDev=3513, P25=2600, P50=3400, P75=5000
```

## Configuration

The tool uses the existing market cache configuration from the main application. Cache files are stored in `cache/derive_markets.json` by default.

## Dependencies

- Go 1.19+
- Existing atomizer market infrastructure
- File-based cache system