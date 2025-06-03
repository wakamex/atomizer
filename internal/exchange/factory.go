package exchange

import (
	"fmt"
	"os"

	"github.com/wakamex/atomizer/internal/exchange/derive"
	"github.com/wakamex/atomizer/internal/types"
)

// NewExchange creates a new exchange instance based on the configuration
func NewExchange(config *types.MarketMakerConfig) (types.MarketMakerExchange, error) {
	switch config.Exchange {
	case "derive":
		return newDeriveExchange(config)
	case "deribit":
		return newDeribitExchange(config)
	default:
		return nil, fmt.Errorf("unsupported exchange: %s", config.Exchange)
	}
}

func newDeriveExchange(config *types.MarketMakerConfig) (types.MarketMakerExchange, error) {
	// Get credentials from environment
	privateKey := os.Getenv("DERIVE_PRIVATE_KEY")
	if privateKey == "" {
		privateKey = os.Getenv("PRIVATE_KEY")
	}
	if privateKey == "" {
		return nil, fmt.Errorf("DERIVE_PRIVATE_KEY or PRIVATE_KEY environment variable not set")
	}

	walletAddress := os.Getenv("DERIVE_WALLET_ADDRESS")
	if walletAddress == "" {
		return nil, fmt.Errorf("DERIVE_WALLET_ADDRESS environment variable not set")
	}

	// Import the derive package
	deriveExchange, err := derive.NewDeriveMarketMakerExchange(privateKey, walletAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create Derive exchange: %w", err)
	}

	return deriveExchange, nil
}

func newDeribitExchange(config *types.MarketMakerConfig) (types.MarketMakerExchange, error) {
	// Get credentials from environment
	apiKey := os.Getenv("DERIBIT_API_KEY")
	apiSecret := os.Getenv("DERIBIT_API_SECRET")

	// Check for private key auth
	privateKey := os.Getenv("DERIBIT_PRIVATE_KEY")
	if privateKey == "" {
		privateKey = os.Getenv("DERIBIT_ED25519_PRIVATE_KEY")
	}

	if privateKey != "" {
		// Use private key authentication
		// TODO: Create Deribit exchange with ED25519 auth
		return nil, fmt.Errorf("deribit ED25519 auth not yet implemented")
	}

	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("DERIBIT_API_KEY and DERIBIT_API_SECRET environment variables not set")
	}

	// TODO: Create DeribitExchange with API key/secret
	return nil, fmt.Errorf("deribit exchange creation not yet implemented")
}
