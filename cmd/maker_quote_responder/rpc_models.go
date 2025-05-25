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

// RFQConfirmation represents the structure of an incoming RFQ confirmation message from ryskv12
// receiving this indicates the trade has been executed on ryskv12 and we need to hedge it on deribit
type RFQConfirmation struct {
	ID              string  `json:"id,omitempty"`	
	Maker           string  `json:"maker" db:"maker,varchar(42)"`
	AssetAddress    string  `json:"assetAddress" db:"asset_address,varchar(42)"`
	ChainID         int     `json:"chainId" db:"chain_id,integer"`
	Expiry          int     `json:"expiry" db:"expiry,integer"`
	IsPut           bool    `json:"isPut" db:"is_put,boolean"`
	Nonce           string  `json:"nonce" db:"nonce,varchar(255)"`
	Price           string  `json:"price" db:"price,numeric(37)"`
	Quantity        string  `json:"quantity" db:"quantity,numeric(37)"`
	QuoteNonce      string  `json:"quoteNonce" db:"quote_nonce,varchar(255)"`
	QuoteValidUntil int     `json:"quoteValidUntil" db:"quote_valid_until,bigint"`
	QuoteSignature  string  `json:"quoteSignature" db:"quote_signature,varchar(132)"`
	Strike          string  `json:"strike" db:"strike,numeric(37)"`
	Taker           string  `json:"taker" db:"taker,varchar(42)"`
	IsTakerBuy      bool    `json:"isTakerBuy" db:"is_taker_buy,boolean"`
	Signature       string  `json:"signature" db:"signature,varchar(132)"`
	ValidUntil      int     `json:"validUntil" db:"valid_until,bigint"`
	CreatedAt       int     `json:"createdAt,omitempty" db:"created_at,bigint"`
	APR             float64 `json:"apr,omitempty" db:"apr,decimal(5,2)"`
}



// CCXTOrderBook represents the order book structure from CCXT exchange
type CCXTOrderBook struct {
	Bids  [][]float64 // Array of [price, amount] pairs
	Asks  [][]float64 // Array of [price, amount] pairs
	Index float64     // Index price of the underlying asset
}

