package config

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/regular_expression"
)

func NewBruteForceProtectionSSH(isNotify bool, notifyCooldown int, notifyEvery int, group *brute_force_protection.Group) ([]*Source, error) {
	var sources []*Source

	journal, err := NewSourceJournal(JournalFieldSystemdUnit, "ssh.service")
	if err != nil {
		return nil, fmt.Errorf("failed to create journal source for SSH: %w", err)
	}

	source := &Source{
		Type:    SourceTypeJournal,
		Journal: journal,
		BruteForceProtectionRule: &brute_force_protection.Rule{
			Name:    "_ssh",
			Message: i18n.Lang.T("alert.bruteForceProtection.ssh.message"),

			IsNotification:       isNotify,
			NotificationCooldown: uint32(notifyCooldown),
			NotificationEvery:    uint32(notifyEvery),

			Patterns: []brute_force_protection.RegexPattern{
				{
					Regexp: regular_expression.NewLazyRegexp(`^Failed password for (invalid user |illegal user )?(\S*) from (\S+)( port \S+ \S+\s*)`),
					Values: []brute_force_protection.PatternValue{
						{
							Name:  i18n.Lang.T("user"),
							Value: 2,
						},
					},
					IP: 3,
				},
			},
			Group: group,
		},
	}

	sources = append(sources, source)

	return sources, nil
}
