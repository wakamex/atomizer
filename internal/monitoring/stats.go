package monitoring

import (
	"context"
	"log"
	"sync"
	"time"
)

// StatsReporter provides periodic statistics reporting
type StatsReporter struct {
	interval   time.Duration
	mu         sync.RWMutex
	stats      map[string]interface{}
	startTime  time.Time
	reportFunc func(map[string]interface{})
}

// NewStatsReporter creates a new stats reporter
func NewStatsReporter(interval time.Duration, reportFunc func(map[string]interface{})) *StatsReporter {
	return &StatsReporter{
		interval:   interval,
		stats:      make(map[string]interface{}),
		startTime:  time.Now(),
		reportFunc: reportFunc,
	}
}

// Start begins periodic reporting
func (s *StatsReporter) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.report()
		}
	}
}

// Set updates a stat value
func (s *StatsReporter) Set(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats[key] = value
}

// Increment increments a numeric stat
func (s *StatsReporter) Increment(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if val, ok := s.stats[key].(int); ok {
		s.stats[key] = val + 1
	} else {
		s.stats[key] = 1
	}
}

// report generates and logs statistics
func (s *StatsReporter) report() {
	s.mu.RLock()
	stats := make(map[string]interface{})
	for k, v := range s.stats {
		stats[k] = v
	}
	s.mu.RUnlock()
	
	// Add runtime stats
	stats["uptime_seconds"] = int64(time.Since(s.startTime).Seconds())
	stats["timestamp"] = time.Now().Unix()
	
	// Call report function
	if s.reportFunc != nil {
		s.reportFunc(stats)
	} else {
		// Default logging
		log.Printf("=== Statistics Report ===")
		log.Printf("Uptime: %v", time.Since(s.startTime).Round(time.Second))
		for k, v := range stats {
			if k != "uptime_seconds" && k != "timestamp" {
				log.Printf("%s: %v", k, v)
			}
		}
	}
}

// GlobalStats provides a global stats instance
var GlobalStats = NewStatsReporter(30*time.Second, nil)

// EnableDebugMode enables verbose logging
var DebugMode = false

// SetDebugMode toggles debug mode
func SetDebugMode(enabled bool) {
	DebugMode = enabled
	if enabled {
		log.Printf("Debug mode ENABLED")
	} else {
		log.Printf("Debug mode DISABLED")
	}
}

// DebugLog logs only in debug mode
func DebugLog(format string, args ...interface{}) {
	if DebugMode {
		log.Printf("[DEBUG] "+format, args...)
	}
}