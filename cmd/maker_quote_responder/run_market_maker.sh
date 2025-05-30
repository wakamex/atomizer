#!/bin/bash

# Market Maker Run Script
# This script runs the market maker with common configurations

set -e

# Default values
EXCHANGE="${EXCHANGE:-derive}"
UNDERLYING="${UNDERLYING:-ETH}"
EXPIRY="${EXPIRY:-20250606}"
STRIKES="${STRIKES:-}"
ALL_STRIKES="${ALL_STRIKES:-false}"
SPREAD_BPS="${SPREAD_BPS:-10}"
SIZE="${SIZE:-0.1}"
REFRESH_SEC="${REFRESH_SEC:-1}"
MAX_POSITION="${MAX_POSITION:-1.0}"
MAX_EXPOSURE="${MAX_EXPOSURE:-10.0}"
MIN_SPREAD_BPS="${MIN_SPREAD_BPS:-1000}"
IMPROVEMENT="${IMPROVEMENT:-0.1}"
IMPROVEMENT_REFERENCE_SIZE="${IMPROVEMENT_REFERENCE_SIZE:-0}"
TEST_MODE="${TEST_MODE:-false}"
DRY_RUN="${DRY_RUN:-false}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Market Maker for ${EXCHANGE} ===${NC}"
echo ""

# Check required environment variables
if [[ "$EXCHANGE" == "derive" ]]; then
    if [[ -z "$DERIVE_PRIVATE_KEY" ]]; then
        echo -e "${RED}Error: DERIVE_PRIVATE_KEY not set${NC}"
        exit 1
    fi
    if [[ -z "$DERIVE_WALLET_ADDRESS" ]]; then
        echo -e "${RED}Error: DERIVE_WALLET_ADDRESS not set${NC}"
        exit 1
    fi
elif [[ "$EXCHANGE" == "deribit" ]]; then
    if [[ -z "$DERIBIT_PRIVATE_KEY" ]] && [[ -z "$DERIBIT_API_KEY" ]]; then
        echo -e "${RED}Error: DERIBIT_PRIVATE_KEY or DERIBIT_API_KEY not set${NC}"
        exit 1
    fi
fi

# Check if strikes are specified
if [[ "$ALL_STRIKES" != "true" ]] && [[ -z "$STRIKES" ]]; then
    echo -e "${YELLOW}No strikes specified. Using default strikes...${NC}"
    STRIKES="2600,2700,2800,2900,3000,3100,3200"
fi

# Build command
CMD="./maker_quote_responder market-maker"
CMD="$CMD -exchange=$EXCHANGE"
CMD="$CMD -underlying=$UNDERLYING"
CMD="$CMD -expiry=$EXPIRY"
CMD="$CMD -spread=$SPREAD_BPS"
CMD="$CMD -size=$SIZE"
CMD="$CMD -refresh=$REFRESH_SEC"
CMD="$CMD -max-position=$MAX_POSITION"
CMD="$CMD -max-exposure=$MAX_EXPOSURE"
CMD="$CMD -min-spread=$MIN_SPREAD_BPS"
CMD="$CMD -improvement=$IMPROVEMENT"
CMD="$CMD -improvement-reference-size=$IMPROVEMENT_REFERENCE_SIZE"

if [[ "$TEST_MODE" == "true" ]]; then
    CMD="$CMD -test"
fi

if [[ "$DRY_RUN" == "true" ]]; then
    CMD="$CMD -dry-run"
fi

if [[ "$ALL_STRIKES" == "true" ]]; then
    CMD="$CMD -all-strikes"
else
    CMD="$CMD -strikes=$STRIKES"
fi

# Show configuration
echo "Configuration:"
echo "  Exchange: $EXCHANGE (test mode: $TEST_MODE)"
echo "  Underlying: $UNDERLYING"
echo "  Expiry: $EXPIRY"
if [[ "$ALL_STRIKES" == "true" ]]; then
    echo "  Strikes: ALL"
else
    echo "  Strikes: $STRIKES"
fi
echo "  Spread: ${SPREAD_BPS}bps"
echo "  Size: $SIZE"
echo "  Refresh: ${REFRESH_SEC}s"
echo "  Max Position: $MAX_POSITION"
echo "  Max Exposure: $MAX_EXPOSURE"
echo "  Improvement: $IMPROVEMENT"
echo "  Improvement Reference Size: $IMPROVEMENT_REFERENCE_SIZE"
echo ""

# Run the market maker
echo -e "${GREEN}Starting market maker...${NC}"
echo "Command: $CMD"
echo ""

# Execute
exec $CMD