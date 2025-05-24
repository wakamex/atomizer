package main

import (
	"fmt"
	"math"
)

// DataQualityMetrics tracks data quality statistics
type DataQualityMetrics struct {
	TotalOptions       int
	FilteredOptions    int
	ExpiredOptions     int
	WideSpreadOptions  int
	LowLiquidityOptions int
	OutlierOptions     int
	BadPriceOptions    int
}

// DataCleaner handles option data cleaning and filtering
type DataCleaner struct {
	MaxSpreadPct     float64 // Maximum bid-ask spread percentage
	MinVolume        int     // Minimum volume for deep OTM/ITM options
	MinMoneyness     float64 // Minimum moneyness to consider
	MaxMoneyness     float64 // Maximum moneyness to consider
	MinIV            float64 // Minimum reasonable IV
	MaxIV            float64 // Maximum reasonable IV
	MinDaysToExpiry  float64 // Minimum days to expiry
	Metrics          DataQualityMetrics
}

// NewDataCleaner creates a new data cleaner with default parameters
func NewDataCleaner() *DataCleaner {
	return &DataCleaner{
		MaxSpreadPct:    20.0,  // 20% max spread
		MinVolume:       10,    // Minimum 10 contracts volume for deep OTM/ITM
		MinMoneyness:    0.5,   // 50% of spot
		MaxMoneyness:    2.0,   // 200% of spot
		MinIV:           0.05,  // 5% minimum IV
		MaxIV:           5.0,   // 500% maximum IV
		MinDaysToExpiry: 0.1,   // 0.1 days minimum
	}
}

// CleanOptionData filters and cleans option data
func (dc *DataCleaner) CleanOptionData(options []OptionData, spotPrice float64) []OptionData {
	dc.Metrics = DataQualityMetrics{TotalOptions: len(options)}
	
	filtered := []OptionData{}
	
	for _, opt := range options {
		// Calculate additional metrics if needed
		if opt.Moneyness == 0 && spotPrice > 0 {
			opt.Moneyness = opt.Strike / spotPrice
		}
		
		// Check various filtering criteria
		if !dc.passesFilters(opt, spotPrice) {
			continue
		}
		
		// Clean the data
		opt = dc.cleanOption(opt)
		
		filtered = append(filtered, opt)
	}
	
	dc.Metrics.FilteredOptions = len(filtered)
	return filtered
}

// passesFilters checks if an option passes all filtering criteria
func (dc *DataCleaner) passesFilters(opt OptionData, spotPrice float64) bool {
	// Remove expired options
	if opt.DaysToExpiry <= dc.MinDaysToExpiry {
		dc.Metrics.ExpiredOptions++
		return false
	}
	
	// Remove options with bad prices
	if math.IsNaN(opt.MidETH) || opt.MidETH <= 0 {
		if math.IsNaN(opt.LastETH) || opt.LastETH <= 0 {
			dc.Metrics.BadPriceOptions++
			return false
		}
	}
	
	// Check bid-ask spread if available
	if opt.BidETH > 0 && opt.AskETH > 0 {
		spread := opt.AskETH - opt.BidETH
		midPrice := (opt.BidETH + opt.AskETH) / 2
		spreadPct := (spread / midPrice) * 100
		
		if spreadPct > dc.MaxSpreadPct {
			dc.Metrics.WideSpreadOptions++
			return false
		}
	}
	
	// Check moneyness and liquidity
	moneyness := opt.Strike / spotPrice
	isDeepOTM := (opt.IsPut && moneyness > 1.3) || (!opt.IsPut && moneyness < 0.7)
	isDeepITM := (opt.IsPut && moneyness < 0.7) || (!opt.IsPut && moneyness > 1.3)
	
	if (isDeepOTM || isDeepITM) && opt.Volume < float64(dc.MinVolume) {
		dc.Metrics.LowLiquidityOptions++
		return false
	}
	
	// Remove extreme moneyness
	if moneyness < dc.MinMoneyness || moneyness > dc.MaxMoneyness {
		dc.Metrics.LowLiquidityOptions++
		return false
	}
	
	// Check IV outliers
	if !math.IsNaN(opt.ImpliedVolatility) {
		if opt.ImpliedVolatility < dc.MinIV || opt.ImpliedVolatility > dc.MaxIV {
			dc.Metrics.OutlierOptions++
			return false
		}
	}
	
	return true
}

