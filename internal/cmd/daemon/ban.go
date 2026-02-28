package daemon

import (
	"context"
	"errors"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/socket"
	"github.com/urfave/cli/v3"
)

func CmdBan() *cli.Command {
	return &cli.Command{
		Name:  "ban",
		Usage: i18n.Lang.T("cmd.daemon.ban.Usage"),
		Commands: []*cli.Command{
			{
				Name:        "clear",
				Usage:       i18n.Lang.T("cmd.daemon.ban.clear.Usage"),
				Description: i18n.Lang.T("cmd.daemon.ban.clear.Description"),
				Action:      CmdBanClear,
			},
		},
	}
}

func CmdBanClear(_ context.Context, _ *cli.Command) error {
	if setting.Config.SocketFile == "" {
		return errors.New(i18n.Lang.T("socket file is not specified"))
	}
	sock, err := socket.NewSocketClient(setting.Config.SocketFile)
	if err != nil {
		return errors.New(i18n.Lang.T("daemon is not running"))
	}
	defer func() {
		_ = sock.Close()
	}()

	result, err := sock.Send("ban_clear")
	if err != nil {
		return err
	}

	if result != "ok" {
		return errors.New(i18n.Lang.T("ban_clear_error"))
	}

	fmt.Println(i18n.Lang.T("ban_clear_success"))

	return nil
}
