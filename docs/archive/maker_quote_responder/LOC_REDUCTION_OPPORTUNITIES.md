# LOC Reduction Opportunities

## 1. Consolidate Helper Functions (~50-100 lines)
- `getString`, `getFloat64`, `getDecimal`, `getInt64` in market_maker_derive.go (28 lines)
- These could be made generic or use reflection
- Similar patterns exist in multiple files

## 2. Remove Duplicate HTTP Client Creation (~20 lines)
- 4 separate HTTP client creations with same timeout
- Could use a single factory function: `newHTTPClient() *http.Client`

## 3. Consolidate Environment Variable Reading (~30 lines)
- `DERIVE_WALLET_ADDRESS` read in 3 places
- `privateKey` read in 3 places  
- Create a config loader function

## 4. Remove Unused Test Files (~200-300 lines)
- Check if all test files are actually run and needed
- `hedge_test.go` (50 lines) - tests removed hedge.go
- Old integration tests may be outdated

## 5. Consolidate Error Handling Patterns (~50 lines)
- 170 `fmt.Errorf` calls with similar patterns
- Many could use error wrapping functions

## 6. Remove getBuildHash Duplication (~15 lines)
- Same function in main.go used by gamma_hedger_pure_main.go
- Could be moved to a common location

## 7. Consolidate WebSocket Message Handling (~100 lines)
- Similar JSON unmarshaling and error handling patterns
- Could use generic message handler

## 8. Remove Empty or Near-Empty Files
- Check files under 50 lines for actual necessity
- Many may just be thin wrappers

## 9. Consolidate Decimal Conversions (~50 lines)
- 109 `decimal.NewFromFloat` calls
- Many with repeated values (0, 1, 100, etc.)
- Could use constants or helper functions

## 10. Test File Cleanup (~200+ lines)
- `arbitrage_test.go` - has compilation errors
- `integration_test.go` - may be outdated
- `final_integration_test.go` - may duplicate other tests

## Total Potential Reduction: 500-800 lines