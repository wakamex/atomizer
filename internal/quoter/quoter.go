package quoter

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wakamex/atomizer/internal/config"
	"github.com/wakamex/atomizer/internal/types"
	"github.com/wakamex/rysk-v12-cli/ryskcore"
)


// MakeQuote generates a signed quote based on the given RFQ request
func MakeQuote(req types.RFQResult, underlying string, originalRfqID string, cfg *config.Config, exchange types.Exchange) (ryskcore.Quote, error) {
	quote, _, err := getExchangeQuote(req, underlying, exchange)
	if err != nil {
		return ryskcore.Quote{}, fmt.Errorf("failed to get quote from exchange: %w", err)
	}

	// Convert price from float to string
	// Rysk uses 10^6 precision, so multiply by 10^6
	priceInPrecision := quote.Price * 1e6
	priceStr := strconv.FormatUint(uint64(priceInPrecision), 10)

	// Validate quantity format (should be in wei)
	_, ok := new(big.Int).SetString(req.Quantity, 10)
	if !ok {
		return ryskcore.Quote{}, fmt.Errorf("failed to parse quantity: %s", req.Quantity)
	}

	// Get asset address
	assetAddress := req.Asset
	if strings.HasPrefix(assetAddress, "0x") && len(assetAddress) == 42 {
		// Valid address, use as is
	} else {
		// Try to convert asset name to address if available
		// This would require a mapping which might not be available
		log.Printf("Warning: asset %s might not be a valid address", assetAddress)
	}

	// Generate nonce (timestamp in microseconds)
	nonce := strconv.FormatInt(time.Now().UnixMicro(), 10)

	// Calculate quote validity (current time + duration)
	validUntil := time.Now().Unix() + cfg.QuoteValidDurationSeconds

	// Create the quote
	ryskQuote := ryskcore.Quote{
		AssetAddress: assetAddress,
		ChainID:      req.ChainID,
		Expiry:       req.Expiry,
		IsPut:        req.IsPut,
		IsTakerBuy:   req.IsTakerBuy,
		Maker:        cfg.MakerAddress,
		Nonce:        nonce,
		Price:        priceStr,
		Quantity:     req.Quantity, // Already in string format
		Strike:       req.Strike,
		ValidUntil:   validUntil,
	}

	// Sign the quote using EIP-712
	messageHash, _, err := ryskcore.CreateQuoteMessage(ryskQuote)
	if err != nil {
		return ryskcore.Quote{}, fmt.Errorf("failed to create quote message: %w", err)
	}

	// Convert private key to hex string for signing
	privateKeyBytes := crypto.FromECDSA(cfg.ParsedPrivateKey)
	privateKeyHex := fmt.Sprintf("%x", privateKeyBytes)
	
	signature, err := ryskcore.Sign(messageHash, privateKeyHex)
	if err != nil {
		return ryskcore.Quote{}, fmt.Errorf("failed to sign quote: %w", err)
	}

	ryskQuote.Signature = signature
	return ryskQuote, nil
}

// getExchangeQuote fetches a quote from the exchange based on the RFQ request
func getExchangeQuote(req types.RFQResult, underlying string, exchange types.Exchange) (Quote, float64, error) {
	// Get the order book
	orderBook, err := exchange.GetOrderBook(req, underlying)
	if err != nil {
		return Quote{}, 0, fmt.Errorf("failed to get order book: %w", err)
	}

	// Calculate the price including slippage
	price, err := getPriceInclSlippage(orderBook, req)
	if err != nil {
		return Quote{}, 0, fmt.Errorf("failed to calculate price with slippage: %w", err)
	}

	// Calculate APR
	expiryTime := time.Unix(req.Expiry, 0)
	daysToExpiry := expiryTime.Sub(time.Now()).Hours() / 24
	strikeFloat, _ := strconv.ParseFloat(req.Strike, 64)
	
	apr := CalculateAPR(price, strikeFloat, orderBook.Index, daysToExpiry, req.IsPut)

	return Quote{
		Price: price,
		APR:   apr,
	}, apr, nil
}

// Quote represents a price quote with APR
type Quote struct {
	Price float64
	APR   float64
}

// getPriceInclSlippage calculates the execution price including slippage based on order book depth
func getPriceInclSlippage(ob types.CCXTOrderBook, req types.RFQResult) (float64, error) {
	// Parse quantity
	quantityWei, ok := new(big.Int).SetString(req.Quantity, 10)
	if !ok {
		return 0, fmt.Errorf("failed to parse quantity: %s", req.Quantity)
	}
	
	// Convert from wei (18 decimals) to float
	quantityFloat := new(big.Float).SetInt(quantityWei)
	divisor := new(big.Float).SetFloat64(1e18)
	quantityFloat.Quo(quantityFloat, divisor)
	
	quantity, _ := quantityFloat.Float64()

	// Determine which side of the order book to use
	var orders [][]float64
	if req.IsTakerBuy {
		// Taker wants to buy, so we look at asks (we sell to them)
		orders = ob.Asks
	} else {
		// Taker wants to sell, so we look at bids (we buy from them)
		orders = ob.Bids
	}

	if len(orders) == 0 {
		return 0, errors.New("order book is empty")
	}

	// Calculate weighted average price based on quantity
	remainingQty := quantity
	totalCost := 0.0

	for _, order := range orders {
		price := order[0]
		size := order[1]

		if remainingQty <= 0 {
			break
		}

		fillQty := math.Min(remainingQty, size)
		totalCost += fillQty * price
		remainingQty -= fillQty
	}

	if remainingQty > 0 {
		// Not enough liquidity in the order book
		return 0, fmt.Errorf("insufficient liquidity: need %f more", remainingQty)
	}

	// Calculate average price
	avgPrice := totalCost / quantity

	return avgPrice, nil
}

// CalculateAPR calculates the annualized percentage rate for an option
func CalculateAPR(optionPrice, strike, spot, daysToExpiry float64, isPut bool) float64 {
	if daysToExpiry <= 0 {
		return 0
	}

	// Calculate intrinsic value
	var intrinsicValue float64
	if isPut {
		intrinsicValue = math.Max(strike-spot, 0)
	} else {
		intrinsicValue = math.Max(spot-strike, 0)
	}

	// Time value = Option Price - Intrinsic Value
	timeValue := optionPrice - intrinsicValue
	if timeValue < 0 {
		timeValue = 0
	}

	// Calculate return based on option type
	var returnPercent float64
	if isPut {
		// For puts: return = time value / strike
		if strike > 0 {
			returnPercent = (timeValue / strike) * 100
		}
	} else {
		// For calls: return = time value / spot
		if spot > 0 {
			returnPercent = (timeValue / spot) * 100
		}
	}

	// Annualize the return
	// APR = (return / days) * 365
	apr := (returnPercent / daysToExpiry) * 365

	return apr
}