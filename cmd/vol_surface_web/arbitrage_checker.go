package main

import (
	"fmt"
	"math"
)

// ArbitrageViolation represents a detected arbitrage opportunity
type ArbitrageViolation struct {
	Type        string  // "calendar", "butterfly", "negative_density"
	Strike      float64
	Expiry1     float64
	Expiry2     float64 // Only used for calendar arbitrage
	Description string
	Severity    float64 // How severe is the violation (0-1)
}

// ArbitrageChecker checks for various arbitrage conditions
type ArbitrageChecker struct {
	Surface      *SVISurface
	SpotPrice    float64
	RiskFreeRate float64
}

// NewArbitrageChecker creates a new arbitrage checker
func NewArbitrageChecker(surface *SVISurface) *ArbitrageChecker {
	return &ArbitrageChecker{
		Surface:      surface,
		SpotPrice:    surface.SpotPrice,
		RiskFreeRate: surface.RiskFreeRate,
	}
}

// CheckAllArbitrage performs all arbitrage checks
func (ac *ArbitrageChecker) CheckAllArbitrage() []ArbitrageViolation {
	violations := []ArbitrageViolation{}
	
	// Check calendar arbitrage
	violations = append(violations, ac.CheckCalendarArbitrage()...)
	
	// Check butterfly arbitrage
	violations = append(violations, ac.CheckButterflyArbitrage()...)
	
	// Check for negative densities
	violations = append(violations, ac.CheckNegativeDensity()...)
	
	return violations
}

// CheckCalendarArbitrage checks if total variance increases with time
func (ac *ArbitrageChecker) CheckCalendarArbitrage() []ArbitrageViolation {
	violations := []ArbitrageViolation{}
	
	// Check at various moneyness levels
	moneyness := []float64{0.5, 0.7, 0.9, 1.0, 1.1, 1.3, 1.5}
	
	for _, k := range moneyness {
		strike := ac.SpotPrice * k
		
		// Check each consecutive pair of expiries
		for i := 0; i < len(ac.Surface.Expiries)-1; i++ {
			t1 := ac.Surface.Expiries[i]
			t2 := ac.Surface.Expiries[i+1]
			
			iv1 := ac.Surface.GetIV(strike, t1)
			iv2 := ac.Surface.GetIV(strike, t2)
			
			// Total variance must be increasing
			var1 := iv1 * iv1 * t1
			var2 := iv2 * iv2 * t2
			
			if var2 < var1 {
				severity := (var1 - var2) / var1
				violations = append(violations, ArbitrageViolation{
					Type:        "calendar",
					Strike:      strike,
					Expiry1:     t1,
					Expiry2:     t2,
					Description: fmt.Sprintf("Total variance decreases from %.4f to %.4f", var1, var2),
					Severity:    severity,
				})
			}
		}
	}
	
	return violations
}

// CheckButterflyArbitrage checks convexity conditions (Durrleman's condition)
func (ac *ArbitrageChecker) CheckButterflyArbitrage() []ArbitrageViolation {
	violations := []ArbitrageViolation{}
	
	// For each expiry, check butterfly spread conditions
	for i, ttm := range ac.Surface.Expiries {
		params := ac.Surface.Parameters[i]
		forward := ac.SpotPrice * math.Exp(ac.RiskFreeRate * ttm)
		
		// Check at various strikes
		for k := -2.0; k <= 2.0; k += 0.1 {
			// g(k) = w(k)/2 where w is total implied variance
			g := func(logK float64) float64 {
				return params.computeTotalVariance(logK) / 2
			}
			
			// Check local convexity: d²g/dk² >= 0
			dk := 0.01
			g0 := g(k)
			g1 := g(k + dk)
			g2 := g(k + 2*dk)
			
			// Second derivative approximation
			d2g := (g2 - 2*g1 + g0) / (dk * dk)
			
			if d2g < -0.0001 { // Small tolerance for numerical errors
				strike := forward * math.Exp(k)
				violations = append(violations, ArbitrageViolation{
					Type:        "butterfly",
					Strike:      strike,
					Expiry1:     ttm,
					Description: fmt.Sprintf("Negative convexity: d²g/dk² = %.6f", d2g),
					Severity:    math.Abs(d2g),
				})
			}
			
			// Also check Durrleman's condition: (1 - kg'(k)/g(k))² - g'(k)²/4 + g''(k) >= 0
			dg := (g(k+dk) - g(k-dk)) / (2 * dk)
			
			if g0 > 0 {
				durrleman := math.Pow(1-k*dg/g0, 2) - dg*dg/4 + d2g
				if durrleman < -0.0001 {
					strike := forward * math.Exp(k)
					violations = append(violations, ArbitrageViolation{
						Type:        "butterfly",
						Strike:      strike,
						Expiry1:     ttm,
						Description: fmt.Sprintf("Durrleman condition violated: %.6f", durrleman),
						Severity:    math.Abs(durrleman),
					})
				}
			}
		}
	}
	
	return violations
}

