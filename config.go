package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	calc "github.com/mnadev/adhango/pkg/calc"
)

// Config holds user settings
type Config struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	Method            *int                    `json:"method,omitempty"`
	FajrAngle         *float64                `json:"fajr_angle,omitempty"`
	IshaAngle         *float64                `json:"isha_angle,omitempty"`
	IshaInterval      *int                    `json:"isha_interval,omitempty"`
	Madhab            *int                    `json:"madhab,omitempty"`
	HighLatitudeRule  *int                    `json:"high_latitude_rule,omitempty"`
	Adjustments       *calc.PrayerAdjustments `json:"adjustments,omitempty"`
	MethodAdjustments *calc.PrayerAdjustments `json:"method_adjustments,omitempty"`

	// User Preferences
	EnableCountdown    bool   `json:"enable_countdown"`
	EnableHighlighting bool   `json:"enable_highlighting"`
	HighlightColour    string `json:"highlight_colour"`
}

const (
	DefaultConfigFileName = "config.json"
	AppName               = "salah-cli"
	unixDefaultConfigDir  = ".config"
)

// For testability, allow overriding environment variable lookup
var getEnv = os.Getenv

// For testability, allow overriding the OS
var getOS = func() string { return runtime.GOOS }

// getConfigPath determines the default path for the config file
func getConfigPath() (string, error) {
	switch getOS() {
	case "windows":
		appData := getEnv("APPDATA")
		if appData == "" {
			userProfile := getEnv("USERPROFILE")
			if userProfile == "" {
				return "", fmt.Errorf("APPDATA and USERPROFILE not set")
			}
			appData = filepath.Join(userProfile, "AppData", "Roaming")
		}
		return filepath.Join(appData, AppName, DefaultConfigFileName), nil
	case "darwin", "linux":
		configHome := getEnv("XDG_CONFIG_HOME")
		if configHome == "" {
			home := getEnv("HOME")
			if home == "" {
				return "", fmt.Errorf("HOME not set")
			}
			configHome = filepath.Join(home, unixDefaultConfigDir)
		}
		return filepath.Join(configHome, AppName, DefaultConfigFileName), nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", getOS())
	}
}

// loadFromFile loads config from a given file path with validation
func loadFromFile(path string) (*Config, error) {
	// Ensure the config directory exists
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory %s: %w", dir, err)
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file %s: %w", path, err)
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields() // fail if unexpected keys are found

	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error decoding JSON from %s: %w", path, err)
	}

	// Run validation checks
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config in %s: %w", path, err)
	}

	return &cfg, nil
}

// load loads config from the default config path
func load() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	return loadFromFile(path)
}

// Validate checks for semantic errors in the configuration
func (c *Config) Validate() error {
	// Latitude must be -90..90
	if c.Latitude < -90 || c.Latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90 (got %f)", c.Latitude)
	}

	// Longitude must be -180..180
	if c.Longitude < -180 || c.Longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180 (got %f)", c.Longitude)
	}

	// Highlight colour must be valid if provided
	if c.EnableHighlighting && c.HighlightColour != "" {
		if _, ok := ansiColors[c.HighlightColour]; !ok {
			return fmt.Errorf("invalid highlight colour '%s'. Allowed: %v", c.HighlightColour, keys(ansiColors))
		}
	}

	// If both isha_angle and isha_interval are set, that’s a conflict
	if c.IshaAngle != nil && c.IshaInterval != nil {
		return fmt.Errorf("only one of isha_angle or isha_interval can be set")
	}

	return nil
}

// keys returns the keys of a string map (helper for error messages)
func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
