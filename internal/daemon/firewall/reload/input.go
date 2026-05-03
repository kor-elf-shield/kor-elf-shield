package reload

import (
	"fmt"
	"net"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/rule"
	nftTable "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/table"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg"
)

func (r *reload) input(builder nft.BatchBuilder, packetfilter chain.PacketFilter) (nftTable.BlockList, error) {
	r.logger.Debug("Reloading input chain")

	batchInput, err := chain.NewBatchInput(
		builder,
		r.table.family,
		r.table.name,
		r.config.MetadataNaming.ChainInputName,
		r.config.Policy.DefaultAllowInput,
		r.config.Policy.InputPriority,
	)
	if err != nil {
		return nil, err
	}

	if err := r.reloadInputDnsNs(batchInput); err != nil {
		return nil, err
	}

	if err := batchInput.AddRule("iifname lo counter accept"); err != nil {
		return nil, err
	}

	beforeLocalInput, err := chain.NewBatchChain(builder, r.table.family, r.table.name, "before-local-input")
	if err != nil {
		return nil, err
	}
	if err := beforeLocalInput.AddRuleIn(batchInput.AddRule); err != nil {
		return nil, err
	}

	localInput, err := chain.NewBatchChain(builder, r.table.family, r.table.name, "local-input")
	if err != nil {
		return nil, err
	}
	if err := r.inputAddIPs(builder, batchInput, localInput); err != nil {
		return nil, err
	}

	afterLocalInput, err := chain.NewBatchChain(builder, r.table.family, r.table.name, "after-local-input")
	if err != nil {
		return nil, err
	}
	if err := afterLocalInput.AddRuleIn(batchInput.AddRule); err != nil {
		return nil, err
	}

	if err := packetfilter.AddRuleIn(batchInput.AddRule); err != nil {
		return nil, err
	}

	if err := r.inputICMP(batchInput); err != nil {
		return nil, err
	}

	if err := batchInput.AddRule("iifname != \"lo\" ct state related,established counter accept"); err != nil {
		return nil, err
	}

	if err := r.inputPorts(batchInput); err != nil {
		return nil, err
	}

	if r.config.Policy.DefaultAllowInput == false {
		drop := r.config.Policy.InputDrop.String()
		if err := batchInput.AddRule("iifname != \"lo\" " + drop); err != nil {
			return nil, err
		}
	}

	return r.blockList(builder, beforeLocalInput)
}

