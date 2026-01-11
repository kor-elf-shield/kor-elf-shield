package chain

import (
	"strings"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
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

	NewLocalForward() error
	LocalForward() LocalForward

	ClearRules() error

	NewNoneChain(chain string) (Chain, error)
	NewChain(chain string, baseChain nftChain.ChainOptions) (Chain, error)
}

type chains struct {
	input        Input
	output       Output
	forward      Forward
	packetFilter PacketFilter

	localInput   LocalInput
	localOutput  LocalOutput
	localForward LocalForward

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

func (c *chains) NewLocalForward() error {
	localForward, err := newLocalForward(c.nft, c.family, c.table)
	if err != nil {
		return err
	}
	c.localForward = localForward
	return nil
}

func (c *chains) LocalForward() LocalForward {
	return c.localForward
}

func (c *chains) ClearRules() error {
	return clearRules(c.nft, c.family, c.table)
}

func (c *chains) NewNoneChain(chainName string) (Chain, error) {
	return c.NewChain(chainName, nftChain.TypeNone)
}

func (c *chains) NewChain(chainName string, baseChain nftChain.ChainOptions) (Chain, error) {
	if err := c.nft.Chain().Add(c.family, c.table, chainName, baseChain); err != nil {
		return nil, err
	}

	return &chain{
		nft:    c.nft,
		family: c.family,
		table:  c.table,
		chain:  chainName,
	}, nil
}

func clearRules(nft nft.NFT, family nftFamily.Type, table string) error {
	if err := nft.Table().Delete(family, table); err != nil {
		if !strings.Contains(string(err.Error()), "delete table "+family.String()+" "+table) {
			return err
		}
	}

	return nil
}
