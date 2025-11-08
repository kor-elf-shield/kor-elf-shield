package chain

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	nftFamily "git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type PacketFilter interface {
	AddRuleIn(AddRuleFunc func(expr ...string) error) error
	AddRuleOut(AddRuleFunc func(expr ...string) error) error
}

type packetFilter struct {
	enable      bool
	invalidName string
}

// newPacketFilter Drop out of order packets and packets in an INVALID state in nftables connection tracking.
func newPacketFilter(nft nft.NFT, family nftFamily.Type, table string, enable bool) (PacketFilter, error) {
	chainInvalidName := "INVALID"
	if !enable {
		return &packetFilter{
			enable:      enable,
			invalidName: chainInvalidName,
		}, nil
	}

	chainName := "INVDROP"

	if err := nft.Chain().Add(family, table, chainName, nftChain.TypeNone); err != nil {
		return nil, err
	}
	if err := nft.Rule().Add(family, table, chainName, "counter drop"); err != nil {
		return nil, err
	}

	if err := nft.Chain().Add(family, table, chainInvalidName, nftChain.TypeNone); err != nil {
		return nil, err
	}

	if err := nft.Rule().Add(family, table, chainInvalidName, "ct state invalid counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := nft.Rule().Add(family, table, chainInvalidName, "tcp flags ! fin,syn,rst,psh,ack,urg counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := nft.Rule().Add(family, table, chainInvalidName, "tcp flags & (fin | syn | rst | psh | ack | urg) == fin | syn | rst | psh | ack | urg counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := nft.Rule().Add(family, table, chainInvalidName, "tcp flags & (fin | syn) == fin | syn counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := nft.Rule().Add(family, table, chainInvalidName, "tcp flags & (syn | rst) == syn | rst counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := nft.Rule().Add(family, table, chainInvalidName, "tcp flags & (fin | rst) == fin | rst counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := nft.Rule().Add(family, table, chainInvalidName, "tcp flags & (fin | ack) == fin counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := nft.Rule().Add(family, table, chainInvalidName, "tcp flags & (psh | ack) == psh counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := nft.Rule().Add(family, table, chainInvalidName, "tcp flags & (ack | urg) == urg counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := nft.Rule().Add(family, table, chainInvalidName, "tcp flags & (fin | syn | rst | ack) != syn ct state new counter jump INVDROP"); err != nil {
		return nil, err
	}

	return &packetFilter{
		enable:      enable,
		invalidName: chainInvalidName,
	}, nil
}

func (f *packetFilter) AddRuleIn(AddRuleFunc func(expr ...string) error) error {
	if !f.enable {
		return nil
	}
	return AddRuleFunc("iifname != \"lo\" meta l4proto tcp counter jump " + f.invalidName)
}

func (f *packetFilter) AddRuleOut(AddRuleFunc func(expr ...string) error) error {
	if !f.enable {
		return nil
	}
	return AddRuleFunc("oifname != \"lo\" meta l4proto tcp counter jump " + f.invalidName)
}
