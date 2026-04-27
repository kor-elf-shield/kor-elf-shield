package firewall

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	port2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/ip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

type portKnockingKnock struct {
	Port     int    `mapstructure:"port"`
	Protocol string `mapstructure:"protocol"`
	Timeout  int32  `mapstructure:"timeout"`
	Action   string `mapstructure:"action"`
}

func (k *portKnockingKnock) ToKnock() (config.ConfigKnock, error) {
	if err := k.validate(); err != nil {
		return config.ConfigKnock{}, err
	}

	protocol, err := port2.ToProtocol(k.Protocol)
	if err != nil {
		return config.ConfigKnock{}, err
	}
	l4Port, err := types.NewL4Port(uint16(k.Port), protocol)
	if err != nil {
		return config.ConfigKnock{}, err
	}

	action, err := port2.ToKnockAction(k.Action)
	if err != nil {
		return config.ConfigKnock{}, err
	}

	return config.ConfigKnock{
		Port:    l4Port,
		Action:  action,
		Timeout: uint32(k.Timeout),
	}, nil
}

func (k *portKnockingKnock) validate() error {
	if err := validate.Port(k.Port, "knock.port"); err != nil {
		return err
	}

	if k.Timeout <= 0 {
		return fmt.Errorf("knock.timeout must be positive")
	}

	return nil
}
