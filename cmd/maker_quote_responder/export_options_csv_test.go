package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
	
	ccxt "github.com/ccxt/ccxt/go/v4"
)

func TestExportOptionsToCSV(t *testing.T) {
	exchange := ccxt.NewDeribit(map[string]interface{}{
		"rateLimit":       10,
		"enableRateLimit": true,
	})

	// Fetch ETH options
	ethTickers, err := exchange.FetchTickers(ccxt.WithFetchTickersParams(map[string]interface{}{
		"code": "ETH",
	}))
	if err != nil {
		t.Fatalf("Failed to fetch ETH tickers: %v", err)
	}

	// Get ETH index price
	perpTicker, err := exchange.FetchTicker("ETH-PERPETUAL")
	if err != nil {
		t.Fatalf("Failed to fetch ETH-PERPETUAL: %v", err)
	}
	
	var ethIndexPrice float64
	if perpTicker.Last != nil {
		ethIndexPrice = *perpTicker.Last
	}

	// Create CSV file
	file, err := os.Create("eth_options_data.csv")
	if err != nil {
		t.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Symbol",
		"Type",
		"Strike",
		"Expiry",
		"Days_To_Expiry",
		"Bid_ETH",
		"Ask_ETH",
		"Mid_ETH",
		"Last_ETH",
		"Spread_ETH",
		"Spread_Pct",
		"Bid_USD",
		"Ask_USD",
		"Mid_USD",
		"Volume",
		"Open_Interest",
		"IV",
		"Delta",
		"Gamma",
		"Theta",
		"Vega",
		"ETH_Index_Price",
		"Moneyness",
		"Time_To_Expiry_Years",
	}
	writer.Write(header)

	// Collect and sort options
	type OptionData struct {
		Symbol string
		Data   []string
	}
	
	options := []OptionData{}

	for symbol, ticker := range ethTickers.Tickers {
		// Remove prefix
		cleanSymbol := symbol
		if idx := strings.Index(symbol, ":"); idx >= 0 {
			cleanSymbol = symbol[idx+1:]
		}

		// Parse option details
		if !strings.HasSuffix(cleanSymbol, "-C") && !strings.HasSuffix(cleanSymbol, "-P") {
			continue
		}

		parts := strings.Split(cleanSymbol, "-")
		if len(parts) != 4 {
			continue
		}

		// Extract components
		// underlying := parts[0] // Not used in this version
		expiryStr := parts[1]
		strikeStr := parts[2]
		optionType := parts[3]

		// Parse strike
		strike, err := strconv.ParseFloat(strikeStr, 64)
		if err != nil {
			continue
		}

		// Parse expiry
		var expiryTime time.Time
		var daysToExpiry float64
		if len(expiryStr) == 6 {
			year, _ := strconv.Atoi("20" + expiryStr[0:2])
			month, _ := strconv.Atoi(expiryStr[2:4])
			day, _ := strconv.Atoi(expiryStr[4:6])
			expiryTime = time.Date(year, time.Month(month), day, 8, 0, 0, 0, time.UTC)
			daysToExpiry = expiryTime.Sub(time.Now()).Hours() / 24
		}

		// Extract price data
		var bid, ask, mid, last, spread, spreadPct float64
		var bidUSD, askUSD, midUSD float64
		
		if ticker.Bid != nil && ticker.Ask != nil {
			bid = *ticker.Bid
			ask = *ticker.Ask
			mid = (bid + ask) / 2
			spread = ask - bid
			if mid > 0 {
				spreadPct = (spread / mid) * 100
			}
			
			// Calculate USD values
			bidUSD = bid * ethIndexPrice
			askUSD = ask * ethIndexPrice
			midUSD = mid * ethIndexPrice
		}
		
		if ticker.Last != nil {
			last = *ticker.Last
		}

		// Extract volume and OI
		var volume, openInterest float64
		if ticker.BaseVolume != nil {
			volume = *ticker.BaseVolume
		}
		if ticker.Info != nil {
			if oi, exists := ticker.Info["open_interest"]; exists {
				if oiFloat, ok := oi.(float64); ok {
					openInterest = oiFloat
				}
			}
		}

		// Extract Greeks from info
		var iv, delta, gamma, theta, vega float64
		if ticker.Info != nil {
			// Try to extract Greeks
			if ivVal, exists := ticker.Info["mark_iv"]; exists {
				if ivFloat, ok := ivVal.(float64); ok {
					iv = ivFloat * 100 // Convert to percentage
				}
			}
			if deltaVal, exists := ticker.Info["greeks.delta"]; exists {
				if deltaFloat, ok := deltaVal.(float64); ok {
					delta = deltaFloat
				}
			}
			if gammaVal, exists := ticker.Info["greeks.gamma"]; exists {
				if gammaFloat, ok := gammaVal.(float64); ok {
					gamma = gammaFloat
				}
			}
			if thetaVal, exists := ticker.Info["greeks.theta"]; exists {
				if thetaFloat, ok := thetaVal.(float64); ok {
					theta = thetaFloat
				}
			}
			if vegaVal, exists := ticker.Info["greeks.vega"]; exists {
				if vegaFloat, ok := vegaVal.(float64); ok {
					vega = vegaFloat
				}
			}
		}

		// Calculate moneyness
		moneyness := ethIndexPrice / strike
		
		// Time to expiry in years
		timeToExpiryYears := daysToExpiry / 365.25

		// Create row
		row := []string{
			cleanSymbol,
			optionType,
			fmt.Sprintf("%.0f", strike),
			expiryTime.Format("2006-01-02"),
			fmt.Sprintf("%.1f", daysToExpiry),
			fmt.Sprintf("%.4f", bid),
			fmt.Sprintf("%.4f", ask),
			fmt.Sprintf("%.4f", mid),
			fmt.Sprintf("%.4f", last),
			fmt.Sprintf("%.4f", spread),
			fmt.Sprintf("%.2f", spreadPct),
			fmt.Sprintf("%.2f", bidUSD),
			fmt.Sprintf("%.2f", askUSD),
			fmt.Sprintf("%.2f", midUSD),
			fmt.Sprintf("%.0f", volume),
			fmt.Sprintf("%.0f", openInterest),
			fmt.Sprintf("%.1f", iv),
			fmt.Sprintf("%.4f", delta),
			fmt.Sprintf("%.4f", gamma),
			fmt.Sprintf("%.4f", theta),
			fmt.Sprintf("%.4f", vega),
			fmt.Sprintf("%.2f", ethIndexPrice),
			fmt.Sprintf("%.3f", moneyness),
			fmt.Sprintf("%.4f", timeToExpiryYears),
		}

		options = append(options, OptionData{
			Symbol: cleanSymbol,
			Data:   row,
		})
	}

	// Sort by symbol
	sort.Slice(options, func(i, j int) bool {
		return options[i].Symbol < options[j].Symbol
	})

	// Write data
	for _, opt := range options {
		writer.Write(opt.Data)
	}

	t.Logf("Exported %d options to eth_options_data.csv", len(options))
	t.Log("CSV file created successfully!")
	
	// Also create a summary
	t.Log("\nSummary Statistics:")
	t.Logf("ETH Index Price: $%.2f", ethIndexPrice)
	t.Logf("Total Options: %d", len(options))
	
	// Count by type
	callCount := 0
	putCount := 0
	for _, opt := range options {
		if opt.Data[1] == "C" {
			callCount++
		} else {
			putCount++
		}
	}
	t.Logf("Calls: %d, Puts: %d", callCount, putCount)
	
	// Sample data for verification
	t.Log("\nSample data (first 5 options):")
	for i := 0; i < 5 && i < len(options); i++ {
		t.Logf("%s: Bid=%.4f, Ask=%.4f, Mid=$%.2f", 
			options[i].Data[0], // Symbol
			parseFloat(options[i].Data[5]), // Bid ETH
			parseFloat(options[i].Data[6]), // Ask ETH
			parseFloat(options[i].Data[13])) // Mid USD
	}
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}