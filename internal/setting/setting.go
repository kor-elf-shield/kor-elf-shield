package setting

import (
	"errors"
	"strings"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

type setting struct {
	Testing          bool   `mapstructure:"testing"`
	TestingInterval  int16  `mapstructure:"testing_interval"`
	Language         string `mapstructure:"language"`
	FallbackLanguage string `mapstructure:"fallback_language"`
	PidFile          string `mapstructure:"pid_file"`

	Log               *log
	BinaryLocations   *binaryLocations
	OtherSettingsPath *otherSettingsPath
}

func settingDefault() *setting {
	return &setting{
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

func (s setting) Validate() error {
	if err := s.validationTestingInterval(); err != nil {
		return err
	}
	if err := s.validateLanguage(); err != nil {
		return err
	}
	if err := s.validatePidFile(); err != nil {
		return err
	}

	return nil
}

func (s setting) validationTestingInterval() error {
	if s.TestingInterval < 1 {
		return errors.New("testing_interval must be greater than 0")
	}
	if s.TestingInterval > 30000 {
		return errors.New("testing_interval must be less than 30000")
	}
	return nil
}

func (s setting) validatePidFile() error {
	if err := validate.PathFile(s.PidFile, "pid_file"); err != nil {
		return err
	}
	if !strings.HasSuffix(strings.ToLower(s.PidFile), ".pid") {
		return errors.New("invalid pid_file. Must be .pid")
	}
	return nil
}

func (s setting) validateLanguage() error {
	if err := validateLanguage(s.Language, "language"); err != nil {
		return err
	}
	if err := validateLanguage(s.FallbackLanguage, "fallback_language"); err != nil {
		return err
	}
	return nil
}

func validateLanguage(language string, parameterName string) error {
	switch language {
	case "ru", "en", "kk":
		return nil
	}
	return errors.New("invalid " + parameterName + ". Must be ru, en, kk")
}
