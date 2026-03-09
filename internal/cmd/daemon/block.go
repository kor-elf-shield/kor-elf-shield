package daemon

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"github.com/urfave/cli/v3"
)

func CmdBlock() *cli.Command {
	return &cli.Command{
		Name:  "block",
		Usage: i18n.Lang.T("cmd.daemon.block.Usage"),
		Commands: []*cli.Command{
			{
				Name:        "add",
				Usage:       i18n.Lang.T("cmd.daemon.block.add.Usage"),
				Description: i18n.Lang.T("cmd.daemon.block.add.Description"),
				Action:      cmdBlockAdd,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "port",
						Usage: i18n.Lang.T("cmd.daemon.block.add.FlagUsage.port"),
					},
					&cli.Uint32Flag{
						Name:  "seconds",
						Usage: i18n.Lang.T("cmd.daemon.block.add.FlagUsage.seconds"),
					},
					&cli.StringFlag{
						Name:  "reason",
						Usage: i18n.Lang.T("cmd.daemon.block.add.FlagUsage.reason"),
					},
				},
			},
			{
				Name:        "delete",
				Usage:       i18n.Lang.T("cmd.daemon.block.delete.Usage"),
				Description: i18n.Lang.T("cmd.daemon.block.delete.Description"),
				Action:      cmdBlockDelete,
			},
			{
				Name:        "clear",
				Usage:       i18n.Lang.T("cmd.daemon.block.clear.Usage"),
				Description: i18n.Lang.T("cmd.daemon.block.clear.Description"),
				Action:      cmdBlockClear,
			},
		},
	}
}

func cmdBlockAdd(_ context.Context, cmd *cli.Command) error {
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

	result, err := sock.SendCommand("block_add_ip", map[string]string{
		"ip":      ip.String(),
		"port":    cmd.String("port"),
		"seconds": strconv.Itoa(int(cmd.Uint32("seconds"))),
		"reason":  cmd.String("reason"),
	})
	if err != nil {
		return err
	}

	if result != "ok" {
		return errors.New(i18n.Lang.T("cmd.error", map[string]any{
			"Error": result,
		}))
	}

	fmt.Println(i18n.Lang.T("block_add_ip_success"))
	return nil
}

func cmdBlockDelete(_ context.Context, cmd *cli.Command) error {
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

	result, err := sock.SendCommand("block_delete_ip", map[string]string{
		"ip": ip.String(),
	})
	if err != nil {
		return err
	}

	if result != "ok" {
		return errors.New(i18n.Lang.T("cmd.error", map[string]any{
			"Error": result,
		}))
	}

	fmt.Println(i18n.Lang.T("block_delete_ip_success"))
	return nil
}

func cmdBlockClear(_ context.Context, _ *cli.Command) error {
	sock, err := newSocket()
	if err != nil {
		return errors.New(i18n.Lang.T("daemon is not running"))
	}
	defer func() {
		_ = sock.Close()
	}()

	result, err := sock.Send("block_clear")
	if err != nil {
		return err
	}

	if result != "ok" {
		return errors.New(i18n.Lang.T("block_clear_error"))
	}

	fmt.Println(i18n.Lang.T("block_clear_success"))

	return nil
}
