package table

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/block"
)

type Table interface {
	DockerChains() firewall.NFTDockerChains
	BlockList() BlockList
}

type BlockList interface {
	ListIP() block.ListIP
	ListIPWithPort() block.ListIPWithPort
}

type table struct {
	dockerChains firewall.NFTDockerChains
	blockList    BlockList
}

func New(blockList BlockList, dockerChains firewall.NFTDockerChains) Table {
	return &table{
		dockerChains: dockerChains,
		blockList:    blockList,
	}
}

func (t *table) DockerChains() firewall.NFTDockerChains {
	return t.dockerChains
}

func (t *table) BlockList() BlockList {
	return t.blockList
}

type blockList struct {
	listIP         block.ListIP
	listIPWithPort block.ListIPWithPort
}

func NewBlockList(listIP block.ListIP, listIPWithPort block.ListIPWithPort) BlockList {
	return &blockList{
		listIP:         listIP,
		listIPWithPort: listIPWithPort,
	}
}

func (b *blockList) ListIP() block.ListIP {
	return b.listIP
}

func (b *blockList) ListIPWithPort() block.ListIPWithPort {
	return b.listIPWithPort
}
