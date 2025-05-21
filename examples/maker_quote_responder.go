package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/wakamex/ryskV12-cli/ryskcore"
)

// JSONRPCError defines the structure for a JSON-RPC error object.
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// JSONRPCResponse defines the structure for a generic JSON-RPC response.
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      string          `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
	// Fields to catch notifications that might also conform to some parts of this structure
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
}

// JSONRPCRequest defines the structure for a JSON-RPC request.
// Note: Duplicated in other examples. Consider moving to a shared package.
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"` // Field name changed to all caps
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// RFQNotification represents the structure of an incoming RFQ message from the server.
// This is based on the example log: {"jsonrpc":"2.0","id":"rfq_id_here","result":{...RFQDetails...}}
type RFQNotification struct {
	JsonRPC string    `json:"jsonrpc"`
	ID      string    `json:"id"` // This is the RFQ ID
	Result  RFQResult `json:"result,omitempty"`
	Error   *RPCError `json:"error,omitempty"`  // For subscription confirmations or errors
	Method  string    `json:"method,omitempty"` // For server-initiated messages like 'rfq'
	Params  RFQResult `json:"params,omitempty"` // If method is 'rfq', params might hold the data
}

// RFQResult contains the actual details of the RFQ.
// Note: The structure might be nested if Method is 'rfq'
// For example: {"method":"rfq","params":{"id":"rfq_id_actual", "asset":"0x...", ...}}
type RFQResult struct {
	ID         string `json:"id,omitempty"` // This might be the actual RFQ ID if nested under params
	Asset      string `json:"asset"`
	AssetName  string `json:"assetName,omitempty"`
	ChainID    int    `json:"chainId"`
	Expiry     int64  `json:"expiry"`
	IsPut      bool   `json:"isPut"`
	IsTakerBuy bool   `json:"isTakerBuy"` // Added to match RFQ
	Quantity   string `json:"quantity"`
	Strike     string `json:"strike"`
	Taker      string `json:"taker,omitempty"`
	ReceivedTS int64  `json:"-"` // Local timestamp when RFQ was processed
}

