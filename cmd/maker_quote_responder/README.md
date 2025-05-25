# Maker Quote Responder Application

This application is a production-ready market maker that connects to the Rysk Finance WebSocket API, listens for Request for Quotes (RFQs), responds with competitive quotes using real-time Deribit prices, and automatically hedges positions on Deribit when trades are executed.

## Prerequisites

1.  **Go**: Ensure Go (version 1.18 or later recommended) is installed. See [Go Installation Guide](https://go.dev/doc/install).
2.  **Environment File (`.env`)**: Create a `.env` file in the current directory (`cmd/maker_quote_responder/`) with the following content:
    ```env
    # Private key for signing quotes (without 0x prefix)
    PRIVATE_KEY=your_private_key_here
    
    # Deribit API credentials (for hedging)
    DERIBIT_API_KEY=your_deribit_testnet_key
    DERIBIT_API_SECRET=your_deribit_testnet_secret
    ```
    Replace the placeholders with your actual credentials.

## Building the Application

To build the `maker_quote_responder` executable:

1.  Navigate to the application directory:
    ```bash
    cd /path/to/atomizer/cmd/maker_quote_responder
    ```
    (Or, if you are already in the project root `/path/to/atomizer`, you can often build with `go build -o cmd/maker_quote_responder/maker_quote_responder ./cmd/maker_quote_responder`)

2.  Run the Go build command (from within `cmd/maker_quote_responder/`):
    ```bash
    go build -o maker_quote_responder .
    ```
    This will create an executable file named `maker_quote_responder` in the `cmd/maker_quote_responder/` directory.

## Running the Application

A `run.sh` script is provided for convenience. It loads environment variables from the `.env` file and executes the built application with the required command-line flags.

1.  **Make the script executable** (if you haven't already):
    ```bash
    chmod +x run.sh
    ```
2.  **Execute the script**:
    ```bash
    ./run.sh
    ```

### Script Configuration

The `run.sh` script is pre-configured with default values for command-line arguments. You can modify `run.sh` to change these values if needed:

*   `MAKER_ADDRESS`: Your Ethereum maker address (e.g., `0x9eAFc0c2b04D96a1C1edAdda8A474a4506752207`).
*   `WEBSOCKET_URL`: The Rysk Finance WebSocket URL (e.g., `wss://rip-testnet.rysk.finance/maker`).
*   `RFQ_ASSET_ADDRESSES`: Comma-separated list of RFQ asset addresses (e.g., `0xb67bfa7b488df4f2efa874f4e59242e9130ae61f`). This will be passed to the `--rfq_asset_addresses` flag.
*   `DUMMY_PRICE`: The fallback price when Deribit pricing fails (e.g., `12500000000000000000`).
*   `QUOTE_VALID_DURATION_SECONDS`: Duration in seconds for how long the quote should be valid (e.g., `45`).
*   `ASSET_MAPPING`: JSON mapping of asset addresses to underlying symbols (e.g., `{"0xb67bfa7b488df4f2efa874f4e59242e9130ae61f":"ETH"}`)

### Manual Execution (without `run.sh`)

If you prefer to run the built executable manually:

```bash
MAKER_ADDRESS="your_maker_address" \
env $(cat .env | grep -v '^#' | xargs) \
./maker_quote_responder \
  --websocket_url="wss://your_websocket_url/maker" \
  --rfq_asset_addresses="your_rfq_asset_address1,your_rfq_asset_address2" \
  --dummy_price="your_dummy_price" \
  --quote_valid_duration_seconds=your_duration
```

## Environment Variables & Command-Line Flags

### Environment Variables (loaded from `.env` by `run.sh` or set manually)

*   `PRIVATE_KEY`: (Required) The private key (without `0x` prefix) used for signing quote responses.

### Command-Line Flags (passed to the executable)

*   `--maker_address` or `MAKER_ADDRESS` env var: (Required) Your Ethereum maker address.
*   `--websocket_url`: (Required) The WebSocket URL for the Rysk Finance maker API.
*   `--rfq_asset_addresses` or `RFQ_ASSET_ADDRESSES` env var: (Required) Comma-separated list of asset addresses to subscribe to for RFQs.
*   `--dummy_price`: (Required) The dummy price to be used in quote responses.
*   `--quote_valid_duration_seconds`: (Optional, default: `30`) The duration in seconds for which the quote is valid.

(Note: The application's `LoadConfig()` function prioritizes command-line flags over environment variables if both are set for the same parameter, except for `PRIVATE_KEY` which is only read from env.)

## Key Features

### 1. Real-time Pricing with Deribit
The application fetches live option prices from Deribit to generate competitive quotes:
- Converts RFQ parameters to Deribit instrument format
- Calculates prices based on order book depth
- Applies configurable market maker premium/discount
- Falls back to dummy pricing if Deribit is unavailable

### 2. Automatic Hedging
When a trade is executed on RyskV12:
- Receives trade confirmation via WebSocket
- Automatically places hedge order on Deribit (testnet)
- Uses fill-or-kill orders with 10% slippage protection
- Logs all hedge attempts and results

### 3. Asset Mapping
Configure which assets to quote on using the `ASSET_MAPPING` environment variable:
```json
{
  "0xb67bfa7b488df4f2efa874f4e59242e9130ae61f": "ETH",
  "0x1234567890123456789012345678901234567890": "BTC"
}
```

## Testing

The application includes comprehensive tests:
- `quoter_test.go`: Unit tests for pricing logic
- `hedge_test.go`: Tests for hedging functionality
- `integration_test.go`: End-to-end testing

Run tests with:
```bash
go test -v ./...
```

## Monitoring

The application logs all important events:
- RFQ receipts and quote responses
- Deribit pricing attempts and failures
- Hedge order placements and results
- Connection status and errors

## Limitations

- Currently supports only call options (puts return an error)
- Hedging is available only on Deribit testnet
- Requires manual asset mapping configuration