// CheckNegativeDensity checks for negative probability density
func (ac *ArbitrageChecker) CheckNegativeDensity() []ArbitrageViolation {
	violations := []ArbitrageViolation{}
	
	// For each expiry
	for i, ttm := range ac.Surface.Expiries {
		params := ac.Surface.Parameters[i]
		
		// Check if any implied volatilities are negative or too small
		for k := -2.0; k <= 2.0; k += 0.1 {
			iv := params.computeIV(k)
			
			if iv < 0.01 { // 1% minimum reasonable IV
				forward := ac.SpotPrice * math.Exp(ac.RiskFreeRate * ttm)
				strike := forward * math.Exp(k)
				
				violations = append(violations, ArbitrageViolation{
					Type:        "negative_density",
					Strike:      strike,
					Expiry1:     ttm,
					Description: fmt.Sprintf("IV too low: %.2f%%", iv*100),
					Severity:    1.0 - iv/0.01,
				})
			}
		}
		
		// Check for extreme IVs that might indicate bad fits
		if params.A < 0 || params.B < 0 || params.Sigma < 0 {
			violations = append(violations, ArbitrageViolation{
				Type:        "negative_density",
				Expiry1:     ttm,
				Description: fmt.Sprintf("Invalid SVI parameters: A=%.4f, B=%.4f, Sigma=%.4f", 
					params.A, params.B, params.Sigma),
				Severity:    1.0,
			})
		}
	}
	
	return violations
}

// FilterViolationsBySeverity returns only violations above a severity threshold
func FilterViolationsBySeverity(violations []ArbitrageViolation, threshold float64) []ArbitrageViolation {
	filtered := []ArbitrageViolation{}
	for _, v := range violations {
		if v.Severity >= threshold {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// PrintArbitrageReport prints a summary of arbitrage violations
func PrintArbitrageReport(violations []ArbitrageViolation) {
	if len(violations) == 0 {
		fmt.Println("\n✅ No arbitrage violations detected")
		return
	}
	
	fmt.Printf("\n⚠️  Found %d arbitrage violations:\n", len(violations))
	
	// Group by type
	byType := make(map[string][]ArbitrageViolation)
	for _, v := range violations {
		byType[v.Type] = append(byType[v.Type], v)
	}
	
	for vType, vList := range byType {
		fmt.Printf("\n%s violations (%d):\n", vType, len(vList))
		
		// Show top 5 most severe
		maxShow := 5
		if len(vList) < maxShow {
			maxShow = len(vList)
		}
		
		// Sort by severity
		for i := 0; i < len(vList); i++ {
			for j := i + 1; j < len(vList); j++ {
				if vList[i].Severity < vList[j].Severity {
					vList[i], vList[j] = vList[j], vList[i]
				}
			}
		}
		
		for i := 0; i < maxShow; i++ {
			v := vList[i]
			if v.Strike > 0 {
				fmt.Printf("  - Strike %.0f, ", v.Strike)
			}
			if v.Expiry2 > 0 {
				fmt.Printf("Expiries %.3f-%.3f: ", v.Expiry1, v.Expiry2)
			} else {
				fmt.Printf("Expiry %.3f: ", v.Expiry1)
			}
			fmt.Printf("%s (severity: %.2f)\n", v.Description, v.Severity)
		}
		
		if len(vList) > maxShow {
			fmt.Printf("  ... and %d more\n", len(vList)-maxShow)
		}
	}
}