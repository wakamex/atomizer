package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/shopspring/decimal"
)

// OptionAnalysis provides analysis of option instruments
type OptionAnalysis struct {
	instruments map[string]DeriveInstrument
}

// TickerResult represents the API response for ticker data
type TickerResult struct {
	InstrumentName         string             `json:"instrument_name"`
	IsActive              bool               `json:"is_active"`
	MarkPrice             decimal.Decimal    `json:"mark_price"`
	BestBidPrice          decimal.Decimal    `json:"best_bid_price"`
	BestBidAmount         decimal.Decimal    `json:"best_bid_amount"`
	BestAskPrice          decimal.Decimal    `json:"best_ask_price"`
	BestAskAmount         decimal.Decimal    `json:"best_ask_amount"`
	IndexPrice            decimal.Decimal    `json:"index_price"`
	Timestamp             int64              `json:"timestamp"`
	Stats                 TickerStats        `json:"stats"`
	OptionPricing         OptionPricing      `json:"option_pricing"`
	OptionDetails         TickerOptionDetails `json:"option_details"`
	OpenInterest          map[string][]OIData `json:"open_interest"`
}

type TickerStats struct {
	ContractVolume decimal.Decimal `json:"contract_volume"`
	NumTrades      string          `json:"num_trades"`
	High           decimal.Decimal `json:"high"`
	Low            decimal.Decimal `json:"low"`
	PercentChange  decimal.Decimal `json:"percent_change"`
	UsdChange      decimal.Decimal `json:"usd_change"`
}

type OptionPricing struct {
	Delta           decimal.Decimal `json:"delta"`
	Theta           decimal.Decimal `json:"theta"`
	Gamma           decimal.Decimal `json:"gamma"`
	Vega            decimal.Decimal `json:"vega"`
	IV              decimal.Decimal `json:"iv"`
	Rho             decimal.Decimal `json:"rho"`
	MarkPrice       decimal.Decimal `json:"mark_price"`
	ForwardPrice    decimal.Decimal `json:"forward_price"`
	DiscountFactor  decimal.Decimal `json:"discount_factor"`
	BidIV           decimal.Decimal `json:"bid_iv"`
	AskIV           decimal.Decimal `json:"ask_iv"`
}

type TickerOptionDetails struct {
	Strike          string  `json:"strike"`
	Expiry          int64   `json:"expiry"`
	OptionType      string  `json:"option_type"`
	SettlementPrice decimal.Decimal `json:"settlement_price"`
}

type OIData struct {
	CurrentOpenInterest decimal.Decimal `json:"current_open_interest"`
	InterestCap        decimal.Decimal `json:"interest_cap"`
	ManagerCurrency    string          `json:"manager_currency"`
}

// CallAnalysisResult holds analysis results for a call option
type CallAnalysisResult struct {
	Ticker          TickerResult
	LiquidityScore  float64
	NormVolume      float64
	NormTrades      float64
	NormTotalOI     float64
	NormTopBook     float64
	NormSpread      float64
	IntrinsicValue  float64
	BidToIntrinsic  float64
	AskToIntrinsic  float64
	CustomBSIV      float64
	DaysToExpiry    float64
}

// NewOptionAnalysis creates a new analysis instance
func NewOptionAnalysis() *OptionAnalysis {
	return &OptionAnalysis{
		instruments: make(map[string]DeriveInstrument),
	}
}

// LoadInstruments fetches all option instruments
func (oa *OptionAnalysis) LoadInstruments() error {
	fmt.Println("Fetching all option instruments from Derive...")
	instruments, err := LoadAllDeriveMarkets()
	if err != nil {
		return fmt.Errorf("failed to load markets: %w", err)
	}
	oa.instruments = instruments
	fmt.Printf("Loaded %d option instruments\n", len(instruments))
	return nil
}

