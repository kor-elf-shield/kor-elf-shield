package chain

import nftChain "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain"

type Chains interface {
	ForwardFilterJump(addRule func(expr ...string) error) error
	PreroutingFilterJump(addRule func(expr ...string) error) error

	PreroutingNatJump(addRule func(expr ...string) error) error
	OutputNatJump(addRule func(expr ...string) error) error
	PostroutingNatJump(addRule func(expr ...string) error) error

	List() *chains
}

type chains struct {
	ForwardFilter Data
	ForwardBridge Data
	ForwardCT     Data

	PreroutingFilter   Data
	DockerFilter       Data
	DockerFilterFirst  Data
	DockerFilterSecond Data

	DockerNat      Data
	PostroutingNat Data
}

type Data struct {
	chain nftChain.Chain
	name  string
}

func (d *chains) ForwardFilterJump(addRule func(expr ...string) error) error {
	return d.ForwardFilter.Jump(addRule, "")
}

func (d *chains) PreroutingFilterJump(addRule func(expr ...string) error) error {
	return d.PreroutingFilter.Jump(addRule, "")
}

func (d *chains) PreroutingNatJump(addRule func(expr ...string) error) error {
	return d.DockerNat.Jump(addRule, "fib daddr type local counter")
}

func (d *chains) OutputNatJump(addRule func(expr ...string) error) error {
	if err := d.DockerNat.Jump(addRule, "ip daddr != 127.0.0.0/8 fib daddr type local counter"); err != nil {
		return err
	}

	return d.DockerNat.Jump(addRule, "ip6 daddr != ::1 fib daddr type local counter")
}

func (d *chains) PostroutingNatJump(addRule func(expr ...string) error) error {
	return d.PostroutingNat.Jump(addRule, "")
}

func (d *chains) List() *chains {
	return d
}

func (d *Data) Jump(addRule func(expr ...string) error, rule string) error {
	args := []string{rule, "jump", d.name}
	return addRule(args...)
}

func (d *Data) JumpTo(data *Data, rule string, comment string) error {
	args := []string{rule, "jump", d.name, comment}
	return data.AddRule(args...)
}

func (d *Data) AddRule(rule ...string) error {
	return d.chain.AddRule(rule...)
}

func (d *Data) RemoveRuleByHandle(handle uint64) error {
	return d.chain.RemoveRuleByHandle(handle)
}

func (d *Data) ListRules() ([]nftChain.Rule, error) {
	return d.chain.ListRules()
}

func (d *Data) Clear() error {
	return d.chain.Clear()
}
