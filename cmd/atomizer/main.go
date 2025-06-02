package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const version = "1.0.0"

var commands = map[string]struct {
	binary      string
	description string
}{
	"rfq":            {"maker_quote_responder", "Run RFQ responder for automated quoting"},
	"market-maker":   {"maker_quote_responder", "Run market maker with continuous quoting"},
	"manual-order":   {"maker_quote_responder", "Place a manual order on Derive exchange"},
	"analyze":        {"analyze_options", "Analyze options market liquidity and pricing"},
	"inventory":      {"inventory", "Show current positions and P&L"},
	"markets":        {"markets", "Display available markets and instruments"},
	"send-quote":     {"send_quote", "Send a single quote to an exchange"},
	"market-monitor": {"market_monitor", "Monitor and store real-time market data"},
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	// Handle version
	if os.Args[1] == "version" || os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Printf("atomizer version %s\n", version)
		os.Exit(0)
	}

	// Handle help
	if os.Args[1] == "help" || os.Args[1] == "--help" || os.Args[1] == "-h" {
		if len(os.Args) > 2 {
			showCommandHelp(os.Args[2])
		} else {
			showHelp()
		}
		os.Exit(0)
	}

	// Get command
	cmdName := os.Args[1]
	cmdInfo, exists := commands[cmdName]
	if !exists {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmdName)
		showHelp()
		os.Exit(1)
	}

	// Find binary path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding executable path: %v\n", err)
		os.Exit(1)
	}
	
	// Look for binary in same directory as atomizer
	binDir := filepath.Dir(execPath)
	binaryPath := filepath.Join(binDir, cmdInfo.binary)

	// Check if binary exists
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Binary not found: %s\n", binaryPath)
		fmt.Fprintf(os.Stderr, "Please run the build script first.\n")
		os.Exit(1)
	}

	// Special handling for market-maker and manual-order commands
	args := os.Args[2:]
	if cmdName == "market-maker" {
		// The maker_quote_responder binary expects "market-maker" as the first argument
		// followed by the actual market maker flags
		args = append([]string{"market-maker"}, args...)
	} else if cmdName == "manual-order" {
		// The maker_quote_responder binary expects "manual-order" as the first argument
		args = append([]string{"manual-order"}, args...)
	}

	// Execute the binary
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "Error running %s: %v\n", cmdInfo.binary, err)
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Atomizer - Unified toolkit for options trading")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  atomizer <command> [arguments]")
	fmt.Println()
	fmt.Println("Available commands:")
	
	// Print commands in a consistent order
	cmds := []string{"rfq", "market-maker", "manual-order", "analyze", "inventory", "markets", "send-quote", "market-monitor"}
	for _, cmd := range cmds {
		info := commands[cmd]
		fmt.Printf("  %-13s %s\n", cmd, info.description)
	}
	
	fmt.Println()
	fmt.Println("Global commands:")
	fmt.Printf("  %-13s Display this help message\n", "help")
	fmt.Printf("  %-13s Show version information\n", "version")
	fmt.Println()
	fmt.Println("Run 'atomizer help <command>' for more information on a command.")
}

