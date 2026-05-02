package chain

import (
	"strconv"
	"strings"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/block"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/rule"
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
	AddRuleIn(AddRuleFunc rule.AddFunc) error
}

type portKnocking struct {
	chain Chain
	sets  block.Sets
}

func NewBatchPortKnocking(builder nft.BatchBuilder, family family.Type, table string, chain string) (PortKnocking, error) {
	batchChain, err := NewBatchChain(builder, family, table, chain)
	if err != nil {
		return nil, err
	}

	return &portKnocking{
		chain: batchChain,
		sets:  block.NewBatchSet(builder, family, table),
	}, nil
}

func (k *portKnocking) AddRuleIn(AddRuleFunc rule.AddFunc) error {
	return k.chain.AddRuleIn(AddRuleFunc)
}

func (k *portKnocking) AddFirstStageRule(
	name string,
	ipVersion ip.Version,
	l4Port types.L4Port,
	timeout uint32,
	action types.KnockAction,
) error {
	if err := k.newPortKnocking(name, ipVersion, timeout); err != nil {
		return err
	}

	expr := []string{
		l4Port.ProtocolString(), "dport", l4Port.NumberString(), "add", "@" + name,
		"{", ipVersion.ToNft(), "saddr timeout", strconv.Itoa(int(timeout)) + "s", "}", action.String(),
	}
	return k.chain.AddRule(expr...)
}

func (k *portKnocking) AddNextStageRule(
	prevName, nextName string,
	ipVersion ip.Version,
	l4Port types.L4Port,
	timeout uint32,
	action types.KnockAction,
) error {
	if err := k.newPortKnocking(nextName, ipVersion, timeout); err != nil {
		return err
	}

	expr := []string{
		ipVersion.ToNft(), "saddr", "@" + prevName,
		l4Port.ProtocolString(), "dport", l4Port.NumberString(), "add", "@" + nextName,
		"{", ipVersion.ToNft(), "saddr}", action.String(),
	}
	return k.chain.AddRule(expr...)
}

func (k *portKnocking) newPortKnocking(name string, ipVersion ip.Version, timeout uint32) error {
	params := []string{"type", ipVersion.ToNftForSet() + ";", "flags timeout; timeout", strconv.Itoa(int(timeout)) + "s;"}
	return k.sets.Add(name, strings.Join(params, " "))
}
