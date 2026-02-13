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
	files    analyzerLog.FileMonitoring
	analysis analyzerLog.Analysis

	logChan chan analysisServices.Entry
}

func New(config config2.Config, logger log.Logger, notify notifications.Notifications) Analyzer {
	var journalMatches []string
	journalMatchesUniq := map[string]struct{}{}

	var files []string
	filesUniq := map[string]struct{}{}

	rulesIndex := analysisServices.NewRulesIndex()

	for _, source := range config.Sources {
		switch source.Type {
		case config2.SourceTypeJournal:
			match := source.Journal.JournalctlMatch()
			if _, ok := journalMatchesUniq[match]; !ok {
				journalMatchesUniq[match] = struct{}{}
				journalMatches = append(journalMatches, match)
			}
		case config2.SourceTypeFile:
			file := source.File.Path
			if _, ok := filesUniq[file]; !ok {
				filesUniq[file] = struct{}{}
				files = append(files, file)
			}
		default:
			logger.Error(fmt.Sprintf("Unknown source type: %s", source.Type))
			continue
		}

		if source.AlertRule != nil {
			err := rulesIndex.Add(source)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to add alert rule: %s", err))
			}
		}
	}

	systemdService := analyzerLog.NewSystemd(config.BinPath.Journalctl, journalMatches, logger)
	filesService := analyzerLog.NewFileMonitoring(files, logger)
	analysisService := analyzerLog.NewAnalysis(rulesIndex, logger, notify)

	return &analyzer{
		config:   config,
		logger:   logger,
		notify:   notify,
		systemd:  systemdService,
		files:    filesService,
		analysis: analysisService,

		logChan: make(chan analysisServices.Entry, 1000),
	}
}

func (a *analyzer) Run(ctx context.Context) {
	go a.processLogs(ctx)
	go a.systemd.Run(ctx, a.logChan)
	go a.files.Run(ctx, a.logChan)

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
		a.logger.Error(err.Error())
	}
	if err := a.files.Close(); err != nil {
		a.logger.Error(err.Error())
	}
	close(a.logChan)

	a.logger.Debug("Analyzer is stop")

	return nil
}
