package firewall

import "github.com/spf13/viper"

type Setting struct {
	SavesRules     bool   `mapstructure:"saves_rules"`
	SavesRulesPath string `mapstructure:"saves_rules_path"`
	MetadataNaming metadataNaming
	Policy         Policy
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
		SavesRules:     false,
		SavesRulesPath: "/etc/nftables.conf",
		MetadataNaming: defaultMetadataNaming(),
		Policy:         defaultPolicy(),
	}
}
