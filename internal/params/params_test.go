package params

import (
	"reflect"
	"salah-cli/internal/config"
	"testing"

	calc "github.com/mnadev/adhango/pkg/calc"
)

func TestBuildCalculationParams_Defaults(t *testing.T) {
	cfg := &config.Config{Latitude: 51.5, Longitude: -0.12}

	params, err := BuildCalculationParams(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if params == nil {
		t.Fatal("expected params, got nil")
	}

	if params.Method != calc.MOON_SIGHTING_COMMITTEE {
		t.Errorf("expected default method %v, got %v", calc.MOON_SIGHTING_COMMITTEE, params.Method)
	}
}

func TestBuildCalculationParams_WithMethod(t *testing.T) {
	method := int(calc.MUSLIM_WORLD_LEAGUE)
	cfg := &config.Config{Method: &method}

	params, err := BuildCalculationParams(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if params.Method != calc.MUSLIM_WORLD_LEAGUE {
		t.Errorf("expected method %v, got %v", calc.MUSLIM_WORLD_LEAGUE, params.Method)
	}
}

func TestBuildCalculationParams_AllOverrides(t *testing.T) {
	method := int(calc.EGYPTIAN)
	fajr := 18.5
	isha := 17.0
	interval := 90
	madhab := int(calc.HANAFI)
	highLat := int(calc.MIDDLE_OF_THE_NIGHT)
	adj := calc.PrayerAdjustments{FajrAdj: 2, DhuhrAdj: 1}

	cfg := &config.Config{
		Method:            &method,
		FajrAngle:         &fajr,
		IshaAngle:         &isha,
		IshaInterval:      &interval,
		Madhab:            &madhab,
		HighLatitudeRule:  &highLat,
		Adjustments:       &adj,
		MethodAdjustments: &adj,
	}

	params, err := BuildCalculationParams(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if params.Method != calc.EGYPTIAN {
		t.Errorf("expected method %v, got %v", calc.EGYPTIAN, params.Method)
	}
	if params.FajrAngle != fajr {
		t.Errorf("expected FajrAngle %v, got %v", fajr, params.FajrAngle)
	}
	if params.IshaAngle != isha {
		t.Errorf("expected IshaAngle %v, got %v", isha, params.IshaAngle)
	}
	if params.IshaInterval != interval {
		t.Errorf("expected IshaInterval %v, got %v", interval, params.IshaInterval)
	}
	if params.Madhab != calc.HANAFI {
		t.Errorf("expected Madhab %v, got %v", calc.HANAFI, params.Madhab)
	}
	if params.HighLatitudeRule != calc.MIDDLE_OF_THE_NIGHT {
		t.Errorf("expected HighLatitudeRule %v, got %v", calc.MIDDLE_OF_THE_NIGHT, params.HighLatitudeRule)
	}
	if !reflect.DeepEqual(params.Adjustments, adj) {
		t.Errorf("expected Adjustments %+v, got %+v", adj, params.Adjustments)
	}
	if !reflect.DeepEqual(params.MethodAdjustments, adj) {
		t.Errorf("expected MethodAdjustments %+v, got %+v", adj, params.MethodAdjustments)
	}
}

func TestBuildCalculationParams_PartialOverrides(t *testing.T) {
	fajr := 19.0
	highLat := int(calc.MIDDLE_OF_THE_NIGHT)
	cfg := &config.Config{
		FajrAngle:        &fajr,
		HighLatitudeRule: &highLat,
	}

	params, err := BuildCalculationParams(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if params.FajrAngle != fajr {
		t.Errorf("expected FajrAngle %v, got %v", fajr, params.FajrAngle)
	}
	if params.HighLatitudeRule != calc.MIDDLE_OF_THE_NIGHT {
		t.Errorf("expected HighLatitudeRule %v, got %v", calc.MIDDLE_OF_THE_NIGHT, params.HighLatitudeRule)
	}
}
