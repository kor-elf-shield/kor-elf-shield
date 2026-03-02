package firewall

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	port2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/ip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

type IP struct {
	IPs        []string `mapstructure:"ips"`
	Action     string   `mapstructure:"action"`
	Directions []string `mapstructure:"directions"`
	Protocols  []string `mapstructure:"protocols"`
	Ports      []int    `mapstructure:"ports"`
	LimitRate  string   `mapstructure:"limit_rate"`
}

func defaultIPs() []IP {
	return []IP{}
}

type IPs struct {
	InIP4  []firewall.ConfigIP
	OutIP4 []firewall.ConfigIP

	InIP6  []firewall.ConfigIP
	OutIP6 []firewall.ConfigIP
}

func (i *IP) ToIPs() (IPs IPs, error error) {
	if err := i.validate(); err != nil {
		error = err
		return
	}
	action, err := port2.ToAction(i.Action)
	if err != nil {
		error = err
		return
	}

	for _, ip := range i.IPs {
		ipNet, ipVersion, err := port2.DetermineIPVersion(ip)
		if err != nil {
			error = err
			return
		}

		baseConfigIP := firewall.ConfigIP{
			IP:        ipNet,
			Action:    action,
			LimitRate: i.LimitRate,
		}
		in, out, err := loopIP(baseConfigIP, i.Directions, i.Protocols, i.Ports)
		if err != nil {
			error = err
			return
		}
		if ipVersion == port2.IPv4 {
			IPs.InIP4 = append(IPs.InIP4, in...)
			IPs.OutIP4 = append(IPs.OutIP4, out...)
			continue
		}

		IPs.InIP6 = append(IPs.InIP6, in...)
		IPs.OutIP6 = append(IPs.OutIP6, out...)
	}

	return
}

func (i *IP) validate() error {
	if len(i.IPs) == 0 {
		return errors.New("ips is required")
	}
	if len(i.Directions) == 0 {
		return errors.New("invalid direction. Must be in or out")
	}
	if len(i.Action) == 0 {
		return errors.New("invalid action. Must be accept, drop or reject")
	}
	if i.LimitRate != "" {
		if err := validate.NftLimitRate(i.LimitRate, "limit_rate"); err != nil {
			return err
		}
	}
	return nil
}

func loopIP(baseConfigIP firewall.ConfigIP, directions []string, protocols []string, ports []int) (in []firewall.ConfigIP, out []firewall.ConfigIP, error error) {
	for _, direction := range directions {
		addDirection, err := port2.ToDirection(direction)
		if err != nil {
			error = err
			return
		}
		addIP := baseConfigIP
		if len(ports) == 0 {
			// If no port is specified, we only allow the IP address to be accepted.
			addIP.OnlyIP = true
			if addDirection == types.DirectionIn {
				in = append(in, addIP)
			} else {
				out = append(out, addIP)
			}
			continue
		}

		if len(protocols) == 0 {
			addIn, addOut, err := loopIPPort(addIP, ports, addDirection, types.ProtocolTCP)
			if err != nil {
				error = err
				return
			}
			if addDirection == types.DirectionIn {
				in = append(in, addIn...)
			} else {
				out = append(out, addOut...)
			}
			continue
		}

		addIn, addOut, err := loopIPProtocol(addIP, protocols, ports, addDirection)
		if err != nil {
			error = err
			return
		}
		if addDirection == types.DirectionIn {
			in = append(in, addIn...)
		} else {
			out = append(out, addOut...)
		}
	}
	return
}

func loopIPProtocol(baseConfigIP firewall.ConfigIP, protocols []string, ports []int, direction types.Direction) (in []firewall.ConfigIP, out []firewall.ConfigIP, error error) {
	for _, protocol := range protocols {
		addProtocol, err := port2.ToProtocol(protocol)
		if err != nil {
			error = err
			return
		}
		addIP := baseConfigIP

		if len(ports) == 0 {
			if direction == types.DirectionIn {
				in = append(in, addIP)
			} else {
				out = append(out, addIP)
			}
			continue
		}

		addIn, addOut, err := loopIPPort(addIP, ports, direction, addProtocol)
		if err != nil {
			error = err
			return
		}
		if direction == types.DirectionIn {
			in = append(in, addIn...)
		} else {
			out = append(out, addOut...)
		}
	}

	return
}

func loopIPPort(baseConfigIP firewall.ConfigIP, ports []int, direction types.Direction, protocol types.Protocol) (in []firewall.ConfigIP, out []firewall.ConfigIP, error error) {
	for _, port := range ports {
		if err := validate.Port(port, "port"); err != nil {
			error = err
			return
		}

		l4Port, err := types.NewL4Port(uint16(port), protocol)
		if err != nil {
			error = err
			return
		}

		addIP := baseConfigIP
		addIP.Port = l4Port
		if direction == types.DirectionIn {
			in = append(in, addIP)
		} else {
			out = append(out, addIP)
		}
	}
	return
}
