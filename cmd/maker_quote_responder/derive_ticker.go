package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DeriveTicker represents ticker data from Derive API
type DeriveTicker struct {
	InstrumentName  string      `json:"instrument_name"`
	BestBidPrice    string      `json:"best_bid_price"`
	BestAskPrice    string      `json:"best_ask_price"`
	BestBidAmount   string      `json:"best_bid_amount"`
	BestAskAmount   string      `json:"best_ask_amount"`
	MarkPrice       string      `json:"mark_price"`
	IndexPrice      string      `json:"index_price"`
	OptionDetails   struct {
		Strike string `json:"strike"`
		OptionType string `json:"option_type"`
	} `json:"option_details"`
}

// GetBidPrice returns bid price as float64
func (t *DeriveTicker) GetBidPrice() float64 {
	var f float64
	fmt.Sscanf(t.BestBidPrice, "%f", &f)
	return f
}

// GetAskPrice returns ask price as float64
func (t *DeriveTicker) GetAskPrice() float64 {
	var f float64
	fmt.Sscanf(t.BestAskPrice, "%f", &f)
	return f
}

// GetBidSize returns bid size as float64
func (t *DeriveTicker) GetBidSize() float64 {
	var f float64
	fmt.Sscanf(t.BestBidAmount, "%f", &f)
	if f == 0 {
		return 1.0 // Default size
	}
	return f
}

// GetAskSize returns ask size as float64
func (t *DeriveTicker) GetAskSize() float64 {
	var f float64
	fmt.Sscanf(t.BestAskAmount, "%f", &f)
	if f == 0 {
		return 1.0 // Default size
	}
	return f
}

// GetIndexPrice returns index price as float64
func (t *DeriveTicker) GetIndexPrice() float64 {
	var f float64
	fmt.Sscanf(t.IndexPrice, "%f", &f)
	return f
}

// DeriveTickerResponse represents the API response
type DeriveTickerResponse struct {
	Result DeriveTicker `json:"result"`
}

// FetchDeriveTicker directly fetches ticker from Derive API, bypassing CCXT
func FetchDeriveTicker(instrumentName string) (*DeriveTicker, error) {
	url := "https://api.lyra.finance/public/get_ticker"
	
	payload := map[string]interface{}{
		"instrument_name": instrumentName,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ticker: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}
	
	var response DeriveTickerResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// Log the raw response for debugging
		return nil, fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
	}
	
	return &response.Result, nil
}