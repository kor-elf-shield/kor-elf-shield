package daemon

import (
	analyzerConfig "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db"
	firewallConfig "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/config"
	GuardConfig "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/guard/config"
)

type DaemonOptions struct {
	PathPidFile         string
	PathSocketFile      string
	DataDir             string
	PathNftables        string
	ConfigFirewall      firewallConfig.Config
	ConfigFirewallGuard GuardConfig.GuardConfig
	ConfigAnalyzer      analyzerConfig.Config
	Repositories        db.Repositories
}
