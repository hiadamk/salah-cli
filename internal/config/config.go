package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"salah-cli/internal/util"
	"strconv"

	"github.com/charmbracelet/huh"
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

// GetConfigPath determines the default path for the config file
func GetConfigPath() (string, error) {
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
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	return loadFromFile(path)
}

func validateLatitude(latitude float64) error {
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("latitude must be between -90 and 90 (got %f)", latitude)
	}
	return nil
}

func validateLongitude(longitude float64) error {
	// Longitude must be -180..180
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("longitude must be between -180 and 180 (got %f)", longitude)
	}
	return nil
}

// Validate checks for semantic errors in the configuration
func (c *Config) Validate() error {
	// Latitude must be -90..90
	if err := validateLatitude(c.Latitude); err != nil {
		return err
	}

	if err := validateLongitude(c.Longitude); err != nil {
		return err
	}

	// Highlight colour must be valid if provided
	if c.EnableHighlighting && c.HighlightColour != "" {
		if _, ok := util.AnsiColors[c.HighlightColour]; !ok {
			return fmt.Errorf("invalid highlight colour '%s'. Allowed: %v", c.HighlightColour, keys(util.AnsiColors))
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

func SetupConfig() (*Config, error) {
	var config Config

	var latitude string
	var longitude string
	var madhab int
	var moonsightingMethod int
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter your latitude:").
				Value(&latitude).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("value can't be empty")
					}
					fl, err := strconv.ParseFloat(str, 64)
					if err != nil {
						return fmt.Errorf("failed to parse latitude value: %v", err)
					}
					return validateLatitude(fl)

				}),

			huh.NewInput().
				Title("Enter your longitude:").
				Value(&longitude).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("value can't be empty")
					}
					fl, err := strconv.ParseFloat(str, 64)
					if err != nil {
						return fmt.Errorf("failed to parse longitude value: %v", err)
					}
					return validateLongitude(fl)

				}),

			huh.NewSelect[int]().
				Title("Choose your Madhab").
				Options(
					huh.NewOption("Shafi/Hanbali/Maliki", 0),
					huh.NewOption("Hanafi", 1),
				).
				Value(&madhab),

			huh.NewSelect[int]().
				Title("Choose your moonsighting method").
				Options(
					huh.NewOption("Other", 0),
					huh.NewOption("Muslim World League", 1),
					huh.NewOption("Egyptian", 2),
					huh.NewOption("Karachi", 3),
					huh.NewOption("Umm al-Qura", 4),
					huh.NewOption("Dubai", 5),
					huh.NewOption("Moon Sighting Committee", 6),
					huh.NewOption("North America (ISNA)", 7),
					huh.NewOption("Kuwait", 8),
					huh.NewOption("Qatar", 9),
					huh.NewOption("Singapore", 10),
					huh.NewOption("UOIF", 11),
				).Value(&moonsightingMethod),
		),
	)

	err := form.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to setup config: %s", err.Error())
	}

	// Ignoring error as this was previously validated
	latFloat, _ := strconv.ParseFloat(latitude, 64)
	config.Latitude = latFloat

	lonFloat, _ := strconv.ParseFloat(longitude, 64)
	config.Longitude = lonFloat
	config.Method = &moonsightingMethod
	config.Madhab = &madhab

	return &config, nil

}

// SaveConfig writes the given config to the specified filepath safely using an atomic rename
func SaveConfig(config *Config, path string) error {
	dir := filepath.Dir(path)
	if statErr := os.MkdirAll(dir, 0o755); statErr != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, statErr)
	}

	tmpFile, err := os.CreateTemp(dir, "config.json.tmp.*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// ensure temp file is removed on any early return
	cleanupTemp := func() { _ = os.Remove(tmpPath) }

	enc := json.NewEncoder(tmpFile)
	enc.SetIndent("", "  ")
	if err := enc.Encode(config); err != nil {
		_ = tmpFile.Close()
		cleanupTemp()
		return fmt.Errorf("failed to encode config to JSON: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		cleanupTemp()
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		cleanupTemp()
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	if err := attemptAtomicRename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to move temp file to final location: %w", err)
	}
	return nil
}

func attemptAtomicRename(tmpPath, targetPath string) error {
	// first attempt
	if err := os.Rename(tmpPath, targetPath); err == nil {
		return nil
	} else {
		// save the first rename error
		primaryErr := err

		// fallback: remove existing target and try again
		if removeErr := os.Remove(targetPath); removeErr == nil {
			if err2 := os.Rename(tmpPath, targetPath); err2 == nil {
				return nil
			} else {
				_ = os.Remove(tmpPath)
				return fmt.Errorf("failed to rename temp config file: %w", err2)
			}
		}

		// couldn't remove existing target; remove temp and report original rename error
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp config file: %w", primaryErr)
	}
}
