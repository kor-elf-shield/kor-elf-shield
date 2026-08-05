package setting

import (
	"errors"

	analyzerConfig "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/blocklist"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor"
	firewallConfig "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/config"
	GuardConfig "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/guard/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/geoip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/notifications"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	logger "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"

	analyzerSetting "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/analyzer"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/blocklists"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/docker"
	firewallSetting "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/firewall"
	geoIPSetting "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/geoip"
	notificationsSetting "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/notifications"

	"github.com/wneessen/go-mail"
)

type otherSettingsPath struct {
	Firewall      string `mapstructure:"firewall"`
	Notifications string `mapstructure:"notifications"`
	Analyzer      string `mapstructure:"analyzer"`
	Docker        string `mapstructure:"docker"`
	Blocklists    string `mapstructure:"blocklists"`
	GeoIP         string `mapstructure:"geoip"`
}

func otherSettingsPathDefault() *otherSettingsPath {
	return &otherSettingsPath{
		Firewall:      "/etc/kor-elf-shield/firewall.toml",
		Notifications: "/etc/kor-elf-shield/notifications.toml",
		Analyzer:      "/etc/kor-elf-shield/analyzer.toml",
		Docker:        "/etc/kor-elf-shield/docker.toml",
		Blocklists:    "/etc/kor-elf-shield/blocklists.toml",
		GeoIP:         "/etc/kor-elf-shield/geoip.toml",
	}
}

