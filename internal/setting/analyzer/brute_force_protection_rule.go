package analyzer

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
)

type BruteForceProtectionRule struct {
	Enabled        bool   `mapstructure:"enabled"`
	Notify         bool   `mapstructure:"notify"`
	NotifyCooldown int    `mapstructure:"notify_cooldown_seconds"`
	NotifyEvery    int    `mapstructure:"notify_every"`
	Name           string `mapstructure:"name"`
	Message        string `mapstructure:"message"`
	Group          string `mapstructure:"group"`
	Source         Source
	Patterns       []BruteForceProtectionPattern
}

func (l *BruteForceProtectionRule) ToSource(isNotify bool, group *brute_force_protection.Group) (*config.Source, error) {
	if err := l.validate(); err != nil {
		return nil, err
	}

	if group == nil {
		return nil, fmt.Errorf("brute force protection group is empty")
	}

	source, err := l.Source.ToSource()
	if err != nil {
		return nil, err
	}

	var patterns []brute_force_protection.RegexPattern

	for _, pattern := range l.Patterns {
		p, err := pattern.ToPattern()
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, p)
	}

	if len(patterns) == 0 {
		return nil, fmt.Errorf("patterns is empty")
	}

	source.BruteForceProtectionRule = &brute_force_protection.Rule{
		Name:                 l.Name,
		Message:              l.Message,
		IsNotification:       isNotify && l.Notify,
		NotificationCooldown: uint32(l.NotifyCooldown),
		NotificationEvery:    uint32(l.NotifyEvery),
		Patterns:             patterns,
		Group:                group,
	}

	return source, nil
}

func (l *BruteForceProtectionRule) validate() error {
	if l.Name == "" {
		return fmt.Errorf("brute force protection name is empty")
	}

	if !reName.MatchString(l.Name) {
		return fmt.Errorf("brute force protection invalid name: %s", l.Name)
	}

	if l.NotifyCooldown < 0 {
		return fmt.Errorf("brute force protection notify cooldown must be positive")
	}

	if l.NotifyEvery < 0 {
		return fmt.Errorf("brute force protection notify every must be positive")
	}

	return nil
}
