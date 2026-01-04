package docker_monitor

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/chain"
	nftChain "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain"
)

type DockerNotSupport struct {
	chains chain.Chains
}

func NewDockerNotSupport() Docker {
	return &DockerNotSupport{
		chains: chain.NewEmptyChains(),
	}
}

func (d *DockerNotSupport) NftReload(_ func(chain string) (nftChain.Chain, error)) error {
	return nil
}

func (d *DockerNotSupport) NftChains() chain.Chains {
	return d.chains
}
