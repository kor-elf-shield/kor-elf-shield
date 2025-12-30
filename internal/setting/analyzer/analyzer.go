package analyzer

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
	"github.com/spf13/viper"
)

type Setting struct {
	Enabled   bool `mapstructure:"enabled"`
	Notify    bool `mapstructure:"notify"`
	SSHEnable bool `mapstructure:"ssh_enable"`
	SSHNotify bool `mapstructure:"ssh_notify"`
}

func InitSetting(path string) (Setting, error) {
	if err := validate.IsTomlFile(path, "otherSettingsPath.analyzer"); err != nil {
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

	return setting, nil
}

func settingDefault() Setting {
	return Setting{
		Enabled:   true,
		Notify:    true,
		SSHEnable: true,
		SSHNotify: true,
	}
}

func (s Setting) Validate() error {

	return nil
}
