package chain

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type LocalInput interface {
	AddRule(expr ...string) error
	AddRuleIn(AddRuleFunc func(expr ...string) error) error
}

type localInput struct {
	nft    nft.NFT
	family family.Type
	table  string
	chain  string
}

func newLocalInput(nft nft.NFT, family family.Type, table string) (LocalInput, error) {
	chain := "local-input"
	if err := nft.Chain().Add(family, table, chain, nftChain.TypeNone); err != nil {
		return nil, err
	}

	return &localInput{
		nft:    nft,
		family: family,
		table:  table,
		chain:  chain,
	}, nil
}

func (c *localInput) AddRule(expr ...string) error {
	return c.nft.Rule().Add(c.family, c.table, c.chain, expr...)
}

func (f *localInput) AddRuleIn(AddRuleFunc func(expr ...string) error) error {
	return AddRuleFunc("iifname != \"lo\" counter jump " + f.chain)
}
