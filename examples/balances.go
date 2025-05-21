package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
)

// JSONRPCRequest defines the structure for a JSON-RPC request.
type JSONRPCRequest struct {
	ID      string          `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// BalancesParams defines the structure for the params field for a 'balances' request.
type BalancesParams struct {
	Account string `json:"account,omitempty"`
}

const (
	requestPipeSuffix      = ".req.pipe"
	defaultDaemonChannelID = "rysk_ipc_default" // Ensure this matches connect.go's default or is passed
)

func main() {
	channelID := flag.String("channel_id", defaultDaemonChannelID, "IPC channel ID of the daemon")
	accountAddress := flag.String("account", "0x000000000000000000000000000000000000dEaD", "Account address for balances")
	requestID := flag.String("request_id", "bal-req-client-1", "Unique ID for this balance request")
	flag.Parse()

	log.SetFlags(log.Ltime) // Minimal logging
	log.Printf("Sending balances request (ID: %s) for account %s to channel %s", *requestID, *accountAddress, *channelID)

	requestPipePath := filepath.Join(os.TempDir(), *channelID+requestPipeSuffix)

	params := BalancesParams{Account: *accountAddress}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("Error marshalling params: %v", err)
	}

	req := JSONRPCRequest{
		ID:      *requestID,
		JSONRPC: "2.0",
		Method:  "balances", // The daemon will look for this method
		Params:  paramsJSON,
	}

	requestJSON, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("Error marshalling request: %v", err)
	}

	pipe, err := os.OpenFile(requestPipePath, os.O_WRONLY, 0) // Open for write-only
	if err != nil {
		log.Fatalf("Error opening pipe %s: %v", requestPipePath, err)
	}
	defer pipe.Close()

	if _, err = pipe.Write(requestJSON); err != nil {
		log.Fatalf("Error writing to pipe: %v", err)
	}

	log.Printf("Request sent to %s. Exiting.", requestPipePath)
}
