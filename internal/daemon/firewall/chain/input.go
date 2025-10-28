package chain

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type Input interface {
	AddRule(expr ...string) error
}

type input struct {
	nft    nft.NFT
	family family.Type
	table  string
	chain  string
}

func newInput(nft nft.NFT, family family.Type, table string, chain string, defaultAllow bool) (Input, error) {
	policy := nftChain.PolicyDrop
	if defaultAllow {
		policy = nftChain.PolicyAccept
	}

	baseChain := nftChain.BaseChainOptions{
		Type:     nftChain.TypeFilter,
		Hook:     nftChain.HookInput,
		Priority: 0,
		Policy:   policy,
		Device:   "",
	}

	if err := nft.Chain().Add(family, table, chain, baseChain); err != nil {
		return nil, err
	}

	return &input{
		nft:    nft,
		family: family,
		table:  table,
		chain:  chain,
	}, nil
}

func (c *input) AddRule(expr ...string) error {
	return c.nft.Rule().Add(c.family, c.table, c.chain, expr...)
}