// RPCError structure for JSON-RPC errors
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var ( // Global config variables set by flags
	websocketURL              string
	rfqAssetAddressesCSV      string
	makerAddress              string
	privateKey                string
	dummyPrice                string
	quoteValidDurationSeconds int64
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	flag.StringVar(&websocketURL, "websocket_url", "wss://rip-testnet.rysk.finance/maker", "WebSocket URL for RFQ stream and quote submission")
	flag.StringVar(&rfqAssetAddressesCSV, "rfq_asset_addresses", "", "Comma-separated list of asset addresses for RFQ streams (e.g., 0xAsset1,0xAsset2)")
	flag.StringVar(&dummyPrice, "dummy_price", "1000000", "Dummy price to quote (ensure format matches Rysk requirements, e.g., units)")
	flag.Int64Var(&quoteValidDurationSeconds, "quote_valid_duration_seconds", 30, "How long your quotes will be valid in seconds")
	flag.Parse()
	makerAddress = os.Getenv("MAKER_ADDRESS")
	privateKey = os.Getenv("PRIVATE_KEY")

	if rfqAssetAddressesCSV == "" {
		log.Fatal("Error: --rfq_asset_addresses is required.")
	}
	if makerAddress == "" {
		log.Fatal("Error: MAKER_ADDRESS environment variable is not set or empty.")
	}
	if privateKey == "" {
		log.Fatal("Error: PRIVATE_KEY environment variable is not set or empty.")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sdkClient, err := ryskcore.NewClient(ctx, websocketURL, nil)
	if err != nil {
		log.Fatalf("Error creating SDK client: %v", err)
	}
	defer sdkClient.Close()

	log.Printf("Connected to WebSocket: %s", websocketURL)

	// Setup message handler for the main sdkClient (for quote confirmations/errors)
	sdkClient.SetHandler(func(message []byte) {
		const clientIdentifier = "MainSdkClient" // Identifier for logging
		// log.Printf("[%s] Raw In: %s", clientIdentifier, string(message)) // Uncomment for deep debugging

		var baseResponse JSONRPCResponse
		if err := json.Unmarshal(message, &baseResponse); err != nil {
			log.Printf("[%s] Error unmarshalling base JSON-RPC response: %v. Message: %s", clientIdentifier, err, string(message))
			return
		}

		if baseResponse.Error != nil {
			log.Printf("[%s] Received RPC Error: ID=%s, Code=%d, Message=%s", clientIdentifier, baseResponse.ID, baseResponse.Error.Code, baseResponse.Error.Message)
			return
		}

		if baseResponse.Result != nil {
			var resultStr string
			if err := json.Unmarshal(baseResponse.Result, &resultStr); err == nil {
				if strings.ToUpper(resultStr) == "OK" {
					log.Printf("[%s] Received 'OK' confirmation for ID %s.", clientIdentifier, baseResponse.ID)
					return
				}
			}
			// If it's not a simple "OK" or an error, just log the result for the main client.
			log.Printf("[%s] Received Result for ID %s: %s", clientIdentifier, baseResponse.ID, string(baseResponse.Result))
		} else if baseResponse.Method != "" {
			log.Printf("[%s] Received unhandled method '%s' with ID %s. Params: %s", clientIdentifier, baseResponse.Method, baseResponse.ID, string(baseResponse.Params))
		} else {
			log.Printf("[%s] Received message with no error, result, or method. ID: %s. Raw: %s", clientIdentifier, baseResponse.ID, string(message))
		}
	})

	// Connect to individual RFQ streams for each asset
	assetAddressesList := strings.Split(rfqAssetAddressesCSV, ",")
	if len(assetAddressesList) == 0 || rfqAssetAddressesCSV == "" {
		log.Println("No asset addresses provided for RFQ subscription. Exiting.")
		// It might make sense to exit if no assets to listen to, or let it run just as a sender if that's a use case.
		// For now, let's assume it needs assets to be useful.
		return
	}

	// Determine base URL for RFQ streams (e.g., wss://host.com from wss://host.com/maker)
	baseURL := strings.TrimSuffix(websocketURL, "/maker") // Common case
	// A more robust way to get base URL might be needed if /maker isn't always the suffix
	if !strings.HasPrefix(baseURL, "ws://") && !strings.HasPrefix(baseURL, "wss://") {
		log.Printf("Warning: Could not reliably determine base URL from %s to construct RFQ stream URLs. Attempting to use as is.", websocketURL)
		baseURL = websocketURL // Fallback, might not be correct
	}

	for _, addr := range assetAddressesList {
		trimmedAddr := strings.TrimSpace(addr)
		if trimmedAddr == "" {
			continue
		}

		rfqStreamURL := fmt.Sprintf("%s/rfqs/%s", baseURL, trimmedAddr)
		log.Printf("Attempting to connect to RFQ Stream for %s at %s", trimmedAddr, rfqStreamURL)

		rfqClientCtx, rfqClientCancel := context.WithCancel(ctx) // Give each RFQ client its own sub-context for independent shutdown if needed
		// defer rfqClientCancel() // This would cancel them all when main returns. Consider if needed or if main client closure handles it.

		rfqSdkClient, err := ryskcore.NewClient(rfqClientCtx, rfqStreamURL, nil)
		if err != nil {
			log.Printf("Failed to connect to RFQ stream %s: %v", rfqStreamURL, err)
			rfqClientCancel() // Cancel context for this failed client
			continue
		}
		log.Printf("Successfully connected to RFQ stream for asset %s at %s", trimmedAddr, rfqStreamURL)
		defer rfqSdkClient.Close() // Ensure client is closed when its goroutine/context finishes

		// Capture sdkClient (main quote sender) and currentAddr for the handler
		mainQuoteSenderClient := sdkClient
		currentAssetAddr := trimmedAddr

		rfqSdkClient.SetHandler(func(message []byte) {
			log.Printf("RFQ In (%s): %s", currentAssetAddr, string(message))
			// Messages on this stream are expected to be RFQs for this asset
			// Attempt to unmarshal as an RFQ notification (assuming structure is consistent)
			var rfqNotification RFQNotification // Assuming RFQNotification is the correct struct for direct RFQ stream messages
			// It's possible the direct stream sends just the RFQData part, or a slightly different wrapper.
			// For now, let's assume it's the same RFQNotification structure and the Method field might be absent or different.

			// Option 1: Try to unmarshal as RFQNotification and check key fields
			// The RFQ data from the stream appears to be in notification.Result
			if err := json.Unmarshal(message, &rfqNotification); err == nil && rfqNotification.Result.Asset != "" {
				rfqID := rfqNotification.ID // The top-level ID seems to be the RFQ ID
				rfqData := rfqNotification.Result
				rfqData.ReceivedTS = time.Now().UnixNano()
				// Note: rfqData (which is RFQResult type) does not have an 'ID' field itself if it's from the 'result' block.
				// The original RFQ ID is captured from rfqNotification.ID.
				log.Printf("Parsed RFQ (ID: %s, Asset: %s from stream %s): ChainID=%d, Expiry=%d, IsPut=%t, Quantity=%s, Strike=%s",
					rfqID, rfqData.Asset, currentAssetAddr, rfqData.ChainID, rfqData.Expiry, rfqData.IsPut, rfqData.Quantity, rfqData.Strike)
				go sendQuoteResponse(mainQuoteSenderClient, rfqData, rfqID) // Send quote using the main client
				return
			} else {
				// Log if unmarshalling failed or if it's not a valid RFQNotification structure we expect
				if err != nil {
					log.Printf("RFQ In (%s) - Error unmarshalling message: %v. Raw: %s", currentAssetAddr, err, string(message))
				} else {
					log.Printf("RFQ In (%s) - Received message not matching expected RFQ structure or missing asset: %s", currentAssetAddr, string(message))
				}
			}

			// Option 2: If direct stream sends RFQData directly (or a simpler wrapper)
			// This section can be enabled if Option 1 consistently fails and you suspect a different format.
			// var rfqData RFQData // Assuming RFQData struct is defined
			// if err := json.Unmarshal(message, &rfqData); err == nil && rfqData.Asset != "" {
			// 	 rfqID := "some_generated_id_or_from_message_if_present" // RFQ ID might not be in the same place
			// 	 rfqData.ReceivedTS = time.Now().UnixNano()
			// 	 log.Printf("Parsed Direct RFQ (Asset: %s from stream %s): ...details...", rfqData.Asset, currentAssetAddr)
			// 	 go sendQuoteResponse(mainQuoteSenderClient, rfqData, rfqID)
			// 	 return
			// }

			// If neither Option 1 nor an enabled Option 2 handles it, it's unexpected.
			// The logging inside the 'else' of Option 1 already covers this for now.
		})

		// Keep track of cancel functions for cleanup if needed, or rely on main context cancellation
		_ = rfqClientCancel // To avoid unused variable error if not used later for explicit individual shutdown
	}

	log.Println("RFQ listener setup complete. Waiting for RFQs and termination signal.")
	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("Received signal: %v. Shutting down...", sig)
	case <-ctx.Done():
		log.Println("SDK client context done. Shutting down...")
	}

	cancel() // Ensure all goroutines using ctx are signaled
	log.Println("Quote Responder Daemon finished.")
}

