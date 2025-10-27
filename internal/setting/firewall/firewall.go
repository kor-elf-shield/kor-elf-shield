package firewall

import "github.com/spf13/viper"

type Setting struct {
	Options        options
	MetadataNaming metadataNaming
	Policy         policy
}

func InitSetting(path string) (Setting, error) {
	setting := settingDefault()

	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		return Setting{}, err
	}

	err = v.Unmarshal(&setting)
	if err != nil {
		return Setting{}, err
	}

	return setting, nil
}

func settingDefault() Setting {
	return Setting{
		Options:        defaultOptions(),
		MetadataNaming: defaultMetadataNaming(),
		Policy:         defaultPolicy(),
	}
}
