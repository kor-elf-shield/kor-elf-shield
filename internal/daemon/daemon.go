package daemon

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/blocklist"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/blocking"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/guard"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/geoip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/info"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/pidfile"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/socket"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/format"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/ip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

type Daemon interface {
	Run(ctx context.Context, isTesting bool, testingInterval uint16) error
	Stop()
}

type daemon struct {
	info               info.Info
	pidFile            pidfile.PidFile
	socket             socket.Socket
	logger             log.Logger
	firewall           firewall.API
	firewallRulesGuard guard.RulesGuard
	notifications      notifications.Notifications
	analyzer           analyzer.Analyzer
	docker             docker_monitor.Docker
	blocklist          blocklist.Blocklist
	geoIPService       geoip.GeoIP

	stopCh chan struct{}
}

func (d *daemon) Run(ctx context.Context, isTesting bool, testingInterval uint16) error {
	if err := d.pidFile.EnsureNoOtherProcess(); err != nil {
		return err
	}
	if err := d.socket.EnsureNoOtherProcess(); err != nil {
		return err
	}
	if err := d.firewall.Reload(d.info); err != nil {
		d.firewall.ClearRules()
		return err
	}
	d.firewall.SavesRules()
	d.firewallRulesGuard.Run(d.info, ctx)
	defer func() {
		_ = d.firewallRulesGuard.Close()
	}()

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

	d.blocklist.Run()
	defer func() {
		_ = d.blocklist.Close()
	}()

	d.geoIPService.Run(ctx)
	defer func() {
		_ = d.geoIPService.Close()
	}()

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

func (d *daemon) socketCommand(command string, args map[string]string, socket socket.Connect) error {
	switch command {
	case "stop":
		d.stopCh <- struct{}{}
		return socket.Write("ok")

	case "status":
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		text := fmt.Sprintf(
			"ok\n\n***\n"+
				"Version:    %s\n"+
				"BuiltWith: %s\n"+
				"Uptime:     %s\n"+
				"Goroutines: %d\n"+
				"Alloc:      %s\n"+
				"HeapAlloc:  %s\n"+
				"Sys:        %s\n"+
				"HeapSys:    %s\n"+
				"NumGC:      %d\n"+
				"***\n",
			d.info.Version(),
			d.info.BuiltWith(),
			format.HumanDuration(d.info.Uptime()),
			runtime.NumGoroutine(),
			format.HumanBytes(m.Alloc),     // Alloc is the total bytes of allocated heap objects.
			format.HumanBytes(m.HeapAlloc), // HeapAlloc is the total bytes of heap memory obtained from the OS.
			format.HumanBytes(m.Sys),       // Sys is the total bytes of memory obtained from the OS.
			format.HumanBytes(m.HeapSys),   // HeapSys is the total bytes of heap memory obtained from the OS.
			m.NumGC,
		)
		return socket.Write(text)

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
	case "block_add_ip":
		if args["ip"] == "" {
			return socket.Write("ip argument is required")
		}
		ipAddr := net.ParseIP(args["ip"])
		if ipAddr == nil {
			_ = socket.Write("invalid ip address")
			return errors.New("invalid ip address")
		}

		port := args["port"]
		if port != "" {
			if err := d.cmdBlockAddIPWithPort(ipAddr, port, args); err != nil {
				_ = socket.Write("block add failed: " + err.Error())
				return err
			}
		} else {
			if err := d.cmdBlockAddIP(ipAddr, args); err != nil {
				_ = socket.Write("block add failed: " + err.Error())
				return err
			}
		}

		return socket.Write("ok")
	case "block_delete_ip":
		if args["ip"] == "" {
			return socket.Write("ip argument is required")
		}
		ipAddr := net.ParseIP(args["ip"])
		if ipAddr == nil {
			_ = socket.Write("invalid ip address")
			return errors.New("invalid ip address")
		}

		if err := d.firewall.UnblockIP(ipAddr); err != nil {
			_ = socket.Write("block delete failed: " + err.Error())
			return err
		}

		return socket.Write("ok")

	case "block_clear":
		if err := d.firewall.UnblockAllIPs(); err != nil {
			_ = socket.Write("block clear failed: " + err.Error())
			return err
		}
		return socket.Write("ok")

	case "geoip_info":
		if args["ip"] == "" {
			return socket.Write("ip argument is required")
		}
		info, err := d.geoIPService.Info(args["ip"])
		if err != nil {
			_ = socket.Write("geoip info failed: " + err.Error())
			return err
		}
		return socket.Write(info)

	case "geoip_refresh":
		ctx := context.Background()
		if err := d.geoIPService.Refresh(ctx); err != nil {
			_ = socket.Write("geoip refresh failed: " + err.Error())
			return err
		}
		_ = socket.Write("ok")
		return nil

	default:
		_ = socket.Write("unknown command")
		return errors.New("unknown command")
	}
}

func (d *daemon) cmdBlockAddIP(ip net.IP, args map[string]string) error {
	blockIP := blocking.BlockIP{
		IP: ip,
	}

	if args["seconds"] != "" {
		seconds, err := strconv.Atoi(args["seconds"])
		if err != nil {
			return err
		}
		blockIP.TimeSeconds = uint32(seconds)
	}

	if args["reason"] != "" {
		blockIP.Reason = args["reason"]
	}

	isBlock, err := d.firewall.BlockIP(blockIP)
	if err != nil {
		return err
	}
	if !isBlock {
		return errors.New("the IP address is not blocked")
	}
	return nil
}

func (d *daemon) cmdBlockAddIPWithPort(ip net.IP, port string, args map[string]string) error {
	l4Port, err := newL4PortFromString(port)
	if err != nil {
		return err
	}

	blockIP := blocking.BlockIPWithPorts{
		IP:    ip,
		Ports: []types.L4Port{l4Port},
	}

	if args["seconds"] != "" {
		seconds, err := strconv.Atoi(args["seconds"])
		if err != nil {
			return err
		}
		blockIP.TimeSeconds = uint32(seconds)
	}

	if args["reason"] != "" {
		blockIP.Reason = args["reason"]
	}

	isBlock, err := d.firewall.BlockIPWithPorts(blockIP)
	if err != nil {
		return err
	}
	if !isBlock {
		return errors.New("the IP address is not blocked")
	}
	return nil
}

func newL4PortFromString(s string) (types.L4Port, error) {
	if s == "" {
		return nil, errors.New("port is empty")
	}

	data := strings.Split(s, "/")
	protocol := types.ProtocolTCP
	port, err := strconv.Atoi(data[0])
	if err != nil {
		return nil, err
	}
	if err := validate.Port(port, "port"); err != nil {
		return nil, err
	}

	if len(data) == 2 {
		protocol, err = ip.ToProtocol(data[1])
		if err != nil {
			return nil, err
		}
	}

	return types.NewL4Port(uint16(port), protocol)
}
