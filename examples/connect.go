package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/wakamex/ryskV12-cli/ryskcore" // SDK client
)

const (
	requestPipeSuffix  = ".req.pipe"
	responsePipeSuffix = ".res.pipe" // Will be used later
	pipePerm           = 0660        // Read/write for user and group
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds) // Add microseconds to log output
	channelID := flag.String("channel_id", "rysk_ipc_default", "Unique ID for the IPC channel (named pipes will be /tmp/<channel_id>.*.pipe)")
	websocketURL := flag.String("url", "wss://rip-testnet.rysk.finance/maker", "WebSocket URL to connect to")
	flag.Parse()

	log.Printf("Starting Rysk Connection Daemon for channel: %s, URL: %s", *channelID, *websocketURL)

	// Define pipe paths
	requestPipePath := filepath.Join(os.TempDir(), *channelID+requestPipeSuffix)
	responsePipePath := filepath.Join(os.TempDir(), *channelID+responsePipeSuffix) // For future use
	log.Printf("Request pipe: %s", requestPipePath)
	log.Printf("Response pipe: %s (will be used later)", responsePipePath)

	// Cleanup pipes on exit
	defer func() {
		log.Println("Cleaning up pipes...")
		os.Remove(requestPipePath)
		os.Remove(responsePipePath)
	}()

	// Create request pipe
	if err := syscall.Mkfifo(requestPipePath, pipePerm); err != nil {
		if !os.IsExist(err) {
			log.Fatalf("Failed to create request pipe %s: %v", requestPipePath, err)
		}
		log.Printf("Warning: Request pipe %s already exists or error creating: %v. Attempting to proceed.", requestPipePath, err)
	} else {
		log.Printf("Request pipe %s created successfully.", requestPipePath)
	}

	// Create response pipe (similar logic for existence) - for future use
	if err := syscall.Mkfifo(responsePipePath, pipePerm); err != nil {
		if !os.IsExist(err) {
			log.Fatalf("Failed to create response pipe %s: %v", responsePipePath, err)
		}
		log.Printf("Warning: Response pipe %s already exists or error creating: %v. Attempting to proceed.", responsePipePath, err)
	} else {
		log.Printf("Response pipe %s created successfully.", responsePipePath)
	}

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called eventually

	// --- SDK Client Initialization ---
	log.Printf("Attempting to connect to %s using Rysk SDK...", *websocketURL)
	sdkClient, err := ryskcore.NewClient(ctx, *websocketURL, nil)
	if err != nil {
		log.Fatalf("Failed to create/connect SDK client: %v", err)
	}
	log.Println("Successfully created SDK client and connected.")
	defer sdkClient.Close()

	// --- Set a Handler for Incoming Messages from WebSocket ---
	sdkClient.SetHandler(func(message []byte) {
		log.Printf("SDK Received (from WebSocket): %s", string(message))
		// TODO: Later, this handler will need to route responses back to the correct client via the response pipe
	})
	log.Println("SDK client is listening for messages from WebSocket.")

	// Start IPC request handler in a new goroutine
	go handleIPCRequests(ctx, requestPipePath, sdkClient)

	// Keep alive, wait for Ctrl+C to terminate this example program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	log.Println("Daemon running. WebSocket connection active. Listening for IPC requests. Press Ctrl+C to disconnect and exit.")

	select {
	case sig := <-interrupt:
		log.Printf("Interrupt signal %v received, initiating shutdown...", sig)
	case <-ctx.Done():
		log.Println("Main context cancelled, initiating shutdown...")
	case <-sdkClient.Ctx.Done(): // Listen for SDK's internal context doneness (e.g., connection lost)
		log.Printf("SDK client context done (%v), initiating shutdown...", sdkClient.Ctx.Err())
	}

	log.Println("Daemon shutting down.")
}

// handleIPCRequests opens the request pipe and processes incoming client requests.
func handleIPCRequests(ctx context.Context, pipePath string, sdkClient *ryskcore.Client) {
	log.Printf("IPC Request Handler: Starting for pipe: %s", pipePath)

	for { // Outer loop to allow reopening the pipe
		select {
		case <-ctx.Done():
			log.Println("IPC Request Handler: Context cancelled, shutting down.")
			return
		default:
		}

		log.Printf("IPC Request Handler: Attempting to open request pipe for reading: %s", pipePath)
		file, err := os.OpenFile(pipePath, os.O_RDONLY, 0) // Use O_RDONLY for reading
		if err != nil {
			log.Printf("IPC Request Handler: Failed to open request pipe %s: %v. Retrying in 5 seconds.", pipePath, err)
			time.Sleep(5 * time.Second)
			continue // Retry opening the pipe
		}
		log.Printf("IPC Request Handler: Request pipe %s opened successfully.", pipePath)

		// Inner loop for reading from the currently open pipe
		readLoop(ctx, file, sdkClient)

		log.Printf("IPC Request Handler: Closing pipe %s (readLoop exited).", pipePath)
		file.Close() // Close the pipe when readLoop exits

		select {
		case <-ctx.Done():
			log.Println("IPC Request Handler: Context cancelled after closing pipe, shutting down.")
			return
		default:
			log.Println("IPC Request Handler: Will attempt to reopen pipe.")
		}
	}
}

// readLoop continuously reads from the provided pipe file.
func readLoop(ctx context.Context, file *os.File, sdkClient *ryskcore.Client) {
	buf := make([]byte, 4096) // 4KB buffer
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
				select {
				case <-ctx.Done():
					log.Println("IPC Read Loop: Context cancelled during/after read error, exiting.")
					return
				default:
				}
				log.Printf("IPC Read Loop: Error reading from request pipe %s: %v. Exiting read loop.", file.Name(), err)
				return
			}

			if n > 0 {
				receivedRequest := string(buf[:n])
				log.Printf("IPC Request Handler: Received request: %s", receivedRequest)

				// TODO: Step 1.3 - Parse this request (ensure it's valid JSON-RPC)
				// TODO: Step 1.4 - Generate a unique ID for this request if not present, or use client's ID
				// TODO: Step 1.5 - Store a mapping of this request ID to the client's response pipe details
				// TODO: Step 1.6 - Forward the request to the WebSocket: sdkClient.Send([]byte(receivedRequest))

				if sdkClient != nil {
					log.Printf("IPC Request Handler: Forwarding to WebSocket: %s", receivedRequest)
					sdkClient.Send([]byte(receivedRequest)) // Actually send the request
				}
			}
		}
	}
}
