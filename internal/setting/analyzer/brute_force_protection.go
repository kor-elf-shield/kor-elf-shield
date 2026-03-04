package analyzer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/ip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

type BruteForceProtection struct {
	Enabled              bool   `mapstructure:"enabled"`
	Notify               bool   `mapstructure:"notify"`
	RateLimitCount       int    `mapstructure:"rate_limit_count"`
	RateLimitPeriod      int    `mapstructure:"rate_limit_period"`
	RateLimitResetPeriod int    `mapstructure:"rate_limit_reset_period"`
	BlockingTime         int    `mapstructure:"blocking_time"`
	SSHEnable            bool   `mapstructure:"ssh_enable"`
	SSHNotify            bool   `mapstructure:"ssh_notify"`
	SSHGroup             string `mapstructure:"ssh_group"`

	Groups []BruteForceProtectionGroup
	Rules  []BruteForceProtectionRule
}

func defaultBruteForceProtection() BruteForceProtection {
	return BruteForceProtection{
		Enabled:              true,
		Notify:               true,
		RateLimitCount:       5,
		RateLimitPeriod:      3600,
		RateLimitResetPeriod: 86400,
		BlockingTime:         3600,
		SSHEnable:            true,
		SSHNotify:            true,
		SSHGroup:             "",

		Groups: []BruteForceProtectionGroup{},
		Rules:  []BruteForceProtectionRule{},
	}
}

func (p *BruteForceProtection) Validate() error {
	if p.RateLimitPeriod <= 0 {
		return errors.New("rate limit period must be greater than 0")
	}

	if p.RateLimitCount <= 0 {
		return errors.New("rate limit count must be greater than 0")
	}

	if p.RateLimitResetPeriod < 0 {
		return errors.New("rate limit reset period must be positive")
	}

	if p.BlockingTime < 0 {
		return errors.New("blocking time must be positive")
	}
	return nil
}

func (p *BruteForceProtection) ToSources() ([]*config.Source, error) {
	var sources []*config.Source

	if !p.Enabled {
		return sources, nil
	}

	groups, err := p.groups()
	if err != nil {
		return nil, err
	}

	if p.SSHEnable {
		sshGroup := "_default"
		if p.SSHGroup != "" {
			if _, ok := groups[p.SSHGroup]; !ok {
				return nil, errors.New("ssh group not found")
			}
			sshGroup = p.SSHGroup
		}
		sshSources, err := config.NewBruteForceProtectionSSH(p.Notify && p.SSHNotify, groups[sshGroup])
		if err != nil {
			return nil, err
		}
		sources = append(sources, sshSources...)
	}

	for _, rule := range p.Rules {
		if !rule.Enabled {
			continue
		}

		var group *brute_force_protection.Group
		groupName := "_default"
		if rule.Group != "" {
			groupName = rule.Group
		}
		if _, ok := groups[groupName]; !ok {
			return nil, fmt.Errorf("group %q not found", rule.Group)
		}
		group = groups[groupName]

		source, err := rule.ToSource(p.Notify, group)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (p *BruteForceProtection) groups() (map[string]*brute_force_protection.Group, error) {
	groups := make(map[string]*brute_force_protection.Group)

	groups["_default"] = &brute_force_protection.Group{
		Name:    "_default",
		Message: i18n.Lang.T("alert.bruteForceProtection.group._default.message"),
		RateLimits: []brute_force_protection.RateLimit{
			{
				Count:               uint32(p.RateLimitCount),
				Period:              uint32(p.RateLimitPeriod),
				BlockingTimeSeconds: uint32(p.BlockingTime),
				BlockConfig:         brute_force_protection.NewBlockOnceIPConfig(),
			},
		},
		RateLimitResetPeriod: uint32(p.RateLimitResetPeriod),
	}

	for _, group := range p.Groups {
		g, err := group.ToGroup()
		if err != nil {
			return nil, err
		}
		groups[g.Name] = g
	}

	return groups, nil
}

func toBlockConfigBySettings(blockType string, ports []string) (brute_force_protection.Block, error) {
	if blockType == "" {
		return nil, errors.New("block type is empty")
	}

	switch blockType {
	case "ip":
		return brute_force_protection.NewBlockOnceIPConfig(), nil
	case "ip_port":
		if len(ports) == 0 {
			return nil, errors.New("ports is empty")
		}

		var blockPorts []types.L4Port
		for _, port := range ports {
			l4Port, err := toL4Port(port)
			if err != nil {
				return nil, err
			}
			blockPorts = append(blockPorts, l4Port)
		}

		return brute_force_protection.NewBlockIPAndPortsConfig(blockPorts), nil
	}

	return nil, errors.New("unknown block type")
}

func toL4Port(portString string) (types.L4Port, error) {
	if portString == "" {
		return nil, errors.New("port is empty")
	}

	data := strings.Split(portString, "/")
	protocol := types.ProtocolTCP
	port, err := strconv.Atoi(data[0])
	if err != nil {
		return nil, err
	}
	if err := validate.Port(port, "port"); err != nil {
		return nil, err
	}

	if len(data) == 2 {
		protocol, err = ip.ToProtocol(data[1])
		if err != nil {
			return nil, err
		}
	}

	return types.NewL4Port(uint16(port), protocol)
}
