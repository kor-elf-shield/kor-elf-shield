package rule_strategy

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
	nftChain "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain"
)

type rebuildStrategy struct {
	chains    chain.Chains
	generator Generator
}

func NewRebuildStrategy(generator Generator) Strategy {
	return &rebuildStrategy{
		generator: generator,
	}
}

func (r *rebuildStrategy) Reload(newNoneChain func(chain string) (nftChain.Chain, error)) error {
	chains, err := chain.NewChains(newNoneChain)
	if err != nil {
		return err
	}
	r.chains = chains

	r.generator.GenerateAll(r.chains, false)

	return nil
}

func (r *rebuildStrategy) Chains() chain.Chains {
	return r.chains
}

func (r *rebuildStrategy) Event(event *client.Event) {
	if event == nil || event.Type != "container" {
		return
	}

	r.generator.ClearChains(r.chains)
	r.generator.GenerateAll(r.chains, false)
}
