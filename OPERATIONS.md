# Atomizer Operations Guide

This guide covers deployment, configuration, trading strategies, monitoring, and troubleshooting for the Atomizer trading system.

## Table of Contents
- [Configuration](#configuration)
- [Deployment](#deployment)
- [Trading Operations](#trading-operations)
- [Market Making](#market-making)
- [Hedging Strategies](#hedging-strategies)
- [Monitoring & Analysis](#monitoring--analysis)
- [Troubleshooting](#troubleshooting)
- [Performance Tuning](#performance-tuning)

## Configuration

### Environment Variables

#### Core Configuration
```bash
# Exchange Selection
EXCHANGE_NAME=derive          # derive, deribit, or ccxt
EXCHANGE_TEST_MODE=false      # Use testnet/sandbox

# Authentication - Derive
DERIVE_PRIVATE_KEY=your_private_key_hex
DERIVE_WALLET_ADDRESS=0x...

# Authentication - Deribit
DERIBIT_API_KEY=your_api_key
DERIBIT_API_SECRET=your_api_secret
# Alternative: Ed25519 authentication
DERIBIT_USE_ED25519=true
DERIBIT_PRIVATE_KEY_ED25519=your_ed25519_private_key

# RFQ Configuration
RFQ_ASSET_ADDRESSES=0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9
WS_URL=wss://api.lyra.finance/ws/v2
QUOTE_VALID_DURATION=30       # seconds
```

#### Risk Management
```bash
# Position Limits
MAX_POSITION_DELTA=10.0       # Maximum delta exposure
MAX_POSITION_SIZE=100.0       # Maximum size per instrument
MIN_LIQUIDITY_SCORE=0.5       # Minimum liquidity for trading

# Hedging Configuration  
ENABLE_GAMMA_HEDGING=true
GAMMA_THRESHOLD=0.5           # Gamma exposure threshold
HEDGE_RATIO=1.0               # 1.0 = full hedge
```

### Asset Mapping

Configure token address to symbol mapping in `asset_mapping.json`:

```json
{
  "0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9": "ETH",
  "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2": "ETH",
  "0x2260FAC5E5542a773Aa44fBCfeDf7C193bc2C599": "BTC"
}
```

### Configuration Files

Create a `config.yaml` for complex configurations:

```yaml
exchange:
  name: derive
  test_mode: false
  
market_maker:
  instruments:
    - underlying: ETH
      expiry: 20250530
      strikes: [2800, 3000, 3200]
  spread: 50                  # basis points
  min_spread: 10             
  size: 0.1
  improvement: 0.05           # For aggressive mode (aggression >= 1.0)
  aggression: 0.7             # 0-0.9: conservative, 1.0+: aggressive
  max_position: 10.0
  
risk:
  max_delta: 10.0
  max_gamma: 5.0
  max_vega: 1000.0
  check_interval: 30s
```

## Deployment

### Local Development
```bash
# Clone and install
git clone https://github.com/wakamex/atomizer.git
cd atomizer
go install ./cmd/atomizer

# Run with environment file
source .env.local
atomizer rfq-responder
```

### Docker Deployment
```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o atomizer ./cmd/atomizer

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/atomizer /usr/local/bin/
CMD ["atomizer", "rfq-responder"]
```

### Production Deployment

#### Systemd Service
```ini
[Unit]
Description=Atomizer Trading System
After=network.target

[Service]
Type=simple
User=atomizer
WorkingDirectory=/opt/atomizer
EnvironmentFile=/etc/atomizer/atomizer.env
ExecStart=/opt/atomizer/bin/atomizer rfq-responder
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

#### High Availability Setup
- Run multiple instances behind a load balancer
- Use Redis for shared state
- Configure leader election for market maker
- Set up database replication

## Trading Operations

### Understanding the Aggression Parameter

The `--aggression` parameter is a unified control for market making behavior:

**Conservative Mode (0.0 - 0.9):**
- Places orders on your side of the spread only
- Formula for bids: `bid_price = best_bid + aggression × (mid - best_bid)`
- Formula for asks: `ask_price = best_ask - aggression × (best_ask - mid)`
- Examples:
  - `0.0`: Join the existing best bid/ask (passive)
  - `0.5`: Place orders halfway between best and mid
  - `0.9`: Place orders very close to mid (aggressive but safe)

**Aggressive Mode (1.0+):**
- Can cross the spread to capture more flow
- Uses the `--improvement` parameter to determine cross amount
- `1.0`: Default aggressive behavior
- `>1.0`: Future use for even more aggressive strategies

```bash
# Examples
atomizer market-maker --aggression 0.0    # Join best bid/ask
atomizer market-maker --aggression 0.5    # Halfway to mid
atomizer market-maker --aggression 0.9    # Near mid, max conservative
atomizer market-maker --aggression 1.0    # Cross spread (default)
```

### Starting the RFQ Responder
```bash
# Basic operation
atomizer rfq-responder

# With specific configuration
atomizer rfq-responder \
  --exchange derive \
  --dummy-price 100 \
  --enable-hedging
```

### Manual Order Placement
```bash
# Place a single order
atomizer manual-order \
  --instrument ETH-PERP \
  --side buy \
  --price 3500 \
  --amount 0.1

# Using environment variables
export ORDER_INSTRUMENT=ETH-20250530-3000-C
export ORDER_SIDE=sell
export ORDER_PRICE=150
export ORDER_AMOUNT=1.0
./scripts/manual_order.sh
```

### Position Management
```bash
# View current positions
atomizer inventory

# Export positions to CSV
atomizer inventory --format csv > positions.csv

# Monitor P&L in real-time
watch -n 5 atomizer inventory --summary
```

## Market Making

### Basic Market Making
```bash
# Conservative mode - join best bid/ask
atomizer market-maker \
  --expiry 20250530 \
  --strikes 3000 \
  --size 0.1 \
  --aggression 0.0

# Conservative mode - halfway to mid
atomizer market-maker \
  --expiry 20250530 \
  --strikes 3000 \
  --size 0.1 \
  --aggression 0.5

# Aggressive mode - cross spread (default)
atomizer market-maker \
  --expiry 20250530 \
  --strikes 2800,3000,3200 \
  --size 0.1 \
  --aggression 1.0
```

### Market Making Strategies

The `--aggression` parameter controls quote placement:

**Conservative Mode (aggression < 1.0):**
- Keeps quotes on your side of the spread
- `0.0`: Join the best bid/ask (least aggressive)
- `0.5`: Place orders halfway between best and mid
- `0.9`: Place orders very close to mid (most aggressive while staying on your side)

**Aggressive Mode (aggression >= 1.0):**
- Traditional mode where quotes can cross the spread
- Uses `--improvement` parameter to tighten spreads
- Can place bids above mid or asks below mid
- Higher fill rates but more adverse selection risk

### Advanced Strategies

#### Volatility-Adjusted Spreads
The market maker automatically adjusts spreads based on:
- Current implied volatility
- Recent price movements  
- Order book depth
- Time to expiry

#### Position-Based Pricing
```bash
# Conservative with position limits
atomizer market-maker \
  --aggression 0.3 \
  --max-position 10.0 \
  --size 0.5

# Aggressive with larger positions
atomizer market-maker \
  --aggression 1.2 \
  --max-position 50.0 \
  --improvement 0.2
```

#### Smart Order Management
- Replaces orders only when necessary
- Minimizes order cancellations
- Groups updates to reduce messages
- Maintains queue position when possible

## Hedging Strategies

### Delta Hedging
Automatically hedges delta exposure from options trades:

```bash
# Configure delta hedging
export ENABLE_DELTA_HEDGE=true
export DELTA_THRESHOLD=0.1    # Minimum delta before hedging
export HEDGE_VENUE=deribit    # Where to hedge
```

**Strategy Details:**
- Monitors net delta across all positions
- Hedges when threshold exceeded
- Uses spot/futures for linear hedging
- Considers transaction costs

### Gamma Hedging
For large positions or volatile markets:

```bash
# Enable gamma hedging
export ENABLE_GAMMA_HEDGING=true
export GAMMA_THRESHOLD=0.5
export GAMMA_HEDGE_RATIO=0.8  # Partial hedge
```

**When to Use:**
- Large option positions
- Near expiry (high gamma)
- Volatile market conditions
- Path-dependent strategies

### Cross-Exchange Hedging
```bash
# Configure multi-venue hedging
atomizer rfq-responder \
  --primary-exchange derive \
  --hedge-exchange deribit \
  --enable-cross-hedge
```

## Monitoring & Analysis

### Real-Time Monitoring

#### Market Data Collection
```bash
# Start market monitor
atomizer market-monitor start \
  --exchanges derive,deribit \
  --instruments ETH-*-C \
  --orderbook \
  --depth 10

# With VictoriaMetrics storage
atomizer market-monitor setup  # First time only
atomizer market-monitor start --store-metrics
```

**Market Monitor Features**:
- Automatic USD conversion for Deribit (display only, stores ETH values)
- Instrument name conversion between exchanges
- File-based caching with configurable TTL
- Real-time WebSocket data collection

#### RFQ Benchmarking
```bash
# Send RFQs to all markets simultaneously
atomizer send-quote \
  --instrument ETH-20250530-3000-C \
  --side buy \
  --size 1.0 \
  --measure-response-time

# Batch RFQ testing
atomizer send-quote \
  --from-inventory \
  --side buy \
  --size 0.1 \
  --summary
```

#### Performance Metrics
Monitor key metrics via Prometheus/Grafana:
- Quote response time
- Fill rate
- P&L by strategy
- Position Greeks
- System latency

### Market Analysis

#### Liquidity Analysis
```bash
# Analyze specific expiry with liquidity scoring
atomizer analyze \
  --underlying ETH \
  --expiry 20250530 \
  --min-moneyness 0.8 \
  --max-moneyness 1.2

# Compare venues
atomizer analyze \
  --underlying ETH \
  --compare-exchanges

# Query ETH calls with real-time data
atomizer analyze query-eth-calls \
  --min-volume 10 \
  --min-liquidity 0.5 \
  --export calls_analysis.csv

# Show near-term options
atomizer analyze nearterm \
  --days 7 \
  --underlying ETH
```

**Liquidity Scoring**:
- Based on: volume, trades, open interest, spreads
- Score range: 0.0 (no liquidity) to 1.0 (highly liquid)
- Includes IV comparison across venues

#### Volatility Surface
```bash
# Export volatility data
atomizer analyze \
  --underlying ETH \
  --export-surface \
  --format csv > vol_surface.csv
```

#### Cross-Exchange Correlation Analysis

Analyze whether market makers are quoting off each other:

```bash
# Prerequisites
cd scripts/market_analysis
pip install -r requirements.txt

# Analyze bid price correlation between exchanges
python analyze_correlation.py \
  --start 24h \
  --metric market_bid_price \
  --plot

# Check spread correlation
python analyze_correlation.py \
  --metric market_spread_percent \
  --instruments "ETH-20250628-3000-C" \
  --lag 5
```

This helps identify:
- Price leadership between venues
- Arbitrage opportunities
- Market efficiency metrics

#### Querying Stored Data

Use PromQL to analyze collected market data:

```bash
# Average spread over time
avg_over_time(market_spread{exchange="derive", instrument="ETH-20250628-3000-C"}[1h])

# Compare bid prices between exchanges
market_bid_price{instrument="ETH-20250628-3000-C"}

# Find instruments with high spreads
topk(10, market_spread_percent)

# Export data for external analysis
curl 'http://localhost:8428/api/v1/export' \
  -d 'match[]={__name__=~"market_.*",instrument="ETH-20250628-3000-C"}' \
  -d 'start=2024-01-01T00:00:00Z' \
  -d 'end=2024-01-02T00:00:00Z' > market_data.jsonl
```

### Trade Analysis

#### P&L Attribution
```bash
# Daily P&L report
atomizer reports pnl \
  --date 2024-01-15 \
  --breakdown-by instrument,strategy

# Historical analysis
atomizer reports pnl \
  --start 2024-01-01 \
  --end 2024-01-31 \
  --export pnl_january.csv
```

## Troubleshooting

### Common Issues

#### WebSocket Connection Issues
```bash
# Debug WebSocket connections
atomizer rfq-responder --debug --log-level trace

# Test connectivity
atomizer test-connection --exchange derive

# Monitor reconnections
tail -f logs/atomizer.log | grep -i "websocket\|reconnect"

# Monitor Derive orders in real-time
./scripts/monitor_derive_orders.sh

# Test specific Derive order
./scripts/test_derive_order.sh
```

#### Quote Generation Failures
1. Check exchange connectivity
2. Verify market data feed
3. Confirm authentication
4. Review error logs

```bash
# Diagnostic mode
atomizer rfq-responder --dry-run --verbose
```

#### Order Rejection
Common causes:
- Insufficient margin
- Price outside limits
- Position limits exceeded
- Invalid instrument

```bash
# Test order placement
atomizer test-order \
  --instrument ETH-20250530-3000-C \
  --side buy \
  --price 100 \
  --amount 0.1
```

#### Aggression Issues

**Orders not filling:**
- Try increasing aggression gradually
- Check if spread is too wide
- Verify market has liquidity

```bash
# Debug aggression settings
atomizer market-maker \
  --aggression 0.0 \
  --size 0.01 \
  --dry-run

# Monitor actual quote placement
atomizer market-maker \
  --aggression 0.5 \
  --debug \
  --log-level trace
```

**Getting adverse selection:**
- Reduce aggression below 1.0
- Use conservative mode (0.3-0.7)
- Increase minimum spread

### Log Analysis

#### Log Locations
```
logs/
├── atomizer.log          # Main application log
├── trades.log            # Trade execution log
├── quotes.log            # Quote generation log
└── errors.log            # Error-only log
```

#### Useful Log Queries
```bash
# Find failed quotes
grep "quote.*failed" logs/atomizer.log

# Track specific RFQ
grep "rfq-id-12345" logs/atomizer.log

# Monitor hedge trades
tail -f logs/trades.log | grep -i hedge
```

### Health Checks

#### System Health
```bash
# Check all components
atomizer health --verbose

# Specific component
atomizer health --component rfq-processor
```

#### Exchange Health
```bash
# Test exchange APIs
atomizer test-exchange --all

# Specific exchange
atomizer test-exchange --name deribit --full
```

## Performance Tuning

### Latency Optimization

#### Quote Generation
- Pre-calculate common values
- Cache market data (100ms TTL)
- Parallel exchange queries
- Connection pooling

```bash
# Enable performance mode
atomizer rfq-responder \
  --performance-mode \
  --cache-ttl 100ms \
  --parallel-quotes 4
```

#### Aggression Tuning for Performance
Different aggression levels have different performance characteristics:

```bash
# High-frequency, low risk (many small trades)
atomizer market-maker \
  --aggression 0.1 \
  --size 0.1 \
  --refresh 1

# Balanced approach
atomizer market-maker \
  --aggression 0.5 \
  --size 0.5 \
  --refresh 3

# Aggressive, fewer but larger trades
atomizer market-maker \
  --aggression 1.0 \
  --size 2.0 \
  --refresh 5 \
  --improvement 0.2
```

#### Network Optimization
```bash
# TCP tuning (Linux)
echo 'net.core.rmem_max = 134217728' >> /etc/sysctl.conf
echo 'net.core.wmem_max = 134217728' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_rmem = 4096 87380 134217728' >> /etc/sysctl.conf
echo 'net.ipv4.tcp_wmem = 4096 65536 134217728' >> /etc/sysctl.conf
sysctl -p
```

### Resource Management

#### Memory Usage
```bash
# Set memory limits
export GOGC=100              # Default GC target
export GOMEMLIMIT=4GiB       # Memory limit

# Monitor memory
atomizer debug memory --interval 60s
```

#### CPU Optimization
```bash
# Set CPU affinity
taskset -c 0-3 atomizer market-maker

# Profile CPU usage
atomizer debug profile --cpu --duration 30s
```

### Scaling Strategies

#### Horizontal Scaling
- RFQ responders: Stateless, scale freely
- Market makers: Coordinate via Redis
- Risk managers: Single instance with backup
- Data collectors: Partition by instrument

#### Load Balancing
```nginx
upstream atomizer_rfq {
    least_conn;
    server atomizer1:8080;
    server atomizer2:8080;
    server atomizer3:8080;
}
```

### Monitoring Performance

#### Key Metrics
- Quote latency (p50, p95, p99)
- WebSocket message rate
- Order placement time
- Cache hit ratio
- GC pause time

#### Alerting Thresholds
```yaml
alerts:
  - name: high_quote_latency
    condition: quote_latency_p99 > 100ms
    severity: warning
    
  - name: low_fill_rate
    condition: fill_rate < 0.8
    severity: critical
    
  - name: position_limit
    condition: position_delta > max_delta * 0.9
    severity: warning
```