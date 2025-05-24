package main

import (
	"flag"
	"log"
)

func main() {
	port := flag.Int("port", 8080, "Port to run the web server on")
	flag.Parse()

	log.Println("🚀 ETH Options Volatility Surface Analyzer")
	log.Println("==========================================")

	// Create and initialize web server
	webServer := NewWebServer()
	
	// Load data and perform analysis
	log.Println("📊 Loading and analyzing volatility surface data...")
	if err := webServer.LoadAndAnalyze(); err != nil {
		log.Fatalf("❌ Failed to load and analyze data: %v", err)
	}
	
	// Start web server
	log.Printf("🌐 Starting web server on port %d", *port)
	log.Printf("✨ Open http://localhost:%d in your browser", *port)
	log.Println("📈 Features: 3D Surface, Term Structure, Volatility Smile, Live Analytics")
	log.Println("🔄 Press Ctrl+C to stop the server")
	
	if err := webServer.Start(*port); err != nil {
		log.Fatalf("❌ Failed to start web server: %v", err)
	}
}