func (r *reload) reloadInputDnsNs(batchInput chain.Chain) error {
	if r.config.Options.DnsStrictNs {
		return nil
	}

	addresses, err := pkg.Resolv.Addresses()
	if err != nil {
		r.logger.Error(fmt.Sprintf("Failed to get nameservers: %s", err))
		return nil
	}
	for _, addr := range addresses {
		ip := net.ParseIP(addr)
		if ip == nil {
			r.logger.Error(fmt.Sprintf("Failed to parse nameserver address: %s", addr))
			continue
		}
		if ip.To4() != nil {
			if err := batchInput.AddRule("ip saddr " + addr + " iifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchInput.AddRule("ip saddr " + addr + " iifname != \"lo\" udp dport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchInput.AddRule("ip saddr " + addr + " iifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchInput.AddRule("ip saddr " + addr + " iifname != \"lo\" udp sport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		if ip.To16() != nil {
			if !r.config.IP6.Enable {
				r.logger.Warn(fmt.Sprintf("IPv6 is disabled, skipping nameserver address: %s", addr))
				continue
			}
			if err := batchInput.AddRule("ip6 saddr " + addr + " iifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchInput.AddRule("ip6 saddr " + addr + " iifname != \"lo\" udp dport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchInput.AddRule("ip6 saddr " + addr + " iifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchInput.AddRule("ip6 saddr " + addr + " iifname != \"lo\" udp sport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		r.logger.Error(fmt.Sprintf("Failed to parse nameserver address: %s", addr))
	}

	return nil
}

func (r *reload) inputAddIPs(builder nft.BatchBuilder, batchInput chain.Chain, localInput chain.Chain) error {
	if err := localInput.AddRuleIn(batchInput.AddRule); err != nil {
		return err
	}

	if err := r.inputPortKnocking(builder, localInput); err != nil {
		return err
	}

	for _, ipConfig := range r.config.IP4.InIPs {
		if err := rule.InputAddIP(localInput.AddRule, ipConfig, "ip"); err != nil {
			return err
		}
	}

	if !r.config.IP6.Enable {
		return nil
	}

	for _, ipConfig := range r.config.IP6.InIPs {
		if err := rule.InputAddIP(localInput.AddRule, ipConfig, "ip6"); err != nil {
			return err
		}
	}
	return nil
}

func (r *reload) inputPortKnocking(builder nft.BatchBuilder, localInput chain.Chain) error {
	if len(r.config.PortKnocking) == 0 {
		return nil
	}

	portKnocking, err := chain.NewBatchPortKnocking(builder, r.table.family, r.table.name, "port_knocking")
	if err != nil {
		return err
	}

	for _, portKnockingConfig := range r.config.PortKnocking {
		var knockName, prevKnockName string
		for index, knock := range portKnockingConfig.Knocks {
			prevKnockName = knockName
			knockName = fmt.Sprintf("knock_%s_%d", portKnockingConfig.Name, index)
			if index == 0 {
				if err := portKnocking.AddFirstStageRule(knockName, portKnockingConfig.IPVersion, knock.Port, knock.Timeout, knock.Action); err != nil {
					return err
				}
				continue
			}

			if err := portKnocking.AddNextStageRule(prevKnockName, knockName, portKnockingConfig.IPVersion, knock.Port, knock.Timeout, knock.Action); err != nil {
				return err
			}
		}

		expr := []string{
			portKnockingConfig.IPVersion.ToNft(), "saddr", "@" + knockName,
			portKnockingConfig.Port.ProtocolString(), "dport", portKnockingConfig.Port.NumberString(), "accept",
		}
		if err := localInput.AddRule(expr...); err != nil {
			return err
		}
	}

	if err := portKnocking.AddRuleIn(localInput.AddRule); err != nil {
		return err
	}

	return nil
}

func (r *reload) inputICMP(batchInput chain.Chain) error {
	drop := r.config.Policy.InputDrop.String()
	if r.config.IP4.IcmpIn == false {
		if err := batchInput.AddRule("iifname != \"lo\" ip protocol icmp icmp type echo-request counter " + drop); err != nil {
			return err
		}
		return r.inputICMPAfter(batchInput)
	}

	if r.config.IP4.IcmpInRate == "0" {
		return r.inputICMPAfter(batchInput)
	}

	if err := batchInput.AddRule("iifname != \"lo\" ip protocol icmp icmp type echo-request limit rate " + r.config.IP4.IcmpInRate + " counter accept"); err != nil {
		return err
	}
	if err := batchInput.AddRule("iifname != \"lo\" ip protocol icmp icmp type echo-request counter " + drop); err != nil {
		return err
	}

	return r.inputICMPAfter(batchInput)
}

func (r *reload) inputICMPAfter(batchInput chain.Chain) error {
	if r.config.IP4.IcmpTimestampDrop == true {
		drop := r.config.Policy.InputDrop.String()
		if err := batchInput.AddRule("iifname != \"lo\" ip protocol icmp icmp type timestamp-request " + drop); err != nil {
			return err
		}
	}

	if err := batchInput.AddRule("iifname != \"lo\" ip protocol icmp counter accept"); err != nil {
		return err
	}

	if r.config.IP6.Enable {
		if r.config.IP6.IcmpStrict {
			return r.inputICMP6Strict(batchInput)
		} else {
			if err := batchInput.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp counter accept"); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *reload) inputICMP6Strict(batchInput chain.Chain) error {
	if err := batchInput.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type destination-unreachable counter accept"); err != nil {
		return err
	}
	if err := batchInput.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type packet-too-big counter accept"); err != nil {
		return err
	}
	if err := batchInput.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type time-exceeded counter accept"); err != nil {
		return err
	}
	if err := batchInput.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type parameter-problem counter accept"); err != nil {
		return err
	}
	if err := batchInput.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type echo-request counter accept"); err != nil {
		return err
	}
	if err := batchInput.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type echo-reply counter accept"); err != nil {
		return err
	}
	if err := batchInput.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type nd-router-advert ip6 hoplimit 255 counter accept"); err != nil {
		return err
	}
	if err := batchInput.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type nd-neighbor-solicit ip6 hoplimit 255 counter accept"); err != nil {
		return err
	}
	if err := batchInput.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type nd-neighbor-advert ip6 hoplimit 255 counter accept"); err != nil {
		return err
	}
	if err := batchInput.AddRule("iifname != \"lo\" meta l4proto ipv6-icmp icmpv6 type nd-redirect ip6 hoplimit 255 counter accept"); err != nil {
		return err
	}
	return nil
}

func (r *reload) inputPorts(batchInput chain.Chain) error {
	for _, port := range r.config.InPorts {
		protocol := port.Port.ProtocolString()
		number := port.Port.NumberString()

		baseRule := "iifname != \"lo\" meta l4proto " + protocol + " ct state new " + protocol + " dport " + number

		if port.LimitRate != "" {
			rule := baseRule + " limit rate " + port.LimitRate + " counter " + port.Action.String()
			if err := batchInput.AddRule(rule); err != nil {
				return err
			}
			ruleDrop := baseRule + " counter " + r.config.Policy.InputDrop.String()
			if err := batchInput.AddRule(ruleDrop); err != nil {
				return err
			}
		} else {
			rule := baseRule + " counter " + port.Action.String()
			if err := batchInput.AddRule(rule); err != nil {
				return err
			}
		}
	}
	return nil
}
