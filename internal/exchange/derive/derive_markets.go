package derive

import (
	"github.com/wakamex/atomizer/internal/exchange/shared"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// DeriveInstrument represents an instrument from Derive/Lyra API
type DeriveInstrument struct {
	InstrumentName string `json:"instrument_name"`
	BaseCurrency   string `json:"base_currency"`
	QuoteCurrency  string `json:"quote_currency"`
	InstrumentType string `json:"instrument_type"`
	IsActive       bool   `json:"is_active"`
	OptionDetails  struct {
		Strike     string `json:"strike"`
		OptionType string `json:"option_type"`
		Expiry     int64  `json:"expiry"`
	} `json:"option_details"`
}

// DerivePaginationResponse represents the pagination info
type DerivePaginationResponse struct {
	NumPages    int `json:"num_pages"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
	TotalCount  int `json:"total_count"`
}

// DeriveInstrumentsResponse represents the API response
type DeriveInstrumentsResponse struct {
	Result struct {
		Instruments []DeriveInstrument       `json:"instruments"`
		Pagination  DerivePaginationResponse `json:"pagination"`
	} `json:"result"`
}

// LoadAllDeriveMarkets fetches all option markets from Derive using pagination
func LoadAllDeriveMarkets() (map[string]DeriveInstrument, error) {
	url := "https://api.lyra.finance/public/get_all_instruments"
	instruments := make(map[string]DeriveInstrument)
	page := 1

	for {
		// Prepare request
		payload := map[string]interface{}{
			"instrument_type": "option",
			"expired":         false,
			"page":            page,
			"page_size":       1000, // Max allowed
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}

		// Make request
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("accept", "application/json")
		req.Header.Set("content-type", "application/json")

		resp, err := shared.NewHTTPClient().Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch instruments: %w", err)
		}
		defer resp.Body.Close()

		// Read response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %w", err)
		}

		// Parse response
		var response DeriveInstrumentsResponse
		if err := json.Unmarshal(body, &response); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		// Add instruments to map
		for _, inst := range response.Result.Instruments {
			instruments[inst.InstrumentName] = inst
		}

		// Log sample instruments to see the format
		if page == 1 && len(response.Result.Instruments) > 0 {
			log.Printf("[Derive] Sample instruments:")
			for i, inst := range response.Result.Instruments {
				if i >= 3 {
					break
				}
				log.Printf("[Derive]   %s (base=%s, strike=%s, type=%s, instType=%s)",
					inst.InstrumentName, inst.BaseCurrency, inst.OptionDetails.Strike, inst.OptionDetails.OptionType, inst.InstrumentType)
			}

			// Debug: print raw JSON for first instrument
			if debugJSON, err := json.MarshalIndent(response.Result.Instruments[0], "", "  "); err == nil {
				log.Printf("[Derive] First instrument raw JSON:\n%s", string(debugJSON))
			}
		}

		log.Printf("[Derive] Loaded page %d/%d (%d instruments)",
			page,
			response.Result.Pagination.NumPages,
			len(response.Result.Instruments))

		// Check if there are more pages
		if page >= response.Result.Pagination.NumPages {
			break
		}
		page++
	}

	return instruments, nil
}

// ConvertDeriveInstrumentToSymbol converts Derive instrument format to CCXT symbol format
func ConvertDeriveInstrumentToSymbol(inst DeriveInstrument) string {
	// Derive uses format like "ETH-27JUN25-4000-C"
	// CCXT expects "ETH/USDC:USDC-25-06-27-4000-C"
	// This is a placeholder - adjust based on actual format requirements
	return fmt.Sprintf("%s/%s:%s", inst.BaseCurrency, inst.QuoteCurrency, inst.InstrumentName)
}
