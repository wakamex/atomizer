package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	ccxt "github.com/ccxt/ccxt/go/v4"
)

// OptionData represents a single option's market data
type OptionData struct {
	Symbol            string
	Type              string // "C" or "P"
	Strike            float64
	Expiry            time.Time
	DaysToExpiry      float64
	BidETH            float64
	AskETH            float64
	MidETH            float64
	LastETH           float64
	SpreadETH         float64
	SpreadPct         float64
	BidUSD            float64
	AskUSD            float64
	MidUSD            float64
	Volume            float64
	OpenInterest      float64
	IV                float64 // Implied Volatility from exchange
	Delta             float64
	Gamma             float64
	Theta             float64
	Vega              float64
	ETHIndexPrice     float64
	Moneyness         float64
	TimeToExpiryYrs   float64
	
	// Additional fields for our analysis
	IsPut             bool
	TimeToExpiry      float64 // In years
	ImpliedVolatility float64 // Our calculated IV
	MarketPriceETH    float64 // The price we use for IV calc
}

// VolSurfaceAnalyzer handles volatility surface analysis
type VolSurfaceAnalyzer struct {
	options      []OptionData
	exchange     *ccxt.Deribit
	spotPrice    float64
	riskFreeRate float64
	surface      *SVISurface
	ssviSurface  *SSVISurface
	cleaner      *DataCleaner
}

// NewVolSurfaceAnalyzer creates a new volatility surface analyzer
func NewVolSurfaceAnalyzer() *VolSurfaceAnalyzer {
	exchange := ccxt.NewDeribit(map[string]interface{}{
		"rateLimit":       10,
		"enableRateLimit": true,
	})
	
	return &VolSurfaceAnalyzer{
		options:      make([]OptionData, 0),
		exchange:     &exchange,
		riskFreeRate: 0.05, // Default 5%
		cleaner:      NewDataCleaner(),
	}
}

// LoadCSV loads options data from existing CSV file
func (vsa *VolSurfaceAnalyzer) LoadCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %v", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file is empty or has no data rows")
	}

	// Skip header row
	for i, record := range records[1:] {
		if len(record) < 24 {
			log.Printf("Skipping row %d: insufficient columns", i+2)
			continue
		}

		option, err := vsa.parseOptionRow(record)
		if err != nil {
			log.Printf("Error parsing row %d: %v", i+2, err)
			continue
		}

		vsa.options = append(vsa.options, option)
	}

	log.Printf("Loaded %d options from CSV", len(vsa.options))
	return nil
}

// parseOptionRow parses a CSV row into OptionData
func (vsa *VolSurfaceAnalyzer) parseOptionRow(record []string) (OptionData, error) {
	option := OptionData{}
	
	option.Symbol = record[0]
	option.Type = record[1]
	
	if strike, err := strconv.ParseFloat(record[2], 64); err != nil {
		return option, fmt.Errorf("invalid strike: %v", err)
	} else {
		option.Strike = strike
	}
	
	if expiry, err := time.Parse("2006-01-02", record[3]); err != nil {
		return option, fmt.Errorf("invalid expiry: %v", err)
	} else {
		option.Expiry = expiry
	}
	
	// Parse numeric fields with NaN handling
	option.DaysToExpiry = parseFloatWithNaN(record[4])
	option.BidETH = parseFloatWithNaN(record[5])
	option.AskETH = parseFloatWithNaN(record[6])
	option.MidETH = parseFloatWithNaN(record[7])
	option.LastETH = parseFloatWithNaN(record[8])
	option.SpreadETH = parseFloatWithNaN(record[9])
	option.SpreadPct = parseFloatWithNaN(record[10])
	option.BidUSD = parseFloatWithNaN(record[11])
	option.AskUSD = parseFloatWithNaN(record[12])
	option.MidUSD = parseFloatWithNaN(record[13])
	option.Volume = parseFloatWithNaN(record[14])
	option.OpenInterest = parseFloatWithNaN(record[15])
	option.IV = parseFloatWithNaN(record[16])
	option.Delta = parseFloatWithNaN(record[17])
	option.Gamma = parseFloatWithNaN(record[18])
	option.Theta = parseFloatWithNaN(record[19])
	option.Vega = parseFloatWithNaN(record[20])
	option.ETHIndexPrice = parseFloatWithNaN(record[21])
	option.Moneyness = parseFloatWithNaN(record[22])
	option.TimeToExpiryYrs = parseFloatWithNaN(record[23])
	
	return option, nil
}

// parseFloatWithNaN handles NaN strings in CSV
func parseFloatWithNaN(s string) float64 {
	if s == "" || s == "NaN" || s == "nan" {
		return math.NaN()
	}
	if val, err := strconv.ParseFloat(s, 64); err != nil {
		return math.NaN()
	} else {
		return val
	}
}

