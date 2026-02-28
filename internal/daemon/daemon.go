package daemon

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/pidfile"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/socket"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Daemon interface {
	Run(ctx context.Context, isTesting bool, testingInterval uint16) error
	Stop()
}

type daemon struct {
	pidFile       pidfile.PidFile
	socket        socket.Socket
	logger        log.Logger
	firewall      firewall.API
	notifications notifications.Notifications
	analyzer      analyzer.Analyzer
	docker        docker_monitor.Docker

	stopCh chan struct{}
}

func (d *daemon) Run(ctx context.Context, isTesting bool, testingInterval uint16) error {
	if err := d.pidFile.EnsureNoOtherProcess(); err != nil {
		return err
	}
	if err := d.socket.EnsureNoOtherProcess(); err != nil {
		return err
	}
	if err := d.firewall.Reload(); err != nil {
		d.firewall.ClearRules()
		return err
	}
	d.firewall.SavesRules()

	if err := d.pidFile.Create(); err != nil {
		return err
	}
	defer func() {
		_ = d.pidFile.Remove()
	}()

	if err := d.socket.Create(); err != nil {
		return err
	}
	defer func() {
		_ = d.socket.Close()
	}()

	d.notifications.Run()
	defer func() {
		_ = d.notifications.Close()
	}()

	d.analyzer.Run(ctx)
	defer func() {
		_ = d.analyzer.Close()
	}()

	if d.firewall.DockerSupport() {
		go d.docker.Run()
		defer func() {
			_ = d.docker.Close()
		}()
	}

	go d.socket.Run(ctx, d.socketCommand)
	d.runWorker(ctx, isTesting, testingInterval)

	return nil
}

func (d *daemon) Stop() {
	d.firewall.ClearRules()
	d.logger.Info("Service stopped")
}

func (d *daemon) runWorker(ctx context.Context, isTesting bool, testingInterval uint16) {
	d.logger.Info("Service started")
	d.stopCh = make(chan struct{}, 1)

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

			if err := d.notifications.DBQueueClear(); err != nil {
				d.logger.Error(fmt.Sprintf("failed to clear notifications queue: %v", err))
			}

			if err := d.analyzer.ClearDBData(); err != nil {
				d.logger.Error(fmt.Sprintf("failed to clear analyzer data: %v", err))
			}

			if err := d.firewall.ClearDBData(); err != nil {
				d.logger.Error(fmt.Sprintf("failed to clear firewall data: %v", err))
			}

			d.Stop()
			return
		case <-d.stopCh:
			d.Stop()
			return
		}
	}
}

func (d *daemon) socketCommand(command string, socket socket.Connect) error {
	switch command {
	case "stop":
		d.stopCh <- struct{}{}
		return socket.Write("ok")
	case "status":
		return socket.Write("ok")
	case "reopen_logger":
		if err := d.logger.ReOpen(); err != nil {
			_ = socket.Write("logger reopen failed: " + err.Error())
			return err
		}
		return socket.Write("ok")
	case "notifications_queue_count":
		count := d.notifications.DBQueueSize()
		return socket.Write(strconv.Itoa(count))
	case "notifications_queue_clear":
		if err := d.notifications.DBQueueClear(); err != nil {
			_ = socket.Write("notifications queue clear failed: " + err.Error())
			return err
		}
		return socket.Write("ok")
	case "ban_clear":
		if err := d.firewall.UnblockAllIPs(); err != nil {
			_ = socket.Write("ban clear failed: " + err.Error())
			return err
		}
		return socket.Write("ok")
	default:
		_ = socket.Write("unknown command")
		return errors.New("unknown command")
	}
}
