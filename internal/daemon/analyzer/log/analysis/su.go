package analysis

import (
	"fmt"
	"regexp"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type su struct {
	login suLogin

	logger log.Logger
	notify notifications.Notifications
}

type suLogin struct {
	enabled bool
	notify  bool
}

func NewSu(config *config.Config, logger log.Logger, notify notifications.Notifications) Analysis {
	if !config.Login.Enabled || !config.Login.Su.Enabled {
		return &EmptyAnalysis{}
	}

	return &su{
		login: suLogin{
			enabled: config.Login.Enabled && config.Login.Su.Enabled,
			notify:  config.Login.Notify && config.Login.Su.Notify,
		},

		logger: logger,
		notify: notify,
	}
}

func (l *su) Process(entry *Entry) error {
	if l.login.enabled {
		result, err := l.login.process(entry)
		if err != nil {
			l.logger.Error(fmt.Sprintf("Failed to process Su login: %s", err))
		} else if result.found {
			if l.login.notify {
				l.notify.SendAsync(notifications.Message{Subject: result.subject, Body: result.body})
			}
			l.logger.Info(fmt.Sprintf("Su login detected: %s", entry.Message))
		}
	}

	return nil
}

func (l *suLogin) process(entry *Entry) (processReturn, error) {
	re := regexp.MustCompile(`^pam_unix\(su:session\): session opened for user (\S+)\(\S+\) by (\S+)\(\S+\)`)
	matches := re.FindStringSubmatch(entry.Message)

	if matches != nil {
		user := matches[1]
		byUser := matches[2]

		return processReturn{
			found: true,
			subject: i18n.Lang.T("alert.login.su.subject", map[string]any{
				"User":   user,
				"ByUser": byUser,
			}),
			body: i18n.Lang.T("alert.login.su.body", map[string]any{
				"User":   user,
				"ByUser": byUser,
				"Log":    entry.Message,
				"Time":   entry.Time,
			}),
		}, nil
	}

	return processReturn{found: false}, nil
}
