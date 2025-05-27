package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"
	"time"

	"github.com/wakamex/rysk-v12-cli/ryskcore"
)

// the purpose of this module is to take in a given option request confirmation, return a quote for that option
// then seperately if  all good execute the trade which triggers the solidity side and then opens an offsetting position
// on deribit
// / the underlying must be a string that is either "ETH", "BTC" or "SOL" (essentially the primary asset, not the LST)
// / the expiry must be provided in the format "DDMMMYY" (e.g. "31JAN25") (note that on deribit if the first number is 0, it will cut it e.g. 6JAN25)
// / quantity in normal numbers
func MakeQuote(req RFQResult, underlying string, originalRfqID string, cfg *AppConfig, exchange Exchange) (ryskcore.Quote, error) {
	quote, _, err := getExchangeQuote(req, underlying, exchange)
	if err != nil {
		return ryskcore.Quote{}, err
	}
	quoteBigInt, ok := new(big.Int).SetString(quote, 10)
	if !ok {
		return ryskcore.Quote{}, fmt.Errorf("bad quote conversion: %s", quote)
	}
	maker := cfg.MakerAddress
	privateKey := cfg.PrivateKey
	finalQuote := ryskcore.Quote{
		AssetAddress: req.Asset,
		ChainID:      req.ChainID,
		IsPut:        req.IsPut,
		Strike:       req.Strike,
		Expiry:       req.Expiry,
		Maker:        maker,
		Nonce:        originalRfqID,
		Price:        quoteBigInt.Mul(quoteBigInt, new(big.Int).SetInt64(1e10)).String(),
		Quantity:     req.Quantity,
		IsTakerBuy:   req.IsTakerBuy,
		ValidUntil:   time.Now().Unix() + cfg.QuoteValidDurationSeconds,
	}
	msgHash, _, err := ryskcore.CreateQuoteMessage(finalQuote)
	if err != nil {
		log.Printf("[Quote %s] Error creating quote message for signing: %v", originalRfqID, err)
		return ryskcore.Quote{}, err
	}

	signature, err := ryskcore.Sign(msgHash, privateKey)
	if err != nil {
		log.Printf("[Quote %s] Error signing quote message: %v", originalRfqID, err)
		return ryskcore.Quote{}, err
	}
	finalQuote.Signature = signature
	return finalQuote, nil
}

func getExchangeQuote(req RFQResult, asset string, exchange Exchange) (string, float64, error) {
	// get the order book using the exchange interface
	book, err := exchange.GetOrderBook(req, asset)
	if err != nil {
		return "", 0.0, err
	}
	dollarPrice, err := getPriceInclSlippage(req, book)
	if err != nil {
		log.Print(err)
		return "", 0, err
	}
	dollarPrice *= 1e6

	_, apr := CalculateAPR(
		big.NewFloat(dollarPrice/1e8),
		big.NewFloat(book.Index),
		req.Expiry,
	)
	// returns dollarPrice in 1e8
	return fmt.Sprintf("%d", int(dollarPrice)), apr, nil
}

// getOrderBook is now deprecated - functionality moved to exchange implementations
// Use exchange.GetOrderBook() instead

func getPriceInclSlippage(req RFQResult, book CCXTOrderBook) (float64, error) {
	amountBigInt, _ := new(big.Int).SetString(req.Quantity, 10)
	amount := amountBigInt.Div(amountBigInt, new(big.Int).SetUint64(1e13)).String()
	// convert the amount string to a float
	amountFloat, _ := strconv.ParseFloat(amount, 64)
	// convert the priceFloat to the correct units
	amountFloat = amountFloat / 1e5
	amount = strconv.FormatFloat(amountFloat, 'f', -1, 64)

	var cumSize float64
	var price float64
	quotes := book.Bids
	if req.IsTakerBuy {
		quotes = book.Asks
	}
	// we need to get a quote that accounts for slippage
	for _, b := range quotes {
		price = b[0]
		size := b[1]
		cumSize += size
		if cumSize >= amountFloat {
			break
		}
	}
	if cumSize < amountFloat {
		return 0.0, fmt.Errorf("cannot quote due to liquidity")
	}
	dollarPrice := math.Round(price * book.Index)
	// consider 20% premium
	premium := float64(10)
	// we take the bid to make more money
	if req.IsTakerBuy {
		dollarPrice = (math.Round(dollarPrice * (100 + premium)))
	} else {
		dollarPrice = (math.Round(dollarPrice * (100 - premium)))
	}
	return dollarPrice, nil
}

func CalculateAPR(nominator *big.Float, denominator *big.Float, maturity int64) (float64, float64) {
	expiryInTimeFormat := time.Unix(maturity, 0)
	timeToExpiryDays := expiryInTimeFormat.Sub(time.Now()).Hours() / 24
	if denominator.Cmp(big.NewFloat(0)) == 0 {
		return timeToExpiryDays, 0
	}
	rate := new(big.Float).Quo(nominator, denominator)
	periods := new(big.Float).Quo(big.NewFloat(365.25), big.NewFloat(timeToExpiryDays))
	apr := new(big.Float).Mul(new(big.Float).Mul(rate, periods), big.NewFloat(100.00))
	a, _ := apr.Float64()
	return timeToExpiryDays, a
}

// convertOptionDetailsToInstrument is now deprecated - functionality moved to exchange implementations
// Use exchange.ConvertToInstrument() instead
