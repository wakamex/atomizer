package main

import (
	"net/http"
	"time"
)

// newHTTPClient creates a standard HTTP client with timeout
func newHTTPClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Second}
}