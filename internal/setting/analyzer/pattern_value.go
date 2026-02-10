package analyzer

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
)

type PatternValue struct {
	Name  string `mapstructure:"name"`
	Value int8   `mapstructure:"value"`
}

func (v *PatternValue) ToPatternValue() (config.PatternValue, error) {
	if err := v.validate(); err != nil {
		return config.PatternValue{}, err
	}

	value := config.PatternValue{
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
