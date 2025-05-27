#!/bin/bash

# Test Deribit Ed25519 asymmetric key authentication

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}=== Deribit Ed25519 Authentication Test ===${NC}"
echo "This tests asymmetric key authentication with Ed25519"
echo ""

# Load .env if exists
if [ -f .env ]; then
    echo -e "${GREEN}Loading from .env${NC}"
    set -a
    source .env
    set +a
fi

# Check for required variables
if [ -z "$DERIBIT_CLIENT_ID" ]; then
    echo -e "${RED}Error: DERIBIT_CLIENT_ID not set${NC}"
    echo "Please set in .env or environment:"
    echo "  export DERIBIT_CLIENT_ID=your_client_id"
    exit 1
fi

if [ -z "$DERIBIT_PRIVATE_KEY" ] && [ -z "$DERIBIT_PRIVATE_KEY_FILE" ]; then
    echo -e "${RED}Error: Private key not found${NC}"
    echo "Please set one of:"
    echo "  DERIBIT_PRIVATE_KEY (key content)"
    echo "  DERIBIT_PRIVATE_KEY_FILE (path to key file)"
    exit 1
fi

echo "Configuration:"
echo "  Client ID: $DERIBIT_CLIENT_ID"
if [ ! -z "$DERIBIT_PRIVATE_KEY_FILE" ]; then
    echo "  Key file: $DERIBIT_PRIVATE_KEY_FILE"
else
    echo "  Key: From DERIBIT_PRIVATE_KEY env var"
fi
echo ""
echo "Usage:"
echo "  ./test_ed25519.sh         # Test mainnet"
echo "  ./test_ed25519.sh --test  # Test testnet"
echo ""

# Run the test
go run test_ed25519_connection.go deribit_ed25519.go "$@"