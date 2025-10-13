package pidfile

import (
	"errors"
	"fmt"
	"kor-elf-shield/internal/log"
	"kor-elf-shield/internal/pkg/filesystem"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

type PidFile interface {
	EnsureNoOtherProcess() error
	Create() error
	Remove() error
}

type pidFile struct {
	path   string
	logger log.Logger
}

func New(path string, logger log.Logger) (PidFile, error) {
	if path == "" {
		return nil, errors.New("path is empty")
	}

	return &pidFile{
		path:   path,
		logger: logger,
	}, nil
}

func (p *pidFile) EnsureNoOtherProcess() error {
	if p.path == "" {
		return errors.New("PID file not specified")
	}

	if _, err := os.Stat(p.path); os.IsNotExist(err) {
		// File not found, so everything is fine.
		return nil
	} else if err != nil {
		return err
	}

	// Determine the pid
	data, err := os.ReadFile(p.path)
	if err != nil {
		return fmt.Errorf("there is a pid file and we couldn't read the file: %w", err)
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return fmt.Errorf("pid file is damaged: %w", err)
	}

	// Check if a process with this PID exists
	// If not, the PID file is outdated and we delete it.
	process, err := os.FindProcess(pid)
	if err == nil {
		// Let's try sending signal 0 (is there a process?).
		err := process.Signal(syscall.Signal(0))
		if err == nil {
			return fmt.Errorf("the daemon is already running with pid %d", pid)
		}
		if err.Error() == "os: process already finished" {
			// The process is complete, you can delete the file below
		}
	}
	// The file is outdated, try deleting it.
	_ = os.Remove(p.path)
	p.logger.Warn(fmt.Sprintf("An obsolete pid file %s with pid %d was found: file removed", p.path, pid))
	return nil
}

func (p *pidFile) Create() error {
	if p.path == "" {
		return errors.New("PID file not specified")
	}

	dir := filepath.Dir(p.path)
	err := filesystem.EnsureDir(dir)
	if err != nil {
		return err
	}

	file, err := os.Create(p.path)
	if err != nil {
		return err
	}
	defer file.Close()

	pid := os.Getpid()
	p.logger.Debug(fmt.Sprintf("write pid file: %d", pid))
	_, err = file.WriteString(strconv.Itoa(pid))

	return err
}

func (p *pidFile) Remove() error {
	if p.path == "" {
		p.logger.Warn("PID file not specified")
		return nil
	}

	p.logger.Debug("remove pid file")
	return os.Remove(p.path)
}
