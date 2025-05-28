package main

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestArbitrageOrchestrator(t *testing.T) {
	// Create test config
	cfg := &AppConfig{
		MaxPositionSize:  "100",
		MaxDeltaExposure: "50",
		ExchangeName:     "test",
	}

	// Create mock exchange
	mockExchange := &MockExchange{}

	// Create orchestrator
	orchestrator := NewArbitrageOrchestrator(cfg, mockExchange)
	if err := orchestrator.Start(); err != nil {
		t.Fatalf("Failed to start orchestrator: %v", err)
	}
	defer orchestrator.Stop()

	// Test manual trade submission
	t.Run("ManualTradeSubmission", func(t *testing.T) {
		req := ManualTradeRequest{
			Instrument: "ETH-20231225-3000-C",
			Strike:     decimal.NewFromInt(3000),
			Expiry:     time.Now().Add(30 * 24 * time.Hour).Unix(),
			IsPut:      false,
			Quantity:   decimal.NewFromFloat(1.0),
			Price:      decimal.NewFromFloat(0.05),
		}

		trade, err := orchestrator.SubmitManualTrade(req)
		if err != nil {
			t.Errorf("Failed to submit manual trade: %v", err)
		}

		if trade.Source != TradeSourceManual {
			t.Errorf("Expected trade source to be MANUAL, got %s", trade.Source)
		}
	})

	// Test RFQ trade submission
	t.Run("RFQTradeSubmission", func(t *testing.T) {
		rfq := RFQResult{
			ID:         "test-rfq-123",
			Strike:     NewBigInt(3000, 8),
			Expiry:     NewBigInt(time.Now().Add(30*24*time.Hour).Unix(), 0),
			IsPut:      false,
			Quantity:   NewBigInt(1, 18), // 1 ETH in wei
			IsTakerBuy: true,
		}

		trade, err := orchestrator.SubmitRFQTrade(rfq)
		if err != nil {
			t.Errorf("Failed to submit RFQ trade: %v", err)
		}

		if trade.Source != TradeSourceRysk {
			t.Errorf("Expected trade source to be RYSK_RFQ, got %s", trade.Source)
		}
	})
}

func TestRiskManager(t *testing.T) {
	cfg := &AppConfig{
		MaxPositionSize:  "10",
		MaxDeltaExposure: "5",
	}

	riskManager := NewRiskManager(cfg)

	t.Run("ValidateTradeWithinLimits", func(t *testing.T) {
		trade := &TradeEvent{
			Instrument: "ETH-20231225-3000-C",
			Quantity:   decimal.NewFromFloat(5),
			IsTakerBuy: false,
		}

		err := riskManager.ValidateTrade(trade)
		if err != nil {
			t.Errorf("Trade should be valid: %v", err)
		}
	})

	t.Run("ValidateTradeExceedsLimits", func(t *testing.T) {
		trade := &TradeEvent{
			Instrument: "ETH-20231225-3000-C",
			Quantity:   decimal.NewFromFloat(15), // Exceeds max position size of 10
			IsTakerBuy: false,
		}

		err := riskManager.ValidateTrade(trade)
		if err == nil {
			t.Error("Trade should be rejected for exceeding position limits")
		}
	})
}

func TestHedgeManager(t *testing.T) {
	cfg := &AppConfig{
		ExchangeName: "test",
	}

	mockExchange := &MockExchange{}
	hedgeManager := NewHedgeManager(mockExchange, cfg)

	t.Run("ExecuteHedge", func(t *testing.T) {
		trade := &TradeEvent{
			ID:         "test-trade-123",
			Instrument: "ETH-20231225-3000-C",
			Strike:     decimal.NewFromInt(3000),
			Expiry:     time.Now().Add(30 * 24 * time.Hour).Unix(),
			IsPut:      false,
			Quantity:   decimal.NewFromFloat(1.0),
			Price:      decimal.NewFromFloat(0.05),
			IsTakerBuy: true, // Taker buys from us, so we need to buy on exchange
		}

		result, err := hedgeManager.ExecuteHedge(trade)
		if err != nil {
			t.Errorf("Failed to execute hedge: %v", err)
		}

		if result.Direction != "buy" {
			t.Errorf("Expected hedge direction to be buy, got %s", result.Direction)
		}
	})
}

// Helper functions
func NewBigInt(value int64, exp int) *big.Int {
	if exp > 0 {
		multiplier := big.NewInt(10)
		multiplier.Exp(multiplier, big.NewInt(int64(exp)), nil)
		result := big.NewInt(value)
		result.Mul(result, multiplier)
		return result
	}
	return big.NewInt(value)
}

// MockExchange for testing
type MockExchange struct{}

func (m *MockExchange) GetOrderBook(req RFQResult, asset string) (CCXTOrderBook, error) {
	return CCXTOrderBook{
		Bids: [][]float64{{2999, 10}, {2998, 20}},
		Asks: [][]float64{{3001, 10}, {3002, 20}},
	}, nil
}

func (m *MockExchange) PlaceOrder(conf RFQConfirmation, underlying string, cfg *AppConfig) error {
	return nil
}

func (m *MockExchange) ConvertToInstrument(asset, strike string, expiry int64, isPut bool) (string, error) {
	optType := "C"
	if isPut {
		optType = "P"
	}
	expiryTime := time.Unix(expiry, 0)
	return fmt.Sprintf("%s-%s-%s-%s", asset, expiryTime.Format("20060102"), strike, optType), nil
}