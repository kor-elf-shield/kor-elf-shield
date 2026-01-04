package chain

type emptyChains struct {
}

func NewEmptyChains() Chains {
	return &emptyChains{}
}

func (c *emptyChains) ForwardFilterJump(_ func(expr ...string) error) error {
	return nil
}

func (c *emptyChains) PreroutingFilterJump(_ func(expr ...string) error) error {
	return nil
}

func (c *emptyChains) PreroutingNatJump(_ func(expr ...string) error) error {
	return nil
}

func (c *emptyChains) OutputNatJump(_ func(expr ...string) error) error {
	return nil
}

func (c *emptyChains) PostroutingNatJump(_ func(expr ...string) error) error {
	return nil
}

func (c *emptyChains) List() *chains {
	return &chains{}
}
