package daemon

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
)

type DaemonOptions struct {
	PathPidFile    string
	PathSocketFile string
	DataDir        string
	PathNftables   string
	ConfigFirewall firewall.Config
	ConfigAnalyzer config.Config
}
