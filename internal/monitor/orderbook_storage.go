package monitor

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// OrderBookStorage handles storing order book data to VictoriaMetrics
type OrderBookStorage struct {
	url    string
	client *http.Client
}

func NewOrderBookStorage(url string) *OrderBookStorage {
	return &OrderBookStorage{
		url:    url,
		client: NewHTTPClient(30 * time.Second),
	}
}

func (s *OrderBookStorage) WriteOrderBooks(orderbooks []OrderBookMetric) error {
	if len(orderbooks) == 0 {
		return nil
	}

	// Convert to Prometheus format
	promData := s.toPrometheusFormat(orderbooks)

	// Send to VictoriaMetrics
	url := fmt.Sprintf("%s/api/v1/import/prometheus", s.url)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(promData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *OrderBookStorage) toPrometheusFormat(orderbooks []OrderBookMetric) string {
	var builder strings.Builder

	for _, ob := range orderbooks {
		timestamp := ob.Timestamp.UnixMilli()
		baseLabels := fmt.Sprintf(`exchange="%s",instrument="%s"`, ob.Exchange, ob.Instrument)

		// Write bid levels
		for i, bid := range ob.Bids {
			level := i + 1
			labels := fmt.Sprintf(`%s,side="bid",level="%d"`, baseLabels, level)
			s.writeMetric(&builder, "orderbook_price", labels, bid.Price, timestamp)
			s.writeMetric(&builder, "orderbook_size", labels, bid.Size, timestamp)
			if bid.Orders > 0 {
				s.writeMetric(&builder, "orderbook_orders", labels, float64(bid.Orders), timestamp)
			}
		}

		// Write ask levels
		for i, ask := range ob.Asks {
			level := i + 1
			labels := fmt.Sprintf(`%s,side="ask",level="%d"`, baseLabels, level)
			s.writeMetric(&builder, "orderbook_price", labels, ask.Price, timestamp)
			s.writeMetric(&builder, "orderbook_size", labels, ask.Size, timestamp)
			if ask.Orders > 0 {
				s.writeMetric(&builder, "orderbook_orders", labels, float64(ask.Orders), timestamp)
			}
		}

		// Calculate and write aggregate metrics
		if len(ob.Bids) > 0 && len(ob.Asks) > 0 {
			// Best bid/ask spread
			spread := ob.Asks[0].Price - ob.Bids[0].Price
			spreadPercent := (spread / ob.Bids[0].Price) * 100
			s.writeMetric(&builder, "orderbook_spread", baseLabels, spread, timestamp)
			s.writeMetric(&builder, "orderbook_spread_percent", baseLabels, spreadPercent, timestamp)

			// Mid price
			midPrice := (ob.Bids[0].Price + ob.Asks[0].Price) / 2
			s.writeMetric(&builder, "orderbook_mid_price", baseLabels, midPrice, timestamp)

			// Total bid/ask depth (sum of all levels)
			var totalBidSize, totalAskSize float64
			for _, bid := range ob.Bids {
				totalBidSize += bid.Size
			}
			for _, ask := range ob.Asks {
				totalAskSize += ask.Size
			}
			s.writeMetric(&builder, "orderbook_total_bid_size", baseLabels, totalBidSize, timestamp)
			s.writeMetric(&builder, "orderbook_total_ask_size", baseLabels, totalAskSize, timestamp)

			// Depth at specific price distances (e.g., within 0.5%, 1%, 2% of mid)
			s.writeDepthMetrics(&builder, baseLabels, ob, midPrice, timestamp)
		}
	}

	return builder.String()
}

func (s *OrderBookStorage) writeDepthMetrics(builder *strings.Builder, baseLabels string, ob OrderBookMetric, midPrice float64, timestamp int64) {
	// Calculate cumulative depth at various price distances
	distances := []float64{0.001, 0.0025, 0.005, 0.01, 0.02} // 0.1%, 0.25%, 0.5%, 1%, 2%

	for _, distance := range distances {
		bidThreshold := midPrice * (1 - distance)
		askThreshold := midPrice * (1 + distance)

		var bidDepth, askDepth float64

		// Sum bid depth within threshold
		for _, bid := range ob.Bids {
			if bid.Price >= bidThreshold {
				bidDepth += bid.Size
			}
		}

		// Sum ask depth within threshold
		for _, ask := range ob.Asks {
			if ask.Price <= askThreshold {
				askDepth += ask.Size
			}
		}

		distanceLabel := fmt.Sprintf(`%s,distance="%.2f%%"`, baseLabels, distance*100)
		s.writeMetric(builder, "orderbook_bid_depth", distanceLabel, bidDepth, timestamp)
		s.writeMetric(builder, "orderbook_ask_depth", distanceLabel, askDepth, timestamp)
	}
}

func (s *OrderBookStorage) writeMetric(builder *strings.Builder, name string, labels string, value float64, timestamp int64) {
	if value != 0 { // Skip zero values to save space
		fmt.Fprintf(builder, "%s{%s} %f %d\n", name, labels, value, timestamp)
	}
}
