package main

import (
	"fmt"
	"time"

	calc "github.com/mnadev/adhango/pkg/calc"
	data "github.com/mnadev/adhango/pkg/data"
	util "github.com/mnadev/adhango/pkg/util"
)

// Dependency injection for current time (can be overridden in tests)
var nowFunc = time.Now

var ansiColors = map[string]string{
	"black":   "\033[30m",
	"red":     "\033[31m",
	"green":   "\033[32m",
	"yellow":  "\033[33m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
	"reset":   "\033[0m",
}

func highlight(text, color string) string {
	code, ok := ansiColors[color]
	if !ok {
		code = ansiColors["green"] // default
	}
	return code + text + ansiColors["reset"]
}

// getPrayerTimesForDate returns prayer times for a given date (testable)
func getPrayerTimesForDate(config *Config, params *calc.CalculationParameters, date time.Time) (*calc.PrayerTimes, error) {
	coordinates, err := util.NewCoordinates(config.Latitude, config.Longitude)
	if err != nil {
		return nil, fmt.Errorf("failed to initialise coordinates: %w", err)
	}
	return calc.NewPrayerTimes(coordinates, data.NewDateComponents(date), params)
}

// getTodaysPrayerTimes returns today's prayer times using nowFunc (testable)
func getTodaysPrayerTimes(config *Config, params *calc.CalculationParameters) (*calc.PrayerTimes, error) {
	return getPrayerTimesForDate(config, params, nowFunc())
}

// getTomorrowsPrayerTimes returns tomorrow's prayer times using nowFunc (testable)
func getTomorrowsPrayerTimes(config *Config, params *calc.CalculationParameters) (*calc.PrayerTimes, error) {
	return getPrayerTimesForDate(config, params, nowFunc().AddDate(0, 0, 1))
}

// formatPrayerTimes returns a string representation of daily prayer times (testable)
func formatPrayerTimes(times *calc.PrayerTimes, config *Config) string {
	nowPrayer := times.CurrentPrayer(time.Now())
	prayers := map[calc.Prayer]string{
		calc.FAJR:    fmt.Sprintf("Fajr %s", times.Fajr.Local().Format("15:04")),
		calc.SUNRISE: fmt.Sprintf("Sunrise %s", times.Sunrise.Local().Format("15:04")),
		calc.DHUHR:   fmt.Sprintf("Dhuhr %s", times.Dhuhr.Local().Format("15:04")),
		calc.ASR:     fmt.Sprintf("Asr %s", times.Asr.Local().Format("15:04")),
		calc.MAGHRIB: fmt.Sprintf("Maghrib %s", times.Maghrib.Local().Format("15:04")),
		calc.ISHA:    fmt.Sprintf("Isha %s", times.Isha.Local().Format("15:04")),
	}
	if config.EnableHighlighting {
		// Highlight the current prayer (name + time)
		if nowPrayer != calc.NO_PRAYER {
			prayers[nowPrayer] = highlight(prayers[nowPrayer], config.HighlightColour)
		}
	}
	return fmt.Sprintf(
		"%s | %s | %s | %s | %s | %s",
		prayers[calc.FAJR],
		prayers[calc.SUNRISE],
		prayers[calc.DHUHR],
		prayers[calc.ASR],
		prayers[calc.MAGHRIB],
		prayers[calc.ISHA],
	)
}

// nextPrayerInfo returns the name and time of the next upcoming prayer (testable)
func nextPrayerInfo(timesToday, timesTomorrow *calc.PrayerTimes, loc *time.Location, prayerNames map[calc.Prayer]string) (string, time.Time, error) {
	if nowFunc().Before(timesToday.Isha.Local()) {
		nextPrayer := timesToday.NextPrayerNow()
		if nextPrayer == calc.NO_PRAYER {
			return "", time.Time{}, fmt.Errorf("no upcoming prayer found for today")
		}
		return prayerNames[nextPrayer], timesToday.TimeForPrayer(nextPrayer).In(loc), nil
	}

	// No more prayers today; fallback to tomorrow's Fajr
	return "Fajr", timesTomorrow.Fajr.Local(), nil
}

func formatNextPrayerInfo(name string, t time.Time, config *Config) string {
	var result string
	result = fmt.Sprintf("%s %s", name, t.Format("15:04"))
	if config.EnableCountdown {
		countdown := formatCountdown(t)
		if countdown != "" {
			result = fmt.Sprintf("%s (%s)\n", result, countdown)
		}
	}
	if config.EnableHighlighting {
		result = highlight(result, config.HighlightColour)
	}
	return result
}

func formatCountdown(t time.Time) string {
	now := nowFunc()
	if t.Before(now) || t.Sub(now) < time.Second {
		return "" // No countdown shown if it's now
	}

	diff := t.Sub(now)
	if diff < time.Minute {
		return fmt.Sprintf("in %d sec", int(diff.Seconds()))
	} else if diff < time.Hour {
		return fmt.Sprintf("in %d min", int(diff.Minutes()))
	}

	hours := int(diff.Hours())
	minutes := int(diff.Minutes()) % 60
	return fmt.Sprintf("in %d hr %d min", hours, minutes)
}
