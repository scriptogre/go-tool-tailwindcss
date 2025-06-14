package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	binary, err := getTailwindBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Run TailwindCSS with all args
	cmd := exec.Command(binary, os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		os.Exit(1)
	}
}

func getTailwindBinary() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, ".cache", "go-tool-tailwindcss")
	binaryPath := filepath.Join(cacheDir, "tailwindcss")
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	// If cached, use it
	if _, err := os.Stat(binaryPath); err == nil {
		return binaryPath, nil
	}

	// Download once on first use
	fmt.Fprintln(os.Stderr, "Downloading TailwindCSS...")
	return downloadAndCache(binaryPath)
}

func downloadAndCache(binaryPath string) (string, error) {
	// Create cache directory
	if err := os.MkdirAll(filepath.Dir(binaryPath), 0755); err != nil {
		return "", fmt.Errorf("could not create cache directory: %w", err)
	}

	// Get latest version
	version, err := getLatestVersion()
	if err != nil {
		return "", fmt.Errorf("could not get latest version: %w", err)
	}

	// Get download URL for current platform
	filename := getTailwindFilename()
	url := fmt.Sprintf("https://github.com/tailwindlabs/tailwindcss/releases/download/%s/%s", version, filename)

	// Download
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("could not download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}

	// Write to cache
	file, err := os.Create(binaryPath)
	if err != nil {
		return "", fmt.Errorf("could not create cache file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("could not write cache file: %w", err)
	}

	// Make executable
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return "", fmt.Errorf("could not make binary executable: %w", err)
	}

	return binaryPath, nil
}

func getTailwindFilename() string {
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return "tailwindcss-linux-x64"
		case "arm64":
			return "tailwindcss-linux-arm64"
		default:
			panic(fmt.Sprintf("unsupported architecture: %s", runtime.GOARCH))
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			return "tailwindcss-macos-x64"
		case "arm64":
			return "tailwindcss-macos-arm64"
		default:
			panic(fmt.Sprintf("unsupported architecture: %s", runtime.GOARCH))
		}
	case "windows":
		if runtime.GOARCH == "amd64" {
			return "tailwindcss-windows-x64.exe"
		}
		panic(fmt.Sprintf("unsupported architecture: %s", runtime.GOARCH))
	default:
		panic(fmt.Sprintf("unsupported platform: %s", runtime.GOOS))
	}
}

func getLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/tailwindlabs/tailwindcss/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}
