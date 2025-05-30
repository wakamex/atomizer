package monitor

import "fmt"

type Exporter struct {
	vmURL string
}

func NewExporter(vmURL string) *Exporter {
	return &Exporter{vmURL: vmURL}
}

func (e *Exporter) Export(query, start, end, step, format, output string) error {
	// TODO: Implement data export from VictoriaMetrics
	fmt.Println("Export not yet implemented")
	fmt.Printf("Would export from: %s\n", e.vmURL)
	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Time range: %s to %s (step: %s)\n", start, end, step)
	fmt.Printf("Format: %s, Output: %s\n", format, output)
	return nil
}