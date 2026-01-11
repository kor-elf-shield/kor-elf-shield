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

	dockerService, dockerSupport, err := newDockerService(ctx, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create docker service: %s", err))
	}

	config, err := setting.Config.ToDaemonOptions(dockerSupport)
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

func newDockerService(ctx context.Context, logger log.Logger) (dockerService docker_monitor.Docker, dockerSupport bool, err error) {
	config, dockerSupport, err := setting.Config.OtherSettingsPath.ToDockerConfig(setting.Config.BinaryLocations)
	if err != nil {
		return docker_monitor.NewDockerNotSupport(), false, err
	}

	if !dockerSupport {
		dockerService = docker_monitor.NewDockerNotSupport()
		return dockerService, false, nil
	}

	dockerService, err = docker_monitor.New(&config, ctx, logger)
	if err != nil {
		return docker_monitor.NewDockerNotSupport(), false, err
	}

	return dockerService, dockerSupport, nil
}
