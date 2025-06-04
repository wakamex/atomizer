package gamma

import (
	"github.com/wakamex/atomizer/internal/types"
	"github.com/shopspring/decimal"
)

// Module implements gamma calculations and hedging decisions
type Module struct {
	gammaThreshold decimal.Decimal
}

// NewModule creates a new gamma module
func NewModule(gammaThreshold float64) *Module {
	return &Module{
		gammaThreshold: decimal.NewFromFloat(gammaThreshold),
	}
}

// CalculateGamma calculates total gamma from positions
func (m *Module) CalculateGamma(positions []types.Position) decimal.Decimal {
	totalGamma := decimal.Zero
	
	for _, pos := range positions {
		if !pos.Quantity.IsZero() {
			totalGamma = totalGamma.Add(pos.Gamma.Mul(pos.Quantity))
		}
	}
	
	return totalGamma
}

// ShouldHedge determines if gamma hedging is needed
func (m *Module) ShouldHedge(gamma decimal.Decimal) bool {
	return gamma.Abs().GreaterThan(m.gammaThreshold)
}