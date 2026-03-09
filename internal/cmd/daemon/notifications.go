package daemon

import (
	"context"
	"errors"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"github.com/urfave/cli/v3"
)

func CmdNotifications() *cli.Command {
	return &cli.Command{
		Name:  "notifications",
		Usage: i18n.Lang.T("cmd.daemon.notifications.Usage"),
		Commands: []*cli.Command{
			{
				Name:  "queue",
				Usage: i18n.Lang.T("cmd.daemon.notifications.queue.Usage"),
				Commands: []*cli.Command{
					{
						Name:        "count",
						Usage:       i18n.Lang.T("cmd.daemon.notifications.queue.count.Usage"),
						Description: i18n.Lang.T("cmd.daemon.notifications.queue.count.Description"),
						Action:      cmdNotificationsQueueCount,
					},
					{
						Name:        "clear",
						Usage:       i18n.Lang.T("cmd.daemon.notifications.queue.clear.Usage"),
						Description: i18n.Lang.T("cmd.daemon.notifications.queue.clear.Description"),
						Action:      cmdNotificationsQueueClear,
					},
				},
			},
		},
	}
}

func cmdNotificationsQueueCount(_ context.Context, _ *cli.Command) error {
	sock, err := newSocket()
	if err != nil {
		return errors.New(i18n.Lang.T("daemon is not running"))
	}
	defer func() {
		_ = sock.Close()
	}()

	result, err := sock.Send("notifications_queue_count")
	if err != nil {
		return err
	}

	fmt.Println(i18n.Lang.T("cmd.daemon.notifications.queue.count.result", map[string]interface{}{
		"Count": result,
	}))

	return nil
}

func cmdNotificationsQueueClear(_ context.Context, _ *cli.Command) error {
	sock, err := newSocket()
	if err != nil {
		return errors.New(i18n.Lang.T("daemon is not running"))
	}
	defer func() {
		_ = sock.Close()
	}()

	result, err := sock.Send("notifications_queue_clear")
	if err != nil {
		return err
	}

	if result != "ok" {
		return errors.New(i18n.Lang.T("notifications_queue_clear_error"))
	}

	fmt.Println(i18n.Lang.T("notifications_queue_clear_success"))

	return nil
}
