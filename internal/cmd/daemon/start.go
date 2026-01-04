package daemon

import (
	"context"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting"

	"github.com/urfave/cli/v3"
)

func CmdStart() *cli.Command {
	return &cli.Command{
		Name:        "start",
		Usage:       i18n.Lang.T("cmd.daemon.start.Usage"),
		Description: i18n.Lang.T("cmd.daemon.start.Description"),
		Action:      runDaemon,
	}
}

func runDaemon(ctx context.Context, _ *cli.Command) error {
	logOptions, err := setting.Config.Log.ToLoggerOptions()
	if err != nil {
		return err
	}
	logger, err := log.NewLogger(logOptions)
	if err != nil {
		return err
	}

	defer func() {
		_ = logger.Sync()
	}()

	config, err := setting.Config.ToDaemonOptions()
	if err != nil {
		logger.Fatal(err.Error())

		// Fatal should call os.Exit(1), but there's a chance that might not happen,
		// so we return err just in case.
		return err
	}

	notificationsService, err := newNotificationsService(logger)
	if err != nil {
		logger.Fatal(err.Error())

		// Fatal should call os.Exit(1), but there's a chance that might not happen,
		// so we return err just in case.return err
		return err
	}

	dockerService, err := newDockerService(ctx, logger, config.ConfigFirewall.Options.DockerSupport)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create docker service: %s", err))
	}

	d, err := daemon.NewDaemon(config, logger, notificationsService, dockerService)
	if err != nil {
		logger.Fatal(err.Error())

		// Fatal should call os.Exit(1), but there's a chance that might not happen,
		// so we return err just in case.return err
		return err
	}

	err = d.Run(ctx, setting.Config.Testing, uint16(setting.Config.TestingInterval))
	if err != nil {
		logger.Fatal(err.Error())

		// Fatal should call os.Exit(1), but there's a chance that might not happen,
		// so we return err just in case.
		return err
	}

	return nil
}

func newNotificationsService(logger log.Logger) (notifications.Notifications, error) {
	config, err := setting.Config.OtherSettingsPath.ToNotificationsConfig()
	if err != nil {
		return nil, err
	}

	return notifications.New(config, logger), nil
}

func newDockerService(ctx context.Context, logger log.Logger, dockerSupport bool) (dockerService docker_monitor.Docker, err error) {
	if dockerSupport {
		dockerPath := setting.Config.BinaryLocations.Docker
		if dockerPath == "" {
			return docker_monitor.NewDockerNotSupport(), fmt.Errorf("docker path is empty")
		}

		dockerService = docker_monitor.New(dockerPath, ctx, logger)
	} else {
		dockerService = docker_monitor.NewDockerNotSupport()
	}

	return dockerService, nil
}
