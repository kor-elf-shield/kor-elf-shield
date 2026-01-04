package setting

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	analyzerSetting "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/analyzer"
	firewallSetting "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/firewall"
	notificationsSetting "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/notifications"
	"github.com/wneessen/go-mail"
)

type otherSettingsPath struct {
	Firewall      string `mapstructure:"firewall"`
	Notifications string `mapstructure:"notifications"`
	Analyzer      string `mapstructure:"analyzer"`
}

func otherSettingsPathDefault() *otherSettingsPath {
	return &otherSettingsPath{
		Firewall:      "/etc/kor-elf-shield/firewall.toml",
		Notifications: "/etc/kor-elf-shield/notifications.toml",
		Analyzer:      "/etc/kor-elf-shield/analyzer.toml",
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

	optionClearMode, err := setting.Options.ToClearMode()
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
			ClearMode:      optionClearMode,
			SavesRules:     setting.Options.SavesRules,
			SavesRulesPath: setting.Options.SavesRulesPath,
			DnsStrict:      setting.Options.DnsStrict,
			DnsStrictNs:    setting.Options.DnsStrictNs,
			PacketFilter:   setting.Options.PacketFilter,
			DockerSupport:  setting.Options.DockerSupport,
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

func (o *otherSettingsPath) ToNotificationsConfig() (notifications.Config, error) {
	setting, err := notificationsSetting.InitSetting(o.Notifications)
	if err != nil {
		return notifications.Config{}, err
	}

	authType := mail.SMTPAuthPlain
	tls := notifications.TLS{}
	if setting.Enabled {
		authType, err = notificationsSetting.ParseAuthType(setting.Email.AuthType)
		if err != nil {
			return notifications.Config{}, err
		}

		tls, err = setting.Email.ToTLSConfig()
		if err != nil {
			return notifications.Config{}, err
		}
	}

	return notifications.Config{
		Enabled:    setting.Enabled,
		ServerName: setting.ServerName,
		Email: notifications.Email{
			Host:     setting.Email.Host,
			Port:     uint(setting.Email.Port),
			Username: setting.Email.Username,
			Password: setting.Email.Password,
			AuthType: authType,
			TLS:      tls,
			From:     setting.Email.From,
			To:       setting.Email.To,
		},
	}, nil
}

func (o *otherSettingsPath) ToAnalyzerConfig(binaryLocations *binaryLocations) (config.Config, error) {
	if binaryLocations.Journalctl == "" {
		return config.Config{}, errors.New(i18n.Lang.T("parameter is not specified", map[string]any{
			"Parameter": "binaryLocations.journalctl",
		}))
	}

	setting, err := analyzerSetting.InitSetting(o.Analyzer)
	if err != nil {
		return config.Config{}, err
	}

	if err := setting.Validate(); err != nil {
		return config.Config{}, err
	}

	binPath := config.BinPath{
		Journalctl: binaryLocations.Journalctl,
	}

	login := config.Login{
		Enabled: setting.Login.Enabled,
		Notify:  setting.Login.Notify,
		SSH: config.LoginSSH{
			Enabled: setting.Login.SSHEnable,
			Notify:  setting.Login.SSHNotify,
		},
	}

	return config.Config{
		BinPath: binPath,
		Login:   login,
	}, nil
}
