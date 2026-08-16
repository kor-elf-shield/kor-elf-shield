package firewall

import (
	"fmt"

	GuardConfig "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/guard/config"
)

type RulesGuard struct {
	Enabled       bool  `mapstructure:"enabled"`
	Notifications bool  `mapstructure:"notifications"`
	Recovery      bool  `mapstructure:"recovery"`
	Interval      int32 `mapstructure:"interval"`
}

func defaultRulesGuard() RulesGuard {
	return RulesGuard{
		Enabled:       true,
		Notifications: true,
		Recovery:      true,
		Interval:      3600,
	}
}

func (r *RulesGuard) Validate() error {
	if r.Interval < 60 {
		return fmt.Errorf("interval must be greater than 60")
	}

	return nil
}

func (r *RulesGuard) ToGuardConfig() GuardConfig.GuardConfig {
	return GuardConfig.GuardConfig{
		Enable:        r.Enabled,
		Notifications: r.Notifications,
		Recovery:      r.Recovery,
		Interval:      uint32(r.Interval),
	}
}
