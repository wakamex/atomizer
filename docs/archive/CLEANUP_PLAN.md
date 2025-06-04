# Cleanup Plan for cmd/maker_quote_responder

## Overview
After successful migration to the new modular structure, the following files from `cmd/maker_quote_responder` can be safely removed.

## Files Migrated and Safe to Delete

### Core Components (Already Migrated)
- ✅ `arbitrage_orchestrator.go` → `internal/arbitrage/orchestrator.go`
- ✅ `http_api.go` → `internal/api/server.go`
- ✅ `rfq_processor.go` → `internal/rfq/processor.go`
- ✅ `risk_manager.go` → `internal/risk/manager.go`
- ✅ `hedge_manager.go` → `internal/hedging/manager.go`
- ✅ `gamma_hedger.go` → `internal/hedging/gamma/hedger.go`
- ✅ `gamma_module.go` → `internal/hedging/gamma/module.go`

### Market Maker Components (Already Migrated)
- ✅ `market_maker.go` → `internal/marketmaker/market_maker.go`
- ✅ `orders.go` → `internal/marketmaker/orders.go`
- ✅ `quotes.go` → `internal/marketmaker/quotes.go`
- ✅ `positions.go` → `internal/marketmaker/positions.go`
- ✅ `reconciliation.go` → `internal/marketmaker/reconciliation.go`
- ✅ `stats.go` → `internal/marketmaker/stats.go`
- ✅ `types.go` → `internal/types/`

### Exchange Components (Already Migrated)
- ✅ `derive_*.go` files → `internal/exchange/derive/`
- ✅ `deribit_*.go` files → `internal/exchange/deribit/`
- ✅ `exchange_*.go` files → `internal/exchange/`

## Files That Need Review Before Deletion

### Main Entry Point
- `main.go` - Contains WebSocket connection logic that needs to be ported for RFQ responder

### Utilities
- `app_config.go` - Configuration loading (partially migrated to `internal/config/`)
- `helpers.go` - Utility functions that might be needed
- `common.go` - Common types and utilities

### Test Files
- `arbitrage_test.go` - Should be adapted to new structure
- `market_maker_test.go` - Should be adapted to new structure
- `market_maker_concurrent_test.go` - Should be adapted to new structure
- `quoter_test.go` - Should be adapted to new structure
- `derive_replace_order_test.go` - Already migrated
- `final_integration_test.go` - Should be adapted

### Documentation
- Various `*.md` files - Should be reviewed and updated
- `*.sh` scripts - Test scripts that might still be useful

## Recommended Deletion Process

1. **Phase 1: Backup**
   ```bash
   tar -czf cmd_maker_quote_responder_backup.tar.gz cmd/maker_quote_responder/
   ```

2. **Phase 2: Delete Core Files**
   Remove files that have been fully migrated and tested

3. **Phase 3: Port Remaining Logic**
   - WebSocket connection from main.go
   - Any unique utilities from helpers.go/common.go

4. **Phase 4: Adapt Tests**
   Update test files to work with new structure

5. **Phase 5: Final Cleanup**
   Remove the entire directory once everything is confirmed working

## Notes
- The `cmd/test_derive/` subdirectory can be removed entirely
- Shell scripts should be reviewed - some might be useful for testing
- Documentation files should be consolidated into project-level docs