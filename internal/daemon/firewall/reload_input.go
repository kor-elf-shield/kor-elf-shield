package firewall

import (
	"fmt"
	"net"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg"
)

func (f *firewall) reloadInput() error {
	f.logger.Debug("Reloading input chain")
	err := f.chains.NewInput(f.config.MetadataNaming.ChainInputName, f.config.Policy.DefaultAllowInput)
	if err != nil {
		return err
	}
	chain := f.chains.Input()

	if err := f.reloadInputDnsNs(); err != nil {
		return err
	}

	if err := chain.AddRule("iifname lo counter accept"); err != nil {
		return err
	}

	if err := f.chains.PacketFilter().AddRuleIn(chain.AddRule); err != nil {
		return err
	}

	if err := f.reloadInputICMP(); err != nil {
		return err
	}

	if err := chain.AddRule("iifname != \"lo\" ct state related,established counter accept"); err != nil {
		return err
	}

	if f.config.Policy.DefaultAllowInput == false {
		drop := f.config.Policy.InputDrop.String()
		if err := chain.AddRule("iifname != \"lo\" " + drop); err != nil {
			return err
		}
	}

	return nil
}

func (f *firewall) reloadInputDnsNs() error {
	if f.config.Options.DnsStrictNs {
		return nil
	}
	chain := f.chains.Input()

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
			if err := chain.AddRule("ip saddr " + addr + " iifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip saddr " + addr + " iifname != \"lo\" udp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip saddr " + addr + " iifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip saddr " + addr + " iifname != \"lo\" udp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		if ip.To16() != nil {
			if !f.config.IP6.Enable {
				f.logger.Warn(fmt.Sprintf("IPv6 is disabled, skipping nameserver address: %s", addr))
				continue
			}
			if err := chain.AddRule("ip6 saddr " + addr + " iifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip6 saddr " + addr + " iifname != \"lo\" udp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip6 saddr " + addr + " iifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip6 saddr " + addr + " iifname != \"lo\" udp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		f.logger.Error(fmt.Sprintf("Failed to parse nameserver address: %s", addr))
	}

	return nil
}

func (f *firewall) reloadInputICMP() error {
	chain := f.chains.Input()
	drop := f.config.Policy.InputDrop.String()
	if f.config.IP4.IcmpIn == false {
		if err := chain.AddRule("iifname != \"lo\" ip protocol icmp icmp type echo-request counter " + drop); err != nil {
			return err
		}
		return f.reloadInputICMPAfter()
	}

	if f.config.IP4.IcmpInRate == "0" {
		return f.reloadInputICMPAfter()
	}

	if err := chain.AddRule("iifname != \"lo\" ip protocol icmp icmp type echo-request limit rate " + f.config.IP4.IcmpInRate + " counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("iifname != \"lo\" ip protocol icmp icmp type echo-request counter " + drop); err != nil {
		return err
	}

	return f.reloadInputICMPAfter()
}
func (f *firewall) reloadInputICMPAfter() error {
	chain := f.chains.Input()

	if f.config.IP4.IcmpTimestampDrop == true {
		drop := f.config.Policy.InputDrop.String()
		if err := chain.AddRule("iifname != \"lo\" ip protocol icmp icmp type timestamp-request " + drop); err != nil {
			return err
		}
	}

	if err := chain.AddRule("iifname != \"lo\" ip protocol icmp counter accept"); err != nil {
		return err
	}

	if f.config.IP6.Enable {
		if f.config.IP6.IcmpStrict {
			return f.reloadInputICMP6Strict()
		} else {
			if err := chain.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp counter accept"); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *firewall) reloadInputICMP6Strict() error {
	chain := f.chains.Input()
	if err := chain.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type destination-unreachable counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type packet-too-big counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type time-exceeded counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type parameter-problem counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type echo-request counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type echo-reply counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type nd-router-advert ip6 hoplimit 255 counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type nd-neighbor-solicit ip6 hoplimit 255 counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type nd-neighbor-advert ip6 hoplimit 255 counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type nd-redirect ip6 hoplimit 255 counter accept"); err != nil {
		return err
	}
	return nil
}
