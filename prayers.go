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
func formatPrayerTimes(times *calc.PrayerTimes) string {
	return fmt.Sprintf(
		"Fajr %s | Sunrise %s | Dhuhr %s | Asr %s | Maghrib %s | Isha %s",
		times.Fajr.Local().Format("15:04"),
		times.Sunrise.Local().Format("15:04"),
		times.Dhuhr.Local().Format("15:04"),
		times.Asr.Local().Format("15:04"),
		times.Maghrib.Local().Format("15:04"),
		times.Isha.Local().Format("15:04"),
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
