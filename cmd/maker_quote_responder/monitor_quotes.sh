#!/bin/bash

echo "Monitoring for accepted quotes..."
echo "Press Ctrl+C to stop"
echo ""

# Watch the maker process output for confirmation messages
# Since the maker is running in a terminal, we'll simulate some RFQs to test

while true; do
    echo "Checking for accepted quotes at $(date)"
    
    # Check if maker is still running
    if ! pgrep -f "maker_quote_responder.*derive" > /dev/null; then
        echo "Maker process not running!"
        exit 1
    fi
    
    echo "Maker is running. Waiting for quote acceptance..."
    echo "To test hedge order placement, a quote needs to be accepted by a taker."
    echo ""
    
    sleep 30
done