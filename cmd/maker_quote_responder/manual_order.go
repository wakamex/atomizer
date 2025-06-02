    package main

import (
    "fmt"
    "log"
    "os"
    "time"
    
    "github.com/shopspring/decimal"
)

// RunManualOrder submits a single order using the existing infrastructure
func RunManualOrder() {
    // Get configuration from environment
    privateKey := os.Getenv("DERIVE_PRIVATE_KEY")
    walletAddress := os.Getenv("DERIVE_WALLET_ADDRESS")
    
    if privateKey == "" || walletAddress == "" {
        log.Fatal("Set DERIVE_PRIVATE_KEY and DERIVE_WALLET_ADDRESS environment variables")
    }
    
    // Get order parameters from environment or use defaults
    instrument := os.Getenv("ORDER_INSTRUMENT")
    if instrument == "" {
        instrument = "ETH-PERP"
    }
    
    side := os.Getenv("ORDER_SIDE")
    if side == "" {
        side = "buy"
    }
    
    priceStr := os.Getenv("ORDER_PRICE")
    priceFloat := 0.1
    if priceStr != "" {
        fmt.Sscanf(priceStr, "%f", &priceFloat)
    }
    
    amountStr := os.Getenv("ORDER_AMOUNT")
    amountFloat := 0.1
    if amountStr != "" {
        fmt.Sscanf(amountStr, "%f", &amountFloat)
    }
    
    // Create WebSocket client using existing code
    wsClient, err := NewDeriveWSClient(privateKey, walletAddress)
    if err != nil {
        log.Fatalf("Failed to create WebSocket client: %v", err)
    }
    defer wsClient.Close()
    
    // Get subaccount ID
    subaccountID := wsClient.GetDefaultSubaccount()
    log.Printf("Using subaccount ID: %d", subaccountID)
    
    // Create market maker exchange adapter
    exchange, err := NewDeriveMarketMakerExchange(privateKey, walletAddress)
    if err != nil {
        log.Fatalf("Failed to create exchange: %v", err)
    }
    
    // Place a limit order
    price := decimal.NewFromFloat(priceFloat)
    amount := decimal.NewFromFloat(amountFloat)
    
    log.Printf("Placing %s order: %s %s @ %s", side, amount, instrument, price)
    
    orderID, err := exchange.PlaceLimitOrder(instrument, side, price, amount)
    if err != nil {
        log.Fatalf("Failed to place order: %v", err)
    }
    
    log.Printf("Order placed successfully! Order ID: %s", orderID)
    
    // Wait a bit then check the order
    time.Sleep(2 * time.Second)
    
    orders, err := exchange.GetOpenOrders()
    if err != nil {
        log.Printf("Failed to get open orders: %v", err)
    } else {
        log.Printf("Open orders: %d", len(orders))
        for _, order := range orders {
            log.Printf("  %s: %s %s %s @ %s", 
                order.OrderID, order.Side, order.Amount, order.Instrument, order.Price)
        }
    }
}