// Derive ReplaceOrder Bug Test (Auto-detects subaccount)
// 
// To run:
//    go mod init test && go mod tidy
//    export DERIVE_PRIVATE_KEY=your_private_key_hex
//    export DERIVE_WALLET=your_wallet_address  
//    go run test.go

package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/websocket"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

func init() {
	fmt.Fprintf(os.Stderr, "INIT: Program starting...\n")
}

func main() {
	// Write to stderr to ensure output isn't buffered
	fmt.Fprintf(os.Stderr, "MAIN: Entered main function\n")
	
	// Also try regular stdout
	fmt.Println("=== DERIVE REPLACEORDER BUG TEST STARTING ===")
	fmt.Printf("Time: %s\n", time.Now().Format(time.RFC3339))
	
	// Force flush stdout
	os.Stdout.Sync()
	
	privateKeyHex := os.Getenv("DERIVE_PRIVATE_KEY")
	walletAddress := os.Getenv("DERIVE_WALLET")
	
	fmt.Printf("Environment check:\n")
	fmt.Printf("  DERIVE_PRIVATE_KEY: %s\n", maskString(privateKeyHex))
	fmt.Printf("  DERIVE_WALLET: %s\n", walletAddress)
	
	if privateKeyHex == "" || walletAddress == "" {
		fmt.Println("ERROR: Set DERIVE_PRIVATE_KEY and DERIVE_WALLET environment variables")
		os.Exit(1)
	}
	
	fmt.Println("\n1. Parsing private key...")
	
	// Parse private key
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	fmt.Printf("   Key length after trim: %d\n", len(privateKeyHex))
	
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		fmt.Printf("ERROR: Invalid private key hex: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("   Decoded %d bytes\n", len(privateKeyBytes))
	
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		fmt.Printf("ERROR: Failed to parse private key: %v\n", err)
		os.Exit(1)
	}
	
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("   Derived address: %s\n", address.Hex())
	
	// Connect to WebSocket
	fmt.Println("\n2. Connecting to WebSocket...")
	fmt.Printf("   URL: wss://api.lyra.finance/ws\n")
	fmt.Printf("   Time: %s\n", time.Now().Format("15:04:05.000"))
	
	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 30 * time.Second
	
	conn, resp, err := dialer.Dial("wss://api.lyra.finance/ws", nil)
	if err != nil {
		fmt.Printf("ERROR: Failed to connect: %v\n", err)
		if resp != nil {
			fmt.Printf("   HTTP Status: %s\n", resp.Status)
		}
		os.Exit(1)
	}
	defer func() {
		fmt.Println("\nClosing WebSocket connection...")
		conn.Close()
	}()
	
	fmt.Println("   Connected successfully!")
	fmt.Printf("   Time: %s\n", time.Now().Format("15:04:05.000"))
	
	// Start message reader
	fmt.Println("\n3. Starting WebSocket message reader...")
	responses := make(chan map[string]interface{}, 10) // Buffer for responses
	errors := make(chan error, 1)
	
	go func() {
		fmt.Println("   Reader goroutine started")
		msgCount := 0
		for {
			var msg map[string]interface{}
			fmt.Printf("   [%s] Waiting for message...\n", time.Now().Format("15:04:05.000"))
			
			if err := conn.ReadJSON(&msg); err != nil {
				fmt.Printf("   ERROR in reader: %v\n", err)
				errors <- err
				return
			}
			
			msgCount++
			fmt.Printf("   [%s] Message #%d received\n", time.Now().Format("15:04:05.000"), msgCount)
			
			// Log all messages for debugging
			if method, ok := msg["method"]; ok {
				fmt.Printf("   Notification: %s\n", method)
			}
			if id, ok := msg["id"]; ok && id != nil {
				fmt.Printf("   Response for ID: %v\n", id)
				responses <- msg
			}
		}
	}()
	
	// Give reader time to start
	time.Sleep(100 * time.Millisecond)
	
	// Login
	fmt.Println("\n4. Attempting login...")
	fmt.Printf("   Time: %s\n", time.Now().Format("15:04:05.000"))
	
	subaccountID, err := login(conn, privateKey, walletAddress, responses, errors)
	if err != nil {
		fmt.Printf("ERROR: Login failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("   Login successful! Subaccount ID: %d\n", subaccountID)
	
	fmt.Println("\n=== STARTING REPLACEORDER BUG TEST ===")
	
	// Step 1: Place order
	fmt.Println("\n5. Placing test order...")
	fmt.Printf("   Time: %s\n", time.Now().Format("15:04:05.000"))
	
	orderID, err := placeOrder(conn, privateKey, walletAddress, subaccountID, responses, errors)
	if err != nil {
		fmt.Printf("ERROR: Failed to place order: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("   Order placed successfully! ID: %s\n", orderID)
	
	time.Sleep(2 * time.Second)
	
	// Step 2: Get open orders before replace
	fmt.Println("\n6. Getting open orders BEFORE replace...")
	fmt.Printf("   Time: %s\n", time.Now().Format("15:04:05.000"))
	
	ordersBefore, err := getOpenOrders(conn, subaccountID, responses, errors)
	if err != nil {
		fmt.Printf("   WARNING: Failed to get orders: %v\n", err)
		ordersBefore = []string{}
	} else {
		fmt.Printf("   Found %d open orders:\n", len(ordersBefore))
		for i, o := range ordersBefore {
			fmt.Printf("     [%d] %s\n", i+1, o)
		}
	}
	
	// Step 3: Replace order
	fmt.Printf("\n7. Attempting to REPLACE order %s...\n", orderID)
	fmt.Printf("   Time: %s\n", time.Now().Format("15:04:05.000"))
	
	newOrderID, err := replaceOrder(conn, privateKey, walletAddress, orderID, subaccountID, responses, errors)
	if err != nil {
		fmt.Printf("   ERROR: Replace failed: %v\n", err)
		newOrderID = ""
	} else {
		fmt.Printf("   SUCCESS: Got new order ID: %s\n", newOrderID)
	}
	
	fmt.Println("\n   Waiting 2 seconds for order to settle...")
	time.Sleep(2 * time.Second)
	
	// Step 4: Get open orders after replace
	fmt.Println("\n8. Getting open orders AFTER replace...")
	fmt.Printf("   Time: %s\n", time.Now().Format("15:04:05.000"))
	
	ordersAfter, err := getOpenOrders(conn, subaccountID, responses, errors)
	if err != nil {
		fmt.Printf("   WARNING: Failed to get orders: %v\n", err)
		ordersAfter = []string{}
	} else {
		fmt.Printf("   Found %d open orders:\n", len(ordersAfter))
		for i, o := range ordersAfter {
			fmt.Printf("     [%d] %s\n", i+1, o)
		}
	}
	
	// Step 5: Analyze results
	fmt.Println("\n=== ANALYZING RESULTS ===")
	oldFound := false
	newFound := false
	
	for _, o := range ordersAfter {
		if o == orderID {
			oldFound = true
			fmt.Printf("   Found original order: %s\n", orderID)
		}
		if o == newOrderID && newOrderID != "" {
			newFound = true
			fmt.Printf("   Found new order: %s\n", newOrderID)
		}
	}
	
	fmt.Println("\n=== FINAL VERDICT ===")
	if oldFound && newFound {
		fmt.Println("❌ BUG CONFIRMED: Both old and new orders exist!")
		fmt.Println("   This means replace created a new order without canceling the old one")
	} else if oldFound && !newFound {
		fmt.Println("❌ BUG CONFIRMED: Old order exists but new doesn't!")
		fmt.Println("   This means replace returned a phantom order ID")
	} else if !oldFound && newFound {
		fmt.Println("✅ SUCCESS: Replace worked correctly")
		fmt.Println("   Old order was canceled and new order was created")
	} else {
		fmt.Println("❌ BUG CONFIRMED: Neither order exists!")
		fmt.Println("   Both orders disappeared or were never created")
	}
	
	// Cleanup
	fmt.Println("\n9. Cleaning up test orders...")
	cleanupCount := 0
	for _, o := range ordersAfter {
		fmt.Printf("   Canceling order %s...\n", o)
		if err := cancelOrder(conn, o, responses, errors); err != nil {
			fmt.Printf("   WARNING: Failed to cancel %s: %v\n", o, err)
		} else {
			cleanupCount++
		}
	}
	fmt.Printf("   Cleaned up %d orders\n", cleanupCount)
	
	fmt.Println("\n=== TEST COMPLETE ===")
	fmt.Printf("Time: %s\n", time.Now().Format(time.RFC3339))
}

func maskString(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "..." + s[len(s)-4:]
}

func login(conn *websocket.Conn, privateKey *ecdsa.PrivateKey, wallet string, responses chan map[string]interface{}, errors chan error) (int64, error) {
	fmt.Println("   Preparing login request...")
	timestamp := time.Now().UTC().UnixMilli()
	msg := fmt.Sprintf("%d", timestamp)
	
	fmt.Printf("   Timestamp: %d\n", timestamp)
	fmt.Printf("   Message to sign: '%s'\n", msg)
	
	// Sign timestamp with Ethereum personal_sign prefix
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msg), msg)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))
	fmt.Printf("   Prefixed message: '\\x19Ethereum Signed Message:\\n%d%s'\n", len(msg), msg)
	fmt.Printf("   Message hash: 0x%s\n", hex.EncodeToString(hash.Bytes()))
	
	sig, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return 0, err
	}
	
	// Fix signature format - transform V from 0/1 to 27/28
	fmt.Printf("   Raw signature v: %d\n", sig[64])
	sig[64] += 27
	fmt.Printf("   Fixed signature v: %d\n", sig[64])
	
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	fmt.Printf("   Signer address: %s\n", address.Hex())
	fmt.Printf("   Wallet address: %s\n", wallet)
	
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "public/login",
		"params": map[string]interface{}{
			"wallet":    wallet,
			"timestamp": fmt.Sprintf("%d", timestamp), // String format
			"signature": "0x" + hex.EncodeToString(sig),
		},
		"id": "login",
	}
	
	fmt.Printf("   Login params: wallet=%s, signer=%s\n", wallet, strings.ToLower(address.Hex()))
	fmt.Printf("   Sending login request at %s\n", time.Now().Format("15:04:05.000"))
	
	if err := conn.WriteJSON(req); err != nil {
		return 0, fmt.Errorf("failed to send login: %v", err)
	}
	
	fmt.Println("   Waiting for login response...")
	
	select {
	case resp := <-responses:
		fmt.Printf("   Got login response at %s\n", time.Now().Format("15:04:05.000"))
		
		if err, ok := resp["error"]; ok && err != nil {
			return 0, fmt.Errorf("login error: %v", err)
		}
		
		// Extract subaccount ID from login response
		// Based on production code, result is an array of subaccount IDs
		if result, ok := resp["result"].([]interface{}); ok && len(result) > 0 {
			if id, ok := result[0].(float64); ok {
				fmt.Printf("   Got subaccount ID: %d\n", int64(id))
				return int64(id), nil
			}
		}
		
		// Default subaccount
		fmt.Println("   No subaccount ID found, using default (0)")
		return 0, nil
		
	case err := <-errors:
		return 0, fmt.Errorf("reader error during login: %v", err)
		
	case <-time.After(10 * time.Second):
		return 0, fmt.Errorf("login timeout after 10 seconds")
	}
}

