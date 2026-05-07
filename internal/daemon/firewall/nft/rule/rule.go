package rule

import "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/config"

type AddFunc func(expr ...string) error

func InputAddIP(addRuleFunc AddFunc, config config.ConfigIP, ipMatch string) error {
	rule := ipMatch + " saddr " + config.IP + " iifname != \"lo\""
	if !config.OnlyIP {
		rule += " " + config.Port.ProtocolString() + " dport " + config.Port.NumberString()
	}
	if config.LimitRate != "" {
		rule += " limit rate " + config.LimitRate
	}
	rule += " counter " + config.Action.String()
	return addRuleFunc(rule)
}

func OutputAddIP(addRuleFunc AddFunc, config config.ConfigIP, ipMatch string) error {
	rule := ipMatch + " daddr " + config.IP + " oifname != \"lo\""
	if !config.OnlyIP {
		rule += " " + config.Port.ProtocolString() + " dport " + config.Port.NumberString()
	}
	if config.LimitRate != "" {
		rule += " limit rate " + config.LimitRate
	}
	rule += " counter " + config.Action.String()
	return addRuleFunc(rule)
}

func ForwardAddIP(addRuleFunc AddFunc, config config.ConfigIP, ipMatch string) error {
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
