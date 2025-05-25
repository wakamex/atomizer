#!/bin/bash

# Default WebSocket URL - taker endpoint for sending RFQs
WS_URL="${WS_URL:-wss://rip-testnet.rysk.finance/taker}"
MARKETS_FILE="${MARKETS_FILE:-markets.json}"

echo "🚀 Market-wide RFQ Performance Test"
echo "URL: $WS_URL"
echo "Markets file: $MARKETS_FILE"
echo ""

# Check if markets.json exists, if not fetch it
if [ ! -f "$MARKETS_FILE" ]; then
    echo "⚠️  Markets file not found. Fetching from API..."
    cd /code/atomizer/cmd/markets && ./markets -output "../send_quote/$MARKETS_FILE"
    cd /code/atomizer/cmd/send_quote
fi

cd /code/atomizer/cmd/send_quote && ./send_quote \
  -url "$WS_URL" \
  -chainId 84532 \
  -quantity 1000000000000000000 \
  -taker 0x0000000000000000000000000000000000000000 \
  -markets "$MARKETS_FILE"