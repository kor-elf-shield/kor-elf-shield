package blocklists

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/blocklist"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
	"github.com/spf13/viper"
)

type Setting struct {
	Enabled bool `mapstructure:"enabled"`
	Sources []Sources
}

func InitSetting(path string) (Setting, error) {
	if err := validate.IsTomlFile(path, "otherSettingsPath.blocklists"); err != nil {
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
		Sources: []Sources{},
	}
}

func (b *Setting) ToSources(logger log.Logger) []*blocklist.SourceConfig {
	var sources []*blocklist.SourceConfig
	if !b.Enabled {
		return sources
	}

	sourceNames := make(map[string]string)

	for _, source := range b.Sources {
		if !source.Enabled {
			continue
		}

		if _, ok := sourceNames[source.Name]; ok {
			logger.Warn(fmt.Sprintf("duplicate source name: %s", source.Name))
			continue
		}
		sourceNames[source.Name] = source.Name

		sourceConfig, err := source.ToSourceConfig()
		if err != nil {
			logger.Warn(fmt.Sprintf("failed to convert source: %s", err))
			continue
		}

		sources = append(sources, sourceConfig)
	}

	return sources
}
