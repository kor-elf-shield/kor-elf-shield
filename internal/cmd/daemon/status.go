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

func CmdStatus() *cli.Command {
	return &cli.Command{
		Name:        "status",
		Usage:       i18n.Lang.T("cmd.daemon.status.Usage"),
		Description: i18n.Lang.T("cmd.daemon.status.Description"),
		Action:      cmdStatus,
	}
}

func cmdStatus(_ context.Context, _ *cli.Command) error {
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

	result, err := sock.Send("status")
	if err != nil {
		return err
	}

	if result != "ok" {
		return errors.New(i18n.Lang.T("daemon is not running"))
	}

	fmt.Println("ok")

	return nil
}
