package analyzer

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
)

type LogAlertGroup struct {
	Name                 string                   `mapstructure:"name"`
	Message              string                   `mapstructure:"message"`
	RateLimitResetPeriod int                      `mapstructure:"rate_limit_reset_period"`
	RateLimits           []LogAlertGroupRateLimit `mapstructure:"rate_limits"`
}

func (g *LogAlertGroup) ToGroup() (*config.AlertGroup, error) {
	if err := g.validate(); err != nil {
		return nil, err
	}

	var rateLimits []config.RateLimit

	for _, rateLimit := range g.RateLimits {
		rLimit, err := rateLimit.ToRateLimit()
		if err != nil {
			return nil, err
		}
		rateLimits = append(rateLimits, rLimit)
	}

	return &config.AlertGroup{
		Name:                 g.Name,
		Message:              g.Message,
		RateLimits:           rateLimits,
		RateLimitResetPeriod: uint32(g.RateLimitResetPeriod),
	}, nil
}

func (g *LogAlertGroup) validate() error {
	if g.Name == "" {
		return fmt.Errorf("alert group name is empty")
	}

	if !reName.MatchString(g.Name) {
		return fmt.Errorf("alert group invalid name: %s", g.Name)
	}

	if g.RateLimitResetPeriod < 0 {
		return fmt.Errorf("alert group rate limit reset period must be positive")
	}

	if len(g.RateLimits) == 0 {
		return fmt.Errorf("alert group rate limits is empty")
	}

	return nil
}
