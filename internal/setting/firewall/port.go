package firewall

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/ip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

type Port struct {
	Numbers    []int    `mapstructure:"numbers"`
	Directions []string `mapstructure:"directions"`
	Protocols  []string `mapstructure:"protocols"`
	Action     string   `mapstructure:"action"`
	LimitRate  string   `mapstructure:"limit_rate"`
}

func defaultPorts() []Port {
	return []Port{}
}

func (p *Port) ToPorts() (InPorts []firewall.ConfigPort, OutPorts []firewall.ConfigPort, error error) {
	if err := p.validate(); err != nil {
		error = err
		return
	}
	action, err := ip.ToAction(p.Action)
	if err != nil {
		error = err
		return
	}

	for _, port := range p.Numbers {
		if err := validate.Port(port, "port"); err != nil {
			error = err
			return
		}
		for _, direction := range p.Directions {
			addDirection, err := ip.ToDirection(direction)
			if err != nil {
				error = err
				return
			}
			for _, protocol := range p.Protocols {
				addProtocol, err := ip.ToProtocol(protocol)
				if err != nil {
					error = err
					return
				}

				l4Port, err := types.NewL4Port(uint16(port), addProtocol)
				if err != nil {
					error = err
					return
				}

				addPort := firewall.ConfigPort{
					Port:      l4Port,
					Action:    action,
					LimitRate: p.LimitRate,
				}
				if addDirection == types.DirectionIn {
					InPorts = append(InPorts, addPort)
				} else {
					OutPorts = append(OutPorts, addPort)
				}
			}
		}
	}

	return
}

func (p *Port) validate() error {
	if len(p.Numbers) == 0 {
		return errors.New("invalid port number. must be 0-65535")
	}
	if len(p.Directions) == 0 {
		return errors.New("invalid direction. Must be in or out")
	}
	if len(p.Protocols) == 0 {
		return errors.New("invalid protocol. Must be tcp or udp")
	}
	if len(p.Action) == 0 {
		return errors.New("invalid action. Must be accept, drop or reject")
	}
	if p.LimitRate != "" {
		if err := validate.NftLimitRate(p.LimitRate, "limit_rate"); err != nil {
			return err
		}
	}
	return nil
}
