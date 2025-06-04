package hedging

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/wakamex/atomizer/internal/config"
	"github.com/wakamex/atomizer/internal/types"
	"github.com/shopspring/decimal"
)

// Manager handles hedge execution on exchanges
type Manager struct {
	exchange      types.Exchange
	config        *config.Config
	maxRetries    int
	retryDelayMs  int
}

// NewManager creates a new hedge manager
func NewManager(exchange types.Exchange, cfg *config.Config) *Manager {
	return &Manager{
		exchange:     exchange,
		config:       cfg,
		maxRetries:   3,
		retryDelayMs: 1000,
	}
}

// ExecuteHedge places a hedge order for the given trade
func (m *Manager) ExecuteHedge(ctx context.Context, trade *types.TradeEvent) error {
	log.Printf("Executing hedge for trade %s on %s", trade.ID, m.config.ExchangeName)
	
	// Convert trade to hedge parameters
	hedgeParams, err := m.buildHedgeParams(trade)
	if err != nil {
		return fmt.Errorf("failed to build hedge params: %w", err)
	}
	
	// Get current order book
	orderBook, err := m.getOrderBookWithRetry(ctx, trade)
	if err != nil {
		return fmt.Errorf("failed to get order book: %w", err)
	}
	
	// Calculate hedge price
	hedgePrice := m.calculateHedgePrice(orderBook, hedgeParams.isBuy)
	
	log.Printf("[HedgeManager] Hedge params - isBuy: %v, calculated price: %s", 
		hedgeParams.isBuy, hedgePrice.String())
	
	// Execute hedge with retries
	var lastErr error
	for attempt := 1; attempt <= m.maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		result, err := m.executeSingleHedge(ctx, hedgeParams, hedgePrice)
		if err == nil {
			log.Printf("Hedge successful on attempt %d: OrderID=%s", attempt, result.OrderID)
			
			// Update trade with hedge information
			trade.HedgeOrderID = result.OrderID
			trade.HedgeExchange = m.config.ExchangeName
			
			return nil
		}
		
		lastErr = err
		log.Printf("Hedge attempt %d failed: %v", attempt, err)
		
		if attempt < m.maxRetries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(m.retryDelayMs) * time.Millisecond):
				// Continue to next attempt
			}
		}
	}
	
	return fmt.Errorf("hedge failed after %d attempts: %w", m.maxRetries, lastErr)
}

// hedgeParams contains parameters for hedge execution
type hedgeParams struct {
	instrument string
	quantity   decimal.Decimal
	isBuy      bool
}

// hedgeResult contains the result of a hedge execution
type hedgeResult struct {
	OrderID    string
	Instrument string
	Direction  string
	Quantity   decimal.Decimal
	Price      decimal.Decimal
	ExecutedAt time.Time
}

// buildHedgeParams converts trade to hedge parameters
func (m *Manager) buildHedgeParams(trade *types.TradeEvent) (*hedgeParams, error) {
	// Convert instrument name if needed
	instrument, err := m.convertInstrumentName(trade)
	if err != nil {
		return nil, err
	}
	
	// Hedge direction is opposite of trade
	// If we sold to taker (taker bought), we need to buy to hedge
	isBuy := trade.IsTakerBuy
	
	return &hedgeParams{
		instrument: instrument,
		quantity:   trade.Quantity,
		isBuy:      isBuy,
	}, nil
}

// convertInstrumentName converts Rysk format to exchange format
func (m *Manager) convertInstrumentName(trade *types.TradeEvent) (string, error) {
	// If we already have the instrument name, use it
	if trade.Instrument != "" && strings.Contains(trade.Instrument, "-") {
		return trade.Instrument, nil
	}
	
	// Otherwise convert from trade parameters
	return m.exchange.ConvertToInstrument(
		"ETH", // Default asset, should be mapped from trade.Instrument
		trade.Strike.String(),
		trade.Expiry,
		trade.IsPut,
	)
}

// getOrderBookWithRetry gets order book with retry logic
func (m *Manager) getOrderBookWithRetry(ctx context.Context, trade *types.TradeEvent) (*types.CCXTOrderBook, error) {
	rfq := types.RFQResult{
		Asset:      trade.Instrument,
		Strike:     trade.Strike.String(),
		Expiry:     trade.Expiry,
		IsPut:      trade.IsPut,
		Quantity:   trade.Quantity.String(),
		IsTakerBuy: trade.IsTakerBuy,
	}
	
	var lastErr error
	for i := 0; i < m.maxRetries; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		orderBook, err := m.exchange.GetOrderBook(rfq, "ETH") // Default asset
		if err == nil {
			return &orderBook, nil
		}
		
		lastErr = err
		if i < m.maxRetries-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}
	
	return nil, fmt.Errorf("failed to get order book after %d attempts: %w", m.maxRetries, lastErr)
}

// calculateHedgePrice determines optimal hedge price
func (m *Manager) calculateHedgePrice(orderBook *types.CCXTOrderBook, isBuy bool) decimal.Decimal {
	// During test mode, use defensive pricing
	if m.config.ExchangeTestMode {
		if isBuy && len(orderBook.Asks) > 0 {
			// Buying - use a price above best ask
			bestAsk := decimal.NewFromFloat(orderBook.Asks[0][0])
			return bestAsk.Mul(decimal.NewFromFloat(1.01)) // 1% above ask
		} else if !isBuy && len(orderBook.Bids) > 0 {
			// Selling - use 2x ask price for defensive testing
			if len(orderBook.Asks) > 0 {
				bestAsk := decimal.NewFromFloat(orderBook.Asks[0][0])
				return bestAsk.Mul(decimal.NewFromFloat(2.0))
			}
		}
	}
	
	// Production pricing - use mid-market or slightly aggressive
	if isBuy && len(orderBook.Asks) > 0 {
		return decimal.NewFromFloat(orderBook.Asks[0][0])
	} else if !isBuy && len(orderBook.Bids) > 0 {
		return decimal.NewFromFloat(orderBook.Bids[0][0])
	}
	
	// Fallback price
	return decimal.NewFromFloat(0.05)
}

// executeSingleHedge executes a single hedge attempt
func (m *Manager) executeSingleHedge(ctx context.Context, params *hedgeParams, price decimal.Decimal) (*hedgeResult, error) {
	// Create RFQ confirmation for order placement
	conf := types.RFQConfirmation{
		Price:    price.String(),
		Quantity: params.quantity.String(),
	}
	
	// Place order
	err := m.exchange.PlaceOrder(conf, params.instrument, m.config)
	if err != nil {
		return nil, fmt.Errorf("failed to place hedge order: %w", err)
	}
	
	// Create result
	direction := "buy"
	if !params.isBuy {
		direction = "sell"
	}
	
	result := &hedgeResult{
		OrderID:    fmt.Sprintf("hedge_%d", time.Now().UnixNano()),
		Instrument: params.instrument,
		Direction:  direction,
		Quantity:   params.quantity,
		Price:      price,
		ExecutedAt: time.Now(),
	}
	
	return result, nil
}

// Ensure Manager implements the HedgeManager interface
var _ types.HedgeManager = (*Manager)(nil)