package chain

import (
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/rule"
)

type Chain interface {
	AddRule(expr ...string) error
	AddRuleIn(AddRuleFunc rule.AddFunc) error
	AddRuleOut(AddRuleFunc rule.AddFunc) error
}

type batchChain struct {
	builder nft.BatchBuilder
	family  family.Type
	table   string
	chain   string
}

func NewBatchChain(builder nft.BatchBuilder, family family.Type, table string, chain string) (Chain, error) {
	if err := builder.Chain().Add(family, table, chain, nftChain.TypeNone); err != nil {
		return nil, err
	}

	return &batchChain{
		builder: builder,
		family:  family,
		table:   table,
		chain:   chain,
	}, nil
}

func NewBatchChainWithOptions(builder nft.BatchBuilder, family family.Type, table string, chain string, baseChain nftChain.ChainOptions) (Chain, error) {
	if err := builder.Chain().Add(family, table, chain, baseChain); err != nil {
		return nil, err
	}

	return &batchChain{
		builder: builder,
		family:  family,
		table:   table,
		chain:   chain,
	}, nil
}

func (b *batchChain) AddRule(expr ...string) error {
	return b.builder.Rule().Add(b.family, b.table, b.chain, expr...)
}

func (b *batchChain) AddRuleIn(AddRuleFunc rule.AddFunc) error {
	return AddRuleFunc("iifname != \"lo\" counter jump " + b.chain)
}

func (b *batchChain) AddRuleOut(AddRuleFunc rule.AddFunc) error {
	return AddRuleFunc("oifname != \"lo\" counter jump " + b.chain)
}
