# Atomizer: ryskV12_py SDK Examples

This directory contains example Python scripts demonstrating basic interactions with the Rysk platform using the `ryskV12_py` SDK.

## Scripts

1.  **`connect_and_identify.py`**:
    Initializes the Rysk client, connects to a specified network (default: Testnet), and displays a placeholder for the wallet address associated with the configured private key. *Note: The wallet address display is a simplified placeholder as the SDK does not provide a direct method; the underlying CLI manages identity via its private key configuration.*

2.  **`check_usdc_balance.py`**:
    Connects to the Rysk network and queries the balance for USDC (or a configured token) for the active wallet.

3.  **`get_eth_put_quote.py`**:
    Connects to the Rysk network and fetches a quote for a pre-defined ETH put option (details like strike and expiry are configurable within the script).

## Prerequisites

Before running these examples, ensure you have the following set up:

1.  **`ryskV12_py` SDK Installed**: The Python SDK must be installed in your environment.
    ```bash
    pip install ryskV12_py # Or however it's specified in your project
    ```

2.  **`ryskV12-cli` Executable**:
    *   The `ryskV12-cli` command-line tool (which `ryskV12_py` wraps) must be installed and executable.
    *   It needs to be configured with your private key for operations on the desired network (Testnet/Mainnet). Refer to the `ryskV12-cli` documentation for setup.
    *   The scripts assume the CLI can be found, either in your system's PATH, in the current directory (`./ryskV12`), or via the `RYSKV12_CLI_PATH` environment variable if you set it to the executable's location.

3.  **Environment Variables (Optional but Recommended for Scripts)**:
    *   `RYSK_PRIVATE_KEY`: The example scripts use this environment variable as a source for the private key string (e.g., `"0xYOUR_PRIVATE_KEY"`). While the CLI manages the actual key, the SDK client constructor in these examples takes it as an argument. **Ensure this is handled securely for any real use.**
    *   `RYSKV12_CLI_PATH`: If your `ryskV12-cli` is not in a standard location, you can set this variable to its absolute path.

## Running the Scripts

Navigate to this `atomizer/examples/` directory and run the scripts directly using Python:

```bash
python connect_and_identify.py
python check_usdc_balance.py
python get_eth_put_quote.py
```

You can modify parameters (like network, token symbols, option details, or private key source) directly within each script for experimentation.
