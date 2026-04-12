package daemon

import (
	"context"
	"errors"
	"fmt"
	"net"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"github.com/urfave/cli/v3"
)

func CmdGeoIP() *cli.Command {
	return &cli.Command{
		Name:  "geoip",
		Usage: i18n.Lang.T("cmd.daemon.geoip.Usage"),
		Commands: []*cli.Command{
			{
				Name:        "info",
				Usage:       i18n.Lang.T("cmd.daemon.geoip.info.Usage"),
				Description: i18n.Lang.T("cmd.daemon.geoip.info.Description"),
				Action:      CmdGeoIPInfo,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "ip",
						Usage: i18n.Lang.T("cmd.daemon.geoip.info.FlagUsage.ip"),
					},
				},
			},
			{
				Name:        "refresh",
				Usage:       i18n.Lang.T("cmd.daemon.geoip.refresh.Usage"),
				Description: i18n.Lang.T("cmd.daemon.geoip.refresh.Description"),
				Action:      CmdGeoIPRefresh,
			},
		},
	}
}

func CmdGeoIPInfo(_ context.Context, cmd *cli.Command) error {
	ip := net.ParseIP(cmd.Args().Get(0))
	if ip == nil {
		return errors.New("invalid ip address")
	}

	sock, err := newSocket()
	if err != nil {
		return errors.New(i18n.Lang.T("daemon is not running"))
	}
	defer func() {
		_ = sock.Close()
	}()

	result, err := sock.SendCommand("geoip_info", map[string]string{
		"ip": ip.String(),
	})
	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil
}

func CmdGeoIPRefresh(_ context.Context, _ *cli.Command) error {
	sock, err := newSocket()
	if err != nil {
		return errors.New(i18n.Lang.T("daemon is not running"))
	}
	defer func() {
		_ = sock.Close()
	}()

	result, err := sock.Send("geoip_refresh")
	if err != nil {
		return err
	}

	if result != "ok" {
		return errors.New(i18n.Lang.T("cmd.error", map[string]any{
			"Error": result,
		}))
	}

	fmt.Println(i18n.Lang.T("geoip_refresh_success"))

	return nil
}
