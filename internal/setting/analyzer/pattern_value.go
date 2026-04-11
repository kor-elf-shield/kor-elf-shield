package analyzer

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
)

type PatternValue struct {
	Name  string `mapstructure:"name"`
	Value int8   `mapstructure:"value"`
	Type  string `mapstructure:"type"`
}

func (v *PatternValue) ToPatternValue() (config.PatternValue, error) {
	if err := v.validate(); err != nil {
		return config.PatternValue{}, err
	}

	value := config.PatternValue{
		Name:  v.Name,
		Value: uint8(v.Value),
	}

	if v.Type != "" {
		t, err := v.toPatternTypeValue()
		if err != nil {
			return value, err
		}
		value.Type = t
	}

	return value, nil
}

func (v *PatternValue) ToPatternValueForBruteForceProtection() (brute_force_protection.PatternValue, error) {
	if err := v.validate(); err != nil {
		return brute_force_protection.PatternValue{}, err
	}

	value := brute_force_protection.PatternValue{
		Name:  v.Name,
		Value: uint8(v.Value),
	}

	return value, nil
}

func (v *PatternValue) validate() error {
	if v.Value <= 0 {
		return fmt.Errorf("value must be greater than 0")
	}

	if v.Name == "" {
		return fmt.Errorf("name is required")
	}

	return nil
}

func (v *PatternValue) toPatternTypeValue() (config.PatternTypeValue, error) {
	if v.Type == "" {
		return "", fmt.Errorf("type is required")
	}

	switch v.Type {
	case "ip":
		return config.PatternValueIP, nil
	default:
		return "", fmt.Errorf("type not support")
	}
}
