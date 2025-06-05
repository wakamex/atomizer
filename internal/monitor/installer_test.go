package monitor

import (
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestGetLatestReleaseURL(t *testing.T) {
	installer := NewVMInstaller()

	url, err := installer.GetLatestReleaseURL()
	if err != nil {
		t.Fatalf("Failed to get latest release URL: %v", err)
	}

	// Check that we got a valid URL
	if !strings.HasPrefix(url, "https://github.com/VictoriaMetrics/VictoriaMetrics/releases/download/") {
		t.Errorf("Invalid URL format: %s", url)
	}

	// Check that it contains victoria-metrics and the architecture
	if !strings.Contains(url, "victoria-metrics") || !strings.Contains(url, runtime.GOARCH) {
		t.Errorf("URL doesn't contain expected pattern: %s", url)
	}

	// Check that it ends with .tar.gz
	if !strings.HasSuffix(url, ".tar.gz") {
		t.Errorf("URL doesn't end with .tar.gz: %s", url)
	}

	t.Logf("Successfully got latest release URL: %s", url)
}

func TestVMInstallerPaths(t *testing.T) {
	installer := NewVMInstaller()

	binaryPath := installer.GetBinaryPath()
	if !strings.Contains(binaryPath, "victoria-metrics-prod") {
		t.Errorf("Binary path doesn't contain expected name: %s", binaryPath)
	}

	dataPath := installer.GetDataPath()
	if !strings.Contains(dataPath, "data") {
		t.Errorf("Data path doesn't contain 'data': %s", dataPath)
	}
}

// Integration test - only run with -integration flag
func TestVMInstallerSetup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	

	installer := NewVMInstaller()

	// This actually downloads VictoriaMetrics
	err := installer.Setup()
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Verify the binary exists
	binaryPath := installer.GetBinaryPath()
	if _, err := os.Stat(binaryPath); err != nil {
		t.Errorf("Binary not found after setup: %s", binaryPath)
	}
}
