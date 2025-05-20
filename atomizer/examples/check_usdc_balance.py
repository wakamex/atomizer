# atomizer/examples/check_usdc_balance.py
#
# Example script to check USDC balance using the ryskV12_py SDK.

import os
import json # For parsing potential JSON output from CLI
from ryskV12.client import Rysk, Env

# Configuration (mirroring connect_and_identify.py for consistency)
EXAMPLE_PRIVATE_KEY = os.environ.get("RYSK_PRIVATE_KEY", "0xYOUR_PRIVATE_KEY_HERE_IF_NOT_SET_IN_CLI_CONFIG")
RYSKV12_CLI_PATH = os.environ.get("RYSKV12_CLI_PATH")
SELECTED_ENV = Env.TESTNET
TOKEN_SYMBOL_TO_CHECK = "USDC" # Token to check balance for

# The account address associated with your private key.
# This might be needed for the balances_args method.
# For this example, users should replace this with their actual account address.
# Similar to connect_and_identify.py, a real app would have a secure way to get this.
EXAMPLE_ACCOUNT_ADDRESS = os.environ.get("RYSK_ACCOUNT_ADDRESS", "0xYOUR_ACCOUNT_ADDRESS_HERE")


def get_placeholder_wallet_address(private_key: str) -> str:
    """
    Placeholder function to derive/get a public address.
    In a real scenario, this would use a proper cryptographic library (e.g., eth_keys, web3.py)
    or be fetched from a secure source if not derivable directly in this context.
    This is used if RYSK_ACCOUNT_ADDRESS is not set.
    """
    if private_key and private_key != "0xYOUR_PRIVATE_KEY_HERE_IF_NOT_SET_IN_CLI_CONFIG" and len(private_key) > 12:
        # This is NOT a real derivation. Just a mock for example if account address not provided.
        return "0x" + private_key[2:12] # Example mock, actual derivation is complex.
    return "0x0000000000000000000000000000000000000000" # Default if no PK

def main():
    print(f"Attempting to check {TOKEN_SYMBOL_TO_CHECK} balance on {SELECTED_ENV.name}...")

    account_address = EXAMPLE_ACCOUNT_ADDRESS
    if account_address == "0xYOUR_ACCOUNT_ADDRESS_HERE":
        print("RYSK_ACCOUNT_ADDRESS environment variable not set, using placeholder derived from private key.")
        account_address = get_placeholder_wallet_address(EXAMPLE_PRIVATE_KEY)
        print(f"Using placeholder account address: {account_address}")


    if account_address == "0x0000000000000000000000000000000000000000" and \
       EXAMPLE_PRIVATE_KEY == "0xYOUR_PRIVATE_KEY_HERE_IF_NOT_SET_IN_CLI_CONFIG":
        print("Error: Cannot determine account address. Please set RYSK_ACCOUNT_ADDRESS or RYSK_PRIVATE_KEY.")
        return

    try:
        if RYSKV12_CLI_PATH:
            client = Rysk(env=SELECTED_ENV, private_key=EXAMPLE_PRIVATE_KEY, v12_cli_path=RYSKV12_CLI_PATH)
        else:
            client = Rysk(env=SELECTED_ENV, private_key=EXAMPLE_PRIVATE_KEY)
        print("Rysk client initialized.")

        # The SDK uses `balances_args(channel_id, account)` and `execute`.
        # It does not seem to have a high-level `get_balance(token='USDC')` method.
        # The `balances_args` method is used to get arguments for the CLI.
        # The output from the CLI needs to be parsed.
        channel_id = "my_balance_channel" # User-defined channel ID for this operation
        print(f"Requesting balances for account: {account_address} via channel: {channel_id}")

        # Form the arguments for the balances command
        balance_args = client.balances_args(channel_id=channel_id, account=account_address)
        
        # Execute the command
        # This returns a Popen object. We need to read stdout.
        print(f"Executing CLI command: {client._cli_path} {' '.join(balance_args)}") # Show what's being run
        process = client.execute(balance_args)
        stdout, stderr = process.communicate() # Wait for command to complete

        if process.returncode == 0:
            print("CLI command executed successfully.")
            print(f"Raw output from CLI:\n{stdout}")

            # Now, parse stdout to find the USDC balance.
            # The format of stdout is unknown. It could be JSON, plain text, etc.
            # We'll assume it might be JSON lines or a simple text format.
            # This parsing logic is a placeholder and needs to be adapted based on actual CLI output.
            found_balance = None
            try:
                # Attempt to parse as JSON (maybe each line is a JSON object, or the whole thing is)
                # This is a common pattern for CLI tools.
                lines = stdout.strip().split('\n')
                for line in lines:
                    try:
                        data = json.loads(line)
                        # Assuming data is a dict and might look like:
                        # {'token': 'USDC', 'balance': '123.45'}
                        # or {'asset': 'USDC', 'amount': '123450000', 'decimals': 6}
                        if isinstance(data, dict) and data.get('token') == TOKEN_SYMBOL_TO_CHECK:
                            found_balance = data.get('balance', 'Not specified')
                            break
                        elif isinstance(data, dict) and data.get('asset') == TOKEN_SYMBOL_TO_CHECK:
                            # If amount and decimals are provided
                            if 'amount' in data and 'decimals' in data:
                                amount = int(data['amount'])
                                decimals = int(data['decimals'])
                                found_balance = str(amount / (10**decimals))
                            else:
                                found_balance = data.get('amount', 'Not specified')
                            break
                        # Add more parsing rules if needed based on actual output
                    except json.JSONDecodeError:
                        # Line is not JSON, try simple string search
                        if TOKEN_SYMBOL_TO_CHECK in line and "balance" in line.lower():
                            # Extremely basic text parsing, highly dependent on format
                            parts = line.split()
                            try:
                                # Try to find a number after the token symbol
                                token_idx = parts.index(TOKEN_SYMBOL_TO_CHECK)
                                found_balance = parts[token_idx + 1] # Very fragile
                            except (ValueError, IndexError):
                                pass # Fallback to line itself if simple parse fails
                            if found_balance is None: found_balance = line # Use the line as a fallback
                            break


                if found_balance is not None:
                    print(f"Balance of {TOKEN_SYMBOL_TO_CHECK}: {found_balance}")
                else:
                    print(f"Could not find {TOKEN_SYMBOL_TO_CHECK} balance in the output. Full output was logged above.")

            except Exception as e:
                print(f"Error parsing balance output: {e}")
                print("The raw output was logged above. You may need to adjust parsing logic.")
        else:
            print(f"CLI command failed with return code {process.returncode}.")
            print(f"Error output (stderr):\n{stderr}")

    except FileNotFoundError:
        cli_path_info = RYSKV12_CLI_PATH if RYSKV12_CLI_PATH else "./ryskV12"
        print(f"ERROR: ryskV12 CLI executable not found at '{cli_path_info}'.")
    except Exception as e:
        print(f"An error occurred while checking balance: {e}")

if __name__ == "__main__":
    main()
