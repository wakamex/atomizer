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

// RFQParams defines the structure for the params field for an 'submit_rfq' request.
type RFQParams struct {
	InstrumentName string `json:"instrument_name"`
	Quantity       string `json:"quantity"`
	Side           string `json:"side"`                      // "buy" or "sell"
	ClientOrderID  string `json:"client_order_id,omitempty"` // Optional, but good practice
}

const (
	requestPipeSuffix      = ".req.pipe"
	defaultDaemonChannelID = "rysk_ipc_default" // Ensure this matches connect.go's default or is passed
)

func main() {
	channelID := flag.String("channel_id", defaultDaemonChannelID, "IPC channel ID of the daemon")
	instrument := flag.String("instrument", "ETH-PERP", "Instrument name for RFQ")
	quantity := flag.String("quantity", "1.0", "Quantity for RFQ")
	side := flag.String("side", "buy", "Side for RFQ (buy or sell)")
	clientOrderID := flag.String("client_order_id", "my-rfq-001", "Client-generated order ID for this RFQ")
	requestID := flag.String("request_id", "rfq-req-client-1", "Unique ID for this JSON-RPC request")
	flag.Parse()

	log.SetFlags(log.Ltime) // Minimal logging
	log.Printf("RFQ Client: Sending RFQ (ID: %s) for %s %s %s, client_oid %s, to channel %s",
		*requestID, *side, *quantity, *instrument, *clientOrderID, *channelID)

	requestPipePath := filepath.Join(os.TempDir(), *channelID+requestPipeSuffix)

	params := RFQParams{
		InstrumentName: *instrument,
		Quantity:       *quantity,
		Side:           *side,
		ClientOrderID:  *clientOrderID,
	}
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("RFQ Client: Error marshalling RFQ params: %v", err)
	}

	req := JSONRPCRequest{
		ID:      *requestID,
		JSONRPC: "2.0",
		Method:  "submit_rfq", // Assumed method name for submitting an RFQ
		Params:  paramsJSON,
	}

	requestJSON, err := json.Marshal(req)
	if err != nil {
		log.Fatalf("RFQ Client: Error marshalling RFQ request: %v", err)
	}

	pipe, err := os.OpenFile(requestPipePath, os.O_WRONLY, 0)
	if err != nil {
		log.Fatalf("RFQ Client: Error opening pipe %s: %v", requestPipePath, err)
	}
	defer pipe.Close()

	if _, err = pipe.Write(requestJSON); err != nil {
		log.Fatalf("RFQ Client: Error writing to pipe: %v", err)
	}

	log.Printf("RFQ Client: Request sent to %s. Exiting.", requestPipePath)
}
