package block

import (
	"fmt"
	"net"
	"strings"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	nftFirewall "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/rule"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
)

type ListIPWithPort interface {
	// AddIP Add an IP address to the list.
	AddIP(addr net.IP, ports []types.L4Port, banSeconds uint32) error

	// AddBatchIP Add an IP address to the list.
	AddBatchIP(builder nft.BatchBuilder, addr net.IP, ports []types.L4Port, banSeconds uint32) error

	// DeleteIP Delete an IP address from the list.
	DeleteIP(addr net.IP, port types.L4Port) error

	// AddRuleToChain Add a rule to the parent chain.
	AddRuleToChain(chainAddRuleFunc rule.AddFunc, action string) error
}

type listIPWithPort struct {
	listIPv4 List
	listIPv6 List
}

func NewListIPWithPort(nft nftFirewall.NFT, builder nft.BatchBuilder, family family.Type, table string, name string) (ListIPWithPort, error) {
	params := "type ipv4_addr . inet_proto . inet_service; flags interval, timeout;"
	listName := name + "_ip4"
	listIPv4, err := newList(nft, builder, family, table, listName, params)
	if err != nil {
		return nil, err
	}

	params = "type ipv6_addr . inet_proto . inet_service; flags interval, timeout;"
	listName = name + "_ip6"
	listIPv6, err := newList(nft, builder, family, table, listName, params)
	if err != nil {
		return nil, err
	}

	return &listIPWithPort{
		listIPv4: listIPv4,
		listIPv6: listIPv6,
	}, nil
}

func (l *listIPWithPort) AddIP(addr net.IP, ports []types.L4Port, banSeconds uint32) error {
	if len(ports) == 0 {
		return fmt.Errorf("ports is empty")
	}

	var elements []string
	for _, port := range ports {
		el := []string{fmt.Sprintf("%s . %s . %d", addr.String(), port.ProtocolString(), port.Number())}
		if banSeconds > 0 {
			el = append(el, "timeout", fmt.Sprintf("%ds", banSeconds))
		}

		elements = append(elements, strings.Join(el, " "))
	}

	element := strings.Join(elements, ",")
	if addr.To4() != nil {
		return l.listIPv4.AddElement(element)
	}

	return l.listIPv6.AddElement(element)
}

func (l *listIPWithPort) AddBatchIP(builder nft.BatchBuilder, addr net.IP, ports []types.L4Port, banSeconds uint32) error {
	if len(ports) == 0 {
		return fmt.Errorf("ports is empty")
	}

	var elements []string
	for _, port := range ports {
		el := []string{fmt.Sprintf("%s . %s . %d", addr.String(), port.ProtocolString(), port.Number())}
		if banSeconds > 0 {
			el = append(el, "timeout", fmt.Sprintf("%ds", banSeconds))
		}

		elements = append(elements, strings.Join(el, " "))
	}

	element := strings.Join(elements, ",")
	if addr.To4() != nil {
		return l.listIPv4.AddBatchElement(builder, element)
	}

	return l.listIPv6.AddBatchElement(builder, element)
}

func (l *listIPWithPort) DeleteIP(addr net.IP, port types.L4Port) error {
	if addr == nil {
		return fmt.Errorf("IP address cannot be nil")
	}
	if port.ToString() == "" {
		return fmt.Errorf("port cannot be empty")
	}

	element := fmt.Sprintf("%s . %s . %d", addr.String(), port.ProtocolString(), port.Number())

	if addr.To4() != nil {
		return l.listIPv4.DeleteElement(element)
	}

	return l.listIPv6.DeleteElement(element)
}

func (l *listIPWithPort) AddRuleToChain(chainAddRuleFunc rule.AddFunc, action string) error {
	rule := "ip saddr . meta l4proto . th dport @" + l.listIPv4.Name() + " " + action
	if err := chainAddRuleFunc(rule); err != nil {
		return err
	}

	rule = "ip6 saddr . meta l4proto . th dport @" + l.listIPv6.Name() + " " + action
	if err := chainAddRuleFunc(rule); err != nil {
		return err
	}

	return nil
}
