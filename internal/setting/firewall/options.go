package firewall

type options struct {
	SavesRules     bool   `mapstructure:"saves_rules"`
	SavesRulesPath string `mapstructure:"saves_rules_path"`
}

func defaultOptions() options {
	return options{
		SavesRules:     false,
		SavesRulesPath: "/etc/nftables.conf",
	}
}
