package daemon

import (
	analyzerConfig "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db"
	firewallConfig "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/config"
)

type DaemonOptions struct {
	PathPidFile    string
	PathSocketFile string
	DataDir        string
	PathNftables   string
	ConfigFirewall firewallConfig.Config
	ConfigAnalyzer analyzerConfig.Config
	Repositories   db.Repositories
}
