package analyzer

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
)

type BruteForceProtection struct {
	Enabled              bool `mapstructure:"enabled"`
	Notify               bool `mapstructure:"notify"`
	RateLimitCount       int  `mapstructure:"rate_limit_count"`
	RateLimitPeriod      int  `mapstructure:"rate_limit_period"`
	RateLimitResetPeriod int  `mapstructure:"rate_limit_reset_period"`
	BlockingTime         int  `mapstructure:"blocking_time"`
	SSHEnable            bool `mapstructure:"ssh_enable"`
	SSHNotify            bool `mapstructure:"ssh_notify"`
}

func defaultBruteForceProtection() BruteForceProtection {
	return BruteForceProtection{
		Enabled:              true,
		Notify:               true,
		RateLimitCount:       5,
		RateLimitPeriod:      3600,
		RateLimitResetPeriod: 86400,
		BlockingTime:         3600,
		SSHEnable:            true,
		SSHNotify:            true,
	}
}

func (p *BruteForceProtection) Validate() error {
	if p.RateLimitPeriod <= 0 {
		return errors.New("rate limit period must be greater than 0")
	}

	if p.RateLimitCount <= 0 {
		return errors.New("rate limit count must be greater than 0")
	}

	if p.RateLimitResetPeriod < 0 {
		return errors.New("rate limit reset period must be positive")
	}

	if p.BlockingTime < 0 {
		return errors.New("blocking time must be positive")
	}
	return nil
}

func (p *BruteForceProtection) ToSources() ([]*config.Source, error) {
	var sources []*config.Source

	if !p.Enabled {
		return sources, nil
	}

	groups, err := p.groups()
	if err != nil {
		return nil, err
	}

	if p.SSHEnable {
		sshSources, err := config.NewBruteForceProtectionSSH(p.Notify && p.SSHNotify, groups["_default"])
		if err != nil {
			return nil, err
		}
		sources = append(sources, sshSources...)
	}

	return sources, nil
}

func (p *BruteForceProtection) groups() (map[string]*brute_force_protection.Group, error) {
	groups := make(map[string]*brute_force_protection.Group)

	groups["_default"] = &brute_force_protection.Group{
		Name:    "_default",
		Message: i18n.Lang.T("alert.bruteForceProtection.group._default.message"),
		RateLimits: []brute_force_protection.RateLimit{
			{
				Count:               uint32(p.RateLimitCount),
				Period:              uint32(p.RateLimitPeriod),
				BlockingTimeSeconds: uint32(p.BlockingTime),
			},
		},
		RateLimitResetPeriod: uint32(p.RateLimitResetPeriod),
	}

	return groups, nil
}
