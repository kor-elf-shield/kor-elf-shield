package analyzer

import (
	"context"
	"fmt"

	analyzerLog "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Analyzer interface {
	Run(ctx context.Context)
	Close() error
}

type analyzer struct {
	config  Config
	logger  log.Logger
	notify  notifications.Notifications
	systemd analyzerLog.Systemd
}

func New(config Config, logger log.Logger, notify notifications.Notifications) Analyzer {
	units := []string{}
	if config.Login.Enabled && config.Login.SSH.Enabled {
		units = append(units, "ssh")
	}

	systemdService := analyzerLog.NewSystemd(config.BinPath.Journalctl, units, logger)

	return &analyzer{
		config:  config,
		logger:  logger,
		notify:  notify,
		systemd: systemdService,
	}
}

func (a *analyzer) Run(ctx context.Context) {
	logChan := make(chan analyzerLog.Entry, 1000)

	go a.systemd.Run(ctx, logChan)
	go a.processLogs(ctx, logChan)
	a.logger.Debug("Analyzer is start")
}

func (a *analyzer) processLogs(ctx context.Context, logChan <-chan analyzerLog.Entry) {
	for {
		select {
		case <-ctx.Done():
			return
		case entry := <-logChan:
			a.logger.Debug(fmt.Sprintf("Received log entry: %s", entry))
		}
	}
}

func (a *analyzer) Close() error {
	if err := a.systemd.Close(); err != nil {
		return err
	}

	a.logger.Debug("Analyzer is stop")

	return nil
}
