package chain

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type Forward interface {
	AddRule(expr ...string) error
}

type forward struct {
	nft    nft.NFT
	family family.Type
	table  string
	chain  string
}

func newForward(nft nft.NFT, family family.Type, table string, chain string, defaultAllow bool, priority int) (Forward, error) {
	policy := nftChain.PolicyDrop
	if defaultAllow {
		policy = nftChain.PolicyAccept
	}

	baseChain := nftChain.BaseChainOptions{
		Type:     nftChain.TypeFilter,
		Hook:     nftChain.HookForward,
		Priority: int32(priority),
		Policy:   policy,
		Device:   "",
	}

	if err := nft.Chain().Add(family, table, chain, baseChain); err != nil {
		return nil, err
	}

	return &forward{
		nft:    nft,
		family: family,
		table:  table,
		chain:  chain,
	}, nil
}

func (c *forward) AddRule(expr ...string) error {
	return c.nft.Rule().Add(c.family, c.table, c.chain, expr...)
}
