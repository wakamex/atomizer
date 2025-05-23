#!/bin/bash

# Navigate to the script's directory to ensure relative paths work correctly
cd "$(dirname "$0")"

# --- Configuration ---
# These values will be used to run the application.
# You can modify them here or override them with environment variables if the script is adapted to do so.

# Required: Your Ethereum maker address
MAKER_ADDRESS_DEFAULT="0x9eAFc0c2b04D96a1C1edAdda8A474a4506752207"
# Required: Rysk Finance WebSocket URL
WEBSOCKET_URL_DEFAULT="wss://rip-testnet.rysk.finance/maker"
# Required: Comma-separated list of RFQ asset addresses
RFQ_ASSET_ADDRESSES_DEFAULT="0xb67bfa7b488df4f2efa874f4e59242e9130ae61f"
# Required: Dummy price for quotes
DUMMY_PRICE_DEFAULT="12500000000000000000"
# Optional: Quote valid duration in seconds
QUOTE_VALID_DURATION_SECONDS_DEFAULT="45"

# Use environment variables if set, otherwise use defaults
MAKER_ADDRESS="${MAKER_ADDRESS:-$MAKER_ADDRESS_DEFAULT}"
WEBSOCKET_URL="${WEBSOCKET_URL:-$WEBSOCKET_URL_DEFAULT}"
RFQ_ASSET_ADDRESSES="${RFQ_ASSET_ADDRESSES:-$RFQ_ASSET_ADDRESSES_DEFAULT}"
DUMMY_PRICE="${DUMMY_PRICE:-$DUMMY_PRICE_DEFAULT}"
QUOTE_VALID_DURATION_SECONDS="${QUOTE_VALID_DURATION_SECONDS:-$QUOTE_VALID_DURATION_SECONDS_DEFAULT}"

# --- Sanity Checks ---

# Check if .env file exists
if [ ! -f .env ]; then
    echo "Error: .env file not found in $(pwd)."
    echo "Please create it with your PRIVATE_KEY as described in README.md."
    exit 1
fi

# Check if the executable exists
EXECUTABLE_NAME="./maker_quote_responder"
if [ ! -f "$EXECUTABLE_NAME" ]; then
    echo "Error: Executable '$EXECUTABLE_NAME' not found in $(pwd)."
    echo "Please build the application first by running:"
    echo "go build -o maker_quote_responder ."
    exit 1
fi

# --- Load .env and Run ---

# Export variables from .env file, filtering out comments and empty lines
# This makes PRIVATE_KEY available to the Go application
export $(grep -v '^#' .env | grep -v '^$' | xargs)

# Check if PRIVATE_KEY is set after attempting to load .env
if [ -z "$PRIVATE_KEY" ]; then
    echo "Error: PRIVATE_KEY is not set in your .env file or as an environment variable."
    exit 1
fi

echo "Starting Maker Quote Responder..."

# Execute the application
# The MAKER_ADDRESS is passed as a command-line flag here for clarity,
# though the app can also pick it up from an environment variable.

MAKER_ADDRESS="$MAKER_ADDRESS" \
./"$EXECUTABLE_NAME" \
    --websocket_url="$WEBSOCKET_URL" \
    --rfq_asset_addresses="$RFQ_ASSET_ADDRESSES" \
    --dummy_price="$DUMMY_PRICE" \
    --quote_valid_duration_seconds="$QUOTE_VALID_DURATION_SECONDS"

echo "Maker Quote Responder finished or was interrupted."
