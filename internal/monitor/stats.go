package monitor

import "fmt"

type StatsClient struct {
	vmURL string
}

func NewStatsClient(vmURL string) *StatsClient {
	return &StatsClient{vmURL: vmURL}
}

func (s *StatsClient) ShowStats() error {
	// TODO: Implement stats queries to VictoriaMetrics
	fmt.Println("Stats not yet implemented")
	fmt.Printf("Would query stats from: %s\n", s.vmURL)
	return nil
}
