package block

import (
	"fmt"
	"net"
	"strings"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type ListIP interface {
	// AddIP Add an IP address to the list.
	AddIP(addr net.IP, banSeconds uint32) error

	// DeleteIP Delete an IP address from the list.
	DeleteIP(addr net.IP) error

	// AddRuleToChain Add a rule to the parent chain.
	AddRuleToChain(chainAddRuleFunc func(expr ...string) error, action string) error
}

type listIP struct {
	listIPv4 List
	listIPv6 List
}

func NewListIP(nft nft.NFT, family family.Type, table string, name string) (ListIP, error) {
	params := "type ipv4_addr; flags interval, timeout;"
	listName := name + "_ip4"
	listIPv4, err := newList(nft, family, table, listName, params)
	if err != nil {
		return nil, err
	}

	params = "type ipv6_addr; flags interval, timeout;"
	listName = name + "_ip6"
	listIPv6, err := newList(nft, family, table, listName, params)
	if err != nil {
		return nil, err
	}

	return &listIP{
		listIPv4: listIPv4,
		listIPv6: listIPv6,
	}, nil
}

func (l *listIP) AddIP(addr net.IP, banSeconds uint32) error {
	el := []string{addr.String()}
	if banSeconds > 0 {
		el = append(el, "timeout", fmt.Sprintf("%ds", banSeconds))
	}

	element := strings.Join(el, " ")

	if addr.To4() != nil {
		return l.listIPv4.AddElement(element)
	}

	return l.listIPv6.AddElement(element)
}

func (l *listIP) DeleteIP(addr net.IP) error {
	if addr == nil {
		return fmt.Errorf("IP address cannot be nil")
	}
	if addr.To4() != nil {
		return l.listIPv4.DeleteElement(addr.String())
	}

	return l.listIPv6.DeleteElement(addr.String())
}

func (l *listIP) AddRuleToChain(chainAddRuleFunc func(expr ...string) error, action string) error {
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
