package rfq

import "time"

// Config holds RFQ responder configuration
type Config struct {
	// Exchange configuration
	Exchange string
	TestMode bool
	
	// Authentication
	DerivePrivateKey string
	DeriveWallet     string
	DeribitKey       string
	DeribitSecret    string
	
	// Operational parameters
	HeartbeatInterval time.Duration
	MaxResponseTime   time.Duration
	
	// Logging
	Verbose bool
	Debug   bool
}