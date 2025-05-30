package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/wakamex/atomizer/internal/monitor"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]

	switch subcommand {
	case "setup":
		runSetup()
	case "start":
		runStart(os.Args[2:])
	case "stats":
		runStats(os.Args[2:])
	case "export":
		runExport(os.Args[2:])
	default:
		fmt.Printf("Unknown subcommand: %s\n", subcommand)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: market_monitor [subcommand] [options]")
	fmt.Println()
	fmt.Println("Subcommands:")
	fmt.Println("  setup    - Download and configure VictoriaMetrics")
	fmt.Println("  start    - Start monitoring market data")
	fmt.Println("  stats    - Show current statistics")
	fmt.Println("  export   - Export data using PromQL")
}

func runSetup() {
	fmt.Println("Setting up VictoriaMetrics...")
	
	installer := monitor.NewVMInstaller()
	if err := installer.Setup(); err != nil {
		log.Fatalf("Setup failed: %v", err)
	}
	
	fmt.Println("Setup complete! You can now run:")
	fmt.Println("  1. Start VictoriaMetrics: ./victoria-metrics-prod -storageDataPath ./vm-data")
	fmt.Println("  2. Start monitoring: atomizer market-monitor start")
}

func runStart(args []string) {
	fs := flag.NewFlagSet("start", flag.ExitOnError)
	
	interval := fs.Duration("interval", 5*time.Second, "Collection interval")
	exchanges := fs.String("exchanges", "derive,deribit", "Comma-separated list of exchanges")
	instruments := fs.String("instruments", "ETH-PERP", "Comma-separated instrument patterns or exact names")
	vmURL := fs.String("vm-url", "http://localhost:8428", "VictoriaMetrics URL")
	workers := fs.Int("workers", 10, "Number of concurrent workers")
	orderbook := fs.Bool("orderbook", false, "Collect order book depth instead of just ticker")
	depth := fs.Int("depth", 10, "Order book depth to collect (for --orderbook mode)")
	debug := fs.Bool("debug", false, "Enable debug logging")
	
	if err := fs.Parse(args); err != nil {
		log.Fatal(err)
	}
	
	// Enable debug mode if requested
	monitor.SetDebug(*debug)
	
	// Parse exchanges
	exchangeList := strings.Split(*exchanges, ",")
	for i := range exchangeList {
		exchangeList[i] = strings.TrimSpace(exchangeList[i])
	}
	
	// Parse instruments
	instrumentPatterns := strings.Split(*instruments, ",")
	for i := range instrumentPatterns {
		instrumentPatterns[i] = strings.TrimSpace(instrumentPatterns[i])
	}
	
	config := &monitor.Config{
		Interval:           *interval,
		Exchanges:          exchangeList,
		InstrumentPatterns: instrumentPatterns,
		VictoriaMetricsURL: *vmURL,
		Workers:            *workers,
	}
	
	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		log.Println("Shutting down...")
	}()
	
	// Start monitoring
	log.Printf("Starting market monitor...")
	log.Printf("  Interval: %v", config.Interval)
	log.Printf("  Exchanges: %v", config.Exchanges)
	log.Printf("  Instruments: %v", config.InstrumentPatterns)
	log.Printf("  VictoriaMetrics: %s", config.VictoriaMetricsURL)
	
	// Create appropriate monitor based on mode
	if *orderbook {
		// Order book mode
		log.Printf("Mode: Order Book Collection (depth=%d)", *depth)
		
		obMon, err := monitor.NewOrderBookMonitor(config, *depth)
		if err != nil {
			log.Fatalf("Failed to create order book monitor: %v", err)
		}
		
		// Start order book monitoring
		if err := obMon.Start(); err != nil {
			log.Fatalf("Failed to start order book monitor: %v", err)
		}
		
		// Wait for interrupt
		<-sigChan
		
		// Stop monitor
		if err := obMon.Stop(); err != nil {
			log.Fatalf("Failed to stop order book monitor: %v", err)
		}
	} else {
		// Regular ticker mode
		log.Printf("Mode: Ticker Collection")
		
		mon, err := monitor.New(config)
		if err != nil {
			log.Fatalf("Failed to create monitor: %v", err)
		}
		
		// Start monitoring
		if err := mon.Start(); err != nil {
			log.Fatalf("Failed to start monitor: %v", err)
		}
		
		// Wait for interrupt
		<-sigChan
		
		// Stop monitor
		if err := mon.Stop(); err != nil {
			log.Fatalf("Failed to stop monitor: %v", err)
		}
	}
}

func runStats(args []string) {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	vmURL := fs.String("vm-url", "http://localhost:8428", "VictoriaMetrics URL")
	
	if err := fs.Parse(args); err != nil {
		log.Fatal(err)
	}
	
	stats := monitor.NewStatsClient(*vmURL)
	if err := stats.ShowStats(); err != nil {
		log.Fatalf("Failed to show stats: %v", err)
	}
}

func runExport(args []string) {
	fs := flag.NewFlagSet("export", flag.ExitOnError)
	
	vmURL := fs.String("vm-url", "http://localhost:8428", "VictoriaMetrics URL")
	query := fs.String("query", "", "PromQL query")
	start := fs.String("start", "1h", "Start time (e.g., 1h, 2d, 2024-01-01)")
	end := fs.String("end", "now", "End time")
	step := fs.String("step", "1m", "Step interval")
	format := fs.String("format", "csv", "Output format (csv, json)")
	output := fs.String("output", "", "Output file (default: stdout)")
	
	if err := fs.Parse(args); err != nil {
		log.Fatal(err)
	}
	
	if *query == "" {
		log.Fatal("Query is required")
	}
	
	exporter := monitor.NewExporter(*vmURL)
	if err := exporter.Export(*query, *start, *end, *step, *format, *output); err != nil {
		log.Fatalf("Export failed: %v", err)
	}
}