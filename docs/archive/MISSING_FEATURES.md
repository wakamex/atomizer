# Missing Features After Refactor

## Critical Missing Features

### 1. **Configuration Validation**
- Private key format validation (64 hex chars)
- Private key to address derivation and matching
- Asset mapping configuration
- Environment variable fallback logic

### 2. **Asset Mapping**
Currently hardcoded empty map in main.go:
```go
AssetMapping: make(map[string]string), // TODO: Load from config
```

Should have:
```go
// Testnet mappings
"0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9": "ETH",  // Sepolia WETH
"0xb67bfa7b488df4f2efa874f4e59242e9130ae61f": "ETH",  // Another testnet ETH
```

### 3. **Helper Functions**
Missing decimal conversion utilities:
- `DecimalFromBigInt(value *big.Int, exp int32) decimal.Decimal`
- `BigIntFromDecimal(value decimal.Decimal, exp int32) *big.Int`
- `DecimalFromString(value string) decimal.Decimal` (already exists in arbitrage)

### 4. **Debug Mode & Monitoring**
- Global debug mode toggle
- Periodic stats reporting (30 second intervals)
- Order tracking consistency checks
- Bid-ask spread tracking per instrument

### 5. **Valkey/Redis Cache Backend**
- Redis client for market data caching
- TTL support
- Key prefixing strategy

### 6. **Manual Order Features**
- Standalone manual order submission
- Environment-based configuration
- Order status checking after placement

### 7. **HTTP Client Configuration**
Standard HTTP client with:
- 30 second timeout
- Keep-alive disabled
- Max idle connections

## Important Missing Features

### 8. **Connection Retry Logic**
- Exponential backoff configuration
- Session-based retry management
- Connection health monitoring

### 9. **Build Information**
- Binary hash generation for version tracking
- Build info logging on startup

### 10. **Exchange-Specific Features**
- Deribit asymmetric key support
- ED25519 authentication for Deribit
- Pending: OKX and Bybit implementations

## Nice-to-Have Features

### 11. **Test Utilities**
- Connection test scripts
- Manual order test scripts
- Debug order verification

### 12. **Monitoring Endpoints**
The HTTP API has these endpoints not fully implemented:
- `/metrics` - Prometheus format metrics
- `/api/risk` - Detailed risk metrics
- `/api/positions` - Position tracking

## Implementation Priority

1. **Asset Mapping** - Critical for RFQ to work
2. **Configuration Validation** - Security and correctness
3. **Helper Functions** - Used throughout
4. **Debug/Stats** - Operations visibility
5. **Cache Backend** - Performance optimization
6. **HTTP Endpoints** - Manual trading capability

## Quick Fixes Needed

### Add to internal/config/config.go:
```go
// DefaultAssetMapping provides default asset mappings for known tokens
var DefaultAssetMapping = map[string]string{
    "0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9": "ETH", // Sepolia WETH
    "0xb67bfa7b488df4f2efa874f4e59242e9130ae61f": "ETH", // Testnet ETH
    // Add more as needed
}
```

### Add to internal/utils/decimal.go:
```go
package utils

import (
    "math/big"
    "github.com/shopspring/decimal"
)

func DecimalFromBigInt(value *big.Int, exp int32) decimal.Decimal {
    if value == nil {
        return decimal.Zero
    }
    return decimal.NewFromBigInt(value, exp)
}

func BigIntFromDecimal(value decimal.Decimal, exp int32) *big.Int {
    multiplier := decimal.New(1, -exp)
    result := value.Mul(multiplier)
    return result.BigInt()
}
```

### Add private key validation to cmd/atomizer/main.go:
```go
func validatePrivateKey(privateKey, makerAddress string) error {
    if len(privateKey) != 64 {
        return fmt.Errorf("private key must be 64 hex characters, got %d", len(privateKey))
    }
    
    // Validate hex
    if !regexp.MustCompile(`^[0-9a-fA-F]+$`).MatchString(privateKey) {
        return fmt.Errorf("private key must be hexadecimal")
    }
    
    // Derive address
    derivedAddr, err := privateKeyToAddress(privateKey)
    if err != nil {
        return err
    }
    
    if !strings.EqualFold(derivedAddr, makerAddress) {
        return fmt.Errorf("private key doesn't match maker address")
    }
    
    return nil
}
```