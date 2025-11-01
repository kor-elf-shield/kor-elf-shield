package chain

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	nftFamily "git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type Chains interface {
	NewPacketFilter(enable bool) error
	PacketFilter() PacketFilter

	NewInput(chain string, defaultAllow bool) error
	Input() Input

	NewOutput(chain string, defaultAllow bool) error
	Output() Output

	NewForward(chain string, defaultAllow bool) error
	Forward() Forward

	NewLocalInput() error
	LocalInput() LocalInput

	NewLocalOutput() error
	LocalOutput() LocalOutput
}

type chains struct {
	input        Input
	output       Output
	forward      Forward
	packetFilter PacketFilter

	localInput  LocalInput
	localOutput LocalOutput

	family nftFamily.Type
	table  string
	nft    nft.NFT
}

func NewChains(nft nft.NFT, table string) (Chains, error) {
	if err := nft.Clear(); err != nil {
		return nil, err
	}

	family := nftFamily.INET
	if err := nft.Table().Add(family, table); err != nil {
		return nil, err
	}

	return &chains{
		nft:    nft,
		table:  table,
		family: family,
	}, nil
}

func (c *chains) NewPacketFilter(enable bool) error {
	filter, err := newPacketFilter(c.nft, c.family, c.table, enable)
	if err != nil {
		return err
	}
	c.packetFilter = filter

	return nil
}

func (c *chains) PacketFilter() PacketFilter {
	return c.packetFilter
}

func (c *chains) NewInput(chain string, defaultAllow bool) error {
	input, err := newInput(c.nft, c.family, c.table, chain, defaultAllow)
	if err != nil {
		return err
	}
	c.input = input

	return nil
}

func (c *chains) Input() Input {
	return c.input
}

func (c *chains) NewOutput(chain string, defaultAllow bool) error {
	output, err := newOutput(c.nft, c.family, c.table, chain, defaultAllow)
	if err != nil {
		return err
	}
	c.output = output

	return nil
}

func (c *chains) Output() Output {
	return c.output
}

func (c *chains) NewForward(chain string, defaultAllow bool) error {
	forward, err := newForward(c.nft, c.family, c.table, chain, defaultAllow)
	if err != nil {
		return err
	}
	c.forward = forward

	return nil
}

func (c *chains) Forward() Forward {
	return c.forward
}

func (c *chains) NewLocalInput() error {
	localInput, err := newLocalInput(c.nft, c.family, c.table)
	if err != nil {
		return err
	}
	c.localInput = localInput
	return nil
}

func (c *chains) LocalInput() LocalInput {
	return c.localInput
}

func (c *chains) NewLocalOutput() error {
	localOutput, err := newLocalOutput(c.nft, c.family, c.table)
	if err != nil {
		return err
	}
	c.localOutput = localOutput
	return nil
}

func (c *chains) LocalOutput() LocalOutput {
	return c.localOutput
}
