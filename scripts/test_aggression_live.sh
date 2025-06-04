#!/bin/bash

# Live test script to see actual quote placement with different aggression levels
# This connects to a real exchange in test mode

echo "Live Aggression Test - Watch Quote Placement"
echo "==========================================="
echo "This will run the market maker with different aggression levels"
echo "Watch the log output to see where quotes are placed"
echo ""

# Check for required env vars
if [ -z "$DERIVE_PRIVATE_KEY" ] || [ -z "$DERIVE_WALLET_ADDRESS" ]; then
    echo "Error: Set DERIVE_PRIVATE_KEY and DERIVE_WALLET_ADDRESS environment variables"
    exit 1
fi

# Test parameters
EXPIRY="${EXPIRY:-20250630}"
STRIKE="${STRIKE:-3000}"
SIZE="${SIZE:-0.01}"  # Small size for testing
EXCHANGE="${EXCHANGE:-derive}"

echo "Testing on $EXCHANGE exchange"
echo "Instrument: ETH-$EXPIRY-$STRIKE-C"
echo ""

# Function to run test
run_test() {
    local aggression=$1
    local duration=${2:-15}
    
    echo "----------------------------------------"
    echo "Testing aggression=$aggression for ${duration}s"
    echo "----------------------------------------"
    
    timeout ${duration}s atomizer market-maker \
        --exchange "$EXCHANGE" \
        --test \
        --expiry "$EXPIRY" \
        --strikes "$STRIKE" \
        --size "$SIZE" \
        --aggression "$aggression" \
        --refresh 5 2>&1 | grep -E "(Starting|Placing|Updated|bid:|ask:|mid:|spread:)" || true
        
    echo ""
    sleep 2
}

# Run tests with increasing aggression
echo "Starting tests..."
run_test 0.0 10   # Join best
run_test 0.5 10   # Halfway
run_test 0.9 10   # Near mid
run_test 1.0 10   # Cross spread

echo "Test complete!"
echo ""
echo "Summary:"
echo "- Aggression 0.0: Should place orders at best bid/ask"
echo "- Aggression 0.5: Should place orders halfway between best and mid"
echo "- Aggression 0.9: Should place orders very close to mid"
echo "- Aggression 1.0: Should cross the spread (if improvement > 0)"