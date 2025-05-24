package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	vol "vol_surface_web"
)

func main() {
	fmt.Println("Loading and analyzing options data...")
	
	// Load data
	options, err := vol.LoadOptionsFromCSV("../../eth_deribit_options.csv")
	if err != nil {
		log.Fatal("Error loading CSV:", err)
	}
	fmt.Printf("Loaded %d options from CSV\n", len(options))
	
	// Calculate IVs
	spotPrice := 1750.0
	riskFreeRate := 0.05
	
	fmt.Println("Calculating implied volatilities...")
	calculateIVs := 0
	for i := range options {
		if !math.IsNaN(options[i].MidETH) && options[i].MidETH > 0 {
			iv := vol.CalculateImpliedVolatility(options[i], spotPrice, riskFreeRate)
			if !math.IsNaN(iv) && iv > 0 {
				options[i].ImpliedVolatility = iv
				calculateIVs++
			}
		}
	}
	fmt.Printf("Calculated implied volatilities for %d options\n", calculateIVs)
	
	// Clean data
	fmt.Println("Cleaning option data...")
	cleaner := vol.NewDataCleaner()
	cleanedOptions := cleaner.CleanOptionData(options, spotPrice)
	cleaner.PrintDataQualityReport()
	
	// Create analyzer
	analyzer := vol.NewVolSurfaceAnalyzer(cleanedOptions, spotPrice, riskFreeRate)
	
	// Fit surfaces
	fmt.Println("Fitting volatility surfaces...")
	if err := analyzer.FitSurfaces(); err != nil {
		log.Printf("Warning: Error fitting surfaces: %v", err)
	}
	
	// Setup web server
	vol.SetupHandlers(analyzer)
	
	fmt.Println("\n🚀 Volatility Surface Web Server started on http://localhost:8080")
	fmt.Println("Open your browser to visualize the surface!")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}