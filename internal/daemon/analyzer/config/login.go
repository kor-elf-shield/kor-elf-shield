package config

import "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"

func NewLoginSSH(isNotify bool) []*Source {
	var sources []*Source

	source := &Source{
		Type: SourceTypeJournal,
		Journal: &SourceJournal{
			Field: JournalFieldSystemdUnit,
			Match: "ssh.service",
		},
		AlertRule: &AlertRule{
			Name:           "_login-ssh",
			Message:        i18n.Lang.T("alert.login.ssh.message"),
			IsNotification: isNotify,
			Patterns: []AlertRegexPattern{
				{
					Regexp: NewLazyRegexp(`^Accepted (\S+) for (\S+) from (\S+) port \S+`),
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

	return sources
}

func NewLoginLocal(isNotify bool) []*Source {
	var sources []*Source

	source := &Source{
		Type: SourceTypeJournal,
		Journal: &SourceJournal{
			Field: JournalFieldSyslogIdentifier,
			Match: "login",
		},

		AlertRule: &AlertRule{
			Name:           "_login-local",
			Message:        i18n.Lang.T("alert.login.local.message"),
			IsNotification: isNotify,
			Patterns: []AlertRegexPattern{
				{
					Regexp: NewLazyRegexp(`^pam_unix\(login:session\): session opened for user (\S+)\(\S+\) by \S+`),
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

	return sources
}

func NewLoginSu(isNotify bool) []*Source {
	var sources []*Source

	source := &Source{
		Type: SourceTypeJournal,
		Journal: &SourceJournal{
			Field: JournalFieldSyslogIdentifier,
			Match: "su",
		},
		AlertRule: &AlertRule{
			Name:           "_login-su",
			Message:        i18n.Lang.T("alert.login.su.message"),
			IsNotification: isNotify,
			Patterns: []AlertRegexPattern{
				{
					Regexp: NewLazyRegexp(`^pam_unix\(su:session\): session opened for user (\S+)\(\S+\) by (\S+)\(\S+\)`),
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

	return sources
}

func NewLoginSudo(isNotify bool) []*Source {
	var sources []*Source

	source := &Source{
		Type: SourceTypeJournal,
		Journal: &SourceJournal{
			Field: JournalFieldSyslogIdentifier,
			Match: "sudo",
		},
		AlertRule: &AlertRule{
			Name:           "_login-sudo",
			Message:        i18n.Lang.T("alert.login.sudo.message"),
			IsNotification: isNotify,
			Patterns: []AlertRegexPattern{
				{
					Regexp: NewLazyRegexp(`^pam_unix\(sudo:session\): session opened for user (\S+)\(\S+\) by (\S+)\(\S+\)`),
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

	return sources
}
