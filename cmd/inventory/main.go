package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"time"

	"github.com/gorilla/websocket"
)

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	ID      string      `json:"id"`
	Params  interface{} `json:"params"`
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

func main() {
	var (
		wsURL = flag.String("url", "wss://rip-testnet.rysk.finance/taker", "WebSocket URL")
	)
	flag.Parse()

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

			if response["id"] == "inventory" && response["result"] != nil {
				var inventory InventoryResponse
				if err := json.Unmarshal(message, &inventory); err == nil {
					inventoryReceived <- inventory
				}
			}
		}
	}()

	// Send inventory request
	inventoryReq := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "inventory",
		ID:      "inventory",
	}

	inventoryBytes, _ := json.Marshal(inventoryReq)
	fmt.Printf("📤 Requesting inventory from: %s\n", *wsURL)
	c.WriteMessage(websocket.TextMessage, inventoryBytes)

	// Wait for inventory response
	select {
	case inventory := <-inventoryReceived:
		displayInventory(inventory)
	case <-time.After(5 * time.Second):
		log.Fatal("Timeout waiting for inventory response")
	case <-interrupt:
		log.Println("Interrupt received")
	}

	// Close connection
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("Write close error:", err)
	}

	select {
	case <-done:
	case <-time.After(time.Second):
	}
}

func displayInventory(inventory InventoryResponse) {
	fmt.Println("\n📊 MARKET INVENTORY")
	fmt.Println("================================================================================")

	// Sort assets for consistent display
	var assets []string
	for asset := range inventory.Result {
		assets = append(assets, asset)
	}
	sort.Strings(assets)

	totalMarkets := 0

	for _, asset := range assets {
		data := inventory.Result[asset]
		fmt.Printf("\n%s Markets:\n", asset)
		fmt.Println("----------------------------------------")

		// Group by expiry
		expiryMarkets := make(map[int64][]InventoryCombination)
		for _, combo := range data.Combinations {
			expiryMarkets[combo.ExpirationTimestamp] = append(expiryMarkets[combo.ExpirationTimestamp], combo)
		}

		// Sort expiries
		var expiries []int64
		for exp := range expiryMarkets {
			expiries = append(expiries, exp)
		}
		sort.Slice(expiries, func(i, j int) bool { return expiries[i] < expiries[j] })

		for _, expiry := range expiries {
			combos := expiryMarkets[expiry]
			if len(combos) == 0 {
				continue
			}

			// Sort by strike
			sort.Slice(combos, func(i, j int) bool { return combos[i].Strike < combos[j].Strike })

			expiryTime := time.Unix(expiry, 0)
			fmt.Printf("\n  Expiry: %s (%s)\n", combos[0].Expiry, expiryTime.Format("2006-01-02"))
			fmt.Printf("  Time to expiry: %.1f days\n", combos[0].TimeToExpiryDays)
			fmt.Printf("  Strikes:\n")

			for _, combo := range combos {
				optionType := "CALL"
				if combo.IsPut {
					optionType = "PUT"
				}

				status := "INACTIVE"
				spread := "-"
				if combo.Bid > 0 || combo.Ask > 0 {
					status = "ACTIVE"
					if combo.Bid > 0 && combo.Ask > 0 {
						spreadPct := ((combo.Ask - combo.Bid) / combo.Ask) * 100
						spread = fmt.Sprintf("%.1f%%", spreadPct)
					}
				}

				fmt.Printf("    %.0f %s: ", combo.Strike, optionType)
				if status == "ACTIVE" {
					fmt.Printf("Bid: %.6f, Ask: %.6f, Spread: %s, Delta: %.3f, APY: %.1f%% [%s]\n",
						combo.Bid, combo.Ask, spread, combo.Delta, combo.APY, status)
					totalMarkets++
				} else {
					fmt.Printf("[%s]\n", status)
				}
			}
		}
	}

	fmt.Printf("\n================================================================================\n")
	fmt.Printf("Total active markets: %d\n", totalMarkets)
	fmt.Printf("Last update: %s\n", time.Now().Format("2006-01-02 15:04:05"))
}