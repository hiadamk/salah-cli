package config

import (
	"encoding/json"
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

	path, err := GetConfigPath()
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

	_, err := GetConfigPath()
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

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		expectErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				Latitude:           40.0,
				Longitude:          -70.0,
				EnableHighlighting: true,
				HighlightColour:    "green",
			},
			expectErr: false,
		},
		{
			name: "latitude out of range",
			cfg: Config{
				Latitude:  100.0,
				Longitude: 0.0,
			},
			expectErr: true,
		},
		{
			name: "longitude out of range",
			cfg: Config{
				Latitude:  0.0,
				Longitude: 200.0,
			},
			expectErr: true,
		},
		{
			name: "invalid highlight colour",
			cfg: Config{
				Latitude:           10.0,
				Longitude:          10.0,
				EnableHighlighting: true,
				HighlightColour:    "notacolour",
			},
			expectErr: true,
		},
		{
			name: "isha angle and interval conflict",
			cfg: Config{
				Latitude:     10.0,
				Longitude:    10.0,
				IshaAngle:    floatPtr(15.0),
				IshaInterval: intPtr(20),
			},
			expectErr: true,
		},
		{
			name: "highlighting disabled ignores colour",
			cfg: Config{
				Latitude:        10.0,
				Longitude:       10.0,
				HighlightColour: "notacolour", // should be fine since highlighting is off
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.expectErr && err == nil {
				t.Errorf("expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("expected no error but got %v", err)
			}
		})
	}
}

func floatPtr(v float64) *float64 {
	return &v
}

func intPtr(v int) *int {
	return &v
}

// -----------------------
// Tests for attemptAtomicRename
// -----------------------
func TestAttemptAtomicRename_Success(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "temp.txt")
	targetFile := filepath.Join(tmpDir, "target.txt")

	if err := os.WriteFile(tmpFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	if err := attemptAtomicRename(tmpFile, targetFile); err != nil {
		t.Fatalf("expected rename to succeed, got %v", err)
	}

	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		t.Fatalf("target file does not exist after rename")
	}

	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Fatalf("temp file should be removed after rename")
	}
}

func TestAttemptAtomicRename_TargetExists(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "temp.txt")
	targetFile := filepath.Join(tmpDir, "target.txt")

	if err := os.WriteFile(tmpFile, []byte("temp"), 0o644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if err := os.WriteFile(targetFile, []byte("existing"), 0o644); err != nil {
		t.Fatalf("failed to create target file: %v", err)
	}

	if err := attemptAtomicRename(tmpFile, targetFile); err != nil {
		t.Fatalf("expected rename with existing target to succeed, got %v", err)
	}

	data, _ := os.ReadFile(targetFile)
	if string(data) != "temp" {
		t.Fatalf("target file content not replaced correctly")
	}

	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Fatalf("temp file should be removed after rename")
	}
}

// -----------------------
// Tests for saveConfig
// -----------------------
func TestSaveConfig_CreatesFileAndDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &Config{Latitude: 51.5, Longitude: -0.1}

	targetFile := filepath.Join(tmpDir, "nested", "config.json")
	if err := SaveConfig(cfg, targetFile); err != nil {
		t.Fatalf("saveConfig failed: %v", err)
	}

	// check file exists
	if _, err := os.Stat(targetFile); os.IsNotExist(err) {
		t.Fatalf("config file not created")
	}

	// check contents
	data, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	var readCfg Config
	if err := json.Unmarshal(data, &readCfg); err != nil {
		t.Fatalf("failed to decode config JSON: %v", err)
	}

	if readCfg.Latitude != cfg.Latitude || readCfg.Longitude != cfg.Longitude {
		t.Fatalf("config contents do not match")
	}
}

func TestSaveConfig_AtomicOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	targetFile := filepath.Join(tmpDir, "config.json")

	cfg1 := &Config{Latitude: 1.1, Longitude: 2.2}
	if err := SaveConfig(cfg1, targetFile); err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	cfg2 := &Config{Latitude: 3.3, Longitude: 4.4}
	if err := SaveConfig(cfg2, targetFile); err != nil {
		t.Fatalf("overwrite save failed: %v", err)
	}

	data, _ := os.ReadFile(targetFile)
	var readCfg Config
	if err := json.Unmarshal(data, &readCfg); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if readCfg.Latitude != cfg2.Latitude || readCfg.Longitude != cfg2.Longitude {
		t.Fatalf("overwrite did not update file correctly")
	}
}
