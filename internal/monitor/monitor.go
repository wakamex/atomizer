package monitor

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Config struct {
	Interval           time.Duration
	Exchanges          []string
	InstrumentPatterns []string
	VictoriaMetricsURL string
	Workers            int
}

type Monitor struct {
	config        *Config
	collectors    map[string]Collector
	spotCollector *DeriveSpotCollector
	storage       *VMStorage
	vmProcess     *exec.Cmd
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

type Collector interface {
	Name() string
	Collect(ctx context.Context, instruments []string) ([]Metric, error)
}

type Metric struct {
	Exchange   string
	Instrument string
	Timestamp  time.Time
	BidPrice   float64
	AskPrice   float64
	BidSize    float64
	AskSize    float64
	LastPrice  float64
	Volume24h  float64
	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
}

func New(config *Config) (*Monitor, error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	m := &Monitor{
		config:     config,
		collectors: make(map[string]Collector),
		ctx:        ctx,
		cancel:     cancel,
	}

	// Initialize collectors based on configured exchanges
	for _, exchange := range config.Exchanges {
		switch exchange {
		case "derive":
			m.collectors[exchange] = NewDeriveCollector()
		case "deribit":
			m.collectors[exchange] = NewDeribitCollector()
		default:
			return nil, fmt.Errorf("unknown exchange: %s", exchange)
		}
	}

	// Initialize storage
	m.storage = NewVMStorage(config.VictoriaMetricsURL)

	// Initialize spot collector
	spotCollector, err := NewDeriveSpotCollector()
	if err != nil {
		log.Printf("Warning: Failed to create spot collector: %v", err)
		// Don't fail completely, just log the error
	} else {
		m.spotCollector = spotCollector
		// Subscribe to ETH and BTC spot feeds
		if err := spotCollector.Subscribe([]string{"ETH", "BTC"}); err != nil {
			log.Printf("Warning: Failed to subscribe to spot feeds: %v", err)
		}
	}

	return m, nil
}

func (m *Monitor) Start() error {
	// Check if VictoriaMetrics is already running
	if m.config.VictoriaMetricsURL == "http://localhost:8428" {
		if !m.isVictoriaMetricsRunning() {
			if err := m.startVictoriaMetrics(); err != nil {
				return fmt.Errorf("failed to start VictoriaMetrics: %w", err)
			}
			// Give VM time to start
			time.Sleep(2 * time.Second)
		} else {
			log.Println("VictoriaMetrics is already running")
		}
	}

	// Start collection loops for each exchange
	for _, collector := range m.collectors {
		m.wg.Add(1)
		go m.collectionLoop(collector)
	}

	// Start spot price collection loop
	if m.spotCollector != nil {
		m.wg.Add(1)
		go m.spotCollectionLoop()
	}

	return nil
}

func (m *Monitor) Stop() error {
	m.cancel()
	m.wg.Wait()

	if m.spotCollector != nil {
		if err := m.spotCollector.Close(); err != nil {
			log.Printf("Error closing spot collector: %v", err)
		}
	}

	if m.vmProcess != nil {
		log.Println("Stopping VictoriaMetrics...")
		if err := m.vmProcess.Process.Signal(os.Interrupt); err != nil {
			m.vmProcess.Process.Kill()
		}
		m.vmProcess.Wait()
	}

	return nil
}

func (m *Monitor) isVictoriaMetricsRunning() bool {
	resp, err := http.Get(m.config.VictoriaMetricsURL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (m *Monitor) startVictoriaMetrics() error {
	installer := NewVMInstaller()
	binaryPath := installer.GetBinaryPath()
	dataPath := installer.GetDataPath()

	if _, err := os.Stat(binaryPath); err != nil {
		return fmt.Errorf("VictoriaMetrics not installed. Run 'atomizer market-monitor setup' first")
	}

	log.Printf("Starting VictoriaMetrics with data path: %s", dataPath)
	m.vmProcess = exec.Command(binaryPath, 
		"-storageDataPath", dataPath,
		"-retentionPeriod", "12",
		"-search.maxStalenessInterval", "5m",
	)
	
	m.vmProcess.Stdout = os.Stdout
	m.vmProcess.Stderr = os.Stderr
	
	if err := m.vmProcess.Start(); err != nil {
		return fmt.Errorf("failed to start VictoriaMetrics: %w", err)
	}

	return nil
}

func (m *Monitor) collectionLoop(collector Collector) {
	defer m.wg.Done()
	
	// Add jitter to avoid all collectors hitting APIs at the same time
	jitter := time.Duration(float64(m.config.Interval) * 0.1)
	time.Sleep(time.Duration(rand.Int63n(int64(jitter))))
	
	ticker := time.NewTicker(m.config.Interval)
	defer ticker.Stop()

	// Rate limiter: max 10 requests per second per collector
	rateLimiter := time.NewTicker(time.Second / 10)
	defer rateLimiter.Stop()

	log.Printf("Starting collection for %s every %v", collector.Name(), m.config.Interval)

	// Initial collection
	<-rateLimiter.C
	m.collect(collector)

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			<-rateLimiter.C
			m.collect(collector)
		}
	}
}

func (m *Monitor) collect(collector Collector) {
	metrics, err := collector.Collect(m.ctx, m.config.InstrumentPatterns)
	if err != nil {
		log.Printf("Collection error for %s: %v", collector.Name(), err)
		return
	}

	if err := m.storage.Write(metrics); err != nil {
		log.Printf("Storage error: %v", err)
		return
	}

	log.Printf("Collected %d metrics from %s", len(metrics), collector.Name())
}

func (m *Monitor) GetStats() (map[string]interface{}, error) {
	// Query VictoriaMetrics for stats
	stats := make(map[string]interface{})
	
	// TODO: Implement stats queries using PromQL
	// - Total data points
	// - Data points per exchange
	// - Latest collection timestamps
	// - Storage size
	
	return stats, nil
}

func (m *Monitor) Export(format string, output string, start, end time.Time) error {
	// TODO: Implement data export
	// - Query data from VictoriaMetrics
	// - Format as CSV, JSON, or Parquet
	// - Write to output file
	
	return fmt.Errorf("export not yet implemented")
}

func (m *Monitor) spotCollectionLoop() {
	defer m.wg.Done()
	
	// Use a faster interval for spot prices as they change frequently
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	log.Printf("Starting spot price collection every 10s")
	
	// Initial collection
	m.collectSpotPrices()
	
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.collectSpotPrices()
		}
	}
}

func (m *Monitor) collectSpotPrices() {
	spotPrices := m.spotCollector.GetAllSpotPrices()
	
	if len(spotPrices) == 0 {
		return
	}
	
	// Convert spot prices to metrics
	metrics := make([]Metric, 0, len(spotPrices))
	for currency, spot := range spotPrices {
		// Create a metric for the spot price
		metric := Metric{
			Exchange:   "derive",
			Instrument: fmt.Sprintf("%s-SPOT", currency),
			Timestamp:  spot.Timestamp,
			BidPrice:   spot.Price,
			AskPrice:   spot.Price,
			LastPrice:  spot.Price,
		}
		metrics = append(metrics, metric)
	}
	
	// Write to storage
	if err := m.storage.Write(metrics); err != nil {
		log.Printf("Failed to write spot prices: %v", err)
		return
	}
	
	log.Printf("Collected spot prices: %v", spotPrices)
}