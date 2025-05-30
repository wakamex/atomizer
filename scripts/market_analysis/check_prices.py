#!/usr/bin/env python3
"""
Quick script to check actual price values from both exchanges
"""

import requests
import pandas as pd
from datetime import datetime, timedelta
import json

def query_prices(vm_url, instrument_derive, instrument_deribit, minutes_back=10):
    """Query recent prices for specific instruments"""
    end_time = datetime.now()
    start_time = end_time - timedelta(minutes=minutes_back)
    
    # Query both instruments
    query = f'''
    {{
        __name__="market_bid_price",
        instrument=~"{instrument_derive}|{instrument_deribit}"
    }}
    '''
    
    url = f"{vm_url}/api/v1/query_range"
    params = {
        "query": f'market_bid_price{{instrument=~"{instrument_derive}|{instrument_deribit}"}}',
        "start": int(start_time.timestamp()),
        "end": int(end_time.timestamp()),
        "step": "30s"
    }
    
    response = requests.get(url, params=params)
    response.raise_for_status()
    
    data = response.json()
    if data["status"] != "success":
        raise Exception(f"Query failed: {data}")
    
    # Parse results
    results = []
    for series in data["data"]["result"]:
        metric = series["metric"]
        exchange = metric.get("exchange", "unknown")
        instrument = metric.get("instrument", "unknown")
        
        for timestamp, value in series["values"]:
            results.append({
                "timestamp": datetime.fromtimestamp(float(timestamp)),
                "exchange": exchange,
                "instrument": instrument,
                "price": float(value)
            })
    
    df = pd.DataFrame(results)
    return df

def main():
    vm_url = "http://localhost:8428"
    
    # Check ETH option prices
    print("Checking ETH-20250601-2600-C prices...")
    df = query_prices(vm_url, "ETH-20250601-2600-C", "ETH-1JUN25-2600-C", minutes_back=5)
    
    if df.empty:
        print("No data found")
        return
    
    # Show recent prices
    print("\nRecent prices:")
    for exchange in df['exchange'].unique():
        exchange_df = df[df['exchange'] == exchange]
        instrument = exchange_df['instrument'].iloc[0]
        print(f"\n{exchange} - {instrument}:")
        recent = exchange_df.tail(5)
        for _, row in recent.iterrows():
            print(f"  {row['timestamp']}: {row['price']:.4f}")
    
    # Calculate rough USD value for Deribit
    eth_price_estimate = 3700  # Rough estimate
    print(f"\nAssuming ETH price ~${eth_price_estimate}:")
    
    deribit_df = df[df['exchange'] == 'deribit']
    if not deribit_df.empty:
        deribit_price_eth = deribit_df['price'].iloc[-1]
        deribit_price_usd = deribit_price_eth * eth_price_estimate
        print(f"Deribit price: {deribit_price_eth:.4f} ETH ≈ ${deribit_price_usd:.2f}")
    
    derive_df = df[df['exchange'] == 'derive']
    if not derive_df.empty:
        derive_price = derive_df['price'].iloc[-1]
        print(f"Derive price: ${derive_price:.2f}")

if __name__ == "__main__":
    main()