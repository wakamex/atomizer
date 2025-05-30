package monitor

import (
	"testing"
)

func TestInstrumentConverter(t *testing.T) {
	ic := NewInstrumentConverter()
	
	tests := []struct {
		name     string
		input    string
		exchange string
		expected string
	}{
		// Derive format to Deribit
		{
			name:     "Derive to Deribit - Call",
			input:    "ETH-20250531-2700-C",
			exchange: "deribit",
			expected: "ETH-31MAY25-2700-C",
		},
		{
			name:     "Derive to Deribit - Put",
			input:    "BTC-20241225-50000-P",
			exchange: "deribit",
			expected: "BTC-25DEC24-50000-P",
		},
		
		// Deribit format to Derive
		{
			name:     "Deribit to Derive - Call",
			input:    "ETH-31MAY25-2700-C",
			exchange: "derive",
			expected: "ETH-20250531-2700-C",
		},
		{
			name:     "Deribit to Derive - Put",
			input:    "BTC-25DEC24-50000-P",
			exchange: "derive",
			expected: "BTC-20241225-50000-P",
		},
		
		// Already in correct format
		{
			name:     "Derive format for Derive",
			input:    "ETH-20250531-2700-C",
			exchange: "derive",
			expected: "ETH-20250531-2700-C",
		},
		{
			name:     "Deribit format for Deribit",
			input:    "ETH-31MAY25-2700-C",
			exchange: "deribit",
			expected: "ETH-31MAY25-2700-C",
		},
		
		// Perpetuals (no conversion)
		{
			name:     "ETH Perpetual",
			input:    "ETH-PERPETUAL",
			exchange: "deribit",
			expected: "ETH-PERPETUAL",
		},
		{
			name:     "BTC Perp",
			input:    "BTC-PERP",
			exchange: "derive",
			expected: "BTC-PERP",
		},
		
		// Simple patterns (no conversion)
		{
			name:     "Simple ETH pattern",
			input:    "ETH",
			exchange: "deribit",
			expected: "ETH",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ic.ConvertForExchange(tt.input, tt.exchange)
			if result != tt.expected {
				t.Errorf("ConvertForExchange(%q, %q) = %q, want %q", 
					tt.input, tt.exchange, result, tt.expected)
			}
		})
	}
}

func TestConvertInstrumentList(t *testing.T) {
	ic := NewInstrumentConverter()
	
	// Test converting mixed formats for Deribit
	instruments := []string{
		"ETH-20250531-2700-C",  // Derive format
		"ETH-31MAY25-3000-C",   // Already Deribit format
		"ETH-PERPETUAL",        // No conversion needed
		"ETH",                  // Simple pattern
	}
	
	expected := []string{
		"ETH-31MAY25-2700-C",
		"ETH-31MAY25-3000-C",
		"ETH-PERPETUAL",
		"ETH",
	}
	
	result := ic.ConvertInstrumentList(instruments, "deribit")
	
	if len(result) != len(expected) {
		t.Fatalf("Expected %d instruments, got %d", len(expected), len(result))
	}
	
	for i, exp := range expected {
		if result[i] != exp {
			t.Errorf("Instrument %d: expected %q, got %q", i, exp, result[i])
		}
	}
}