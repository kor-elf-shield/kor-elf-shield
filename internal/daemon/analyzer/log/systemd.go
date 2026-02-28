package log

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	analysisServices "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/log/analysis"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Systemd interface {
	Run(ctx context.Context, logChan chan<- analysisServices.Entry)
	Close() error
}

type systemd struct {
	path    string
	matches []string
	logger  log.Logger

	cmd *exec.Cmd
	mu  sync.Mutex
}

type journalRawEntry struct {
	Message           string `json:"MESSAGE"`
	Unit              string `json:"_SYSTEMD_UNIT"`
	PID               string `json:"_PID"`
	SyslogIdentifier  string `json:"SYSLOG_IDENTIFIER"`
	SourceTimestamp   string `json:"_SOURCE_REALTIME_TIMESTAMP"`
	RealtimeTimestamp string `json:"__REALTIME_TIMESTAMP"`
}

func NewSystemd(path string, matches []string, logger log.Logger) Systemd {
	return &systemd{
		path:    path,
		matches: matches,
		logger:  logger,
	}
}

func (s *systemd) Run(ctx context.Context, logChan chan<- analysisServices.Entry) {
	if len(s.matches) == 0 {
		s.logger.Debug("No matches specified for journalctl")
		return
	}

	s.logger.Debug("Journalctl started")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := s.watch(ctx, logChan); err != nil {
				s.logger.Error(fmt.Sprintf("Journalctl exited with error: %v", err))
			}

			// Pause before restarting to avoid CPU load during persistent errors
			select {
			case <-ctx.Done():
				return
			case <-time.After(5 * time.Second):
				s.logger.Warn("Journalctl connection lost. Restarting in 5s...")
				continue
			}
		}
	}
}

func (s *systemd) watch(ctx context.Context, logChan chan<- analysisServices.Entry) error {
	args := []string{"-f", "-n", "0", "-o", "json"}
	for index, match := range s.matches {
		if index > 0 {
			args = append(args, "+")
		}
		args = append(args, match)
	}
	cmd := exec.CommandContext(ctx, s.path, args...)

	s.mu.Lock()
	s.cmd = cmd
	s.mu.Unlock()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe error: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start error: %w", err)
	}

	decoder := json.NewDecoder(stdout)
	for {
		var raw journalRawEntry
		if err := decoder.Decode(&raw); err != nil {
			if err == io.EOF {
				break // The process terminated normally or was killed.
			}
			return fmt.Errorf("decode error: %w", err)
		}

		tsStr := raw.SourceTimestamp
		if tsStr == "" {
			tsStr = raw.RealtimeTimestamp
		}

		var entryTime time.Time
		if usec, err := strconv.ParseInt(tsStr, 10, 64); err == nil {
			entryTime = time.Unix(0, usec*int64(time.Microsecond))
		} else {
			entryTime = time.Now()
		}

		entry := analysisServices.Entry{
			Source:           config.SourceTypeJournal,
			Message:          raw.Message,
			Unit:             raw.Unit,
			PID:              raw.PID,
			SyslogIdentifier: raw.SyslogIdentifier,
			Time:             entryTime,
		}

		select {
		case <-ctx.Done():
			break
		case logChan <- entry:
		}
	}

	return cmd.Wait()
}

func (s *systemd) Close() error {
	if s.matches == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cmd != nil && s.cmd.Process != nil {
		s.logger.Debug("Stopping journalctl")

		// Force journalctl to quit on shutdown
		return s.cmd.Process.Kill()
	}

	s.logger.Debug("Journalctl stopped")

	return nil
}
