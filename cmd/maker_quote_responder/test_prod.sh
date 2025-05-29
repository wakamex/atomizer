#!/bin/bash
# Production testing script for multi-exchange market maker
# Tests with $1 orders at 2x best ask price for safety

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Production Testing Script ===${NC}"
echo -e "${YELLOW}This will test with REAL MONEY on Deribit mainnet${NC}"
echo -e "${YELLOW}Orders will be placed at 2x best ask price for safety${NC}"
echo ""

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found!${NC}"
    echo "Please create .env with:"
    echo "  PRIVATE_KEY=your_private_key_here"
    echo "  DERIBIT_API_KEY=your_deribit_api_key"
    echo "  DERIBIT_API_SECRET=your_deribit_api_secret"
    exit 1
fi

# Load environment variables
set -a
source .env
set +a

# Verify required env vars
if [ -z "$PRIVATE_KEY" ] || [ -z "$DERIBIT_API_KEY" ] || [ -z "$DERIBIT_API_SECRET" ]; then
    echo -e "${RED}Error: Missing required environment variables in .env${NC}"
    exit 1
fi

# Set maker address (you should set this to your actual address)
export MAKER_ADDRESS=${MAKER_ADDRESS:-"0xYourMakerAddress"}

# Production Rysk mainnet URL (update this to actual mainnet URL)
WEBSOCKET_URL="wss://rip-testnet.rysk.finance/maker"  # TODO: Update to mainnet URL

# Asset addresses to quote on (update these for mainnet)
RFQ_ASSET_ADDRESSES="0xb67bfa7b488df4f2efa874f4e59242e9130ae61f"  # TODO: Update to mainnet addresses

# Build the application
echo -e "${GREEN}Building application...${NC}"
go build -o maker_quote_responder .

echo -e "${GREEN}Starting market maker in PRODUCTION mode...${NC}"
echo -e "${YELLOW}Configuration:${NC}"
echo "  Exchange: Deribit (mainnet)"
echo "  Hedge strategy: Sell at 2x best ask"
echo "  WebSocket: $WEBSOCKET_URL"
echo "  Assets: $RFQ_ASSET_ADDRESSES"
echo ""
echo -e "${YELLOW}Press Ctrl+C to stop${NC}"
echo ""

# Run the market maker
./maker_quote_responder \
    --websocket_url="$WEBSOCKET_URL" \
    --rfq_asset_addresses="$RFQ_ASSET_ADDRESSES" \
    --exchange="deribit" \
    --exchange_test_mode=false \
    --quote_valid_duration_seconds=30