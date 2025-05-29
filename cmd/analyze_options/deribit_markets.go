package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	
	"github.com/shopspring/decimal"
)

// DeribitInstrument represents an instrument from Deribit API
type DeribitInstrument struct {
	InstrumentName     string  `json:"instrument_name"`
	BaseCurrency       string  `json:"base_currency"`
	QuoteCurrency      string  `json:"quote_currency"`
	IsActive           bool    `json:"is_active"`
	Kind               string  `json:"kind"`
	OptionType         string  `json:"option_type"`
	Strike             float64 `json:"strike"`
	ExpirationTimestamp int64  `json:"expiration_timestamp"`
	SettlementPeriod   string  `json:"settlement_period"`
	ContractSize       float64 `json:"contract_size"`
	MinTradeAmount     float64 `json:"min_trade_amount"`
}

// DeribitTickerResponse represents ticker data from Deribit
type DeribitTickerResponse struct {
	Result DeribitTicker `json:"result"`
}

type DeribitTicker struct {
	InstrumentName    string  `json:"instrument_name"`
	BestBidPrice      float64 `json:"best_bid_price"`
	BestBidAmount     float64 `json:"best_bid_amount"`
	BestAskPrice      float64 `json:"best_ask_price"`
	BestAskAmount     float64 `json:"best_ask_amount"`
	LastPrice         float64 `json:"last_price"`
	MarkPrice         float64 `json:"mark_price"`
	IndexPrice        float64 `json:"index_price"`
	OpenInterest      float64 `json:"open_interest"`
	Stats             struct {
		Volume float64 `json:"volume"`
		High   float64 `json:"high"`
		Low    float64 `json:"low"`
	} `json:"stats"`
	Greeks            DeribitGreeks `json:"greeks"`
}

type DeribitGreeks struct {
	Delta float64 `json:"delta"`
	Gamma float64 `json:"gamma"`
	Theta float64 `json:"theta"`
	Vega  float64 `json:"vega"`
	Rho   float64 `json:"rho"`
}

// DeribitClient handles communication with Deribit API
type DeribitClient struct {
	baseURL    string
	apiKey     string
	apiSecret  string
	testMode   bool
}

// NewDeribitClient creates a new Deribit client
func NewDeribitClient(testMode bool) *DeribitClient {
	baseURL := "https://www.deribit.com/api/v2"
	if testMode {
		baseURL = "https://test.deribit.com/api/v2"
	}
	
	return &DeribitClient{
		baseURL:   baseURL,
		apiKey:    os.Getenv("DERIBIT_API_KEY"),
		apiSecret: os.Getenv("DERIBIT_API_SECRET"),
		testMode:  testMode,
	}
}

// LoadAllDeribitMarkets fetches all option markets from Deribit
func (dc *DeribitClient) LoadAllDeribitMarkets() (map[string]DeriveInstrument, error) {
	instruments := make(map[string]DeriveInstrument)
	
	// Fetch for each currency
	currencies := []string{"BTC", "ETH"}
	
	for _, currency := range currencies {
		fmt.Printf("Fetching %s options from Deribit...\n", currency)
		
		url := fmt.Sprintf("%s/public/get_instruments?currency=%s&kind=option&expired=false", dc.baseURL, currency)
		
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch %s instruments: %w", currency, err)
		}
		defer resp.Body.Close()
		
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}
		
		var response struct {
			Result []DeribitInstrument `json:"result"`
		}
		
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		
		// Convert Deribit instruments to common format
		for _, inst := range response.Result {
			if !inst.IsActive {
				continue
			}
			
			deriveInst := DeriveInstrument{
				InstrumentName: inst.InstrumentName,
				BaseCurrency:   inst.BaseCurrency,
				QuoteCurrency:  inst.QuoteCurrency,
				InstrumentType: "option",
				IsActive:       inst.IsActive,
			}
			
			// Set option details
			deriveInst.OptionDetails.Strike = fmt.Sprintf("%.0f", inst.Strike)
			// Convert Deribit's "call"/"put" to "C"/"P"
			if inst.OptionType == "call" {
				deriveInst.OptionDetails.OptionType = "C"
			} else if inst.OptionType == "put" {
				deriveInst.OptionDetails.OptionType = "P"
			} else {
				deriveInst.OptionDetails.OptionType = inst.OptionType
			}
			deriveInst.OptionDetails.Expiry = inst.ExpirationTimestamp / 1000 // Convert to seconds
			
			// Debug first few ETH options
			if currency == "ETH" && inst.OptionType == "call" && len(instruments) < 3 {
				expiryTime := time.Unix(deriveInst.OptionDetails.Expiry, 0)
				log.Printf("[Deribit] Sample %s option: %s, expiry: %s (%d)", 
					currency, inst.InstrumentName, expiryTime.Format("2006-01-02 15:04"), deriveInst.OptionDetails.Expiry)
			}
			
			instruments[inst.InstrumentName] = deriveInst
		}
		
		log.Printf("[Deribit] Loaded %d %s options", len(response.Result), currency)
	}
	
	return instruments, nil
}

