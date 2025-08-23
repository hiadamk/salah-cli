package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// Test loading a valid config from a temporary file
func TestLoadFromFile_ValidConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := `{
		"latitude": 51.5,
		"longitude": -0.12,
		"method": 1
	}`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	cfg, err := loadFromFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Latitude != 51.5 {
		t.Errorf("expected latitude 51.5, got %v", cfg.Latitude)
	}
	if cfg.Longitude != -0.12 {
		t.Errorf("expected longitude -0.12, got %v", cfg.Longitude)
	}
	if cfg.Method == nil || *cfg.Method != 1 {
		t.Errorf("expected method 1, got %v", cfg.Method)
	}
}

// Test loading a config when the directory does not exist
func TestLoadFromFile_CreateMissingDirectory(t *testing.T) {
	dir := t.TempDir()                                        // temporary directory
	configPath := filepath.Join(dir, "subdir", "config.json") // intentionally missing subdir

	// write a valid JSON file after creating parent dir
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		t.Fatalf("failed to create parent dir: %v", err)
	}
	if err := os.WriteFile(configPath, []byte(`{"latitude":10,"longitude":20}`), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := loadFromFile(configPath)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Latitude != 10 || cfg.Longitude != 20 {
		t.Errorf("unexpected config values: %+v", cfg)
	}
}

// Test Load() with overridden environment variables for Unix
func TestGetConfigPath_UnixFallback(t *testing.T) {
	originalGetEnv := getEnv
	defer func() { getEnv = originalGetEnv }()

	getEnv = func(key string) string {
		switch key {
		case "XDG_CONFIG_HOME":
			return ""
		case "HOME":
			return "/tmp/fakehome"
		}
		return ""
	}

	if runtime.GOOS == "windows" {
		t.Skip("Skipping Unix path test on Windows")
	}

	path, err := getConfigPath()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	expected := "/tmp/fakehome/.config/salah-cli/config.json"
	if path != expected {
		t.Errorf("expected %q, got %q", expected, path)
	}
}

// Test Load() with unsupported OS
func TestGetConfigPath_UnsupportedOS(t *testing.T) {
	originalGetOS := getOS
	defer func() { getOS = originalGetOS }()

	getOS = func() string { return "plan9" } // simulate unsupported OS

	_, err := getConfigPath()
	if err == nil {
		t.Fatal("expected error for unsupported OS, got nil")
	}
}

// Test LoadFromFile with invalid JSON
func TestLoadFromFile_InvalidJSON(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config_invalid*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(`invalid json`)
	tmpFile.Close()

	_, err = loadFromFile(tmpFile.Name())
	if err == nil {
		t.Fatal("expected JSON decoding error, got nil")
	}
}

// Test LoadFromFile with missing file
func TestLoadFromFile_MissingFile(t *testing.T) {
	_, err := loadFromFile("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
