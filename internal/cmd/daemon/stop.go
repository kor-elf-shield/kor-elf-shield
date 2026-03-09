package daemon

import (
	"context"
	"errors"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
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
	sock, err := newSocket()
	if err != nil {
		return errors.New(i18n.Lang.T("daemon is not running"))
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
