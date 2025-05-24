package main

import (
	"fmt"
	"log"
	"math"
	"sort"
)

// VolSurfacePoint represents a point on the volatility surface
type VolSurfacePoint struct {
	Strike       float64 `json:"strike"`
	TimeToExpiry float64 `json:"timeToExpiry"`
	Moneyness    float64 `json:"moneyness"`
	IV           float64 `json:"iv"`
	OptionType   string  `json:"optionType"`
}

// VolSurfaceModel represents a fitted volatility surface
type VolSurfaceModel struct {
	Points    []VolSurfacePoint
	Summary   VolSurfaceSummary
	ByExpiry  map[string][]VolSurfacePoint // Grouped by expiry date
	ByStrike  map[float64][]VolSurfacePoint // Grouped by strike
}

// VolSurfaceSummary provides summary statistics of the volatility surface
type VolSurfaceSummary struct {
	TotalPoints      int                 `json:"totalPoints"`
	CallPoints       int                 `json:"callPoints"`
	PutPoints        int                 `json:"putPoints"`
	AvgIV            float64             `json:"avgIV"`
	MinIV            float64             `json:"minIV"`
	MaxIV            float64             `json:"maxIV"`
	AvgTimeToExpiry  float64             `json:"avgTimeToExpiry"`
	MinTimeToExpiry  float64             `json:"minTimeToExpiry"`
	MaxTimeToExpiry  float64             `json:"maxTimeToExpiry"`
	ATMVolsByExpiry  map[string]float64  `json:"atmVolsByExpiry"` // ATM volatilities by expiry
	VolSkewByExpiry  map[string]float64  `json:"volSkewByExpiry"` // Put-call vol skew by expiry
}

// FitVolatilitySurface analyzes and fits the volatility surface
func (vsa *VolSurfaceAnalyzer) FitVolatilitySurface() VolSurfaceModel {
	log.Println("Fitting volatility surface...")
	
	var points []VolSurfacePoint
	
	// Extract valid volatility points
	for _, option := range vsa.options {
		if !math.IsNaN(option.IV) && option.IV > 0 && option.DaysToExpiry > 0 {
			point := VolSurfacePoint{
				Strike:       option.Strike,
				TimeToExpiry: option.TimeToExpiryYrs,
				Moneyness:    option.Moneyness,
				IV:           option.IV,
				OptionType:   option.Type,
			}
			points = append(points, point)
		}
	}
	
	// Create model
	model := VolSurfaceModel{
		Points:   points,
		ByExpiry: make(map[string][]VolSurfacePoint),
		ByStrike: make(map[float64][]VolSurfacePoint),
	}
	
	// Group by expiry and strike
	for _, point := range points {
		expiryKey := fmt.Sprintf("%.4f", point.TimeToExpiry)
		model.ByExpiry[expiryKey] = append(model.ByExpiry[expiryKey], point)
		model.ByStrike[point.Strike] = append(model.ByStrike[point.Strike], point)
	}
	
	// Calculate summary statistics
	model.Summary = vsa.calculateSummary(points)
	
	log.Printf("Fitted volatility surface with %d points", len(points))
	return model
}

