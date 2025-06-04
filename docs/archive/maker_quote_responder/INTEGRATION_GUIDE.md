# Integration Guide for Arbitrage Bot Enhancements

This guide shows how to integrate the new arbitrage components into the existing maker_quote_responder without breaking current functionality.

## Overview

The enhancements add:
- Arbitrage orchestrator for coordinating trades
- HTTP API for manual trade submission
- Risk management with position and delta limits
- Gamma dynamic delta hedging
- Async trade processing
- Valkey cache support

## Integration Steps

### 1. Environment Variables

Add these to your `.env` file:

```bash
# Risk limits
MAX_POSITION_SIZE=1000
MAX_DELTA_EXPOSURE=500
GAMMA_MAX_DELTA=100

# Features
ENABLE_MANUAL_TRADES=true
GAMMA_HEDGING=false

# Cache
CACHE_BACKEND=valkey  # or "file"
VALKEY_ADDR=localhost:6379

# HTTP API
HTTP_PORT=8080
```

### 2. Minimal Code Changes

In your `main.go`, add these lines after exchange creation:

```go
// After creating the exchange (line ~77)
exchange, err := factory.CreateExchange(cfg)

// Add: Initialize arbitrage system
orchestrator, err := InitializeArbitrageSystem(cfg, exchange)
if err != nil {
    log.Fatalf("Failed to initialize arbitrage system: %v", err)
}
defer orchestrator.Stop()
```

### 3. Integrate RFQ Processing

In the RFQ processing section (around line 241), add:

```go
// When processing RFQ
if orchestrator != nil {
    // Track the RFQ in the orchestrator
    if trade, err := orchestrator.SubmitRFQTrade(rfqReq); err == nil {
        log.Printf("RFQ %s tracked as trade %s", rfqReq.RfqId, trade.ID)
    }
}
```

### 4. Update Confirmation Handling

In the confirmation handler (around line 188), modify to use orchestrator:

```go
// Instead of calling HedgeOrder directly
if orchestrator != nil {
    orchestrator.OnRFQConfirmation(confirmation)
} else {
    // Fallback to original behavior
    HedgeOrder(confirmation, underlying, cfg, exchange)
}
```

## API Endpoints

Once integrated, these endpoints are available:

### Submit Manual Trade
```bash
curl -X POST http://localhost:8080/api/trade \
  -H "Content-Type: application/json" \
  -d '{
    "instrument": "ETH-20231225-3000-C",
    "strike": "3000",
    "expiry": 1703462400,
    "is_put": false,
    "quantity": "1.0",
    "price": "0.05"
  }'
```

### View Active Trades
```bash
curl http://localhost:8080/api/trades
```

### Check Risk Metrics
```bash
curl http://localhost:8080/api/risk
```

### View Positions
```bash
curl http://localhost:8080/api/positions
```

### Prometheus Metrics
```bash
curl http://localhost:8080/metrics
```

## Testing the Integration

1. **Start with existing mode** - Run without any changes to verify it still works
2. **Enable HTTP API** - Add the orchestrator initialization and test manual trades
3. **Enable risk checks** - Set position limits and verify they're enforced
4. **Enable gamma hedging** - Set `GAMMA_HEDGING=true` and monitor delta adjustments
5. **Switch to Valkey** - Set `CACHE_BACKEND=valkey` for production caching

## Migration Path

### Phase 1: Monitoring Only
- Initialize orchestrator
- Track trades but don't modify execution
- Use HTTP API to monitor

### Phase 2: Risk Management
- Enable position and delta limits
- Add pre-trade validation

### Phase 3: Full Integration
- Use HedgeManager for all hedging
- Enable gamma hedging
- Switch to Valkey cache

## Rollback Plan

If issues arise, simply:
1. Remove orchestrator initialization
2. Existing HedgeOrder calls continue working
3. No data migration needed

## Configuration Reference

All new options with defaults:

```yaml
# Command line flags
--http_port=8080
--enable_manual_trades=true
--gamma_hedging=false
--cache_backend=file
--valkey_addr=localhost:6379

# Environment variables
MAX_POSITION_SIZE=1000
MAX_DELTA_EXPOSURE=500
GAMMA_MAX_DELTA=100
```

## Troubleshooting

### HTTP API not accessible
- Check `ENABLE_MANUAL_TRADES=true`
- Verify port is not in use
- Check firewall rules

### Trades not being hedged
- Verify orchestrator is initialized
- Check logs for error messages
- Ensure exchange connection is active

### Cache errors
- For Valkey: ensure Redis/Valkey is running
- For file cache: check directory permissions
- Verify `CACHE_BACKEND` setting

## Next Steps

1. **YAML Configuration** - Add support for config files
2. **Deribit Integration** - Add when Derive is stable
3. **Advanced Strategies** - Implement spreads, straddles
4. **Performance Tuning** - Optimize for high-frequency trading