// FetchDeribitTicker fetches ticker data for a single instrument from Deribit
func (dc *DeribitClient) FetchDeribitTicker(instrumentName string) (*TickerResult, error) {
	url := fmt.Sprintf("%s/public/ticker?instrument_name=%s", dc.baseURL, instrumentName)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	
	var deribitResp DeribitTickerResponse
	if err := json.NewDecoder(resp.Body).Decode(&deribitResp); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w", err)
	}
	
	ticker := deribitResp.Result
	
	// Convert to common TickerResult format
	result := &TickerResult{
		InstrumentName: ticker.InstrumentName,
		IsActive:      true,
		MarkPrice:     decimal.NewFromFloat(ticker.MarkPrice),
		BestBidPrice:  decimal.NewFromFloat(ticker.BestBidPrice),
		BestBidAmount: decimal.NewFromFloat(ticker.BestBidAmount),
		BestAskPrice:  decimal.NewFromFloat(ticker.BestAskPrice),
		BestAskAmount: decimal.NewFromFloat(ticker.BestAskAmount),
		IndexPrice:    decimal.NewFromFloat(ticker.IndexPrice),
		Timestamp:     time.Now().Unix() * 1000,
	}
	
	// Parse instrument name to get option details
	// Deribit format: BTC-27DEC24-100000-C
	parts := strings.Split(instrumentName, "-")
	if len(parts) == 4 {
		result.OptionDetails.Strike = parts[2]
		result.OptionDetails.OptionType = parts[3]
		
		// Parse expiry date
		expiryStr := parts[1]
		// Convert format like "27DEC24" to timestamp
		expiry, err := parseDeribitExpiry(expiryStr)
		if err == nil {
			result.OptionDetails.Expiry = expiry.Unix()
		}
	}
	
	// Set stats
	result.Stats.ContractVolume = decimal.NewFromFloat(ticker.Stats.Volume)
	result.Stats.High = decimal.NewFromFloat(ticker.Stats.High)
	result.Stats.Low = decimal.NewFromFloat(ticker.Stats.Low)
	
	// Set option pricing
	result.OptionPricing.Delta = decimal.NewFromFloat(ticker.Greeks.Delta)
	result.OptionPricing.Gamma = decimal.NewFromFloat(ticker.Greeks.Gamma)
	result.OptionPricing.Theta = decimal.NewFromFloat(ticker.Greeks.Theta)
	result.OptionPricing.Vega = decimal.NewFromFloat(ticker.Greeks.Vega)
	result.OptionPricing.Rho = decimal.NewFromFloat(ticker.Greeks.Rho)
	result.OptionPricing.MarkPrice = decimal.NewFromFloat(ticker.MarkPrice)
	
	// Set open interest
	result.OpenInterest = make(map[string][]OIData)
	result.OpenInterest["total"] = []OIData{
		{
			CurrentOpenInterest: decimal.NewFromFloat(ticker.OpenInterest),
			ManagerCurrency:    parts[0], // BTC or ETH
		},
	}
	
	return result, nil
}

// parseDeribitExpiry parses Deribit expiry format (e.g., "27DEC24") to time.Time
func parseDeribitExpiry(expiryStr string) (time.Time, error) {
	// Deribit uses format like "27DEC24" which means 27 December 2024
	// Options expire at 08:00 UTC
	
	if len(expiryStr) < 7 {
		return time.Time{}, fmt.Errorf("invalid expiry format: %s", expiryStr)
	}
	
	day := expiryStr[:2]
	month := expiryStr[2:5]
	year := "20" + expiryStr[5:7]
	
	// Convert month abbreviation to number
	monthMap := map[string]string{
		"JAN": "01", "FEB": "02", "MAR": "03", "APR": "04",
		"MAY": "05", "JUN": "06", "JUL": "07", "AUG": "08",
		"SEP": "09", "OCT": "10", "NOV": "11", "DEC": "12",
	}
	
	monthNum, ok := monthMap[month]
	if !ok {
		return time.Time{}, fmt.Errorf("invalid month: %s", month)
	}
	
	// Parse as UTC time at 08:00
	timeStr := fmt.Sprintf("%s-%s-%s 08:00:00", year, monthNum, day)
	return time.Parse("2006-01-02 15:04:05", timeStr)
}