// calculateSummary computes summary statistics for the volatility surface
func (vsa *VolSurfaceAnalyzer) calculateSummary(points []VolSurfacePoint) VolSurfaceSummary {
	if len(points) == 0 {
		return VolSurfaceSummary{}
	}
	
	summary := VolSurfaceSummary{
		TotalPoints:     len(points),
		ATMVolsByExpiry: make(map[string]float64),
		VolSkewByExpiry: make(map[string]float64),
	}
	
	// Calculate basic statistics
	totalIV := 0.0
	totalTime := 0.0
	summary.MinIV = math.Inf(1)
	summary.MaxIV = math.Inf(-1)
	summary.MinTimeToExpiry = math.Inf(1)
	summary.MaxTimeToExpiry = math.Inf(-1)
	
	for _, point := range points {
		if point.OptionType == "C" {
			summary.CallPoints++
		} else {
			summary.PutPoints++
		}
		
		totalIV += point.IV
		totalTime += point.TimeToExpiry
		
		if point.IV < summary.MinIV {
			summary.MinIV = point.IV
		}
		if point.IV > summary.MaxIV {
			summary.MaxIV = point.IV
		}
		if point.TimeToExpiry < summary.MinTimeToExpiry {
			summary.MinTimeToExpiry = point.TimeToExpiry
		}
		if point.TimeToExpiry > summary.MaxTimeToExpiry {
			summary.MaxTimeToExpiry = point.TimeToExpiry
		}
	}
	
	summary.AvgIV = totalIV / float64(len(points))
	summary.AvgTimeToExpiry = totalTime / float64(len(points))
	
	// Calculate ATM volatilities and skew by expiry
	expiryGroups := make(map[string][]VolSurfacePoint)
	for _, point := range points {
		expiryKey := fmt.Sprintf("%.4f", point.TimeToExpiry)
		expiryGroups[expiryKey] = append(expiryGroups[expiryKey], point)
	}
	
	for expiryKey, expiryPoints := range expiryGroups {
		atmVol, skew := vsa.calculateATMVolAndSkew(expiryPoints)
		summary.ATMVolsByExpiry[expiryKey] = atmVol
		summary.VolSkewByExpiry[expiryKey] = skew
	}
	
	return summary
}

// calculateATMVolAndSkew calculates ATM volatility and put-call skew for a given expiry
func (vsa *VolSurfaceAnalyzer) calculateATMVolAndSkew(points []VolSurfacePoint) (float64, float64) {
	if len(points) == 0 {
		return math.NaN(), math.NaN()
	}
	
	// Find ATM options (closest to current spot price)
	atmVol := math.NaN()
	minMoneynessDistance := math.Inf(1)
	
	for _, point := range points {
		distance := math.Abs(point.Moneyness - 1.0) // Distance from ATM (moneyness = 1)
		if distance < minMoneynessDistance {
			minMoneynessDistance = distance
			atmVol = point.IV
		}
	}
	
	// Calculate vol skew (difference between OTM put and call volatilities)
	otmPutVol := math.NaN()
	otmCallVol := math.NaN()
	
	// Look for OTM puts (moneyness < 0.9) and OTM calls (moneyness > 1.1)
	for _, point := range points {
		if point.OptionType == "P" && point.Moneyness < 0.9 && (math.IsNaN(otmPutVol) || point.IV > otmPutVol) {
			otmPutVol = point.IV
		}
		if point.OptionType == "C" && point.Moneyness > 1.1 && (math.IsNaN(otmCallVol) || point.IV > otmCallVol) {
			otmCallVol = point.IV
		}
	}
	
	var skew float64
	if !math.IsNaN(otmPutVol) && !math.IsNaN(otmCallVol) {
		skew = otmPutVol - otmCallVol // Positive skew means puts are more expensive
	} else {
		skew = math.NaN()
	}
	
	return atmVol, skew
}

// PrintVolSurfaceAnalysis prints detailed analysis of the volatility surface
func (vsa *VolSurfaceAnalyzer) PrintVolSurfaceAnalysis(model VolSurfaceModel) {
	fmt.Println("\n=== VOLATILITY SURFACE ANALYSIS ===")
	
	s := model.Summary
	fmt.Printf("Total Points: %d (Calls: %d, Puts: %d)\n", s.TotalPoints, s.CallPoints, s.PutPoints)
	fmt.Printf("Average IV: %.2f%%\n", s.AvgIV*100)
	fmt.Printf("IV Range: %.2f%% - %.2f%%\n", s.MinIV*100, s.MaxIV*100)
	fmt.Printf("Time to Expiry Range: %.4f - %.4f years\n", s.MinTimeToExpiry, s.MaxTimeToExpiry)
	
	fmt.Println("\n--- ATM Volatilities by Expiry ---")
	// Sort expiries for ordered display
	var expiries []string
	for expiry := range s.ATMVolsByExpiry {
		expiries = append(expiries, expiry)
	}
	sort.Strings(expiries)
	
	for _, expiry := range expiries {
		atmVol := s.ATMVolsByExpiry[expiry]
		skew := s.VolSkewByExpiry[expiry]
		if !math.IsNaN(atmVol) {
			fmt.Printf("  %s years: ATM %.2f%%", expiry, atmVol*100)
			if !math.IsNaN(skew) {
				fmt.Printf(", Skew: %.2f%%", skew*100)
			}
			fmt.Println()
		}
	}
	
	fmt.Println("\n--- Volatility Term Structure ---")
	vsa.printTermStructure(model)
	
	fmt.Println("\n--- Volatility Smile Analysis ---")
	vsa.printVolatilitySmile(model)
}

