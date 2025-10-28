package firewall

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
)

type policy struct {
	DefaultAllowInput   bool   `mapstructure:"default_allow_input"`
	DefaultAllowOutput  bool   `mapstructure:"default_allow_output"`
	DefaultAllowForward bool   `mapstructure:"default_allow_forward"`
	InputDrop           string `mapstructure:"input_drop"`
	OutputDrop          string `mapstructure:"output_drop"`
	ForwardDrop         string `mapstructure:"forward_drop"`
}

func defaultPolicy() policy {
	return policy{
		DefaultAllowInput:   false,
		DefaultAllowOutput:  false,
		DefaultAllowForward: false,
		InputDrop:           "drop",
		OutputDrop:          "reject",
		ForwardDrop:         "drop",
	}
}

func (p policy) ToConfigPolicy() (firewall.ConfigPolicy, error) {
	inputDrop, err := p.dropToPolicyDrop(p.InputDrop, "input_drop")
	if err != nil {
		return firewall.ConfigPolicy{}, err
	}

	outputDrop, err := p.dropToPolicyDrop(p.OutputDrop, "output_drop")
	if err != nil {
		return firewall.ConfigPolicy{}, err
	}

	forwardDrop, err := p.dropToPolicyDrop(p.ForwardDrop, "forward_drop")
	if err != nil {
		return firewall.ConfigPolicy{}, err
	}

	return firewall.ConfigPolicy{
		DefaultAllowInput:   p.DefaultAllowInput,
		DefaultAllowOutput:  p.DefaultAllowOutput,
		DefaultAllowForward: p.DefaultAllowForward,
		InputDrop:           inputDrop,
		OutputDrop:          outputDrop,
		ForwardDrop:         forwardDrop,
	}, nil
}

func (p policy) dropToPolicyDrop(drop string, parametrName string) (firewall.PolicyDrop, error) {
	if drop == "" {
		return 0, fmt.Errorf("%s is empty", parametrName)
	}
	switch drop {
	case "drop":
		return firewall.Drop, nil
	case "reject":
		return firewall.Reject, nil
	default:
		return 0, fmt.Errorf("invalid %s . Must be drop or reject", parametrName)
	}
}
