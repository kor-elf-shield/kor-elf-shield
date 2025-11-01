package firewall

import (
	"errors"
	"strings"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
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
	if len(p.Numbers) == 0 {
		error = errors.New("invalid port number. must be 0-65535")
		return
	}
	if len(p.Directions) == 0 {
		error = errors.New("invalid direction. Must be in or out")
		return
	}
	if len(p.Protocols) == 0 {
		error = errors.New("invalid protocol. Must be tcp or udp")
		return
	}
	if len(p.Action) == 0 {
		error = errors.New("invalid action. Must be accept, drop or reject")
		return
	}
	action, err := toAction(p.Action)
	if err != nil {
		error = err
		return
	}
	if p.LimitRate != "" {
		if err := validate.NftLimitRate(p.LimitRate, "limit_rate"); err != nil {
			error = err
			return
		}
	}

	for _, port := range p.Numbers {
		if port < 0 || port > 65535 {
			error = errors.New("invalid port number. must be 0-65535")
			return
		}
		for _, direction := range p.Directions {
			addDirection, err := toDirection(direction)
			if err != nil {
				error = err
				return
			}
			for _, protocol := range p.Protocols {
				addProtocol, err := toProtocol(protocol)
				if err != nil {
					error = err
					return
				}

				addPort := firewall.ConfigPort{
					Number:    uint16(port),
					Direction: addDirection,
					Protocol:  addProtocol,
					Action:    action,
					LimitRate: p.LimitRate,
				}
				if addDirection == firewall.DirectionIn {
					InPorts = append(InPorts, addPort)
				} else {
					OutPorts = append(OutPorts, addPort)
				}
			}
		}
	}

	return
}

func toDirection(direction string) (firewall.Direction, error) {
	switch strings.ToLower(direction) {
	case "in":
		return firewall.DirectionIn, nil
	case "out":
		return firewall.DirectionOut, nil
	default:
		return firewall.DirectionIn, errors.New("invalid direction. Must be in or out")
	}
}

func toProtocol(protocol string) (firewall.Protocol, error) {
	switch strings.ToLower(protocol) {
	case "tcp":
		return firewall.ProtocolTCP, nil
	case "udp":
		return firewall.ProtocolUDP, nil
	default:
		return firewall.ProtocolTCP, errors.New("invalid protocol. Must be tcp or udp")
	}
}

func toAction(action string) (firewall.Action, error) {
	switch strings.ToLower(action) {
	case "accept":
		return firewall.ActionAccept, nil
	case "drop":
		return firewall.ActionDrop, nil
	case "reject":
		return firewall.ActionReject, nil
	default:
		return firewall.ActionAccept, errors.New("invalid action. Must be accept, drop or reject")
	}
}
