package analyzer

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
)

type BruteForceProtectionGroup struct {
	Name                 string                               `mapstructure:"name"`
	Message              string                               `mapstructure:"message"`
	RateLimitResetPeriod int                                  `mapstructure:"rate_limit_reset_period"`
	BlockType            string                               `mapstructure:"block_type"`
	Ports                []string                             `mapstructure:"ports"`
	RateLimits           []BruteForceProtectionGroupRateLimit `mapstructure:"rate_limits"`
}

func (g *BruteForceProtectionGroup) ToGroup() (*brute_force_protection.Group, error) {
	if err := g.validate(); err != nil {
		return nil, err
	}

	var rateLimits []brute_force_protection.RateLimit

	blockType := g.BlockType
	if blockType == "" {
		blockType = "ip"
	}
	blockConfig, err := toBlockConfigBySettings(blockType, g.Ports)
	if err := err; err != nil {
		return nil, err
	}

	for _, rateLimit := range g.RateLimits {
		rLimit, err := rateLimit.ToRateLimit(blockConfig)
		if err != nil {
			return nil, err
		}
		rateLimits = append(rateLimits, rLimit)
	}

	return &brute_force_protection.Group{
		Name:                 g.Name,
		Message:              g.Message,
		RateLimits:           rateLimits,
		RateLimitResetPeriod: uint32(g.RateLimitResetPeriod),
	}, nil
}

func (g *BruteForceProtectionGroup) validate() error {
	if g.Name == "" {
		return fmt.Errorf("brute force protection group name is empty")
	}

	if !reName.MatchString(g.Name) {
		return fmt.Errorf("brute force protection group invalid name: %s", g.Name)
	}

	if g.RateLimitResetPeriod < 0 {
		return fmt.Errorf("brute force protection group rate limit reset period must be positive")
	}

	if len(g.RateLimits) == 0 {
		return fmt.Errorf("brute force protection group rate limits is empty")
	}

	return nil
}
