#!/bin/bash

# Test script for running market maker with asymmetric keys

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}=== Market Maker with Asymmetric Keys ===${NC}"
echo ""

# Load .env if exists
if [ -f .env ]; then
    echo -e "${GREEN}Loading from .env${NC}"
    set -a
    source .env
    set +a
fi

# Check for asymmetric key credentials
if [ -z "$DERIBIT_CLIENT_ID" ]; then
    echo -e "${RED}Error: DERIBIT_CLIENT_ID not set${NC}"
    echo "Add to .env:"
    echo "  DERIBIT_CLIENT_ID=your_client_id"
    exit 1
fi

if [ -z "$DERIBIT_PRIVATE_KEY" ] && [ -z "$DERIBIT_PRIVATE_KEY_FILE" ]; then
    echo -e "${RED}Error: Private key not found${NC}"
    echo "Add to .env one of:"
    echo "  DERIBIT_PRIVATE_KEY=\"-----BEGIN PRIVATE KEY-----"
    echo "  ...your key..."
    echo "  -----END PRIVATE KEY-----\""
    echo "  OR"
    echo "  DERIBIT_PRIVATE_KEY_FILE=/path/to/key.pem"
    exit 1
fi

# Check other required vars
if [ -z "$MAKER_ADDRESS" ]; then
    echo -e "${RED}Error: MAKER_ADDRESS not set${NC}"
    exit 1
fi

# Set defaults
WEBSOCKET_URL=${WEBSOCKET_URL:-"wss://rip-testnet.rysk.finance/maker"}
RFQ_ASSET_ADDRESSES=${RFQ_ASSET_ADDRESSES:-"0xb67bfa7b488df4f2efa874f4e59242e9130ae61f"}

echo "Configuration:"
echo "  Exchange: Deribit (Asymmetric Auth)"
echo "  Client ID: $DERIBIT_CLIENT_ID"
echo "  Maker Address: $MAKER_ADDRESS"
echo "  WebSocket: $WEBSOCKET_URL"
echo "  Assets: $RFQ_ASSET_ADDRESSES"
echo ""
echo -e "${YELLOW}Starting market maker...${NC}"
echo ""

# Run the market maker
./maker_quote_responder \
    --websocket_url="$WEBSOCKET_URL" \
    --rfq_asset_addresses="$RFQ_ASSET_ADDRESSES" \
    --exchange="deribit" \
    --exchange_test_mode=false \
    --quote_valid_duration_seconds=30