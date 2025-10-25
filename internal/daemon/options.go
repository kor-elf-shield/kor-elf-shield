package daemon

import "kor-elf-shield/internal/daemon/firewall"

type DaemonOptions struct {
	PathPidFile    string
	PathNftables   string
	ConfigFirewall firewall.Config
}