func placeOrder(conn *websocket.Conn, privateKey *ecdsa.PrivateKey, wallet string, subaccountID int64, responses chan map[string]interface{}, errors chan error) (string, error) {
	fmt.Println("   Preparing order request with EIP-712 signing...")
	
	// Order parameters
	instrumentName := "ETH-PERP"
	side := "buy"
	price := 2000.0  // Much lower price that won't get filled
	amount := 0.1    // Minimum order size
	maxFee := int64(1000)
	
	// Create nonce and expiry
	nonce := uint64(time.Now().UnixMilli())*1000 + 1
	expiry := time.Now().Unix() + 3600
	
	// Create EIP-712 signature
	signer := crypto.PubkeyToAddress(privateKey.PublicKey)
	
	// For ETH-PERP, use the correct values from the API
	// These are hardcoded for ETH-PERP on mainnet
	baseAssetAddress := "0xAf65752C4643E25C02F693f9D4FE19cF23a095E3" // ETH-PERP base asset
	baseAssetSubID := "0"
	moduleAddress := "0xB8D20c2B7a1Ad2EE33Bc50eF10876eD3035b5e7b" // Trade module
	
	// Convert amounts to wei (18 decimals) - using big.Float for precision
	limitPriceWei, _ := new(big.Float).Mul(big.NewFloat(price), big.NewFloat(1e18)).Int(nil)
	amountWei, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e18)).Int(nil)
	maxFeeWei := new(big.Int).Mul(big.NewInt(maxFee), big.NewInt(1e18))
	
	// Debug output
	fmt.Printf("   Signing params:\n")
	fmt.Printf("     Owner: %s\n", wallet)
	fmt.Printf("     Signer: %s\n", signer.Hex())
	fmt.Printf("     SubaccountID: %d\n", subaccountID)
	fmt.Printf("     Nonce: %d\n", nonce)
	fmt.Printf("     Expiry: %d\n", expiry)
	fmt.Printf("     LimitPrice (wei): %s\n", limitPriceWei.String())
	fmt.Printf("     Amount (wei): %s\n", amountWei.String())
	fmt.Printf("     MaxFee (wei): %s\n", maxFeeWei.String())
	
	// Create signature
	signature, err := signOrder(
		privateKey,
		wallet,
		signer.Hex(),
		subaccountID,
		nonce,
		expiry,
		moduleAddress,
		baseAssetAddress,
		baseAssetSubID,
		limitPriceWei,
		amountWei,
		maxFeeWei,
		side == "buy",
	)
	if err != nil {
		return "", fmt.Errorf("failed to sign order: %v", err)
	}
	
	fmt.Printf("     Signature: %s\n", signature)
	
	// Create order request
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "private/order",
		"params": map[string]interface{}{
			"instrument_name":      instrumentName,
			"direction":           side,
			"order_type":         "limit",
			"time_in_force":      "gtc",
			"amount":             fmt.Sprintf("%.6f", amount),
			"limit_price":        fmt.Sprintf("%.6f", price),
			"max_fee":            fmt.Sprintf("%d", maxFee),
			"subaccount_id":      subaccountID,
			"nonce":              nonce,
			"signature_expiry_sec": expiry,
			"owner":              wallet,  // Add owner field
			"signer":             signer.Hex(),  // Use checksummed address
			"signature":          signature,
			"mmp":                false,
		},
		"id": fmt.Sprintf("order_%d", time.Now().UnixNano()),
	}
	
	fmt.Printf("   Order params: ETH-PERP, buy, %.2f @ %.2f, max_fee=1000, subaccount_id=%d\n", amount, price, subaccountID)
	fmt.Printf("   Full order request:\n%s\n", jsonString(req))
	fmt.Printf("   Sending order at %s\n", time.Now().Format("15:04:05.000"))
	
	if err := conn.WriteJSON(req); err != nil {
		return "", fmt.Errorf("failed to send order: %v", err)
	}
	
	fmt.Println("   Waiting for order response...")
	
	select {
	case resp := <-responses:
		fmt.Printf("   Got order response at %s\n", time.Now().Format("15:04:05.000"))
		if err, ok := resp["error"]; ok && err != nil {
			return "", fmt.Errorf("order error: %v", err)
		}
		
		// Extract order ID from response
		if result, ok := resp["result"].(map[string]interface{}); ok {
			if order, ok := result["order"].(map[string]interface{}); ok {
				if orderID, ok := order["order_id"].(string); ok {
					return orderID, nil
				}
			}
		}
		return "", fmt.Errorf("no order ID in response: %v", resp)
		
	case err := <-errors:
		return "", fmt.Errorf("reader error during order: %v", err)
		
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("order timeout after 10 seconds")
	}
}

