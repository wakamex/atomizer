# Maker Quote Responder Application

This application demonstrates how to connect to the Rysk Finance WebSocket API, listen for Request for Quotes (RFQs) on specific asset streams, and respond with dummy quotes. It showcases a modular Go application structure with persistent connection management and retry logic.

## Prerequisites

1.  **Go**: Ensure Go (version 1.18 or later recommended) is installed. See [Go Installation Guide](https://go.dev/doc/install).
2.  **Environment File (`.env`)**: Create a `.env` file in the current directory (`cmd/maker_quote_responder/`) with the following content:
    ```env
    # Private key for signing quotes (without 0x prefix)
    PRIVATE_KEY=your_private_key_here 
    ```
    Replace `your_private_key_here` with your actual private key.

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
*   `RFQ_ASSET_ADDRESSES_CSV`: Comma-separated list of RFQ asset addresses (e.g., `0xb67bfa7b488df4f2efa874f4e59242e9130ae61f`). This will be passed to the `--rfq_asset_addresses_csv` flag.
*   `DUMMY_PRICE`: The dummy price to use for quotes (e.g., `12500000000000000000`).
*   `QUOTE_VALID_DURATION_SECONDS`: Duration in seconds for how long the quote should be valid (e.g., `45`).

### Manual Execution (without `run.sh`)

If you prefer to run the built executable manually:

```bash
MAKER_ADDRESS="your_maker_address" \
env $(cat .env | grep -v '^#' | xargs) \
./maker_quote_responder \
  --websocket_url="wss://your_websocket_url/maker" \
  --rfq_asset_addresses_csv="your_rfq_asset_address1,your_rfq_asset_address2" \
  --dummy_price="your_dummy_price" \
  --quote_valid_duration_seconds=your_duration
```

## Environment Variables & Command-Line Flags

### Environment Variables (loaded from `.env` by `run.sh` or set manually)

*   `PRIVATE_KEY`: (Required) The private key (without `0x` prefix) used for signing quote responses.

### Command-Line Flags (passed to the executable)

*   `--maker_address` or `MAKER_ADDRESS` env var: (Required) Your Ethereum maker address.
*   `--websocket_url`: (Required) The WebSocket URL for the Rysk Finance maker API.
*   `--rfq_asset_addresses_csv` or `RFQ_ASSET_ADDRESSES_CSV` env var: (Required) Comma-separated list of asset addresses to subscribe to for RFQs.
*   `--dummy_price`: (Required) The dummy price to be used in quote responses.
*   `--quote_valid_duration_seconds`: (Optional, default: `30`) The duration in seconds for which the quote is valid.

(Note: The application's `LoadConfig()` function prioritizes command-line flags over environment variables if both are set for the same parameter, except for `PRIVATE_KEY` which is only read from env.)
