package docker

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
	"github.com/spf13/viper"
)

type Setting struct {
	Enabled bool `mapstructure:"enabled"`
}

func InitSetting(path string) (Setting, error) {
	if err := validate.IsTomlFile(path, "otherSettingsPath.docker"); err != nil {
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
		Enabled: false,
	}
}

func (s Setting) Validate() error {

	return nil
}
