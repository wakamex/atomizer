#!/usr/bin/env python3
"""
Inspect orderbook data from VictoriaMetrics
"""

import argparse
import requests
import pandas as pd
from datetime import datetime, timedelta
from tabulate import tabulate
import numpy as np


class OrderBookInspector:
    def __init__(self, vm_url="http://localhost:8428"):
        self.vm_url = vm_url
    
    def query_current(self, query):
        """Query current values from VictoriaMetrics"""
        url = f"{self.vm_url}/api/v1/query"
        params = {"query": query}
        
        response = requests.get(url, params=params)
        response.raise_for_status()
        
        data = response.json()
        if data["status"] != "success":
            raise Exception(f"Query failed: {data}")
        
        return data["data"]["result"]
    
    def query_range(self, query, start_time, end_time, step="5s"):
        """Query time range from VictoriaMetrics"""
        url = f"{self.vm_url}/api/v1/query_range"
        params = {
            "query": query,
            "start": int(start_time.timestamp()),
            "end": int(end_time.timestamp()),
            "step": step
        }
        
        response = requests.get(url, params=params)
        response.raise_for_status()
        
        data = response.json()
        if data["status"] != "success":
            raise Exception(f"Query failed: {data}")
        
        return data["data"]["result"]
    
    def get_current_orderbook(self, instrument, exchange):
        """Get current orderbook snapshot"""
        # Get prices and sizes for all levels
        price_query = f'orderbook_price{{instrument="{instrument}",exchange="{exchange}"}}'
        size_query = f'orderbook_size{{instrument="{instrument}",exchange="{exchange}"}}'
        
        price_results = self.query_current(price_query)
        size_results = self.query_current(size_query)
        
        # Organize by side and level
        orderbook = {"bid": {}, "ask": {}, "metadata": {}}
        
        # Process prices
        for result in price_results:
            side = result["metric"]["side"]
            level = int(result["metric"]["level"])
            price = float(result["value"][1])
            orderbook[side][level] = {"price": price, "size": 0}
        
        # Process sizes
        for result in size_results:
            side = result["metric"]["side"]
            level = int(result["metric"]["level"])
            size = float(result["value"][1])
            if level in orderbook[side]:
                orderbook[side][level]["size"] = size
        
        # Get ETH spot price for Deribit conversion
        if exchange == "deribit" and "ETH" in instrument:
            eth_spot = self.get_eth_spot_price()
            if eth_spot:
                orderbook["metadata"]["eth_spot"] = eth_spot
                # Convert prices to USD
                for side in ["bid", "ask"]:
                    for level in orderbook[side]:
                        orderbook[side][level]["price_usd"] = orderbook[side][level]["price"] * eth_spot
        
        return orderbook
    
    def get_eth_spot_price(self):
        """Get current ETH spot price"""
        query = 'market_bid_price{instrument="ETH-SPOT"}'
        results = self.query_current(query)
        
        if results:
            return float(results[0]["value"][1])
        return None
    
    def display_orderbook(self, orderbook, instrument, exchange):
        """Display orderbook in a nice format"""
        print(f"\n{'='*60}")
        print(f"Order Book: {instrument} on {exchange}")
        print(f"Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        
        # Check if we have ETH spot price for conversion
        has_usd_conversion = "eth_spot" in orderbook.get("metadata", {})
        if has_usd_conversion:
            eth_spot = orderbook["metadata"]["eth_spot"]
            print(f"ETH Spot: ${eth_spot:.2f}")
        
        print(f"{'='*60}")
        
        # Prepare bid and ask data
        bid_data = []
        ask_data = []
        
        # Get max levels
        max_bid_level = max(orderbook["bid"].keys()) if orderbook["bid"] else 0
        max_ask_level = max(orderbook["ask"].keys()) if orderbook["ask"] else 0
        max_level = max(max_bid_level, max_ask_level)
        
        for level in range(1, max_level + 1):
            bid_price = ""
            bid_size = ""
            ask_price = ""
            ask_size = ""
            
            if level in orderbook["bid"]:
                bid_size = f"{orderbook['bid'][level]['size']:.2f}"
                if has_usd_conversion and "price_usd" in orderbook["bid"][level]:
                    # Show both ETH and USD prices for Deribit
                    eth_price = orderbook['bid'][level]['price']
                    usd_price = orderbook['bid'][level]['price_usd']
                    bid_price = f"${usd_price:.2f} ({eth_price:.4f} ETH)"
                else:
                    bid_price = f"${orderbook['bid'][level]['price']:.4f}"
            
            if level in orderbook["ask"]:
                ask_size = f"{orderbook['ask'][level]['size']:.2f}"
                if has_usd_conversion and "price_usd" in orderbook["ask"][level]:
                    # Show both ETH and USD prices for Deribit
                    eth_price = orderbook['ask'][level]['price']
                    usd_price = orderbook['ask'][level]['price_usd']
                    ask_price = f"${usd_price:.2f} ({eth_price:.4f} ETH)"
                else:
                    ask_price = f"${orderbook['ask'][level]['price']:.4f}"
            
            bid_data.append([bid_size, bid_price])
            ask_data.append([ask_price, ask_size])
        
        # Create combined table
        headers = ["Bid Size", "Bid Price", "Level", "Ask Price", "Ask Size"]
        table_data = []
        
        for i, (bid, ask) in enumerate(zip(bid_data, ask_data)):
            table_data.append(bid + [i+1] + ask)
        
        print(tabulate(table_data, headers=headers, tablefmt="grid", floatfmt=".4f"))
        
        # Show spread
        if 1 in orderbook["bid"] and 1 in orderbook["ask"]:
            if has_usd_conversion and "price_usd" in orderbook["bid"][1]:
                # Use USD prices for calculations
                best_bid = orderbook["bid"][1]["price_usd"]
                best_ask = orderbook["ask"][1]["price_usd"]
                best_bid_eth = orderbook["bid"][1]["price"]
                best_ask_eth = orderbook["ask"][1]["price"]
                
                spread = best_ask - best_bid
                spread_pct = (spread / best_bid) * 100
                mid_price = (best_bid + best_ask) / 2
                
                print(f"\nBest Bid: ${best_bid:.2f} ({best_bid_eth:.4f} ETH)")
                print(f"Best Ask: ${best_ask:.2f} ({best_ask_eth:.4f} ETH)")
                print(f"Spread: ${spread:.2f} ({spread_pct:.2f}%)")
                print(f"Mid Price: ${mid_price:.2f}")
            else:
                best_bid = orderbook["bid"][1]["price"]
                best_ask = orderbook["ask"][1]["price"]
                spread = best_ask - best_bid
                spread_pct = (spread / best_bid) * 100
                mid_price = (best_bid + best_ask) / 2
                
                print(f"\nBest Bid: ${best_bid:.4f}")
                print(f"Best Ask: ${best_ask:.4f}")
                print(f"Spread: ${spread:.4f} ({spread_pct:.2f}%)")
                print(f"Mid Price: ${mid_price:.4f}")
    
    def get_orderbook_history(self, instrument, exchange, minutes_back=5):
        """Get orderbook history"""
        end_time = datetime.now()
        start_time = end_time - timedelta(minutes=minutes_back)
        
        # Query mid price history
        query = f'orderbook_mid_price{{instrument="{instrument}",exchange="{exchange}"}}'
        results = self.query_range(query, start_time, end_time)
        
        if not results:
            print(f"No data found for {instrument} on {exchange}")
            return pd.DataFrame()
        
        # Convert to DataFrame
        data = []
        for timestamp, value in results[0]["values"]:
            data.append({
                "timestamp": datetime.fromtimestamp(float(timestamp)),
                "mid_price": float(value)
            })
        
        df = pd.DataFrame(data)
        
        # Add spread data if available
        spread_query = f'orderbook_spread{{instrument="{instrument}",exchange="{exchange}"}}'
        spread_results = self.query_range(spread_query, start_time, end_time)
        
        if spread_results:
            spread_data = []
            for timestamp, value in spread_results[0]["values"]:
                spread_data.append({
                    "timestamp": datetime.fromtimestamp(float(timestamp)),
                    "spread": float(value)
                })
            spread_df = pd.DataFrame(spread_data)
            df = pd.merge(df, spread_df, on="timestamp", how="left")
        
        return df
    
    def compare_orderbooks(self, instrument1, exchange1, instrument2, exchange2):
        """Compare two orderbooks side by side"""
        ob1 = self.get_current_orderbook(instrument1, exchange1)
        ob2 = self.get_current_orderbook(instrument2, exchange2)
        
        print(f"\n{'='*80}")
        print(f"Order Book Comparison")
        print(f"Left:  {instrument1} on {exchange1}")
        print(f"Right: {instrument2} on {exchange2}")
        print(f"Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        print(f"{'='*80}")
        
        # Get max levels
        max_level = max(
            max(ob1["bid"].keys()) if ob1["bid"] else 0,
            max(ob1["ask"].keys()) if ob1["ask"] else 0,
            max(ob2["bid"].keys()) if ob2["bid"] else 0,
            max(ob2["ask"].keys()) if ob2["ask"] else 0
        )
        
        headers = ["Bid Size", "Bid $", "Ask $", "Ask Size", "Level", 
                   "Bid Size", "Bid $", "Ask $", "Ask Size"]
        table_data = []
        
        for level in range(1, max_level + 1):
            row = []
            
            # Left orderbook bid
            if level in ob1["bid"]:
                size = f"{ob1['bid'][level]['size']:.2f}"
                if "price_usd" in ob1["bid"][level]:
                    price = f"{ob1['bid'][level]['price_usd']:.2f}"
                else:
                    price = f"{ob1['bid'][level]['price']:.2f}"
                row.extend([size, price])
            else:
                row.extend(["", ""])
            
            # Left orderbook ask
            if level in ob1["ask"]:
                if "price_usd" in ob1["ask"][level]:
                    price = f"{ob1['ask'][level]['price_usd']:.2f}"
                else:
                    price = f"{ob1['ask'][level]['price']:.2f}"
                size = f"{ob1['ask'][level]['size']:.2f}"
                row.extend([price, size])
            else:
                row.extend(["", ""])
            
            # Level
            row.append(level)
            
            # Right orderbook bid
            if level in ob2["bid"]:
                size = f"{ob2['bid'][level]['size']:.2f}"
                if "price_usd" in ob2["bid"][level]:
                    price = f"{ob2['bid'][level]['price_usd']:.2f}"
                else:
                    price = f"{ob2['bid'][level]['price']:.2f}"
                row.extend([size, price])
            else:
                row.extend(["", ""])
            
            # Right orderbook ask
            if level in ob2["ask"]:
                if "price_usd" in ob2["ask"][level]:
                    price = f"{ob2['ask'][level]['price_usd']:.2f}"
                else:
                    price = f"{ob2['ask'][level]['price']:.2f}"
                size = f"{ob2['ask'][level]['size']:.2f}"
                row.extend([price, size])
            else:
                row.extend(["", ""])
            
            table_data.append(row)
        
        print(tabulate(table_data, headers=headers, tablefmt="grid"))
        
        # Compare spreads and mid prices
        if 1 in ob1["bid"] and 1 in ob1["ask"] and 1 in ob2["bid"] and 1 in ob2["ask"]:
            # Get prices (use USD if available)
            if "price_usd" in ob1["bid"][1]:
                bid1 = ob1["bid"][1]["price_usd"]
                ask1 = ob1["ask"][1]["price_usd"]
            else:
                bid1 = ob1["bid"][1]["price"]
                ask1 = ob1["ask"][1]["price"]
                
            if "price_usd" in ob2["bid"][1]:
                bid2 = ob2["bid"][1]["price_usd"]
                ask2 = ob2["ask"][1]["price_usd"]
            else:
                bid2 = ob2["bid"][1]["price"]
                ask2 = ob2["ask"][1]["price"]
            
            mid1 = (bid1 + ask1) / 2
            mid2 = (bid2 + ask2) / 2
            spread1 = ask1 - bid1
            spread2 = ask2 - bid2
            
            print(f"\n{exchange1} - Mid: ${mid1:.2f}, Spread: ${spread1:.2f}")
            print(f"{exchange2} - Mid: ${mid2:.2f}, Spread: ${spread2:.2f}")
            print(f"Price Difference: ${abs(mid1 - mid2):.2f} ({abs(mid1 - mid2)/mid1*100:.2f}%)")


def main():
    parser = argparse.ArgumentParser(description='Inspect orderbook data from VictoriaMetrics')
    parser.add_argument('--vm-url', default='http://localhost:8428', help='VictoriaMetrics URL')
    parser.add_argument('--instrument', help='Instrument to inspect')
    parser.add_argument('--exchange', help='Exchange (derive/deribit)')
    parser.add_argument('--compare', help='Compare with another instrument (format: instrument:exchange)')
    parser.add_argument('--history', type=int, help='Show price history for N minutes')
    parser.add_argument('--all', action='store_true', help='Show all available orderbooks')
    
    args = parser.parse_args()
    
    inspector = OrderBookInspector(args.vm_url)
    
    if args.compare:
        # Compare two orderbooks
        parts = args.compare.split(':')
        if len(parts) != 2:
            print("Compare format should be: instrument:exchange")
            return
        
        instrument2, exchange2 = parts
        if not args.exchange:
            print("Please specify --exchange for the first instrument")
            return
            
        inspector.compare_orderbooks(args.instrument, args.exchange, instrument2, exchange2)
    
    elif args.history:
        # Show historical data
        if not args.exchange:
            print("Please specify --exchange")
            return
            
        df = inspector.get_orderbook_history(args.instrument, args.exchange, args.history)
        if not df.empty:
            print(f"\nPrice History for {args.instrument} on {args.exchange} (last {args.history} minutes)")
            print("="*60)
            
            # Show summary statistics
            print(f"Mid Price - Mean: ${df['mid_price'].mean():.4f}, Std: ${df['mid_price'].std():.4f}")
            print(f"Mid Price - Min: ${df['mid_price'].min():.4f}, Max: ${df['mid_price'].max():.4f}")
            
            if 'spread' in df.columns:
                print(f"Spread - Mean: ${df['spread'].mean():.4f}, Std: ${df['spread'].std():.4f}")
            
            # Show recent values
            print(f"\nRecent values:")
            print(df.tail(10).to_string(index=False))
    
    elif args.all:
        # Show all available instruments
        query = 'group by (instrument, exchange) (orderbook_mid_price)'
        results = inspector.query_current(query)
        
        print("\nAvailable Order Books:")
        print("="*40)
        for result in results:
            instrument = result["metric"]["instrument"]
            exchange = result["metric"]["exchange"]
            print(f"  {instrument} on {exchange}")
    
    else:
        # Show single orderbook
        if not args.exchange:
            print("Please specify --exchange")
            return
            
        orderbook = inspector.get_current_orderbook(args.instrument, args.exchange)
        inspector.display_orderbook(orderbook, args.instrument, args.exchange)


if __name__ == "__main__":
    main()