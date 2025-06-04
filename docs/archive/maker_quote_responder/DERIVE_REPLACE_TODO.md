# TODO: Fix Derive Replace Order in Production

## Current Status
The Derive replace order functionality is not working in production. Our investigation revealed several critical issues that need to be fixed.

## Key Differences Between Test and Production

### 1. **CRITICAL: Wrong Domain Separator and Action Typehash**
**Current Production (`derive_order.go`):**
```go
domainSeparator := common.HexToHash("0xd96e5f90797da7ec8dc4e276260c7f3f87fedf68775fbe1ef116e996fc60441b")
actionTypehash := common.HexToHash("0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17")
```

**What Actually Works (from our test):**
```go
domainSeparator := common.HexToHash("0x8f06151ae86a1e59b6cf39212fb0978551b4dcafefd44e3a8b860a9c0b1e6141")
actionTypehash := common.HexToHash("0x5147386a8c7e1c2fb020f0ad9cd5e9c6e28dbce75c86560d0956e15bc50e3041")
```

**Action Required:** Update the constants in `derive_order.go` to match what actually works.

### 2. **Replace is NOT Atomic**
**Current Assumption:** Replace atomically swaps orders
**Reality:** Replace cancels first, then creates. If signature fails, you lose the original order.

**Action Required:** 
- Add error handling for partial failures
- Consider implementing retry logic
- Document this behavior clearly in comments

### 3. **Field Name Mismatch**
**Current Production (`market_maker_derive.go` line 308):** Uses `order_id_to_cancel` ✓ (correct)
**Documentation/Comments:** May reference `order_id` (incorrect)

**Action Required:** Ensure all documentation uses `order_id_to_cancel`

## Implementation Checklist

### Immediate Fixes Needed:

1. **[ ] Update EIP-712 Constants**
   - File: `derive_order.go`
   - Update domain separator to: `0x8f06151ae86a1e59b6cf39212fb0978551b4dcafefd44e3a8b860a9c0b1e6141`
   - Update action typehash to: `0x5147386a8c7e1c2fb020f0ad9cd5e9c6e28dbce75c86560d0956e15bc50e3041`

2. **[ ] Add Non-Atomic Handling**
   - File: `market_maker_derive.go` 
   - In `ReplaceOrder` method, add handling for partial failures
   - Log when cancel succeeds but create fails
   - Implement recovery strategy (retry create or place new order)

3. **[ ] Verify All Required Fields**
   - Ensure replace request includes ALL fields:
     - `order_id_to_cancel` (not `order_id`)
     - `instrument_name`
     - `direction`
     - `order_type`
     - `time_in_force`
     - `amount`
     - `limit_price`
     - `max_fee`
     - `subaccount_id`
     - `nonce`
     - `signature_expiry_sec`
     - `owner`
     - `signer`
     - `signature`
     - `mmp` (optional)

4. **[ ] Add Integration Tests**
   - Add test that actually places and replaces orders
   - Add test for failure scenarios (invalid signature)
   - Verify non-atomic behavior

5. **[ ] Update Error Messages**
   - Make it clear when a replace partially fails
   - Log both cancelled order ID and failed create reason

### Code Changes Needed:

```go
// In derive_order.go, update the Sign method:
func (a *DeriveAction) Sign(privateKey *ecdsa.PrivateKey) error {
    // UPDATE THESE VALUES:
    domainSeparator := common.HexToHash("0x8f06151ae86a1e59b6cf39212fb0978551b4dcafefd44e3a8b860a9c0b1e6141")
    actionTypehash := common.HexToHash("0x5147386a8c7e1c2fb020f0ad9cd5e9c6e28dbce75c86560d0956e15bc50e3041")
    // ... rest of signing logic
}
```

```go
// In market_maker_derive.go, enhance error handling:
func (d *DeriveMarketMakerExchange) ReplaceOrder(...) (string, error) {
    // ... existing code ...
    
    // After sending replace request:
    if result.Result != nil {
        if result.Result.CancelledOrder != nil && result.Result.Order == nil {
            // Cancel succeeded but create failed
            log.Printf("WARNING: Replace cancelled order %s but failed to create new order: %v",
                result.Result.CancelledOrder.OrderID, result.Result.CreateOrderError)
            // Attempt recovery: try to place a new order
            return d.PlaceOrder(instrument, side, price, amount)
        }
    }
}
```

## Testing Instructions

1. **Manual Test:**
   ```bash
   # Set environment
   export DERIVE_PRIVATE_KEY=your_key
   export DERIVE_WALLET=your_wallet
   
   # Run the working test
   go run simple_replace_test_auto.go
   ```

2. **Integration Test:**
   ```bash
   # After fixes, run the integration test
   go test -v -run TestDeriveReplaceOrder
   ```

## Notes

- The `simple_replace_test_auto.go` file contains a working implementation with the correct constants
- The test successfully places and replaces orders when using the right domain separator and typehash
- Always ensure orders are placed with prices that won't execute immediately (well below market for buys)

## Priority: HIGH
This is a critical fix as incorrect replace behavior can lead to:
- Lost orders (cancelled but not replaced)
- Incorrect risk exposure
- Market maker being unable to update quotes efficiently