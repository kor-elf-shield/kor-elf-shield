package firewall

import (
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	nftFirewall "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
)

type NFTDocker interface {
	Chains() NFTDockerChains
	NFT() nftFirewall.NFT
}

type NFTDockerChains interface {
	List() []chain.Docker

	ForwardFilter() chain.Docker
	ForwardBridge() chain.Docker
	ForwardCT() chain.Docker

	PreroutingFilter() chain.Docker
	DockerFilter() chain.Docker
	DockerFilterFirst() chain.Docker
	DockerFilterSecond() chain.Docker

	DockerNat() chain.Docker
	PostroutingNat() chain.Docker
}

type nftDocker struct {
	chains NFTDockerChains
	nft    nftFirewall.NFT
}

func NewNFT(nft nftFirewall.NFT, chains NFTDockerChains) NFTDocker {
	return &nftDocker{
		chains: chains,
		nft:    nft,
	}
}

func (n *nftDocker) NFT() nftFirewall.NFT {
	return n.nft
}

func (n *nftDocker) Chains() NFTDockerChains {
	return n.chains
}

type nftDockerChains struct {
	forwardFilter chain.Docker
	forwardBridge chain.Docker
	forwardCT     chain.Docker

	preroutingFilter   chain.Docker
	dockerFilter       chain.Docker
	dockerFilterFirst  chain.Docker
	dockerFilterSecond chain.Docker

	dockerNat      chain.Docker
	postroutingNat chain.Docker
}

func NewNFTChains(nft nftFirewall.NFT, family family.Type, table string) NFTDockerChains {
	return &nftDockerChains{
		forwardFilter: chain.NewDocker(nft.NFT(), family, table, "docker_forward_filter"),
		forwardBridge: chain.NewDocker(nft.NFT(), family, table, "docker_forward_bridge"),
		forwardCT:     chain.NewDocker(nft.NFT(), family, table, "docker_forward_ct"),

		preroutingFilter:   chain.NewDocker(nft.NFT(), family, table, "docker_prerouting_filter"),
		dockerFilter:       chain.NewDocker(nft.NFT(), family, table, "docker_filter"),
		dockerFilterFirst:  chain.NewDocker(nft.NFT(), family, table, "docker_filter_first"),
		dockerFilterSecond: chain.NewDocker(nft.NFT(), family, table, "docker_filter_second"),

		dockerNat:      chain.NewDocker(nft.NFT(), family, table, "docker_nat"),
		postroutingNat: chain.NewDocker(nft.NFT(), family, table, "docker_postrouting_nat"),
	}
}

func (n *nftDockerChains) ForwardFilter() chain.Docker {
	return n.forwardFilter
}

func (n *nftDockerChains) ForwardBridge() chain.Docker {
	return n.forwardBridge
}

func (n *nftDockerChains) ForwardCT() chain.Docker {
	return n.forwardCT
}

func (n *nftDockerChains) PreroutingFilter() chain.Docker {
	return n.preroutingFilter
}

func (n *nftDockerChains) DockerFilter() chain.Docker {
	return n.dockerFilter
}

func (n *nftDockerChains) DockerFilterFirst() chain.Docker {
	return n.dockerFilterFirst
}

func (n *nftDockerChains) DockerFilterSecond() chain.Docker {
	return n.dockerFilterSecond
}

func (n *nftDockerChains) DockerNat() chain.Docker {
	return n.dockerNat
}

func (n *nftDockerChains) PostroutingNat() chain.Docker {
	return n.postroutingNat
}

func (n *nftDockerChains) List() []chain.Docker {
	return []chain.Docker{
		n.forwardFilter,
		n.forwardBridge,
		n.forwardCT,
		n.preroutingFilter,
		n.dockerFilter,
		n.dockerFilterFirst,
		n.dockerFilterSecond,
		n.dockerNat,
		n.postroutingNat,
	}
}
