package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// DeribitClient is a complete Deribit API client using Ed25519 authentication
type DeribitClient struct {
	ClientID    string
	PrivateKey  ed25519.PrivateKey
	BaseURL     string
	HTTPClient  *http.Client
	
	// Token management
	accessToken  string
	refreshToken string
	tokenExpiry  time.Time
	tokenMutex   sync.RWMutex
}

// DeribitResponse is the standard JSON-RPC response
type DeribitResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *DeribitError   `json:"error,omitempty"`
}

// DeribitError represents API errors
type DeribitError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewDeribitClient creates a new Deribit client with Ed25519 authentication
func NewDeribitClient(clientID string, privateKeyPEM string, testnet bool) (*DeribitClient, error) {
	// Parse the private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	// Try parsing as PKCS8 first
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try as raw Ed25519 seed
		if len(block.Bytes) == ed25519.SeedSize {
			key = ed25519.NewKeyFromSeed(block.Bytes)
		} else {
			return nil, fmt.Errorf("failed to parse private key: %v", err)
		}
	}

	ed25519Key, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an Ed25519 private key")
	}

	baseURL := "https://www.deribit.com"
	if testnet {
		baseURL = "https://test.deribit.com"
	}

	client := &DeribitClient{
		ClientID:   clientID,
		PrivateKey: ed25519Key,
		BaseURL:    baseURL,
		HTTPClient: newHTTPClient(),
	}

	// Authenticate immediately
	if err := client.authenticate(); err != nil {
		return nil, fmt.Errorf("initial authentication failed: %w", err)
	}

	return client, nil
}

// authenticate performs the Ed25519 authentication
func (c *DeribitClient) authenticate() error {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	nonce := fmt.Sprintf("%d", timestamp)
	data := ""

	// Create signature
	stringToSign := fmt.Sprintf("%d\n%s\n%s", timestamp, nonce, data)
	signature := ed25519.Sign(c.PrivateKey, []byte(stringToSign))
	signatureB64 := base64.StdEncoding.EncodeToString(signature)

	// Build request
	params := map[string]interface{}{
		"grant_type": "client_signature",
		"client_id":  c.ClientID,
		"timestamp":  timestamp,
		"signature":  signatureB64,
		"nonce":      nonce,
		"data":       data,
	}

	resp, err := c.call("public/auth", params, false)
	if err != nil {
		return err
	}

	// Parse auth response
	var authResult struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.Unmarshal(resp.Result, &authResult); err != nil {
		return fmt.Errorf("failed to parse auth response: %w", err)
	}

	// Store tokens
	c.tokenMutex.Lock()
	c.accessToken = authResult.AccessToken
	c.refreshToken = authResult.RefreshToken
	c.tokenExpiry = time.Now().Add(time.Duration(authResult.ExpiresIn-60) * time.Second) // 1 minute buffer
	c.tokenMutex.Unlock()

	log.Printf("Authenticated successfully, token expires at %s", c.tokenExpiry.Format(time.RFC3339))
	return nil
}

// ensureAuthenticated checks and refreshes token if needed
func (c *DeribitClient) ensureAuthenticated() error {
	c.tokenMutex.RLock()
	needsRefresh := time.Now().After(c.tokenExpiry)
	c.tokenMutex.RUnlock()

	if needsRefresh {
		return c.authenticate()
	}
	return nil
}

// call makes an API call
func (c *DeribitClient) call(method string, params map[string]interface{}, private bool) (*DeribitResponse, error) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      time.Now().UnixNano(),
		"method":  method,
		"params":  params,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/api/v2/", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	
	if private {
		c.tokenMutex.RLock()
		token := c.accessToken
		c.tokenMutex.RUnlock()
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result DeribitResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w, body: %s", err, string(body))
	}

	if result.Error != nil {
		return nil, fmt.Errorf("API error %d: %s", result.Error.Code, result.Error.Message)
	}

	return &result, nil
}

// Public API Methods

// GetOrderBook fetches the order book for an instrument
func (c *DeribitClient) GetOrderBook(instrument string) (*OrderBook, error) {
	params := map[string]interface{}{
		"instrument_name": instrument,
		"depth":           10,
	}

	resp, err := c.call("public/get_order_book", params, false)
	if err != nil {
		return nil, err
	}

	var book OrderBook
	if err := json.Unmarshal(resp.Result, &book); err != nil {
		return nil, err
	}

	return &book, nil
}

// Private API Methods

// GetAccountSummary gets account summary for a currency
func (c *DeribitClient) GetAccountSummary(currency string) (map[string]interface{}, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"currency": currency,
		"extended": true,
	}

	resp, err := c.call("private/get_account_summary", params, true)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// PlaceOrder places a new order
func (c *DeribitClient) PlaceOrder(instrument string, amount float64, orderType string, side string, price float64) (*Order, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"instrument_name": instrument,
		"amount":          amount,
		"type":            orderType,
		"label":           "atomizer",
		"advanced":        "usd",
	}

	// Add price for limit orders
	if orderType == "limit" {
		params["price"] = price
	}

	method := fmt.Sprintf("private/%s", side) // private/buy or private/sell
	resp, err := c.call(method, params, true)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal(resp.Result, &order); err != nil {
		return nil, err
	}

	return &order, nil
}

// GetOpenOrders gets all open orders
func (c *DeribitClient) GetOpenOrders() ([]Order, error) {
	if err := c.ensureAuthenticated(); err != nil {
		return nil, err
	}

	resp, err := c.call("private/get_open_orders", map[string]interface{}{}, true)
	if err != nil {
		return nil, err
	}

	var orders []Order
	if err := json.Unmarshal(resp.Result, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// CancelOrder cancels an order
func (c *DeribitClient) CancelOrder(orderID string) error {
	if err := c.ensureAuthenticated(); err != nil {
		return err
	}

	params := map[string]interface{}{
		"order_id": orderID,
	}

	_, err := c.call("private/cancel", params, true)
	return err
}

// Data structures

type OrderBook struct {
	Timestamp       int64       `json:"timestamp"`
	Stats           interface{} `json:"stats"`
	State           string      `json:"state"`
	InstrumentName  string      `json:"instrument_name"`
	Bids            [][]float64 `json:"bids"`
	Asks            [][]float64 `json:"asks"`
	MarkPrice       float64     `json:"mark_price"`
	IndexPrice      float64     `json:"index_price"`
	UnderlyingPrice float64     `json:"underlying_price"`
	UnderlyingIndex string      `json:"underlying_index"`
}

type Order struct {
	OrderID        string  `json:"order_id"`
	InstrumentName string  `json:"instrument_name"`
	Direction      string  `json:"direction"`
	Amount         float64 `json:"amount"`
	Price          float64 `json:"price"`
	OrderType      string  `json:"order_type"`
	OrderState     string  `json:"order_state"`
	FilledAmount   float64 `json:"filled_amount"`
	AveragePrice   float64 `json:"average_price"`
	LastUpdateTime int64   `json:"last_update_timestamp"`
}