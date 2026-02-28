package analyzer

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
)

type BruteForceProtectionGroupRateLimit struct {
	Count        int `mapstructure:"count"`
	Period       int `mapstructure:"period"`
	BlockingTime int `mapstructure:"blocking_time"`
}

func (l *BruteForceProtectionGroupRateLimit) ToRateLimit() (brute_force_protection.RateLimit, error) {
	if err := l.validate(); err != nil {
		return brute_force_protection.RateLimit{}, err
	}

	return brute_force_protection.RateLimit{
		Count:               uint32(l.Count),
		Period:              uint32(l.Period),
		BlockingTimeSeconds: uint32(l.BlockingTime),
	}, nil
}

func (l *BruteForceProtectionGroupRateLimit) validate() error {
	if l.Count <= 0 {
		return fmt.Errorf("count must be greater than 0")
	}
	if l.Period <= 0 {
		return fmt.Errorf("period must be greater than 0")
	}
	if l.BlockingTime < 0 {
		return fmt.Errorf("blocking time must be non-negative")
	}
	return nil
}
