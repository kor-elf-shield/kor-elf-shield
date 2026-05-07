package table

import (
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/block"
)

type Table interface {
	Clear() error
	DockerChains() firewall.NFTDockerChains
	BlockList() BlockList
}

type BlockList interface {
	ListIP() block.ListIP
	ListIPWithPort() block.ListIPWithPort
	Blocks() map[string]block.Blocklist
}

type table struct {
	nft    nft.NFT
	family family.Type
	name   string

	dockerChains firewall.NFTDockerChains
	blockList    BlockList
}

func New(
	nft nft.NFT, family family.Type, name string,
	blockList BlockList, dockerChains firewall.NFTDockerChains,
) Table {
	return &table{
		nft:    nft,
		family: family,
		name:   name,

		dockerChains: dockerChains,
		blockList:    blockList,
	}
}

func (t *table) Clear() error {
	// clear does not clean completely
	return t.nft.NFT().Table().Delete(t.family, t.name)
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
	blocks         map[string]block.Blocklist
}

func NewBlockList(listIP block.ListIP, listIPWithPort block.ListIPWithPort, blocks map[string]block.Blocklist) BlockList {
	return &blockList{
		listIP:         listIP,
		listIPWithPort: listIPWithPort,
		blocks:         blocks,
	}
}

func (b *blockList) ListIP() block.ListIP {
	return b.listIP
}

func (b *blockList) ListIPWithPort() block.ListIPWithPort {
	return b.listIPWithPort
}

func (b *blockList) Blocks() map[string]block.Blocklist {
	return b.blocks
}
