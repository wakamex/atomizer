# Data Analysis: Derive vs Deribit Order Book Correlation

This document explains how to analyze whether Derive market makers are quoting off the Deribit order book.

## Background

We collect real-time order book data from both Derive and Deribit exchanges using the market monitor, storing it in VictoriaMetrics. This data includes:

- `market_bid_price` - Best bid prices
- `market_ask_price` - Best ask prices  
- `market_bid_size` - Best bid sizes
- `market_ask_size` - Best ask sizes
- `market_spread` - Absolute spread
- `market_spread_percent` - Spread as percentage

## Running the Correlation Analysis

### Prerequisites

1. **Market Monitor Running**: Ensure the market monitor is collecting data:
   ```bash
   atomizer monitor --derive --deribit --instruments "ETH-*"
   ```

2. **VictoriaMetrics Accessible**: Default URL is `http://localhost:8428`

3. **Python Dependencies**:
   ```bash
   cd /code/atomizer/cmd/analyze_correlation
   pip install -r requirements.txt
   ```

### Basic Analysis

```bash
cd /code/atomizer/cmd/analyze_correlation

# Analyze bid prices for the last hour
python analyze_correlation.py
```

### Common Analysis Scenarios

#### 1. Check if Derive follows Deribit prices

```bash
# Analyze 24 hours of bid prices with visualization
python analyze_correlation.py --start 24h --metric market_bid_price --plot
```

High correlation (>0.95) indicates Derive is likely using Deribit as reference.

#### 2. Analyze lead/lag relationship

```bash
# Check which exchange leads price movements
python analyze_correlation.py --start 6h --specific-instrument ETH-20250601-2600-C
```

Positive lag means Deribit leads (Derive follows).

#### 3. Compare spreads

```bash
# See if spreads are correlated
python analyze_correlation.py --metric market_spread_percent --start 12h --plot
```

#### 4. Full analysis for specific expiry

```bash
# Comprehensive analysis with saved plots
python analyze_correlation.py \
  --instrument ETH-20250601- \
  --start 24h \
  --metric market_bid_price \
  --save-plot eth_june_correlation.png
```

## Interpreting Results

### Correlation Values
- **>0.95**: Very high correlation - strong evidence of reference pricing
- **0.8-0.95**: High correlation - likely some reference pricing
- **0.5-0.8**: Moderate correlation - partial influence
- **<0.5**: Low correlation - independent pricing

### Lead/Lag Analysis
- **Positive lag**: Deribit leads by N steps (Derive follows)
- **Negative lag**: Derive leads by N steps (rare)
- **Zero lag**: Simultaneous price movements

### Price Differences
- **Positive mean diff**: Derive prices higher than Deribit
- **Negative mean diff**: Derive prices lower than Deribit
- **High std dev**: More price dispersion between exchanges

## Example Output

```
Instrument: ETH-20250601-2600-C
  Data points: 360
  Correlation: 0.9842 (p-value: 1.23e-287)
  Mean difference: 0.8521 (2.41%)
  Std deviation of diff: 0.3214
  Best lag: 3 steps (correlation: 0.9856)
    -> Deribit leads Derive by 3 steps
  Derive: mean=36.2143, std=2.1532
  Deribit: mean=35.3622, std=2.0891
```

This indicates:
- 98.42% correlation - Derive strongly follows Deribit
- Derive prices average 85¢ higher (2.41% premium)
- Deribit price changes appear 30 seconds before Derive (3 × 10s steps)

## Automated Analysis

Run all metrics for an instrument:

```bash
#!/bin/bash
INSTRUMENT="ETH-20250601-2600-C"

for metric in market_bid_price market_ask_price market_spread_percent; do
  echo "=== Analyzing $metric ==="
  python analyze_correlation.py \
    --metric $metric \
    --specific-instrument $INSTRUMENT \
    --start 6h
  echo
done
```

## Visualization Features

The `--plot` option generates 4 subplots:

1. **Time Series**: Overlaid prices from both exchanges
2. **Scatter Plot**: Direct price correlation visualization
3. **Price Difference Histogram**: Distribution of price gaps
4. **Rolling Correlation**: How correlation changes over time

## Advanced Usage

### Custom Time Ranges
```bash
# Specific date range
python analyze_correlation.py --start "2024-05-30T12:00:00"

# Different granularity
python analyze_correlation.py --step 1m  # 1-minute intervals
```

### Batch Analysis
```bash
# Analyze all June expiries
python analyze_correlation.py --instrument "ETH-202506"
```

### Export Results
```bash
# Save detailed results (redirect output)
python analyze_correlation.py --start 24h > correlation_report.txt
```

## Key Findings to Look For

1. **Market Making Strategy**:
   - Correlation >0.95 across multiple instruments = systematic Deribit reference
   - Consistent positive lag = reactive pricing strategy
   - Systematic price offset = markup/discount strategy

2. **Arbitrage Opportunities**:
   - Low correlation periods = pricing discrepancies
   - Large price differences = potential arbitrage
   - Variable lag = unstable reference pricing

3. **Market Efficiency**:
   - Decreasing lag over time = improving efficiency
   - Tightening spreads = competitive market making
   - Correlation breakdown = independent price discovery

## Troubleshooting

### No Data Found
```bash
# Check if monitor is running and collecting data
curl -s "http://localhost:8428/api/v1/query?query=market_bid_price" | jq .

# Verify specific instruments
curl -s "http://localhost:8428/api/v1/label/__name__/values" | jq .
```

### Connection Errors
```bash
# Specify VictoriaMetrics URL
python analyze_correlation.py --vm-url http://your-server:8428
```

### Missing Instruments
Ensure the market monitor is configured to collect both exchanges:
```bash
atomizer monitor --derive --deribit --instruments "ETH-*"
```

## Additional Tools

### Orderbook Inspector
View real-time orderbook data:
```bash
cd cmd/analyze_correlation
python3 inspect_orderbook.py --instrument "ETH-20250601-2600-C" --exchange derive
```

Compare orderbooks between exchanges:
```bash
python3 inspect_orderbook.py --instrument "ETH-20250601-2600-C" --exchange derive \
    --compare "ETH-1JUN25-2600-C:deribit"
```

See [cmd/analyze_correlation/README.md](cmd/analyze_correlation/README.md) for more examples.

## Next Steps

1. Start collecting data from both exchanges
2. Let it run for at least 30 minutes to build up history  
3. Run correlation analysis during active trading hours
4. Experiment with different time windows and instruments