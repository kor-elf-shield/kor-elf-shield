package daemon

import (
	"context"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/pidfile"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Daemon interface {
	Run(ctx context.Context, isTesting bool, testingInterval int) error
	Stop()
}

type daemon struct {
	pidFile  pidfile.PidFile
	logger   log.Logger
	firewall firewall.API
}

func (d *daemon) Run(ctx context.Context, isTesting bool, testingInterval int) error {
	if err := d.pidFile.EnsureNoOtherProcess(); err != nil {
		return err
	}
	if err := d.firewall.Reload(); err != nil {
		return err
	}
	d.firewall.SavesRules()

	if err := d.pidFile.Create(); err != nil {
		return err
	}
	defer func() {
		_ = d.pidFile.Remove()
	}()

	d.runWorker(ctx, isTesting, testingInterval)

	return nil
}

func (d *daemon) Stop() {
	d.firewall.ClearRules()
}

func (d *daemon) runWorker(ctx context.Context, isTesting bool, testingInterval int) {
	d.logger.Info("Service started")

	// Channel timer for auto-completion in test mode
	var stopTestingCh <-chan time.Time
	if isTesting && testingInterval > 0 {
		d.logger.Info("Testing mode enabled")
		stopTestingCh = time.After(time.Duration(testingInterval) * time.Minute)
	}

	for {
		select {
		case <-ctx.Done():
			d.logger.Info("Service stopped")
			return
		case <-stopTestingCh:
			d.logger.Info("Testing interval expired, stopping service")
			d.Stop()
			return
		}
	}
}
