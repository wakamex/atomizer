#!/usr/bin/env python3
import json
import time
from web3 import Web3
from eth_account.messages import encode_defunct
import websocket
import os

# Configuration
PRIVATE_KEY = os.getenv("DERIVE_PRIVATE_KEY") or os.getenv("PRIVATE_KEY")
DERIVE_WALLET = os.getenv("DERIVE_WALLET_ADDRESS")
WEBSOCKET_URL = "wss://api.lyra.finance/ws"

def sign_message(message, private_key):
    """Sign a message with the private key"""
    w3 = Web3()
    account = w3.eth.account.from_key(private_key)
    message_hash = encode_defunct(text=message)
    signed = account.sign_message(message_hash)
    return signed.signature.hex()

def main():
    print(f"Private Key: {PRIVATE_KEY[:10]}...")
    print(f"Derive Wallet: {DERIVE_WALLET}")
    
    # Generate timestamp
    timestamp = str(int(time.time() * 1000))
    print(f"Timestamp: {timestamp}")
    
    # Sign the timestamp
    signature = sign_message(timestamp, PRIVATE_KEY)
    print(f"Signature: 0x{signature}")
    
    # Create login message
    login_msg = {
        "jsonrpc": "2.0",
        "method": "public/login",
        "params": {
            "wallet": DERIVE_WALLET,
            "timestamp": timestamp,
            "signature": f"0x{signature}"
        },
        "id": timestamp
    }
    
    print(f"\nLogin message:")
    print(json.dumps(login_msg, indent=2))
    
    # Connect and send
    try:
        ws = websocket.create_connection(WEBSOCKET_URL)
        print(f"\nConnected to {WEBSOCKET_URL}")
        
        # Send login
        ws.send(json.dumps(login_msg))
        print("Sent login message")
        
        # Get response
        response = ws.recv()
        print(f"\nResponse: {response}")
        
        ws.close()
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    main()