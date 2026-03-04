package entity

import (
	"errors"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/ip"
)

type Blocking struct {
	IP           string `json:"IP"`
	Ports        []BlockingPort
	ExpireAtUnix int64  `json:"ExpireAtUnix"`
	Reason       string `json:"Reason"`
}

func (b *Blocking) IsPorts() bool {
	return len(b.Ports) > 0
}

func (b *Blocking) ToL4Ports() ([]types.L4Port, error) {
	if !b.IsPorts() {
		return nil, fmt.Errorf("ports is empty")
	}

	l4Ports := make([]types.L4Port, 0, len(b.Ports))
	for _, port := range b.Ports {
		l4port, err := port.ToL4Port()
		if err != nil {
			return nil, err
		}
		l4Ports = append(l4Ports, l4port)
	}

	return l4Ports, nil
}

type BlockingPort struct {
	Number   uint16 `json:"Port"`
	Protocol string `json:"Protocol"`
}

func (p *BlockingPort) ToL4Port() (types.L4Port, error) {
	if p.Protocol == "" {
		return nil, errors.New("protocol is empty")
	}

	protocol, err := ip.ToProtocol(p.Protocol)
	if err != nil {
		return nil, err
	}

	return types.NewL4Port(p.Number, protocol)
}
