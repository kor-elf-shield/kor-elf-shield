package rule_strategy

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type rebuildStrategy struct {
	nftDocker firewall.NFTDocker
	generator Generator
	logger    log.Logger
}

func NewRebuildStrategy(generator Generator, logger log.Logger) Strategy {
	return &rebuildStrategy{
		generator: generator,
		logger:    logger,
	}
}

func (r *rebuildStrategy) Reload(nftDocker firewall.NFTDocker) error {
	r.nftDocker = nftDocker

	batchBuilder, err := r.nftDocker.NFT().NewBuildBatch()
	if err != nil {
		return err
	}
	defer func() {
		if err := batchBuilder.Close(); err != nil {
			r.logger.Warn(err.Error())
		}
	}()

	r.generator.GenerateAll(batchBuilder, r.nftDocker.Chains(), false)

	return r.nftDocker.NFT().RunBatch(batchBuilder)
}

func (r *rebuildStrategy) Chains() firewall.NFTDockerChains {
	return r.nftDocker.Chains()
}

func (r *rebuildStrategy) Event(event *client.Event) {
	if event == nil || event.Type != "container" {
		return
	}

	batchBuilder, err := r.nftDocker.NFT().NewBuildBatch()
	if err != nil {
		r.logger.Error(err.Error())
		return
	}
	defer func() {
		if err := batchBuilder.Close(); err != nil {
			r.logger.Warn(err.Error())
		}
	}()

	r.generator.ClearChains(batchBuilder, r.nftDocker.Chains())
	r.generator.GenerateAll(batchBuilder, r.nftDocker.Chains(), false)

	if err := r.nftDocker.NFT().RunBatch(batchBuilder); err != nil {
		r.logger.Error(err.Error())
	}
}
