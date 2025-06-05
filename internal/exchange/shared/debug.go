package shared

import (
    "log"
    "os"
)

// DebugLog logs debug messages if the specified environment variable is set
func DebugLog(envVar, format string, args ...interface{}) {
    if os.Getenv(envVar) != "" {
        log.Printf("[DEBUG] "+format, args...)
    }
}

// DeriveDebugLog logs debug messages if DERIVE_DEBUG is set
func DeriveDebugLog(format string, args ...interface{}) {
    DebugLog("DERIVE_DEBUG", "[DERIVE] "+format, args...)
}

// DeriveWSDebugLog logs WebSocket-specific debug messages if DERIVE_WS_DEBUG is set
func DeriveWSDebugLog(format string, args ...interface{}) {
    DebugLog("DERIVE_WS_DEBUG", "[DERIVE] "+format, args...)
}

// DeribitDebugLog logs debug messages if DERIBIT_DEBUG is set
func DeribitDebugLog(format string, args ...interface{}) {
    DebugLog("DERIBIT_DEBUG", "[DERIBIT] "+format, args...)
}