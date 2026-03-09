package brute_force_protection

import "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"

type Block interface {
	PortsBlocked() (bool, []types.L4Port)
}

type block struct {
	shouldPortsBlocked bool
	ports              []types.L4Port
}

func NewBlockOnceIPConfig() Block {
	return &block{
		shouldPortsBlocked: false,
		ports:              nil,
	}
}

func NewBlockIPAndPortsConfig(ports []types.L4Port) Block {
	return &block{
		shouldPortsBlocked: true,
		ports:              ports,
	}
}

func (b *block) PortsBlocked() (bool, []types.L4Port) {
	if !b.shouldPortsBlocked {
		return false, nil
	}

	return true, b.ports
}
