package main

import "encoding/json"

// JSONRPCError defines the structure for a JSON-RPC error object.
// Based on structure from sdk/jsonrpc.go
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // To match sdk/jsonrpc.go ErrorData
}

// JSONRPCResponse defines the structure for a generic JSON-RPC response.
// Based on structure from sdk/jsonrpc.go
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"` // Using json.RawMessage for flexibility
	Error   *JSONRPCError   `json:"error,omitempty"`
	// Fields for notifications that might also conform to parts of this structure
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"` // Using json.RawMessage for flexibility
}

// JSONRPCRequest defines the structure for a JSON-RPC request.
// Based on structure from sdk/jsonrpc.go
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// RFQNotification represents the structure of an incoming RFQ message from the server.
// This is specific to how the maker_quote_responder example interprets RFQ messages.
type RFQNotification struct {
	JsonRPC string    `json:"jsonrpc"`
	ID      string    `json:"id"` // This is the RFQ ID or subscription ID
	Result  RFQResult `json:"result,omitempty"` // For some message structures
	Error   *JSONRPCError `json:"error,omitempty"` // Standardized to JSONRPCError
	Method  string    `json:"method,omitempty"` // For server-initiated messages like 'rfq'
	Params  RFQResult `json:"params,omitempty"` // If method is 'rfq', params might hold the data
}

// RFQResult contains the actual details of the RFQ.
// This is specific to how the maker_quote_responder example interprets RFQ data.
type RFQResult struct {
	ID         string `json:"id,omitempty"` // This might be the actual RFQ ID if nested under params
	Asset      string `json:"asset"`
	AssetName  string `json:"assetName,omitempty"`
	ChainID    int    `json:"chainId"`
	Expiry     int64  `json:"expiry"`
	IsPut      bool   `json:"isPut"`
	IsTakerBuy bool   `json:"isTakerBuy"`
	Quantity   string `json:"quantity"`
	Strike     string `json:"strike"`
	Taker      string `json:"taker,omitempty"`
	ReceivedTS int64  `json:"-"` // Local timestamp when RFQ was processed
}

