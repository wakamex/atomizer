package shared

import (
    "net/http"
    "time"
)

// NewHTTPClient creates a new HTTP client with default timeout
func NewHTTPClient() *http.Client {
    return &http.Client{
        Timeout: 30 * time.Second,
    }
}