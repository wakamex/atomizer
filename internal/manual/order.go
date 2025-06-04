package manual

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wakamex/atomizer/internal/config"
	"github.com/wakamex/atomizer/internal/types"
)

// OrderConfig contains configuration for manual order placement
type OrderConfig struct {
	Instrument string
	Side       string
	Price      float64
	Amount     float64
}

// OrderService handles manual order placement
type OrderService struct {
	config   *config.Config
	exchange types.MarketMakerExchange
}

// NewOrderService creates a new manual order service
func NewOrderService(cfg *config.Config, exchange types.MarketMakerExchange) *OrderService {
	return &OrderService{
		config:   cfg,
		exchange: exchange,
	}
}

// PlaceOrder submits a manual order to the exchange
func (s *OrderService) PlaceOrder(orderCfg OrderConfig) (string, error) {
	// Validate inputs
	if orderCfg.Instrument == "" {
		return "", fmt.Errorf("instrument is required")
	}
	
	if orderCfg.Side != "buy" && orderCfg.Side != "sell" {
		return "", fmt.Errorf("side must be 'buy' or 'sell'")
	}
	
	if orderCfg.Price <= 0 {
		return "", fmt.Errorf("price must be positive")
	}
	
	if orderCfg.Amount <= 0 {
		return "", fmt.Errorf("amount must be positive")
	}
	
	// Convert to decimal for precision
	price := decimal.NewFromFloat(orderCfg.Price)
	amount := decimal.NewFromFloat(orderCfg.Amount)
	
	log.Printf("Placing %s order: %s %s @ %s", orderCfg.Side, amount, orderCfg.Instrument, price)
	
	// Place the order
	orderID, err := s.exchange.PlaceLimitOrder(orderCfg.Instrument, orderCfg.Side, price, amount)
	if err != nil {
		return "", fmt.Errorf("failed to place order: %w", err)
	}
	
	log.Printf("Order placed successfully! Order ID: %s", orderID)
	
	return orderID, nil
}

// GetOpenOrders retrieves all open orders
func (s *OrderService) GetOpenOrders() ([]types.MarketMakerOrder, error) {
	orders, err := s.exchange.GetOpenOrders()
	if err != nil {
		return nil, fmt.Errorf("failed to get open orders: %w", err)
	}
	
	return orders, nil
}

// CancelOrder cancels a specific order
func (s *OrderService) CancelOrder(orderID string) error {
	return s.exchange.CancelOrder(orderID)
}

// CancelAllOrders cancels all open orders
func (s *OrderService) CancelAllOrders() error {
	orders, err := s.exchange.GetOpenOrders()
	if err != nil {
		return fmt.Errorf("failed to get open orders: %w", err)
	}
	
	for _, order := range orders {
		if err := s.exchange.CancelOrder(order.OrderID); err != nil {
			log.Printf("Failed to cancel order %s: %v", order.OrderID, err)
		}
	}
	
	return nil
}

// RunManualOrder is a standalone function for command-line usage
func RunManualOrder(cfg *config.Config, exchange types.MarketMakerExchange, orderCfg OrderConfig) error {
	service := NewOrderService(cfg, exchange)
	
	// Place the order
	_, err := service.PlaceOrder(orderCfg)
	if err != nil {
		return err
	}
	
	// Wait a bit then check the order
	time.Sleep(2 * time.Second)
	
	orders, err := service.GetOpenOrders()
	if err != nil {
		log.Printf("Failed to get open orders: %v", err)
	} else {
		log.Printf("Open orders: %d", len(orders))
		for _, order := range orders {
			log.Printf("  %s: %s %s %s @ %s", 
				order.OrderID, order.Side, order.Amount, order.Instrument, order.Price)
		}
	}
	
	return nil
}

// ParseOrderFromEnv parses order configuration from environment variables
func ParseOrderFromEnv(defaults OrderConfig) OrderConfig {
	cfg := defaults
	
	// Override from environment if set
	if instrument := getEnv("ORDER_INSTRUMENT"); instrument != "" {
		cfg.Instrument = instrument
	}
	
	if side := getEnv("ORDER_SIDE"); side != "" {
		cfg.Side = side
	}
	
	if priceStr := getEnv("ORDER_PRICE"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			cfg.Price = price
		}
	}
	
	if amountStr := getEnv("ORDER_AMOUNT"); amountStr != "" {
		if amount, err := strconv.ParseFloat(amountStr, 64); err == nil {
			cfg.Amount = amount
		}
	}
	
	return cfg
}

func getEnv(key string) string {
	// This is a helper to make testing easier
	return os.Getenv(key)
}