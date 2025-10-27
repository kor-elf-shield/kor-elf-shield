package firewall

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
)

type policy struct {
	Input   string `mapstructure:"input"`
	Output  string `mapstructure:"output"`
	Forward string `mapstructure:"forward"`
}

func defaultPolicy() policy {
	return policy{
		Input:   "drop",
		Output:  "reject",
		Forward: "drop",
	}
}

func (p policy) ToConfigPolicy() (firewall.ConfigPolicy, error) {
	input, err := p.input()
	if err != nil {
		return firewall.ConfigPolicy{}, err
	}

	output, err := p.output()
	if err != nil {
		return firewall.ConfigPolicy{}, err
	}

	forward, err := p.forward()
	if err != nil {
		return firewall.ConfigPolicy{}, err
	}

	return firewall.ConfigPolicy{
		Input:   input,
		Output:  output,
		Forward: forward,
	}, nil
}

func (p policy) input() (firewall.Policy, error) {
	if p.Input == "" {
		return 0, errors.New("input policy is empty")
	}
	switch p.Input {
	case "drop":
		return firewall.PolicyDrop, nil
	case "reject":
		return firewall.PolicyReject, nil
	case "accept":
		return firewall.PolicyAccept, nil
	default:
		return 0, errors.New("invalid input policy. Must be drop, reject or accept")
	}
}

func (p policy) output() (firewall.Policy, error) {
	if p.Output == "" {
		return 0, errors.New("output policy is empty")
	}
	switch p.Output {
	case "drop":
		return firewall.PolicyDrop, nil
	case "reject":
		return firewall.PolicyReject, nil
	case "accept":
		return firewall.PolicyAccept, nil
	default:
		return 0, errors.New("invalid output policy. Must be drop, reject or accept")
	}
}

func (p policy) forward() (firewall.Policy, error) {
	if p.Forward == "" {
		return 0, errors.New("forward policy is empty")
	}
	switch p.Forward {
	case "drop":
		return firewall.PolicyDrop, nil
	case "reject":
		return firewall.PolicyReject, nil
	case "accept":
		return firewall.PolicyAccept, nil
	default:
		return 0, errors.New("invalid forward policy. Must be drop, reject or accept")
	}
}
