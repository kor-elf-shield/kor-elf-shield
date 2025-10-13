package daemon

import (
	"kor-elf-shield/internal/daemon/pidfile"
	"kor-elf-shield/internal/log"
)

type Daemon interface {
	Run() error
}

type daemon struct {
	pidFile pidfile.PidFile
	logger  log.Logger
}

func (d *daemon) Run() error {
	if err := d.pidFile.EnsureNoOtherProcess(); err != nil {
		return err
	}

	if err := d.pidFile.Create(); err != nil {
		return err
	}
	defer d.pidFile.Remove()

	return nil
}
