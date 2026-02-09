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
	var matches []string
	alertRuleIndex := analysisServices.NewAlertRuleIndex()

	for _, source := range config.Sources {
		switch source.Type {
		case config2.SourceTypeJournal:
			match := source.Journal.JournalctlMatch()
			matches = append(matches, match)
		default:
			logger.Error(fmt.Sprintf("Unknown source type: %s", source.Type))
			continue
		}

		if source.AlertRule != nil {
			err := alertRuleIndex.Add(source)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to add alert rule: %s", err))
			}
		}
	}

	systemdService := analyzerLog.NewSystemd(config.BinPath.Journalctl, matches, logger)
	analysisService := analyzerLog.NewAnalysis(alertRuleIndex, logger, notify)

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
			a.logger.Debug(fmt.Sprintf("Received log entry: %v", entry))

			a.analysis.Alert(&entry)
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
