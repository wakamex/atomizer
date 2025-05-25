package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	ID      string      `json:"id"`
	Params  interface{} `json:"params"`
}

type RFQParams struct {
	Asset     string `json:"asset"`
	AssetName string `json:"assetName"`
	ChainID   int    `json:"chainId"`
	Expiry    int64  `json:"expiry"`
	IsPut     bool   `json:"isPut"`
	Quantity  string `json:"quantity"`
	Strike    string `json:"strike"`
	Taker     string `json:"taker"`
}

type QuoteResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Result  QuoteResult `json:"result"`
}

type QuoteResult struct {
	AssetAddress string  `json:"assetAddress"`
	ChainID      int     `json:"chainId"`
	IsPut        bool    `json:"isPut"`
	Strike       string  `json:"strike"`
	Expiry       int64   `json:"expiry"`
	Maker        string  `json:"maker"`
	Nonce        string  `json:"nonce"`
	Price        string  `json:"price"`
	Quantity     string  `json:"quantity"`
	IsTakerBuy   bool    `json:"isTakerBuy"`
	Signature    string  `json:"signature"`
	ValidUntil   int64   `json:"validUntil"`
	APY          float64 `json:"apy"`
}

type InventoryResponse struct {
	JSONRPC string                       `json:"jsonrpc"`
	ID      string                       `json:"id"`
	Result  map[string]InventoryAsset    `json:"result"`
}

type InventoryAsset struct {
	Strikes      map[string][]float64              `json:"strikes"`
	Expiries     []int64                           `json:"expiries"`
	Combinations map[string]InventoryCombination   `json:"combinations"`
}

type InventoryCombination struct {
	Expiry               string  `json:"expiry"`
	ExpirationTimestamp  int64   `json:"expiration_timestamp"`
	TimeToExpiryDays     float64 `json:"timeToExpiryDays"`
	Strike               float64 `json:"strike"`
	IsPut                bool    `json:"isPut"`
	Delta                float64 `json:"delta"`
	Bid                  float64 `json:"bid"`
	Ask                  float64 `json:"ask"`
	Index                float64 `json:"index"`
	APY                  float64 `json:"apy"`
	Timestamp            int64   `json:"timestamp"`
}

type RFQTiming struct {
	Asset      string
	Strike     float64
	Expiry     string
	SentTime   time.Time
	ResponseTime time.Duration
	Responded  bool
}

type Market struct {
	Symbol                 string `json:"symbol"`
	Address                string `json:"address"`
	Decimals               int    `json:"decimals"`
	ChainID                int    `json:"chainId"`
	Active                 bool   `json:"active"`
	Price                  string `json:"price"`
	Underlying             string `json:"underlying"`
	UnderlyingAssetAddress string `json:"underlyingAssetAddress"`
	MinTradeSize           string `json:"minTradeSize"`
	MaxTradeSize           string `json:"maxTradeSize"`
}

