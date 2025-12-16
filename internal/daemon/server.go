package daemon

import (
	"errors"

	firewall2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/pidfile"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/socket"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

func NewDaemon(opts DaemonOptions, logger log.Logger, notifications notifications.Notifications) (Daemon, error) {
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	pidFile, err := pidfile.New(opts.PathPidFile, logger)
	if err != nil {
		return nil, err
	}

	sock, err := socket.New(opts.PathSocketFile, logger)
	if err != nil {
		return nil, err
	}

	firewall, err := firewall2.New(opts.PathNftables, logger, opts.ConfigFirewall)

	return &daemon{
		pidFile:       pidFile,
		socket:        sock,
		logger:        logger,
		firewall:      firewall,
		notifications: notifications,
	}, nil
}
