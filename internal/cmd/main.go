package cmd

import (
	"context"
	"fmt"
	"kor-elf-shield/internal/cmd/daemon"
	"kor-elf-shield/internal/i18n"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/urfave/cli/v3"
)

type AppVersion struct {
	Version string
	Extra   string
}

func NewMainApp(appVer AppVersion, defaultConfigPath string) *cli.Command {
	app := &cli.Command{}
	app.Name = "kor-elf-shield"
	app.Usage = i18n.Lang.T("app.Usage")
	app.Description = i18n.Lang.T("app.Description")
	app.Version = appVer.Version + appVer.Extra
	app.EnableShellCompletion = true
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: i18n.Lang.T("app.flags.config.Usage"),
			Value: defaultConfigPath,
		},
	}
	app.Commands = []*cli.Command{
		daemon.CmdStart(),
	}

	return app
}

func RunMainApp(app *cli.Command, args ...string) error {
	ctx, cancel := installSignals()
	defer cancel()
	err := app.Run(ctx, args)
	if err == nil {
		return nil
	}
	if strings.HasPrefix(err.Error(), "flag provided but not defined:") {
		// the cli package should already have output the error message, so just exit
		cli.OsExiter(1)
		return err
	}
	_, _ = fmt.Fprintf(app.ErrWriter, "\u001B[31m%s\u001B[0m:\n", i18n.Lang.T("Command error"))
	_, _ = fmt.Fprintf(app.ErrWriter, "\u001B[31m%v\u001B[0m\n", err)
	cli.OsExiter(1)
	return err
}

func installSignals() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// install notify
		signalChannel := make(chan os.Signal, 1)

		signal.Notify(
			signalChannel,
			syscall.SIGINT,
			syscall.SIGTERM,
		)
		select {
		case <-signalChannel:
		case <-ctx.Done():
		}
		cancel()
		signal.Reset()
	}()

	return ctx, cancel
}
