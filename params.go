package main

import (
	"github.com/mnadev/adhango/pkg/calc"
)

func buildCalculationParams(config *Config) (*calc.CalculationParameters, error) {

	var params *calc.CalculationParameters

	if config.Method == nil {
		params = calc.GetMethodParameters(calc.MOON_SIGHTING_COMMITTEE)
	} else {
		params = calc.GetMethodParameters(calc.CalculationMethod(*config.Method))
	}

	if config.FajrAngle != nil {
		params.FajrAngle = *config.FajrAngle
	}
	if config.IshaAngle != nil {
		params.IshaAngle = *config.IshaAngle
	}
	if config.IshaInterval != nil {
		params.IshaInterval = *config.IshaInterval
	}
	if config.Madhab != nil {
		params.Madhab = calc.AsrJuristicMethod(*config.Madhab)
	}
	if config.HighLatitudeRule != nil {
		params.HighLatitudeRule = calc.HighLatitudeRule(*config.HighLatitudeRule)
	}
	if config.Adjustments != nil {
		params.Adjustments = *config.Adjustments
	}
	if config.MethodAdjustments != nil {
		params.MethodAdjustments = *config.MethodAdjustments
	}

	return params, nil

}