func (o *otherSettingsPath) ToFirewallConfig(dockerSupport bool) (firewallConfig.Config, GuardConfig.GuardConfig, error) {
	setting, err := firewallSetting.InitSetting(o.Firewall)
	if err != nil {
		return firewallConfig.Config{}, GuardConfig.GuardConfig{}, err
	}

	configPolicy, err := setting.Policy.ToConfigPolicy()
	if err != nil {
		return firewallConfig.Config{}, GuardConfig.GuardConfig{}, err
	}

	inPorts, outPorts, err := setting.ToPorts()
	if err != nil {
		return firewallConfig.Config{}, GuardConfig.GuardConfig{}, err
	}

	IPs, err := setting.ToIPs()
	if err != nil {
		return firewallConfig.Config{}, GuardConfig.GuardConfig{}, err
	}

	optionClearMode, err := setting.Options.ToClearMode()
	if err != nil {
		return firewallConfig.Config{}, GuardConfig.GuardConfig{}, err
	}

	portKnocking, err := setting.ToConfigPortKnocking()
	if err != nil {
		return firewallConfig.Config{}, GuardConfig.GuardConfig{}, err
	}

	firewall := firewallConfig.Config{
		InPorts:  inPorts,
		OutPorts: outPorts,
		IP4: firewallConfig.ConfigIP4{
			IcmpIn:            setting.IP4.IcmpIn,
			IcmpInRate:        setting.IP4.IcmpInRate,
			IcmpOut:           setting.IP4.IcmpOut,
			IcmpOutRate:       setting.IP4.IcmpOutRate,
			IcmpTimestampDrop: setting.IP4.IcmpTimestampDrop,
			InIPs:             IPs.InIP4,
			OutIPs:            IPs.OutIP4,
		},
		IP6: firewallConfig.ConfigIP6{
			Enable:     setting.IP6.Enable,
			IcmpStrict: setting.IP6.IcmpStrict,
			InIPs:      IPs.InIP6,
			OutIPs:     IPs.OutIP6,
		},
		Options: firewallConfig.ConfigOptions{
			Cache:          setting.Options.Cache,
			ClearMode:      optionClearMode,
			SavesRules:     setting.Options.SavesRules,
			SavesRulesPath: setting.Options.SavesRulesPath,
			DnsStrict:      setting.Options.DnsStrict,
			DnsStrictNs:    setting.Options.DnsStrictNs,
			PacketFilter:   setting.Options.PacketFilter,
			DockerSupport:  dockerSupport,
		},
		MetadataNaming: firewallConfig.ConfigMetadata{
			TableName:        setting.MetadataNaming.TableName,
			ChainInputName:   setting.MetadataNaming.ChainInputName,
			ChainOutputName:  setting.MetadataNaming.ChainOutputName,
			ChainForwardName: setting.MetadataNaming.ChainForwardName,
		},
		Policy:       configPolicy,
		PortKnocking: portKnocking,
	}

	rulesGuard := setting.RulesGuard.ToGuardConfig()

	return firewall, rulesGuard, nil
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
		Enabled:       setting.Enabled,
		EnableRetries: setting.EnableRetries,
		RetryInterval: uint16(setting.RetryInterval),
		ServerName:    setting.ServerName,
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

func (o *otherSettingsPath) ToAnalyzerConfig(binaryLocations *binaryLocations) (analyzerConfig.Config, error) {
	if binaryLocations.Journalctl == "" {
		return analyzerConfig.Config{}, errors.New(i18n.Lang.T("parameter is not specified", map[string]any{
			"Parameter": "binaryLocations.journalctl",
		}))
	}

	setting, err := analyzerSetting.InitSetting(o.Analyzer)
	if err != nil {
		return analyzerConfig.Config{}, err
	}

	if err := setting.Validate(); err != nil {
		return analyzerConfig.Config{}, err
	}

	binPath := analyzerConfig.BinPath{
		Journalctl: binaryLocations.Journalctl,
	}

	sources, err := setting.ToSources()
	if err != nil {
		return analyzerConfig.Config{}, err
	}

	return analyzerConfig.Config{
		BinPath: binPath,
		Sources: sources,
	}, nil
}

func (o *otherSettingsPath) ToDockerConfig(binaryLocations *binaryLocations) (config docker_monitor.Config, dockerSupport bool, err error) {
	if binaryLocations.Docker == "" {
		return docker_monitor.Config{}, false, errors.New(i18n.Lang.T("parameter is not specified", map[string]any{
			"Parameter": "binaryLocations.docker",
		}))
	}

	setting, err := docker.InitSetting(o.Docker)
	if err != nil {
		return docker_monitor.Config{}, false, err
	}

	if err := setting.Validate(); err != nil {
		return docker_monitor.Config{}, false, err
	}

	ruleStrategy, err := setting.ToRuleStrategy()
	if err != nil {
		return docker_monitor.Config{}, false, err
	}

	return docker_monitor.Config{
		Path:         binaryLocations.Docker,
		RuleStrategy: ruleStrategy,
	}, setting.Enabled, nil
}

func (o *otherSettingsPath) ToBlocklistConfig(logger logger.Logger) (sources []*blocklist.SourceConfig, blocklistSupport bool, err error) {
	setting, err := blocklists.InitSetting(o.Blocklists)
	if err != nil {
		return []*blocklist.SourceConfig{}, false, err
	}

	sources = setting.ToSources(logger)

	if setting.Enabled && len(sources) == 0 {
		return []*blocklist.SourceConfig{}, false, errors.New(i18n.Lang.T("blocklist sources are empty"))
	}

	return sources, setting.Enabled, nil
}

func (o *otherSettingsPath) ToConfig(dataDir string, logger logger.Logger) (geoIPService *geoip.Config, geoIPSupport bool, err error) {
	setting, err := geoIPSetting.InitSetting(o.GeoIP)
	if err != nil {
		return &geoip.Config{}, false, err
	}

	if !setting.Enabled {
		return nil, false, nil
	}

	geoIPService, err = setting.ToConfig(dataDir, logger)
	if err != nil {
		return geoIPService, false, err
	}

	return geoIPService, setting.Enabled, nil
}

func (o *otherSettingsPath) ListPathFiles() map[string]string {
	return map[string]string{
		"firewall":      o.Firewall,
		"notifications": o.Notifications,
		"analyzer":      o.Analyzer,
		"docker":        o.Docker,
		"blocklists":    o.Blocklists,
		"geoip":         o.GeoIP,
	}
}
