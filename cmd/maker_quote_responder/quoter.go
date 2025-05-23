package main

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/wakamex/rysk-v12-cli/ryskcore"
	ccxt "github.com/ccxt/ccxt/go/v4"
)

// the purpose of this module is to take in a given option request confirmation, return a quote for that option
// then seperately if  all good execute the trade which triggers the solidity side and then opens an offsetting position
// on deribit
// / the underlying must be a string that is either "ETH", "BTC" or "SOL" (essentially the primary asset, not the LST)
// / the expiry must be provided in the format "DDMMMYY" (e.g. "31JAN25") (note that on deribit if the first number is 0, it will cut it e.g. 6JAN25)
// / quantity in normal numbers
func MakeQuote(req RFQResult, underlying string, originalRfqID string) (ryskcore.Quote, error) {
	quote, _, err := getDeribitQuote(req, underlying)
	if err != nil {
		return ryskcore.Quote{}, err
	}
	quoteBigInt, ok := new(big.Int).SetString(quote, 10)
	if !ok {
		return ryskcore.Quote{}, fmt.Errorf("bad quote conversion: %s", quote)
	}
	maker := os.Getenv("MAKER_ADDRESS")
	privateKey := os.Getenv("PRIVATE_KEY")
	finalQuote := ryskcore.Quote{
		AssetAddress: req.Asset,
		ChainID:      req.ChainID,
		IsPut:        req.IsPut,
		Strike:       req.Strike,
		Expiry:       req.Expiry,
		Maker:        maker,
		Nonce:        strconv.FormatInt(time.Now().UnixMicro(), 10),
		Price:        quoteBigInt.Mul(quoteBigInt, new(big.Int).SetInt64(1e10)).String(),
		Quantity:     req.Quantity,
		IsTakerBuy:   req.IsTakerBuy,
		ValidUntil:   time.Now().Unix() + 60,
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

func getDeribitQuote(req RFQResult, asset string) (string, float64, error) {
	// get the order book using CCXT
	book, err := getOrderBook(req, asset)
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

func getOrderBook(req RFQResult, asset string) (CCXTOrderBook, error) {
	// Convert option details to instrument name
	instrumentName, err := convertOptionDetailsToInstrument(asset, req.Strike, req.Expiry, req.IsPut)
	if err != nil {
		return CCXTOrderBook{}, err
	}

	// Initialize CCXT exchange
	exchange := ccxt.NewDeribit(map[string]interface{}{
		"rateLimit": 10,
		"enableRateLimit": true,
	})

	// CCXT requires the symbol format with prefix for options
	ccxtSymbol := fmt.Sprintf("%s/USD:%s", asset, instrumentName)
	
	// Fetch order book
	orderBook, err := exchange.FetchOrderBook(ccxtSymbol)
	if err != nil {
		return CCXTOrderBook{}, fmt.Errorf("failed to fetch order book: %v", err)
	}

	// Fetch ticker for the option instrument to get its price
	optionTicker, err := exchange.FetchTicker(ccxtSymbol)
	if err != nil {
		return CCXTOrderBook{}, fmt.Errorf("failed to fetch option ticker: %v", err)
	}

	// Get the underlying index price from the option ticker
	indexPrice := 0.0
	if optionTicker.Info != nil {
		// Deribit option tickers include underlying price in the info
		if underlyingPrice, exists := optionTicker.Info["underlying_price"]; exists {
			if price, ok := underlyingPrice.(float64); ok {
				indexPrice = price
			}
		}
	}
	
	// If we couldn't get index price from option ticker, try to fetch spot/futures
	if indexPrice == 0.0 {
		// Try different instrument names that Deribit might support
		possibleIndexInstruments := []string{
			asset + "-PERPETUAL",      // This is the correct format for Deribit
			asset + "_USDC-PERPETUAL",
			asset + "_USD-PERPETUAL",
		}
		
		for _, instrument := range possibleIndexInstruments {
			indexTicker, err := exchange.FetchTicker(instrument)
			if err == nil && indexTicker.Last != nil {
				indexPrice = *indexTicker.Last
				break
			}
		}
	}

	// Convert order book to our structure
	book := CCXTOrderBook{
		Bids:  orderBook.Bids,
		Asks:  orderBook.Asks,
		Index: indexPrice,
	}

	return book, nil
}

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

func convertOptionDetailsToInstrument(
	asset string,
	strike string,
	expiry int64,
	isPut bool,
) (string, error) {
	// convert the strike from a big.Int string to a normal number
	strikeBigInt, ok := new(big.Int).SetString(strike, 10)
	if !ok {
		return "", fmt.Errorf("invalid strike")
	}
	strike = strikeBigInt.Div(strikeBigInt, new(big.Int).SetUint64(1e8)).String()

	// convert the expiry from a timestamp seconds into a deribit compatible date time
	// Deribit uses YYMMDD format
	deribitExpiry := time.Unix(expiry, 0).Format("060102")
	// convert isPut to "C" or "P"
	optionType := "C"
	if isPut {
		optionType = "P"
		return "", fmt.Errorf("puts not supported")
	}

	// construct the instrument name in Deribit format: ASSET-YYMMDD-STRIKE-C/P
	// e.g., ETH-250530-3000-C for an ETH call option expiring May 30, 2025
	instrumentName := fmt.Sprintf("%s-%s-%s-%s", asset, deribitExpiry, strike, optionType)
	return instrumentName, nil
}
