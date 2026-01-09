package rule_strategy

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
	nftChain "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain"
)

type RebuildStrategy struct {
	chains    chain.Chains
	generator Generator
}

func NewRebuildStrategy(generator Generator) Strategy {
	return &RebuildStrategy{
		generator: generator,
	}
}

func (r *RebuildStrategy) Reload(newNoneChain func(chain string) (nftChain.Chain, error)) error {
	chains, err := chain.NewChains(newNoneChain)
	if err != nil {
		return err
	}
	r.chains = chains

	r.generator.GenerateAll(r.chains)

	return nil
}

func (r *RebuildStrategy) Chains() chain.Chains {
	return r.chains
}

func (r *RebuildStrategy) Event(_ *client.Event) {
	r.generator.ClearChains(r.chains)
	r.generator.GenerateAll(r.chains)
}
