package setting

import (
	"time"

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

	// Default config
	Config = settingDefault()
}

func InitSetting(path string) error {
	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&Config)
	if err != nil {
		return err
	}

	return nil
}
