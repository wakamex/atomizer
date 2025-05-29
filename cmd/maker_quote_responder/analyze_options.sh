#\!/bin/bash

# Options Analysis Tool
# This script analyzes Derive/Lyra options to help with manual trading decisions

echo "=== Derive Options Analysis Tool ==="
echo ""

# Default to running all analyses
COMMAND="${1:-all}"

# Extract exchange option if provided
EXCHANGE_FLAG=""
for arg in "$@"; do
    if [[ $arg == --exchange=* ]]; then
        EXCHANGE_FLAG="$arg"
    fi
done

case "$COMMAND" in
    "all")
        echo "Running complete options analysis..."
        go run ../analyze_options/ all $EXCHANGE_FLAG
        ;;
    "expiry")
        echo "Analyzing options by expiry date..."
        go run ../analyze_options/ expiry $EXCHANGE_FLAG
        ;;
    "nearterm")
        DAYS="${2:-30}"
        echo "Showing options expiring in next $DAYS days..."
        go run ../analyze_options/ nearterm $DAYS $EXCHANGE_FLAG
        ;;
    "export")
        DAYS="${2:-1}"
        echo "Exporting ETH call options expiring in $DAYS days..."
        go run ../analyze_options/ export $DAYS $EXCHANGE_FLAG
        ;;
    "stats")
        echo "Showing active options statistics and strike distribution..."
        go run ../analyze_options/ stats $EXCHANGE_FLAG
        ;;
    "active")
        echo "Showing active percentage by expiry..."
        go run ../analyze_options/ active $EXCHANGE_FLAG
        ;;
    "query")
        EXPIRY_INDEX="${2:-1}"
        echo "Querying ETH calls for expiry #$EXPIRY_INDEX..."
        go run ../analyze_options/ query $EXPIRY_INDEX $EXCHANGE_FLAG
        ;;
    "help"|*)
        echo "Usage: ./analyze_options.sh [command] [args] [options]"
        echo ""
        echo "Commands:"
        echo "  all                    - Run all analyses (default)"
        echo "  expiry                 - Analyze options by expiry date"
        echo "  nearterm [days]        - Show near-term options (default: 30 days)"
        echo "  export [days]          - Export ETH calls to CSV (default: 1 day)"
        echo "  stats                  - Show active options statistics + strike distribution"
        echo "  active                 - Show active percentage by expiry with ✓/✗ indicators"
        echo "  query [N]              - Query ETH calls for Nth expiry (default: 1=nearest, 2=second-nearest)"
        echo ""
        echo "Options:"
        echo "  --exchange=deribit     - Use Deribit instead of Derive/Lyra (default)"
        echo ""
        echo "Environment Variables:"
        echo "  DERIBIT_TEST_MODE=true - Use Deribit testnet"
        echo "  ANALYZE_BTC=true       - Analyze BTC options instead of ETH (Deribit only)"
        echo ""
        echo "Examples:"
        echo "  ./analyze_options.sh                        # Run all analyses on Derive"
        echo "  ./analyze_options.sh query --exchange=deribit  # Query Deribit options"
        echo "  ./analyze_options.sh nearterm 7             # Show options expiring in 7 days"
        echo "  ./analyze_options.sh export 2               # Export ETH calls expiring in 2 days"
        echo "  ./analyze_options.sh query                  # Query nearest expiry ETH calls"
        echo "  ./analyze_options.sh query 2                # Query second-nearest expiry ETH calls"
        ;;
esac

echo ""
echo "=== Analysis Complete ==="
