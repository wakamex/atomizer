#!/usr/bin/env python3
"""
Check presence of specific order sizes in the orderbook over time.
Useful for identifying market maker patterns and relationships.
"""

import argparse
import requests
import pandas as pd
import numpy as np
from datetime import datetime, timedelta
from collections import defaultdict
import matplotlib.pyplot as plt
import seaborn as sns
from scipy.stats import pearsonr


class PresenceChecker:
    def __init__(self, vm_url="http://localhost:8428"):
        self.vm_url = vm_url
    
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
    
    def check_size_presence(self, instrument, exchange, sizes, start_time, end_time, step="5s"):
        """Check presence of specific sizes in orderbook over time"""
        # Initialize presence tracking
        presence_data = defaultdict(lambda: defaultdict(bool))
        all_timestamps = set()
        
        # Query each size
        for size in sizes:
            # Try both exact match and close match (within 0.1)
            query = f'orderbook_size{{instrument="{instrument}",exchange="{exchange}"}} == {size} or (orderbook_size{{instrument="{instrument}",exchange="{exchange}"}} > {size - 0.1} and orderbook_size{{instrument="{instrument}",exchange="{exchange}"}} < {size + 0.1})'
            
            results = self.query_range(query, start_time, end_time, step)
            
            # Process results
            for series in results:
                for timestamp, value in series["values"]:
                    ts = float(timestamp)
                    all_timestamps.add(ts)
                    if float(value) > 0:  # Size is present
                        presence_data[size][ts] = True
        
        # Convert to DataFrame
        timestamps = sorted(all_timestamps)
        data = []
        
        for ts in timestamps:
            row = {"timestamp": datetime.fromtimestamp(ts)}
            for size in sizes:
                row[f"size_{size}"] = presence_data[size].get(ts, False)
            data.append(row)
        
        return pd.DataFrame(data)
    
    def analyze_presence(self, df, sizes):
        """Analyze presence patterns"""
        results = {}
        
        # Calculate presence percentage for each size
        for size in sizes:
            col = f"size_{size}"
            if col in df.columns:
                presence_pct = (df[col].sum() / len(df)) * 100
                results[size] = {
                    "presence_pct": presence_pct,
                    "present_count": df[col].sum(),
                    "total_count": len(df)
                }
        
        # Calculate correlation between presence patterns
        correlations = {}
        for i, size1 in enumerate(sizes):
            for j, size2 in enumerate(sizes):
                if i < j:  # Only calculate once for each pair
                    col1 = f"size_{size1}"
                    col2 = f"size_{size2}"
                    if col1 in df.columns and col2 in df.columns:
                        # Convert boolean to int for correlation
                        corr, p_value = pearsonr(df[col1].astype(int), df[col2].astype(int))
                        correlations[f"{size1}_vs_{size2}"] = {
                            "correlation": corr,
                            "p_value": p_value
                        }
        
        # Find co-occurrence patterns
        co_occurrences = {}
        for i, size1 in enumerate(sizes):
            for j, size2 in enumerate(sizes):
                if i < j:
                    col1 = f"size_{size1}"
                    col2 = f"size_{size2}"
                    if col1 in df.columns and col2 in df.columns:
                        both_present = ((df[col1]) & (df[col2])).sum()
                        either_present = ((df[col1]) | (df[col2])).sum()
                        co_occurrence_rate = both_present / either_present if either_present > 0 else 0
                        co_occurrences[f"{size1}_with_{size2}"] = {
                            "both_present": both_present,
                            "either_present": either_present,
                            "co_occurrence_rate": co_occurrence_rate
                        }
        
        return results, correlations, co_occurrences
    
    def plot_presence_timeline(self, df, sizes, instrument, exchange):
        """Plot presence timeline for each size"""
        fig, axes = plt.subplots(len(sizes), 1, figsize=(12, 2 * len(sizes)), sharex=True)
        if len(sizes) == 1:
            axes = [axes]
        
        for i, size in enumerate(sizes):
            col = f"size_{size}"
            if col in df.columns:
                # Create presence plot
                ax = axes[i]
                presence = df[col].astype(int)
                ax.fill_between(df.index, 0, presence, alpha=0.7, label=f"Size {size}")
                ax.set_ylabel(f"{size}")
                ax.set_ylim(-0.1, 1.1)
                ax.set_yticks([0, 1])
                ax.set_yticklabels(["Absent", "Present"])
                ax.grid(True, alpha=0.3)
        
        axes[-1].set_xlabel("Time")
        fig.suptitle(f"Order Size Presence Timeline\n{instrument} on {exchange}")
        plt.tight_layout()
        return fig
    
    def plot_correlation_heatmap(self, correlations, sizes):
        """Plot correlation heatmap between size presences"""
        # Create correlation matrix
        n = len(sizes)
        corr_matrix = np.eye(n)
        
        for i, size1 in enumerate(sizes):
            for j, size2 in enumerate(sizes):
                if i != j:
                    key = f"{min(size1, size2)}_vs_{max(size1, size2)}"
                    if key in correlations:
                        corr_matrix[i, j] = correlations[key]["correlation"]
        
        # Plot heatmap
        fig, ax = plt.subplots(figsize=(8, 6))
        sns.heatmap(corr_matrix, annot=True, fmt=".3f", cmap="coolwarm", center=0,
                    xticklabels=sizes, yticklabels=sizes, ax=ax,
                    cbar_kws={"label": "Presence Correlation"})
        ax.set_title("Order Size Presence Correlation Matrix")
        plt.tight_layout()
        return fig


