package daemon

import (
	"errors"
	firewall2 "kor-elf-shield/internal/daemon/firewall"
	"kor-elf-shield/internal/daemon/pidfile"
	"kor-elf-shield/internal/log"
)

func NewDaemon(opts DaemonOptions, logger log.Logger) (Daemon, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	pidFile, err := pidfile.New(opts.PathPidFile, logger)
	if err != nil {
		return nil, err
	}

	firewall, err := firewall2.New(opts.PathNftables, logger, opts.ConfigFirewall)

	return &daemon{
		pidFile:  pidFile,
		logger:   logger,
		firewall: firewall,
	}, nil
}
