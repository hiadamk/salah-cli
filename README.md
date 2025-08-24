# Salah CLI

A developer-focused command-line tool for calculating Islamic prayer
times using configurable parameters. Built in **Go**, leveraging
[`github.com/mnadev/adhango`](https://github.com/mnadev/adhango) for
core prayer time calculations.

------------------------------------------------------------------------

## Project Structure

    .
    ├── config.go        # Configuration handling (loads/saves user config)
    ├── prayers.go       # Core prayer time retrieval and formatting logic
    ├── params.go        # Builds calculation parameters from config
    ├── main.go          # CLI entry point
    └── *_test.go        # Unit tests for all modules

------------------------------------------------------------------------

## Building

### Requirements

-   **Go 1.21+**

### Build Binary

``` bash
git clone https://github.com/<your-username>/salah-cli.git
cd salah-cli
go build -o salah-cli
```

Run the CLI:

``` bash
./salah-cli --help
```

------------------------------------------------------------------------

## Configuration

### Default Config Location

-   **Linux/macOS:** `~/.config/salah-cli/config.json`
-   **Windows:** `%APPDATA%\salah-cli\config.json`

### Configuration Options

  --------------------------------------------------------------------------------------
  Field                  Type      Required   Description
  ---------------------- --------- ---------- ------------------------------------------
  `latitude`             float64   Yes        Latitude of the location for prayer times.

  `longitude`            float64   Yes        Longitude of the location for prayer
                                              times.

  `method`               int       No         Calculation method (default: Muslim World
                                              League).

  `fajr_angle`           float64   No         Custom Fajr angle in degrees.

  `isha_angle`           float64   No         Custom Isha angle in degrees.

  `isha_interval`        int       No         Interval (minutes) after Maghrib for Isha.

  `madhab`               int       No         Asr juristic method (0 = Shafi, 1 =
                                              Hanafi).

  `high_latitude_rule`   int       No         Rule for high latitude adjustments.

  `adjustments`          object    No         Prayer-specific adjustments (in minutes).

  `method_adjustments`   object    No         Adjustments specific to calculation
                                              method.
  --------------------------------------------------------------------------------------

### Example Config

``` json
{
  "latitude": 51.5074,
  "longitude": -0.1278,
  "method": 2,
  "fajr_angle": 18.0,
  "isha_angle": 18.0,
  "madhab": 0,
  "high_latitude_rule": 1,
  "adjustments": {
    "FajrAdj": 2,
    "SunriseAdj": 0,
    "DhuhrAdj": 1,
    "AsrAdj": 0,
    "MaghribAdj": 2,
    "IshaAdj": 3
  }
}
```

------------------------------------------------------------------------

## Usage

``` bash
salah-cli today    # Show today's prayer times
salah-cli next     # Show next upcoming prayer
salah-cli --help   # Show usage instructions
```

Example:

``` bash
$ salah-cli today
Fajr 05:15 | Sunrise 06:48 | Dhuhr 12:30 | Asr 15:45 | Maghrib 18:10 | Isha 19:30

$ salah-cli next
Upcoming: Dhuhr 12:30
```

------------------------------------------------------------------------

## Development

Run all tests:

``` bash
go test ./...
```

Run locally:

``` bash
go run . today
```

Linting (optional but recommended):

``` bash
go vet ./...
```

------------------------------------------------------------------------

## Testing Notes

-   **Unit tests:** Use real values (no mocks) for integration-like
    coverage of `prayers.go` and `config.go`.
-   **Config tests:** Write to temporary directories using
    `t.TempDir()`.
-   **Future improvement:** Refactor packages (`config`, `prayers`,
    `params`) to improve testability and reuse.

------------------------------------------------------------------------

## Roadmap

See ROADMAP.md for upcoming planned features
