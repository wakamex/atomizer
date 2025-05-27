package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var (
		clientID       = flag.String("client-id", os.Getenv("DERIBIT_CLIENT_ID"), "Deribit Client ID")
		privateKeyFile = flag.String("key-file", os.Getenv("DERIBIT_PRIVATE_KEY_FILE"), "Path to Ed25519 private key file")
		testMode       = flag.Bool("test", false, "Use testnet")
	)
	flag.Parse()

	if *clientID == "" {
		log.Fatal("Error: Client ID is required (--client-id or DERIBIT_CLIENT_ID env var)")
	}

	// Read private key
	var privateKeyPEM string
	if *privateKeyFile != "" {
		keyData, err := ioutil.ReadFile(*privateKeyFile)
		if err != nil {
			log.Fatalf("Error reading private key file: %v", err)
		}
		privateKeyPEM = string(keyData)
	} else if envKey := os.Getenv("DERIBIT_PRIVATE_KEY"); envKey != "" {
		privateKeyPEM = envKey
	} else {
		log.Fatal("Error: Private key is required (--key-file or DERIBIT_PRIVATE_KEY env var)")
	}

	fmt.Println("=== Deribit Ed25519 Authentication Test ===")
	fmt.Printf("Client ID: %s\n", *clientID)
	if *testMode {
		fmt.Println("Network: TESTNET")
	} else {
		fmt.Println("Network: MAINNET")
	}
	fmt.Println()

	// Create client
	client, err := NewDeribitEd25519Client(*clientID, privateKeyPEM, *testMode)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	fmt.Println("✅ Ed25519 client created successfully")

	// Authenticate
	fmt.Println("\nAuthenticating...")
	accessToken, err := client.Authenticate()
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
	fmt.Println("✅ Authentication successful")
	fmt.Printf("Access token: %s...%s\n", accessToken[:10], accessToken[len(accessToken)-10:])

	// Test API call - Get account summary
	fmt.Println("\nTesting API call - Get account summary...")
	result, err := client.CallPrivateMethod(accessToken, "get_account_summary", map[string]interface{}{
		"currency": "ETH",
		"extended": true,
	})
	if err != nil {
		log.Printf("API call failed: %v", err)
		// Try with BTC
		fmt.Println("\nRetrying with BTC...")
		result, err = client.CallPrivateMethod(accessToken, "get_account_summary", map[string]interface{}{
			"currency": "BTC",
			"extended": true,
		})
		if err != nil {
			log.Fatalf("API call failed: %v", err)
		}
	}

	// Pretty print result
	jsonData, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println("Account Summary:")
	fmt.Println(string(jsonData))

	fmt.Println("\n✅ Ed25519 authentication is working correctly!")
}