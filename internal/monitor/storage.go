package monitor

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type VMStorage struct {
	url    string
	client *http.Client
}

func NewVMStorage(url string) *VMStorage {
	return &VMStorage{
		url:    url,
		client: NewHTTPClient(30 * time.Second),
	}
}

func (v *VMStorage) Write(metrics []Metric) error {
	if len(metrics) == 0 {
		return nil
	}
	
	// Convert metrics to Prometheus format
	promData := v.toPrometheusFormat(metrics)
	
	// Send to VictoriaMetrics
	url := fmt.Sprintf("%s/api/v1/import/prometheus", v.url)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(promData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "text/plain")
	
	resp, err := v.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	
	return nil
}

func (v *VMStorage) toPrometheusFormat(metrics []Metric) string {
	var builder strings.Builder
	
	for _, m := range metrics {
		timestamp := m.Timestamp.UnixMilli()
		
		// Create labels
		labels := fmt.Sprintf(`exchange="%s",instrument="%s"`, m.Exchange, m.Instrument)
		
		// Write each metric field
		v.writeMetric(&builder, "market_bid_price", labels, m.BidPrice, timestamp)
		v.writeMetric(&builder, "market_ask_price", labels, m.AskPrice, timestamp)
		v.writeMetric(&builder, "market_bid_size", labels, m.BidSize, timestamp)
		v.writeMetric(&builder, "market_ask_size", labels, m.AskSize, timestamp)
		v.writeMetric(&builder, "market_last_price", labels, m.LastPrice, timestamp)
		v.writeMetric(&builder, "market_volume_24h", labels, m.Volume24h, timestamp)
		v.writeMetric(&builder, "market_open_price", labels, m.OpenPrice, timestamp)
		v.writeMetric(&builder, "market_high_price", labels, m.HighPrice, timestamp)
		v.writeMetric(&builder, "market_low_price", labels, m.LowPrice, timestamp)
		
		// Calculate and write spread
		if m.AskPrice > 0 && m.BidPrice > 0 {
			spread := m.AskPrice - m.BidPrice
			spreadPercent := (spread / m.BidPrice) * 100
			v.writeMetric(&builder, "market_spread", labels, spread, timestamp)
			v.writeMetric(&builder, "market_spread_percent", labels, spreadPercent, timestamp)
		}
	}
	
	return builder.String()
}

func (v *VMStorage) writeMetric(builder *strings.Builder, name string, labels string, value float64, timestamp int64) {
	if value != 0 { // Skip zero values to save space
		fmt.Fprintf(builder, "%s{%s} %f %d\n", name, labels, value, timestamp)
	}
}

func sanitizeMetricName(name string) string {
	// Replace non-alphanumeric characters with underscores
	var result strings.Builder
	for _, ch := range name {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' {
			result.WriteRune(ch)
		} else {
			result.WriteRune('_')
		}
	}
	return result.String()
}