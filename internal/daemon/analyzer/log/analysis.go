package log

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	analysisServices "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Analysis interface {
	SSH(entry *analysisServices.Entry) error
	Locale(entry *analysisServices.Entry) error
	Su(entry *analysisServices.Entry) error
	Sudo(entry *analysisServices.Entry) error
}

type analysis struct {
	sshService    analysisServices.Analysis
	localeService analysisServices.Analysis
	suService     analysisServices.Analysis
	sudoService   analysisServices.Analysis

	logger log.Logger
	notify notifications.Notifications
}

func NewAnalysis(config *config.Config, logger log.Logger, notify notifications.Notifications) Analysis {
	return &analysis{
		sshService:    analysisServices.NewSSH(config, logger, notify),
		localeService: analysisServices.NewLocale(config, logger, notify),
		suService:     analysisServices.NewSu(config, logger, notify),
		sudoService:   analysisServices.NewSudo(config, logger, notify),
		logger:        logger,
		notify:        notify,
	}
}

func (a *analysis) SSH(entry *analysisServices.Entry) error {
	return a.sshService.Process(entry)
}

func (a *analysis) Locale(entry *analysisServices.Entry) error {
	return a.localeService.Process(entry)
}

func (a *analysis) Su(entry *analysisServices.Entry) error {
	return a.suService.Process(entry)
}

func (a *analysis) Sudo(entry *analysisServices.Entry) error {
	return a.sudoService.Process(entry)
}
