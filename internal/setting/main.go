package setting

import (
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
	"github.com/spf13/viper"
)

var (
	// AppVer is the version of the current build of Kor-Elf-Shield. It is set in main.go from main.Version.
	AppVer string

	// AppBuiltWith represents a human-readable version go runtime build version and build tags. (See main.go formatBuiltWith().)
	AppBuiltWith string

	// AppStartTime program start time
	AppStartTime time.Time
	Config       *setting
)

func init() {
	if AppVer == "" {
		AppVer = "development"
	}
}

func InitSetting(path string) error {
	if err := validate.IsTomlFile(path, "config"); err != nil {
		return err
	}

	// Default config
	Config = settingDefault(path)

	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&Config); err != nil {
		return err
	}

	if err := Config.Validate(); err != nil {
		return err
	}

	return nil
}
