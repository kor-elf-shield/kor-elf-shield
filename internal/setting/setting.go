package setting

import (
	"errors"
	"kor-elf-shield/internal/daemon"
	"kor-elf-shield/internal/i18n"
)

type setting struct {
	Testing          bool   `mapstructure:"testing"`
	TestingInterval  int    `mapstructure:"testing_interval"`
	Language         string `mapstructure:"language"`
	FallbackLanguage string `mapstructure:"fallback_language"`
	PidFile          string `mapstructure:"pid_file"`

	Log               *log
	BinaryLocations   *binaryLocations
	OtherSettingsPath *otherSettingsPath

	path string
}

func settingDefault(path string) *setting {
	return &setting{
		path: path,

		Testing:          true,
		TestingInterval:  5,
		Language:         "ru",
		FallbackLanguage: "ru",
		PidFile:          "/var/run/kor-elf-shield/kor-elf-shield.pid",

		Log:               logDefault(),
		BinaryLocations:   binaryLocationsDefault(),
		OtherSettingsPath: otherSettingsPathDefault(),
	}
}

func (s setting) ToDaemonOptions() (daemon.DaemonOptions, error) {
	if s.PidFile == "" {
		return daemon.DaemonOptions{}, errors.New(i18n.Lang.T("parameter is not specified", map[string]any{
			"Parameter": "pid_file",
		}))
	}

	if s.BinaryLocations.Nftables == "" {
		return daemon.DaemonOptions{}, errors.New(i18n.Lang.T("parameter is not specified", map[string]any{
			"Parameter": "binaryLocations.nftables",
		}))
	}

	firewallConfig, err := s.OtherSettingsPath.ToFirewallConfig()
	if err != nil {
		return daemon.DaemonOptions{}, err
	}

	return daemon.DaemonOptions{
		PathPidFile:    s.PidFile,
		PathNftables:   s.BinaryLocations.Nftables,
		ConfigFirewall: firewallConfig,
	}, nil
}

func (s setting) Path() string {
	return s.path
}
