package analysis

import (
	"fmt"
	"regexp"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type locale struct {
	login localeLogin

	logger log.Logger
	notify notifications.Notifications
}

type localeLogin struct {
	enabled bool
	notify  bool
}

func NewLocale(config *config.Config, logger log.Logger, notify notifications.Notifications) Analysis {
	if !config.Login.Enabled || !config.Login.Local.Enabled {
		return &EmptyAnalysis{}
	}

	return &locale{
		login: localeLogin{
			enabled: config.Login.Enabled && config.Login.SSH.Enabled,
			notify:  config.Login.Notify && config.Login.SSH.Notify,
		},

		logger: logger,
		notify: notify,
	}
}

func (l *locale) Process(entry *Entry) error {
	if l.login.enabled {
		result, err := l.login.process(entry)
		if err != nil {
			l.logger.Error(fmt.Sprintf("Failed to process TTY login: %s", err))
		} else if result.found {
			if l.login.notify {
				l.notify.SendAsync(notifications.Message{Subject: result.subject, Body: result.body})
			}
			l.logger.Info(fmt.Sprintf("TTY login detected: %s", entry.Message))
		}
	}

	return nil
}

func (l *localeLogin) process(entry *Entry) (processReturn, error) {
	re := regexp.MustCompile(`^pam_unix\(login:session\): session opened for user (\S+)\(\S+\) by \S+`)
	matches := re.FindStringSubmatch(entry.Message)

	if matches != nil {
		user := matches[1]

		return processReturn{
			found: true,
			subject: i18n.Lang.T("alert.login.locale.subject", map[string]any{
				"User": user,
			}),
			body: i18n.Lang.T("alert.login.locale.body", map[string]any{
				"User": user,
				"Log":  entry.Message,
				"Time": entry.Time,
			}),
		}, nil
	}

	return processReturn{found: false}, nil
}
