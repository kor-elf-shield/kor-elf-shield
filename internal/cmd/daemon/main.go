package daemon

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	logger2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/socket"
)

func newSocket() (socket.Client, error) {
	if setting.Config.SocketFile == "" {
		return nil, errors.New(i18n.Lang.T("socket file is not specified"))
	}

	logger := logger2.NewFalseLogger()
	secret, err := setting.Config.ToSecret(logger)
	if err != nil {
		return nil, err
	}

	return socket.NewSocketClient(setting.Config.SocketFile, secret)
}
