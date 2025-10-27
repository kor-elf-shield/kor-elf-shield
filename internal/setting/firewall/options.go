package firewall

type options struct {
	SavesRules     bool   `mapstructure:"saves_rules"`
	SavesRulesPath string `mapstructure:"saves_rules_path"`
	DnsStrict      bool   `mapstructure:"dns_strict"`
	DnsStrictNs    bool   `mapstructure:"dns_strict_ns"`
	PacketFilter   bool   `mapstructure:"packet_filter"`
}

func defaultOptions() options {
	return options{
		SavesRules:     false,
		SavesRulesPath: "/etc/nftables.conf",
		DnsStrict:      false,
		DnsStrictNs:    false,
		PacketFilter:   true,
	}
}
