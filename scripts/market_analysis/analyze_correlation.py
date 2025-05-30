#!/usr/bin/env python3
"""
Analyze correlation between Derive and Deribit order books using VictoriaMetrics data.
"""

import argparse
import requests
import pandas as pd
import numpy as np
from datetime import datetime, timedelta
import matplotlib.pyplot as plt
import seaborn as sns
from scipy import stats
import re
import warnings
warnings.filterwarnings('ignore')


class InstrumentConverter:
    """Convert between Derive and Deribit instrument naming conventions."""
    
    def __init__(self):
        # Regex patterns for different formats
        self.derive_pattern = re.compile(r'^([A-Z]+)-(\d{8})-(\d+)-([CP])$')  # ETH-20250531-2700-C
        self.deribit_pattern = re.compile(r'^([A-Z]+)-(\d{1,2})([A-Z]{3})(\d{2})-(\d+)-([CP])$')  # ETH-31MAY25-2700-C
        
        # Month mapping
        self.months = {
            'JAN': 1, 'FEB': 2, 'MAR': 3, 'APR': 4, 'MAY': 5, 'JUN': 6,
            'JUL': 7, 'AUG': 8, 'SEP': 9, 'OCT': 10, 'NOV': 11, 'DEC': 12
        }
        self.month_names = {v: k for k, v in self.months.items()}
    
    def derive_to_deribit(self, instrument):
        """Convert Derive format (ETH-20250531-2700-C) to Deribit format (ETH-31MAY25-2700-C)."""
        match = self.derive_pattern.match(instrument)
        if not match:
            return None
            
        asset, date_str, strike, option_type = match.groups()
        
        # Parse YYYYMMDD
        try:
            date = datetime.strptime(date_str, '%Y%m%d')
        except ValueError:
            return None
        
        # Format as DDMMMYY
        day = date.day
        month = self.month_names[date.month]
        year = date.strftime('%y')
        
        return f"{asset}-{day}{month}{year}-{strike}-{option_type}"
    
    def deribit_to_derive(self, instrument):
        """Convert Deribit format (ETH-31MAY25-2700-C) to Derive format (ETH-20250531-2700-C)."""
        match = self.deribit_pattern.match(instrument)
        if not match:
            return None
            
        asset, day, month, year, strike, option_type = match.groups()
        
        # Convert 2-digit year to 4-digit
        year_int = int(year)
        full_year = 2000 + year_int
        
        # Get month number
        month_num = self.months.get(month)
        if not month_num:
            return None
        
        # Create date
        try:
            date = datetime(full_year, month_num, int(day))
        except ValueError:
            return None
        
        # Format as YYYYMMDD
        date_str = date.strftime('%Y%m%d')
        
        return f"{asset}-{date_str}-{strike}-{option_type}"
    
    def convert_for_exchange(self, instrument, target_exchange):
        """Convert instrument to the format used by the target exchange."""
        # Try to detect format and convert
        if self.derive_pattern.match(instrument):
            if target_exchange == 'derive':
                return instrument
            else:
                return self.derive_to_deribit(instrument)
        elif self.deribit_pattern.match(instrument):
            if target_exchange == 'deribit':
                return instrument
            else:
                return self.deribit_to_derive(instrument)
        
        # Return original if no pattern matches
        return instrument
    
    def get_canonical_name(self, instrument):
        """Get a canonical name for matching (using Derive format)."""
        if self.deribit_pattern.match(instrument):
            return self.deribit_to_derive(instrument)
        return instrument


