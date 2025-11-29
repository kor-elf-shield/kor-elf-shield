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
	InputPriority       int    `mapstructure:"input_priority"`
	OutputDrop          string `mapstructure:"output_drop"`
	OutputPriority      int    `mapstructure:"output_priority"`
	ForwardDrop         string `mapstructure:"forward_drop"`
	ForwardPriority     int    `mapstructure:"forward_priority"`
}

func defaultPolicy() policy {
	return policy{
		DefaultAllowInput:   false,
		DefaultAllowOutput:  false,
		DefaultAllowForward: false,
		InputDrop:           "drop",
		InputPriority:       -10,
		OutputDrop:          "reject",
		OutputPriority:      -10,
		ForwardDrop:         "drop",
		ForwardPriority:     -10,
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
		InputPriority:       p.InputPriority,
		OutputDrop:          outputDrop,
		OutputPriority:      p.OutputPriority,
		ForwardDrop:         forwardDrop,
		ForwardPriority:     p.ForwardPriority,
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

func (p policy) Validate() error {
	if err := validateDrop(p.InputDrop, "input_drop"); err != nil {
		return err
	}
	if err := validatePriority(p.InputPriority, "input_priority"); err != nil {
		return err
	}

	if err := validateDrop(p.OutputDrop, "output_drop"); err != nil {
		return err
	}
	if err := validatePriority(p.OutputPriority, "output_priority"); err != nil {
		return err
	}

	if err := validateDrop(p.ForwardDrop, "forward_drop"); err != nil {
		return err
	}
	if err := validatePriority(p.ForwardPriority, "forward_priority"); err != nil {
		return err
	}
	return nil
}

func validateDrop(drop string, parameterName string) error {
	switch drop {
	case "drop", "reject":
		return nil
	}
	return fmt.Errorf("invalid %s. Must be drop or reject", parameterName)
}

func validatePriority(priority int, parameterName string) error {
	if priority < -50 || priority > 50 {
		return fmt.Errorf("%s must be in range -50-50", parameterName)
	}
	return nil
}
