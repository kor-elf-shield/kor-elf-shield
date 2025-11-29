package chain

import (
	"strings"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	nftFamily "git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type Chains interface {
	NewPacketFilter(enable bool) error
	PacketFilter() PacketFilter

	NewInput(chain string, defaultAllow bool, priority int) error
	Input() Input

	NewOutput(chain string, defaultAllow bool, priority int) error
	Output() Output

	NewForward(chain string, defaultAllow bool, priority int) error
	Forward() Forward

	NewLocalInput() error
	LocalInput() LocalInput

	NewLocalOutput() error
	LocalOutput() LocalOutput

	ClearRules() error
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
	family := nftFamily.INET

	if err := clearRules(nft, family, table); err != nil {
		return nil, err
	}

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

func (c *chains) NewInput(chain string, defaultAllow bool, priority int) error {
	input, err := newInput(c.nft, c.family, c.table, chain, defaultAllow, priority)
	if err != nil {
		return err
	}
	c.input = input

	return nil
}

func (c *chains) Input() Input {
	return c.input
}

func (c *chains) NewOutput(chain string, defaultAllow bool, priority int) error {
	output, err := newOutput(c.nft, c.family, c.table, chain, defaultAllow, priority)
	if err != nil {
		return err
	}
	c.output = output

	return nil
}

func (c *chains) Output() Output {
	return c.output
}

func (c *chains) NewForward(chain string, defaultAllow bool, priority int) error {
	forward, err := newForward(c.nft, c.family, c.table, chain, defaultAllow, priority)
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

func (c *chains) ClearRules() error {
	return clearRules(c.nft, c.family, c.table)
}

func clearRules(nft nft.NFT, family nftFamily.Type, table string) error {
	if err := nft.Table().Delete(family, table); err != nil {
		if !strings.Contains(string(err.Error()), "delete table "+family.String()+" "+table) {
			return err
		}
	}

	return nil
}
