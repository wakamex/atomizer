package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/wakamex/rysk-v12-cli/ryskcore"
)

// Debouncing prevents duplicate quotes for the same RFQ within a time window
var (
	lastQuoteTime      = make(map[string]time.Time)
	lastQuoteTimeMutex = &sync.Mutex{}
	debounceDuration   = 5 * time.Second
)

func isDebounced(rfqID string) bool {
	lastQuoteTimeMutex.Lock()
	defer lastQuoteTimeMutex.Unlock()
	if t, ok := lastQuoteTime[rfqID]; ok && time.Since(t) < debounceDuration {
		log.Printf("[Debounce] Skipping quote for RFQ ID %s, already quoted recently.", rfqID)
		return true
	}
	lastQuoteTime[rfqID] = time.Now()
	return false
}

// Generates and sends a quote response with real-time exchange pricing or dummy fallback
func sendQuoteResponse(client *ryskcore.Client, rfq RFQResult, originalRfqID string, cfg *AppConfig, exchange Exchange) {
	// Create base quote params
	quoteParams := ryskcore.Quote{
		AssetAddress: rfq.Asset, ChainID: rfq.ChainID, Expiry: rfq.Expiry, IsPut: rfq.IsPut,
		IsTakerBuy: rfq.IsTakerBuy, Maker: cfg.MakerAddress, Nonce: originalRfqID,
		Price: cfg.DummyPrice, Quantity: rfq.Quantity, Strike: rfq.Strike,
		ValidUntil: time.Now().Unix() + cfg.QuoteValidDurationSeconds,
	}

	// Try real-time pricing from exchange if asset mapping exists
	if underlying, hasMapping := cfg.AssetMapping[rfq.Asset]; hasMapping {
		if exchangeQuote, err := MakeQuote(rfq, underlying, originalRfqID, cfg, exchange); err != nil {
			log.Printf("[Quote %s] Error getting %s quote: %v. Falling back to dummy price.", originalRfqID, cfg.ExchangeName, err)
		} else {
			quoteParams = exchangeQuote
			quoteParams.ValidUntil = time.Now().Unix() + cfg.QuoteValidDurationSeconds
		}
	} else {
		log.Printf("[Quote %s] No asset mapping for %s. Using dummy price.", originalRfqID, rfq.Asset)
	}

	// Sign quote with EIP-712 if not already signed by MakeQuote
	if quoteParams.Signature == "" {
		if msgHash, _, err := ryskcore.CreateQuoteMessage(quoteParams); err != nil {
			log.Printf("[Quote %s] Error creating quote message for signing: %v", originalRfqID, err)
			return
		} else if signature, err := ryskcore.Sign(msgHash, cfg.PrivateKey); err != nil {
			log.Printf("[Quote %s] Error signing quote message: %v", originalRfqID, err)
			return
		} else {
			quoteParams.Signature = signature
		}
	}

	// Send JSON-RPC quote response to the maker endpoint
	if quoteParamsBytes, err := json.Marshal(quoteParams); err != nil {
		log.Printf("[Quote %s] Error marshalling quoteParams: %v", originalRfqID, err)
	} else if requestBytes, err := json.Marshal(JSONRPCRequest{
		JSONRPC: "2.0", ID: originalRfqID, Method: "quote", Params: json.RawMessage(quoteParamsBytes),
	}); err != nil {
		log.Printf("[Quote %s] Error marshalling quote request: %v", originalRfqID, err)
	} else {
		client.Send(requestBytes)
		log.Printf("[Quote %s] Sent quote response: Asset=%s, Strike=%s, Expiry=%d, IsPut=%t, Price=%s, Quantity=%s, ValidUntil=%d, Nonce=%s",
			originalRfqID, quoteParams.AssetAddress, quoteParams.Strike, quoteParams.Expiry, quoteParams.IsPut, quoteParams.Price, quoteParams.Quantity, quoteParams.ValidUntil, quoteParams.Nonce)
	}
}


// Processes incoming RFQ (Request for Quote) messages and responds with quotes
func HandleRfqMessage(message []byte, currentAssetAddr string, mainQuoteSenderClient *ryskcore.Client, cfg *AppConfig, exchange Exchange) {
	log.Printf("RFQ In (%s): %s", currentAssetAddr, string(message))

	var rfqNotification RFQNotification
	if err := json.Unmarshal(message, &rfqNotification); err != nil {
		log.Printf("RFQ In (%s) - Error unmarshalling message: %v. Raw: %s", currentAssetAddr, err, string(message))
	} else if rfqNotification.Result.Asset == "" {
		log.Printf("RFQ In (%s) - Missing asset in RFQ: %s", currentAssetAddr, string(message))
	} else if isDebounced(rfqNotification.ID) {
		return // Skip duplicate RFQs within debounce window
	} else {
		// Process valid RFQ and respond with quote
		rfqData := rfqNotification.Result
		rfqData.ReceivedTS = time.Now().UnixNano()
		log.Printf("Parsed RFQ (ID: %s, Asset: %s from stream %s): ChainID=%d, Expiry=%d, IsPut=%t, Quantity=%s, Strike=%s",
			rfqNotification.ID, rfqData.Asset, currentAssetAddr, rfqData.ChainID, rfqData.Expiry, rfqData.IsPut, rfqData.Quantity, rfqData.Strike)
		go sendQuoteResponse(mainQuoteSenderClient, rfqData, rfqNotification.ID, cfg, exchange)
	}
}
