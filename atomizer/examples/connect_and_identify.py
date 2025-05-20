# atomizer/examples/connect_and_identify.py
#
# Example script to connect to a Rysk network and confirm the wallet identity
# using the ryskV12_py SDK.

import os
from ryskV12.client import Rysk, Env

# Configuration (users might need to adjust this)
# Ensure ryskV12-cli is installed and configured.
# The SDK might find the CLI automatically if it's in the PATH or ./ryskV12,
# or you might need to provide the path explicitly.

# For real usage, the private key should be loaded securely, e.g., from an environment variable or a secure vault.
# For this example, we'll use a placeholder. If the ryskV12-cli is configured
# with a default private key, the SDK might pick that up, or it might require it here.
# The SDK's Rysk constructor requires a private_key.
EXAMPLE_PRIVATE_KEY = os.environ.get("RYSK_PRIVATE_KEY", "0xYOUR_PRIVATE_KEY_HERE_IF_NOT_SET_IN_CLI_CONFIG")

# Optional: Path to the ryskV12 CLI executable.
# If None, the SDK defaults to "./ryskV12".
RYSKV12_CLI_PATH = os.environ.get("RYSKV12_CLI_PATH") # Or set to an absolute path e.g., "/usr/local/bin/ryskV12"

# Network environment: Env.TESTNET, Env.MAINNET, or Env.LOCAL
SELECTED_ENV = Env.TESTNET # Example: Use Rysk Testnet

def get_wallet_address_from_private_key(private_key: str) -> str:
    """
    Placeholder function to derive a public address from a private key.
    In a real scenario, this would use a proper cryptographic library (e.g., eth_keys, web3.py).
    This is a simplified mock for example purposes, as ryskV12_py SDK doesn't directly expose a get_address method.
    The CLI likely handles this internally.
    """
    if private_key and private_key != "0xYOUR_PRIVATE_KEY_HERE_IF_NOT_SET_IN_CLI_CONFIG" and len(private_key) > 5:
        # This is NOT a real derivation. Just a mock.
        return "0x" + private_key[2:12] + "..." + private_key[-10:]
    return "unknown_address_due_to_missing_pk"

def main():
    print(f"Attempting to initialize Rysk client for environment: {SELECTED_ENV.name}...")
    wallet_address = "N/A"

    try:
        # Initialize the Rysk client
        if RYSKV12_CLI_PATH:
            print(f"Using CLI path: {RYSKV12_CLI_PATH}")
            client = Rysk(env=SELECTED_ENV, private_key=EXAMPLE_PRIVATE_KEY, v12_cli_path=RYSKV12_CLI_PATH)
        else:
            print("Using default CLI path ('./ryskV12')")
            client = Rysk(env=SELECTED_ENV, private_key=EXAMPLE_PRIVATE_KEY)

        print("Rysk client initialized.")

        # The ryskV12_py SDK does not seem to have a direct method like `client.get_wallet_address()`.
        # The CLI it wraps likely uses the configured private key to determine the wallet address.
        # For this script, we'll simulate deriving it or acknowledge it's managed by the CLI.
        # In a real application, you might already know your public address corresponding to the private key.
        print("Attempting to identify wallet address...")
        # wallet_address = client.get_some_identifier_method() # Replace if an actual method is discovered

        # As a placeholder, we can "derive" it from the private key for display purposes.
        # Or, if the CLI has a command to show current address, the SDK could wrap that.
        # The SDK's README example used a `public_address` variable separately.
        wallet_address = get_wallet_address_from_private_key(EXAMPLE_PRIVATE_KEY)

        if wallet_address != "unknown_address_due_to_missing_pk":
            print(f"Successfully initialized client for {SELECTED_ENV.name}.")
            print(f"Associated Wallet Address (derived for example): {wallet_address}")
            print("Note: The SDK interacts with the CLI, which uses the private key for operations.")
            print("A direct 'get_wallet_address' SDK method was not identified during exploration.")
        else:
            print(f"Client initialized for {SELECTED_ENV.name}, but wallet address cannot be displayed (private key placeholder used).")
            print("Ensure your private key is correctly configured for actual operations.")


    except FileNotFoundError:
        cli_path_info = RYSKV12_CLI_PATH if RYSKV12_CLI_PATH else "./ryskV12"
        print(f"ERROR: ryskV12 CLI executable not found at '{cli_path_info}'.")
        print("Ensure it's installed, executable, and the path is correct (RYSKV12_CLI_PATH or in ./).")
    except Exception as e:
        print(f"An error occurred: {e}")
        print("This could be due to incorrect CLI setup, network issues, or other configuration problems.")

if __name__ == "__main__":
    main()
