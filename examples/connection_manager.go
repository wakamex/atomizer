package main

import (
	"context"
	"log"
	"time"

	"github.com/wakamex/ryskV12-cli/ryskcore"
)

const (
	initialBackoff = 1 * time.Second
	maxBackoff     = 30 * time.Second
)

// EstablishConnectionWithRetry attempts to establish a WebSocket connection using ryskcore.NewClient,
// retrying with exponential backoff on failure.
// It returns the client or an error if the context is cancelled or retries are exhausted (though this impl retries indefinitely until context cancel).
func EstablishConnectionWithRetry(ctx context.Context, url string, clientIdentifier string) (*ryskcore.Client, error) {
	currentBackoff := initialBackoff
	for {
		select {
		case <-ctx.Done():
			log.Printf("[%s] Context cancelled, stopping connection attempts to %s.", clientIdentifier, url)
			return nil, ctx.Err()
		default:
		}

		log.Printf("[%s] Attempting to connect to %s...", clientIdentifier, url)
		client, err := ryskcore.NewClient(ctx, url, nil) // Pass the main context
		if err == nil {
			log.Printf("[%s] Successfully connected to %s.", clientIdentifier, url)
			return client, nil
		}

		log.Printf("[%s] Failed to connect to %s: %v. Retrying in %v...", clientIdentifier, url, err, currentBackoff)

		// Wait for currentBackoff duration or until context is cancelled
		timer := time.NewTimer(currentBackoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Printf("[%s] Context cancelled during backoff for %s.", clientIdentifier, url)
			return nil, ctx.Err()
		case <-timer.C:
		}

		currentBackoff *= 2
		if currentBackoff > maxBackoff {
			currentBackoff = maxBackoff
		}
	}
}

// SetupRfqStream attempts to establish a WebSocket connection for a specific RFQ stream.
// It returns the client, a context cancel function for this specific client, and an error.
// The caller is responsible for calling the returned cancel function and closing the client when done.
func SetupRfqStream(parentCtx context.Context, rfqStreamURL string, assetAddr string) (*ryskcore.Client, context.CancelFunc, error) {
	// Create a new context for this specific RFQ client, derived from the parent context.
	// This allows individual RFQ client goroutines to be cancelled without affecting others or the main client.
	rfqClientCtx, rfqClientCancel := context.WithCancel(parentCtx)

	log.Printf("[RFQ %s] Attempting to connect to stream %s...", assetAddr, rfqStreamURL)
	client, err := ryskcore.NewClient(rfqClientCtx, rfqStreamURL, nil)
	if err != nil {
		log.Printf("[RFQ %s] Failed to connect to stream %s: %v", assetAddr, rfqStreamURL, err)
		rfqClientCancel() // Important: Cancel the context if the connection failed immediately.
		return nil, nil, err
	}

	log.Printf("[RFQ %s] Successfully connected to stream %s.", assetAddr, rfqStreamURL)
	return client, rfqClientCancel, nil
}
