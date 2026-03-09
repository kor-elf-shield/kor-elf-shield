package daemon

import (
	"context"
	"errors"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"github.com/urfave/cli/v3"
)

func CmdBlock() *cli.Command {
	return &cli.Command{
		Name:  "block",
		Usage: i18n.Lang.T("cmd.daemon.block.Usage"),
		Commands: []*cli.Command{
			{
				Name:        "clear",
				Usage:       i18n.Lang.T("cmd.daemon.block.clear.Usage"),
				Description: i18n.Lang.T("cmd.daemon.block.clear.Description"),
				Action:      cmdBlockClear,
			},
		},
	}
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
