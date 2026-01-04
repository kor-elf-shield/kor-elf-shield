package log

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	analysisServices "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Analysis interface {
	SSH(entry *analysisServices.Entry) error
}

type analysis struct {
	sshService analysisServices.Analysis

	logger log.Logger
	notify notifications.Notifications
}

func NewAnalysis(config *config.Config, logger log.Logger, notify notifications.Notifications) Analysis {
	return &analysis{
		sshService: analysisServices.NewSSH(config, logger, notify),
		logger:     logger,
		notify:     notify,
	}
}

func (a *analysis) SSH(entry *analysisServices.Entry) error {
	return a.sshService.Process(entry)
}
