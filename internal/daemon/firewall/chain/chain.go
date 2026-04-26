package chain

import (
	"encoding/json"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type Chain interface {
	AddRule(expr ...string) error
	ListRules() ([]Rule, error)
	RemoveRuleByHandle(handle uint64) error
	Clear() error
}

type chain struct {
	nft    nft.NFT
	family family.Type
	table  string
	chain  string
}

type NftOutput struct {
	Nftables []NftElement `json:"nftables"`
}
type NftElement struct {
	Rule *Rule `json:"rule,omitempty"`
}

type Rule struct {
	Handle  uint64 `json:"handle"`
	Comment string `json:"comment"`
}

func (c *chain) AddRule(expr ...string) error {
	return c.nft.Rule().Add(c.family, c.table, c.chain, expr...)
}

func (c *chain) ListRules() ([]Rule, error) {
	args := []string{"-a", "-j", "list", "chain", c.family.String(), c.table, c.chain}
	jsonData, err := c.nft.Command().RunWithOutput(args...)
	if err != nil {
		return nil, err
	}
	var output NftOutput
	if err := json.Unmarshal([]byte(jsonData), &output); err != nil {
		return nil, err
	}

	var rules []Rule
	for _, el := range output.Nftables {
		if el.Rule != nil {
			rules = append(rules, *el.Rule)
		}
	}

	return rules, nil
}

func (c *chain) RemoveRuleByHandle(handle uint64) error {
	return c.nft.Rule().Delete(c.family, c.table, c.chain, handle)
}

func (c *chain) Clear() error {
	return c.nft.Chain().Clear(c.family, c.table, c.chain)
}
