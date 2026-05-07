package reload

import (
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
)

func (r *reload) packetFilter(builder nft.BatchBuilder) (chain.PacketFilter, error) {
	if !r.config.Options.PacketFilter {
		return chain.NewPacketFilterFalse(), nil
	}

	chainInvalidName := "INVALID"

	chainName := "INVDROP"

	if err := r.addChain(builder, chainName, nftChain.TypeNone); err != nil {
		return nil, err
	}
	if err := r.addRule(builder, chainName, "counter drop"); err != nil {
		return nil, err
	}

	if err := r.addChain(builder, chainInvalidName, nftChain.TypeNone); err != nil {
		return nil, err
	}

	if err := r.addRule(builder, chainInvalidName, "ct state invalid counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := r.addRule(builder, chainInvalidName, "tcp flags ! fin,syn,rst,psh,ack,urg counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := r.addRule(builder, chainInvalidName, "tcp flags & (fin | syn | rst | psh | ack | urg) == fin | syn | rst | psh | ack | urg counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := r.addRule(builder, chainInvalidName, "tcp flags & (fin | syn) == fin | syn counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := r.addRule(builder, chainInvalidName, "tcp flags & (syn | rst) == syn | rst counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := r.addRule(builder, chainInvalidName, "tcp flags & (fin | rst) == fin | rst counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := r.addRule(builder, chainInvalidName, "tcp flags & (fin | ack) == fin counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := r.addRule(builder, chainInvalidName, "tcp flags & (psh | ack) == psh counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := r.addRule(builder, chainInvalidName, "tcp flags & (ack | urg) == urg counter jump INVDROP"); err != nil {
		return nil, err
	}

	if err := r.addRule(builder, chainInvalidName, "tcp flags & (fin | syn | rst | ack) != syn ct state new counter jump INVDROP"); err != nil {
		return nil, err
	}

	return chain.NewPacketFilter(chainInvalidName), nil
}
