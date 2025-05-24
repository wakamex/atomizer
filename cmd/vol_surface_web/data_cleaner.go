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
		MaxSpreadPct:    10.0,  // 10% max spread (more strict for crypto)
		MinVolume:       5,     // Minimum 5 contracts volume for deep OTM/ITM
		MinMoneyness:    0.6,   // 60% of spot (based on Deribit data)
		MaxMoneyness:    1.8,   // 180% of spot (based on Deribit data)
		MinIV:           0.3,   // 30% minimum IV (Deribit rarely below this)
		MaxIV:           2.5,   // 250% maximum IV (realistic for crypto)
		MinDaysToExpiry: 1.0,   // 1 day minimum
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

// PreprocessOptionsForSVI applies strict filtering for SVI/SSVI fitting
func PreprocessOptionsForSVI(options []OptionData, spotPrice float64) []OptionData {
	filtered := []OptionData{}
	
	for _, opt := range options {
		// Skip if no valid IV
		if math.IsNaN(opt.ImpliedVolatility) || opt.ImpliedVolatility <= 0 {
			continue
		}
		
		// Remove extreme moneyness (based on Deribit market data)
		// Active trading typically from 0.6 to 1.8 moneyness
		moneyness := opt.Strike / spotPrice
		if moneyness < 0.6 || moneyness > 1.8 {
			continue
		}
		
		// Remove suspiciously low IV (10% minimum)
		if opt.ImpliedVolatility < 0.1 {
			continue
		}
		
		// Remove extremely high IV
		if opt.ImpliedVolatility > 5.0 {
			continue
		}
		
		// Remove illiquid options more aggressively
		if opt.Volume < 5 && opt.OpenInterest < 20 {
			continue
		}
		
		// Remove options with wide spreads
		if opt.SpreadPct > 10 {
			continue
		}
		
		// Remove very short-dated options
		if opt.DaysToExpiry < 1 {
			continue
		}
		
		// Remove very long-dated options (> 1 year)
		if opt.DaysToExpiry > 365 {
			continue
		}
		
		filtered = append(filtered, opt)
	}
	
	// Additional filtering: remove statistical outliers within each expiry group
	filtered = removeOutliersPerExpiry(filtered)
	
	return filtered
}

// removeOutliersPerExpiry removes IV outliers within each expiry group
func removeOutliersPerExpiry(options []OptionData) []OptionData {
	// Group by expiry
	expiryGroups := make(map[float64][]OptionData)
	for _, opt := range options {
		// Round expiry to nearest day for grouping
		expiryKey := math.Round(opt.DaysToExpiry)
		expiryGroups[expiryKey] = append(expiryGroups[expiryKey], opt)
	}
	
	filtered := []OptionData{}
	
	for _, group := range expiryGroups {
		if len(group) < 5 {
			// Keep all if too few points
			filtered = append(filtered, group...)
			continue
		}
		
		// Calculate IV statistics for this expiry
		ivs := make([]float64, len(group))
		for i, opt := range group {
			ivs[i] = opt.ImpliedVolatility
		}
		
		// Calculate quartiles
		q1, q3 := calculateQuartiles(ivs)
		iqr := q3 - q1
		
		// More conservative bounds for SVI fitting
		lowerBound := q1 - 1.0*iqr // Instead of 1.5*IQR
		upperBound := q3 + 1.0*iqr
		
		// Filter based on bounds
		for _, opt := range group {
			if opt.ImpliedVolatility >= lowerBound && opt.ImpliedVolatility <= upperBound {
				filtered = append(filtered, opt)
			}
		}
	}
	
	return filtered
}