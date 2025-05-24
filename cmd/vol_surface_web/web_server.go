package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

// WebServer serves the volatility surface visualization
type WebServer struct {
	analyzer        *VolSurfaceAnalyzer
	model           *VolSurfaceModel
	filteredOptions []OptionData
}

// VolSurfaceData represents data for web visualization
type VolSurfaceData struct {
	Points    []VolSurfacePoint   `json:"points"`
	Summary   VolSurfaceSummary   `json:"summary"`
	SpotPrice float64             `json:"spotPrice"`
}

// NewWebServer creates a new web server instance
func NewWebServer() *WebServer {
	return &WebServer{}
}

// LoadAndAnalyze loads data and performs volatility surface analysis
func (ws *WebServer) LoadAndAnalyze() error {
	log.Println("Loading data for web visualization...")
	
	// Create analyzer and load data
	ws.analyzer = NewVolSurfaceAnalyzer()
	
	if err := ws.analyzer.LoadCSV("eth_options_data.csv"); err != nil {
		return fmt.Errorf("failed to load CSV: %v", err)
	}
	
	// Update pricing and calculate IVs
	if err := ws.analyzer.UpdatePricing(); err != nil {
		log.Printf("Warning: Failed to update pricing: %v", err)
	}
	
	ws.analyzer.CalculateImpliedVolatilities()
	
	// Clean the data and track filtered options
	_, ws.filteredOptions = ws.analyzer.CleanData()
	
	// Fit SVI surface
	if err := ws.analyzer.FitSVISurface(); err != nil {
		log.Printf("Warning: Failed to fit SVI surface: %v", err)
		// Fall back to raw data visualization
	}
	
	// Check for arbitrage
	violations := ws.analyzer.CheckArbitrage()
	if len(violations) > 0 {
		log.Printf("⚠️  Found %d arbitrage violations", len(violations))
	}
	
	// Fit volatility surface (for compatibility with existing visualization)
	model := ws.analyzer.FitVolatilitySurface()
	ws.model = &model
	
	log.Printf("Loaded %d options with %d valid IV points", len(ws.analyzer.options), len(model.Points))
	return nil
}

