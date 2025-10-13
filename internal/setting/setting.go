package setting

import (
	"errors"
	"kor-elf-shield/internal/daemon"
	"kor-elf-shield/internal/i18n"
)

type setting struct {
	Language         string `mapstructure:"language"`
	FallbackLanguage string `mapstructure:"fallback_language"`
	PidFile          string `mapstructure:"pid_file"`

	Log *log
}

func settingDefault() *setting {
	return &setting{
		Language:         "ru",
		FallbackLanguage: "ru",
		PidFile:          "/var/run/kor-elf-shield/kor-elf-shield.pid",

		Log: logDefault(),
	}
}

func (s setting) ToDaemonOptions() (daemon.DaemonOptions, error) {
	if s.PidFile == "" {
		return daemon.DaemonOptions{}, errors.New(i18n.Lang.T("parameter is not specified", map[string]any{
			"Parameter": "pid_file",
		}))
	}

	return daemon.DaemonOptions{
		PathPidFile: s.PidFile,
	}, nil
}
