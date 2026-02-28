package analyzer

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
)

type LogAlertGroupRateLimit struct {
	Count  int `mapstructure:"count"`
	Period int `mapstructure:"period"`
}

func (l *LogAlertGroupRateLimit) ToRateLimit() (config.RateLimit, error) {
	if err := l.validate(); err != nil {
		return config.RateLimit{}, err
	}

	return config.RateLimit{
		Count:  uint32(l.Count),
		Period: uint32(l.Period),
	}, nil
}

func (l *LogAlertGroupRateLimit) validate() error {
	if l.Count <= 0 {
		return fmt.Errorf("count must be greater than 0")
	}
	if l.Period <= 0 {
		return fmt.Errorf("period must be greater than 0")
	}
	return nil
}