// printTermStructure analyzes the volatility term structure
func (vsa *VolSurfaceAnalyzer) printTermStructure(model VolSurfaceModel) {
	// Group ATM volatilities by time to expiry
	type TermPoint struct {
		TimeToExpiry float64
		ATMVol       float64
	}
	
	var termPoints []TermPoint
	for expiryKey, atmVol := range model.Summary.ATMVolsByExpiry {
		if !math.IsNaN(atmVol) {
			var time float64
			fmt.Sscanf(expiryKey, "%f", &time)
			termPoints = append(termPoints, TermPoint{
				TimeToExpiry: time,
				ATMVol:       atmVol,
			})
		}
	}
	
	// Sort by time to expiry
	sort.Slice(termPoints, func(i, j int) bool {
		return termPoints[i].TimeToExpiry < termPoints[j].TimeToExpiry
	})
	
	for _, point := range termPoints {
		fmt.Printf("  %.2f years: %.2f%% ATM vol\n", point.TimeToExpiry, point.ATMVol*100)
	}
	
	// Analyze term structure shape
	if len(termPoints) >= 2 {
		shortTerm := termPoints[0].ATMVol
		longTerm := termPoints[len(termPoints)-1].ATMVol
		fmt.Printf("Term Structure: %.2f%% (short) → %.2f%% (long)", shortTerm*100, longTerm*100)
		if longTerm > shortTerm {
			fmt.Printf(" [Upward Sloping]\n")
		} else if longTerm < shortTerm {
			fmt.Printf(" [Downward Sloping]\n")
		} else {
			fmt.Printf(" [Flat]\n")
		}
	}
}

// printVolatilitySmile analyzes the volatility smile for each expiry
func (vsa *VolSurfaceAnalyzer) printVolatilitySmile(model VolSurfaceModel) {
	// Analyze smile for each expiry with sufficient data
	for expiryKey, expiryPoints := range model.ByExpiry {
		if len(expiryPoints) < 3 {
			continue // Need at least 3 points for smile analysis
		}
		
		var time float64
		fmt.Sscanf(expiryKey, "%f", &time)
		
		// Sort by moneyness
		sort.Slice(expiryPoints, func(i, j int) bool {
			return expiryPoints[i].Moneyness < expiryPoints[j].Moneyness
		})
		
		fmt.Printf("  %.3f years expiry:\n", time)
		
		// Find wings (OTM options)
		leftWing := math.NaN()
		rightWing := math.NaN()
		atmVol := math.NaN()
		atmMoneyness := 1.0
		
		minATMDistance := math.Inf(1)
		for _, point := range expiryPoints {
			distance := math.Abs(point.Moneyness - atmMoneyness)
			if distance < minATMDistance {
				minATMDistance = distance
				atmVol = point.IV
			}
			
			if point.Moneyness < 0.9 {
				leftWing = point.IV
			}
			if point.Moneyness > 1.1 && math.IsNaN(rightWing) {
				rightWing = point.IV
			}
		}
		
		// Print smile characteristics
		if !math.IsNaN(atmVol) {
			fmt.Printf("    ATM (%.2f): %.2f%%\n", atmMoneyness, atmVol*100)
		}
		if !math.IsNaN(leftWing) {
			fmt.Printf("    Left Wing: %.2f%%\n", leftWing*100)
		}
		if !math.IsNaN(rightWing) {
			fmt.Printf("    Right Wing: %.2f%%\n", rightWing*100)
		}
		
		// Calculate smile asymmetry
		if !math.IsNaN(leftWing) && !math.IsNaN(rightWing) {
			asymmetry := leftWing - rightWing
			fmt.Printf("    Asymmetry: %.2f%% ", asymmetry*100)
			if asymmetry > 0.05 {
				fmt.Printf("[Put Skew]\n")
			} else if asymmetry < -0.05 {
				fmt.Printf("[Call Skew]\n")
			} else {
				fmt.Printf("[Balanced]\n")
			}
		}
	}
}