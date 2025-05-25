package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type MarketsResponse map[string][]Asset

type Asset struct {
	Symbol                 string   `json:"symbol"`
	Address                string   `json:"address"`
	Decimals               int      `json:"decimals"`
	ChainID                int      `json:"chainId"`
	Active                 bool     `json:"active"`
	Price                  string   `json:"price"`
	Underlying             string   `json:"underlying"`
	UnderlyingAssetAddress string   `json:"underlyingAssetAddress"`
	MinTradeSize           string   `json:"minTradeSize"`
	MaxTradeSize           string   `json:"maxTradeSize"`
}

func main() {
	var (
		apiURL     = flag.String("url", "https://rip-testnet.rysk.finance/api/assets", "API URL to fetch markets")
		outputFile = flag.String("output", "markets.json", "Output JSON file")
		pretty     = flag.Bool("pretty", true, "Pretty print JSON output")
	)
	flag.Parse()

	fmt.Printf("📊 Fetching markets data from: %s\n", *apiURL)

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the request
	resp, err := client.Get(*apiURL)
	if err != nil {
		log.Fatalf("Failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("API returned status %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	// Parse into generic structure first to see what we get
	var rawData map[string]interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Save the raw data
	var outputData []byte
	if *pretty {
		outputData, err = json.MarshalIndent(rawData, "", "  ")
	} else {
		outputData, err = json.Marshal(rawData)
	}
	
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
	}

	// Write to file
	if err := os.WriteFile(*outputFile, outputData, 0644); err != nil {
		log.Fatalf("Failed to write file: %v", err)
	}

	fmt.Printf("✅ Markets data saved to: %s\n", *outputFile)

	// Print summary
	fmt.Println("\n📈 MARKETS SUMMARY:")
	fmt.Println("================================================================================")
	
	for chainID, data := range rawData {
		if assets, ok := data.([]interface{}); ok {
			fmt.Printf("\nChain ID %s: %d assets\n", chainID, len(assets))
			
			// Try to extract some details about the first few assets
			for i, asset := range assets {
				if i >= 3 {
					fmt.Println("  ...")
					break
				}
				if assetMap, ok := asset.(map[string]interface{}); ok {
					symbol := "?"
					address := "?"
					underlying := "?"
					active := false
					
					if s, ok := assetMap["symbol"].(string); ok {
						symbol = s
					}
					if addr, ok := assetMap["address"].(string); ok {
						address = addr
					}
					if u, ok := assetMap["underlying"].(string); ok {
						underlying = u
					}
					if a, ok := assetMap["active"].(bool); ok {
						active = a
					}
					
					status := "INACTIVE"
					if active {
						status = "ACTIVE"
					}
					
					fmt.Printf("  - %-8s (%s): %s [%s]\n", symbol, underlying, address, status)
				}
			}
		}
	}
	
	fmt.Println("================================================================================")
	fmt.Printf("Total chains: %d\n", len(rawData))
	fmt.Printf("Output file: %s (%.2f KB)\n", *outputFile, float64(len(outputData))/1024.0)
}