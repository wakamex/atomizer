package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/wakamex/rysk-v12-cli/ryskcore"
)

const (
	requestPipeSuffix  = ".req.pipe"
	responsePipeSuffix = ".res.pipe"
	pipePerm           = 0660
)

type jsonRPCResponse struct {
	ID json.RawMessage `json:"id"`
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
	channelID := flag.String("channel_id", "rysk_ipc_default", "Unique ID for the IPC channel")
	websocketURL := flag.String("url", "wss://rip-testnet.rysk.finance/maker", "WebSocket URL for the maker connection")
	rfqAssetAddressesCSV := flag.String("rfq_asset_addresses", "", "Comma-separated list of asset addresses for RFQ streams (e.g., ETH-PERP,BTC-PERP)")
	flag.Parse()

	log.Printf("Starting Rysk Connection Daemon for channel: %s, Maker URL: %s", *channelID, *websocketURL)

	requestPipePath := filepath.Join(os.TempDir(), *channelID+requestPipeSuffix)
	responsePipePath := filepath.Join(os.TempDir(), *channelID+responsePipeSuffix)
	log.Printf("Request pipe: %s", requestPipePath)
	log.Printf("Response pipe: %s", responsePipePath)

	defer func() {
		log.Println("Cleaning up pipes...")
		os.Remove(requestPipePath)
		os.Remove(responsePipePath)
	}()

	createPipeIfNotExists(requestPipePath, "Request")
	createPipeIfNotExists(responsePipePath, "Response")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var allClients []*ryskcore.Client
	defer func() {
		log.Println("Closing all SDK clients...")
		for _, client := range allClients {
			if client != nil {
				client.Close()
			}
		}
		log.Println("All SDK clients closed.")
	}()

	log.Printf("Attempting to connect to %s using Rysk SDK...", *websocketURL)
	sdkClient, err := ryskcore.NewClient(ctx, *websocketURL, nil)
	if err != nil {
		log.Fatalf("Failed to create/connect SDK client: %v", err)
	}
	allClients = append(allClients, sdkClient)
	log.Println("Successfully created SDK client and connected.")

	sdkClient.SetHandler(func(message []byte) {
		log.Printf("Maker SDK Received: %s", string(message))

		var resp jsonRPCResponse
		if err := json.Unmarshal(message, &resp); err == nil && resp.ID != nil {
			log.Printf("Detected JSON-RPC response with ID. Attempting to write to response pipe: %s", responsePipePath)

			pipe, openErr := os.OpenFile(responsePipePath, os.O_WRONLY, 0)
			if openErr != nil {
				log.Printf("Error opening response pipe %s for writing: %v", responsePipePath, openErr)
				return
			}
			defer pipe.Close()

			if _, writeErr := pipe.Write(message); writeErr != nil {
				log.Printf("Error writing response to pipe %s: %v", responsePipePath, writeErr)
			} else {
				log.Printf("Successfully wrote response to pipe %s", responsePipePath)
			}
		} else if err != nil {
			log.Printf("Message is not a JSON-RPC response with an ID, or unmarshal error: %v. Not routing to response pipe.", err)
		} else {
			// Message unmarshalled but had no ID, or ID was null.
			// This is normal for some messages like subscription confirmations.
		}
	})
	log.Println("Maker SDK client is listening for messages.")

	if *rfqAssetAddressesCSV != "" {
		baseURL := strings.TrimSuffix(*websocketURL, "/maker")
		if baseURL == *websocketURL {
			log.Printf("Warning: Could not reliably determine base URL from %s to construct RFQ stream URLs. Assuming it's a base URL.", *websocketURL)
		}
		assetAddresses := strings.Split(*rfqAssetAddressesCSV, ",")
		for _, addr := range assetAddresses {
			trimmedAddr := strings.TrimSpace(addr)
			if trimmedAddr == "" {
				continue
			}
			rfqStreamURL := fmt.Sprintf("%s/rfqs/%s", baseURL, trimmedAddr)
			log.Printf("Attempting to connect to RFQ Stream for %s: %s", trimmedAddr, rfqStreamURL)

			rfqClient, rfqErr := ryskcore.NewClient(ctx, rfqStreamURL, nil)
			if rfqErr != nil {
				log.Printf("Failed to create RFQ Listener SDK client for %s (%s): %v", trimmedAddr, rfqStreamURL, rfqErr)
				continue
			}
			allClients = append(allClients, rfqClient)
			currentAddr := trimmedAddr
			rfqClient.SetHandler(func(message []byte) {
				log.Printf("RFQ SDK Received (%s): %s", currentAddr, string(message))
				// Note: RFQ messages are typically not responses to client requests via IPC in this design
				// So, not attempting to write them to the main responsePipePath here.
			})
			log.Printf("RFQ Listener for %s connected and listening.", currentAddr)
		}
	}

	go handleIPCRequests(ctx, requestPipePath, sdkClient)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	log.Println("Daemon running. WebSocket connection active. Listening for IPC requests. Press Ctrl+C to disconnect and exit.")

	select {
	case sig := <-interrupt:
		log.Printf("Interrupt signal %v received, initiating shutdown...", sig)
	case <-ctx.Done():
		log.Println("Main context cancelled, initiating shutdown...")
	case <-sdkClient.Ctx.Done():
		log.Printf("SDK client context done (%v), initiating shutdown...", sdkClient.Ctx.Err())
	}

	log.Println("Daemon shutting down.")
}

