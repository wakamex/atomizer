package main

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// DeribitEd25519Client handles asymmetric key authentication for Deribit
type DeribitEd25519Client struct {
	ClientID   string
	PrivateKey ed25519.PrivateKey
	BaseURL    string
	HTTPClient *http.Client
}

// NewDeribitEd25519Client creates a new client with Ed25519 authentication
func NewDeribitEd25519Client(clientID string, privateKeyPEM string, testnet bool) (*DeribitEd25519Client, error) {
	// Parse the private key
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	// Try parsing as PKCS8 first (most common for Ed25519)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Try parsing as Ed25519 seed
		if len(block.Bytes) == ed25519.SeedSize {
			privateKey := ed25519.NewKeyFromSeed(block.Bytes)
			key = privateKey
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

	return &DeribitEd25519Client{
		ClientID:   clientID,
		PrivateKey: ed25519Key,
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// Authenticate performs the client_signature authentication
func (c *DeribitEd25519Client) Authenticate() (string, error) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	nonce := fmt.Sprintf("%d", timestamp) // Simple nonce using timestamp
	data := ""

	// Create the string to sign
	stringToSign := fmt.Sprintf("%d\n%s\n%s", timestamp, nonce, data)
	
	// Sign with Ed25519
	signature := ed25519.Sign(c.PrivateKey, []byte(stringToSign))
	
	// Encode signature as base64
	signatureB64 := base64.StdEncoding.EncodeToString(signature)

	// Create auth request
	authRequest := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "public/auth",
		"params": map[string]interface{}{
			"grant_type": "client_signature",
			"client_id":  c.ClientID,
			"timestamp":  timestamp,
			"signature":  signatureB64,
			"nonce":      nonce,
			"data":       data,
		},
	}

	// Send request
	jsonData, err := json.Marshal(authRequest)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/v2/public/auth",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Check for error
	if errObj, ok := result["error"].(map[string]interface{}); ok {
		return "", fmt.Errorf("auth error: %v", errObj["message"])
	}

	// Extract access token
	if res, ok := result["result"].(map[string]interface{}); ok {
		if token, ok := res["access_token"].(string); ok {
			return token, nil
		}
	}

	return "", fmt.Errorf("failed to extract access token from response: %s", string(body))
}

// CallPrivateMethod calls a private API method using the access token
func (c *DeribitEd25519Client) CallPrivateMethod(accessToken string, method string, params map[string]interface{}) (interface{}, error) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  method,
		"params":  params,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/api/v2/private/"+method, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if errObj, ok := result["error"].(map[string]interface{}); ok {
		return nil, fmt.Errorf("API error: %v", errObj["message"])
	}

	return result["result"], nil
}