// AnalyzeByExpiry shows options grouped by expiry date
func (oa *OptionAnalysis) AnalyzeByExpiry() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("OPTIONS ANALYSIS BY EXPIRY")
	fmt.Println(strings.Repeat("=", 80))

	// Group by expiry date
	expiryMap := make(map[string][]DeriveInstrument)
	for _, inst := range oa.instruments {
		expiryDate := time.Unix(inst.OptionDetails.Expiry, 0).Format("2006-01-02")
		expiryMap[expiryDate] = append(expiryMap[expiryDate], inst)
	}

	// Sort expiry dates
	var expiries []string
	for expiry := range expiryMap {
		expiries = append(expiries, expiry)
	}
	sort.Strings(expiries)

	// Display analysis
	for _, expiry := range expiries {
		instruments := expiryMap[expiry]
		active := 0
		currencies := make(map[string]int)
		strikes := make(map[string]bool)
		optionTypes := make(map[string]int)

		for _, inst := range instruments {
			if inst.IsActive {
				active++
			}
			currencies[inst.BaseCurrency]++
			strikes[inst.OptionDetails.Strike] = true
			optionTypes[inst.OptionDetails.OptionType]++
		}

		fmt.Printf("\nExpiry: %s\n", expiry)
		fmt.Printf("  Total Options: %d\n", len(instruments))
		fmt.Printf("  Active Options: %d\n", active)
		fmt.Printf("  Currencies: %v\n", currencies)
		fmt.Printf("  Unique Strikes: %d\n", len(strikes))
		fmt.Printf("  Option Types: %v\n", optionTypes)
	}
}

// AnalyzeNearTermOptions shows options expiring in the next N days
func (oa *OptionAnalysis) AnalyzeNearTermOptions(days int) {
	fmt.Printf("\n" + strings.Repeat("=", 80))
	fmt.Printf("\nNEAR-TERM ACTIVE OPTIONS (Next %d days)\n", days)
	fmt.Println(strings.Repeat("=", 80))

	now := time.Now()
	cutoff := now.Add(time.Duration(days) * 24 * time.Hour)

	// Collect near-term active options
	nearTerm := make(map[string][]DeriveInstrument)
	for _, inst := range oa.instruments {
		expiry := time.Unix(inst.OptionDetails.Expiry, 0)
		if inst.IsActive && expiry.After(now) && expiry.Before(cutoff) {
			key := fmt.Sprintf("%s_%s", inst.BaseCurrency, expiry.Format("2006-01-02"))
			nearTerm[key] = append(nearTerm[key], inst)
		}
	}

	// Sort and display
	var keys []string
	for key := range nearTerm {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		instruments := nearTerm[key]
		parts := strings.Split(key, "_")
		currency := parts[0]
		expiryStr := parts[1]
		
		// Calculate days to expiry
		expiry, _ := time.Parse("2006-01-02", expiryStr)
		daysToExpiry := int(expiry.Sub(now).Hours() / 24)

		// Find strike range
		minStrike := decimal.NewFromFloat(999999)
		maxStrike := decimal.Zero
		for _, inst := range instruments {
			strike, _ := decimal.NewFromString(inst.OptionDetails.Strike)
			if strike.LessThan(minStrike) {
				minStrike = strike
			}
			if strike.GreaterThan(maxStrike) {
				maxStrike = strike
			}
		}

		fmt.Printf("%s %s (%dd): %d options, strikes $%s-$%s\n",
			currency, expiryStr, daysToExpiry, len(instruments),
			minStrike.StringFixed(0), maxStrike.StringFixed(0))
	}
}

