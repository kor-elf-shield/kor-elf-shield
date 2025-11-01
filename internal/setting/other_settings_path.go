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

	inPorts, outPorts, err := setting.ToPorts()
	if err != nil {
		return firewall.Config{}, err
	}

	IPs, err := setting.ToIPs()
	if err != nil {
		return firewall.Config{}, err
	}

	return firewall.Config{
		InPorts:  inPorts,
		OutPorts: outPorts,
		IP4: firewall.ConfigIP4{
			IcmpIn:            setting.IP4.IcmpIn,
			IcmpInRate:        setting.IP4.IcmpInRate,
			IcmpOut:           setting.IP4.IcmpOut,
			IcmpOutRate:       setting.IP4.IcmpOutRate,
			IcmpTimestampDrop: setting.IP4.IcmpTimestampDrop,
			InIPs:             IPs.InIP4,
			OutIPs:            IPs.OutIP4,
		},
		IP6: firewall.ConfigIP6{
			Enable:     setting.IP6.Enable,
			IcmpStrict: setting.IP6.IcmpStrict,
			InIPs:      IPs.InIP6,
			OutIPs:     IPs.OutIP6,
		},
		Options: firewall.ConfigOptions{
			SavesRules:     setting.Options.SavesRules,
			SavesRulesPath: setting.Options.SavesRulesPath,
			DnsStrict:      setting.Options.DnsStrict,
			DnsStrictNs:    setting.Options.DnsStrictNs,
			PacketFilter:   setting.Options.PacketFilter,
		},
		MetadataNaming: firewall.ConfigMetadata{
			TableName:        setting.MetadataNaming.TableName,
			ChainInputName:   setting.MetadataNaming.ChainInputName,
			ChainOutputName:  setting.MetadataNaming.ChainOutputName,
			ChainForwardName: setting.MetadataNaming.ChainForwardName,
		},
		Policy: configPolicy,
	}, nil
}
