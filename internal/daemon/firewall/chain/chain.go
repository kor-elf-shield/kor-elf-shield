package chain

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type Chain interface {
	AddRule(expr ...string) error
	Clear() error
}

type chain struct {
	nft    nft.NFT
	family family.Type
	table  string
	chain  string
}

func (c *chain) AddRule(expr ...string) error {
	return c.nft.Rule().Add(c.family, c.table, c.chain, expr...)
}

func (c *chain) Clear() error {
	return c.nft.Chain().Clear(c.family, c.table, c.chain)
}
