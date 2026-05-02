package chain

import "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/rule"

type PacketFilter interface {
	AddRuleIn(AddRuleFunc rule.AddFunc) error
	AddRuleOut(AddRuleFunc rule.AddFunc) error
}

type packetFilter struct {
	chainName string
}

// NewPacketFilter drop-out-of-order packets and packets in an INVALID state in nftables connection tracking.
func NewPacketFilter(chainName string) PacketFilter {
	return &packetFilter{
		chainName: chainName,
	}
}

func (pf *packetFilter) AddRuleIn(addRuleFunc rule.AddFunc) error {
	return addRuleFunc("iifname != \"lo\" meta l4proto tcp counter jump " + pf.chainName)
}

func (pf *packetFilter) AddRuleOut(addRuleFunc rule.AddFunc) error {
	return addRuleFunc("oifname != \"lo\" meta l4proto tcp counter jump " + pf.chainName)
}

type packetFilterFalse struct{}

// NewPacketFilterFalse returns a PacketFilter that does nothing.
func NewPacketFilterFalse() PacketFilter {
	return &packetFilterFalse{}
}

func (pf *packetFilterFalse) AddRuleIn(_ rule.AddFunc) error {
	return nil
}

func (pf *packetFilterFalse) AddRuleOut(_ rule.AddFunc) error {
	return nil
}
