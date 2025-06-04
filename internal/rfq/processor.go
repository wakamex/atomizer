package rfq

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/wakamex/atomizer/internal/config"
	"github.com/wakamex/atomizer/internal/quoter"
	"github.com/wakamex/atomizer/internal/types"
	"github.com/wakamex/rysk-v12-cli/ryskcore"
)

// Processor handles RFQ processing and quote generation
type Processor struct {
	config               *config.Config
	exchange             types.Exchange
	lastQuoteTime        map[string]time.Time
	lastQuoteTimeMutex   sync.Mutex
	debounceDuration     time.Duration
}

// NewProcessor creates a new RFQ processor
func NewProcessor(cfg *config.Config, exchange types.Exchange) *Processor {
	return &Processor{
		config:           cfg,
		exchange:         exchange,
		lastQuoteTime:    make(map[string]time.Time),
		debounceDuration: 5 * time.Second,
	}
}

// ProcessRFQ handles an incoming RFQ and generates a quote response
func (p *Processor) ProcessRFQ(client RyskClient, rfq types.RFQResult, originalRfqID string) error {
	// Check debounce
	if p.isDebounced(originalRfqID) {
		return nil
	}

	// Generate quote
	quote, err := p.generateQuote(rfq, originalRfqID)
	if err != nil {
		return fmt.Errorf("failed to generate quote: %w", err)
	}

	// Send quote response
	return p.sendQuoteResponse(client, quote, originalRfqID)
}

// isDebounced checks if we've recently quoted this RFQ
func (p *Processor) isDebounced(rfqID string) bool {
	p.lastQuoteTimeMutex.Lock()
	defer p.lastQuoteTimeMutex.Unlock()
	
	if t, ok := p.lastQuoteTime[rfqID]; ok && time.Since(t) < p.debounceDuration {
		log.Printf("[Debounce] Skipping quote for RFQ ID %s, already quoted recently.", rfqID)
		return true
	}
	
	p.lastQuoteTime[rfqID] = time.Now()
	return false
}

// generateQuote creates a quote for the RFQ
func (p *Processor) generateQuote(rfq types.RFQResult, originalRfqID string) (*ryskcore.Quote, error) {
	// Create base quote params
	quoteParams := &ryskcore.Quote{
		AssetAddress: rfq.Asset,
		ChainID:      rfq.ChainID,
		Expiry:       rfq.Expiry,
		IsPut:        rfq.IsPut,
		IsTakerBuy:   rfq.IsTakerBuy,
		Maker:        p.config.MakerAddress,
		Nonce:        originalRfqID,
		Price:        p.config.DummyPrice,
		Quantity:     rfq.Quantity,
		Strike:       rfq.Strike,
		ValidUntil:   time.Now().Unix() + p.config.QuoteValidDurationSeconds,
	}

	// Try real-time pricing from exchange if asset mapping exists
	if underlying, hasMapping := p.config.AssetMapping[rfq.Asset]; hasMapping {
		exchangeQuote, err := p.makeExchangeQuote(rfq, underlying, originalRfqID)
		if err != nil {
			log.Printf("[Quote %s] Error getting %s quote: %v. Falling back to dummy price.", 
				originalRfqID, p.config.ExchangeName, err)
			// Sign the dummy quote since we're falling back
			if err := p.signQuote(quoteParams); err != nil {
				return nil, fmt.Errorf("failed to sign quote: %w", err)
			}
		} else {
			// Use exchange quote (already signed by quoter.MakeQuote)
			quoteParams = exchangeQuote
		}
	} else {
		log.Printf("[Quote %s] No asset mapping for %s. Using dummy price.", originalRfqID, rfq.Asset)
		// Sign the dummy quote
		if err := p.signQuote(quoteParams); err != nil {
			return nil, fmt.Errorf("failed to sign quote: %w", err)
		}
	}

	return quoteParams, nil
}

// makeExchangeQuote generates a quote using exchange pricing
func (p *Processor) makeExchangeQuote(rfq types.RFQResult, underlying string, rfqID string) (*ryskcore.Quote, error) {
	// Use the quoter module to generate a properly signed quote
	quote, err := quoter.MakeQuote(rfq, underlying, rfqID, p.config, p.exchange)
	if err != nil {
		return nil, fmt.Errorf("failed to make quote: %w", err)
	}
	
	// Convert to pointer for compatibility
	return &quote, nil
}

// signQuote signs the quote using EIP-712
func (p *Processor) signQuote(quote *ryskcore.Quote) error {
	// Create the message hash
	messageHash, _, err := ryskcore.CreateQuoteMessage(*quote)
	if err != nil {
		return fmt.Errorf("failed to create quote message: %w", err)
	}
	
	// Convert private key to hex string for signing
	privateKeyBytes := crypto.FromECDSA(p.config.ParsedPrivateKey)
	privateKeyHex := fmt.Sprintf("%x", privateKeyBytes)
	
	// Sign the message
	signature, err := ryskcore.Sign(messageHash, privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to sign quote: %w", err)
	}
	
	// Update the quote with the signature
	quote.Signature = signature
	return nil
}

// sendQuoteResponse sends the quote back to the RFQ requester
func (p *Processor) sendQuoteResponse(client RyskClient, quote *ryskcore.Quote, rfqID string) error {
	// Marshal quote params
	quoteParamsBytes, err := json.Marshal(quote)
	if err != nil {
		return fmt.Errorf("failed to marshal quote params: %w", err)
	}

	// Create JSON-RPC request
	request := types.JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      rfqID,
		Method:  "quote",
		Params:  json.RawMessage(quoteParamsBytes),
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Send response
	client.Send(requestBytes)
	
	log.Printf("[Quote %s] Sent quote response: Asset=%s, Strike=%s, Expiry=%d, IsPut=%t, Price=%s, Quantity=%s, ValidUntil=%d",
		rfqID, quote.AssetAddress, quote.Strike, quote.Expiry, quote.IsPut, 
		quote.Price, quote.Quantity, quote.ValidUntil)

	return nil
}

// RyskClient interface for sending responses
type RyskClient interface {
	Send(data []byte)
}