// UpdatePricing fetches current market prices for options missing pricing data
func (vsa *VolSurfaceAnalyzer) UpdatePricing() error {
	log.Println("Updating pricing data from Deribit...")
	
	// For now, just set a dummy spot price to test the CSV functionality
	vsa.spotPrice = 3500.0
	log.Printf("Using dummy ETH spot price: $%.2f", vsa.spotPrice)
	
	// Skip Deribit API calls for now to test basic functionality
	log.Println("Skipping Deribit API calls for basic testing")
	return nil
}

// CalculateImpliedVolatilities calculates IV for options with market prices
func (vsa *VolSurfaceAnalyzer) CalculateImpliedVolatilities() {
	log.Println("Calculating implied volatilities...")
	
	const riskFreeRate = 0.05 // 5% risk-free rate
	calculatedCount := 0
	
	for i := range vsa.options {
		option := &vsa.options[i]
		
		// Populate additional fields
		option.IsPut = (option.Type == "P")
		option.TimeToExpiry = option.TimeToExpiryYrs
		
		// Get market price (prefer MidETH, fall back to LastETH)
		var marketPriceETH float64
		if !math.IsNaN(option.MidETH) && option.MidETH > 0 {
			marketPriceETH = option.MidETH
		} else if !math.IsNaN(option.LastETH) && option.LastETH > 0 {
			marketPriceETH = option.LastETH
		} else {
			continue // No usable price data
		}
		
		option.MarketPriceETH = marketPriceETH
		
		// Skip if expired or price too low
		if option.DaysToExpiry <= 0 || marketPriceETH <= 0.001 {
			continue
		}
		
		// Convert market price to USD
		marketPriceUSD := marketPriceETH * vsa.spotPrice
		
		// Set up Black-Scholes parameters
		params := BlackScholesParams{
			SpotPrice:    vsa.spotPrice,
			StrikePrice:  option.Strike,
			TimeToExpiry: option.TimeToExpiryYrs,
			RiskFreeRate: riskFreeRate,
			IsCall:       option.Type == "C",
		}
		
		// Calculate implied volatility
		iv := ImpliedVolatility(marketPriceUSD, params)
		
		if !math.IsNaN(iv) && iv > 0 {
			option.ImpliedVolatility = iv
			option.IV = iv // Keep both for compatibility
			
			// Calculate Greeks using the computed IV
			params.Volatility = iv
			option.Delta = BlackScholesDelta(params)
			option.Gamma = BlackScholesGamma(params)
			option.Theta = BlackScholesTheta(params)
			option.Vega = BlackScholesVega(params)
			
			calculatedCount++
		}
	}
	
	log.Printf("Calculated implied volatilities for %d options", calculatedCount)
}

// fetchOptionPricing gets current market data for a specific option
func (vsa *VolSurfaceAnalyzer) fetchOptionPricing(option *OptionData) error {
	// Convert symbol to CCXT format: ETH/USD:ETH-YYMMDD-STRIKE-C
	ccxtSymbol := fmt.Sprintf("ETH/USD:%s", option.Symbol)
	
	// Fetch order book
	orderBook, err := vsa.exchange.FetchOrderBook(ccxtSymbol)
	if err != nil {
		return err
	}
	
	// Update pricing data
	if len(orderBook.Bids) > 0 {
		option.BidETH = orderBook.Bids[0][0]
		option.BidUSD = option.BidETH * vsa.spotPrice
	}
	
	if len(orderBook.Asks) > 0 {
		option.AskETH = orderBook.Asks[0][0]
		option.AskUSD = option.AskETH * vsa.spotPrice
	}
	
	if !math.IsNaN(option.BidETH) && !math.IsNaN(option.AskETH) {
		option.MidETH = (option.BidETH + option.AskETH) / 2
		option.MidUSD = option.MidETH * vsa.spotPrice
		option.SpreadETH = option.AskETH - option.BidETH
		option.SpreadPct = (option.SpreadETH / option.MidETH) * 100
	}
	
	// Update current ETH index price
	option.ETHIndexPrice = vsa.spotPrice
	
	// Recalculate time-based fields
	now := time.Now()
	option.DaysToExpiry = option.Expiry.Sub(now).Hours() / 24
	option.TimeToExpiryYrs = option.DaysToExpiry / 365.25
	option.Moneyness = vsa.spotPrice / option.Strike
	
	return nil
}

