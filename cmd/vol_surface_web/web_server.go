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
	analyzer *VolSurfaceAnalyzer
	model    *VolSurfaceModel
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
	
	// Fit volatility surface
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
                
                document.getElementById('loading').style.display = 'none';
                document.getElementById('content').style.display = 'block';
                
                createStats(data);
                create3DSurface(data);
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
        
        function create3DSurface(data) {
            const points = data.points || [];
            console.log('Creating 3D surface with', points.length, 'points');
            
            if (points.length === 0) {
                document.getElementById('surface3d').innerHTML = '<p>No data available for 3D surface</p>';
                return;
            }
            
            // Prepare data for 3D surface
            const strikes = [...new Set(points.map(p => p.strike))].sort((a,b) => a-b);
            const times = [...new Set(points.map(p => p.timeToExpiry))].sort((a,b) => a-b);
            
            // Create Z matrix
            const z = times.map(time => {
                return strikes.map(strike => {
                    const point = points.find(p => 
                        Math.abs(p.strike - strike) < 0.1 && 
                        Math.abs(p.timeToExpiry - time) < 0.001
                    );
                    return point ? point.iv * 100 : null;
                });
            });
            
            const trace = {
                x: strikes,
                y: times,
                z: z,
                type: 'surface',
                colorscale: 'Viridis',
                showscale: true,
                colorbar: { title: 'IV (%)' }
            };
            
            const layout = {
                scene: {
                    xaxis: { title: 'Strike Price ($)' },
                    yaxis: { title: 'Time to Expiry (years)' },
                    zaxis: { title: 'Implied Volatility (%)' }
                },
                margin: { l: 0, r: 0, b: 0, t: 0 }
            };
            
            Plotly.newPlot('surface3d', [trace], layout);
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