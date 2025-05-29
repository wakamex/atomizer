# Hedging Strategy

## Overview
The system has two starting points for acquiring call options, but the hedging strategy is always the same.

## Starting Points

### A) RFQ Trade (Automated)
- User sells call to us via Rysk RFQ
- We buy the call from the user
- This triggers automatic hedging

### B) Manual Trade
- We manually buy a call (via API endpoint)
- This simulates having acquired a call position
- This triggers the same hedging flow

## Hedging Strategy
**Regardless of how we acquire the call (A or B), the hedge is always the same:**
- We SELL the call on Derive at 2x ask price
- This is a defensive strategy to provide liquidity at a safe distance from the market

## Implementation Details
- When `IsTakerBuy = true`: User/we bought → We need to hedge by SELLING
- When `IsTakerBuy = false`: User sold to us → We need to hedge by SELLING
- **Result: We always SELL on Derive for hedging**

## Summary
Buy Call (from user or manually) → Hedge by Selling Call on Derive at 2x ask