// ExportETHCallsNearExpiry exports ETH call options expiring soon
func (oa *OptionAnalysis) ExportETHCallsNearExpiry(maxDays int) error {
	fmt.Printf("\n" + strings.Repeat("=", 80))
	fmt.Printf("\nEXPORTING ETH CALL OPTIONS (0-%d DAY EXPIRY)\n", maxDays)
	fmt.Println(strings.Repeat("=", 80))

	now := time.Now()
	cutoff := now.Add(time.Duration(maxDays) * 24 * time.Hour)

	// Filter ETH calls
	var ethCalls []DeriveInstrument
	for _, inst := range oa.instruments {
		expiry := time.Unix(inst.OptionDetails.Expiry, 0)
		if inst.BaseCurrency == "ETH" && 
		   inst.OptionDetails.OptionType == "C" &&
		   expiry.After(now) && expiry.Before(cutoff) {
			ethCalls = append(ethCalls, inst)
		}
	}

	// Sort by expiry and strike
	sort.Slice(ethCalls, func(i, j int) bool {
		if ethCalls[i].OptionDetails.Expiry != ethCalls[j].OptionDetails.Expiry {
			return ethCalls[i].OptionDetails.Expiry < ethCalls[j].OptionDetails.Expiry
		}
		strikeI, _ := decimal.NewFromString(ethCalls[i].OptionDetails.Strike)
		strikeJ, _ := decimal.NewFromString(ethCalls[j].OptionDetails.Strike)
		return strikeI.LessThan(strikeJ)
	})

	// Create CSV file
	filename := fmt.Sprintf("eth_calls_0_%d_day_%s.csv", maxDays, now.Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"instrument_name", "is_active", "expiry_date", "strike",
		"days_to_expiry", "hours_to_expiry",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write data
	activeCount := 0
	for _, inst := range ethCalls {
		expiry := time.Unix(inst.OptionDetails.Expiry, 0)
		hoursToExpiry := expiry.Sub(now).Hours()
		daysToExpiry := int(hoursToExpiry / 24)

		if inst.IsActive {
			activeCount++
		}

		record := []string{
			inst.InstrumentName,
			strconv.FormatBool(inst.IsActive),
			expiry.Format("2006-01-02 15:04:05"),
			inst.OptionDetails.Strike,
			strconv.Itoa(daysToExpiry),
			fmt.Sprintf("%.1f", hoursToExpiry),
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	fmt.Printf("\nFound %d ETH call options expiring in 0-%d days\n", len(ethCalls), maxDays)
	fmt.Printf("Saved to: %s\n", filename)

	// Show summary
	if len(ethCalls) > 0 {
		fmt.Println("\nSummary:")
		fmt.Printf("  Active: %d / %d\n", activeCount, len(ethCalls))
		
		// Group by expiry
		expiryGroups := make(map[string]int)
		for _, inst := range ethCalls {
			expiry := time.Unix(inst.OptionDetails.Expiry, 0)
			expiryStr := expiry.Format("2006-01-02 15:04")
			expiryGroups[expiryStr]++
		}

		fmt.Println("  Expiries:")
		var expiryKeys []string
		for k := range expiryGroups {
			expiryKeys = append(expiryKeys, k)
		}
		sort.Strings(expiryKeys)

		for _, expiryStr := range expiryKeys {
			expiry, _ := time.Parse("2006-01-02 15:04", expiryStr)
			hours := expiry.Sub(now).Hours()
			fmt.Printf("    %s (%.1f hours): %d options\n", 
				expiryStr, hours, expiryGroups[expiryStr])
		}
	}

	return nil
}

// ShowActiveOptionsStats shows statistics for active options
func (oa *OptionAnalysis) ShowActiveOptionsStats() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("ACTIVE OPTIONS ANALYSIS (Currently Tradeable)")
	fmt.Println(strings.Repeat("=", 80))

	total := len(oa.instruments)
	active := 0
	for _, inst := range oa.instruments {
		if inst.IsActive {
			active++
		}
	}

	fmt.Printf("\nTotal Options: %d\n", total)
	fmt.Printf("Active (Tradeable): %d (%.1f%%)\n", active, float64(active)/float64(total)*100)
	fmt.Printf("Inactive: %d (%.1f%%)\n", total-active, float64(total-active)/float64(total)*100)
}

// ShowActivePercentageByExpiry shows percentage of active options by expiry
func (oa *OptionAnalysis) ShowActivePercentageByExpiry() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("PERCENTAGE OF ACTIVE OPTIONS BY EXPIRY")
	fmt.Println(strings.Repeat("=", 80))

	type expiryStats struct {
		expiry    string
		active    int
		total     int
		pctActive float64
	}

	// Calculate stats
	expiryMap := make(map[string]*expiryStats)
	for _, inst := range oa.instruments {
		expiryDate := time.Unix(inst.OptionDetails.Expiry, 0).Format("2006-01-02")
		if _, exists := expiryMap[expiryDate]; !exists {
			expiryMap[expiryDate] = &expiryStats{expiry: expiryDate}
		}
		stats := expiryMap[expiryDate]
		stats.total++
		if inst.IsActive {
			stats.active++
		}
	}

	// Convert to slice and calculate percentages
	var expiries []expiryStats
	for _, stats := range expiryMap {
		stats.pctActive = float64(stats.active) / float64(stats.total) * 100
		expiries = append(expiries, *stats)
	}

	// Sort by percentage active (descending)
	sort.Slice(expiries, func(i, j int) bool {
		return expiries[i].pctActive > expiries[j].pctActive
	})

	// Display with visual indicators
	for _, stats := range expiries {
		status := "✗ INACTIVE"
		if stats.pctActive == 100 {
			status = "✓ ACTIVE"
		} else if stats.pctActive > 0 {
			status = "✗ MIXED"
		}
		
		fmt.Printf("%s: %d/%d (%.0f%%) %s\n",
			stats.expiry,
			stats.active,
			stats.total,
			stats.pctActive,
			status,
		)
	}
}

// ShowStrikeDistribution shows strike distribution statistics
func (oa *OptionAnalysis) ShowStrikeDistribution() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("STRIKE DISTRIBUTION FOR ACTIVE OPTIONS")
	fmt.Println(strings.Repeat("=", 80))

	// Group strikes by currency
	strikesByCurrency := make(map[string][]float64)
	for _, inst := range oa.instruments {
		if inst.IsActive {
			strike, _ := strconv.ParseFloat(inst.OptionDetails.Strike, 64)
			strikesByCurrency[inst.BaseCurrency] = append(
				strikesByCurrency[inst.BaseCurrency], strike)
		}
	}

	// Calculate and display stats for each currency
	var currencies []string
	for currency := range strikesByCurrency {
		currencies = append(currencies, currency)
	}
	sort.Strings(currencies)

	for _, currency := range currencies {
		strikes := strikesByCurrency[currency]
		if len(strikes) > 0 {
			sort.Float64s(strikes)
			
			// Calculate statistics
			sum := 0.0
			for _, s := range strikes {
				sum += s
			}
			mean := sum / float64(len(strikes))
			
			// Calculate std deviation
			sumSquaredDiff := 0.0
			for _, s := range strikes {
				diff := s - mean
				sumSquaredDiff += diff * diff
			}
			std := 0.0
			if len(strikes) > 1 {
				std = math.Sqrt(sumSquaredDiff / float64(len(strikes)))
			}
			
			// Calculate percentiles
			p25 := strikes[len(strikes)*25/100]
			p50 := strikes[len(strikes)*50/100]
			p75 := strikes[len(strikes)*75/100]
			
			fmt.Printf("\n%s:\n", currency)
			fmt.Printf("  Count: %d\n", len(strikes))
			fmt.Printf("  Mean: %.0f\n", mean)
			fmt.Printf("  Std: %.0f\n", std)
			fmt.Printf("  Min: %.0f\n", strikes[0])
			fmt.Printf("  25%%: %.0f\n", p25)
			fmt.Printf("  50%%: %.0f\n", p50)
			fmt.Printf("  75%%: %.0f\n", p75)
			fmt.Printf("  Max: %.0f\n", strikes[len(strikes)-1])
		}
	}
}

