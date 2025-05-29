#!/bin/bash

# Test script for running market maker with Derive exchange

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${YELLOW}=== Market Maker with Derive Exchange ===${NC}"
echo ""

# Load .env if exists
if [ -f .env ]; then
    echo -e "${GREEN}Loading from .env${NC}"
    set -a
    source .env
    set +a
fi

# Check for Derive credentials
if [ -z "$DERIVE_API_KEY" ] || [ -z "$DERIVE_API_SECRET" ]; then
    echo -e "${RED}Error: Derive credentials not set${NC}"
    echo "Add to .env:"
    echo "  DERIVE_API_KEY=your_api_key"
    echo "  DERIVE_API_SECRET=your_api_secret"
    exit 1
fi

# Check other required vars
if [ -z "$MAKER_ADDRESS" ]; then
    echo -e "${RED}Error: MAKER_ADDRESS not set${NC}"
    exit 1
fi

if [ -z "$PRIVATE_KEY" ]; then
    echo -e "${RED}Error: PRIVATE_KEY not set${NC}"
    exit 1
fi

# Set defaults
WEBSOCKET_URL=${WEBSOCKET_URL:-"wss://rip-testnet.rysk.finance/maker"}
RFQ_ASSET_ADDRESSES=${RFQ_ASSET_ADDRESSES:-"0xb67bfa7b488df4f2efa874f4e59242e9130ae61f"}

echo "Configuration:"
echo "  Exchange: Derive"
echo "  Maker Address: $MAKER_ADDRESS"
echo "  WebSocket: $WEBSOCKET_URL"
echo "  Assets: $RFQ_ASSET_ADDRESSES"
echo ""
echo -e "${YELLOW}Note: Derive integration is not yet fully implemented${NC}"
echo -e "${YELLOW}This will show the structure but API calls will fail${NC}"
echo ""
echo -e "${YELLOW}Starting market maker...${NC}"
echo ""

# Run the market maker with Derive
./maker_quote_responder \
    --websocket_url="$WEBSOCKET_URL" \
    --rfq_asset_addresses="$RFQ_ASSET_ADDRESSES" \
    --exchange="derive" \
    --exchange_test_mode=false \
    --quote_valid_duration_seconds=30