def main():
    parser = argparse.ArgumentParser(description='Check presence of specific order sizes over time')
    parser.add_argument('--vm-url', default='http://localhost:8428', help='VictoriaMetrics URL')
    parser.add_argument('--instrument', required=True, help='Instrument to analyze')
    parser.add_argument('--exchange', required=True, help='Exchange (derive/deribit)')
    parser.add_argument('--sizes', required=True, help='Comma-separated list of sizes to track')
    parser.add_argument('--start', default='1h', help='Start time (e.g., 30m, 1h, 24h)')
    parser.add_argument('--step', default='5s', help='Time step for analysis')
    parser.add_argument('--plot', action='store_true', help='Generate visualization plots')
    parser.add_argument('--save-plots', help='Save plots to directory')
    
    args = parser.parse_args()
    
    # Parse sizes
    sizes = [float(s.strip()) for s in args.sizes.split(',')]
    
    # Parse time range
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
    
    # Create checker
    checker = PresenceChecker(args.vm_url)
    
    print(f"Analyzing order size presence for {args.instrument} on {args.exchange}")
    print(f"Tracking sizes: {sizes}")
    print(f"Time range: {start_time} to {end_time}")
    print(f"Step: {args.step}")
    print("-" * 80)
    
    # Check presence
    df = checker.check_size_presence(args.instrument, args.exchange, sizes, start_time, end_time, args.step)
    
    if df.empty:
        print("No data found")
        return
    
    # Set timestamp as index for plotting
    df.set_index('timestamp', inplace=True)
    
    # Analyze presence
    results, correlations, co_occurrences = checker.analyze_presence(df, sizes)
    
    # Display results
    print("\nPresence Analysis:")
    print("=" * 80)
    
    # Sort by presence percentage
    sorted_sizes = sorted(results.items(), key=lambda x: x[1]["presence_pct"], reverse=True)
    
    for size, stats in sorted_sizes:
        print(f"\nSize {size}:")
        print(f"  Presence: {stats['presence_pct']:.1f}% ({stats['present_count']}/{stats['total_count']} observations)")
    
    # Group sizes by similar presence
    print("\n" + "=" * 80)
    print("Presence Grouping (likely same market maker if similar):")
    print("=" * 80)
    
    # Group by presence percentage (within 5% considered similar)
    groups = defaultdict(list)
    for size, stats in results.items():
        pct = stats["presence_pct"]
        group_key = round(pct / 5) * 5  # Round to nearest 5%
        groups[group_key].append((size, pct))
    
    for group_pct in sorted(groups.keys(), reverse=True):
        group_sizes = groups[group_pct]
        if len(group_sizes) > 1:
            print(f"\n~{group_pct}% presence: {', '.join([f'{s[0]} ({s[1]:.1f}%)' for s in group_sizes])}")
    
    # Display correlations
    print("\n" + "=" * 80)
    print("Presence Correlations (high correlation = likely same market maker):")
    print("=" * 80)
    
    sorted_corrs = sorted(correlations.items(), key=lambda x: abs(x[1]["correlation"]), reverse=True)
    for pair, stats in sorted_corrs:
        if stats["p_value"] < 0.05:  # Only show significant correlations
            print(f"\n{pair}: {stats['correlation']:.3f} (p={stats['p_value']:.3e})")
            if abs(stats["correlation"]) > 0.8:
                print("  -> STRONG correlation - likely same market maker")
            elif abs(stats["correlation"]) > 0.5:
                print("  -> Moderate correlation - possibly related")
            else:
                print("  -> Weak correlation - likely different market makers")
    
    # Display co-occurrences
    print("\n" + "=" * 80)
    print("Co-occurrence Analysis:")
    print("=" * 80)
    
    sorted_cooc = sorted(co_occurrences.items(), key=lambda x: x[1]["co_occurrence_rate"], reverse=True)
    for pair, stats in sorted_cooc:
        if stats["either_present"] > 10:  # Only show if enough data
            print(f"\n{pair}:")
            print(f"  Co-occurrence rate: {stats['co_occurrence_rate']:.1%}")
            print(f"  Both present: {stats['both_present']} times")
            print(f"  Either present: {stats['either_present']} times")
    
    # Market maker hypothesis
    print("\n" + "=" * 80)
    print("Market Maker Hypothesis:")
    print("=" * 80)
    
    # Find highly correlated groups
    mm_groups = []
    used_sizes = set()
    
    for pair, stats in correlations.items():
        if stats["correlation"] > 0.8 and stats["p_value"] < 0.05:
            sizes_in_pair = [float(s) for s in pair.replace("_vs_", "_").split("_")]
            
            # Find or create group
            group_found = False
            for group in mm_groups:
                if any(s in group for s in sizes_in_pair):
                    group.update(sizes_in_pair)
                    group_found = True
                    break
            
            if not group_found:
                mm_groups.append(set(sizes_in_pair))
    
    # Add uncorrelated sizes as individual groups
    for size in sizes:
        if not any(size in group for group in mm_groups):
            mm_groups.append({size})
    
    print(f"\nDetected {len(mm_groups)} likely market maker(s):")
    for i, group in enumerate(mm_groups, 1):
        sorted_group = sorted(group)
        avg_presence = np.mean([results[s]["presence_pct"] for s in sorted_group])
        print(f"\nMarket Maker {i}:")
        print(f"  Sizes: {sorted_group}")
        print(f"  Average presence: {avg_presence:.1f}%")
    
    # Generate plots if requested
    if args.plot or args.save_plots:
        # Timeline plot
        fig1 = checker.plot_presence_timeline(df, sizes, args.instrument, args.exchange)
        if args.save_plots:
            fig1.savefig(f"{args.save_plots}/presence_timeline.png", dpi=300, bbox_inches='tight')
        elif args.plot:
            # Auto-save as plot.png when using --plot
            fig1.savefig("plot.png", dpi=300, bbox_inches='tight')
            print(f"\nPlot saved as plot.png")
        
        # Correlation heatmap
        if len(sizes) > 1:
            fig2 = checker.plot_correlation_heatmap(correlations, sizes)
            if args.save_plots:
                fig2.savefig(f"{args.save_plots}/presence_correlation.png", dpi=300, bbox_inches='tight')
            elif args.plot:
                # Auto-save correlation heatmap
                fig2.savefig("plot_correlation.png", dpi=300, bbox_inches='tight')
                print(f"Correlation heatmap saved as plot_correlation.png")
        
        if args.plot:
            plt.show()


if __name__ == "__main__":
    main()