package main

import (
    "encoding/json"
    "fmt"
    "time"
    "github.com/gorilla/websocket"
)

func main() {
    conn, _, _ := websocket.DefaultDialer.Dial("wss://api.lyra.finance/ws", nil)
    defer conn.Close()
    
    // Subscribe
    conn.WriteJSON(map[string]interface{}{
        "method": "subscribe",
        "params": map[string]interface{}{
            "channels": []string{"orderbook.ETH-20250603-2600-C.1.10"},
        },
    })
    
    var lastPrices []string
    lineCount := 0
    
    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Error:", err)
            break
        }
        
        // Parse message
        var data map[string]interface{}
        if err := json.Unmarshal(msg, &data); err != nil {
            continue
        }
        
        // Check if it's a subscription update
        if method, ok := data["method"].(string); ok && method == "subscription" {
            if params, ok := data["params"].(map[string]interface{}); ok {
                if orderbook, ok := params["data"].(map[string]interface{}); ok {
                    // Get bids and asks
                    bids, _ := orderbook["bids"].([]interface{})
                    asks, _ := orderbook["asks"].([]interface{})
                    
                    // Collect current prices for header
                    var currentPrices []string
                    
                    // Get bid prices (in reverse order)
                    for i := 2; i >= 0; i-- {
                        if i < len(bids) {
                            if bid, ok := bids[i].([]interface{}); ok && len(bid) >= 2 {
                                currentPrices = append(currentPrices, fmt.Sprintf("%v", bid[0]))
                            } else {
                                currentPrices = append(currentPrices, "-")
                            }
                        } else {
                            currentPrices = append(currentPrices, "-")
                        }
                    }
                    
                    // Get ask prices
                    for i := 0; i < 3; i++ {
                        if i < len(asks) {
                            if ask, ok := asks[i].([]interface{}); ok && len(ask) >= 2 {
                                currentPrices = append(currentPrices, fmt.Sprintf("%v", ask[0]))
                            } else {
                                currentPrices = append(currentPrices, "-")
                            }
                        } else {
                            currentPrices = append(currentPrices, "-")
                        }
                    }
                    
                    // Check if prices changed or print header every 20 lines
                    pricesChanged := false
                    if lastPrices == nil || len(lastPrices) != len(currentPrices) {
                        pricesChanged = true
                    } else {
                        for i, price := range currentPrices {
                            if price != lastPrices[i] {
                                pricesChanged = true
                                break
                            }
                        }
                    }
                    
                    // Print header if prices changed or every 20 lines
                    if pricesChanged || lineCount%20 == 0 {
                        fmt.Println()
                        header := "Time     |"
                        separator := "---------|"
                        for _, price := range currentPrices {
                            if len(price) > 7 {
                                price = price[:7]
                            }
                            header += fmt.Sprintf(" %7s |", price)
                            separator += "---------|"
                        }
                        header += " Spread"
                        separator += "-------"
                        fmt.Println(header)
                        fmt.Println(separator)
                        lastPrices = currentPrices
                    }
                    
                    // Build data line with sizes
                    line := fmt.Sprintf("%s |", time.Now().Format("15:04:05"))
                    
                    // Add bid sizes (in reverse order)
                    for i := 2; i >= 0; i-- {
                        if i < len(bids) {
                            if bid, ok := bids[i].([]interface{}); ok && len(bid) >= 2 {
                                line += fmt.Sprintf(" %7v |", bid[1])
                            } else {
                                line += "       - |"
                            }
                        } else {
                            line += "       - |"
                        }
                    }
                    
                    // Add ask sizes
                    for i := 0; i < 3; i++ {
                        if i < len(asks) {
                            if ask, ok := asks[i].([]interface{}); ok && len(ask) >= 2 {
                                line += fmt.Sprintf(" %7v |", ask[1])
                            } else {
                                line += "       - |"
                            }
                        } else {
                            line += "       - |"
                        }
                    }
                    
                    // Calculate and add spread
                    // Best bid is at currentPrices[2] (third column)
                    // Best ask is at currentPrices[3] (fourth column)
                    if len(currentPrices) >= 4 && currentPrices[2] != "-" && currentPrices[3] != "-" {
                        var bestBid, bestAsk float64
                        fmt.Sscanf(currentPrices[2], "%f", &bestBid)
                        fmt.Sscanf(currentPrices[3], "%f", &bestAsk)
                        spread := bestAsk - bestBid
                        line += fmt.Sprintf(" %.2f", spread)
                    } else {
                        line += "   -"
                    }
                    
                    fmt.Println(line)
                    lineCount++
                }
            }
        }
    }
}