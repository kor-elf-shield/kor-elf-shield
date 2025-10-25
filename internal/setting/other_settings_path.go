package setting

import (
	"kor-elf-shield/internal/daemon/firewall"
	firewallSetting "kor-elf-shield/internal/setting/firewall"
)

type otherSettingsPath struct {
	Firewall string `mapstructure:"firewall"`
}

func otherSettingsPathDefault() *otherSettingsPath {
	return &otherSettingsPath{
		Firewall: "/etc/kor-elf-shield/firewall.toml",
	}
}

func (o *otherSettingsPath) ToFirewallConfig() (firewall.Config, error) {
	setting, err := firewallSetting.InitSetting(o.Firewall)
	if err != nil {
		return firewall.Config{}, err
	}

	return firewall.Config{
		SavesRules:     setting.SavesRules,
		SavesRulesPath: setting.SavesRulesPath,
		TableName:      setting.MetadataNaming.TableName,
	}, nil
}
