#\!/bin/bash

# Options Analysis Tool
# This script analyzes Derive/Lyra options to help with manual trading decisions

echo "=== Derive Options Analysis Tool ==="
echo ""

# Default to running all analyses
COMMAND="${1:-all}"

case "$COMMAND" in
    "all")
        echo "Running complete options analysis..."
        go run ../analyze_options/ all
        ;;
    "expiry")
        echo "Analyzing options by expiry date..."
        go run ../analyze_options/ expiry
        ;;
    "nearterm")
        DAYS="${2:-30}"
        echo "Showing options expiring in next $DAYS days..."
        go run ../analyze_options/ nearterm $DAYS
        ;;
    "export")
        DAYS="${2:-1}"
        echo "Exporting ETH call options expiring in $DAYS days..."
        go run ../analyze_options/ export $DAYS
        ;;
    "stats")
        echo "Showing active options statistics and strike distribution..."
        go run ../analyze_options/ stats
        ;;
    "active")
        echo "Showing active percentage by expiry..."
        go run ../analyze_options/ active
        ;;
    "query")
        EXPIRY_INDEX="${2:-1}"
        echo "Querying ETH calls for expiry #$EXPIRY_INDEX..."
        go run ../analyze_options/ query $EXPIRY_INDEX
        ;;
    "help"|*)
        echo "Usage: ./analyze_options.sh [command] [args]"
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
        echo "Examples:"
        echo "  ./analyze_options.sh                   # Run all analyses"
        echo "  ./analyze_options.sh nearterm 7        # Show options expiring in 7 days"
        echo "  ./analyze_options.sh export 2          # Export ETH calls expiring in 2 days"
        echo "  ./analyze_options.sh query             # Query nearest expiry ETH calls"
        echo "  ./analyze_options.sh query 2           # Query second-nearest expiry ETH calls"
        ;;
esac

echo ""
echo "=== Analysis Complete ==="