// FetchTicker fetches ticker data for a single instrument
func (oa *OptionAnalysis) FetchTicker(instrumentName string) (*TickerResult, error) {
	url := "https://api.lyra.finance/public/get_ticker"
	payload := map[string]string{"instrument_name": instrumentName}
	jsonData, _ := json.Marshal(payload)
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var result struct {
		Result TickerResult `json:"result"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}
	
	return &result.Result, nil
}

// normalizeMetric normalizes a slice of values between 0 and 1
func normalizeMetric(values []float64, higherIsBetter bool) []float64 {
	if len(values) == 0 {
		return values
	}
	
	minVal := values[0]
	maxVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}
	
	if minVal == maxVal {
		result := make([]float64, len(values))
		for i := range result {
			result[i] = 0.5
		}
		return result
	}
	
	normalized := make([]float64, len(values))
	for i, v := range values {
		if higherIsBetter {
			normalized[i] = (v - minVal) / (maxVal - minVal)
		} else {
			normalized[i] = 1.0 - (v - minVal) / (maxVal - minVal)
		}
	}
	
	return normalized
}

// CalculateLiquidityScores calculates liquidity scores for a set of results
func (oa *OptionAnalysis) CalculateLiquidityScores(results []CallAnalysisResult) {
	n := len(results)
	if n == 0 {
		return
	}
	
	// Extract metrics
	volumes := make([]float64, n)
	trades := make([]float64, n)
	totalOIs := make([]float64, n)
	topBookAmounts := make([]float64, n)
	spreads := make([]float64, n)
	
	for i, r := range results {
		volumes[i] = r.Ticker.Stats.ContractVolume.InexactFloat64()
		numTrades, _ := strconv.ParseFloat(r.Ticker.Stats.NumTrades, 64)
		trades[i] = numTrades
		
		// Calculate total OI
		totalOI := 0.0
		for _, oiList := range r.Ticker.OpenInterest {
			for _, oi := range oiList {
				totalOI += oi.CurrentOpenInterest.InexactFloat64()
			}
		}
		totalOIs[i] = totalOI
		
		// Top of book amounts
		topBookAmounts[i] = r.Ticker.BestBidAmount.Add(r.Ticker.BestAskAmount).InexactFloat64()
		
		// Bid-ask spread
		if r.Ticker.BestAskPrice.GreaterThan(decimal.Zero) && r.Ticker.BestBidPrice.GreaterThan(decimal.Zero) {
			spreads[i] = r.Ticker.BestAskPrice.Sub(r.Ticker.BestBidPrice).InexactFloat64()
		}
	}
	
	// Normalize metrics
	normVolumes := normalizeMetric(volumes, true)
	normTrades := normalizeMetric(trades, true)
	normOIs := normalizeMetric(totalOIs, true)
	normTopBooks := normalizeMetric(topBookAmounts, true)
	normSpreads := normalizeMetric(spreads, false)
	
	// Calculate liquidity scores
	for i := range results {
		results[i].NormVolume = normVolumes[i]
		results[i].NormTrades = normTrades[i]
		results[i].NormTotalOI = normOIs[i]
		results[i].NormTopBook = normTopBooks[i]
		results[i].NormSpread = normSpreads[i]
		
		// Multiplicative score
		results[i].LiquidityScore = normVolumes[i] * normTrades[i] * normOIs[i] * 
			normTopBooks[i] * normSpreads[i]
	}
}

// QueryETHCalls fetches and analyzes ETH call options expiring soon
func (oa *OptionAnalysis) QueryETHCalls(maxDays int) error {
	now := time.Now()
	
	// Find all unique expiries for ETH calls
	expiryMap := make(map[int64][]DeriveInstrument)
	for _, inst := range oa.instruments {
		if inst.BaseCurrency == "ETH" && 
		   inst.OptionDetails.OptionType == "C" &&
		   inst.OptionDetails.Expiry > now.Unix() {
			expiryMap[inst.OptionDetails.Expiry] = append(expiryMap[inst.OptionDetails.Expiry], inst)
		}
	}
	
	// Sort expiries to find nearest
	var expiries []int64
	for exp := range expiryMap {
		expiries = append(expiries, exp)
	}
	sort.Slice(expiries, func(i, j int) bool { return expiries[i] < expiries[j] })
	
	// Use the nearest expiry
	if len(expiries) == 0 {
		return fmt.Errorf("no ETH call options found")
	}
	
	nearestExpiry := expiries[0]
	nearestExpiryTime := time.Unix(nearestExpiry, 0)
	daysToExpiry := nearestExpiryTime.Sub(now).Hours() / 24
	
	fmt.Printf("\n" + strings.Repeat("=", 80))
	fmt.Printf("\nFETCHING ETH CALL OPTIONS - NEAREST EXPIRY: %s (%.1f days)\n", 
		nearestExpiryTime.Format("2006-01-02 15:04"), daysToExpiry)
	fmt.Println(strings.Repeat("=", 80))
	
	// Get ETH calls for nearest expiry
	ethCalls := expiryMap[nearestExpiry]
	
	if len(ethCalls) == 0 {
		fmt.Println("No ETH calls found in the specified time range")
		return nil
	}
	
	fmt.Printf("Fetching ticker data for %d ETH calls...\n", len(ethCalls))
	
	// Fetch ticker data with rate limiting
	results := make([]CallAnalysisResult, 0, len(ethCalls))
	for i, call := range ethCalls {
		fmt.Printf("\rFetching %d/%d: %s", i+1, len(ethCalls), call.InstrumentName)
		
		ticker, err := oa.FetchTicker(call.InstrumentName)
		if err != nil {
			fmt.Printf(" - Error: %v\n", err)
			continue
		}
		
		// Calculate days to expiry
		expiry := time.Unix(ticker.OptionDetails.Expiry, 0)
		daysToExpiry := expiry.Sub(now).Hours() / 24
		
		// Calculate intrinsic value
		spotPrice := ticker.IndexPrice.InexactFloat64()
		strike, _ := strconv.ParseFloat(ticker.OptionDetails.Strike, 64)
		intrinsicValue := math.Max(0, spotPrice - strike)
		
		result := CallAnalysisResult{
			Ticker:         *ticker,
			DaysToExpiry:   daysToExpiry,
			IntrinsicValue: intrinsicValue,
		}
		
		// Calculate bid/ask to intrinsic ratios
		if intrinsicValue > 0.001 {
			result.BidToIntrinsic = ticker.BestBidPrice.InexactFloat64() / intrinsicValue
			result.AskToIntrinsic = ticker.BestAskPrice.InexactFloat64() / intrinsicValue
		}
		
		results = append(results, result)
		
		// Rate limiting
		time.Sleep(50 * time.Millisecond)
	}
	
	fmt.Printf("\n\nSuccessfully fetched %d ticker results\n", len(results))
	
	// Calculate liquidity scores
	oa.CalculateLiquidityScores(results)
	
	// Sort by liquidity score for display
	sort.Slice(results, func(i, j int) bool {
		return results[i].LiquidityScore > results[j].LiquidityScore
	})
	
	// Display top 10 most liquid
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("TOP 10 MOST LIQUID ETH CALLS")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\n%-30s %8s %8s %8s %10s %6s %8s\n", 
		"Instrument", "Bid", "Ask", "Score", "Volume", "Trades", "Delta")
	fmt.Println(strings.Repeat("-", 80))
	
	for i := 0; i < 10 && i < len(results); i++ {
		r := results[i]
		fmt.Printf("%-30s %8.2f %8.2f %8.6f %10.0f %6d %8.4f\n",
			r.Ticker.InstrumentName,
			r.Ticker.BestBidPrice.InexactFloat64(),
			r.Ticker.BestAskPrice.InexactFloat64(),
			r.LiquidityScore,
			r.Ticker.Stats.ContractVolume.InexactFloat64(),
			func() int { n, _ := strconv.Atoi(r.Ticker.Stats.NumTrades); return n }(),
			r.Ticker.OptionPricing.Delta.InexactFloat64(),
		)
	}
	
	// Show liquidity score components
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("LIQUIDITY SCORE COMPONENTS (Top 10)")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\n%-30s %8s %8s %8s %8s %8s\n",
		"Instrument", "Volume", "Trades", "OI", "TopBook", "Spread")
	fmt.Println(strings.Repeat("-", 80))
	
	for i := 0; i < 10 && i < len(results); i++ {
		r := results[i]
		fmt.Printf("%-30s %8.4f %8.4f %8.4f %8.4f %8.4f\n",
			r.Ticker.InstrumentName,
			r.NormVolume,
			r.NormTrades,
			r.NormTotalOI,
			r.NormTopBook,
			r.NormSpread,
		)
	}
	
	// Show IV comparison - sort all results
	ivResults := make([]CallAnalysisResult, len(results))
	copy(ivResults, results)
	
	// Sort by: 
	// 1. Options with B/I < 1 first (sorted by ascending B/I ratio)
	// 2. Then all others (sorted by ascending bid IV)
	sort.Slice(ivResults, func(i, j int) bool {
		// Both have B/I < 1: sort by ascending B/I ratio
		if ivResults[i].BidToIntrinsic > 0 && ivResults[i].BidToIntrinsic < 1 &&
		   ivResults[j].BidToIntrinsic > 0 && ivResults[j].BidToIntrinsic < 1 {
			return ivResults[i].BidToIntrinsic < ivResults[j].BidToIntrinsic
		}
		
		// One has B/I < 1, the other doesn't: B/I < 1 comes first
		if ivResults[i].BidToIntrinsic > 0 && ivResults[i].BidToIntrinsic < 1 {
			return true
		}
		if ivResults[j].BidToIntrinsic > 0 && ivResults[j].BidToIntrinsic < 1 {
			return false
		}
		
		// Neither has B/I < 1: sort by bid IV
		return ivResults[i].Ticker.OptionPricing.BidIV.LessThan(ivResults[j].Ticker.OptionPricing.BidIV)
	})
	
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("IMPLIED VOLATILITY COMPARISON (Sorted by B/I Ratio, then Bid IV)")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("\n%-30s %8s %8s %8s %8s %8s %10s %8s\n",
		"Instrument", "Bid", "Ask", "API IV", "Bid IV", "Ask IV", "Intrinsic", "B/I Ratio")
	fmt.Println(strings.Repeat("-", 100))
	
	displayCount := 10
	if len(ivResults) < displayCount {
		displayCount = len(ivResults)
	}
	
	for i := 0; i < displayCount; i++ {
		r := ivResults[i]
		fmt.Printf("%-30s %8.2f %8.2f %8.4f %8.4f %8.4f %10.2f %8.2f\n",
			r.Ticker.InstrumentName,
			r.Ticker.BestBidPrice.InexactFloat64(),
			r.Ticker.BestAskPrice.InexactFloat64(),
			r.Ticker.OptionPricing.IV.InexactFloat64(),
			r.Ticker.OptionPricing.BidIV.InexactFloat64(),
			r.Ticker.OptionPricing.AskIV.InexactFloat64(),
			r.IntrinsicValue,
			r.BidToIntrinsic,
		)
	}
	
	// Export results to CSV
	filename := fmt.Sprintf("eth_calls_liquidity_%s.csv", now.Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV: %w", err)
	}
	defer file.Close()
	
	writer := csv.NewWriter(file)
	defer writer.Flush()
	
	// Write header
	header := []string{
		"instrument_name", "days_to_expiry", "liquidity_score",
		"bid_price", "ask_price", "spread", "volume", "num_trades",
		"delta", "gamma", "iv", "intrinsic_value", "bid_to_intrinsic",
	}
	writer.Write(header)
	
	// Write data
	for _, r := range results {
		record := []string{
			r.Ticker.InstrumentName,
			fmt.Sprintf("%.1f", r.DaysToExpiry),
			fmt.Sprintf("%.6f", r.LiquidityScore),
			fmt.Sprintf("%.2f", r.Ticker.BestBidPrice.InexactFloat64()),
			fmt.Sprintf("%.2f", r.Ticker.BestAskPrice.InexactFloat64()),
			fmt.Sprintf("%.2f", r.Ticker.BestAskPrice.Sub(r.Ticker.BestBidPrice).InexactFloat64()),
			fmt.Sprintf("%.0f", r.Ticker.Stats.ContractVolume.InexactFloat64()),
			r.Ticker.Stats.NumTrades,
			fmt.Sprintf("%.4f", r.Ticker.OptionPricing.Delta.InexactFloat64()),
			fmt.Sprintf("%.4f", r.Ticker.OptionPricing.Gamma.InexactFloat64()),
			fmt.Sprintf("%.4f", r.Ticker.OptionPricing.IV.InexactFloat64()),
			fmt.Sprintf("%.2f", r.IntrinsicValue),
			fmt.Sprintf("%.2f", r.BidToIntrinsic),
		}
		writer.Write(record)
	}
	
	fmt.Printf("\n\nResults saved to: %s\n", filename)
	
	return nil
}

// Main function for standalone execution
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Options Analysis Tool")
		fmt.Println("Usage: go run analyze_options.go [command]")
		fmt.Println("\nCommands:")
		fmt.Println("  all       - Run all analyses")
		fmt.Println("  expiry    - Analyze by expiry date")
		fmt.Println("  nearterm  - Show near-term options (30 days)")
		fmt.Println("  export    - Export ETH calls (0-1 day)")
		fmt.Println("  stats     - Show active options statistics + strike distribution")
		fmt.Println("  active    - Show active percentage by expiry with indicators")
		fmt.Println("  query     - Query ETH calls with liquidity analysis (nearest expiry)")
		return
	}

	analyzer := NewOptionAnalysis()
	if err := analyzer.LoadInstruments(); err != nil {
		log.Fatalf("Failed to load instruments: %v", err)
	}

	command := os.Args[1]
	switch command {
	case "all":
		analyzer.ShowActiveOptionsStats()
		analyzer.ShowActivePercentageByExpiry()
		analyzer.AnalyzeByExpiry()
		analyzer.ShowStrikeDistribution()
		analyzer.AnalyzeNearTermOptions(30)
		analyzer.ExportETHCallsNearExpiry(1)
	case "expiry":
		analyzer.AnalyzeByExpiry()
		analyzer.ShowActivePercentageByExpiry()
	case "nearterm":
		days := 30
		if len(os.Args) > 2 {
			if d, err := strconv.Atoi(os.Args[2]); err == nil {
				days = d
			}
		}
		analyzer.AnalyzeNearTermOptions(days)
	case "export":
		days := 1
		if len(os.Args) > 2 {
			if d, err := strconv.Atoi(os.Args[2]); err == nil {
				days = d
			}
		}
		if err := analyzer.ExportETHCallsNearExpiry(days); err != nil {
			log.Fatalf("Export failed: %v", err)
		}
	case "stats":
		analyzer.ShowActiveOptionsStats()
		analyzer.ShowStrikeDistribution()
	case "active":
		analyzer.ShowActivePercentageByExpiry()
	case "query":
		days := 1
		if len(os.Args) > 2 {
			if d, err := strconv.Atoi(os.Args[2]); err == nil {
				days = d
			}
		}
		if err := analyzer.QueryETHCalls(days); err != nil {
			log.Fatalf("Query failed: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}