package chain

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type LocalOutput interface {
	AddRule(expr ...string) error
	AddRuleOut(AddRuleFunc func(expr ...string) error) error
}

type localOutput struct {
	nft    nft.NFT
	family family.Type
	table  string
	chain  string
}

func newLocalOutput(nft nft.NFT, family family.Type, table string) (LocalOutput, error) {
	chain := "local-output"
	if err := nft.Chain().Add(family, table, chain, nftChain.TypeNone); err != nil {
		return nil, err
	}

	return &localOutput{
		nft:    nft,
		family: family,
		table:  table,
		chain:  chain,
	}, nil
}

func (l *localOutput) AddRule(expr ...string) error {
	return l.nft.Rule().Add(l.family, l.table, l.chain, expr...)
}

func (l *localOutput) AddRuleOut(AddRuleFunc func(expr ...string) error) error {
	return AddRuleFunc("oifname != \"lo\" counter jump " + l.chain)
}
