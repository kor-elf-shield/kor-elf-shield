package setting

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
	firewallSetting "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/firewall"
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

	configPolicy, err := setting.Policy.ToConfigPolicy()
	if err != nil {
		return firewall.Config{}, err
	}

	return firewall.Config{
		SavesRules:     setting.SavesRules,
		SavesRulesPath: setting.SavesRulesPath,
		MetadataNaming: firewall.ConfigMetadata{
			TableName:        setting.MetadataNaming.TableName,
			ChainInputName:   setting.MetadataNaming.ChainInputName,
			ChainOutputName:  setting.MetadataNaming.ChainOutputName,
			ChainForwardName: setting.MetadataNaming.ChainForwardName,
		},
		Policy: configPolicy,
	}, nil
}
