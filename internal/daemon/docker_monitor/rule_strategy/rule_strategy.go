package rule_strategy

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
)

type Strategy interface {
	Reload(nftDocker firewall.NFTDocker) error
	Chains() firewall.NFTDockerChains
	Event(event *client.Event)
}