// Start starts the web server
func (ws *WebServer) Start(port int) error {
	// Serve static files and API endpoints
	http.HandleFunc("/", ws.handleHome)
	http.HandleFunc("/api/surface", ws.handleSurfaceData)
	http.HandleFunc("/api/points", ws.handlePointsData)
	http.HandleFunc("/api/fitted-surface", ws.handleFittedSurface)
	
	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting web server on http://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleHome serves the main HTML page
func (ws *WebServer) handleHome(w http.ResponseWriter, r *http.Request) {
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <title>ETH Options Volatility Surface</title>
    <script src="https://cdn.plot.ly/plotly-latest.min.js"></script>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; background: #f5f5f5; }
        .container { max-width: 1400px; margin: 0 auto; }
        .header { background: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .charts { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; margin-bottom: 20px; }
        .chart-container { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .full-width { grid-column: 1 / -1; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 15px; }
        .stat-card { background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center; }
        .stat-value { font-size: 24px; font-weight: bold; color: #2563eb; }
        .stat-label { font-size: 14px; color: #6b7280; margin-top: 5px; }
        .loading { text-align: center; padding: 50px; }
        h1 { color: #1f2937; margin: 0; }
        h2 { color: #374151; margin-top: 0; }
        .legend { background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); margin-bottom: 20px; }
        .legend-item { display: inline-flex; align-items: center; margin-right: 20px; margin-bottom: 5px; }
        .legend-symbol { width: 20px; height: 20px; margin-right: 8px; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 ETH Options Volatility Surface</h1>
            <p>Real-time analysis of Ethereum options market structure and implied volatility dynamics</p>
        </div>
        
        <div id="loading" class="loading">
            <h2>Loading volatility surface data...</h2>
        </div>
        
        <div id="content" style="display: none;">
            <div class="stats" id="stats"></div>
            
            <div class="legend">
                <strong>Point Legend:</strong>
                <div style="margin-top: 10px;">
                    <div class="legend-item">
                        <div class="legend-symbol" style="background-color: #3b82f6;"></div>
                        <span>Included Calls (Blue circles)</span>
                    </div>
                    <div class="legend-item">
                        <div class="legend-symbol" style="background-color: #10b981;"></div>
                        <span>Included Puts (Green squares)</span>
                    </div>
                    <div class="legend-item">
                        <div class="legend-symbol" style="background-color: #ef4444;"></div>
                        <span>Filtered Calls (Red X)</span>
                    </div>
                    <div class="legend-item">
                        <div class="legend-symbol" style="background-color: #f97316;"></div>
                        <span>Filtered Puts (Orange X)</span>
                    </div>
                </div>
            </div>
            
            <div class="charts">
                <div class="chart-container full-width">
                    <h2>3D Volatility Surface</h2>
                    <div id="surface3d"></div>
                </div>
                
                <div class="chart-container">
                    <h2>Volatility vs Strike (by Expiry)</h2>
                    <div id="volStrike"></div>
                </div>
                
                <div class="chart-container">
                    <h2>Volatility Term Structure</h2>
                    <div id="termStructure"></div>
                </div>
                
                <div class="chart-container full-width">
                    <h2>Individual IV Points</h2>
                    <div id="scatterPoints"></div>
                </div>
            </div>
        </div>
    </div>

    <script>
        // Load and visualize data
        async function loadData() {
            try {
                const response = await fetch('/api/surface');
                const data = await response.json();
                
                // Also fetch fitted surface
                const fittedResponse = await fetch('/api/fitted-surface');
                const fittedData = await fittedResponse.json();
                
                document.getElementById('loading').style.display = 'none';
                document.getElementById('content').style.display = 'block';
                
                createStats(data);
                create3DSurface(data, fittedData);
                createVolStrikeChart(data);
                createTermStructure(data);
                createScatterPoints(data);
            } catch (error) {
                console.error('Error loading data:', error);
                document.getElementById('loading').innerHTML = '<h2>Error loading data: ' + error.message + '</h2>';
            }
        }
        
        function createStats(data) {
            const stats = document.getElementById('stats');
            const summary = data.summary || {};
            
            // Helper function to safely format numbers
            const safeFormat = (val, decimals = 1, suffix = '') => {
                if (val == null || isNaN(val) || !isFinite(val)) return 'N/A';
                return (val * (suffix === '%' ? 100 : 1)).toFixed(decimals) + suffix;
            };
            
            const safeInt = (val) => {
                if (val == null || isNaN(val)) return 'N/A';
                return Math.round(val).toString();
            };
            
            stats.innerHTML = ` + "`" + `
                <div class="stat-card">
                    <div class="stat-value">${safeInt(summary.totalPoints)}</div>
                    <div class="stat-label">Total IV Points</div>
                </div>
                <div class="stat-card">
                    <div class="stat-value">${safeFormat(summary.avgIV, 1, '%')}</div>
                    <div class="stat-label">Average IV</div>
                </div>
                <div class="stat-card">
                    <div class="stat-value">${safeFormat(summary.minIV, 1, '%')}</div>
                    <div class="stat-label">Min IV</div>
                </div>
                <div class="stat-card">
                    <div class="stat-value">${safeFormat(summary.maxIV, 1, '%')}</div>
                    <div class="stat-label">Max IV</div>
                </div>
                <div class="stat-card">
                    <div class="stat-value">$${safeFormat(data.spotPrice, 0)}</div>
                    <div class="stat-label">ETH Spot Price</div>
                </div>
                <div class="stat-card">
                    <div class="stat-value">${safeInt(summary.callPoints)}/${safeInt(summary.putPoints)}</div>
                    <div class="stat-label">Calls/Puts</div>
                </div>
            ` + "`" + `;
        }
        
        function create3DSurface(data, fittedData) {
            const points = data.points || [];
            console.log('Creating 3D surface with', points.length, 'points');
            console.log('Has SVI fit:', fittedData.hasSVIFit);
            
            const traces = [];
            
            // Add fitted surface if available
            if (fittedData && fittedData.hasSVIFit && fittedData.ivGrid) {
                const gridSize = fittedData.gridSize;
                const strikes = [];
                const times = [];
                
                // Generate strike and time arrays
                for (let i = 0; i < gridSize; i++) {
                    strikes.push(fittedData.strikeRange[0] + i * (fittedData.strikeRange[1] - fittedData.strikeRange[0]) / (gridSize - 1));
                    times.push(fittedData.timeRange[0] + i * (fittedData.timeRange[1] - fittedData.timeRange[0]) / (gridSize - 1));
                }
                
                traces.push({
                    x: strikes,
                    y: times,
                    z: fittedData.ivGrid,
                    type: 'surface',
                    colorscale: 'Viridis',
                    showscale: true,
                    colorbar: { title: 'Fitted IV (%)' },
                    name: 'SVI Fitted Surface',
                    opacity: 0.8
                });
            }
            
            // Add included data points (used for surface fitting)
            if (fittedData && fittedData.rawPoints && fittedData.rawPoints.length > 0) {
                const rawPoints = fittedData.rawPoints;
                traces.push({
                    x: rawPoints.map(p => p.strike),
                    y: rawPoints.map(p => p.timeToExpiry),
                    z: rawPoints.map(p => p.iv),
                    mode: 'markers',
                    type: 'scatter3d',
                    marker: {
                        size: 4,
                        color: rawPoints.map(p => p.optionType === 'C' ? '#3b82f6' : '#10b981'), // Blue for calls, green for puts
                        symbol: rawPoints.map(p => p.optionType === 'C' ? 'circle' : 'square'),
                        line: {
                            color: 'white',
                            width: 0.5
                        }
                    },
                    name: 'Included Points',
                    text: rawPoints.map(p => 'Strike: ' + p.strike + '<br>TTM: ' + p.timeToExpiry.toFixed(3) + '<br>IV: ' + p.iv.toFixed(1) + '%<br>Type: ' + p.optionType + '<br>Status: INCLUDED')
                });
            }
            
            // Add filtered data points (excluded from surface fitting)
            if (fittedData && fittedData.filteredPoints && fittedData.filteredPoints.length > 0) {
                const filteredPoints = fittedData.filteredPoints;
                traces.push({
                    x: filteredPoints.map(p => p.strike),
                    y: filteredPoints.map(p => p.timeToExpiry),
                    z: filteredPoints.map(p => p.iv),
                    mode: 'markers',
                    type: 'scatter3d',
                    marker: {
                        size: 3,
                        color: filteredPoints.map(p => p.optionType === 'C' ? '#ef4444' : '#f97316'), // Red for calls, orange for puts
                        symbol: 'x',
                        opacity: 0.6
                    },
                    name: 'Filtered Out',
                    text: filteredPoints.map(p => 'Strike: ' + p.strike + '<br>TTM: ' + p.timeToExpiry.toFixed(3) + '<br>IV: ' + p.iv.toFixed(1) + '%<br>Type: ' + p.optionType + '<br>Status: FILTERED')
                });
            }
            
            const layout = {
                scene: {
                    xaxis: { title: 'Strike Price ($)' },
                    yaxis: { title: 'Time to Expiry (years)' },
                    zaxis: { title: 'Implied Volatility (%)' },
                    camera: {
                        eye: { x: 1.5, y: 1.5, z: 1.5 }
                    }
                },
                margin: { l: 0, r: 0, b: 0, t: 0 },
                showlegend: true
            };
            
            Plotly.newPlot('surface3d', traces, layout);
        }
        
        function createVolStrikeChart(data) {
            const points = data.points;
            const expiryGroups = {};
            
            points.forEach(p => {
                const expiry = p.timeToExpiry.toFixed(3);
                if (!expiryGroups[expiry]) expiryGroups[expiry] = [];
                expiryGroups[expiry].push(p);
            });
            
            const traces = Object.entries(expiryGroups).map(([expiry, pts]) => ({
                x: pts.map(p => p.strike),
                y: pts.map(p => p.iv * 100),
                mode: 'markers+lines',
                name: expiry + 'Y',
                type: 'scatter'
            }));
            
            const layout = {
                xaxis: { title: 'Strike Price ($)' },
                yaxis: { title: 'Implied Volatility (%)' },
                showlegend: true,
                margin: { l: 50, r: 50, b: 50, t: 50 }
            };
            
            Plotly.newPlot('volStrike', traces, layout);
        }
        
        function createTermStructure(data) {
            const summary = data.summary;
            const atmVols = summary.atmVolsByExpiry;
            
            const x = Object.keys(atmVols).map(k => parseFloat(k)).sort((a,b) => a-b);
            const y = x.map(time => atmVols[time.toFixed(4)] * 100);
            
            const trace = {
                x: x,
                y: y,
                mode: 'markers+lines',
                type: 'scatter',
                name: 'ATM Volatility',
                line: { width: 3 },
                marker: { size: 8 }
            };
            
            const layout = {
                xaxis: { title: 'Time to Expiry (years)' },
                yaxis: { title: 'ATM Implied Volatility (%)' },
                showlegend: false,
                margin: { l: 50, r: 50, b: 50, t: 50 }
            };
            
            Plotly.newPlot('termStructure', [trace], layout);
        }
        
        function createScatterPoints(data) {
            const points = data.points;
            const calls = points.filter(p => p.optionType === 'C');
            const puts = points.filter(p => p.optionType === 'P');
            
            const traces = [
                {
                    x: calls.map(p => p.moneyness),
                    y: calls.map(p => p.iv * 100),
                    mode: 'markers',
                    type: 'scatter',
                    name: 'Calls',
                    marker: { color: 'green', size: 6 }
                },
                {
                    x: puts.map(p => p.moneyness),
                    y: puts.map(p => p.iv * 100),
                    mode: 'markers',
                    type: 'scatter',
                    name: 'Puts',
                    marker: { color: 'red', size: 6 }
                }
            ];
            
            const layout = {
                xaxis: { title: 'Moneyness (Spot/Strike)' },
                yaxis: { title: 'Implied Volatility (%)' },
                showlegend: true,
                margin: { l: 50, r: 50, b: 50, t: 50 }
            };
            
            Plotly.newPlot('scatterPoints', traces, layout);
        }
        
        // Load data when page loads
        loadData();
    </script>
</body>
</html>`

	tmpl, err := template.New("home").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, nil)
}

// handleSurfaceData serves the volatility surface data as JSON
func (ws *WebServer) handleSurfaceData(w http.ResponseWriter, r *http.Request) {
	if ws.model == nil {
		http.Error(w, "No data available", http.StatusNotFound)
		return
	}

	// Clean data to remove NaN values
	cleanSummary := ws.cleanSummaryForJSON(ws.model.Summary)
	
	data := VolSurfaceData{
		Points:    ws.model.Points,
		Summary:   cleanSummary,
		SpotPrice: ws.analyzer.spotPrice,
	}

	log.Printf("Preparing to encode %d points, spot price: %.2f", len(data.Points), data.SpotPrice)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
	log.Printf("Successfully encoded surface data")
}

// handlePointsData serves individual point data
func (ws *WebServer) handlePointsData(w http.ResponseWriter, r *http.Request) {
	if ws.model == nil {
		http.Error(w, "No data available", http.StatusNotFound)
		return
	}

	// Optional filtering by expiry
	expiryParam := r.URL.Query().Get("expiry")
	points := ws.model.Points

	if expiryParam != "" {
		if expiry, err := strconv.ParseFloat(expiryParam, 64); err == nil {
			filteredPoints := []VolSurfacePoint{}
			for _, point := range points {
				if abs(point.TimeToExpiry-expiry) < 0.01 { // tolerance for floating point comparison
					filteredPoints = append(filteredPoints, point)
				}
			}
			points = filteredPoints
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(points)
}

// abs returns absolute value of float64
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// cleanSummaryForJSON removes NaN/Inf values from summary to make it JSON serializable
func (ws *WebServer) cleanSummaryForJSON(summary VolSurfaceSummary) VolSurfaceSummary {
	cleaned := summary
	
	// Replace NaN/Inf values with proper defaults
	if math.IsNaN(cleaned.AvgIV) || math.IsInf(cleaned.AvgIV, 0) {
		cleaned.AvgIV = 0
	}
	if math.IsNaN(cleaned.MinIV) || math.IsInf(cleaned.MinIV, 0) {
		cleaned.MinIV = 0
	}
	if math.IsNaN(cleaned.MaxIV) || math.IsInf(cleaned.MaxIV, 0) {
		cleaned.MaxIV = 0
	}
	if math.IsNaN(cleaned.AvgTimeToExpiry) || math.IsInf(cleaned.AvgTimeToExpiry, 0) {
		cleaned.AvgTimeToExpiry = 0
	}
	if math.IsNaN(cleaned.MinTimeToExpiry) || math.IsInf(cleaned.MinTimeToExpiry, 0) {
		cleaned.MinTimeToExpiry = 0
	}
	if math.IsNaN(cleaned.MaxTimeToExpiry) || math.IsInf(cleaned.MaxTimeToExpiry, 0) {
		cleaned.MaxTimeToExpiry = 0
	}
	
	// Clean maps - ensure they're initialized
	if cleaned.ATMVolsByExpiry == nil {
		cleaned.ATMVolsByExpiry = make(map[string]float64)
	}
	if cleaned.VolSkewByExpiry == nil {
		cleaned.VolSkewByExpiry = make(map[string]float64)
	}
	
	cleanATMVols := make(map[string]float64)
	for k, v := range summary.ATMVolsByExpiry {
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			cleanATMVols[k] = v
		}
	}
	cleaned.ATMVolsByExpiry = cleanATMVols
	
	cleanSkews := make(map[string]float64)
	for k, v := range summary.VolSkewByExpiry {
		if !math.IsNaN(v) && !math.IsInf(v, 0) {
			cleanSkews[k] = v
		}
	}
	cleaned.VolSkewByExpiry = cleanSkews
	
	log.Printf("Cleaned summary: TotalPoints=%d, AvgIV=%.4f, CallPoints=%d, PutPoints=%d", 
		cleaned.TotalPoints, cleaned.AvgIV, cleaned.CallPoints, cleaned.PutPoints)
	
	return cleaned
}

// FittedSurfaceData represents fitted surface data for visualization
type FittedSurfaceData struct {
	StrikeRange     []float64   `json:"strikeRange"`
	TimeRange       []float64   `json:"timeRange"`
	GridSize        int         `json:"gridSize"`
	IVGrid          [][]float64 `json:"ivGrid"`
	RawPoints       []VolSurfacePoint `json:"rawPoints"`
	FilteredPoints  []VolSurfacePoint `json:"filteredPoints"`
	SpotPrice       float64     `json:"spotPrice"`
	HasSVIFit       bool        `json:"hasSVIFit"`
}

// handleFittedSurface serves the fitted SVI surface data
func (ws *WebServer) handleFittedSurface(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving fitted surface data...")
	
	// Generate grid for fitted surface
	gridSize := 50
	minStrike := ws.analyzer.spotPrice * 0.5
	maxStrike := ws.analyzer.spotPrice * 2.0
	minTime := 0.01
	maxTime := 1.0
	
	data := FittedSurfaceData{
		StrikeRange: []float64{minStrike, maxStrike},
		TimeRange:   []float64{minTime, maxTime},
		GridSize:    gridSize,
		SpotPrice:   ws.analyzer.spotPrice,
		HasSVIFit:   ws.analyzer.surface != nil,
	}
	
	// Generate strike and time grids
	strikes := make([]float64, gridSize)
	times := make([]float64, gridSize)
	
	for i := 0; i < gridSize; i++ {
		strikes[i] = minStrike + float64(i)*(maxStrike-minStrike)/float64(gridSize-1)
		times[i] = minTime + float64(i)*(maxTime-minTime)/float64(gridSize-1)
	}
	
	// Generate IV grid
	data.IVGrid = make([][]float64, gridSize)
	for i := 0; i < gridSize; i++ {
		data.IVGrid[i] = make([]float64, gridSize)
		for j := 0; j < gridSize; j++ {
			if ws.analyzer.surface != nil {
				// Use fitted surface
				iv := ws.analyzer.GetFittedIV(strikes[j], times[i])
				data.IVGrid[i][j] = iv * 100 // Convert to percentage
			} else {
				// Fallback: simple interpolation
				data.IVGrid[i][j] = 50.0 // Default 50% vol
			}
		}
	}
	
	// Add raw data points for comparison (included in surface fitting)
	for _, opt := range ws.analyzer.options {
		if !math.IsNaN(opt.ImpliedVolatility) && opt.ImpliedVolatility > 0 && opt.TimeToExpiry > 0 {
			data.RawPoints = append(data.RawPoints, VolSurfacePoint{
				Strike:       opt.Strike,
				TimeToExpiry: opt.TimeToExpiry,
				IV:           opt.ImpliedVolatility * 100,
				OptionType:   opt.Type,
			})
		}
	}
	
	// Add filtered data points (excluded from surface fitting)
	for _, opt := range ws.filteredOptions {
		if !math.IsNaN(opt.ImpliedVolatility) && opt.ImpliedVolatility > 0 && opt.TimeToExpiry > 0 {
			data.FilteredPoints = append(data.FilteredPoints, VolSurfacePoint{
				Strike:       opt.Strike,
				TimeToExpiry: opt.TimeToExpiry,
				IV:           opt.ImpliedVolatility * 100,
				OptionType:   opt.Type,
			})
		}
	}
	
	// Clean NaN values for JSON
	data = cleanFittedSurfaceForJSON(data)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// cleanFittedSurfaceForJSON replaces NaN/Inf values
func cleanFittedSurfaceForJSON(data FittedSurfaceData) FittedSurfaceData {
	for i := range data.IVGrid {
		for j := range data.IVGrid[i] {
			if math.IsNaN(data.IVGrid[i][j]) || math.IsInf(data.IVGrid[i][j], 0) {
				data.IVGrid[i][j] = 0
			}
		}
	}
	return data
}
