package daemon

import "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"

type DaemonOptions struct {
	PathPidFile    string
	PathSocketFile string
	PathNftables   string
	ConfigFirewall firewall.Config
}
