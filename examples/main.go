package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rysk-finance/rysk-v12-cli/ryskcore" // SDK client
)

func main() {
	targetURL := "wss://rip-testnet.rysk.finance/maker"

	log.Println("Starting minimal Rysk SDK example...")

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called eventually

	// --- SDK Client Initialization ---
	log.Printf("Attempting to connect to %s using Rysk SDK...", targetURL)
	// Assuming NewClient doesn't require specific headers for a basic connection (passing nil)
	sdkClient, err := ryskcore.NewClient(ctx, targetURL, nil)
	if err != nil {
		log.Fatalf("Failed to create/connect SDK client: %v", err)
	}
	log.Println("Successfully created SDK client and connected.")

	// --- Set a Handler for Incoming Messages ---
	sdkClient.SetHandler(func(message []byte) {
		log.Printf("SDK Received: %s", string(message))
	})
	log.Println("SDK client is listening for messages.")

	// Keep alive, wait for Ctrl+C to terminate this example program
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	log.Println("SDK Connection active. Press Ctrl+C to disconnect and exit.")

	select {
	case <-interrupt:
		log.Println("Interrupt signal received, initiating shutdown...")
	case <-ctx.Done(): // Handles cancellation if the main context passed to NewClient is cancelled elsewhere
		log.Println("Main context cancelled, initiating shutdown...")
	case <-sdkClient.Ctx.Done(): // Listen for SDK's internal context doneness (e.g., connection lost)
		log.Println("SDK client context done, initiating shutdown...")
	}

	// --- Graceful Shutdown of SDK Client ---
	log.Println("Closing SDK client...")
	if err := sdkClient.Close(); err != nil {
		log.Printf("Error closing SDK client: %v", err)
	} else {
		log.Println("SDK client closed successfully.")
	}

	log.Println("Program exited.")
}
