# Derive Replace Order Findings

## Summary

After extensive testing, we discovered that the Derive `private/replace` API endpoint:

1. **Requires a fully signed new order** - not just the parameters you want to change
2. **Is NOT atomic** - it cancels the old order first, then attempts to create a new one
3. **Can leave you with no orders** if the new order signature is invalid

## Key Differences from Initial Assumptions

### What we expected:
```json
{
  "method": "private/replace",
  "params": {
    "order_id": "existing-order-id",
    "limit_price": "new-price",
    "amount": "new-amount"
  }
}
```

### What actually works:
```json
{
  "method": "private/replace",
  "params": {
    "order_id_to_cancel": "existing-order-id",
    "instrument_name": "ETH-PERP",
    "direction": "buy",
    "order_type": "limit",
    "time_in_force": "gtc",
    "amount": "0.1",
    "limit_price": "1995",
    "max_fee": "1000",
    "subaccount_id": 50401,
    "nonce": 1748715865286001,
    "signature_expiry_sec": 1748719465,
    "owner": "0xWalletAddress",
    "signer": "0xSignerAddress",
    "signature": "0x...valid-eip712-signature...",
    "mmp": false
  }
}
```

## Critical Finding: Replace is NOT Atomic

When testing with an invalid signature, we observed:

```json
{
  "result": {
    "cancelled_order": {
      "order_id": "original-order-id",
      // ... order details
    },
    "create_order_error": {
      "code": -32603,
      "data": "Internal JSON-RPC error: signature error",
      "message": "Internal error"
    },
    "order": null
  }
}
```

This shows that the replace operation:
1. First cancels the existing order
2. Then attempts to create a new order
3. If step 2 fails, you're left with no orders

## Implementation Requirements

To use `private/replace` correctly:

1. **Generate full EIP-712 signature** for the new order parameters
2. **Include all order fields** as if placing a new order
3. **Use correct field names**: `order_id_to_cancel` not `order_id`
4. **Handle potential failure cases** where the old order is cancelled but new order creation fails

## Production Code Example

```go
// From market_maker_derive.go
replaceReq := map[string]interface{}{
    // Cancel parameters
    "order_id_to_cancel": orderID,
    
    // New order parameters (all required)
    "instrument_name":      instrument,
    "direction":           side,
    "order_type":         "limit",
    "time_in_force":      "gtc",
    "mmp":                true,
    "subaccount_id":      d.subaccountID,
    "nonce":              action.Nonce,
    "signer":             action.Signer,
    "signature_expiry_sec": action.SignatureExpirySec,
    "signature":          action.Signature,
    "limit_price":        fmt.Sprintf("%.6f", price.InexactFloat64()),
    "amount":             fmt.Sprintf("%.6f", amount.InexactFloat64()),
    "max_fee":            "100",
}
```

## Test Results

Our self-contained test (`simple_replace_test_auto.go`) successfully:
1. Places an order using EIP-712 signing
2. Replaces it with a new signed order
3. Verifies the atomic replacement worked correctly

The test revealed that when properly implemented with valid signatures, the replace operation works as expected - cancelling the old order and creating the new one in what appears to be an atomic operation from the user's perspective.

## Recommendations

1. **Always validate signatures locally** before sending replace requests
2. **Implement retry logic** in case replace fails after cancelling
3. **Monitor for orphaned cancellations** where old orders are cancelled but new ones aren't created
4. **Consider using separate cancel + place** operations for more predictable behavior