# Market Data Analysis Tools

Tools for analyzing correlation between exchanges and inspecting order book data from VictoriaMetrics.

## analyze_correlation.py

Analyzes price correlation between Derive and Deribit exchanges, with automatic instrument name conversion and ETH→USD price conversion.

### Examples

```bash
# Basic correlation analysis
python3 analyze_correlation.py --instrument "ETH-" --start 5m --step 10s --metric orderbook_mid_price

# Compare percentage returns instead of prices
python3 analyze_correlation.py --instrument "ETH-" --start 15m --step 5s --compare-returns

# Generate visualization plots
python3 analyze_correlation.py --instrument "ETH-" --start 30m --plot --save-plot correlation.png

# Analyze without USD conversion
python3 analyze_correlation.py --instrument "ETH-" --start 10m --no-convert-usd
```

### Key Features
- Automatic instrument name conversion (ETH-20250601-2600-C ↔ ETH-1JUN25-2600-C)
- ETH spot price lookup for USD conversion of Deribit prices
- Lead/lag analysis with statistical significance (p-values)
- Support for both ticker and orderbook metrics
- Percentage returns comparison mode

## inspect_orderbook.py

Inspects real-time and historical order book data from VictoriaMetrics.

### Examples

```bash
# View current orderbook
python3 inspect_orderbook.py --instrument "ETH-20250601-2600-C" --exchange derive

# Compare orderbooks side-by-side
python3 inspect_orderbook.py --instrument "ETH-20250601-2600-C" --exchange derive \
    --compare "ETH-1JUN25-2600-C:deribit"

# Show price history (last 30 minutes)
python3 inspect_orderbook.py --instrument "ETH-20250601-2600-C" --exchange derive --history 30

# List all available orderbooks
python3 inspect_orderbook.py --all
```

### Output Includes
- Best bid/ask prices and sizes
- Multiple orderbook levels
- Spread analysis (absolute and percentage)
- Historical price statistics
- Side-by-side exchange comparison

## check_prices.py

Quick utility to check recent spot and option prices.

```bash
python3 check_prices.py
```