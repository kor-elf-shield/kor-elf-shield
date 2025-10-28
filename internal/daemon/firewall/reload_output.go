package firewall

import (
	"fmt"
	"net"

	nftableCchain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	nftablesFamily "git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg"
)

func (f *firewall) reloadOutput() error {
	family := nftablesFamily.INET
	tableName := f.config.MetadataNaming.TableName
	chainName := f.config.MetadataNaming.ChainOutputName

	baseChain := nftableCchain.BaseChainOptions{
		Type:     nftableCchain.TypeFilter,
		Hook:     nftableCchain.HookOutput,
		Priority: 0,
		Policy:   f.config.Policy.Output.ChainDefaultPolicy(),
		Device:   "",
	}
	if err := f.nft.Chain().Add(family, tableName, chainName, baseChain); err != nil {
		return err
	}

	if err := f.reloadOutputDnsNs(); err != nil {
		return err
	}
	if err := f.reloadOutputDns(); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "oifname lo counter accept"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "oifname != \"lo\" meta l4proto tcp counter jump INVALID"); err != nil {
		return err
	}

	if err := f.reloadOutputICMP(); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "oifname != \"lo\" ct state related,established counter accept"); err != nil {
		return err
	}

	if f.config.Policy.Output == PolicyReject {
		if err := f.nft.Rule().Add(family, tableName, chainName, f.config.Policy.Output.String()); err != nil {
			return err
		}
	}

	return nil
}

func (f *firewall) reloadOutputDns() error {
	if f.config.Options.DnsStrict {
		return nil
	}
	family := nftablesFamily.INET
	tableName := f.config.MetadataNaming.TableName
	chainName := f.config.MetadataNaming.ChainOutputName

	if err := f.nft.Rule().Add(family, tableName, chainName, "oifname != \"lo\" tcp dport 53 counter accept"); err != nil {
		return err
	}
	if err := f.nft.Rule().Add(family, tableName, chainName, "oifname != \"lo\" udp dport 53 counter accept"); err != nil {
		return err
	}
	if err := f.nft.Rule().Add(family, tableName, chainName, "oifname != \"lo\" tcp sport 53 counter accept"); err != nil {
		return err
	}
	if err := f.nft.Rule().Add(family, tableName, chainName, "oifname != \"lo\" udp sport 53 counter accept"); err != nil {
		return err
	}

	return nil
}

func (f *firewall) reloadOutputDnsNs() error {
	if f.config.Options.DnsStrictNs {
		return nil
	}
	family := nftablesFamily.INET
	tableName := f.config.MetadataNaming.TableName
	chainName := f.config.MetadataNaming.ChainOutputName

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
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip daddr "+addr+" oifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip daddr "+addr+" oifname != \"lo\" udp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip daddr "+addr+" oifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip daddr "+addr+" oifname != \"lo\" udp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		if ip.To16() != nil {
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip6 daddr "+addr+" oifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip6 daddr "+addr+" oifname != \"lo\" udp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip6 daddr "+addr+" oifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := f.nft.Rule().Add(family, tableName, chainName, "ip6 daddr "+addr+" oifname != \"lo\" udp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		f.logger.Error(fmt.Sprintf("Failed to parse nameserver address: %s", addr))
	}

	return nil
}

func (f *firewall) reloadOutputICMP() error {
	family := nftablesFamily.INET
	tableName := f.config.MetadataNaming.TableName
	chainName := f.config.MetadataNaming.ChainOutputName
	if f.config.IP4.IcmpOut == false {
		if err := f.nft.Rule().Add(family, tableName, chainName, "oifname != \"lo\" ip protocol icmp icmp type echo-request counter drop"); err != nil {
			return err
		}
		return f.reloadOutputICMPAfter()
	}

	if f.config.IP4.IcmpOutRate == "0" {
		return f.reloadOutputICMPAfter()
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "oifname != \"lo\" ip protocol icmp icmp type echo-request limit rate "+f.config.IP4.IcmpInRate+" counter accept"); err != nil {
		return err
	}
	if err := f.nft.Rule().Add(family, tableName, chainName, "oifname != \"lo\" ip protocol icmp icmp type echo-request counter drop"); err != nil {
		return err
	}

	return f.reloadOutputICMPAfter()
}
func (f *firewall) reloadOutputICMPAfter() error {
	family := nftablesFamily.INET
	tableName := f.config.MetadataNaming.TableName
	chainName := f.config.MetadataNaming.ChainOutputName

	if f.config.IP4.IcmpTimestampDrop == true {
		if err := f.nft.Rule().Add(family, tableName, chainName, "oifname != \"lo\" ip protocol icmp icmp type timestamp-request drop"); err != nil {
			return err
		}
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "oifname != \"lo\" ip protocol icmp counter accept"); err != nil {
		return err
	}
	return nil
}
