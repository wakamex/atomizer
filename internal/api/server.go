package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wakamex/atomizer/internal/types"
)

// Server provides REST API for manual operations
type Server struct {
	orchestrator Orchestrator
	riskManager  types.RiskManager
	port         int
	server       *http.Server
}

// Orchestrator interface for the arbitrage orchestrator
type Orchestrator interface {
	SubmitManualTrade(req types.ManualTradeRequest) (*types.TradeEvent, error)
	GetActiveTrades() []types.TradeEvent
}

// NewServer creates a new HTTP server
func NewServer(orchestrator Orchestrator, riskManager types.RiskManager, port int) *Server {
	return &Server{
		orchestrator: orchestrator,
		riskManager:  riskManager,
		port:         port,
	}
}

// Start begins serving HTTP requests
func (s *Server) Start() error {
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

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return s.server.Shutdown(ctx)
	}
	return nil
}

// CORS middleware
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
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
func (s *Server) handleTrade(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var req types.ManualTradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Submit trade to orchestrator
	trade, err := s.orchestrator.SubmitManualTrade(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to submit trade: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Return trade details
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trade)
}

// handleGetTrades returns active trades
func (s *Server) handleGetTrades(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	trades := s.orchestrator.GetActiveTrades()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trades)
}

// handleGetRisk returns current risk metrics
func (s *Server) handleGetRisk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	delta, gamma := s.riskManager.GetGreeks()
	
	response := map[string]interface{}{
		"delta":     delta.String(),
		"gamma":     gamma.String(),
		"timestamp": time.Now().Unix(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleGetPositions returns current positions
func (s *Server) handleGetPositions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	positions := s.riskManager.GetPositions()
	
	// Convert positions map to slice for JSON response
	positionList := make([]PositionResponse, 0, len(positions))
	for instrument, pos := range positions {
		positionList = append(positionList, PositionResponse{
			Instrument:  instrument,
			Quantity:    pos.Quantity.String(),
			AvgPrice:    pos.AvgPrice.String(),
			Delta:       pos.Delta.String(),
			Gamma:       pos.Gamma.String(),
			LastUpdated: pos.LastUpdated.Unix(),
		})
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(positionList)
}

// handleHealth returns server health status
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleMetrics returns Prometheus-style metrics
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	delta, gamma := s.riskManager.GetGreeks()
	positions := s.riskManager.GetPositions()
	trades := s.orchestrator.GetActiveTrades()
	
	// Format metrics in Prometheus format
	fmt.Fprintf(w, "# HELP portfolio_delta Current portfolio delta\n")
	fmt.Fprintf(w, "# TYPE portfolio_delta gauge\n")
	fmt.Fprintf(w, "portfolio_delta %s\n", delta.StringFixed(4))
	
	fmt.Fprintf(w, "# HELP portfolio_gamma Current portfolio gamma\n")
	fmt.Fprintf(w, "# TYPE portfolio_gamma gauge\n")
	fmt.Fprintf(w, "portfolio_gamma %s\n", gamma.StringFixed(4))
	
	fmt.Fprintf(w, "# HELP active_positions Number of active positions\n")
	fmt.Fprintf(w, "# TYPE active_positions gauge\n")
	fmt.Fprintf(w, "active_positions %d\n", len(positions))
	
	fmt.Fprintf(w, "# HELP active_trades Number of active trades\n")
	fmt.Fprintf(w, "# TYPE active_trades gauge\n")
	fmt.Fprintf(w, "active_trades %d\n", len(trades))
}

// PositionResponse represents a position in the API response
type PositionResponse struct {
	Instrument  string `json:"instrument"`
	Quantity    string `json:"quantity"`
	AvgPrice    string `json:"avg_price"`
	Delta       string `json:"delta"`
	Gamma       string `json:"gamma"`
	LastUpdated int64  `json:"last_updated"`
}