func sendQuoteResponse(client *ryskcore.Client, rfq RFQResult, originalRfqID string) {
	// Debounce: Simple check to avoid re-quoting too quickly for the same RFQ ID if messages are duplicated by chance
	// A more robust solution would use a timed cache.
	// For now, this is a placeholder if needed; actual RFQ IDs should be unique.

	quoteParams := ryskcore.Quote{
		AssetAddress: rfq.Asset,
		ChainID:      rfq.ChainID,
		Expiry:       rfq.Expiry,
		IsPut:        rfq.IsPut,
		IsTakerBuy:   rfq.IsTakerBuy, // Use the value from the parsed RFQ
		Maker:        makerAddress,
		Nonce:        fmt.Sprintf("%s-%d", originalRfqID, time.Now().UnixNano()), // Unique nonce for our quote
		Price:        dummyPrice,
		Quantity:     rfq.Quantity,
		Strike:       rfq.Strike,
		ValidUntil:   time.Now().Unix() + quoteValidDurationSeconds,
	}

	msgHash, _, err := ryskcore.CreateQuoteMessage(quoteParams)
	if err != nil {
		log.Printf("[Quote %s] Error creating quote message for signing: %v", originalRfqID, err)
		return
	}

	signature, err := ryskcore.Sign(msgHash, privateKey)
	if err != nil {
		log.Printf("[Quote %s] Error signing quote message: %v", originalRfqID, err)
		return
	}
	quoteParams.Signature = signature

	// Marshal quoteParams to json.RawMessage for the Params field
	quoteParamsBytes, err := json.Marshal(quoteParams)
	if err != nil {
		log.Printf("[Quote %s] Error marshalling quoteParams: %v", originalRfqID, err)
		return
	}

	quoteRequest := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      originalRfqID, // Use the original RFQ ID as the ID for our quote response
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

	// Optional: Add a small delay or use a rate limiter if sending too many quotes too fast
	// time.Sleep(100 * time.Millisecond)
}

// Mutex and map for debouncing (optional, simple example)
var (
	lastQuoteTime      = make(map[string]time.Time)
	lastQuoteTimeMutex = &sync.Mutex{}
	debounceDuration   = 5 * time.Second // Don't re-quote same RFQ ID within 5s
)

// isDebounced checks if we should skip quoting for this RFQ ID based on time.
// This is a simple example and might not be perfectly robust for all scenarios.
// The Rysk system should ideally handle duplicate RFQ IDs or provide unique ones.
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
