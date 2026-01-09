package docker_monitor

import (
	"context"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/rule_strategy"
	nftChain "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Docker interface {
	NftReload(newNoneChain func(chain string) (nftChain.Chain, error)) error
	NftChains() chain.Chains
	Run()
	Close() error
}

type docker struct {
	dockerClient client.Docker
	ruleStrategy rule_strategy.Strategy
	logger       log.Logger
	ctx          context.Context
}

func New(config *Config, ctx context.Context, logger log.Logger) (Docker, error) {
	dockerClient := client.NewDocker(config.Path, ctx, logger)
	ruleStrategy, err := newRuleStrategy(config, dockerClient, logger)
	if err != nil {
		return nil, err
	}

	return &docker{
		dockerClient: dockerClient,
		logger:       logger,
		ctx:          ctx,
		ruleStrategy: ruleStrategy,
	}, nil
}

func (d *docker) NftReload(newNoneChain func(chain string) (nftChain.Chain, error)) error {
	return d.ruleStrategy.Reload(newNoneChain)
}

func (d *docker) NftChains() chain.Chains {
	return d.ruleStrategy.Chains()
}

func (d *docker) Run() {
	events := d.dockerClient.Events()
	for {
		select {
		case <-d.ctx.Done():
			return
		case event := <-events:
			if event.Message == "" {
				continue
			}
			d.logger.Debug("Docker event received: " + event.Message)
			d.ruleStrategy.Event(&event)
		}
	}
}

func (d *docker) Close() error {
	return d.dockerClient.EventsClose()
}

func (d *docker) chainCommand(chainData chain.Data, rule string) {
	if err := chainData.AddRule(rule); err != nil {
		d.logger.Error(err.Error())
	}
}
