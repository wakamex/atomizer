package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Direction represents order direction
type Direction int

const (
	Buy Direction = iota
	Sell
)

func (d Direction) Sign() decimal.Decimal {
	if d == Buy {
		return decimal.NewFromInt(1)
	}
	return decimal.NewFromInt(-1)
}

func (d Direction) String() string {
	if d == Buy {
		return "buy"
	}
	return "sell"
}

// OrderType represents the type of order
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

// TimeInForce represents order time in force
type TimeInForce string

const (
	TimeInForcePostOnly TimeInForce = "post_only"
	TimeInForceGTC      TimeInForce = "gtc"
	TimeInForceIOC      TimeInForce = "ioc"
)

// TickerData represents ticker information
type TickerData struct {
	InstrumentName string
	MarkPrice      decimal.Decimal
	MinimumAmount  decimal.Decimal
	AmountStep     decimal.Decimal
	TickSize       decimal.Decimal
	OptionPricing  *OptionPricing
	OptionDetails  *OptionDetails
}

// OptionPricing contains option Greeks
type OptionPricing struct {
	Delta decimal.Decimal
	Gamma decimal.Decimal
	Vega  decimal.Decimal
	Theta decimal.Decimal
}

// OptionDetails contains option specifics
type OptionDetails struct {
	Expiry int64 // Unix timestamp
}

// OrderbookData represents order book state
type OrderbookData struct {
	Bids [][]decimal.Decimal // [price, amount]
	Asks [][]decimal.Decimal // [price, amount]
}

// GammaPosition represents a trading position for gamma hedging
type GammaPosition struct {
	InstrumentName string
	Amount         decimal.Decimal
}

// OrderArgs represents order parameters
type OrderArgs struct {
	Amount      decimal.Decimal
	LimitPrice  decimal.Decimal
	Direction   Direction
	TimeInForce TimeInForce
	OrderType   OrderType
	Label       string
	MMP         bool
}

// MarketData interface for market state
type MarketData interface {
	GetTicker(instrumentName string) (*TickerData, error)
	GetOrderbookExcludeMyOrders(instrumentName string) (*OrderbookData, error)
	GetOrders(instrumentName string) []GammaOrder
	IterPositions() []GammaPosition
}

// GammaOrder represents an order for gamma hedging
type GammaOrder struct {
	OrderID   string
	Direction Direction
	Price     decimal.Decimal
	Amount    decimal.Decimal
	Label     string
	Status    string
}

// WsClient interface for WebSocket client
type WsClient interface {
	Login() error
	EnableCancelOnDisconnect() error
	SendOrder(ticker *TickerData, subaccountID int64, args OrderArgs) error
	SendReplace(ticker *TickerData, subaccountID int64, cancelID uuid.UUID, args OrderArgs) error
	CancelAll(subaccountID int64) error
	CancelByLabel(subaccountID int64, label string) error
	Ping() error
}

// GammaDDHAlgo implements gamma dynamic delta hedging
type GammaDDHAlgo struct {
	SubaccountID int64
	PerpName     string
	MaxAbsDelta  decimal.Decimal
	MaxAbsSpread decimal.Decimal
	ActionWaitMS uint64
	PriceTol     decimal.Decimal
	AmountTol    decimal.Decimal
	MidPriceTol  decimal.Decimal
}

// GammaDDHState holds the current state for the algorithm
type GammaDDHState struct {
	PerpTicker  *TickerData
	PerpBestBid *decimal.Decimal
	PerpBestAsk *decimal.Decimal
	NetDelta    decimal.Decimal
	NetGamma    decimal.Decimal
	BidIDs      []OrderInfo
	AskIDs      []OrderInfo
}

// OrderInfo holds order identification and details
type OrderInfo struct {
	ID     string
	Price  decimal.Decimal
	Amount decimal.Decimal
}

