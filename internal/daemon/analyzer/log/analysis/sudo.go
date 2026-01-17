package analysis

import (
	"fmt"
	"regexp"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type sudo struct {
	login sudoLogin

	logger log.Logger
	notify notifications.Notifications
}

type sudoLogin struct {
	enabled bool
	notify  bool
}

func NewSudo(config *config.Config, logger log.Logger, notify notifications.Notifications) Analysis {
	if !config.Login.Enabled || !config.Login.Su.Enabled {
		return &EmptyAnalysis{}
	}

	return &sudo{
		login: sudoLogin{
			enabled: config.Login.Enabled && config.Login.Sudo.Enabled,
			notify:  config.Login.Notify && config.Login.Sudo.Notify,
		},

		logger: logger,
		notify: notify,
	}
}

func (s *sudo) Process(entry *Entry) error {
	if s.login.enabled {
		result, err := s.login.process(entry)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to process Sudo login: %s", err))
		} else if result.found {
			if s.login.notify {
				s.notify.SendAsync(notifications.Message{Subject: result.subject, Body: result.body})
			}
			s.logger.Info(fmt.Sprintf("Sudo login detected: %s", entry.Message))
		}
	}

	return nil
}

func (s *sudoLogin) process(entry *Entry) (processReturn, error) {
	re := regexp.MustCompile(`^pam_unix\(sudo:session\): session opened for user (\S+)\(\S+\) by (\S+)\(\S+\)`)
	matches := re.FindStringSubmatch(entry.Message)

	if matches != nil {
		user := matches[1]
		byUser := matches[2]

		return processReturn{
			found: true,
			subject: i18n.Lang.T("alert.login.sudo.subject", map[string]any{
				"User":   user,
				"ByUser": byUser,
			}),
			body: i18n.Lang.T("alert.login.sudo.body", map[string]any{
				"User":   user,
				"ByUser": byUser,
				"Log":    entry.Message,
				"Time":   entry.Time,
			}),
		}, nil
	}

	return processReturn{found: false}, nil
}
