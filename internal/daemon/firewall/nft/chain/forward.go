package chain

import (
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

func NewBatchForward(builder nft.BatchBuilder, family family.Type, table string, chain string, defaultAllow bool, priority int) (Chain, error) {
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