// GetIndexChangeToTargetDelta calculates index price change needed to reach target delta
func (s *GammaDDHState) GetIndexChangeToTargetDelta(targetDelta, maxAbsSpread decimal.Decimal) (decimal.Decimal, decimal.Decimal) {
	zero := decimal.Zero

	if targetDelta.GreaterThan(zero) && s.NetDelta.GreaterThan(targetDelta) {
		return zero, s.NetDelta
	}
	if targetDelta.LessThan(zero) && s.NetDelta.LessThan(targetDelta) {
		return zero, s.NetDelta
	}
	if s.NetGamma.LessThanOrEqual(zero) {
		return zero, s.NetDelta
	}

	dIndex := targetDelta.Sub(s.NetDelta).Div(s.NetGamma)
	if dIndex.Abs().LessThan(maxAbsSpread) {
		return dIndex, targetDelta
	}

	sign := decimal.NewFromInt(1)
	if dIndex.IsNegative() {
		sign = decimal.NewFromInt(-1)
	}
	dIndex = maxAbsSpread.Mul(sign)
	targetDelta = s.NetDelta.Add(s.NetGamma.Mul(dIndex))
	return dIndex, targetDelta
}

// GetSmoothMid returns the smooth mid-price
func (s *GammaDDHState) GetSmoothMid(midPriceTol decimal.Decimal) decimal.Decimal {
	if s.PerpBestBid != nil && s.PerpBestAsk != nil {
		spread := s.PerpBestAsk.Sub(*s.PerpBestBid)
		if spread.LessThan(midPriceTol) {
			return s.PerpBestBid.Add(*s.PerpBestAsk).Div(decimal.NewFromInt(2))
		}
	}
	return s.PerpTicker.MarkPrice
}

// HedgerAction performs hedging action in a given direction
func (a *GammaDDHAlgo) HedgerAction(state *GammaDDHState, client WsClient, direction Direction) error {
	ticker := state.PerpTicker
	targetDelta := a.MaxAbsDelta.Neg().Mul(direction.Sign())
	dIndex, delta := state.GetIndexChangeToTargetDelta(targetDelta, a.MaxAbsSpread)
	amount := delta.Neg().Mul(direction.Sign())
	limitPrice := dIndex.Add(ticker.MarkPrice)

	// Prevent post-only cross errors by clipping limit price to inside BBO
	if direction == Buy && state.PerpBestAsk != nil {
		maxPrice := state.PerpBestAsk.Sub(ticker.TickSize)
		if limitPrice.GreaterThan(maxPrice) {
			limitPrice = maxPrice
		}
	} else if direction == Sell && state.PerpBestBid != nil {
		minPrice := state.PerpBestBid.Add(ticker.TickSize)
		if limitPrice.LessThan(minPrice) {
			limitPrice = minPrice
		}
	}

	// Clip to smooth mid if too aggressive
	midPrice := state.GetSmoothMid(a.MidPriceTol)
	if direction == Buy && limitPrice.GreaterThan(midPrice) {
		limitPrice = midPrice
	} else if direction == Sell && limitPrice.LessThan(midPrice) {
		limitPrice = midPrice
	}

	label := fmt.Sprintf("gamma-ddh-%s", direction)
	var openIDs []OrderInfo
	if direction == Buy {
		openIDs = state.BidIDs
	} else {
		openIDs = state.AskIDs
	}

	if amount.LessThan(ticker.MinimumAmount) {
		log.Printf("Amount calculated too small: %s", amount)
		if len(openIDs) >= 1 {
			return client.CancelByLabel(a.SubaccountID, label)
		}
		return nil
	}

	// Round values
	amountScale := int32(ticker.AmountStep.Exponent() * -1)
	priceScale := int32(ticker.TickSize.Exponent() * -1)

	orderArgs := OrderArgs{
		Amount:      amount.Round(amountScale),
		LimitPrice:  limitPrice.Round(priceScale),
		Direction:   direction,
		TimeInForce: TimeInForcePostOnly,
		OrderType:   OrderTypeLimit,
		Label:       label,
		MMP:         false,
	}

	switch len(openIDs) {
	case 0:
		return client.SendOrder(ticker, a.SubaccountID, orderArgs)
	case 1:
		isPriceNew := openIDs[0].Price.Sub(limitPrice).Abs().GreaterThan(a.PriceTol)
		isAmountNew := openIDs[0].Amount.Sub(amount).Abs().GreaterThan(a.AmountTol)
		if !isPriceNew && !isAmountNew {
			return client.Ping()
		}
		cancelID, err := uuid.Parse(openIDs[0].ID)
		if err != nil {
			return err
		}
		return client.SendReplace(ticker, a.SubaccountID, cancelID, orderArgs)
	default:
		log.Printf("Open orders: %+v", openIDs)
		return client.CancelAll(a.SubaccountID)
	}
}

