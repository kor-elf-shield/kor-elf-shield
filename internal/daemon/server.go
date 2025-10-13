package daemon

import (
	"kor-elf-shield/internal/daemon/pidfile"
	"kor-elf-shield/internal/log"
)

func NewDaemon(opts DaemonOptions, logger log.Logger) (Daemon, error) {
	pidFile, err := pidfile.New(opts.PathPidFile, logger)
	if err != nil {
		return nil, err
	}

	return &daemon{
		pidFile: pidFile,
		logger:  logger,
	}, nil
}
