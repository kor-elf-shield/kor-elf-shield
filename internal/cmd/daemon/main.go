package daemon

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/socket"
)

func newSocket() (socket.Client, error) {
	if setting.Config.SocketFile == "" {
		return nil, errors.New(i18n.Lang.T("socket file is not specified"))
	}
	return socket.NewSocketClient(setting.Config.SocketFile)
}
