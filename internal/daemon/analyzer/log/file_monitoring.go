package log

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	analysisServices "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/file_monitoring"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"github.com/nxadm/tail"
)

type FileMonitoring interface {
	Run(ctx context.Context, logChan chan<- analysisServices.Entry)
	Close() error
}

type fileMonitoring struct {
	paths  []string
	logger log.Logger

	tailers []*tail.Tail
	mu      sync.Mutex
}

func NewFileMonitoring(paths []string, logger log.Logger) FileMonitoring {
	return &fileMonitoring{
		paths:   paths,
		logger:  logger,
		tailers: []*tail.Tail{},
	}
}

func (fm *fileMonitoring) Run(ctx context.Context, logChan chan<- analysisServices.Entry) {
	pathsCount := len(fm.paths)
	if pathsCount == 0 {
		fm.logger.Debug("No paths specified for file monitoring")
		return
	}

	fm.logger.Debug("Starting file monitoring")

	tailLogger := file_monitoring.NewLogger(fm.logger)

	for _, path := range fm.paths {
		path := path
		go func() {
			fm.monitorFile(path, ctx, logChan, tailLogger)
		}()
	}
}

func (fm *fileMonitoring) Close() error {
	for _, t := range fm.tailers {
		_ = t.Stop()
		fm.logger.Debug(fmt.Sprintf("Stopped monitoring file %s", t.Filename))
	}

	return nil
}

func (fm *fileMonitoring) monitorFile(path string, ctx context.Context, logChan chan<- analysisServices.Entry, tailLogger file_monitoring.Logger) {
	fm.logger.Debug(fmt.Sprintf("Monitoring file %s", path))
	t, err := tail.TailFile(path, tail.Config{
		Follow:   true,
		ReOpen:   true,
		Poll:     true,
		Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
		Logger:   tailLogger,
	})

	fm.mu.Lock()
	fm.tailers = append(fm.tailers, t)
	fm.mu.Unlock()

	if err != nil {
		fm.logger.Error(fmt.Sprintf("Failed to tail file %s: %s", path, err))
		return
	}

	for {
		select {
		case <-ctx.Done():
			return

		case line, ok := <-t.Lines:
			if !ok {
				return
			}
			if line == nil {
				continue
			}

			entry := analysisServices.Entry{
				Source:  config.SourceTypeFile,
				File:    path,
				Message: line.Text,
				Time:    time.Now(),
			}
			select {
			case <-ctx.Done():
				return
			case logChan <- entry:
			}
		}
	}
}