func replaceOrder(conn *websocket.Conn, privateKey *ecdsa.PrivateKey, wallet string, orderID string, subaccountID int64, responses chan map[string]interface{}, errors chan error) (string, error) {
	fmt.Println("   Preparing replace request with proper signature...")
	
	// Replace requires a fully signed new order
	nonce := uint64(time.Now().UnixMilli())*1000 + 2
	expiry := time.Now().Unix() + 3600
	signer := crypto.PubkeyToAddress(privateKey.PublicKey)
	
	// Sign the new order
	price := 1995.0
	amount := 0.1
	maxFee := int64(1000)
	
	// Create signature for the replacement order
	baseAssetAddress := "0xAf65752C4643E25C02F693f9D4FE19cF23a095E3"
	baseAssetSubID := "0"
	moduleAddress := "0xB8D20c2B7a1Ad2EE33Bc50eF10876eD3035b5e7b"
	
	limitPriceWei, _ := new(big.Float).Mul(big.NewFloat(price), big.NewFloat(1e18)).Int(nil)
	amountWei, _ := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1e18)).Int(nil)
	maxFeeWei := new(big.Int).Mul(big.NewInt(maxFee), big.NewInt(1e18))
	
	signature, err := signOrder(
		privateKey,
		wallet,
		signer.Hex(),
		subaccountID,
		nonce,
		expiry,
		moduleAddress,
		baseAssetAddress,
		baseAssetSubID,
		limitPriceWei,
		amountWei,
		maxFeeWei,
		true, // buy
	)
	if err != nil {
		return "", fmt.Errorf("failed to sign replacement order: %v", err)
	}
	
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "private/replace",
		"params": map[string]interface{}{
			"order_id_to_cancel": orderID,
			"instrument_name": "ETH-PERP",
			"direction": "buy",
			"order_type": "limit",
			"time_in_force": "gtc",
			"amount": fmt.Sprintf("%.6f", amount),
			"limit_price": fmt.Sprintf("%.6f", price),
			"max_fee": fmt.Sprintf("%d", maxFee),
			"subaccount_id": subaccountID,
			"nonce": nonce,
			"signature_expiry_sec": expiry,
			"owner": wallet,
			"signer": signer.Hex(),
			"signature": signature,
			"mmp": false,
		},
		"id": fmt.Sprintf("replace_%d", time.Now().UnixNano()),
	}
	
	fmt.Printf("   Replace params: order_id=%s, new_price=1995\n", orderID)
	fmt.Printf("   Full request: %s\n", jsonString(req))
	fmt.Printf("   Sending replace at %s\n", time.Now().Format("15:04:05.000"))
	
	if err := conn.WriteJSON(req); err != nil {
		return "", fmt.Errorf("failed to send replace: %v", err)
	}
	
	fmt.Println("   Waiting for replace response...")
	
	select {
	case resp := <-responses:
		fmt.Printf("   Got replace response at %s\n", time.Now().Format("15:04:05.000"))
		fmt.Printf("   Full response: %s\n", jsonString(resp))
		
		if err, ok := resp["error"]; ok && err != nil {
			return "", fmt.Errorf("replace error: %v", err)
		}
		
		// Try to extract new order ID
		if result, ok := resp["result"].(map[string]interface{}); ok {
			// Check for order.order_id
			if order, ok := result["order"].(map[string]interface{}); ok {
				if orderID, ok := order["order_id"].(string); ok {
					return orderID, nil
				}
			}
			// Check for direct order_id
			if orderID, ok := result["order_id"].(string); ok {
				return orderID, nil
			}
		}
		return "", fmt.Errorf("no order ID in replace response")
		
	case err := <-errors:
		return "", fmt.Errorf("reader error during replace: %v", err)
		
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("replace timeout after 10 seconds")
	}
}

