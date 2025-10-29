package firewall

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
	"github.com/spf13/viper"
)

type Setting struct {
	IP4            ip4
	IP6            ip6
	Options        options
	MetadataNaming metadataNaming
	Policy         policy
}

func InitSetting(path string) (Setting, error) {
	if err := validate.IsTomlFile(path, "otherSettingsPath.firewall"); err != nil {
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
	if err := setting.Validate(); err != nil {
		return Setting{}, err
	}

	return setting, nil
}

func settingDefault() Setting {
	return Setting{
		IP4:            defaultIp4(),
		IP6:            defaultIp6(),
		Options:        defaultOptions(),
		MetadataNaming: defaultMetadataNaming(),
		Policy:         defaultPolicy(),
	}
}

func (s Setting) Validate() error {
	if err := s.IP4.Validate(); err != nil {
		return err
	}
	if err := s.IP6.Validate(); err != nil {
		return err
	}
	if err := s.MetadataNaming.Validate(); err != nil {
		return err
	}
	if err := s.Policy.Validate(); err != nil {
		return err
	}
	if err := s.Options.Validate(); err != nil {
		return err
	}
	return nil
}
