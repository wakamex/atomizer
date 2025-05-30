#!/bin/bash

# Default values
VM_URL="${VM_URL:-http://localhost:8428}"
METRIC="${METRIC:-market_bid_price}"
INSTRUMENT="${INSTRUMENT:-ETH-}"
START="${START:-1h}"
STEP="${STEP:-10s}"

# Build the tool
echo "Building correlation analyzer..."
go build -o analyze_correlation main.go

# Run analysis
echo "Running correlation analysis..."
./analyze_correlation \
  -vm-url="$VM_URL" \
  -metric="$METRIC" \
  -instrument="$INSTRUMENT" \
  -start="$START" \
  -step="$STEP"

# Example queries for different analyses:
echo -e "\n\n=== Additional Analysis Examples ==="
echo "1. Analyze ask prices:"
echo "   METRIC=market_ask_price ./analyze.sh"
echo ""
echo "2. Analyze spreads:"
echo "   METRIC=market_spread_percent ./analyze.sh"
echo ""
echo "3. Analyze specific instrument:"
echo "   INSTRUMENT=ETH-20250601-2600-C ./analyze.sh"
echo ""
echo "4. Analyze longer time period:"
echo "   START=24h ./analyze.sh"
echo ""
echo "5. Compare bid-ask spread correlation:"
echo "   METRIC=market_spread ./analyze.sh"