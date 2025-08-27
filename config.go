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

// loadFromFile loads config from a given file path
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
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error decoding JSON from %s: %w", path, err)
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
