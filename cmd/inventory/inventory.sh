#!/bin/bash

# Default WebSocket URL
WS_URL="${WS_URL:-wss://rip-testnet.rysk.finance/taker}"

echo "📊 Fetching Market Inventory"
echo "URL: $WS_URL"
echo ""

cd /code/atomizer/cmd/inventory && go run . -url "$WS_URL"