package reload

import (
	"fmt"
	"net"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/rule"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg"
)

func (r *reload) output(builder nft.BatchBuilder, packetfilter chain.PacketFilter) error {
	r.logger.Debug("Reloading output chain")

	batchOutput, err := chain.NewBatchOutput(
		builder,
		r.table.family,
		r.table.name,
		r.config.MetadataNaming.ChainOutputName,
		r.config.Policy.DefaultAllowOutput,
		r.config.Policy.OutputPriority,
	)
	if err != nil {
		return err
	}

	if err := r.outputDnsNs(batchOutput); err != nil {
		return err
	}
	if err := r.outputDns(batchOutput); err != nil {
		return err
	}
	if err := batchOutput.AddRule("oifname lo counter accept"); err != nil {
		return err
	}

	localOutput, err := chain.NewBatchChain(builder, r.table.family, r.table.name, "local-output")
	if err != nil {
		return err
	}
	if err := r.outputAddIPs(batchOutput, localOutput); err != nil {
		return err
	}

	if err := packetfilter.AddRuleOut(batchOutput.AddRule); err != nil {
		return err
	}

	if err := r.outputICMP(batchOutput); err != nil {
		return err
	}

	if err := batchOutput.AddRule("oifname != \"lo\" ct state related,established counter accept"); err != nil {
		return err
	}

	if err := r.outputPorts(batchOutput); err != nil {
		return err
	}

	if r.config.Policy.DefaultAllowOutput == false {
		drop := r.config.Policy.OutputDrop.String()
		if err := batchOutput.AddRule("oifname != \"lo\" " + drop); err != nil {
			return err
		}
	}

	return nil
}

func (r *reload) outputDnsNs(batchOutput chain.Chain) error {
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
			if err := batchOutput.AddRule("ip daddr " + addr + " oifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchOutput.AddRule("ip daddr " + addr + " oifname != \"lo\" udp dport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchOutput.AddRule("ip daddr " + addr + " oifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchOutput.AddRule("ip daddr " + addr + " oifname != \"lo\" udp sport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		if ip.To16() != nil {
			if err := batchOutput.AddRule("ip6 daddr " + addr + " oifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchOutput.AddRule("ip6 daddr " + addr + " oifname != \"lo\" udp dport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchOutput.AddRule("ip6 daddr " + addr + " oifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := batchOutput.AddRule("ip6 daddr " + addr + " oifname != \"lo\" udp sport 53 counter accept"); err != nil {
				r.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		r.logger.Error(fmt.Sprintf("Failed to parse nameserver address: %s", addr))
	}

	return nil
}

func (r *reload) outputDns(batchOutput chain.Chain) error {
	if r.config.Options.DnsStrict {
		return nil
	}

	if err := batchOutput.AddRule("oifname != \"lo\" tcp dport 53 counter accept"); err != nil {
		return err
	}
	if err := batchOutput.AddRule("oifname != \"lo\" udp dport 53 counter accept"); err != nil {
		return err
	}
	if err := batchOutput.AddRule("oifname != \"lo\" tcp sport 53 counter accept"); err != nil {
		return err
	}
	if err := batchOutput.AddRule("oifname != \"lo\" udp sport 53 counter accept"); err != nil {
		return err
	}

	return nil
}

func (r *reload) outputAddIPs(batchOutput chain.Chain, localOutput chain.Chain) error {
	if err := localOutput.AddRuleOut(batchOutput.AddRule); err != nil {
		return err
	}

	for _, ipConfig := range r.config.IP4.OutIPs {
		if err := rule.OutputAddIP(localOutput.AddRule, ipConfig, "ip"); err != nil {
			return err
		}
	}

	if !r.config.IP6.Enable {
		return nil
	}

	for _, ipConfig := range r.config.IP6.OutIPs {
		if err := rule.OutputAddIP(localOutput.AddRule, ipConfig, "ip6"); err != nil {
			return err
		}
	}

	return nil
}

func (r *reload) outputICMP(batchOutput chain.Chain) error {
	drop := r.config.Policy.OutputDrop.String()
	if r.config.IP4.IcmpOut == false {
		if err := batchOutput.AddRule("oifname != \"lo\" ip protocol icmp icmp type echo-request counter " + drop); err != nil {
			return err
		}
		return r.outputICMPAfter(batchOutput)
	}

	if r.config.IP4.IcmpOutRate == "0" {
		return r.outputICMPAfter(batchOutput)
	}

	if err := batchOutput.AddRule("oifname != \"lo\" ip protocol icmp icmp type echo-request limit rate " + r.config.IP4.IcmpInRate + " counter accept"); err != nil {
		return err
	}
	if err := batchOutput.AddRule("oifname != \"lo\" ip protocol icmp icmp type echo-request counter " + drop); err != nil {
		return err
	}

	return r.outputICMPAfter(batchOutput)
}

func (r *reload) outputICMPAfter(batchOutput chain.Chain) error {
	if r.config.IP4.IcmpTimestampDrop == true {
		drop := r.config.Policy.OutputDrop.String()
		if err := batchOutput.AddRule("oifname != \"lo\" ip protocol icmp icmp type timestamp-reply " + drop); err != nil {
			return err
		}
	}

	if err := batchOutput.AddRule("oifname != \"lo\" ip protocol icmp counter accept"); err != nil {
		return err
	}
	return nil
}

func (r *reload) outputPorts(batchOutput chain.Chain) error {
	for _, port := range r.config.OutPorts {
		protocol := port.Port.ProtocolString()
		number := port.Port.NumberString()
		baseRule := "oifname != \"lo\" meta l4proto " + protocol + " ct state new " + protocol + " dport " + number

		if port.LimitRate != "" {
			addRule := baseRule + " limit rate " + port.LimitRate + " counter " + port.Action.String()
			if err := batchOutput.AddRule(addRule); err != nil {
				return err
			}
			addRuleDrop := baseRule + " counter " + r.config.Policy.InputDrop.String()
			if err := batchOutput.AddRule(addRuleDrop); err != nil {
				return err
			}
		} else {
			addRule := baseRule + " counter " + port.Action.String()
			if err := batchOutput.AddRule(addRule); err != nil {
				return err
			}
		}
	}
	return nil
}
