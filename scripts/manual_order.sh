#!/bin/bash

# Script to place manual orders using the atomizer CLI
# Usage: ./manual_order.sh [options]
#
# Environment variables:
#   DERIVE_PRIVATE_KEY    - Private key for Derive exchange
#   DERIVE_WALLET_ADDRESS - Wallet address for Derive exchange
#   DERIBIT_API_KEY       - API key for Deribit exchange
#   DERIBIT_API_SECRET    - API secret for Deribit exchange
#
# Can also override with environment:
#   ORDER_INSTRUMENT - Instrument to trade (e.g., ETH-PERP)
#   ORDER_SIDE       - Side (buy/sell)
#   ORDER_PRICE      - Price
#   ORDER_AMOUNT     - Amount

set -e

# Default values
EXCHANGE="${EXCHANGE:-derive}"
INSTRUMENT="${ORDER_INSTRUMENT:-ETH-PERP}"
SIDE="${ORDER_SIDE:-buy}"
PRICE="${ORDER_PRICE:-0.1}"
AMOUNT="${ORDER_AMOUNT:-0.1}"

# Build the binary if needed
if [ ! -f "./cmd/atomizer/atomizer" ]; then
    echo "Building atomizer..."
    go build -o ./cmd/atomizer/atomizer ./cmd/atomizer
fi

# Run the manual order command
./cmd/atomizer/atomizer manual-order \
    --exchange "$EXCHANGE" \
    --instrument "$INSTRUMENT" \
    --side "$SIDE" \
    --price "$PRICE" \
    --amount "$AMOUNT" \
    "$@"