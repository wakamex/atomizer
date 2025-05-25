#!/bin/bash

# Default API URL
API_URL="${API_URL:-https://rip-testnet.rysk.finance/api/assets}"
OUTPUT_FILE="${OUTPUT_FILE:-markets.json}"

echo "🌐 Fetching Markets Data"
echo "API: $API_URL"
echo "Output: $OUTPUT_FILE"
echo ""

cd /code/atomizer/cmd/markets && go run . -url "$API_URL" -output "$OUTPUT_FILE"