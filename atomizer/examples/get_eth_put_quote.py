# atomizer/examples/get_eth_put_quote.py
#
# Example script to get a quote for an ETH put option using the ryskV12_py SDK.

import os
import json # For parsing JSON string output if necessary
from ryskV12.client import Rysk, Env
# from ryskV12.models import OptionType, Series # OptionType/Series not seen in models.py exploration. Using boolean for isPut.
# The Quote model is for sending quotes. For getting, we assume a CLI command.

# Configuration
EXAMPLE_PRIVATE_KEY = os.environ.get("RYSK_PRIVATE_KEY", "0xYOUR_PRIVATE_KEY_HERE_IF_NOT_SET_IN_CLI_CONFIG")
RYSKV12_CLI_PATH = os.environ.get("RYSKV12_CLI_PATH")
SELECTED_ENV = Env.TESTNET

# Option Parameters for the quote
OPTION_PARAMS = {
    "underlying_symbol": "ETH",
    "is_put_option": True, # True for Put, False for Call
    "strike_price_usd": "2000", # Ensure type matches SDK/CLI (string for instrument name)
    "expiry_date": "20241231"  # YYYYMMDD format
}

def main():
    option_type_str = "PUT" if OPTION_PARAMS["is_put_option"] else "CALL"
    print(f"Attempting to get quote for {OPTION_PARAMS['underlying_symbol']} {option_type_str} "
          f"Strike: {OPTION_PARAMS['strike_price_usd']} Expiry: {OPTION_PARAMS['expiry_date']} on {SELECTED_ENV.name}...")

    try:
        if RYSKV12_CLI_PATH:
            client = Rysk(env=SELECTED_ENV, private_key=EXAMPLE_PRIVATE_KEY, v12_cli_path=RYSKV12_CLI_PATH)
        else:
            client = Rysk(env=SELECTED_ENV, private_key=EXAMPLE_PRIVATE_KEY)
        print("Rysk client initialized.")

        # Construct the instrument name string.
        # The actual format depends on how the ryskV12-cli expects it.
        # This is a common convention: SYMBOL-YYYYMMDD-STRIKE-P_or_C
        instrument_name = (
            f"{OPTION_PARAMS['underlying_symbol']}-"
            f"{OPTION_PARAMS['expiry_date']}-"
            f"{OPTION_PARAMS['strike_price_usd']}-"
            f"{'P' if OPTION_PARAMS['is_put_option'] else 'C'}"
        )
        print(f"Requesting quote for instrument: {instrument_name}")
        amount_to_quote = "1" # Example: 1 contract

        # Assumption: The ryskV12-cli has a command like `quote <instrument_name> --amount <amount>`
        # to get a quote from the market. The ryskV12_py SDK's `quote_args` and `Quote` model
        # are for *sending* a quote (acting as a maker), not for requesting one (acting as a taker).
        # We use the generic `execute` method to call this assumed CLI command.
        
        cli_command_args = ["quote", instrument_name, "--amount", amount_to_quote]
        
        print(f"Executing CLI command: {client._cli_path} {' '.join(cli_command_args)}")
        process = client.execute(cli_command_args)
        stdout, stderr = process.communicate() # Wait for command to complete

        if process.returncode == 0:
            print("CLI command executed successfully.")
            print(f"Raw output from CLI:\n{stdout}")

            # Parse the stdout to find the quote details.
            # The format of stdout is unknown (JSON, plain text, etc.).
            # This parsing logic is a placeholder.
            try:
                # Attempt to parse as JSON
                quote_data = json.loads(stdout)
                print("Successfully parsed quote data (assumed JSON).")
                # Example: quote_data might be {'bid': '10.0', 'ask': '10.5', 'iv': '0.65'}
                # These keys are illustrative.
                print(f"  Bid: {quote_data.get('bid', 'N/A')}")
                print(f"  Ask: {quote_data.get('ask', 'N/A')}")
                print(f"  Mid Price: {quote_data.get('mid_price', 'N/A')}") # If available
                print(f"  Implied Volatility: {quote_data.get('iv', 'N/A')}") # If available
                print(f"  Timestamp: {quote_data.get('timestamp', 'N/A')}") # If available

            except json.JSONDecodeError:
                print("Could not parse quote output as JSON.")
                print("You may need to implement custom parsing for the CLI's output format if it's not JSON.")
                # Add more robust parsing here if the format is known (e.g., regex for plain text)
            except Exception as e:
                print(f"Error parsing quote output: {e}")
        else:
            print(f"CLI command failed with return code {process.returncode}.")
            print(f"Error output (stderr):\n{stderr}")

    except FileNotFoundError:
        cli_path_info = RYSKV12_CLI_PATH if RYSKV12_CLI_PATH else "./ryskV12"
        print(f"ERROR: ryskV12 CLI executable not found at '{cli_path_info}'.")
    except Exception as e:
        print(f"An error occurred while getting quote: {e}")

if __name__ == "__main__":
    main()