// cleanOption cleans individual option data
func (dc *DataCleaner) cleanOption(opt OptionData) OptionData {
	// Use mid price if available, otherwise last price
	if math.IsNaN(opt.MidETH) || opt.MidETH <= 0 {
		if !math.IsNaN(opt.LastETH) && opt.LastETH > 0 {
			opt.MidETH = opt.LastETH
		}
	}
	
	// Ensure time to expiry is positive
	if opt.TimeToExpiry < 0 {
		opt.TimeToExpiry = opt.DaysToExpiry / 365.0
	}
	
	// Cap extreme IVs
	if !math.IsNaN(opt.ImpliedVolatility) {
		if opt.ImpliedVolatility > dc.MaxIV {
			opt.ImpliedVolatility = dc.MaxIV
		} else if opt.ImpliedVolatility < dc.MinIV {
			opt.ImpliedVolatility = dc.MinIV
		}
	}
	
	return opt
}

// DetectOutliers uses statistical methods to detect IV outliers
func (dc *DataCleaner) DetectOutliers(options []OptionData) []int {
	outlierIndices := []int{}
	
	// Group by expiry and option type
	groups := make(map[string][]int)
	for i, opt := range options {
		if math.IsNaN(opt.ImpliedVolatility) {
			continue
		}
		
		// Create group key
		expiryBucket := int(opt.TimeToExpiry * 12) // Monthly buckets
		optType := "C"
		if opt.IsPut {
			optType = "P"
		}
		key := fmt.Sprintf("%d_%s", expiryBucket, optType)
		
		groups[key] = append(groups[key], i)
	}
	
	// For each group, detect outliers using IQR method
	for _, indices := range groups {
		if len(indices) < 4 {
			continue // Need at least 4 points
		}
		
		// Extract IVs for this group
		ivs := make([]float64, len(indices))
		for i, idx := range indices {
			ivs[i] = options[idx].ImpliedVolatility
		}
		
		// Calculate quartiles
		q1, q3 := calculateQuartiles(ivs)
		iqr := q3 - q1
		
		// Outlier thresholds
		lowerBound := q1 - 1.5*iqr
		upperBound := q3 + 1.5*iqr
		
		// Mark outliers
		for i, idx := range indices {
			if ivs[i] < lowerBound || ivs[i] > upperBound {
				outlierIndices = append(outlierIndices, idx)
			}
		}
	}
	
	return outlierIndices
}

// calculateQuartiles calculates Q1 and Q3 for a slice of values
func calculateQuartiles(values []float64) (q1, q3 float64) {
	n := len(values)
	if n < 4 {
		return values[0], values[n-1]
	}
	
	// Sort values
	sorted := make([]float64, n)
	copy(sorted, values)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	
	// Calculate quartile positions
	q1Pos := float64(n+1) / 4.0
	q3Pos := 3.0 * float64(n+1) / 4.0
	
	// Interpolate if necessary
	q1 = interpolateValue(sorted, q1Pos)
	q3 = interpolateValue(sorted, q3Pos)
	
	return q1, q3
}

// interpolateValue interpolates a value at a given position
func interpolateValue(sorted []float64, pos float64) float64 {
	lower := int(pos) - 1
	upper := lower + 1
	
	if upper >= len(sorted) {
		return sorted[len(sorted)-1]
	}
	if lower < 0 {
		return sorted[0]
	}
	
	weight := pos - float64(lower+1)
	return sorted[lower]*(1-weight) + sorted[upper]*weight
}

// PrintDataQualityReport prints a summary of data quality metrics
func (dc *DataCleaner) PrintDataQualityReport() {
	m := dc.Metrics
	fmt.Println("\n📊 Data Quality Report")
	fmt.Printf("Total options: %d\n", m.TotalOptions)
	fmt.Printf("Filtered options: %d (%.1f%% kept)\n", 
		m.FilteredOptions, 
		float64(m.FilteredOptions)/float64(m.TotalOptions)*100)
	
	fmt.Println("\nFiltered out:")
	if m.ExpiredOptions > 0 {
		fmt.Printf("  - Expired: %d\n", m.ExpiredOptions)
	}
	if m.BadPriceOptions > 0 {
		fmt.Printf("  - Bad prices: %d\n", m.BadPriceOptions)
	}
	if m.WideSpreadOptions > 0 {
		fmt.Printf("  - Wide spreads: %d\n", m.WideSpreadOptions)
	}
	if m.LowLiquidityOptions > 0 {
		fmt.Printf("  - Low liquidity: %d\n", m.LowLiquidityOptions)
	}
	if m.OutlierOptions > 0 {
		fmt.Printf("  - IV outliers: %d\n", m.OutlierOptions)
	}
}