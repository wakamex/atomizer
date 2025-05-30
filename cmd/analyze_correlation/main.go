package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// DataPoint represents a time series data point
type DataPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// QueryResult represents VictoriaMetrics query result
type QueryResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Values [][]interface{}   `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// CorrelationAnalyzer analyzes correlation between exchange quotes
type CorrelationAnalyzer struct {
	vmURL string
}

func NewCorrelationAnalyzer(vmURL string) *CorrelationAnalyzer {
	return &CorrelationAnalyzer{vmURL: vmURL}
}

// QueryMetric queries VictoriaMetrics for a specific metric
func (c *CorrelationAnalyzer) QueryMetric(query string, start, end time.Time, step time.Duration) (*QueryResult, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("start", fmt.Sprintf("%d", start.Unix()))
	params.Set("end", fmt.Sprintf("%d", end.Unix()))
	params.Set("step", fmt.Sprintf("%ds", int(step.Seconds())))

	url := fmt.Sprintf("%s/api/v1/query_range?%s", c.vmURL, params.Encode())
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query VictoriaMetrics: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result QueryResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("query failed with status: %s", result.Status)
	}

	return &result, nil
}

// ExtractTimeSeries extracts time series data from query result
func ExtractTimeSeries(result *QueryResult) map[string][]DataPoint {
	series := make(map[string][]DataPoint)
	
	for _, r := range result.Data.Result {
		key := fmt.Sprintf("%s_%s", r.Metric["exchange"], r.Metric["instrument"])
		points := make([]DataPoint, 0, len(r.Values))
		
		for _, v := range r.Values {
			if len(v) >= 2 {
				timestamp, _ := v[0].(float64)
				valueStr, _ := v[1].(string)
				value := 0.0
				fmt.Sscanf(valueStr, "%f", &value)
				
				points = append(points, DataPoint{
					Timestamp: int64(timestamp),
					Value:     value,
				})
			}
		}
		
		series[key] = points
	}
	
	return series
}

// AlignTimeSeries aligns two time series by timestamp
func AlignTimeSeries(series1, series2 []DataPoint) ([]DataPoint, []DataPoint) {
	// Create maps for quick lookup
	s1Map := make(map[int64]float64)
	s2Map := make(map[int64]float64)
	
	for _, p := range series1 {
		s1Map[p.Timestamp] = p.Value
	}
	
	for _, p := range series2 {
		s2Map[p.Timestamp] = p.Value
	}
	
	// Find common timestamps
	var aligned1, aligned2 []DataPoint
	
	for ts, v1 := range s1Map {
		if v2, exists := s2Map[ts]; exists {
			aligned1 = append(aligned1, DataPoint{Timestamp: ts, Value: v1})
			aligned2 = append(aligned2, DataPoint{Timestamp: ts, Value: v2})
		}
	}
	
	// Sort by timestamp
	sort.Slice(aligned1, func(i, j int) bool {
		return aligned1[i].Timestamp < aligned1[j].Timestamp
	})
	sort.Slice(aligned2, func(i, j int) bool {
		return aligned2[i].Timestamp < aligned2[j].Timestamp
	})
	
	return aligned1, aligned2
}

// CalculateCorrelation calculates Pearson correlation coefficient
func CalculateCorrelation(x, y []DataPoint) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}
	
	n := float64(len(x))
	
	// Calculate means
	var sumX, sumY float64
	for i := range x {
		sumX += x[i].Value
		sumY += y[i].Value
	}
	meanX := sumX / n
	meanY := sumY / n
	
	// Calculate correlation
	var num, denX, denY float64
	for i := range x {
		dx := x[i].Value - meanX
		dy := y[i].Value - meanY
		num += dx * dy
		denX += dx * dx
		denY += dy * dy
	}
	
	if denX == 0 || denY == 0 {
		return 0
	}
	
	return num / math.Sqrt(denX*denY)
}

// CalculateLag finds the lag that maximizes correlation
func CalculateLag(x, y []DataPoint, maxLag int) (int, float64) {
	bestLag := 0
	bestCorr := 0.0
	
	for lag := -maxLag; lag <= maxLag; lag++ {
		// Shift y by lag
		var xAligned, yAligned []DataPoint
		
		if lag >= 0 {
			// Positive lag: y leads x
			if lag < len(y) && lag < len(x) {
				xAligned = x[lag:]
				yAligned = y[:len(y)-lag]
			}
		} else {
			// Negative lag: x leads y
			absLag := -lag
			if absLag < len(x) && absLag < len(y) {
				xAligned = x[:len(x)-absLag]
				yAligned = y[absLag:]
			}
		}
		
		if len(xAligned) > 10 { // Need enough points
			corr := CalculateCorrelation(xAligned, yAligned)
			if math.Abs(corr) > math.Abs(bestCorr) {
				bestCorr = corr
				bestLag = lag
			}
		}
	}
	
	return bestLag, bestCorr
}

func main() {
	var (
		vmURL      = flag.String("vm-url", "http://localhost:8428", "VictoriaMetrics URL")
		start      = flag.String("start", "1h", "Start time (e.g., 1h, 24h, 2023-01-01T00:00:00Z)")
		end        = flag.String("end", "now", "End time")
		step       = flag.String("step", "10s", "Step interval")
		instrument = flag.String("instrument", "ETH-", "Instrument pattern to analyze")
		metric     = flag.String("metric", "market_bid_price", "Metric to analyze (market_bid_price, market_ask_price, etc.)")
	)
	flag.Parse()

	analyzer := NewCorrelationAnalyzer(*vmURL)

	// Parse time range
	endTime := time.Now()
	if *end != "now" {
		var err error
		endTime, err = time.Parse(time.RFC3339, *end)
		if err != nil {
			log.Fatalf("Failed to parse end time: %v", err)
		}
	}

	startTime := endTime
	if strings.HasSuffix(*start, "h") {
		hours := 1
		fmt.Sscanf(*start, "%dh", &hours)
		startTime = endTime.Add(-time.Duration(hours) * time.Hour)
	} else if strings.HasSuffix(*start, "d") {
		days := 1
		fmt.Sscanf(*start, "%dd", &days)
		startTime = endTime.Add(-time.Duration(days) * 24 * time.Hour)
	} else {
		var err error
		startTime, err = time.Parse(time.RFC3339, *start)
		if err != nil {
			log.Fatalf("Failed to parse start time: %v", err)
		}
	}

	stepDuration, err := time.ParseDuration(*step)
	if err != nil {
		log.Fatalf("Failed to parse step duration: %v", err)
	}

	// Query data
	query := fmt.Sprintf(`%s{instrument=~"%s.*"}`, *metric, *instrument)
	log.Printf("Querying: %s", query)
	log.Printf("Time range: %s to %s (step: %s)", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339), stepDuration)

	result, err := analyzer.QueryMetric(query, startTime, endTime, stepDuration)
	if err != nil {
		log.Fatalf("Failed to query metrics: %v", err)
	}

	// Extract time series
	series := ExtractTimeSeries(result)
	
	// Group by instrument
	instrumentMap := make(map[string]map[string][]DataPoint) // instrument -> exchange -> data
	for key, data := range series {
		parts := strings.Split(key, "_")
		if len(parts) >= 2 {
			exchange := parts[0]
			instrument := strings.Join(parts[1:], "_")
			
			if instrumentMap[instrument] == nil {
				instrumentMap[instrument] = make(map[string][]DataPoint)
			}
			instrumentMap[instrument][exchange] = data
		}
	}

	// Analyze correlations
	fmt.Println("\n=== Correlation Analysis ===")
	fmt.Printf("Metric: %s\n", *metric)
	fmt.Printf("Time range: %s to %s\n", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
	fmt.Println("\nInstruments found:")
	
	for instrument, exchanges := range instrumentMap {
		if len(exchanges) >= 2 {
			fmt.Printf("\n%s:\n", instrument)
			
			// Check if we have both Derive and Deribit data
			deriveData, hasDerive := exchanges["derive"]
			deribitData, hasDeribit := exchanges["deribit"]
			
			if hasDerive && hasDeribit {
				// Align time series
				alignedDerive, alignedDeribit := AlignTimeSeries(deriveData, deribitData)
				
				if len(alignedDerive) > 0 {
					// Calculate correlation
					correlation := CalculateCorrelation(alignedDerive, alignedDeribit)
					
					// Calculate lag
					lag, lagCorr := CalculateLag(alignedDerive, alignedDeribit, 10)
					
					fmt.Printf("  Derive vs Deribit:\n")
					fmt.Printf("    Data points: %d\n", len(alignedDerive))
					fmt.Printf("    Correlation: %.4f\n", correlation)
					fmt.Printf("    Best lag: %d steps (correlation: %.4f)\n", lag, lagCorr)
					
					if lag > 0 {
						fmt.Printf("    -> Deribit leads Derive by %d steps (~%s)\n", lag, time.Duration(lag)*stepDuration)
					} else if lag < 0 {
						fmt.Printf("    -> Derive leads Deribit by %d steps (~%s)\n", -lag, time.Duration(-lag)*stepDuration)
					} else {
						fmt.Printf("    -> No significant lag detected\n")
					}
					
					// Calculate average spread
					if len(alignedDerive) > 0 {
						var sumDiff, sumAbsDiff float64
						for i := range alignedDerive {
							diff := alignedDerive[i].Value - alignedDeribit[i].Value
							sumDiff += diff
							sumAbsDiff += math.Abs(diff)
						}
						avgDiff := sumDiff / float64(len(alignedDerive))
						avgAbsDiff := sumAbsDiff / float64(len(alignedDerive))
						
						fmt.Printf("    Average difference: %.4f (absolute: %.4f)\n", avgDiff, avgAbsDiff)
					}
				}
			} else {
				fmt.Printf("  Missing data from one exchange (Derive: %v, Deribit: %v)\n", hasDerive, hasDeribit)
			}
		}
	}
}