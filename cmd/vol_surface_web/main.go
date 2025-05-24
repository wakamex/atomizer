package main

import (
	"fmt"
	"log"
	"math"
	"os"
)

func testSurfaceRange(analyzer *VolSurfaceAnalyzer, testName string) {
	// Test the surface at various points
	minStrike, maxStrike := 1750.0, 7000.0
	minTime, maxTime := 0.01, 1.0
	gridSize := 20
	
	minIV, maxIV := 999.0, 0.0
	
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			strike := minStrike + float64(j)*(maxStrike-minStrike)/float64(gridSize-1)
			ttm := minTime + float64(i)*(maxTime-minTime)/float64(gridSize-1)
			
			iv := analyzer.GetFittedIV(strike, ttm)
			if iv > 0 {
				ivPct := iv * 100
				if ivPct < minIV {
					minIV = ivPct
				}
				if ivPct > maxIV {
					maxIV = ivPct
				}
			}
		}
	}
	
	fmt.Printf("%s: IV range %.1f%% - %.1f%% (spread: %.1f%%)\n", 
		testName, minIV, maxIV, maxIV-minIV)
	
	// Also print SSVI params if available
	if analyzer.ssviSurface != nil {
		p := analyzer.ssviSurface.Parameters
		fmt.Printf("  SSVI params: Theta=%.4f, Rho=%.4f, Phi=%.4f\n", p.Theta, p.Rho, p.Phi)
	}
}

func main() {
	log.SetFlags(0) // Remove timestamps for cleaner output
	
	// Check if we should run the web server
	if len(os.Args) > 1 && os.Args[1] == "server" {
		runWebServer()
		return
	}
	
	// Otherwise run tests
	fmt.Println("=== Testing different fixes for flat volatility surface ===\n")
	
	// Test 1: Original implementation
	fmt.Println("Test 1: Original implementation")
	analyzer1 := NewVolSurfaceAnalyzer()
	if err := analyzer1.LoadCSV("eth_options_data.csv"); err != nil {
		log.Fatal(err)
	}
	analyzer1.spotPrice = 3500.0
	analyzer1.CalculateImpliedVolatilities()
	analyzer1.CleanData()
	analyzer1.FitSVISurface()
	testSurfaceRange(analyzer1, "Original")
	
	// Test 2: Skip PreprocessOptionsForSVI (use all cleaned data)
	fmt.Println("\nTest 2: Skip aggressive preprocessing")
	analyzer2 := NewVolSurfaceAnalyzer()
	analyzer2.LoadCSV("eth_options_data.csv")
	analyzer2.spotPrice = 3500.0
	analyzer2.CalculateImpliedVolatilities()
	cleanedOpts, _ := analyzer2.CleanData()
	
	// Directly fit SSVI without extra preprocessing
	fmt.Printf("  Using %d options (vs %d after preprocessing)\n", len(cleanedOpts), 92)
	ssviSurface2, _ := FitSSVISurface(cleanedOpts, analyzer2.spotPrice, analyzer2.riskFreeRate)
	analyzer2.ssviSurface = ssviSurface2
	testSurfaceRange(analyzer2, "No preprocessing")
	
	// Test 3: Better initial SSVI parameters
	fmt.Println("\nTest 3: Better SSVI initial parameters (higher Phi)")
	// Temporarily modify the initial params in FitSSVISurface
	// We'll need to create a modified version
	analyzer3 := NewVolSurfaceAnalyzer()
	analyzer3.LoadCSV("eth_options_data.csv")
	analyzer3.spotPrice = 3500.0
	analyzer3.CalculateImpliedVolatilities()
	analyzer3.CleanData()
	
	// Custom SSVI fit with better initial params
	preprocessedOptions3 := PreprocessOptionsForSVI(analyzer3.options, analyzer3.spotPrice)
	ssviSurface3, _ := FitSSVISurfaceWithInitialParams(preprocessedOptions3, analyzer3.spotPrice, 
		analyzer3.riskFreeRate, 0.8) // Higher initial Phi
	analyzer3.ssviSurface = ssviSurface3
	testSurfaceRange(analyzer3, "High initial Phi=0.8")
	
	// Test 4: Minimum Phi constraint
	fmt.Println("\nTest 4: Minimum Phi constraint")
	analyzer4 := NewVolSurfaceAnalyzer()
	analyzer4.LoadCSV("eth_options_data.csv")
	analyzer4.spotPrice = 3500.0
	analyzer4.CalculateImpliedVolatilities()
	analyzer4.CleanData()
	
	preprocessedOptions4 := PreprocessOptionsForSVI(analyzer4.options, analyzer4.spotPrice)
	ssviSurface4, _ := FitSSVISurfaceWithMinPhi(preprocessedOptions4, analyzer4.spotPrice, 
		analyzer4.riskFreeRate, 0.3) // Min Phi = 0.3
	analyzer4.ssviSurface = ssviSurface4
	testSurfaceRange(analyzer4, "Min Phi=0.3")
	
	// Test 5: Use only SVI (not SSVI)
	fmt.Println("\nTest 5: SVI only (no SSVI)")
	analyzer5 := NewVolSurfaceAnalyzer()
	analyzer5.LoadCSV("eth_options_data.csv")
	analyzer5.spotPrice = 3500.0
	analyzer5.CalculateImpliedVolatilities()
	analyzer5.CleanData()
	analyzer5.FitSVISurface()
	analyzer5.ssviSurface = nil // Force use of SVI
	testSurfaceRange(analyzer5, "SVI only")
	
	// Test 6: Check data distribution
	fmt.Println("\nTest 6: Data distribution analysis")
	checkDataDistribution(analyzer1.options)
}

