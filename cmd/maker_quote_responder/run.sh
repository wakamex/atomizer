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
# Optional: Exchange to use for hedging
EXCHANGE_DEFAULT="deribit"
# Optional: Exchange test mode
EXCHANGE_TEST_MODE_DEFAULT="false"
# Optional: HTTP API port for manual trades
HTTP_PORT_DEFAULT="8080"
# Optional: Maximum position delta exposure
MAX_POSITION_DELTA_DEFAULT="10.0"
# Optional: Minimum liquidity score for trades
MIN_LIQUIDITY_SCORE_DEFAULT="0.001"
# Optional: Enable manual trades via HTTP API
ENABLE_MANUAL_TRADES_DEFAULT="true"
# Optional: Enable gamma hedging
ENABLE_GAMMA_HEDGING_DEFAULT="false"
# Optional: Gamma threshold for hedging
GAMMA_THRESHOLD_DEFAULT="0.1"
# Optional: Cache backend (file or valkey)
CACHE_BACKEND_DEFAULT="file"
# Optional: Valkey server address
VALKEY_ADDR_DEFAULT="localhost:6379"

# Use environment variables if set, otherwise use defaults
MAKER_ADDRESS="${MAKER_ADDRESS:-$MAKER_ADDRESS_DEFAULT}"
WEBSOCKET_URL="${WEBSOCKET_URL:-$WEBSOCKET_URL_DEFAULT}"
RFQ_ASSET_ADDRESSES="${RFQ_ASSET_ADDRESSES:-$RFQ_ASSET_ADDRESSES_DEFAULT}"
DUMMY_PRICE="${DUMMY_PRICE:-$DUMMY_PRICE_DEFAULT}"
QUOTE_VALID_DURATION_SECONDS="${QUOTE_VALID_DURATION_SECONDS:-$QUOTE_VALID_DURATION_SECONDS_DEFAULT}"
EXCHANGE="${EXCHANGE:-$EXCHANGE_DEFAULT}"
EXCHANGE_TEST_MODE="${EXCHANGE_TEST_MODE:-$EXCHANGE_TEST_MODE_DEFAULT}"
HTTP_PORT="${HTTP_PORT:-$HTTP_PORT_DEFAULT}"
MAX_POSITION_DELTA="${MAX_POSITION_DELTA:-$MAX_POSITION_DELTA_DEFAULT}"
MIN_LIQUIDITY_SCORE="${MIN_LIQUIDITY_SCORE:-$MIN_LIQUIDITY_SCORE_DEFAULT}"
ENABLE_MANUAL_TRADES="${ENABLE_MANUAL_TRADES:-$ENABLE_MANUAL_TRADES_DEFAULT}"
ENABLE_GAMMA_HEDGING="${ENABLE_GAMMA_HEDGING:-$ENABLE_GAMMA_HEDGING_DEFAULT}"
GAMMA_THRESHOLD="${GAMMA_THRESHOLD:-$GAMMA_THRESHOLD_DEFAULT}"
CACHE_BACKEND="${CACHE_BACKEND:-$CACHE_BACKEND_DEFAULT}"
VALKEY_ADDR="${VALKEY_ADDR:-$VALKEY_ADDR_DEFAULT}"

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

# Print details before running
echo "Maker Address: $MAKER_ADDRESS"
echo "WebSocket URL: $WEBSOCKET_URL"
echo "RFQ Asset Addresses: $RFQ_ASSET_ADDRESSES"
echo "Dummy Price: $DUMMY_PRICE"
echo "Quote Valid Duration Seconds: $QUOTE_VALID_DURATION_SECONDS"
echo "Exchange: $EXCHANGE"
echo "Exchange Test Mode: $EXCHANGE_TEST_MODE"
echo "HTTP API Port: $HTTP_PORT"
echo "Max Position Delta: $MAX_POSITION_DELTA"
echo "Min Liquidity Score: $MIN_LIQUIDITY_SCORE"
echo "Enable Manual Trades: $ENABLE_MANUAL_TRADES"
echo "Enable Gamma Hedging: $ENABLE_GAMMA_HEDGING"
echo "Gamma Threshold: $GAMMA_THRESHOLD"
echo "Cache Backend: $CACHE_BACKEND"
echo "Valkey Address: $VALKEY_ADDR"

# Print Derive-specific variables if using Derive exchange
if [ "$EXCHANGE" = "derive" ]; then
    echo "Derive Wallet Address: ${DERIVE_WALLET_ADDRESS:-Not set}"
    echo "Derive Subaccount ID: ${DERIVE_SUBACCOUNT_ID:-Not set}"
    if [ -n "$DERIVE_PRIVATE_KEY" ]; then
        echo "Derive Private Key: Set (using DERIVE_PRIVATE_KEY)"
    else
        echo "Derive Private Key: Using PRIVATE_KEY"
    fi
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
    --quote_valid_duration_seconds="$QUOTE_VALID_DURATION_SECONDS" \
    --exchange="$EXCHANGE" \
    --exchange_test_mode="$EXCHANGE_TEST_MODE" \
    --http_port="$HTTP_PORT" \
    --max_position_delta="$MAX_POSITION_DELTA" \
    --min_liquidity_score="$MIN_LIQUIDITY_SCORE" \
    --enable_manual_trades="$ENABLE_MANUAL_TRADES" \
    --enable_gamma_hedging="$ENABLE_GAMMA_HEDGING" \
    --gamma_threshold="$GAMMA_THRESHOLD" \
    --cache_backend="$CACHE_BACKEND" \
    --valkey_addr="$VALKEY_ADDR"

echo "Maker Quote Responder finished or was interrupted."
