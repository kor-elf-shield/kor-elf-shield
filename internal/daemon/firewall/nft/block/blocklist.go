package block

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	nftFirewall "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/rule"
)

type Blocklist interface {
	// ReplaceElementsIPv4 Replacing IP addresses.
	ReplaceElementsIPv4(ips []string) error

	// ReplaceElementsIPv6 Replacing IP addresses.
	ReplaceElementsIPv6(ips []string) error

	// AddRuleToChain Add a rule to the parent chain.
	AddRuleToChain(chainAddRuleFunc rule.AddFunc, action string) error
}

type blocklist struct {
	listIPv4 List
	listIPv6 List
}

func NewBlocklist(nft nftFirewall.NFT, builder nft.BatchBuilder, family family.Type, table string, name string) (Blocklist, error) {
	params := "type ipv4_addr; flags interval; auto-merge;"
	listName := name + "_ip4"
	listIPv4, err := newList(nft, builder, family, table, listName, params)
	if err != nil {
		return nil, err
	}

	params = "type ipv6_addr; flags interval; auto-merge;"
	listName = name + "_ip6"
	listIPv6, err := newList(nft, builder, family, table, listName, params)
	if err != nil {
		return nil, err
	}

	return &blocklist{
		listIPv4: listIPv4,
		listIPv6: listIPv6,
	}, nil
}

func (l *blocklist) ReplaceElementsIPv4(ips []string) error {
	return l.listIPv4.ReplaceElements(ips)
}

func (l *blocklist) ReplaceElementsIPv6(ips []string) error {
	return l.listIPv6.ReplaceElements(ips)
}

func (l *blocklist) AddRuleToChain(chainAddRuleFunc rule.AddFunc, action string) error {
	rule := "ip saddr @" + l.listIPv4.Name() + " " + action
	if err := chainAddRuleFunc(rule); err != nil {
		return err
	}

	rule = "ip6 saddr @" + l.listIPv6.Name() + " " + action
	if err := chainAddRuleFunc(rule); err != nil {
		return err
	}

	return nil
}
