package monitor

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

type VMInstaller struct {
	installPath string
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func NewVMInstaller() *VMInstaller {
	homeDir, _ := os.UserHomeDir()
	return &VMInstaller{
		installPath: filepath.Join(homeDir, ".atomizer", "victoria-metrics"),
	}
}

func (v *VMInstaller) GetLatestReleaseURL() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/VictoriaMetrics/VictoriaMetrics/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	// Determine the asset name based on OS and architecture
	// VictoriaMetrics uses the pattern: victoria-metrics-{os}-{arch}-{version}.tar.gz
	// We want the single-node version (not cluster, not enterprise)
	assetName := fmt.Sprintf("victoria-metrics-%s-%s-%s.tar.gz", 
		runtime.GOOS, runtime.GOARCH, release.TagName)

	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("no suitable release found for: %s", assetName)
}

func (v *VMInstaller) Setup() error {
	if err := os.MkdirAll(v.installPath, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	binaryPath := filepath.Join(v.installPath, "victoria-metrics-prod")
	if _, err := os.Stat(binaryPath); err == nil {
		fmt.Println("VictoriaMetrics is already installed at:", binaryPath)
		return nil
	}

	fmt.Println("Fetching latest VictoriaMetrics release...")
	downloadURL, err := v.GetLatestReleaseURL()
	if err != nil {
		return fmt.Errorf("failed to get latest release: %w", err)
	}

	fmt.Printf("Downloading VictoriaMetrics from %s...\n", downloadURL)
	if err := v.download(downloadURL); err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	fmt.Printf("VictoriaMetrics installed successfully at: %s\n", binaryPath)
	fmt.Printf("To start VictoriaMetrics manually:\n")
	fmt.Printf("  %s -storageDataPath=%s/data\n", binaryPath, v.installPath)
	return nil
}

func (v *VMInstaller) download(downloadURL string) error {
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	tmpFile := filepath.Join(v.installPath, "victoria-metrics.tar.gz")
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return err
	}

	return v.extract(tmpFile)
}

func (v *VMInstaller) extract(tarPath string) error {
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Name == "victoria-metrics-prod" {
			outPath := filepath.Join(v.installPath, header.Name)
			outFile, err := os.Create(outPath)
			if err != nil {
				return err
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				return err
			}

			if err := os.Chmod(outPath, 0755); err != nil {
				return err
			}
			break
		}
	}

	return nil
}

func (v *VMInstaller) GetBinaryPath() string {
	return filepath.Join(v.installPath, "victoria-metrics-prod")
}

func (v *VMInstaller) GetDataPath() string {
	return filepath.Join(v.installPath, "data")
}
