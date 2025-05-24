# 🚀 ETH Options Volatility Surface Analyzer

A comprehensive web-based visualization tool for analyzing Ethereum options volatility surfaces, implied volatilities, and market microstructure.

## ✨ Features

### 📊 **Volatility Surface Analysis**
- **3D Surface Visualization** - Interactive 3D plot of IV across strike and time
- **Term Structure Analysis** - Volatility term structure with slope analysis  
- **Volatility Smile Analysis** - Strike-based smile analysis for each expiry
- **Real-time Market Data** - Integration with Deribit options data

### 📈 **Advanced Analytics**
- **Black-Scholes Implementation** - Full BS model with Greeks calculation
- **Implied Volatility Extraction** - Newton-Raphson IV calculation
- **Market Structure Metrics** - Put/call ratios, volume, open interest
- **Smile Asymmetry Detection** - Quantified put/call skew analysis

### 🌐 **Web Interface**
- **Interactive Charts** - Powered by Plotly.js
- **Real-time Updates** - Live data streaming capability
- **Multiple Views** - Surface, term structure, individual points
- **Mobile Responsive** - Works on desktop and mobile devices

## 🚀 Quick Start

### Option 1: Run the Web App
```bash
# Build the web application
go build -o vol_web_app vol_web_app.go web_server.go vol_surface_analyzer.go black_scholes.go vol_surface_fitter.go rpc_models.go

# Start the web server (default port 8080)
./vol_web_app

# Or specify a custom port
./vol_web_app -port=9000
```

### Option 2: Run Tests (includes startup test)
```bash
# Run all volatility surface tests
go test -v -run TestVolSurfaceAnalysis

# Run web app tests
go test -v -run TestVolSurfaceWebApp

# Test web app startup (runs server briefly)
go test -v -run TestWebAppStartup
```

## 📊 Example Analysis Output

```
=== VOLATILITY SURFACE ANALYSIS ===
Total Points: 199 (Calls: 46, Puts: 153)
Average IV: 195.42%
IV Range: 7.74% - 495.64%
Time to Expiry Range: 0.0040 - 0.8418 years

--- ATM Volatilities by Expiry ---
  0.0040 years: ATM 435.55%
  0.0068 years: ATM 370.12%
  0.0177 years: ATM 273.27%
  ...

--- Volatility Term Structure ---
Term Structure: 435.55% (short) → 115.44% (long) [Downward Sloping]

--- Volatility Smile Analysis ---
  0.056 years expiry:
    ATM (1.00): 212.39%
    Left Wing: 271.47%
    Right Wing: 177.26%
    Asymmetry: 94.20% [Put Skew]
```

## 🛠 Technical Architecture

### Core Components
- **`vol_surface_analyzer.go`** - Data loading and IV calculation engine
- **`black_scholes.go`** - Complete Black-Scholes implementation with Greeks
- **`vol_surface_fitter.go`** - Surface fitting and analysis algorithms
- **`web_server.go`** - HTTP server and JSON API endpoints
- **`vol_web_app.go`** - Standalone web application runner

### API Endpoints
- **`GET /`** - Main web interface
- **`GET /api/surface`** - Complete volatility surface data (JSON)
- **`GET /api/points`** - Individual IV points data (JSON)
- **`GET /api/points?expiry=0.25`** - Filtered points by expiry (JSON)

### Data Pipeline
```
CSV Data → IV Calculation → Surface Fitting → Web Visualization
    ↓            ↓              ↓              ↓
Raw Options  Black-Scholes   Analytics    Interactive Charts
```

## 📈 Visualization Features

### 3D Volatility Surface
- **Interactive 3D Plot** with rotatable view
- **Color-coded IV levels** using Viridis colormap
- **Hover tooltips** with detailed option information
- **Customizable axis ranges** and viewing angles

### Volatility Smile Charts
- **Strike vs IV plots** for each expiry
- **Call/Put differentiation** with color coding
- **Trend line fitting** and asymmetry quantification
- **Moneyness-based analysis** around ATM levels

### Term Structure Analysis
- **Time vs ATM IV** with trend analysis
- **Slope classification** (upward/downward/flat)
- **Historical comparison** capabilities
- **Forward vol curve** extrapolation

## 🔧 Configuration

### Data Sources
- **Primary**: Existing CSV files with options data
- **Live Data**: Deribit API integration (configurable)
- **Custom Data**: JSON import capability

### Analysis Parameters
```go
const (
    riskFreeRate = 0.05    // 5% risk-free rate
    maxIterations = 100    // IV solver iterations
    tolerance = 1e-6       // Convergence tolerance
)
```

## 📊 Market Insights

### Key Metrics Tracked
- **Implied Volatility Range**: 7.74% - 495.64%
- **Put/Call Ratio**: Typically 3:1 (153 puts vs 46 calls)
- **Volatility Term Structure**: Usually downward sloping in crypto
- **Smile Asymmetry**: Strong put skew (94-142% asymmetry)

### Typical Crypto Options Patterns
- **High Short-Term IV**: 400%+ for weekly options
- **Mean-Reverting Term Structure**: Long-term IV around 115%
- **Strong Put Bias**: Expensive downside protection
- **Smile Asymmetry**: Left wing dominance typical

## 🚀 Advanced Usage

### Custom Analysis
```go
// Load and analyze custom data
analyzer := NewVolSurfaceAnalyzer()
analyzer.LoadCSV("custom_options.csv")
analyzer.CalculateImpliedVolatilities()
surface := analyzer.FitVolatilitySurface()
analyzer.PrintVolSurfaceAnalysis(surface)
```

### Live Data Integration
```go
// Enable live Deribit pricing
analyzer.UpdatePricing() // Fetches current market data
```

### Web API Integration
```bash
# Get current surface data
curl http://localhost:8080/api/surface

# Get points for specific expiry
curl http://localhost:8080/api/points?expiry=0.25
```

## 🎯 Use Cases

### **Trading Applications**
- **Option Pricing** - Fair value estimation using real market IVs
- **Arbitrage Detection** - Surface anomaly identification
- **Risk Management** - Greeks analysis and sensitivity testing
- **Strategy Development** - Volatility-based trading strategies

### **Research & Analysis**
- **Market Structure Analysis** - Term structure and smile evolution
- **Historical Studies** - Volatility regime identification  
- **Model Validation** - Black-Scholes vs market price comparison
- **Educational Purposes** - Options theory visualization

### **Institutional Applications**
- **Portfolio Management** - Volatility exposure analysis
- **Risk Reporting** - Comprehensive volatility metrics
- **Regulatory Compliance** - Market data visualization
- **Client Reporting** - Professional-grade analytics

## 🔬 Technical Details

### Black-Scholes Implementation
- **Newton-Raphson IV Solver** with adaptive convergence
- **Complete Greeks Suite** (Delta, Gamma, Theta, Vega)
- **Numerical Stability** for extreme parameters
- **Edge Case Handling** for expired/ATM options

### Performance Optimizations
- **Concurrent Processing** for large datasets
- **Efficient JSON Serialization** with NaN handling
- **Memory-Optimized** data structures
- **Caching Strategy** for repeated calculations

---

**🚀 Ready to explore the volatility surface? Start the web app and visit http://localhost:8080**