func showCommandHelp(command string) {
	cmdInfo, exists := commands[command]
	if !exists {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	fmt.Printf("%s - %s\n\n", command, cmdInfo.description)
	
	// Command-specific help
	switch command {
	case "rfq":
		fmt.Println("Usage: atomizer rfq [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --derive-key KEY        Derive private key (or set DERIVE_PRIVATE_KEY)")
		fmt.Println("  --derive-wallet ADDR    Derive wallet address (or set DERIVE_WALLET_ADDRESS)")
		fmt.Println("  --deribit-key KEY       Deribit API key (or set DERIBIT_API_KEY)")
		fmt.Println("  --deribit-secret SECRET Deribit API secret (or set DERIBIT_API_SECRET)")
		fmt.Println("  -e, --exchange NAME     Exchange to use (derive or deribit)")
		fmt.Println("  -t, --test              Use testnet")
		
	case "market-maker":
		fmt.Println("Usage: atomizer market-maker [options]")
		fmt.Println()
		fmt.Println("Required:")
		fmt.Println("  --expiry DATE           Expiry date (e.g., 20250530)")
		fmt.Println("  --strikes STRIKES       Comma-separated strikes OR --all-strikes")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -e, --exchange NAME     Exchange to use (derive or deribit)")
		fmt.Println("  --underlying ASSET      Underlying asset (default: ETH)")
		fmt.Println("  --size SIZE             Quote size (default: 0.1)")
		fmt.Println("  --improvement AMOUNT    Price improvement (default: 0.1)")
		fmt.Println("  --improvement-reference-size SIZE")
		fmt.Println("                          Min size for best bid/ask selection (default: 0)")
		fmt.Println("  --max-position SIZE     Max position per instrument (default: 1.0)")
		fmt.Println("  --max-exposure SIZE     Max total exposure (default: 10.0)")
		fmt.Println("  --min-spread BPS        Min spread in basis points (default: 1000)")
		fmt.Println("  --refresh SECONDS       Refresh interval (default: 1)")
		fmt.Println("  --dry-run               Print config without starting")
		
	case "analyze":
		fmt.Println("Usage: atomizer analyze [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -e, --expiry INDEX      Expiry index (0=nearest, default: 0)")
		fmt.Println("  -u, --underlying ASSET  Underlying asset (default: ETH)")
		fmt.Println("  --exchanges LIST        Comma-separated exchanges")
		fmt.Println("  --compare               Compare across exchanges")
		
	case "inventory":
		fmt.Println("Usage: atomizer inventory [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -e, --exchange NAME     Exchange to query (or 'all')")
		fmt.Println("  -f, --format FORMAT     Output format (table, json, csv)")
		fmt.Println("  -r, --refresh           Auto-refresh every 5 seconds")
		
	case "markets":
		fmt.Println("Usage: atomizer markets [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -e, --exchange NAME     Exchange to query")
		fmt.Println("  -u, --underlying ASSET  Filter by underlying")
		fmt.Println("  --expiry DATE           Filter by expiry")
		fmt.Println("  --type TYPE             Filter by type (call, put)")
		fmt.Println("  --active                Show only active markets")
		
	case "send-quote":
		fmt.Println("Usage: atomizer send-quote [options]")
		fmt.Println()
		fmt.Println("Required:")
		fmt.Println("  -i, --instrument NAME   Instrument to quote")
		fmt.Println("  -s, --side SIDE         Side (buy or sell)")
		fmt.Println("  -p, --price PRICE       Limit price")
		fmt.Println("  --size SIZE             Order size")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  -e, --exchange NAME     Exchange to use")
		fmt.Println("  --type TYPE             Order type (limit, market)")
		
	case "manual-order":
		fmt.Println("Usage: atomizer manual-order")
		fmt.Println()
		fmt.Println("Places a manual order on Derive exchange.")
		fmt.Println()
		fmt.Println("Required Environment Variables:")
		fmt.Println("  DERIVE_PRIVATE_KEY      Your private key")
		fmt.Println("  DERIVE_WALLET_ADDRESS   Your Derive wallet address")
		fmt.Println()
		fmt.Println("Optional Environment Variables:")
		fmt.Println("  ORDER_INSTRUMENT        Instrument to trade (default: ETH-PERP)")
		fmt.Println("  ORDER_SIDE              Order side: buy/sell (default: buy)")
		fmt.Println("  ORDER_PRICE             Order price (default: 2000)")
		fmt.Println("  ORDER_AMOUNT            Order amount (default: 0.1)")
		fmt.Println("  DERIVE_SUBACCOUNT_ID    Subaccount ID (default: auto)")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  # Place default order (buy 0.1 ETH-PERP @ $2000)")
		fmt.Println("  atomizer manual-order")
		fmt.Println()
		fmt.Println("  # Place custom order")
		fmt.Println("  ORDER_INSTRUMENT=BTC-PERP ORDER_SIDE=sell ORDER_PRICE=45000 ORDER_AMOUNT=0.5 atomizer manual-order")
		
	case "market-monitor":
		fmt.Println("Usage: atomizer market-monitor [subcommand] [options]")
		fmt.Println()
		fmt.Println("Subcommands:")
		fmt.Println("  setup                   Download and configure VictoriaMetrics")
		fmt.Println("  start                   Start monitoring market data")
		fmt.Println("  stats                   Show current statistics")
		fmt.Println("  export                  Export data using PromQL")
		fmt.Println()
		fmt.Println("Start Options:")
		fmt.Println("  --interval DURATION     Collection interval (default: 5s)")
		fmt.Println("  --exchanges LIST        Comma-separated exchanges (default: derive,deribit)")
		fmt.Println("  --instruments PATTERN   Instrument patterns (default: ETH-*,BTC-*)")
		fmt.Println("  --vm-url URL            VictoriaMetrics URL (default: http://localhost:8428)")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  atomizer market-monitor setup")
		fmt.Println("  atomizer market-monitor start --interval 10s --instruments ETH-*")
	}
}