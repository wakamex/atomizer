package monitor

import "log"

// Logger provides debug logging functionality
type Logger struct {
	debugEnabled bool
}

// NewLogger creates a new logger instance
func NewLogger(debug bool) *Logger {
	return &Logger{
		debugEnabled: debug,
	}
}

// Debugf logs a debug message if debug mode is enabled
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.debugEnabled {
		log.Printf(format, args...)
	}
}

// Printf always logs (for non-debug messages)
func (l *Logger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Global logger instance
var logger = NewLogger(false)

// SetDebug enables or disables debug logging
func SetDebug(enabled bool) {
	logger.debugEnabled = enabled
}

// Debugf logs a debug message if debug mode is enabled
func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

// Printf always logs (for non-debug messages)
func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}
