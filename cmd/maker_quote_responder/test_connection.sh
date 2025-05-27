#!/bin/bash
# Test Deribit connection and basic API operations

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}=== Deribit Connection Test ===${NC}"
echo "This will test your Deribit API connection"
echo ""

# Load .env if it exists
if [ -f .env ]; then
    echo -e "${GREEN}Loading credentials from .env${NC}"
    set -a
    source .env
    set +a
else
    echo -e "${RED}Warning: .env file not found${NC}"
    echo "Using environment variables DERIBIT_API_KEY and DERIBIT_API_SECRET"
fi

# Check for API credentials
if [ -z "$DERIBIT_API_KEY" ] || [ -z "$DERIBIT_API_SECRET" ]; then
    echo -e "${RED}Error: API credentials not found${NC}"
    echo "Please set DERIBIT_API_KEY and DERIBIT_API_SECRET"
    exit 1
fi

echo ""
echo "Available options:"
echo "  ./test_connection.sh         # Test mainnet connection"
echo "  ./test_connection.sh --test  # Test testnet connection"
echo ""

# Run the connection test
go run test_deribit_connection.go "$@"