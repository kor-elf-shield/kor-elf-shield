package chain

import (
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type BeforeLocalInput interface {
	AddRule(expr ...string) error
	AddRuleIn(AddRuleFunc func(expr ...string) error) error
}

type beforeLocalInput struct {
	nft    nft.NFT
	family family.Type
	table  string
	chain  string
}

func newBeforeLocalInput(nft nft.NFT, family family.Type, table string) (LocalInput, error) {
	chain := "before-local-input"
	if err := nft.Chain().Add(family, table, chain, nftChain.TypeNone); err != nil {
		return nil, err
	}

	return &beforeLocalInput{
		nft:    nft,
		family: family,
		table:  table,
		chain:  chain,
	}, nil
}

func (l *beforeLocalInput) AddRule(expr ...string) error {
	return l.nft.Rule().Add(l.family, l.table, l.chain, expr...)
}

func (l *beforeLocalInput) AddRuleIn(AddRuleFunc func(expr ...string) error) error {
	return AddRuleFunc("iifname != \"lo\" counter jump " + l.chain)
}