func main() {
	var (
		wsURL       = flag.String("url", "wss://rip-testnet.rysk.finance/taker", "WebSocket URL")
		chainID     = flag.Int("chainId", 84532, "Chain ID")
		quantity    = flag.String("quantity", "1000000000000000000", "Quantity in wei")
		taker       = flag.String("taker", "0x0000000000000000000000000000000000000000", "Taker address")
		marketsFile = flag.String("markets", "markets.json", "Path to markets.json file")
	)
	flag.Parse()

	rfqTimings := make(map[string]*RFQTiming)
	var timingsMutex sync.Mutex

	u, err := url.Parse(*wsURL)
	if err != nil {
		log.Fatal("Invalid URL:", err)
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("WebSocket dial error:", err)
	}
	defer c.Close()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan struct{})

	inventoryReceived := make(chan InventoryResponse)
	
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}

			var response map[string]interface{}
			if err := json.Unmarshal(message, &response); err != nil {
				log.Println("JSON parse error:", err)
				continue
			}

			if idStr, ok := response["id"].(string); ok && strings.HasPrefix(idStr, "rfq-") && response["result"] != nil {
				timingsMutex.Lock()
				if timing, exists := rfqTimings[idStr]; exists {
					timing.ResponseTime = time.Since(timing.SentTime)
					timing.Responded = true
				}
				timingsMutex.Unlock()
			} else if response["id"] == "inventory" && response["result"] != nil {
				var inventory InventoryResponse
				if err := json.Unmarshal(message, &inventory); err == nil {
					inventoryReceived <- inventory
				}
			}
		}
	}()

	// Load markets data from file
	fmt.Printf("📊 Loading markets from %s...\n", *marketsFile)
	
	marketsData, err := loadMarketsData(*marketsFile)
	if err != nil {
		log.Fatalf("Failed to load markets data: %v", err)
	}
	
	// Get markets for the specified chain
	markets, exists := marketsData[strconv.Itoa(*chainID)]
	if !exists {
		log.Fatalf("No markets found for chain ID %d", *chainID)
	}
	
	fmt.Printf("✅ Loaded %d markets for chain %d\n", len(markets), *chainID)
	
	// First, send inventory request to get available strikes and expiries
	fmt.Println("\n📊 Fetching current inventory for strike/expiry data...")
	inventoryReq := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "inventory",
		ID:      "inventory",
	}

	inventoryBytes, _ := json.Marshal(inventoryReq)
	c.WriteMessage(websocket.TextMessage, inventoryBytes)

	// Wait for inventory response
	var inventory InventoryResponse
	select {
	case inventory = <-inventoryReceived:
		fmt.Println("✅ Inventory received")
	case <-time.After(5 * time.Second):
		log.Fatal("Timeout waiting for inventory response")
	}

	// Send RFQs to all available markets
	fmt.Println("\n🚀 Sending RFQs to all available markets simultaneously...")
	rfqCount := 0
	
	// Create a map of underlying to asset addresses from markets data
	underlyingToAddress := make(map[string]string)
	for _, market := range markets {
		if market.Active && market.Symbol == "W" + market.Underlying {
			underlyingToAddress[market.Underlying] = market.Address
		}
	}
	
	// Also check for exact matches
	for _, market := range markets {
		if market.Active {
			underlyingToAddress[market.Symbol] = market.Address
		}
	}

	for asset, data := range inventory.Result {
		assetAddr, hasAddr := underlyingToAddress[asset]
		if !hasAddr {
			fmt.Printf("⚠️  No active market found for %s, skipping\n", asset)
			continue
		}
		
		for _, combo := range data.Combinations {
			if combo.Bid > 0 || combo.Ask > 0 { // Only send RFQ if market is active
				rfqID := fmt.Sprintf("rfq-%s-%.0f-%d", asset, combo.Strike, combo.ExpirationTimestamp)
				
				rfq := JSONRPCRequest{
					JSONRPC: "2.0",
					Method:  "request",
					ID:      rfqID,
					Params: RFQParams{
						Asset:     assetAddr,
						AssetName: asset,
						ChainID:   *chainID,
						Expiry:    combo.ExpirationTimestamp,
						IsPut:     combo.IsPut,
						Quantity:  *quantity,
						Strike:    strconv.FormatFloat(combo.Strike * 1e8, 'f', 0, 64), // Convert to wei format
						Taker:     *taker,
					},
				}

				timingsMutex.Lock()
				rfqTimings[rfqID] = &RFQTiming{
					Asset:    asset,
					Strike:   combo.Strike,
					Expiry:   combo.Expiry,
					SentTime: time.Now(),
				}
				timingsMutex.Unlock()

				rfqBytes, _ := json.Marshal(rfq)
				c.WriteMessage(websocket.TextMessage, rfqBytes)
				rfqCount++
				
				// Small delay to avoid overwhelming the server
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	fmt.Printf("\n📤 Sent %d RFQs, waiting for responses...\n", rfqCount)
	
	// Wait for responses
	time.Sleep(3 * time.Second)

	// Print summary
	printResponseTimeSummary(rfqTimings)

	select {
	case <-interrupt:
		log.Println("Interrupt received, closing connection...")

		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("Write close error:", err)
		}

		select {
		case <-done:
		case <-time.After(time.Second):
		}
	}
}

func printResponseTimeSummary(timings map[string]*RFQTiming) {
	fmt.Println("\n📈 RESPONSE TIME SUMMARY:")
	fmt.Println(strings.Repeat("=", 80))
	
	// Group by asset
	assetTimings := make(map[string][]time.Duration)
	totalResponded := 0
	
	for _, timing := range timings {
		if timing.Responded {
			assetTimings[timing.Asset] = append(assetTimings[timing.Asset], timing.ResponseTime)
			totalResponded++
		}
	}
	
	// Sort assets for consistent output
	var assets []string
	for asset := range assetTimings {
		assets = append(assets, asset)
	}
	sort.Strings(assets)
	
	// Print per-asset statistics
	var allTimings []time.Duration
	for _, asset := range assets {
		timingList := assetTimings[asset]
		allTimings = append(allTimings, timingList...)
		
		min, max, avg := calculateStats(timingList)
		fmt.Printf("\n%s Markets (%d responses):\n", asset, len(timingList))
		fmt.Printf("  Min: %v\n", min)
		fmt.Printf("  Max: %v\n", max)
		fmt.Printf("  Avg: %v\n", avg)
	}
	
	// Overall statistics
	if len(allTimings) > 0 {
		min, max, avg := calculateStats(allTimings)
		fmt.Printf("\nOVERALL (%d/%d responded):\n", totalResponded, len(timings))
		fmt.Printf("  Min: %v\n", min)
		fmt.Printf("  Max: %v\n", max)
		fmt.Printf("  Avg: %v\n", avg)
	} else {
		fmt.Println("\n❌ No responses received")
	}
	
	fmt.Println(strings.Repeat("=", 80))
}

func calculateStats(timings []time.Duration) (min, max, avg time.Duration) {
	if len(timings) == 0 {
		return
	}
	
	min = timings[0]
	max = timings[0]
	var sum time.Duration
	
	for _, t := range timings {
		if t < min {
			min = t
		}
		if t > max {
			max = t
		}
		sum += t
	}
	
	avg = sum / time.Duration(len(timings))
	return
}

func loadMarketsData(filename string) (map[string][]Market, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open markets file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read markets file: %w", err)
	}

	var markets map[string][]Market
	if err := json.Unmarshal(data, &markets); err != nil {
		return nil, fmt.Errorf("failed to parse markets JSON: %w", err)
	}

	return markets, nil
}