package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wakamex/rysk-v12-cli/ryskcore" // Re-add ryskcore for type definitions if needed, or ensure types are fully self-contained elsewhere
)

const (
	sessionInitialBackoff = 1 * time.Second
	sessionMaxBackoff     = 30 * time.Second
)

// getBuildHash returns a hash of the binary to identify different builds
func getBuildHash() string {
	// Try to get build info
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	
	// Create a hash from build info
	h := sha256.New()
	h.Write([]byte(info.Main.Version))
	h.Write([]byte(info.Main.Sum))
	h.Write([]byte(info.GoVersion))
	
	// Add VCS info if available
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" || setting.Key == "vcs.time" {
			h.Write([]byte(setting.Value))
		}
	}
	
	// If we can read the binary itself, add its hash
	if exePath, err := os.Executable(); err == nil {
		if data, err := os.ReadFile(exePath); err == nil {
			// Just hash first 1MB to avoid memory issues with large binaries
			maxLen := 1024 * 1024
			if len(data) > maxLen {
				data = data[:maxLen]
			}
			h.Write(data)
		}
	}
	
	return hex.EncodeToString(h.Sum(nil))[:16] // Return first 16 chars for brevity
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	
	// Print build hash as first output
	buildHash := getBuildHash()
	log.Printf("========================================")
	log.Printf("Build hash: %s", buildHash)
	log.Printf("========================================")
	
	cfg := LoadConfig()
	
	// Create exchange instance using factory
	factory := NewExchangeFactory()
	exchange, err := factory.CreateExchange(cfg)
	if err != nil {
		log.Fatalf("Failed to create exchange: %v", err)
	}
	if cfg.ExchangeTestMode {
		log.Printf("Successfully initialized %s exchange in TEST MODE", cfg.ExchangeName)
	} else {
		log.Printf("Successfully initialized %s exchange in PRODUCTION MODE", cfg.ExchangeName)
	}

	// Create orchestrator
	orchestrator := NewArbitrageOrchestrator(cfg, exchange)
	
	// Start the orchestrator
	go orchestrator.Start()
	
	// Create and start HTTP server for manual trades
	httpPort := 8080 // Default port
	if cfg.HTTPPort != "" {
		// Parse port from string
		if port, err := strconv.Atoi(cfg.HTTPPort); err == nil {
			httpPort = port
		}
	}
	httpServer := NewHTTPServer(orchestrator, orchestrator.riskManager, httpPort)
	go func() {
		log.Printf("Starting HTTP server on port %d", httpPort)
		if err := httpServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()
	defer httpServer.Stop()

	// appCtx governs the entire application lifecycle, cancelled by OS signals.
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	
	// Stop orchestrator when app stops
	defer orchestrator.Stop()

	assetAddressesList := strings.Split(cfg.RFQAssetAddressesCSV, ",")
	if len(assetAddressesList) == 0 || cfg.RFQAssetAddressesCSV == "" {
		log.Println("No asset addresses provided for RFQ subscription. Application will not start RFQ listeners.")
		// Decide if the app should exit or run without RFQ streams if that's a valid use case.
		// For now, we allow it to proceed, potentially only having the main client (if it had other purposes).
	}

	baseURL := ""
	if cfg.RFQAssetAddressesCSV != "" { // Only determine baseURL if we have assets to connect to
		baseURL = strings.TrimSuffix(cfg.WebSocketURL, "/maker")
		if !strings.HasPrefix(baseURL, "ws://") && !strings.HasPrefix(baseURL, "wss://") {
			log.Printf("Warning: Could not reliably determine base URL from %s to construct RFQ stream URLs. Attempting to use as is.", cfg.WebSocketURL)
			baseURL = cfg.WebSocketURL
		}
	}

	currentSessionBackoff := sessionInitialBackoff

appRunningLoop:
	for {
		select {
		case <-appCtx.Done():
			log.Printf("Application context cancelled (%v). Exiting main loop.", appCtx.Err())
			break appRunningLoop
		default:
		}

		log.Println("Starting new session...")
		// sessionCtx governs the lifecycle of the main SDK client and its associated RFQ streams for this attempt.
		sessionCtx, sessionCancel := context.WithCancel(appCtx)

		mainSdkClient, err := EstablishConnectionWithRetry(sessionCtx, cfg.WebSocketURL, "MainConnection")
		if err != nil {
			log.Printf("Failed to establish main SDK connection for session: %v. Will retry after backoff.", err)
			sessionCancel() // Ensure session resources are cleaned up if main client fails
			// Backoff before retrying the session
			if !waitForBackoffOrSignal(appCtx, currentSessionBackoff) {
				break appRunningLoop // App context cancelled during backoff
			}
			currentSessionBackoff = nextBackoff(currentSessionBackoff, sessionInitialBackoff, sessionMaxBackoff)
			continue appRunningLoop
		}
		log.Println("Main SDK client connected for session.")
		// Reset session backoff on successful connection
		currentSessionBackoff = sessionInitialBackoff

		// Setup message handler for the main mainSdkClient
		mainSdkClient.SetHandler(func(message []byte) {
			const clientIdentifier = "MainSdkClient"
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
				if err := json.Unmarshal(baseResponse.Result, &resultStr); err == nil && strings.ToUpper(resultStr) == "OK" {
					log.Printf("[%s] Received 'OK' confirmation for ID %s. Full message: %s", clientIdentifier, baseResponse.ID, string(message))
					return
				}
				
				// Check if this is a trade confirmation
				if baseResponse.ID == "trade" {
					var confirmation RFQConfirmation
					if err := json.Unmarshal(baseResponse.Result, &confirmation); err != nil {
						log.Printf("[%s] Error unmarshalling trade confirmation: %v", clientIdentifier, err)
					} else {
						log.Printf("[%s] Received trade confirmation for Quote ID %s", clientIdentifier, confirmation.QuoteNonce)
						
						// Get the underlying asset from the asset mapping
						if _, hasMapping := cfg.AssetMapping[confirmation.AssetAddress]; hasMapping {
							// Convert confirmation to TradeEvent and send to orchestrator
							quantity, _ := decimal.NewFromString(confirmation.Quantity)
							price, _ := decimal.NewFromString(confirmation.Price)
							tradeEvent := &TradeEvent{
								ID:             confirmation.QuoteNonce,
								Source:         TradeSourceRysk,
								Status:         TradeStatusExecuted,
								RFQId:          confirmation.QuoteNonce,
								Instrument:     confirmation.AssetAddress,
								Strike:         decimal.RequireFromString(confirmation.Strike),
								Expiry:         int64(confirmation.Expiry),
								IsPut:          confirmation.IsPut,
								Quantity:       quantity,
								Price:          price,
								IsTakerBuy:     confirmation.IsTakerBuy,
								Timestamp:      time.Now(),
							}
							
							// Send to orchestrator
							if err := orchestrator.HandleManualTrade(tradeEvent); err != nil {
								log.Printf("[%s] Error processing trade through orchestrator: %v", clientIdentifier, err)
							} else {
								log.Printf("[%s] Successfully sent trade to orchestrator for Quote ID %s", clientIdentifier, confirmation.QuoteNonce)
							}
						} else {
							log.Printf("[%s] No asset mapping for %s. Cannot hedge.", clientIdentifier, confirmation.AssetAddress)
						}
						return
					}
				}
				
				log.Printf("[%s] Received Result for ID %s: %s", clientIdentifier, baseResponse.ID, string(baseResponse.Result))
			} else if baseResponse.Method != "" {
				// Handle different methods
				switch baseResponse.Method {
				case "rfq_confirmation":
					// Handle RFQ confirmation - this means a trade was executed
					var confirmation RFQConfirmation
					if err := json.Unmarshal(baseResponse.Params, &confirmation); err != nil {
						log.Printf("[%s] Error unmarshalling rfq_confirmation params: %v", clientIdentifier, err)
						return
					}
					log.Printf("[%s] Received RFQ confirmation for Quote ID %s", clientIdentifier, confirmation.QuoteNonce)
					
					// Get the underlying asset from the asset mapping
					if _, hasMapping := cfg.AssetMapping[confirmation.AssetAddress]; hasMapping {
						// Convert confirmation to TradeEvent and send to orchestrator
						quantity, _ := decimal.NewFromString(confirmation.Quantity)
						price, _ := decimal.NewFromString(confirmation.Price)
						tradeEvent := &TradeEvent{
							ID:             confirmation.QuoteNonce,
							Source:         TradeSourceRysk,
							Status:         TradeStatusExecuted,
							RFQId:          confirmation.QuoteNonce,
							Instrument:     confirmation.AssetAddress,
							Strike:         decimal.RequireFromString(confirmation.Strike),
							Expiry:         int64(confirmation.Expiry),
							IsPut:          confirmation.IsPut,
							Quantity:       quantity,
							Price:          price,
							IsTakerBuy:     confirmation.IsTakerBuy,
							Timestamp:      time.Now(),
						}
						
						// Send to orchestrator
						if err := orchestrator.HandleManualTrade(tradeEvent); err != nil {
							log.Printf("[%s] Error processing trade through orchestrator: %v", clientIdentifier, err)
						} else {
							log.Printf("[%s] Successfully sent trade to orchestrator for Quote ID %s", clientIdentifier, confirmation.QuoteNonce)
						}
					} else {
						log.Printf("[%s] No asset mapping for %s. Cannot hedge.", clientIdentifier, confirmation.AssetAddress)
					}
				default:
					log.Printf("[%s] Received unhandled method '%s' with ID %s. Params: %s", clientIdentifier, baseResponse.Method, baseResponse.ID, string(baseResponse.Params))
				}
			} else {
				log.Printf("[%s] Received message with no error, result, or method. ID: %s. Raw: %s", clientIdentifier, baseResponse.ID, string(message))
			}
		})

		// Manage RFQ stream clients for this session
		var rfqClients []*ryskcore.Client
		var rfqClientCancels []context.CancelFunc

		if cfg.RFQAssetAddressesCSV != "" && baseURL != "" {
			for _, addr := range assetAddressesList {
				trimmedAddr := strings.TrimSpace(addr)
				if trimmedAddr == "" {
					continue
				}
				rfqStreamURL := fmt.Sprintf("%s/rfqs/%s", baseURL, trimmedAddr)

				// Use sessionCtx for RFQ streams
				rfqSdkClient, rfqClientCancel, err := SetupRfqStream(sessionCtx, rfqStreamURL, trimmedAddr)
				if err != nil {
					// Log and continue, one failed RFQ stream shouldn't stop the session if others can connect.
					// SetupRfqStream already logs its errors.
					continue
				}
				rfqClients = append(rfqClients, rfqSdkClient)
				rfqClientCancels = append(rfqClientCancels, rfqClientCancel)

				// Capture necessary variables for the handler closure
				currentAssetAddr := trimmedAddr       // Important to capture per iteration
				localRfqSdkClient := rfqSdkClient // Capture client for this handler

				localRfqSdkClient.SetHandler(func(message []byte) {
					HandleRfqMessage(message, currentAssetAddr, mainSdkClient, cfg, exchange) // mainSdkClient is from the outer scope
				})
			}
			log.Printf("RFQ stream setup complete for session. %d streams active.", len(rfqClients))
		}

		// Wait for session to end (main client disconnects) or app to terminate
		select {
		case <-appCtx.Done():
			log.Printf("Application context cancelled during active session (%v). Shutting down session.", appCtx.Err())
			// App will break from outer loop next
		case <-sessionCtx.Done(): // Could be triggered by sessionCancel() or if a critical error occurs within session logic
			log.Printf("Session context done (%v). Ending session and will attempt to restart.", sessionCtx.Err())
		case <-mainSdkClient.Ctx.Done(): // Main SDK client disconnected
			log.Printf("Main SDK client disconnected (%v). Ending session and will attempt to restart.", mainSdkClient.Ctx.Err())
		}

		// Cleanup for the current session
		log.Println("Cleaning up session...")
		for i, rfqClient := range rfqClients {
			if rfqClient != nil {
				rfqClient.Close() // Close the client
			}
			if rfqClientCancels[i] != nil {
				rfqClientCancels[i]() // Cancel its context
			}
		}
		mainSdkClient.Close() // Close main client for this session
		sessionCancel()       // Cancel the session context itself, ensuring all its children are also cancelled

		// If appCtx is done, we don't need to backoff and retry
		if appCtx.Err() != nil {
			break appRunningLoop
		}

		// Backoff before retrying the session
		log.Printf("Session ended. Retrying new session after backoff: %v", currentSessionBackoff)
			if !waitForBackoffOrSignal(appCtx, currentSessionBackoff) {
				break appRunningLoop // App context cancelled during backoff
			}
		currentSessionBackoff = nextBackoff(currentSessionBackoff, sessionInitialBackoff, sessionMaxBackoff)
	}

	log.Println("Quote Responder Daemon finished.")
}

// waitForBackoffOrSignal waits for a timer or context cancellation.
// Returns false if appCtx is cancelled during wait, true otherwise.
func waitForBackoffOrSignal(appCtx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-appCtx.Done():
		log.Printf("Application context cancelled during backoff wait: %v", appCtx.Err())
		return false
	case <-timer.C:
		return true
	}
}

// nextBackoff calculates the next backoff duration.
func nextBackoff(current, initial, max time.Duration) time.Duration {
	next := current * 2
	if next > max {
		return max
	}
	if next < initial { // Should not happen if current >= initial
	    return initial
	}
	return next
}
