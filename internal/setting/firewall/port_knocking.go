package firewall

import (
	"fmt"
	"strings"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	port2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/ip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

type portKnocking struct {
	Name      string              `mapstructure:"name"`
	IPVersion string              `mapstructure:"ip_version"`
	Port      int                 `mapstructure:"port"`
	Protocol  string              `mapstructure:"protocol"`
	Knocks    []portKnockingKnock `mapstructure:"knock"`
}

func defaultPortKnocking() []portKnocking {
	return []portKnocking{}
}

func (p *portKnocking) ToPortKnocking() (firewall.ConfigPortKnocking, error) {
	if len(p.Knocks) == 0 {
		return firewall.ConfigPortKnocking{}, fmt.Errorf("port knocking must have at least one knock")
	}

	if err := p.validate(); err != nil {
		return firewall.ConfigPortKnocking{}, err
	}

	protocol, err := port2.ToProtocol(p.Protocol)
	if err != nil {
		return firewall.ConfigPortKnocking{}, err
	}

	l4Port, err := types.NewL4Port(uint16(p.Port), protocol)
	if err != nil {
		return firewall.ConfigPortKnocking{}, err
	}

	ipVersion, err := toVersionIP(p.IPVersion)
	if err != nil {
		return firewall.ConfigPortKnocking{}, err
	}

	knocks := make([]*firewall.ConfigKnock, 0, len(p.Knocks))
	for _, knock := range p.Knocks {
		knock, err := knock.ToKnock()
		if err != nil {
			return firewall.ConfigPortKnocking{}, err
		}
		knocks = append(knocks, &knock)
	}

	return firewall.ConfigPortKnocking{
		Name:      p.Name,
		Port:      l4Port,
		IPVersion: ipVersion,
		Knocks:    knocks,
	}, nil
}

func (p *portKnocking) validate() error {
	if err := validate.Name(p.Name, "portKnocking.name"); err != nil {
		return err
	}

	if err := validate.Port(p.Port, "portKnocking.port"); err != nil {
		return err
	}

	return nil
}

func toVersionIP(versionIP string) (port2.Version, error) {
	switch strings.ToLower(versionIP) {
	case "ip4":
		return port2.IPv4, nil
	case "ip6":
		return port2.IPv6, nil
	default:
		return port2.IPv4, fmt.Errorf("invalid version_ip. Must be ip4 or ip6")
	}
}
