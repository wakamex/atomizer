package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// HTTPServer provides REST API for manual operations
type HTTPServer struct {
	orchestrator *ArbitrageOrchestrator
	riskManager  *RiskManager
	port         int
	server       *http.Server
}

// NewHTTPServer creates a new HTTP server
func NewHTTPServer(orchestrator *ArbitrageOrchestrator, riskManager *RiskManager, port int) *HTTPServer {
	return &HTTPServer{
		orchestrator: orchestrator,
		riskManager:  riskManager,
		port:         port,
	}
}

// Start begins serving HTTP requests
func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()
	
	// Trade endpoints
	mux.HandleFunc("/api/trade", s.handleTrade)
	mux.HandleFunc("/api/trades", s.handleGetTrades)
	
	// Risk endpoints
	mux.HandleFunc("/api/risk", s.handleGetRisk)
	mux.HandleFunc("/api/positions", s.handleGetPositions)
	
	// Health check
	mux.HandleFunc("/health", s.handleHealth)
	
	// Metrics endpoint (Prometheus format)
	mux.HandleFunc("/metrics", s.handleMetrics)
	
	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting HTTP server on %s", addr)
	
	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.corsMiddleware(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	
	return s.server.ListenAndServe()
}

// CORS middleware
func (s *HTTPServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// handleTrade processes manual trade submissions
func (s *HTTPServer) handleTrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req ManualTradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}
	
	// Validate request
	if err := s.validateTradeRequest(&req); err != nil {
		http.Error(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
		return
	}
	
	// Submit trade
	trade, err := s.orchestrator.SubmitManualTrade(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to submit trade: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Return response
	resp := map[string]interface{}{
		"trade_id":   trade.ID,
		"status":     trade.Status,
		"instrument": trade.Instrument,
		"quantity":   trade.Quantity.String(),
		"price":      trade.Price.String(),
		"timestamp":  trade.Timestamp,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleGetTrades returns active trades
func (s *HTTPServer) handleGetTrades(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	trades := s.orchestrator.GetActiveTrades()
	
	// Convert to API response format
	resp := make([]map[string]interface{}, len(trades))
	for i, trade := range trades {
		resp[i] = map[string]interface{}{
			"trade_id":       trade.ID,
			"source":         trade.Source,
			"status":         trade.Status,
			"instrument":     trade.Instrument,
			"quantity":       trade.Quantity.String(),
			"price":          trade.Price.String(),
			"hedge_order_id": trade.HedgeOrderID,
			"hedge_exchange": trade.HedgeExchange,
			"timestamp":      trade.Timestamp,
		}
		
		if trade.Error != nil {
			resp[i]["error"] = trade.Error.Error()
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleGetRisk returns current risk metrics
func (s *HTTPServer) handleGetRisk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	metrics := s.riskManager.GetRiskMetrics()
	
	resp := map[string]interface{}{
		"total_delta":       metrics.TotalDelta.String(),
		"total_gamma":       metrics.TotalGamma.String(),
		"total_positions":   metrics.TotalPositions,
		"max_position_size": metrics.MaxPositionSize.String(),
		"updated_at":        metrics.UpdatedAt,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleGetPositions returns current positions
func (s *HTTPServer) handleGetPositions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	positions := s.riskManager.GetPositions()
	
	resp := make([]map[string]interface{}, 0, len(positions))
	for instrument, pos := range positions {
		resp = append(resp, map[string]interface{}{
			"instrument":   instrument,
			"quantity":     pos.Quantity.String(),
			"avg_price":    pos.AvgPrice.String(),
			"delta":        pos.Delta.String(),
			"gamma":        pos.Gamma.String(),
			"last_updated": pos.LastUpdated,
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleHealth returns health status
func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"status": "healthy",
		"time":   time.Now(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleMetrics returns Prometheus-formatted metrics
func (s *HTTPServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := s.riskManager.GetRiskMetrics()
	trades := s.orchestrator.GetActiveTrades()
	
	// Count trades by status
	tradesByStatus := make(map[TradeStatus]int)
	for _, trade := range trades {
		tradesByStatus[trade.Status]++
	}
	
	// Format as Prometheus metrics
	fmt.Fprintf(w, "# HELP arbitrage_total_delta Current total delta exposure\n")
	fmt.Fprintf(w, "# TYPE arbitrage_total_delta gauge\n")
	fmt.Fprintf(w, "arbitrage_total_delta %s\n", metrics.TotalDelta.String())
	
	fmt.Fprintf(w, "\n# HELP arbitrage_total_gamma Current total gamma exposure\n")
	fmt.Fprintf(w, "# TYPE arbitrage_total_gamma gauge\n")
	fmt.Fprintf(w, "arbitrage_total_gamma %s\n", metrics.TotalGamma.String())
	
	fmt.Fprintf(w, "\n# HELP arbitrage_active_positions Number of active positions\n")
	fmt.Fprintf(w, "# TYPE arbitrage_active_positions gauge\n")
	fmt.Fprintf(w, "arbitrage_active_positions %d\n", metrics.TotalPositions)
	
	fmt.Fprintf(w, "\n# HELP arbitrage_trades_total Total number of trades by status\n")
	fmt.Fprintf(w, "# TYPE arbitrage_trades_total counter\n")
	for status, count := range tradesByStatus {
		fmt.Fprintf(w, "arbitrage_trades_total{status=\"%s\"} %d\n", status, count)
	}
	
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
}

// validateTradeRequest validates a manual trade request
func (s *HTTPServer) validateTradeRequest(req *ManualTradeRequest) error {
	if req.Instrument == "" {
		return fmt.Errorf("instrument is required")
	}
	
	if req.Side != "buy" && req.Side != "sell" {
		return fmt.Errorf("side must be 'buy' or 'sell'")
	}
	
	if req.Quantity.IsZero() || req.Quantity.IsNegative() {
		return fmt.Errorf("quantity must be positive")
	}
	
	if req.Price.IsZero() || req.Price.IsNegative() {
		return fmt.Errorf("price must be positive")
	}
	
	if req.Strike.IsZero() || req.Strike.IsNegative() {
		return fmt.Errorf("strike must be positive")
	}
	
	if req.Expiry <= 0 {
		return fmt.Errorf("expiry must be a valid timestamp")
	}
	
	// Check if expiry is in the future
	if time.Unix(req.Expiry, 0).Before(time.Now()) {
		return fmt.Errorf("expiry must be in the future")
	}
	
	return nil
}

// Stop gracefully stops the HTTP server
func (s *HTTPServer) Stop() error {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.server.Shutdown(ctx)
	}
	return nil
}

