package analyzer

import (
	"context"
	"fmt"

	config2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	analyzerLog "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log"
	analysisServices "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Analyzer interface {
	Run(ctx context.Context)
	Close() error
}

type analyzer struct {
	config   config2.Config
	logger   log.Logger
	notify   notifications.Notifications
	systemd  analyzerLog.Systemd
	analysis analyzerLog.Analysis

	logChan chan analysisServices.Entry
}

func New(config config2.Config, logger log.Logger, notify notifications.Notifications) Analyzer {
	var units []string
	if config.Login.Enabled && config.Login.SSH.Enabled {
		units = append(units, "_SYSTEMD_UNIT=ssh.service")
	}

	if config.Login.Enabled && config.Login.Local.Enabled {
		units = append(units, "SYSLOG_IDENTIFIER=login")
	}

	systemdService := analyzerLog.NewSystemd(config.BinPath.Journalctl, units, logger)
	analysisService := analyzerLog.NewAnalysis(&config, logger, notify)

	return &analyzer{
		config:   config,
		logger:   logger,
		notify:   notify,
		systemd:  systemdService,
		analysis: analysisService,

		logChan: make(chan analysisServices.Entry, 1000),
	}
}

func (a *analyzer) Run(ctx context.Context) {
	go a.systemd.Run(ctx, a.logChan)
	go a.processLogs(ctx)
	a.logger.Debug("Analyzer is start")
}

func (a *analyzer) processLogs(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-a.logChan:
			if !ok {
				// Channel closed
				return
			}
			a.logger.Debug(fmt.Sprintf("Received log entry: %s", entry))

			switch {
			case entry.Unit == "ssh.service":
				if err := a.analysis.SSH(&entry); err != nil {
					a.logger.Error(fmt.Sprintf("Failed to analyze SSH logs: %s", err))
				}
			case entry.SyslogIdentifier == "login":
				if err := a.analysis.Locale(&entry); err != nil {
					a.logger.Error(fmt.Sprintf("Failed to analyze locale logs: %s", err))
				}
			default:
				a.logger.Debug(fmt.Sprintf("Unknown unit or SyslogIdentifier: %s", entry.Unit))
			}
		}
	}
}

func (a *analyzer) Close() error {
	if err := a.systemd.Close(); err != nil {
		return err
	}
	close(a.logChan)

	a.logger.Debug("Analyzer is stop")

	return nil
}
