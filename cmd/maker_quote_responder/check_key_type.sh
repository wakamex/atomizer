#!/bin/bash

# Check what type of Deribit API key you have

echo "=== Deribit API Key Type Checker ==="
echo ""
echo "Check your API key format:"
echo ""
echo "1. Standard API Keys (HMAC-SHA256):"
echo "   - API Key: Short alphanumeric string (e.g., 'AbCdEfGh')"
echo "   - API Secret: Longer alphanumeric string (e.g., 'AbCdEfGhIjKlMnOpQrStUvWxYz123456')"
echo ""
echo "2. Asymmetric API Keys (Ed25519/RSA):"
echo "   - Client ID: Short identifier (e.g., 'ABC123')"
echo "   - Private Key: Long key starting with '-----BEGIN PRIVATE KEY-----'"
echo "   - You would have generated this key pair yourself"
echo ""

# Load .env if exists
if [ -f .env ]; then
    source .env
fi

if [ ! -z "$DERIBIT_API_KEY" ]; then
    echo "Found DERIBIT_API_KEY:"
    echo "  Length: ${#DERIBIT_API_KEY} characters"
    echo "  First 4 chars: ${DERIBIT_API_KEY:0:4}..."
    
    if [[ "$DERIBIT_API_KEY" == *"-----BEGIN"* ]]; then
        echo "  ✓ This looks like an asymmetric PRIVATE KEY"
    else
        echo "  ✓ This looks like a standard API KEY"
    fi
fi

if [ ! -z "$DERIBIT_API_SECRET" ]; then
    echo ""
    echo "Found DERIBIT_API_SECRET:"
    echo "  Length: ${#DERIBIT_API_SECRET} characters"
    
    if [[ "$DERIBIT_API_SECRET" == *"-----BEGIN"* ]]; then
        echo "  ✓ This looks like an asymmetric key"
    else
        echo "  ✓ This looks like a standard API SECRET"
    fi
fi

echo ""
echo "If you have asymmetric keys (Ed25519/RSA):"
echo "  - CCXT doesn't support this natively"
echo "  - You'll need to use Deribit's REST API directly"
echo "  - Or generate standard API keys in Deribit settings"
echo ""
echo "To generate standard API keys:"
echo "  1. Log into Deribit"
echo "  2. Go to Account → API Management"
echo "  3. Create new API key"
echo "  4. Choose 'Deribit-generated key' (not 'Self-generated key')"