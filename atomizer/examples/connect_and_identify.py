"""Connect to a Rysk network and confirm the wallet identity."""

import os
from ryskV12.client import Rysk, Env
from atomizer.utils.loader import load_env_vars
from atomizer.utils.keys import get_wallet_address_from_private_key

load_env_vars()

PRIVATE_KEY = os.environ.get("PRIVATE_KEY")
SELECTED_ENV = Env.TESTNET

def main():
    print(f"Attempting to initialize Rysk client for environment: {SELECTED_ENV.name}...")
    wallet_address = "N/A"

    try:
        # Initialize the Rysk client
        print("Initializing Rysk client. The SDK will attempt to locate the ryskV12 CLI.")
        client = Rysk(env=SELECTED_ENV, private_key=PRIVATE_KEY)
        print("Rysk client initialized.")

        print("Attempting to identify wallet address...")
        wallet_address = get_wallet_address_from_private_key(PRIVATE_KEY)
        print(f"Wallet address: {wallet_address}")

        if wallet_address != "unknown_address_due_to_missing_pk":
            print(f"Successfully initialized client for {SELECTED_ENV.name}.")
            print(f"Associated Wallet Address: {wallet_address}")
        else:
            print(f"Client initialized for {SELECTED_ENV.name}, but wallet address cannot be displayed (private key placeholder used).")
            print("Ensure your private key is correctly configured for actual operations.")
    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    main()