func createPipeIfNotExists(pipePath string, pipeNameForLog string) {
	if err := syscall.Mkfifo(pipePath, pipePerm); err != nil {
		if !os.IsExist(err) {
			log.Fatalf("Failed to create %s pipe %s: %v", pipeNameForLog, pipePath, err)
		}
		log.Printf("Warning: %s pipe %s already exists or error creating: %v. Attempting to proceed.", pipeNameForLog, pipePath, err)
	} else {
		log.Printf("%s pipe %s created successfully.", pipeNameForLog, pipePath)
	}
}

func handleIPCRequests(ctx context.Context, pipePath string, sdkClient *ryskcore.Client) {
	log.Printf("IPC Request Handler: Starting for pipe: %s", pipePath)

	for {
		select {
		case <-ctx.Done():
			log.Println("IPC Request Handler: Context cancelled, shutting down.")
			return
		default:
		}

		log.Printf("IPC Request Handler: Attempting to open request pipe for reading: %s", pipePath)
		file, err := os.OpenFile(pipePath, os.O_RDONLY, 0)
		if err != nil {
			log.Printf("IPC Request Handler: Failed to open request pipe %s: %v. Retrying in 5 seconds.", pipePath, err)
			time.Sleep(5 * time.Second)
			continue
		}
		log.Printf("IPC Request Handler: Request pipe %s opened successfully.", pipePath)

		readLoop(ctx, file, sdkClient)

		log.Printf("IPC Request Handler: Closing pipe %s (readLoop exited).", pipePath)
		file.Close()

		select {
		case <-ctx.Done():
			log.Println("IPC Request Handler: Context cancelled after closing pipe, shutting down.")
			return
		default:
			log.Println("IPC Request Handler: Will attempt to reopen pipe.")
		}
	}
}

func readLoop(ctx context.Context, file *os.File, sdkClient *ryskcore.Client) {
	buf := make([]byte, 4096)
	log.Printf("IPC Read Loop: Starting for pipe: %s", file.Name())
	for {
		select {
		case <-ctx.Done():
			log.Println("IPC Read Loop: Context cancelled, exiting.")
			return
		default:
			n, err := file.Read(buf)
			if err != nil {
				if err == io.EOF {
					log.Println("IPC Read Loop: EOF on request pipe. Assuming writer closed. Exiting read loop to allow pipe reopen.")
					return
				}
				log.Printf("IPC Read Loop: Error reading from request pipe %s: %v. Exiting read loop.", file.Name(), err)
				return
			}

			if n > 0 {
				receivedRequest := string(buf[:n])
				log.Printf("IPC Request Handler: Received request: %s", receivedRequest)
				if sdkClient != nil {
					log.Printf("IPC Request Handler: Forwarding to WebSocket: %s", receivedRequest)
					sdkClient.Send([]byte(receivedRequest))
				}
			}
		}
	}
}
