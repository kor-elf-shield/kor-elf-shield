package firewall

import (
	"fmt"
	"net"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg"
)

func (f *firewall) reloadOutput() error {
	f.logger.Debug("Reloading output chain")
	err := f.chains.NewOutput(f.config.MetadataNaming.ChainOutputName, f.config.Policy.DefaultAllowOutput, f.config.Policy.OutputPriority)
	if err != nil {
		return err
	}
	chain := f.chains.Output()

	if err := f.reloadOutputDnsNs(); err != nil {
		return err
	}
	if err := f.reloadOutputDns(); err != nil {
		return err
	}
	if err := chain.AddRule("oifname lo counter accept"); err != nil {
		return err
	}

	if err := f.reloadOutputAddIPs(); err != nil {
		return err
	}

	if err := f.chains.PacketFilter().AddRuleOut(chain.AddRule); err != nil {
		return err
	}

	if err := f.reloadOutputICMP(); err != nil {
		return err
	}

	if err := chain.AddRule("oifname != \"lo\" ct state related,established counter accept"); err != nil {
		return err
	}

	if err := f.reloadOutputPorts(); err != nil {
		return err
	}

	if f.config.Policy.DefaultAllowOutput == false {
		drop := f.config.Policy.OutputDrop.String()
		if err := chain.AddRule("oifname != \"lo\" " + drop); err != nil {
			return err
		}
	}

	return nil
}

func (f *firewall) reloadOutputDns() error {
	if f.config.Options.DnsStrict {
		return nil
	}
	chain := f.chains.Output()

	if err := chain.AddRule("oifname != \"lo\" tcp dport 53 counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("oifname != \"lo\" udp dport 53 counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("oifname != \"lo\" tcp sport 53 counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("oifname != \"lo\" udp sport 53 counter accept"); err != nil {
		return err
	}

	return nil
}

func (f *firewall) reloadOutputDnsNs() error {
	if f.config.Options.DnsStrictNs {
		return nil
	}
	chain := f.chains.Output()

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
			if err := chain.AddRule("ip daddr " + addr + " oifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip daddr " + addr + " oifname != \"lo\" udp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip daddr " + addr + " oifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip daddr " + addr + " oifname != \"lo\" udp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		if ip.To16() != nil {
			if err := chain.AddRule("ip6 daddr " + addr + " oifname != \"lo\" tcp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip6 daddr " + addr + " oifname != \"lo\" udp dport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip6 daddr " + addr + " oifname != \"lo\" tcp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			if err := chain.AddRule("ip6 daddr " + addr + " oifname != \"lo\" udp sport 53 counter accept"); err != nil {
				f.logger.Error(fmt.Sprintf("Failed to add rule: %s", err))
			}
			continue
		}

		f.logger.Error(fmt.Sprintf("Failed to parse nameserver address: %s", addr))
	}

	return nil
}

func (f *firewall) reloadOutputICMP() error {
	chain := f.chains.Output()
	drop := f.config.Policy.OutputDrop.String()
	if f.config.IP4.IcmpOut == false {
		if err := chain.AddRule("oifname != \"lo\" ip protocol icmp icmp type echo-request counter " + drop); err != nil {
			return err
		}
		return f.reloadOutputICMPAfter()
	}

	if f.config.IP4.IcmpOutRate == "0" {
		return f.reloadOutputICMPAfter()
	}

	if err := chain.AddRule("oifname != \"lo\" ip protocol icmp icmp type echo-request limit rate " + f.config.IP4.IcmpInRate + " counter accept"); err != nil {
		return err
	}
	if err := chain.AddRule("oifname != \"lo\" ip protocol icmp icmp type echo-request counter " + drop); err != nil {
		return err
	}

	return f.reloadOutputICMPAfter()
}

func (f *firewall) reloadOutputICMPAfter() error {
	chain := f.chains.Output()

	if f.config.IP4.IcmpTimestampDrop == true {
		drop := f.config.Policy.OutputDrop.String()
		if err := chain.AddRule("oifname != \"lo\" ip protocol icmp icmp type timestamp-request " + drop); err != nil {
			return err
		}
	}

	if err := chain.AddRule("oifname != \"lo\" ip protocol icmp counter accept"); err != nil {
		return err
	}
	return nil
}

func (f *firewall) reloadOutputPorts() error {
	chain := f.chains.Output()
	for _, port := range f.config.OutPorts {
		protocol := port.Port.ProtocolString()
		number := port.Port.NumberString()
		baseRule := "oifname != \"lo\" meta l4proto " + protocol + " ct state new " + protocol + " dport " + number

		if port.LimitRate != "" {
			rule := baseRule + " limit rate " + port.LimitRate + " counter " + port.Action.String()
			if err := chain.AddRule(rule); err != nil {
				return err
			}
			ruleDrop := baseRule + " counter " + f.config.Policy.InputDrop.String()
			if err := chain.AddRule(ruleDrop); err != nil {
				return err
			}
		} else {
			rule := baseRule + " counter " + port.Action.String()
			if err := chain.AddRule(rule); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *firewall) reloadOutputAddIPs() error {
	if err := f.chains.NewLocalOutput(); err != nil {
		return err
	}
	chain := f.chains.LocalOutput()
	if err := chain.AddRuleOut(f.chains.Output().AddRule); err != nil {
		return err
	}

	for _, ipConfig := range f.config.IP4.OutIPs {
		if err := outputAddIP(chain.AddRule, ipConfig, "ip"); err != nil {
			return err
		}
	}

	if !f.config.IP6.Enable {
		return nil
	}

	for _, ipConfig := range f.config.IP6.OutIPs {
		if err := outputAddIP(chain.AddRule, ipConfig, "ip6"); err != nil {
			return err
		}
	}

	return nil
}

func outputAddIP(addRuleFunc func(expr ...string) error, config ConfigIP, ipMatch string) error {

	rule := ipMatch + " daddr " + config.IP + " oifname != \"lo\""
	if !config.OnlyIP {
		rule += " " + config.Port.ProtocolString() + " dport " + config.Port.NumberString()
	}
	if config.LimitRate != "" {
		rule += " limit rate " + config.LimitRate
	}
	rule += " counter " + config.Action.String()
	if err := addRuleFunc(rule); err != nil {
		return err
	}

	return nil
}
