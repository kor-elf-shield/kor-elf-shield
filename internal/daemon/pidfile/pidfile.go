package pidfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/filesystem"
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
	if logger == nil {
		return nil, errors.New("logger is nil")
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

	// Check if the file is a regular file
	info, err := os.Lstat(p.path)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return errors.New("PID file is not a regular file")
	}

	// Determine the pid
	data, err := os.ReadFile(p.path)
	if err != nil {
		return fmt.Errorf("there is a PID file and we couldn't read the file: %w", err)
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return fmt.Errorf("PID file is damaged: %w", err)
	}

	// Check if a process with this PID exists
	// If not, the PID file is outdated and we delete it.
	process, err := os.FindProcess(pid)
	if err == nil {
		// Let's try sending signal 0 (is there a process?).
		err := process.Signal(syscall.Signal(0))
		if err == nil {
			return fmt.Errorf("the daemon is already running with PID %d", pid)
		}
		if err.Error() == "os: process already finished" {
			// The process is complete, you can delete the file below
			_ = os.Remove(p.path)
			p.logger.Warn(fmt.Sprintf("An obsolete PID file %s with PID %d was found: file removed", p.path, pid))
		}
	}

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
	defer func() {
		_ = file.Close()
	}()

	pid := os.Getpid()
	p.logger.Debug(fmt.Sprintf("Write PID file: %d", pid))
	_, err = file.WriteString(strconv.Itoa(pid))

	return err
}

func (p *pidFile) Remove() error {
	if p.path == "" {
		p.logger.Warn("PID file not specified")
		return nil
	}

	p.logger.Debug("Remove PID file")
	return os.Remove(p.path)
}
