package analyzer

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Analyzer interface {
	Run()
	Close() error
}

type analyzer struct {
	config Config
	logger log.Logger
	notify notifications.Notifications
}

func New(config Config, logger log.Logger, notify notifications.Notifications) Analyzer {
	return &analyzer{
		config: config,
		logger: logger,
		notify: notify,
	}
}

func (a *analyzer) Run() {
	a.logger.Debug("Analyzer is start")
}

func (a *analyzer) Close() error {
	a.logger.Debug("Analyzer is stop")

	return nil
}
