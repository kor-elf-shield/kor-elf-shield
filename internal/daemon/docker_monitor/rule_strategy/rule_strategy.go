package rule_strategy

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
	nftChain "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain"
)

type Strategy interface {
	Reload(newNoneChain func(chain string) (nftChain.Chain, error)) error
	Chains() chain.Chains
	Event(event *client.Event)
}
