package firewall

import (
	"fmt"
	"net"

	nftableCchain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	nftablesFamily "git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg"
)

func (f *firewall) reloadInput() error {
	family := nftablesFamily.INET
	tableName := f.config.MetadataNaming.TableName
	chainName := f.config.MetadataNaming.ChainInputName

	baseChain := nftableCchain.BaseChainOptions{
		Type:     nftableCchain.TypeFilter,
		Hook:     nftableCchain.HookInput,
		Priority: 0,
		Policy:   f.config.Policy.Input.ChainDefaultPolicy(),
		Device:   "",
	}
	if err := f.nft.Chain().Add(family, tableName, chainName, baseChain); err != nil {
		return err
	}

	if err := f.reloadInputDnsNs(); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "iifname lo counter accept"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "iifname != \"lo\" ct state related,established counter accept"); err != nil {
		return err
	}

	if f.config.Policy.Input == PolicyReject {
		if err := f.nft.Rule().Add(family, tableName, chainName, "reject"); err != nil {
			return err
		}
	}

	return nil
}

func (f *firewall) reloadInputDnsNs() error {
	if f.config.Options.DnsStrictNs {
		return nil
	}
	family := nftablesFamily.INET
	tableName := f.config.MetadataNaming.TableName
	chainName := f.config.MetadataNaming.ChainInputName

	addresses, err := pkg.Resolv.Addresses()
	if err != nil {
		f.logger.Error(fmt.Sprintf("Failed to get nameservers: %s", err))
		return nil
	}
	for _, addr := range addresses {
		ip := net.ParseIP(addr)
		if ip == nil {
			f.logger.Error(fmt.Sprintf("Failed to parse nameserver address: %s", addr))
			continue
		}
		if ip.To4() != nil {
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip saddr "+addr+" iifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip saddr "+addr+" iifname != \"lo\" udp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip saddr "+addr+" iifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip saddr "+addr+" iifname != \"lo\" udp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		if ip.To16() != nil {
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip6 saddr "+addr+" iifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip6 saddr "+addr+" iifname != \"lo\" udp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip6 saddr "+addr+" iifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip6 saddr "+addr+" iifname != \"lo\" udp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		f.logger.Error(fmt.Sprintf("Failed to parse nameserver address: %s", addr))
	}

	return nil
}