func checkDataDistribution(options []OptionData) {
	// Count options by moneyness buckets
	buckets := make(map[string]int)
	spotPrice := 3500.0
	
	for _, opt := range options {
		if math.IsNaN(opt.ImpliedVolatility) || opt.ImpliedVolatility <= 0 {
			continue
		}
		
		moneyness := opt.Strike / spotPrice
		var bucket string
		switch {
		case moneyness < 0.7:
			bucket = "Deep OTM Put (<0.7)"
		case moneyness < 0.9:
			bucket = "OTM Put (0.7-0.9)"
		case moneyness < 1.1:
			bucket = "ATM (0.9-1.1)"
		case moneyness < 1.3:
			bucket = "OTM Call (1.1-1.3)"
		default:
			bucket = "Deep OTM Call (>1.3)"
		}
		buckets[bucket]++
	}
	
	fmt.Println("  Strike distribution:")
	for bucket, count := range buckets {
		fmt.Printf("    %s: %d options\n", bucket, count)
	}
}

// Modified FitSSVISurface with configurable initial Phi
func FitSSVISurfaceWithInitialParams(options []OptionData, spotPrice, riskFreeRate, initialPhi float64) (*SSVISurface, error) {
	validOptions := []OptionData{}
	for _, opt := range options {
		if !math.IsNaN(opt.ImpliedVolatility) && opt.ImpliedVolatility > 0 && 
		   opt.TimeToExpiry > 0 && opt.Moneyness > 0.5 && opt.Moneyness < 2.0 {
			validOptions = append(validOptions, opt)
		}
	}
	
	if len(validOptions) < 10 {
		return nil, fmt.Errorf("insufficient valid options: %d", len(validOptions))
	}
	
	atmVol := estimateATMVol(validOptions)
	params := SSVIParameters{
		Theta: atmVol * atmVol,
		Rho:   -0.3,
		Phi:   initialPhi, // Use provided initial Phi
	}
	
	// Simple optimization (abbreviated for testing)
	learningRate := 0.01
	for iter := 0; iter < 100; iter++ {
		currentError, gradTheta, gradRho, gradPhi := calculateSSVIGradients(
			params, validOptions, spotPrice, riskFreeRate)
		
		params.Theta -= learningRate * gradTheta
		params.Rho -= learningRate * gradRho
		params.Phi -= learningRate * gradPhi
		
		// Apply constraints
		params.Theta = math.Max(0.01, math.Min(10.0, params.Theta))
		params.Rho = math.Max(-0.95, math.Min(0.95, params.Rho))
		params.Phi = math.Max(0.05, math.Min(2.0, params.Phi))
		
		if iter > 0 && currentError < 1e-6 {
			break
		}
	}
	
	return &SSVISurface{
		Parameters:   params,
		SpotPrice:    spotPrice,
		RiskFreeRate: riskFreeRate,
	}, nil
}

// Modified FitSSVISurface with minimum Phi constraint
func FitSSVISurfaceWithMinPhi(options []OptionData, spotPrice, riskFreeRate, minPhi float64) (*SSVISurface, error) {
	validOptions := []OptionData{}
	for _, opt := range options {
		if !math.IsNaN(opt.ImpliedVolatility) && opt.ImpliedVolatility > 0 && 
		   opt.TimeToExpiry > 0 && opt.Moneyness > 0.5 && opt.Moneyness < 2.0 {
			validOptions = append(validOptions, opt)
		}
	}
	
	if len(validOptions) < 10 {
		return nil, fmt.Errorf("insufficient valid options: %d", len(validOptions))
	}
	
	atmVol := estimateATMVol(validOptions)
	params := SSVIParameters{
		Theta: atmVol * atmVol,
		Rho:   -0.3,
		Phi:   math.Max(minPhi, 0.3),
	}
	
	// Simple optimization with minimum Phi constraint
	learningRate := 0.01
	for iter := 0; iter < 100; iter++ {
		currentError, gradTheta, gradRho, gradPhi := calculateSSVIGradients(
			params, validOptions, spotPrice, riskFreeRate)
		
		params.Theta -= learningRate * gradTheta
		params.Rho -= learningRate * gradRho
		params.Phi -= learningRate * gradPhi
		
		// Apply constraints with higher minimum Phi
		params.Theta = math.Max(0.01, math.Min(10.0, params.Theta))
		params.Rho = math.Max(-0.95, math.Min(0.95, params.Rho))
		params.Phi = math.Max(minPhi, math.Min(2.0, params.Phi)) // Enforce minimum Phi
		
		if iter > 0 && currentError < 1e-6 {
			break
		}
	}
	
	return &SSVISurface{
		Parameters:   params,
		SpotPrice:    spotPrice,
		RiskFreeRate: riskFreeRate,
	}, nil
}
func runWebServer() {
	// Create and start the web server
	ws := NewWebServer()
	
	fmt.Println("Loading and analyzing options data...")
	if err := ws.LoadAndAnalyze(); err != nil {
		log.Fatal("Error loading data:", err)
	}
	
	// Start the server on port 8080
	if err := ws.Start(8080); err != nil {
		log.Fatal("Server error:", err)
	}
}
