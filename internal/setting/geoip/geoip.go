package geoip

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/geoip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
	"github.com/spf13/viper"
)

type Setting struct {
	Enabled bool   `mapstructure:"enabled"`
	Service string `mapstructure:"service"`
	Maxmind Maxmind
}

func InitSetting(path string) (Setting, error) {
	if err := validate.IsTomlFile(path, "otherSettingsPath.geoip"); err != nil {
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

	return setting, nil
}

func settingDefault() Setting {
	return Setting{
		Enabled: false,
		Service: "maxmind",
		Maxmind: Maxmind{},
	}
}

func (s *Setting) ToConfig(dataDir string, logger log.Logger) (*geoip.Config, error) {
	if s.Service == "maxmind" {
		return s.Maxmind.ToConfig(dataDir, logger)
	}

	return nil, fmt.Errorf("unknown service: %s", s.Service)
}