// SaveEnhancedCSV saves the updated options data to a new CSV file
func (vsa *VolSurfaceAnalyzer) SaveEnhancedCSV(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// Write header
	header := []string{
		"Symbol", "Type", "Strike", "Expiry", "Days_To_Expiry",
		"Bid_ETH", "Ask_ETH", "Mid_ETH", "Last_ETH", "Spread_ETH", "Spread_Pct",
		"Bid_USD", "Ask_USD", "Mid_USD", "Volume", "Open_Interest",
		"IV", "Delta", "Gamma", "Theta", "Vega",
		"ETH_Index_Price", "Moneyness", "Time_To_Expiry_Years",
	}
	
	if err := writer.Write(header); err != nil {
		return err
	}
	
	// Write data rows
	for _, option := range vsa.options {
		record := []string{
			option.Symbol,
			option.Type,
			fmt.Sprintf("%.0f", option.Strike),
			option.Expiry.Format("2006-01-02"),
			formatFloat(option.DaysToExpiry, 1),
			formatFloat(option.BidETH, 4),
			formatFloat(option.AskETH, 4),
			formatFloat(option.MidETH, 4),
			formatFloat(option.LastETH, 4),
			formatFloat(option.SpreadETH, 4),
			formatFloat(option.SpreadPct, 2),
			formatFloat(option.BidUSD, 2),
			formatFloat(option.AskUSD, 2),
			formatFloat(option.MidUSD, 2),
			formatFloat(option.Volume, 0),
			formatFloat(option.OpenInterest, 0),
			formatFloat(option.IV, 4),
			formatFloat(option.Delta, 4),
			formatFloat(option.Gamma, 4),
			formatFloat(option.Theta, 4),
			formatFloat(option.Vega, 4),
			formatFloat(option.ETHIndexPrice, 2),
			formatFloat(option.Moneyness, 3),
			formatFloat(option.TimeToExpiryYrs, 4),
		}
		
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	
	log.Printf("Saved enhanced CSV with %d options to %s", len(vsa.options), filename)
	return nil
}

// formatFloat formats float64 with NaN handling
func formatFloat(val float64, precision int) string {
	if math.IsNaN(val) {
		return "NaN"
	}
	return fmt.Sprintf("%.*f", precision, val)
}

// AnalyzeMarkets provides summary statistics of the options market
func (vsa *VolSurfaceAnalyzer) AnalyzeMarkets() {
	if len(vsa.options) == 0 {
		log.Println("No options data to analyze")
		return
	}
	
	log.Println("\n=== MARKET ANALYSIS ===")
	
	calls, puts := 0, 0
	activeOptions := 0
	totalVolume := 0.0
	totalOI := 0.0
	validIVCount := 0
	avgIV := 0.0
	
	for _, option := range vsa.options {
		if option.Type == "C" {
			calls++
		} else {
			puts++
		}
		
		if option.DaysToExpiry > 0 {
			activeOptions++
		}
		
		if !math.IsNaN(option.Volume) {
			totalVolume += option.Volume
		}
		
		if !math.IsNaN(option.OpenInterest) {
			totalOI += option.OpenInterest
		}
		
		if !math.IsNaN(option.IV) && option.IV > 0 {
			avgIV += option.IV
			validIVCount++
		}
	}
	
	if validIVCount > 0 {
		avgIV /= float64(validIVCount)
	}
	
	log.Printf("Total Options: %d (Calls: %d, Puts: %d)", len(vsa.options), calls, puts)
	log.Printf("Active Options: %d", activeOptions)
	log.Printf("Total Volume: %.0f", totalVolume)
	log.Printf("Total Open Interest: %.0f", totalOI)
	log.Printf("Average IV: %.2f%% (%d valid)", avgIV*100, validIVCount)
	log.Printf("Current ETH Price: $%.2f", vsa.spotPrice)
}

// CleanData applies data cleaning filters to the options
func (vsa *VolSurfaceAnalyzer) CleanData() (cleanedOptions, filteredOptions []OptionData) {
	log.Println("Cleaning option data...")
	
	// Keep a copy of all original options
	allOptions := make([]OptionData, len(vsa.options))
	copy(allOptions, vsa.options)
	
	// Clean the data
	vsa.options = vsa.cleaner.CleanOptionData(vsa.options, vsa.spotPrice)
	vsa.cleaner.PrintDataQualityReport()
	
	// Determine which options were filtered out
	cleanedMap := make(map[string]bool)
	for _, opt := range vsa.options {
		cleanedMap[opt.Symbol] = true
	}
	
	filteredOptions = []OptionData{}
	for _, opt := range allOptions {
		if !cleanedMap[opt.Symbol] {
			filteredOptions = append(filteredOptions, opt)
		}
	}
	
	return vsa.options, filteredOptions
}

// FitSVISurface fits the SVI model to the cleaned data
func (vsa *VolSurfaceAnalyzer) FitSVISurface() error {
	log.Println("Fitting volatility surfaces...")
	
	// First, apply strict preprocessing for SVI fitting
	// TEMPORARY: Use all cleaned options instead of aggressive preprocessing
	preprocessedOptions := vsa.options // PreprocessOptionsForSVI(vsa.options, vsa.spotPrice)
	log.Printf("Using %d options for SVI fitting (from %d cleaned)", len(preprocessedOptions), len(vsa.options))
	
	// Debug: show distribution of preprocessed data
	if len(preprocessedOptions) > 0 {
		minIV, maxIV := preprocessedOptions[0].ImpliedVolatility, preprocessedOptions[0].ImpliedVolatility
		sumIV := 0.0
		for _, opt := range preprocessedOptions {
			if opt.ImpliedVolatility < minIV {
				minIV = opt.ImpliedVolatility
			}
			if opt.ImpliedVolatility > maxIV {
				maxIV = opt.ImpliedVolatility
			}
			sumIV += opt.ImpliedVolatility
		}
		avgIV := sumIV / float64(len(preprocessedOptions))
		log.Printf("Preprocessed data IV range: %.1f%% - %.1f%%, avg: %.1f%%", 
			minIV*100, maxIV*100, avgIV*100)
		
		// Show moneyness distribution
		moneynessCount := make(map[string]int)
		for _, opt := range preprocessedOptions {
			moneyness := opt.Strike / vsa.spotPrice
			switch {
			case moneyness < 0.8:
				moneynessCount["<0.8"]++
			case moneyness < 0.95:
				moneynessCount["0.8-0.95"]++
			case moneyness < 1.05:
				moneynessCount["0.95-1.05 (ATM)"]++
			case moneyness < 1.2:
				moneynessCount["1.05-1.2"]++
			default:
				moneynessCount[">1.2"]++
			}
		}
		log.Println("Moneyness distribution:")
		for k, v := range moneynessCount {
			log.Printf("  %s: %d options", k, v)
		}
	}
	
	// Try SSVI fitting first (global fit)
	log.Println("Attempting SSVI (Surface SVI) fitting...")
	ssviSurface, err := FitSSVISurface(preprocessedOptions, vsa.spotPrice, vsa.riskFreeRate)
	if err != nil {
		log.Printf("SSVI fitting failed: %v", err)
	} else {
		vsa.ssviSurface = ssviSurface
		log.Printf("✅ SSVI fitted successfully: Theta=%.4f, Rho=%.4f, Phi=%.4f",
			ssviSurface.Parameters.Theta, ssviSurface.Parameters.Rho, ssviSurface.Parameters.Phi)
	}
	
	// Also fit slice-by-slice SVI with improved constraints
	log.Println("Fitting slice-by-slice SVI with regularization...")
	surface, err := FitVolatilitySurfaceWithConstraints(preprocessedOptions, vsa.spotPrice, vsa.riskFreeRate)
	if err != nil {
		return fmt.Errorf("failed to fit volatility surface: %v", err)
	}
	
	vsa.surface = surface
	log.Printf("Fitted SVI surface with %d expiry slices", len(surface.Expiries))
	
	// Print surface parameters
	for i, exp := range surface.Expiries {
		params := surface.Parameters[i]
		log.Printf("  Expiry %.3f: A=%.4f, B=%.4f, Rho=%.4f, M=%.4f, Sigma=%.4f",
			exp, params.A, params.B, params.Rho, params.M, params.Sigma)
	}
	
	return nil
}

// CheckArbitrage runs arbitrage checks on the fitted surface
func (vsa *VolSurfaceAnalyzer) CheckArbitrage() []ArbitrageViolation {
	if vsa.surface == nil {
		log.Println("No fitted surface available for arbitrage checking")
		return []ArbitrageViolation{}
	}
	
	checker := NewArbitrageChecker(vsa.surface)
	violations := checker.CheckAllArbitrage()
	
	// Filter by severity threshold
	significantViolations := FilterViolationsBySeverity(violations, 0.01)
	
	PrintArbitrageReport(significantViolations)
	
	return significantViolations
}

// GetFittedIV returns the fitted implied volatility for any strike and expiry
func (vsa *VolSurfaceAnalyzer) GetFittedIV(strike, timeToExpiry float64) float64 {
	// Prefer SSVI if available (it's arbitrage-free by construction)
	if vsa.ssviSurface != nil {
		iv := vsa.ssviSurface.GetIV(strike, timeToExpiry)
		// Debug first few calls
		if strike == 1750.0 && timeToExpiry < 0.02 {
			log.Printf("SSVI GetIV: strike=%.0f, tte=%.3f -> IV=%.4f", strike, timeToExpiry, iv)
		}
		return iv
	}
	// Fall back to slice-by-slice SVI
	if vsa.surface != nil {
		return vsa.surface.GetIV(strike, timeToExpiry)
	}
	return 0
}

// GetSurface returns the fitted SVI surface
func (vsa *VolSurfaceAnalyzer) GetSurface() *SVISurface {
	return vsa.surface
}
