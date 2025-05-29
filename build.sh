#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building Atomizer unified CLI...${NC}"

# Create bin directory if it doesn't exist
mkdir -p bin

# Get git commit hash for version info
GIT_HASH=$(git rev-parse --short HEAD 2>/dev/null || echo "dev")
export BUILD_HASH=$GIT_HASH

# Build main atomizer CLI
echo -e "${YELLOW}Building atomizer CLI...${NC}"
cd cmd/atomizer
go build -o ../../bin/atomizer
cd ../..
echo -e "${GREEN}✓ atomizer CLI${NC}"

# Build maker_quote_responder (for RFQ and market maker)
echo -e "${YELLOW}Building maker_quote_responder...${NC}"
cd cmd/maker_quote_responder
go build -o ../../bin/maker_quote_responder
cd ../..
echo -e "${GREEN}✓ maker_quote_responder${NC}"

# Build analyze_options
echo -e "${YELLOW}Building analyze_options...${NC}"
cd cmd/analyze_options
go build -o ../../bin/analyze_options
cd ../..
echo -e "${GREEN}✓ analyze_options${NC}"

# Build inventory
echo -e "${YELLOW}Building inventory...${NC}"
cd cmd/inventory
go build -o ../../bin/inventory
cd ../..
echo -e "${GREEN}✓ inventory${NC}"

# Build markets
echo -e "${YELLOW}Building markets...${NC}"
cd cmd/markets
go build -o ../../bin/markets
cd ../..
echo -e "${GREEN}✓ markets${NC}"

# Build send_quote
echo -e "${YELLOW}Building send_quote...${NC}"
cd cmd/send_quote
go build -o ../../bin/send_quote
cd ../..
echo -e "${GREEN}✓ send_quote${NC}"

echo -e "${GREEN}Build complete!${NC}"
echo ""
echo "To use the unified CLI, add the bin directory to your PATH:"
echo "  export PATH=\$PATH:$(pwd)/bin"
echo ""
echo "Or run directly:"
echo "  ./bin/atomizer help"