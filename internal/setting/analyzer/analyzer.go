package analyzer

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
	"github.com/spf13/viper"
)

type Setting struct {
	Login Login
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
		Login: defaultLogin(),
	}
}

func (s Setting) ToSources() ([]*config.Source, error) {
	var sources []*config.Source

	loginSources, err := s.Login.ToSources()
	if err != nil {
		return sources, err
	}
	sources = append(sources, loginSources...)

	return sources, nil
}

func (s Setting) Validate() error {
	if err := s.Login.Validate(); err != nil {
		return err
	}

	return nil
}