func getOpenOrders(conn *websocket.Conn, subaccountID int64, responses chan map[string]interface{}, errors chan error) ([]string, error) {
	fmt.Println("   Preparing get_open_orders request...")
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "private/get_open_orders",
		"params": map[string]interface{}{
			"instrument_name": "ETH-PERP",
			"subaccount_id": subaccountID,
		},
		"id": fmt.Sprintf("orders_%d", time.Now().UnixNano()),
	}
	
	fmt.Printf("   Get orders params: instrument=ETH-PERP, subaccount_id=%d\n", subaccountID)
	fmt.Printf("   Sending request at %s\n", time.Now().Format("15:04:05.000"))
	
	if err := conn.WriteJSON(req); err != nil {
		return nil, fmt.Errorf("failed to send get_open_orders: %v", err)
	}
	
	fmt.Println("   Waiting for orders response...")
	
	select {
	case resp := <-responses:
		fmt.Printf("   Got orders response at %s\n", time.Now().Format("15:04:05.000"))
		if err, ok := resp["error"]; ok && err != nil {
			return nil, fmt.Errorf("get orders error: %v", err)
		}
		
		var orderIDs []string
		if result, ok := resp["result"].(map[string]interface{}); ok {
			if orders, ok := result["orders"].([]interface{}); ok {
				for _, o := range orders {
					if order, ok := o.(map[string]interface{}); ok {
						if orderID, ok := order["order_id"].(string); ok {
							orderIDs = append(orderIDs, orderID)
						}
					}
				}
			}
		}
		fmt.Printf("   Found %d order IDs\n", len(orderIDs))
		return orderIDs, nil
		
	case err := <-errors:
		return nil, fmt.Errorf("reader error during get_open_orders: %v", err)
		
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("get orders timeout after 10 seconds")
	}
}

