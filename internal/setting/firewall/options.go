package firewall

import (
	"errors"
	"strings"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

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

func (o options) Validate() error {
	if err := o.ValidateSavesRulesPath(); err != nil {
		return err
	}
	return nil
}

func (o options) ValidateSavesRulesPath() error {
	if o.SavesRulesPath == "" {
		return errors.New("saves_rules_path is empty")
	}
	if err := validate.PathFile(o.SavesRulesPath, "saves_rules_path"); err != nil {
		return err
	}
	if !strings.HasSuffix(strings.ToLower(o.SavesRulesPath), "nftables.conf") {
		return errors.New("invalid saves_rules_path. Must be nftables.conf")
	}

	return nil
}
