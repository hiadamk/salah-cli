package prayers

import (
	"salah-cli/internal/config"
	"salah-cli/internal/params"
	"strings"
	"testing"
	"time"

	calc "github.com/mnadev/adhango/pkg/calc"
)

func TestGetTodaysAndTomorrowsPrayerTimes(t *testing.T) {
	cfg := &config.Config{Latitude: 51.5, Longitude: -0.12} // London
	params, err := params.BuildCalculationParams(cfg)
	if err != nil {
		t.Fatalf("failed to build params: %v", err)
	}

	today, err := GetTodaysPrayerTimes(cfg, params)
	if err != nil {
		t.Fatalf("failed to get today's prayer times: %v", err)
	}
	if today.Fajr.IsZero() {
		t.Errorf("expected non-zero Fajr time")
	}

	tomorrow, err := GetTomorrowsPrayerTimes(cfg, params)
	if err != nil {
		t.Fatalf("failed to get tomorrow's prayer times: %v", err)
	}
	if tomorrow.Fajr.IsZero() {
		t.Errorf("expected non-zero Fajr time for tomorrow")
	}
}

func TestFormatPrayerTimes(t *testing.T) {
	cfg := &config.Config{Latitude: 51.5, Longitude: -0.12}
	params, _ := params.BuildCalculationParams(cfg)

	times, _ := GetTodaysPrayerTimes(cfg, params)
	out := FormatPrayerTimes(times, cfg)
	if !strings.Contains(out, "Fajr") || !strings.Contains(out, "Isha") {
		t.Errorf("expected formatted string to contain prayer names, got %q", out)
	}
}

func TestNextPrayerInfo_TodayAndTomorrow(t *testing.T) {
	cfg := &config.Config{Latitude: 51.5, Longitude: -0.12}
	params, _ := params.BuildCalculationParams(cfg)

	today, _ := GetTodaysPrayerTimes(cfg, params)
	tomorrow, _ := GetTomorrowsPrayerTimes(cfg, params)

	prayerNames := map[calc.Prayer]string{
		calc.FAJR:    "Fajr",
		calc.DHUHR:   "Dhuhr",
		calc.ASR:     "Asr",
		calc.MAGHRIB: "Maghrib",
		calc.ISHA:    "Isha",
	}

	name, tNext, err := NextPrayerInfo(today, tomorrow, time.Local, prayerNames)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if name == "" {
		t.Errorf("expected a next prayer name, got empty string")
	}
	if tNext.IsZero() {
		t.Errorf("expected a valid next prayer time, got zero")
	}
}

func TestNextPrayerInfo_AtEndOfDay(t *testing.T) {
	originalNow := nowFunc
	defer func() { nowFunc = originalNow }()

	// Simulate current time after Isha
	nowFunc = func() time.Time {
		return time.Date(2025, 8, 23, 23, 0, 0, 0, time.UTC)
	}

	cfg := &config.Config{Latitude: 51.5, Longitude: -0.12}
	params, _ := params.BuildCalculationParams(cfg)

	today, _ := GetTodaysPrayerTimes(cfg, params)
	tomorrow, _ := GetTomorrowsPrayerTimes(cfg, params)

	prayerNames := map[calc.Prayer]string{
		calc.FAJR:    "Fajr",
		calc.DHUHR:   "Dhuhr",
		calc.ASR:     "Asr",
		calc.MAGHRIB: "Maghrib",
		calc.ISHA:    "Isha",
	}

	name, tNext, err := NextPrayerInfo(today, tomorrow, time.UTC, prayerNames)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if name != "Fajr" {
		t.Errorf("expected next prayer Fajr tomorrow, got %v", name)
	}
	if tNext.IsZero() {
		t.Errorf("expected valid Fajr time, got zero")
	}
}

func TestFormatCountdown(t *testing.T) {
	// Override nowFunc for predictable testing
	baseTime := time.Date(2025, 8, 27, 12, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return baseTime }

	tests := []struct {
		name     string
		target   time.Time
		expected string
	}{
		{
			name:     "More than 1 hour",
			target:   baseTime.Add(2*time.Hour + 15*time.Minute),
			expected: "in 2 hr 15 min",
		},
		{
			name:     "Less than 1 hour",
			target:   baseTime.Add(45 * time.Minute),
			expected: "in 45 min",
		},
		{
			name:     "Less than 1 minute",
			target:   baseTime.Add(30 * time.Second),
			expected: "in 30 sec",
		},
		{
			name:     "Exactly now",
			target:   baseTime,
			expected: "",
		},
		{
			name:     "In the past",
			target:   baseTime.Add(-10 * time.Minute),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatCountdown(tt.target)
			if got != tt.expected {
				t.Errorf("formatCountdown(%v) = %q; want %q", tt.target, got, tt.expected)
			}
		})
	}
}

func TestFormatNextPrayerInfo(t *testing.T) {
	// Freeze "now" so countdowns are deterministic
	fixedNow := time.Date(2025, 8, 27, 18, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time { return fixedNow }

	tests := []struct {
		name     string
		prayer   string
		prayerAt time.Time
		config   *config.Config
		expected string
	}{
		{
			name:     "no countdown, no highlight",
			prayer:   "Maghrib",
			prayerAt: fixedNow.Add(2 * time.Hour),
			config:   &config.Config{EnableCountdown: false, EnableHighlighting: false},
			expected: "Maghrib 20:00",
		},
		{
			name:     "with countdown, no highlight",
			prayer:   "Maghrib",
			prayerAt: fixedNow.Add(90 * time.Minute),
			config:   &config.Config{EnableCountdown: true, EnableHighlighting: false},
			expected: "Maghrib 19:30 (in 1 hr 30 min)\n",
		},
		{
			name:     "with countdown <1hr, no highlight",
			prayer:   "Isha",
			prayerAt: fixedNow.Add(30 * time.Minute),
			config:   &config.Config{EnableCountdown: true, EnableHighlighting: false},
			expected: "Isha 18:30 (in 30 min)\n",
		},
		{
			name:     "highlight only",
			prayer:   "Maghrib",
			prayerAt: fixedNow.Add(2 * time.Hour),
			config:   &config.Config{EnableCountdown: false, EnableHighlighting: true, HighlightColour: "red"},
			expected: "\033[31mMaghrib 20:00\033[0m",
		},
		{
			name:     "highlight + countdown",
			prayer:   "Isha",
			prayerAt: fixedNow.Add(45 * time.Minute),
			config:   &config.Config{EnableCountdown: true, EnableHighlighting: true, HighlightColour: "blue"},
			expected: "\033[34mIsha 18:45 (in 45 min)\n\033[0m",
		},
		{
			name:     "highlight fallback to green if invalid",
			prayer:   "Fajr",
			prayerAt: fixedNow.Add(8 * time.Hour),
			config:   &config.Config{EnableCountdown: false, EnableHighlighting: true, HighlightColour: "invalidColor"},
			expected: "\033[32mFajr 02:00\033[0m", // wraps in green by default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatNextPrayerInfo(tt.prayer, tt.prayerAt, tt.config)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
