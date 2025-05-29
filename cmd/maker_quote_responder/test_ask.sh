#!/bin/bash
# Simple test script for placing a high ask on Deribit

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}=== Deribit Ask Order Test ===${NC}"
echo "This will place a SELL order at 2x the best ask price"
echo ""

# Default to a far out option that's likely to have low value
DEFAULT_INSTRUMENT="ETH-28MAR25-5000-C"

# Load .env if it exists
if [ -f .env ]; then
    set -a
    source .env
    set +a
fi

# Run the test with safe defaults
echo -e "${GREEN}Using defaults:${NC}"
echo "  Instrument: $DEFAULT_INSTRUMENT (far OTM call)"
echo "  Quantity: 0.1 ETH"
echo "  Price: 2x best ask"
echo "  Mode: Mainnet (use --test for testnet)"
echo ""
echo "Example commands:"
echo "  ./test_ask.sh                    # Mainnet with defaults"
echo "  ./test_ask.sh --test             # Testnet with defaults" 
echo "  ./test_ask.sh --qty 0.01         # Smaller quantity"
echo "  ./test_ask.sh --mult 3.0         # 3x best ask"
echo "  ./test_ask.sh --instrument ETH-31JAN25-4000-C"
echo ""

# Build and run
go run test_deribit_ask.go "$@"