#!/bin/bash

# Test script for manual trade submission
# This script demonstrates how to submit a manual trade to the arbitrage bot

# Configuration
API_HOST="${API_HOST:-localhost}"
API_PORT="${API_PORT:-8080}"
BASE_URL="http://${API_HOST}:${API_PORT}"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Manual Trade Test Script ===${NC}"
echo "API Endpoint: ${BASE_URL}"
echo ""

# Check if server is running
echo -n "Checking if API server is running... "
if curl -s "${BASE_URL}/health" > /dev/null; then
    echo -e "${GREEN}OK${NC}"
else
    echo -e "${RED}FAILED${NC}"
    echo "Please ensure the maker_quote_responder is running with HTTP API enabled"
    echo "Run with: ENABLE_MANUAL_TRADES=true ./run.sh"
    exit 1
fi

# Function to submit a trade
submit_trade() {
    local instrument=$1
    local strike=$2
    local expiry=$3
    local is_put=$4
    local quantity=$5
    local price=$6
    
    echo -e "\n${BLUE}Submitting trade:${NC}"
    echo "  Instrument: $instrument"
    echo "  Strike: $strike"
    echo "  Expiry: $(date -d @$expiry '+%Y-%m-%d %H:%M:%S')"
    echo "  Type: $([ "$is_put" = "true" ] && echo "PUT" || echo "CALL")"
    echo "  Quantity: $quantity"
    echo "  Price: $price"
    
    response=$(curl -s -X POST "${BASE_URL}/api/trade" \
        -H "Content-Type: application/json" \
        -d "{
            \"instrument\": \"$instrument\",
            \"strike\": \"$strike\",
            \"expiry\": $expiry,
            \"is_put\": $is_put,
            \"quantity\": \"$quantity\",
            \"price\": \"$price\"
        }")
    
    echo -e "\n${GREEN}Response:${NC}"
    echo "$response" | jq '.' 2>/dev/null || echo "$response"
}

# Test trades
echo -e "\n${BLUE}=== Test 1: ETH Call Option ===${NC}"
# Current time + 30 days
EXPIRY_1=$(($(date +%s) + 30*24*60*60))
submit_trade "ETH-CALL" "3500" "$EXPIRY_1" "false" "1.0" "0.05"

echo -e "\n${BLUE}=== Test 2: ETH Put Option ===${NC}"
# Current time + 60 days
EXPIRY_2=$(($(date +%s) + 60*24*60*60))
submit_trade "ETH-PUT" "2800" "$EXPIRY_2" "true" "0.5" "0.03"

echo -e "\n${BLUE}=== Test 3: Small ETH Call (Testing minimum size) ===${NC}"
# Current time + 7 days
EXPIRY_3=$(($(date +%s) + 7*24*60*60))
submit_trade "ETH-CALL" "3200" "$EXPIRY_3" "false" "0.1" "0.01"

# Check active trades
echo -e "\n${BLUE}=== Active Trades ===${NC}"
curl -s "${BASE_URL}/api/trades" | jq '.' 2>/dev/null || curl -s "${BASE_URL}/api/trades"

# Check risk metrics
echo -e "\n${BLUE}=== Risk Metrics ===${NC}"
curl -s "${BASE_URL}/api/risk" | jq '.' 2>/dev/null || curl -s "${BASE_URL}/api/risk"

# Check positions
echo -e "\n${BLUE}=== Current Positions ===${NC}"
curl -s "${BASE_URL}/api/positions" | jq '.' 2>/dev/null || curl -s "${BASE_URL}/api/positions"

# Test invalid trade (should be rejected)
echo -e "\n${BLUE}=== Test 4: Invalid Trade (Negative quantity) ===${NC}"
submit_trade "ETH-CALL" "3000" "$EXPIRY_1" "false" "-1.0" "0.05"

# Test risk limit (if MAX_POSITION_SIZE is set)
echo -e "\n${BLUE}=== Test 5: Large Trade (Testing position limits) ===${NC}"
submit_trade "ETH-CALL" "3000" "$EXPIRY_1" "false" "10000.0" "0.05"

echo -e "\n${BLUE}=== Test Complete ===${NC}"