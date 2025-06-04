package arbitrage

import (
	"context"

	"github.com/wakamex/atomizer/internal/types"
	"github.com/shopspring/decimal"
)

// GammaModule interface for gamma calculations
type GammaModule interface {
	CalculateGamma(positions []types.Position) decimal.Decimal
	ShouldHedge(gamma decimal.Decimal) bool
}

// GammaHedger interface for gamma hedging
type GammaHedger interface {
	Start(ctx context.Context)
	Stop()
}