func cancelOrder(conn *websocket.Conn, orderID string, responses chan map[string]interface{}, errors chan error) error {
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method": "private/cancel_order",
		"params": map[string]interface{}{
			"order_id": orderID,
		},
		"id": fmt.Sprintf("cancel_%d", time.Now().UnixNano()),
	}
	
	return conn.WriteJSON(req)
}

func jsonString(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// EIP-712 signing implementation
func signOrder(
	privateKey *ecdsa.PrivateKey,
	owner string,
	signer string,
	subaccountID int64,
	nonce uint64,
	expiry int64,
	moduleAddress string,
	baseAssetAddress string,
	baseAssetSubID string,
	limitPrice *big.Int,
	amount *big.Int,
	maxFee *big.Int,
	isBid bool,
) (string, error) {
	// Domain separator for Derive/Lyra mainnet - FROM PRODUCTION CODE
	domainSeparator := common.HexToHash("0xd96e5f90797da7ec8dc4e276260c7f3f87fedf68775fbe1ef116e996fc60441b")
	
	// Action typehash - FROM PRODUCTION CODE
	actionTypehash := common.HexToHash("0x4d7a9f27c403ff9c0f19bce61d76d82f9aa29f8d6d4b0c5474607d9770d1af17")
	
	fmt.Printf("   EIP-712 Debug:\n")
	fmt.Printf("     Domain Separator: %s\n", domainSeparator.Hex())
	fmt.Printf("     Action Typehash: %s\n", actionTypehash.Hex())
	fmt.Printf("     Module Address: %s\n", moduleAddress)
	fmt.Printf("     Base Asset: %s\n", baseAssetAddress)
	fmt.Printf("     SubID: %s\n", baseAssetSubID)
	
	// Encode module data (Trade module)
	moduleData, err := encodeTradeModuleData(
		common.HexToAddress(baseAssetAddress),
		mustParseBig(baseAssetSubID),
		limitPrice,
		amount,
		maxFee,
		uint64(subaccountID),  // RecipientID should be the same as subaccountID
		isBid,
	)
	if err != nil {
		return "", err
	}
	
	// Hash module data
	moduleDataHash := crypto.Keccak256Hash(moduleData)
	fmt.Printf("     Module Data Length: %d\n", len(moduleData))
	fmt.Printf("     Module Data Hash: %s\n", moduleDataHash.Hex())
	
	// Encode action struct
	actionData, err := encodeActionStruct(
		actionTypehash,
		uint64(subaccountID),
		nonce,
		common.HexToAddress(moduleAddress),
		moduleDataHash,
		expiry,
		common.HexToAddress(owner),
		common.HexToAddress(signer),
	)
	if err != nil {
		return "", err
	}
	
	// Create typed data hash
	actionHash := crypto.Keccak256Hash(actionData)
	fmt.Printf("     Action Data Length: %d\n", len(actionData))
	fmt.Printf("     Action Hash: %s\n", actionHash.Hex())
	
	// Final hash: 0x1901 + domain separator + action hash
	message := append([]byte{0x19, 0x01}, domainSeparator.Bytes()...)
	message = append(message, actionHash.Bytes()...)
	typedDataHash := crypto.Keccak256Hash(message)
	fmt.Printf("     Final Typed Data Hash: %s\n", typedDataHash.Hex())
	
	// Sign
	signature, err := crypto.Sign(typedDataHash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}
	
	// Transform V from 0/1 to 27/28
	signature[64] += 27
	
	return "0x" + hex.EncodeToString(signature), nil
}

func encodeTradeModuleData(asset common.Address, subID, limitPrice, amount, maxFee *big.Int, recipientID uint64, isBid bool) ([]byte, error) {
	// Use ABI encoding like the production code
	addressType, _ := abi.NewType("address", "", nil)
	uintType, _ := abi.NewType("uint256", "", nil)
	intType, _ := abi.NewType("int256", "", nil)
	boolType, _ := abi.NewType("bool", "", nil)
	
	args := abi.Arguments{
		{Type: addressType},
		{Type: uintType},
		{Type: intType},
		{Type: intType},
		{Type: uintType},
		{Type: uintType},
		{Type: boolType},
	}
	
	return args.Pack(asset, subID, limitPrice, amount, maxFee, big.NewInt(int64(recipientID)), isBid)
}

func encodeActionStruct(typehash common.Hash, subaccountID, nonce uint64, module common.Address, moduleDataHash common.Hash, expiry int64, owner, signer common.Address) ([]byte, error) {
	// Use ABI encoding for consistency
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	uintType, _ := abi.NewType("uint256", "", nil)
	addressType, _ := abi.NewType("address", "", nil)
	
	args := abi.Arguments{
		{Type: bytes32Type},
		{Type: uintType},
		{Type: uintType},
		{Type: addressType},
		{Type: bytes32Type},
		{Type: uintType},
		{Type: addressType},
		{Type: addressType},
	}
	
	return args.Pack(
		typehash,
		new(big.Int).SetUint64(subaccountID),
		new(big.Int).SetUint64(nonce),
		module,
		moduleDataHash,
		big.NewInt(expiry),
		owner,
		signer,
	)
}

func mustParseBig(s string) *big.Int {
	n, _ := new(big.Int).SetString(s, 10)
	if n == nil {
		return big.NewInt(0)
	}
	return n
}