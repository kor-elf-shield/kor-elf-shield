package chain

import (
	"encoding/json"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/data"
)

type Docker interface {
	Name() string
	AddRule(builder nft.BatchBuilder, expr ...string) error
	JumpTo(builder nft.BatchBuilder, chain Docker, rule string, comment string) error
	ListRules() ([]data.Rule, error)
	RemoveRuleByHandle(builder nft.BatchBuilder, handle uint64) error
	Clear(builder nft.BatchBuilder) error
}

type docker struct {
	nft    nft.NFT
	family family.Type
	table  string
	chain  string
}

func NewDocker(nft nft.NFT, family family.Type, table, chain string) Docker {
	return &docker{
		nft:    nft,
		family: family,
		table:  table,
		chain:  chain,
	}
}

func (d *docker) Name() string {
	return d.chain
}

func (d *docker) AddRule(builder nft.BatchBuilder, expr ...string) error {
	return builder.Rule().Add(d.family, d.table, d.chain, expr...)
}

func (d *docker) JumpTo(builder nft.BatchBuilder, chain Docker, rule string, comment string) error {
	args := []string{rule, "jump", d.chain, comment}
	return chain.AddRule(builder, args...)
}

func (d *docker) ListRules() ([]data.Rule, error) {
	args := []string{"-a", "-j", "list", "chain", d.family.String(), d.table, d.chain}
	jsonData, err := d.nft.Command().RunWithOutput(args...)
	if err != nil {
		return nil, err
	}
	var output data.NftOutput
	if err := json.Unmarshal([]byte(jsonData), &output); err != nil {
		return nil, err
	}

	var rules []data.Rule
	for _, el := range output.Nftables {
		if el.Rule != nil {
			rules = append(rules, *el.Rule)
		}
	}

	return rules, nil
}

func (d *docker) RemoveRuleByHandle(builder nft.BatchBuilder, handle uint64) error {
	return builder.Rule().Delete(d.family, d.table, d.chain, handle)
}

func (d *docker) Clear(builder nft.BatchBuilder) error {
	return builder.Chain().Clear(d.family, d.table, d.chain)
}
