package main

import (
	"testing"
	"time"
)

func TestHedgeOrder(t *testing.T) {
	cfg := &AppConfig{
		DeribitApiKey:    "PCZQT974",    // Replace with your Deribit testnet API key
		DeribitApiSecret: "rGLJ0WL5micLZzzj9blp5LAen_AgDi0ProJsFEK2woI", // Replace with your Deribit testnet API secret
	}

	tests := []struct {
		name        string
		conf        RFQConfirmation
		underlying  string
		expectError bool
	}{
		{
			name: "hedge taker sell",
			conf: RFQConfirmation{
				QuoteNonce: "test-nonce-2",
				Strike:     "260000000000",
				Expiry:     int(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()+1, 8, 0, 0, 0, time.UTC).Unix()),
				IsPut:      false,
				Quantity:   "2000000000000000000", // 2 ETH in wei
				Price:      "100000000000000000",
				IsTakerBuy: false,
			},
			underlying:  "ETH",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HedgeOrder(tt.conf, tt.underlying, cfg)
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
