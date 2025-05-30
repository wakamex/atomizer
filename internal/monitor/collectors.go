package monitor

import (
	"net/http"
	"strings"
	"time"
)

// BaseCollector provides common functionality for all collectors
type BaseCollector struct {
	client  *http.Client
	baseURL string
}

// NewHTTPClient creates a configured HTTP client for collectors
func NewHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

// MatchesPattern checks if a symbol matches any of the provided patterns
func MatchesPattern(symbol string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}
	
	symbolUpper := strings.ToUpper(symbol)
	for _, pattern := range patterns {
		patternUpper := strings.ToUpper(pattern)
		
		// Exact match
		if symbolUpper == patternUpper {
			return true
		}
		
		// Contains match (for patterns like "ETH")
		if strings.Contains(symbolUpper, patternUpper) {
			return true
		}
	}
	return false
}

// FilterInstruments filters a list of instruments based on patterns
func FilterInstruments(instruments []string, patterns []string) []string {
	if len(patterns) == 0 {
		return instruments
	}
	
	filtered := []string{}
	for _, inst := range instruments {
		if MatchesPattern(inst, patterns) {
			filtered = append(filtered, inst)
		}
	}
	
	return filtered
}