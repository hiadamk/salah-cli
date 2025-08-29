package main

import (
	"fmt"
	"os"
	"time"

	calc "github.com/mnadev/adhango/pkg/calc"
)

var prayerNames = map[calc.Prayer]string{
	calc.FAJR:    "Fajr",
	calc.DHUHR:   "Dhuhr",
	calc.ASR:     "Asr",
	calc.MAGHRIB: "Maghrib",
	calc.ISHA:    "Isha",
}

func printHelp() {
	fmt.Println("Usage:")
	fmt.Println("  salah-cli today             Show today's prayer times")
	fmt.Println("  salah-cli next              Show the next upcoming prayer time")
	fmt.Println("  salah-cli validate-config   Validate the config file")
	os.Exit(0)
}

func runValidateConfig() {
	cfg, err := load()
	if err != nil {
		fmt.Printf("❌ Failed to load config: %v\n", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		fmt.Printf("❌ Invalid config: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Config is valid!")
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		printHelp()
		return
	}
	command := os.Args[1]
	switch command {
	case "today":
		// Load configuration
		config, err := load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}
		params, err := buildCalculationParams(config)
		if err != nil {
			fmt.Println("Error building calculation parameters:", err)
			os.Exit(1)
		}
		todays, err := getTodaysPrayerTimes(config, params)
		if err != nil {
			fmt.Println("Failed to get today's prayer times:", err)
			os.Exit(1)
		}
		fmt.Println(formatPrayerTimes(todays, config))
	case "next":
		// Load configuration
		config, err := load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
			os.Exit(1)
		}
		params, err := buildCalculationParams(config)
		if err != nil {
			fmt.Println("Error building calculation parameters:", err)
			os.Exit(1)
		}
		todays, err := getTodaysPrayerTimes(config, params)
		if err != nil {
			fmt.Println("Failed to get today's prayer times:", err)
			os.Exit(1)
		}
		tomorrows, err := getTomorrowsPrayerTimes(config, params)
		if err != nil {
			fmt.Println("Failed to get tomorrow's prayer times:", err)
			os.Exit(1)
		}
		name, t, err := nextPrayerInfo(todays, tomorrows, time.Local, prayerNames)
		if err != nil {
			fmt.Println("Error determining next prayer:", err)
			os.Exit(1)
		}
		fmt.Println(formatNextPrayerInfo(name, t, config))
	case "validate-config":
		runValidateConfig()
	case "setup":
		config, err := setupConfig()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		configPath, err := getConfigPath()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if err := saveConfig(config, configPath); err != nil {
			fmt.Printf("failed to save created config: %v\n", err.Error())
		}
		fmt.Printf("Successfully written config file to %s\n", configPath)
		os.Exit(0)
	case "help", "--help", "-h":
		printHelp()

	default:
		fmt.Println("Unknown command:", command)
		printHelp()
		os.Exit(1)
	}
}
