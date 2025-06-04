#!/bin/bash

# Script to test different aggression levels
# This helps visualize quote placement behavior

echo "Testing Market Maker Aggression Levels"
echo "====================================="

# Default values
EXPIRY="${EXPIRY:-20250630}"
STRIKES="${STRIKES:-3000}"
SIZE="${SIZE:-0.1}"

# Color codes
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to run market maker with specific aggression
test_aggression() {
    local aggression=$1
    local description=$2
    
    echo -e "\n${YELLOW}Testing aggression=$aggression - $description${NC}"
    echo "Command: atomizer market-maker --expiry $EXPIRY --strikes $STRIKES --size $SIZE --aggression $aggression --dry-run"
    
    # Run for 10 seconds then kill
    timeout 10s atomizer market-maker \
        --expiry "$EXPIRY" \
        --strikes "$STRIKES" \
        --size "$SIZE" \
        --aggression "$aggression" \
        --dry-run 2>&1 | grep -E "(Starting|Placing|bid:|ask:|spread:|WARN)" || true
        
    echo "---"
}

# Run tests
echo -e "${GREEN}Conservative Mode Tests:${NC}"
test_aggression 0.0 "Join best bid/ask (most passive)"
test_aggression 0.3 "30% toward mid"
test_aggression 0.5 "Halfway to mid"
test_aggression 0.7 "70% toward mid"
test_aggression 0.9 "Very close to mid (max conservative)"

echo -e "\n${GREEN}Aggressive Mode Tests:${NC}"
test_aggression 1.0 "Cross spread (default aggressive)"
test_aggression 1.5 "Cross spread (future use)"

echo -e "\n${GREEN}Edge Cases:${NC}"
test_aggression -0.5 "Negative (should clamp to 0.0)"
test_aggression 0.95 "Above 0.9 (should clamp to 0.9 in conservative)"

echo -e "\n${YELLOW}Test complete!${NC}"
echo "Note: Use --debug flag for more detailed output"