package analyzer

import (
	"errors"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
)

type BruteForceProtection struct {
	Enabled              bool   `mapstructure:"enabled"`
	Notify               bool   `mapstructure:"notify"`
	RateLimitCount       int    `mapstructure:"rate_limit_count"`
	RateLimitPeriod      int    `mapstructure:"rate_limit_period"`
	RateLimitResetPeriod int    `mapstructure:"rate_limit_reset_period"`
	BlockingTime         int    `mapstructure:"blocking_time"`
	SSHEnable            bool   `mapstructure:"ssh_enable"`
	SSHNotify            bool   `mapstructure:"ssh_notify"`
	SSHGroup             string `mapstructure:"ssh_group"`

	Groups []BruteForceProtectionGroup
	Rules  []BruteForceProtectionRule
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
		SSHGroup:             "",

		Groups: []BruteForceProtectionGroup{},
		Rules:  []BruteForceProtectionRule{},
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
		sshGroup := "_default"
		if p.SSHGroup != "" {
			if _, ok := groups[p.SSHGroup]; !ok {
				return nil, errors.New("ssh group not found")
			}
			sshGroup = p.SSHGroup
		}
		sshSources, err := config.NewBruteForceProtectionSSH(p.Notify && p.SSHNotify, groups[sshGroup])
		if err != nil {
			return nil, err
		}
		sources = append(sources, sshSources...)
	}

	for _, rule := range p.Rules {
		if !rule.Enabled {
			continue
		}

		var group *brute_force_protection.Group
		groupName := "_default"
		if rule.Group != "" {
			groupName = rule.Group
		}
		if _, ok := groups[groupName]; !ok {
			return nil, fmt.Errorf("group %q not found", rule.Group)
		}
		group = groups[groupName]

		source, err := rule.ToSource(p.Notify, group)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
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

	for _, group := range p.Groups {
		g, err := group.ToGroup()
		if err != nil {
			return nil, err
		}
		groups[g.Name] = g
	}

	return groups, nil
}
