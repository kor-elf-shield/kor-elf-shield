package config

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/regular_expression"
)

func NewLoginSSH(isNotify bool) ([]*Source, error) {
	var sources []*Source

	journal, err := NewSourceJournal(JournalFieldSystemdUnit, "ssh.service")
	if err != nil {
		return nil, fmt.Errorf("failed to create journal source for SSH login: %w", err)
	}

	source := &Source{
		Type:    SourceTypeJournal,
		Journal: journal,
		AlertRule: &AlertRule{
			Name:           "_login-ssh",
			Message:        i18n.Lang.T("alert.login.ssh.message"),
			IsNotification: isNotify,
			Patterns: []AlertRegexPattern{
				{
					Regexp: regular_expression.NewLazyRegexp(`^Accepted (\S+) for (\S+) from (\S+) port \S+`),
					Values: []PatternValue{
						{
							Name:  i18n.Lang.T("user"),
							Value: 2,
						},
						{
							Name:  "IP",
							Value: 3,
						},
					},
				},
			},
			Group: nil,
		},
	}

	sources = append(sources, source)

	return sources, nil
}

func NewLoginLocal(isNotify bool) ([]*Source, error) {
	var sources []*Source

	journal, err := NewSourceJournal(JournalFieldSyslogIdentifier, "login")
	if err != nil {
		return nil, fmt.Errorf("failed to create journal source for local login: %w", err)
	}

	source := &Source{
		Type:    SourceTypeJournal,
		Journal: journal,
		AlertRule: &AlertRule{
			Name:           "_login-local",
			Message:        i18n.Lang.T("alert.login.local.message"),
			IsNotification: isNotify,
			Patterns: []AlertRegexPattern{
				{
					Regexp: regular_expression.NewLazyRegexp(`^pam_unix\(login:session\): session opened for user (\S+)\(\S+\) by \S+`),
					Values: []PatternValue{
						{
							Name:  i18n.Lang.T("user"),
							Value: 1,
						},
					},
				},
			},
			Group: nil,
		},
	}

	sources = append(sources, source)

	return sources, nil
}

func NewLoginSu(isNotify bool) ([]*Source, error) {
	var sources []*Source

	journal, err := NewSourceJournal(JournalFieldSyslogIdentifier, "su")
	if err != nil {
		return nil, fmt.Errorf("failed to create journal source for su login: %w", err)
	}

	source := &Source{
		Type:    SourceTypeJournal,
		Journal: journal,
		AlertRule: &AlertRule{
			Name:           "_login-su",
			Message:        i18n.Lang.T("alert.login.su.message"),
			IsNotification: isNotify,
			Patterns: []AlertRegexPattern{
				{
					Regexp: regular_expression.NewLazyRegexp(`^pam_unix\(su:session\): session opened for user (\S+)\(\S+\) by (\S+)\(\S+\)`),
					Values: []PatternValue{
						{
							Name:  i18n.Lang.T("user"),
							Value: 2,
						},
						{
							Name:  i18n.Lang.T("access to user has been gained"),
							Value: 1,
						},
					},
				},
			},
			Group: nil,
		},
	}

	sources = append(sources, source)

	return sources, nil
}

func NewLoginSudo(isNotify bool) ([]*Source, error) {
	var sources []*Source

	journal, err := NewSourceJournal(JournalFieldSyslogIdentifier, "sudo")
	if err != nil {
		return nil, fmt.Errorf("failed to create journal source for sudo login: %w", err)
	}

	source := &Source{
		Type:    SourceTypeJournal,
		Journal: journal,
		AlertRule: &AlertRule{
			Name:           "_login-sudo",
			Message:        i18n.Lang.T("alert.login.sudo.message"),
			IsNotification: isNotify,
			Patterns: []AlertRegexPattern{
				{
					Regexp: regular_expression.NewLazyRegexp(`^pam_unix\(sudo:session\): session opened for user (\S+)\(\S+\) by (\S+)\(\S+\)`),
					Values: []PatternValue{
						{
							Name:  i18n.Lang.T("user"),
							Value: 2,
						},
						{
							Name:  i18n.Lang.T("access to user has been gained"),
							Value: 1,
						},
					},
				},
			},
			Group: nil,
		},
	}

	sources = append(sources, source)

	return sources, nil
}
