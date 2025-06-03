package marketmaker

import "log"

// debugMode controls whether debug logging is enabled
var debugMode bool

// SetDebugMode enables or disables debug logging
func SetDebugMode(enabled bool) {
	debugMode = enabled
}

// DebugLog logs a message only if debug mode is enabled
func DebugLog(format string, args ...interface{}) {
	if debugMode {
		log.Printf("[DEBUG] "+format, args...)
	}
}

// EnableDebugMode enables debug logging
func EnableDebugMode() {
	SetDebugMode(true)
	log.Println("Debug mode enabled")
}

// DisableDebugMode disables debug logging
func DisableDebugMode() {
	SetDebugMode(false)
	log.Println("Debug mode disabled")
}
