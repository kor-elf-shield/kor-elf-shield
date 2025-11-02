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

func CmdStop() *cli.Command {
	return &cli.Command{
		Name:        "stop",
		Usage:       i18n.Lang.T("cmd.daemon.stop.Usage"),
		Description: i18n.Lang.T("cmd.daemon.stop.Description"),
		Action:      stopDaemon,
	}
}

func stopDaemon(_ context.Context, _ *cli.Command) error {
	if setting.Config.SocketFile == "" {
		return errors.New(i18n.Lang.T("socket file is not specified"))
	}

	sock, err := socket.NewSocketClient(setting.Config.SocketFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = sock.Close()
	}()

	result, err := sock.Send("stop")
	if err != nil {
		return err
	}

	if result != "ok" {
		return errors.New(i18n.Lang.T("daemon stop failed"))
	}

	fmt.Println(i18n.Lang.T("daemon stopped"))

	return nil
}
