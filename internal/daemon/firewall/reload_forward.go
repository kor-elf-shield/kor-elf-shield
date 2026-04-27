package firewall

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
)

func (f *firewall) reloadForward() error {
	f.logger.Debug("Reloading forward chain")
	err := f.chains.NewForward(f.config.MetadataNaming.ChainForwardName, f.config.Policy.DefaultAllowForward, f.config.Policy.ForwardPriority)
	if err != nil {
		return err
	}
	chain := f.chains.Forward()

	if err := f.reloadForwardAddIPs(); err != nil {
		return err
	}

	if f.config.Options.DockerSupport {
		if err := f.docker.NftChains().ForwardFilterJump(chain.AddRule); err != nil {
			return err
		}
	}

	if f.config.Policy.DefaultAllowForward == false {
		drop := f.config.Policy.ForwardDrop.String()
		if err := chain.AddRule(drop); err != nil {
			return err
		}
	}

	return nil
}

func (f *firewall) reloadForwardAddIPs() error {
	if err := f.chains.NewLocalForward(); err != nil {
		return err
	}
	chain := f.chains.LocalForward()
	if err := chain.AddRuleIn(f.chains.Forward().AddRule); err != nil {
		return err
	}

	for _, ipConfig := range f.config.IP4.InIPs {
		if ipConfig.Action != types.ActionDrop && ipConfig.Action != types.ActionReject {
			continue
		}
		if err := forwardAddIP(chain.AddRule, ipConfig, "ip"); err != nil {
			return err
		}
	}

	if !f.config.IP6.Enable {
		return nil
	}

	for _, ipConfig := range f.config.IP6.InIPs {
		if ipConfig.Action != types.ActionDrop && ipConfig.Action != types.ActionReject {
			continue
		}
		if err := forwardAddIP(chain.AddRule, ipConfig, "ip6"); err != nil {
			return err
		}
	}

	return nil
}

func forwardAddIP(addRuleFunc func(expr ...string) error, config config.ConfigIP, ipMatch string) error {
	rule := ipMatch + " saddr " + config.IP + " iifname != \"lo\""

	// There, during routing, the port changes and then the IP blocking rule will not work.
	//if !config.OnlyIP {
	//	rule += " " + config.Protocol.String() + " dport " + strconv.Itoa(int(config.Port))
	//}

	rule += " counter " + config.Action.String()
	if err := addRuleFunc(rule); err != nil {
		return err
	}

	return nil
}
