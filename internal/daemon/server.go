package daemon

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor"
	firewall2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/blocking"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/pidfile"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/socket"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

func NewDaemon(opts DaemonOptions, logger log.Logger, notifications notifications.Notifications, docker docker_monitor.Docker) (Daemon, error) {
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

	blockingService := blocking.New(opts.Repositories.Blocking(), logger)
	firewall, err := firewall2.New(opts.PathNftables, blockingService, logger, opts.ConfigFirewall, docker)

	analyzerService := analyzer.New(opts.ConfigAnalyzer, firewall.BlockIP, opts.Repositories, logger, notifications)

	return &daemon{
		pidFile:       pidFile,
		socket:        sock,
		logger:        logger,
		firewall:      firewall,
		notifications: notifications,
		analyzer:      analyzerService,
		docker:        docker,
	}, nil
}
