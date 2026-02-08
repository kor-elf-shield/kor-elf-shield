package notifications

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
	"github.com/spf13/viper"
)

type Setting struct {
	Enabled       bool   `mapstructure:"enabled"`
	EnableRetries bool   `mapstructure:"enable_retries"`
	RetryInterval int16  `mapstructure:"retry_interval"`
	ServerName    string `mapstructure:"server_name"`
	Email         Email
}

func InitSetting(path string) (Setting, error) {
	if err := validate.IsTomlFile(path, "otherSettingsPath.notifications"); err != nil {
		return Setting{}, err
	}

	setting := settingDefault()

	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return Setting{}, err
	}
	if err := v.Unmarshal(&setting); err != nil {
		return Setting{}, err
	}

	if !setting.Enabled {
		return setting, nil
	}

	if err := setting.Validate(); err != nil {
		return Setting{}, err
	}

	return setting, nil
}

func settingDefault() Setting {
	return Setting{
		Enabled:       false,
		EnableRetries: true,
		RetryInterval: 600,
		ServerName:    "server",
		Email:         defaultEmail(),
	}
}

func (s Setting) Validate() error {
	if s.ServerName == "" {
		return errors.New("server_name is not specified")
	}

	if err := validate.NameWithDot(s.ServerName, "server_name"); err != nil {
		return err
	}

	if err := s.Email.Validate(); err != nil {
		return err
	}

	if s.RetryInterval < 1 {
		return errors.New("retry_interval must be greater than 0")
	}

	return nil
}