// GetVariablePct calculates variable percentage based on time to expiry
func (a *GammaDDHAlgo) GetVariablePct(expiry int64) decimal.Decimal {
	secToExpiry := float64(expiry - time.Now().Unix())
	const feedTWAPSec = 1800.0

	if secToExpiry < feedTWAPSec {
		pct := decimal.NewFromFloat(secToExpiry / feedTWAPSec)
		if pct.LessThan(decimal.Zero) {
			return decimal.Zero
		}
		return pct
	}
	return decimal.NewFromInt(1)
}

// FilterOpenIDs filters open orders by direction
func FilterOpenIDs(orders []GammaOrder, direction Direction) []OrderInfo {
	var result []OrderInfo
	for _, order := range orders {
		if order.Direction == direction && order.Status == "open" {
			result = append(result, OrderInfo{
				ID:     order.OrderID,
				Price:  order.Price,
				Amount: order.Amount,
			})
		}
	}
	return result
}

// GetAlgoState retrieves current algorithm state
func (a *GammaDDHAlgo) GetAlgoState(market MarketData) (*GammaDDHState, error) {
	perpTicker, err := market.GetTicker(a.PerpName)
	if err != nil {
		return nil, fmt.Errorf("perp ticker not found: %w", err)
	}

	perpOrderbook, err := market.GetOrderbookExcludeMyOrders(a.PerpName)
	if err != nil {
		return nil, fmt.Errorf("perp orderbook not found: %w", err)
	}

	orders := market.GetOrders(a.PerpName)

	state := &GammaDDHState{
		PerpTicker: perpTicker,
		NetDelta:   decimal.Zero,
		NetGamma:   decimal.Zero,
		BidIDs:     FilterOpenIDs(orders, Buy),
		AskIDs:     FilterOpenIDs(orders, Sell),
	}

	if len(perpOrderbook.Bids) > 0 {
		state.PerpBestBid = &perpOrderbook.Bids[0][0]
	}
	if len(perpOrderbook.Asks) > 0 {
		state.PerpBestAsk = &perpOrderbook.Asks[0][0]
	}

	// Calculate net delta and gamma from positions
	for _, position := range market.IterPositions() {
		log.Printf("Position: %s, %s", position.InstrumentName, position.Amount)
		ticker, err := market.GetTicker(position.InstrumentName)
		if err != nil {
			continue
		}

		if ticker.OptionPricing != nil && ticker.OptionDetails != nil {
			pctVar := a.GetVariablePct(ticker.OptionDetails.Expiry)
			state.NetGamma = state.NetGamma.Add(
				position.Amount.Mul(ticker.OptionPricing.Gamma).Mul(pctVar))
			state.NetDelta = state.NetDelta.Add(
				position.Amount.Mul(ticker.OptionPricing.Delta).Mul(pctVar))
		} else if ticker.InstrumentName == a.PerpName {
			state.NetDelta = state.NetDelta.Add(position.Amount)
		}
	}

	log.Printf("Net Delta: %s, Net Gamma: %s", state.NetDelta, state.NetGamma)
	return state, nil
}

// StartHedger starts the hedging algorithm
func (a *GammaDDHAlgo) StartHedger(ctx context.Context, state MarketData, client WsClient) error {
	log.Println("Starting hedger task")

	err := client.Login()
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	err = client.EnableCancelOnDisconnect()
	if err != nil {
		return fmt.Errorf("enable cancel on disconnect failed: %w", err)
	}

	ticker := time.NewTicker(time.Duration(a.ActionWaitMS) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			algoState, err := a.GetAlgoState(state)
			if err != nil {
				return err
			}

			// Execute bid and ask actions concurrently
			var wg sync.WaitGroup
			var bidErr, askErr error

			wg.Add(2)
			go func() {
				defer wg.Done()
				bidErr = a.HedgerAction(algoState, client, Buy)
			}()
			go func() {
				defer wg.Done()
				askErr = a.HedgerAction(algoState, client, Sell)
			}()
			wg.Wait()

			if bidErr != nil {
				return bidErr
			}
			if askErr != nil {
				return askErr
			}
		}
	}
}