class OrderBookAnalyzer:
    def __init__(self, vm_url="http://localhost:8428"):
        self.vm_url = vm_url
        self.converter = InstrumentConverter()
        
    def query_metric(self, query, start_time, end_time, step="10s"):
        """Query VictoriaMetrics for time series data."""
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
    
    def parse_time_series(self, result):
        """Parse VictoriaMetrics result into pandas DataFrame."""
        all_data = []
        
        for series in result:
            metric = series["metric"]
            exchange = metric.get("exchange", "unknown")
            instrument = metric.get("instrument", "unknown")
            
            for timestamp, value in series["values"]:
                all_data.append({
                    "timestamp": datetime.fromtimestamp(float(timestamp)),
                    "exchange": exchange,
                    "instrument": instrument,
                    "value": float(value)
                })
        
        return pd.DataFrame(all_data)
    
    def get_eth_spot_prices(self, start_time, end_time, step="10s"):
        """Get ETH spot prices for converting Deribit ETH prices to USD."""
        # Query ETH spot price - try different possible instrument names
        # Including perpetual futures which are close to spot
        spot_instruments = [
            "ETH-SPOT",  # Our new spot price from Derive
            "ETH-USD", "ETH-USDT", "ETH", "ETHUSD",
            "ETH-PERPETUAL", "ETH-PERP", "ETH_USDT", "ETH_USD",
            "ETH-FS-30MAY25_PERP", "ETH-FS-6JUN25_PERP"  # Deribit perpetuals
        ]
        
        for instrument in spot_instruments:
            query = f'market_bid_price{{instrument="{instrument}"}}'
            result = self.query_metric(query, start_time, end_time, step)
            
            if result:
                df = self.parse_time_series(result)
                if not df.empty:
                    print(f"Found ETH spot/perpetual prices using instrument: {instrument} from {df['exchange'].iloc[0]}")
                    return df[['timestamp', 'value']].rename(columns={'value': 'eth_spot'})
        
        # Try to find any ETH perpetual from Deribit - use regex to match any PERP instrument
        query = 'market_bid_price{instrument=~"ETH.*PERP.*"}'
        result = self.query_metric(query, start_time, end_time, step)
        if result:
            df = self.parse_time_series(result)
            if not df.empty:
                # Use the first available perpetual
                instrument = df['instrument'].iloc[0]
                exchange = df['exchange'].iloc[0]
                print(f"Found ETH perpetual prices using instrument: {instrument} from {exchange}")
                return df[['timestamp', 'value']].rename(columns={'value': 'eth_spot'})
        
        print("Warning: Could not find ETH spot prices")
        print("Available instruments might not include ETH spot/perpetual data")
        print("Make sure the market monitor is running with spot collection enabled")
        return pd.DataFrame()
    
    def get_paired_data(self, metric="market_bid_price", instrument_pattern="ETH-", 
                       start_time=None, end_time=None, step="10s", convert_deribit_to_usd=True):
        """Get paired data for Derive and Deribit."""
        if end_time is None:
            end_time = datetime.now()
        if start_time is None:
            start_time = end_time - timedelta(hours=1)
            
        # Query options data
        query = f'{metric}{{instrument=~"{instrument_pattern}.*"}}'
        result = self.query_metric(query, start_time, end_time, step)
        
        # Parse into DataFrame
        df = self.parse_time_series(result)
        
        if df.empty:
            return pd.DataFrame()
        
        # Get ETH spot prices if we need to convert Deribit prices
        eth_spot_df = pd.DataFrame()
        if convert_deribit_to_usd and instrument_pattern.startswith("ETH"):
            eth_spot_df = self.get_eth_spot_prices(start_time, end_time, step)
        
        # Add canonical instrument names for matching
        df['canonical_instrument'] = df['instrument'].apply(self.converter.get_canonical_name)
        
        # Create separate dataframes for each exchange
        derive_df = df[df['exchange'] == 'derive'].copy()
        deribit_df = df[df['exchange'] == 'deribit'].copy()
        
        # Merge on canonical instrument name and timestamp
        if not derive_df.empty and not deribit_df.empty:
            merged_df = pd.merge(
                derive_df[['timestamp', 'canonical_instrument', 'value', 'instrument']],
                deribit_df[['timestamp', 'canonical_instrument', 'value', 'instrument']],
                on=['timestamp', 'canonical_instrument'],
                suffixes=('_derive', '_deribit'),
                how='inner'
            )
            
            # If we have ETH spot prices, merge them in
            if not eth_spot_df.empty and convert_deribit_to_usd:
                merged_df = pd.merge(
                    merged_df,
                    eth_spot_df,
                    on='timestamp',
                    how='left'
                )
                
                # Forward fill ETH spot prices for any missing values
                merged_df['eth_spot'] = merged_df['eth_spot'].ffill()
                
                # Convert Deribit prices from ETH to USD
                merged_df['value_deribit_usd'] = merged_df['value_deribit'] * merged_df['eth_spot']
                merged_df['value_deribit_original'] = merged_df['value_deribit']
                merged_df['value_deribit'] = merged_df['value_deribit_usd']
            
            # Rename columns
            merged_df.rename(columns={
                'value_derive': 'derive',
                'value_deribit': 'deribit',
                'canonical_instrument': 'instrument'
            }, inplace=True)
            
            # Add original instrument names for reference
            merged_df['derive_instrument'] = merged_df['instrument_derive']
            merged_df['deribit_instrument'] = merged_df['instrument_deribit']
            
            # Select final columns
            columns = ['timestamp', 'instrument', 'derive', 'deribit', 
                      'derive_instrument', 'deribit_instrument']
            if 'eth_spot' in merged_df.columns:
                columns.extend(['eth_spot', 'value_deribit_original'])
            
            return merged_df[columns]
        
        return pd.DataFrame()
    
    def analyze_correlation(self, df, instrument=None):
        """Analyze correlation between Derive and Deribit prices."""
        if instrument:
            df = df[df['instrument'] == instrument]
        
        if df.empty:
            return None
            
        results = []
        
        # Group by instrument
        for inst, group in df.groupby('instrument'):
            if len(group) < 10:  # Need enough data points
                continue
                
            derive_prices = group['derive'].values
            deribit_prices = group['deribit'].values
            
            # Calculate correlation
            correlation, p_value = stats.pearsonr(derive_prices, deribit_prices)
            
            # Calculate price differences
            price_diff = derive_prices - deribit_prices
            mean_diff = np.mean(price_diff)
            std_diff = np.std(price_diff)
            
            # Calculate percentage differences
            pct_diff = (price_diff / deribit_prices) * 100
            mean_pct_diff = np.mean(pct_diff)
            
            # Find lead/lag using cross-correlation
            max_lag = min(20, len(derive_prices) // 4)
            correlations = []
            lag_p_values = []
            lags = range(-max_lag, max_lag + 1)
            
            for lag in lags:
                if lag < 0:
                    # Derive leads
                    x = derive_prices[:lag]
                    y = deribit_prices[-lag:]
                elif lag > 0:
                    # Deribit leads
                    x = derive_prices[lag:]
                    y = deribit_prices[:-lag]
                else:
                    x = derive_prices
                    y = deribit_prices
                
                if len(x) > 2:  # Need at least 3 points for correlation
                    corr, p_val = stats.pearsonr(x, y)
                    correlations.append(corr)
                    lag_p_values.append(p_val)
                else:
                    correlations.append(0)
                    lag_p_values.append(1.0)
            
            # Find best lag
            best_lag_idx = np.argmax(np.abs(correlations))
            best_lag = lags[best_lag_idx]
            best_lag_corr = correlations[best_lag_idx]
            best_lag_p_value = lag_p_values[best_lag_idx]
            
            results.append({
                'instrument': inst,
                'n_points': len(group),
                'correlation': correlation,
                'p_value': p_value,
                'mean_diff': mean_diff,
                'std_diff': std_diff,
                'mean_pct_diff': mean_pct_diff,
                'best_lag': best_lag,
                'best_lag_corr': best_lag_corr,
                'best_lag_p_value': best_lag_p_value,
                'derive_mean': np.mean(derive_prices),
                'deribit_mean': np.mean(deribit_prices),
                'derive_std': np.std(derive_prices),
                'deribit_std': np.std(deribit_prices)
            })
        
        return pd.DataFrame(results)
    
    def plot_analysis(self, df, instrument=None, save_path=None):
        """Create visualization of the analysis."""
        if instrument:
            df = df[df['instrument'] == instrument]
            
        if df.empty:
            print("No data to plot")
            return
            
        # Create figure with subplots
        fig, axes = plt.subplots(2, 2, figsize=(15, 10))
        fig.suptitle(f'Derive vs Deribit Analysis{" - " + instrument if instrument else ""}')
        
        # 1. Time series comparison
        ax = axes[0, 0]
        for inst, group in df.groupby('instrument'):
            if instrument and inst != instrument:
                continue
            group_sorted = group.sort_values('timestamp')
            ax.plot(group_sorted['timestamp'], group_sorted['derive'], 
                   label=f'{inst} (Derive)', alpha=0.7)
            ax.plot(group_sorted['timestamp'], group_sorted['deribit'], 
                   label=f'{inst} (Deribit)', alpha=0.7, linestyle='--')
        ax.set_xlabel('Time')
        ax.set_ylabel('Price')
        ax.set_title('Price Time Series')
        ax.legend()
        ax.grid(True, alpha=0.3)
        
        # 2. Scatter plot
        ax = axes[0, 1]
        for inst, group in df.groupby('instrument'):
            if instrument and inst != instrument:
                continue
            ax.scatter(group['deribit'], group['derive'], alpha=0.5, label=inst)
        
        # Add diagonal line
        min_val = min(df['derive'].min(), df['deribit'].min())
        max_val = max(df['derive'].max(), df['deribit'].max())
        ax.plot([min_val, max_val], [min_val, max_val], 'k--', alpha=0.5)
        
        ax.set_xlabel('Deribit Price')
        ax.set_ylabel('Derive Price')
        ax.set_title('Price Correlation')
        ax.legend()
        ax.grid(True, alpha=0.3)
        
        # 3. Price difference histogram
        ax = axes[1, 0]
        for inst, group in df.groupby('instrument'):
            if instrument and inst != instrument:
                continue
            diff = group['derive'] - group['deribit']
            ax.hist(diff, bins=30, alpha=0.5, label=inst)
        ax.set_xlabel('Price Difference (Derive - Deribit)')
        ax.set_ylabel('Frequency')
        ax.set_title('Price Difference Distribution')
        ax.legend()
        ax.grid(True, alpha=0.3)
        
        # 4. Rolling correlation
        ax = axes[1, 1]
        for inst, group in df.groupby('instrument'):
            if instrument and inst != instrument:
                continue
            group_sorted = group.sort_values('timestamp')
            if len(group_sorted) > 20:
                rolling_corr = group_sorted['derive'].rolling(20).corr(group_sorted['deribit'])
                ax.plot(group_sorted['timestamp'], rolling_corr, label=inst, alpha=0.7)
        ax.set_xlabel('Time')
        ax.set_ylabel('Correlation')
        ax.set_title('Rolling Correlation (20 points)')
        ax.legend()
        ax.grid(True, alpha=0.3)
        ax.set_ylim([-1, 1])
        
        plt.tight_layout()
        
        if save_path:
            plt.savefig(save_path, dpi=300, bbox_inches='tight')
            if save_path != "plot.png":  # Only print if not the default
                print(f"Plot saved to {save_path}")
        else:
            plt.show()

def main():
    parser = argparse.ArgumentParser(description='Analyze Derive vs Deribit order book correlation')
    parser.add_argument('--vm-url', default='http://localhost:8428', help='VictoriaMetrics URL')
    parser.add_argument('--metric', default='market_bid_price', 
                       choices=['market_bid_price', 'market_ask_price', 'market_spread', 
                               'market_spread_percent', 'market_bid_size', 'market_ask_size',
                               'orderbook_mid_price', 'orderbook_price', 'orderbook_spread',
                               'orderbook_spread_percent', 'orderbook_total_bid_size', 
                               'orderbook_total_ask_size'],
                       help='Metric to analyze')
    parser.add_argument('--instrument', default='ETH-', help='Instrument pattern (e.g., ETH-, ETH-20250601-)')
    parser.add_argument('--start', default='1h', help='Start time (e.g., 30m, 1h, 24h, 2023-01-01T00:00:00)')
    parser.add_argument('--step', default='10s', help='Step interval')
    parser.add_argument('--plot', action='store_true', help='Generate plots')
    parser.add_argument('--save-plot', help='Save plot to file')
    parser.add_argument('--specific-instrument', help='Analyze specific instrument only')
    parser.add_argument('--no-convert-usd', action='store_true', help='Do not convert Deribit ETH prices to USD')
    parser.add_argument('--compare-returns', action='store_true', help='Compare percentage returns instead of absolute prices')
    
    args = parser.parse_args()
    
    # Parse start time
    end_time = datetime.now()
    if args.start.endswith('h'):
        hours = int(args.start[:-1])
        start_time = end_time - timedelta(hours=hours)
    elif args.start.endswith('m'):
        minutes = int(args.start[:-1])
        start_time = end_time - timedelta(minutes=minutes)
    elif args.start.endswith('d'):
        days = int(args.start[:-1])
        start_time = end_time - timedelta(days=days)
    else:
        start_time = datetime.fromisoformat(args.start)
    
    # Create analyzer
    analyzer = OrderBookAnalyzer(args.vm_url)
    
    print(f"Analyzing {args.metric} for instruments matching '{args.instrument}'")
    print(f"Time range: {start_time} to {end_time}")
    print(f"Step: {args.step}")
    print("-" * 80)
    
    # Get data
    df = analyzer.get_paired_data(
        metric=args.metric,
        instrument_pattern=args.instrument,
        start_time=start_time,
        end_time=end_time,
        step=args.step,
        convert_deribit_to_usd=not args.no_convert_usd
    )
    
    if df.empty:
        print("No paired data found for Derive and Deribit")
        return
    
    print(f"Found {len(df)} paired data points across {df['instrument'].nunique()} instruments")
    
    # Show sample of matched instruments and ETH spot price
    if not df.empty:
        print("\nSample instrument mappings:")
        sample_instruments = df[['derive_instrument', 'deribit_instrument']].drop_duplicates().head(5)
        for _, row in sample_instruments.iterrows():
            print(f"  Derive: {row['derive_instrument']} <-> Deribit: {row['deribit_instrument']}")
        
        if 'eth_spot' in df.columns:
            print(f"\nETH Spot Price Range: ${df['eth_spot'].min():.2f} - ${df['eth_spot'].max():.2f}")
            print(f"Average ETH Spot: ${df['eth_spot'].mean():.2f}")
    print()
    
    # If comparing returns, calculate percentage changes
    if args.compare_returns and not df.empty:
        print("\nCalculating percentage returns for correlation analysis...")
        # Group by instrument and calculate returns
        df_returns = []
        for instrument, group in df.groupby('instrument'):
            group = group.sort_values('timestamp')
            if len(group) > 1:
                group['derive_return'] = group['derive'].pct_change() * 100
                group['deribit_return'] = group['deribit'].pct_change() * 100
                # Drop the first row (NaN from pct_change)
                group = group.dropna()
                df_returns.append(group)
        
        if df_returns:
            df_returns = pd.concat(df_returns)
            # Replace price columns with returns for analysis
            df_returns['derive'] = df_returns['derive_return']
            df_returns['deribit'] = df_returns['deribit_return']
            df = df_returns
            print(f"Using {len(df)} return data points for analysis")
    
    # Analyze correlation
    results = analyzer.analyze_correlation(df, args.specific_instrument)
    
    if results is not None and not results.empty:
        # Sort by correlation
        results = results.sort_values('correlation', ascending=False)
        
        analysis_type = "Return" if args.compare_returns else "Price"
        print(f"{analysis_type} Correlation Analysis Results:")
        print("=" * 80)
        
        for _, row in results.iterrows():
            print(f"\nInstrument: {row['instrument']}")
            print(f"  Data points: {row['n_points']}")
            print(f"  Correlation: {row['correlation']:.4f} (p-value: {row['p_value']:.4e})")
            if args.compare_returns:
                print(f"  Mean return difference: {row['mean_diff']:.4f}%")
                print(f"  Std deviation of return diff: {row['std_diff']:.4f}%")
            else:
                print(f"  Mean difference: {row['mean_diff']:.4f} ({row['mean_pct_diff']:.2f}%)")
                print(f"  Std deviation of diff: {row['std_diff']:.4f}")
            print(f"  Best lag: {row['best_lag']} steps (correlation: {row['best_lag_corr']:.4f}, p-value: {row['best_lag_p_value']:.4e})")
            
            # Determine if lag is statistically significant
            if row['best_lag_p_value'] < 0.05:
                if row['best_lag'] > 0:
                    print(f"    -> Deribit leads Derive by {row['best_lag']} steps (statistically significant)")
                elif row['best_lag'] < 0:
                    print(f"    -> Derive leads Deribit by {-row['best_lag']} steps (statistically significant)")
                else:
                    print(f"    -> No lag detected (synchronous movement)")
            else:
                print(f"    -> Lead/lag not statistically significant (p={row['best_lag_p_value']:.3f})")
                
            if args.compare_returns:
                print(f"  Derive returns: mean={row['derive_mean']:.4f}%, std={row['derive_std']:.4f}%")
                print(f"  Deribit returns: mean={row['deribit_mean']:.4f}%, std={row['deribit_std']:.4f}%")
            else:
                print(f"  Derive: mean={row['derive_mean']:.4f}, std={row['derive_std']:.4f}")
                print(f"  Deribit: mean={row['deribit_mean']:.4f}, std={row['deribit_std']:.4f}")
        
        # Summary statistics
        print("\n" + "=" * 80)
        print("Summary Statistics:")
        print(f"  Average correlation: {results['correlation'].mean():.4f}")
        print(f"  Instruments with correlation > 0.9: {(results['correlation'] > 0.9).sum()}")
        print(f"  Instruments with correlation > 0.95: {(results['correlation'] > 0.95).sum()}")
        print(f"  Average absolute price difference: {results['mean_diff'].abs().mean():.4f}")
        print(f"  Average percentage difference: {results['mean_pct_diff'].abs().mean():.2f}%")
    
    # Generate plots if requested
    if args.plot or args.save_plot:
        save_path = args.save_plot if args.save_plot else "plot.png" if args.plot else None
        analyzer.plot_analysis(df, args.specific_instrument, save_path)
        if args.plot and not args.save_plot:
            print("\nPlot saved as plot.png")

if __name__ == "__main__":
    main()