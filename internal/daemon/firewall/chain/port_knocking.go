package chain

import (
	"strconv"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain/block"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/ip"
)

type PortKnocking interface {
	AddFirstStageRule(
		name string,
		ipVersion ip.Version,
		l4Port types.L4Port,
		timeout uint32,
		action types.KnockAction,
	) error
	AddNextStageRule(
		prevName, nextName string,
		ipVersion ip.Version,
		l4Port types.L4Port,
		timeout uint32,
		action types.KnockAction,
	) error
	AddRuleIn(AddRuleFunc func(expr ...string) error) error
}

type portKnocking struct {
	nft    nft.NFT
	family family.Type
	table  string
	chain  string
}

func newPortKnocking(nft nft.NFT, family family.Type, table string, chain string) (PortKnocking, error) {
	if err := nft.Chain().Add(family, table, chain, nftChain.TypeNone); err != nil {
		return nil, err
	}

	return &portKnocking{
		nft:    nft,
		family: family,
		table:  table,
		chain:  chain,
	}, nil
}

func (k *portKnocking) AddRuleIn(AddRuleFunc func(expr ...string) error) error {
	return AddRuleFunc("iifname != \"lo\" counter jump " + k.chain)
}

func (k *portKnocking) AddFirstStageRule(
	name string,
	ipVersion ip.Version,
	l4Port types.L4Port,
	timeout uint32,
	action types.KnockAction,
) error {
	if err := block.NewPortKnocking(k.nft, k.family, k.table, name, ipVersion, timeout); err != nil {
		return err
	}

	expr := []string{
		l4Port.ProtocolString(), "dport", l4Port.NumberString(), "add", "@" + name,
		"{", ipVersion.ToNft(), "saddr timeout", strconv.Itoa(int(timeout)) + "s", "}", action.String(),
	}
	return k.nft.Rule().Add(k.family, k.table, k.chain, expr...)
}

func (k *portKnocking) AddNextStageRule(
	prevName, nextName string,
	ipVersion ip.Version,
	l4Port types.L4Port,
	timeout uint32,
	action types.KnockAction,
) error {
	if err := block.NewPortKnocking(k.nft, k.family, k.table, nextName, ipVersion, timeout); err != nil {
		return err
	}

	expr := []string{
		ipVersion.ToNft(), "saddr", "@" + prevName,
		l4Port.ProtocolString(), "dport", l4Port.NumberString(), "add", "@" + nextName,
		"{", ipVersion.ToNft(), "saddr}", action.String(),
	}
	return k.nft.Rule().Add(k.family, k.table, k.chain, expr...)
}
