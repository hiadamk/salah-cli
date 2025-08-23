package main

import (
	"strings"
	"testing"
	"time"

	calc "github.com/mnadev/adhango/pkg/calc"
)

func TestGetTodaysAndTomorrowsPrayerTimes(t *testing.T) {
	cfg := &Config{Latitude: 51.5, Longitude: -0.12} // London
	params, err := buildCalculationParams(cfg)
	if err != nil {
		t.Fatalf("failed to build params: %v", err)
	}

	today, err := getTodaysPrayerTimes(cfg, params)
	if err != nil {
		t.Fatalf("failed to get today's prayer times: %v", err)
	}
	if today.Fajr.IsZero() {
		t.Errorf("expected non-zero Fajr time")
	}

	tomorrow, err := getTomorrowsPrayerTimes(cfg, params)
	if err != nil {
		t.Fatalf("failed to get tomorrow's prayer times: %v", err)
	}
	if tomorrow.Fajr.IsZero() {
		t.Errorf("expected non-zero Fajr time for tomorrow")
	}
}

func TestFormatPrayerTimes(t *testing.T) {
	cfg := &Config{Latitude: 51.5, Longitude: -0.12}
	params, _ := buildCalculationParams(cfg)

	times, _ := getTodaysPrayerTimes(cfg, params)
	out := formatPrayerTimes(times)
	if !strings.Contains(out, "Fajr") || !strings.Contains(out, "Isha") {
		t.Errorf("expected formatted string to contain prayer names, got %q", out)
	}
}

func TestNextPrayerInfo_TodayAndTomorrow(t *testing.T) {
	cfg := &Config{Latitude: 51.5, Longitude: -0.12}
	params, _ := buildCalculationParams(cfg)

	today, _ := getTodaysPrayerTimes(cfg, params)
	tomorrow, _ := getTomorrowsPrayerTimes(cfg, params)

	prayerNames := map[calc.Prayer]string{
		calc.FAJR:    "Fajr",
		calc.DHUHR:   "Dhuhr",
		calc.ASR:     "Asr",
		calc.MAGHRIB: "Maghrib",
		calc.ISHA:    "Isha",
	}

	name, tNext, err := nextPrayerInfo(today, tomorrow, time.Local, prayerNames)
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

	cfg := &Config{Latitude: 51.5, Longitude: -0.12}
	params, _ := buildCalculationParams(cfg)

	today, _ := getTodaysPrayerTimes(cfg, params)
	tomorrow, _ := getTomorrowsPrayerTimes(cfg, params)

	prayerNames := map[calc.Prayer]string{
		calc.FAJR:    "Fajr",
		calc.DHUHR:   "Dhuhr",
		calc.ASR:     "Asr",
		calc.MAGHRIB: "Maghrib",
		calc.ISHA:    "Isha",
	}

	name, tNext, err := nextPrayerInfo(today, tomorrow, time.UTC, prayerNames)
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
