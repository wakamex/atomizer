package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/wakamex/rysk-v12-cli/ryskcore"
)

// Mutex and map for debouncing
var (
	lastQuoteTime      = make(map[string]time.Time)
	lastQuoteTimeMutex = &sync.Mutex{}
	debounceDuration   = 5 * time.Second // Don't re-quote same RFQ ID within 5s
)

// isDebounced checks if we should skip quoting for this RFQ ID based on time.
func isDebounced(rfqID string) bool {
	lastQuoteTimeMutex.Lock()
	defer lastQuoteTimeMutex.Unlock()

	if t, ok := lastQuoteTime[rfqID]; ok {
		if time.Since(t) < debounceDuration {
			log.Printf("[Debounce] Skipping quote for RFQ ID %s, already quoted recently.", rfqID)
			return true
		}
	}
	lastQuoteTime[rfqID] = time.Now()
	return false
}

// sendQuoteResponse constructs and sends a quote response.
func sendQuoteResponse(client *ryskcore.Client, rfq RFQResult, originalRfqID string, cfg *AppConfig) {
	var quoteParams ryskcore.Quote
	
	// Try to get underlying asset from mapping
	underlying, hasMapping := cfg.AssetMapping[rfq.Asset]
	if hasMapping {
		// Try to use Deribit pricing
		deribitQuote, err := MakeQuote(rfq, underlying, originalRfqID, cfg)
		if err != nil {
			log.Printf("[Quote %s] Error getting Deribit quote: %v. Falling back to dummy price.", originalRfqID, err)
			// Fall back to dummy price
			quoteParams = ryskcore.Quote{
				AssetAddress: rfq.Asset,
				ChainID:      rfq.ChainID,
				Expiry:       rfq.Expiry,
				IsPut:        rfq.IsPut,
				IsTakerBuy:   rfq.IsTakerBuy,
				Maker:        cfg.MakerAddress,
				Nonce:        originalRfqID,
				Price:        cfg.DummyPrice,
				Quantity:     rfq.Quantity,
				Strike:       rfq.Strike,
				ValidUntil:   time.Now().Unix() + cfg.QuoteValidDurationSeconds,
			}
		} else {
			// Use Deribit quote
			quoteParams = deribitQuote
			// Override ValidUntil to use config value
			quoteParams.ValidUntil = time.Now().Unix() + cfg.QuoteValidDurationSeconds
		}
	} else {
		// No mapping available, use dummy price
		log.Printf("[Quote %s] No asset mapping for %s. Using dummy price.", originalRfqID, rfq.Asset)
		quoteParams = ryskcore.Quote{
			AssetAddress: rfq.Asset,
			ChainID:      rfq.ChainID,
			Expiry:       rfq.Expiry,
			IsPut:        rfq.IsPut,
			IsTakerBuy:   rfq.IsTakerBuy,
			Maker:        cfg.MakerAddress,
			Nonce:        originalRfqID,
			Price:        cfg.DummyPrice,
			Quantity:     rfq.Quantity,
			Strike:       rfq.Strike,
			ValidUntil:   time.Now().Unix() + cfg.QuoteValidDurationSeconds,
		}
	}

	// If the quote already has a signature (from MakeQuote), skip signing
	if quoteParams.Signature == "" {
		msgHash, _, err := ryskcore.CreateQuoteMessage(quoteParams)
		if err != nil {
			log.Printf("[Quote %s] Error creating quote message for signing: %v", originalRfqID, err)
			return
		}

		signature, err := ryskcore.Sign(msgHash, cfg.PrivateKey)
		if err != nil {
			log.Printf("[Quote %s] Error signing quote message: %v", originalRfqID, err)
			return
		}
		quoteParams.Signature = signature
	}

	quoteParamsBytes, err := json.Marshal(quoteParams)
	if err != nil {
		log.Printf("[Quote %s] Error marshalling quoteParams: %v", originalRfqID, err)
		return
	}

	quoteRequest := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      originalRfqID,
		Method:  "quote",
		Params:  json.RawMessage(quoteParamsBytes),
	}

	requestBytes, err := json.Marshal(quoteRequest)
	if err != nil {
		log.Printf("[Quote %s] Error marshalling quote request: %v", originalRfqID, err)
		return
	}

	client.Send(requestBytes)
	log.Printf("[Quote %s] Sent quote response: Asset=%s, Strike=%s, Expiry=%d, IsPut=%t, Price=%s, Quantity=%s, ValidUntil=%d, Nonce=%s",
		originalRfqID, quoteParams.AssetAddress, quoteParams.Strike, quoteParams.Expiry, quoteParams.IsPut, quoteParams.Price, quoteParams.Quantity, quoteParams.ValidUntil, quoteParams.Nonce)
}

// HandleRfqMessage processes an incoming RFQ message from a stream.
// It unmarshals the message, performs debouncing, and then sends a quote response.
func HandleRfqMessage(message []byte, currentAssetAddr string, mainQuoteSenderClient *ryskcore.Client, cfg *AppConfig) {
	log.Printf("RFQ In (%s): %s", currentAssetAddr, string(message))

	var rfqNotification RFQNotification
	// The RFQ data from the stream appears to be in notification.Result
	if err := json.Unmarshal(message, &rfqNotification); err == nil && rfqNotification.Result.Asset != "" {
		rfqID := rfqNotification.ID // The top-level ID seems to be the RFQ ID
		rfqData := rfqNotification.Result
		rfqData.ReceivedTS = time.Now().UnixNano()

		log.Printf("Parsed RFQ (ID: %s, Asset: %s from stream %s): ChainID=%d, Expiry=%d, IsPut=%t, Quantity=%s, Strike=%s",
			rfqID, rfqData.Asset, currentAssetAddr, rfqData.ChainID, rfqData.Expiry, rfqData.IsPut, rfqData.Quantity, rfqData.Strike)

		// Debounce before sending quote
		if isDebounced(rfqID) {
			return
		}

		// Send quote using the main client
		// Note: This is called as a goroutine in the original main's handler. 
		// Consider if HandleRfqMessage itself should be launched as a goroutine by its caller,
		// or if sendQuoteResponse should remain a goroutine call if it's long-running.
		// For now, keeping the go sendQuoteResponse pattern.
		go sendQuoteResponse(mainQuoteSenderClient, rfqData, rfqID, cfg)
		return
	} else {
		// Log if unmarshalling failed or if it's not a valid RFQNotification structure we expect
		if err != nil {
			log.Printf("RFQ In (%s) - Error unmarshalling message: %v. Raw: %s", currentAssetAddr, err, string(message))
		} else {
			log.Printf("RFQ In (%s) - Received message not matching expected RFQ structure or missing asset: %s", currentAssetAddr, string(message))
		}
	}
	// Add any other parsing logic if RFQ messages can come in different formats
}
