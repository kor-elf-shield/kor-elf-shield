package daemon

import (
	"context"
	"errors"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"github.com/urfave/cli/v3"
)

func CmdReopenLogger() *cli.Command {
	return &cli.Command{
		Name:        "reopen_logger",
		Usage:       i18n.Lang.T("cmd.daemon.reopen_logger.Usage"),
		Description: i18n.Lang.T("cmd.daemon.reopen_logger.Description"),
		Action:      cmdReopenLogger,
	}
}

func cmdReopenLogger(_ context.Context, _ *cli.Command) error {
	sock, err := newSocket()
	if err != nil {
		return errors.New(i18n.Lang.T("daemon is not running"))
	}
	defer func() {
		_ = sock.Close()
	}()

	result, err := sock.Send("reopen_logger")
	if err != nil {
		return err
	}

	if result != "ok" {
		return errors.New(i18n.Lang.T("daemon is not reopening logger"))
	}

	fmt.Println("ok")

	return nil
}
