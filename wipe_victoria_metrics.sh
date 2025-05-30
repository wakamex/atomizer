#!/bin/bash

echo "This script will wipe all VictoriaMetrics data and restart with a clean database."
echo ""

# Check for running VictoriaMetrics processes
PIDS=$(pgrep -f victoria-metrics-prod)
if [ ! -z "$PIDS" ]; then
    echo "Found VictoriaMetrics running with PIDs: $PIDS"
    for pid in $PIDS; do
        echo "  Process $pid running from: $(pwdx $pid 2>/dev/null | cut -d' ' -f2)"
    done
fi

# Get the official installation path
VM_INSTALL_PATH="$HOME/.atomizer/victoria-metrics"

echo ""
echo "This will remove ALL data from the following locations:"
echo "  - $VM_INSTALL_PATH/data (official market monitor installation)"
echo "  - ./data (local directory)"
echo "  - ./victoria-metrics/data (local subdirectory)"
echo "  - ./vm-data (alternative local directory)"
echo ""
echo "Press Ctrl+C to cancel, or Enter to continue..."
read

# Find and kill VictoriaMetrics process
echo "Stopping all VictoriaMetrics processes..."
pkill -f victoria-metrics-prod || echo "No VictoriaMetrics processes found"

# Wait a moment for it to shut down
sleep 3

# Remove the data directories
echo "Removing data directories..."

# Remove the official installation data
if [ -d "$VM_INSTALL_PATH/data" ]; then
    rm -rf "$VM_INSTALL_PATH/data"
    echo "  Removed $VM_INSTALL_PATH/data"
fi

# Remove any local data directories
[ -d "./data" ] && rm -rf ./data && echo "  Removed ./data"
[ -d "./victoria-metrics/data" ] && rm -rf ./victoria-metrics/data && echo "  Removed ./victoria-metrics/data"
[ -d "./vm-data" ] && rm -rf ./vm-data && echo "  Removed ./vm-data"

# Create fresh data directory in the official location
mkdir -p "$VM_INSTALL_PATH/data"
echo "  Created fresh $VM_INSTALL_PATH/data"

echo ""
echo "Data wiped successfully!"
echo ""
echo "To start collecting data again, use the market monitor:"
echo "  ./cmd/market_monitor/market_monitor start --orderbook --instruments \"ETH-20250601-2600-C\""
echo ""
echo "The market monitor will automatically:"
echo "  - Start VictoriaMetrics from $VM_INSTALL_PATH"
echo "  - Collect order book data"
echo "  - Collect ETH and BTC spot prices"