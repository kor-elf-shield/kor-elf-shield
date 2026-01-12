package analysis

import (
	"fmt"
	"regexp"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type ssh struct {
	login sshLogin

	logger log.Logger
	notify notifications.Notifications
}

type sshLogin struct {
	enabled bool
	notify  bool
}

func NewSSH(config *config.Config, logger log.Logger, notify notifications.Notifications) Analysis {
	if !config.Login.Enabled || !config.Login.SSH.Enabled {
		return &EmptyAnalysis{}
	}

	return &ssh{
		login: sshLogin{
			enabled: config.Login.Enabled && config.Login.SSH.Enabled,
			notify:  config.Login.Notify && config.Login.SSH.Notify,
		},

		logger: logger,
		notify: notify,
	}
}

func (s *ssh) Process(entry *Entry) error {
	if s.login.enabled {
		result, err := s.login.processLogin(entry)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to process ssh login: %s", err))
		} else if result.found {
			if s.login.notify {
				s.notify.SendAsync(notifications.Message{Subject: result.subject, Body: result.body})
			}
			s.logger.Info(fmt.Sprintf("SSH login detected: %s", entry.Message))
		}
	}

	return nil
}

func (l *sshLogin) processLogin(entry *Entry) (processReturn, error) {
	re := regexp.MustCompile(`^Accepted (\S+) for (\S+) from (\S+) port \S+`)
	matches := re.FindStringSubmatch(entry.Message)

	if matches != nil {
		user := matches[2]
		ip := matches[3]

		return processReturn{
			found: true,
			subject: i18n.Lang.T("alert.login.subject", map[string]any{
				"User": user,
				"IP":   ip,
			}),
			body: i18n.Lang.T("alert.login.body", map[string]any{
				"User": user,
				"IP":   ip,
				"Log":  entry.Message,
				"Time": entry.Time,
			}),
		}, nil
	}

	return processReturn{found: false}, nil
}
