package main

import (
	"encoding/json"
	"flag"
	"io" // Added for reading from pipe
	"log"
	"os"
	"path/filepath"
	"strings" // Added for error checking
	"time"    // Added for timeout/retry logic
	
	"github.com/wakamex/atomizer/internal/types"
)

// BalancesParams defines the structure for the params field for a 'balances' request.
type BalancesParams struct {
	Account string `json:"account,omitempty"`
}

const (
	requestPipeSuffix      = ".req.pipe"
	responsePipeSuffix     = ".res.pipe" // Added to read responses
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

	req := types.JSONRPCRequest{
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
		// pipe.Close() is handled by defer, but if we fatal here, it won't run.
		// Explicitly close before fatal if not using defer for this specific scope.
		log.Fatalf("Error writing to request pipe: %v", err)
	}
	// Explicitly close the request pipe as we are done writing to it.
	// The defer pipe.Close() would also handle this, but being explicit is fine.
	pipe.Close()

	log.Printf("Request sent to %s.", requestPipePath)

	// Now, attempt to read the response from the response pipe
	responsePipePath := filepath.Join(os.TempDir(), *channelID+responsePipeSuffix)
	log.Printf("Attempting to read response from %s", responsePipePath)

	var resPipe *os.File
	var openErr error
	// Retry opening the response pipe for a few seconds, as connect.go might not have written to it yet.
	for i := 0; i < 10; i++ { // Try for up to 5 seconds (10 * 500ms)
		resPipe, openErr = os.OpenFile(responsePipePath, os.O_RDONLY, 0)
		if openErr == nil {
			break // Successfully opened
		}
		// Common errors when pipe isn't ready or writer hasn't opened its end yet.
		if os.IsNotExist(openErr) || strings.Contains(openErr.Error(), "no such file or directory") || strings.Contains(openErr.Error(), "device not configured") {
			log.Printf("Response pipe not ready yet (%v), retrying in 500ms...", openErr)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		// For other errors, fail faster
		log.Fatalf("Error opening response pipe %s: %v", responsePipePath, openErr)
	}
	if openErr != nil {
		log.Fatalf("Failed to open response pipe %s after retries: %v", responsePipePath, openErr)
	}
	defer resPipe.Close()

	log.Printf("Response pipe %s opened. Reading...", responsePipePath)
	responseData, err := io.ReadAll(resPipe) // Read all data from the pipe until EOF (writer closes)
	if err != nil {
		log.Fatalf("Error reading from response pipe %s: %v", responsePipePath, err)
	}

	if len(responseData) > 0 {
		log.Printf("Received response data (length: %d):", len(responseData))
		// Attempt to unmarshal and pretty print if JSON
		var prettyJSON interface{}
		if jsonErr := json.Unmarshal(responseData, &prettyJSON); jsonErr == nil {
			formattedOutput, _ := json.MarshalIndent(prettyJSON, "", "  ")
			log.Printf("Formatted Response:\n%s", string(formattedOutput))
		} else {
			log.Printf("Raw Response (not valid JSON or unmarshal error %v):\n%s", jsonErr, string(responseData))
		}
	} else {
		log.Println("No data received from response pipe.")
	}

	log.Printf("Exiting after processing response.")
}
