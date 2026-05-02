package rule_strategy

import (
	"fmt"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type incrementalStrategy struct {
	dockerClient client.Docker
	nftDocker    firewall.NFTDocker
	generator    Generator
	logger       log.Logger
}

func NewIncrementalStrategy(generator Generator, dockerClient client.Docker, logger log.Logger) Strategy {
	return &incrementalStrategy{
		dockerClient: dockerClient,
		generator:    generator,
		logger:       logger,
	}
}

func (i *incrementalStrategy) Reload(nftDocker firewall.NFTDocker) error {
	i.nftDocker = nftDocker

	batchBuilder, err := i.nftDocker.NFT().NewBuildBatch()
	if err != nil {
		return err
	}
	defer func() {
		if err := batchBuilder.Close(); err != nil {
			i.logger.Warn(err.Error())
		}
	}()

	i.generator.GenerateAll(batchBuilder, i.nftDocker.Chains(), true)

	return i.nftDocker.NFT().RunBatch(batchBuilder)
}

func (i *incrementalStrategy) Chains() firewall.NFTDockerChains {
	return i.nftDocker.Chains()
}

func (i *incrementalStrategy) Event(event *client.Event) {
	if event == nil || event.ID == "" {
		return
	}

	if event.Type == "container" {
		if event.Action == "start" {
			if err := i.eventContainerStart(event.ID); err != nil {
				i.logger.Error(fmt.Sprintf("failed to handle container start event: %s", err))
			}
			return
		}

		if event.Action == "die" {
			if err := i.eventContainerStop(event.ID); err != nil {
				i.logger.Error(fmt.Sprintf("failed to handle container stop event: %s", err))
			}
			return
		}

		return
	}

	if event.Type == "network" {
		if event.Action == "create" {
			if err := i.eventNetworkCreate(event.ID); err != nil {
				i.logger.Error(fmt.Sprintf("failed to handle network create event: %s", err))
			}
			return
		}

		if event.Action == "destroy" {
			if err := i.eventNetworkDestroy(event.ID); err != nil {
				i.logger.Error(fmt.Sprintf("failed to handle network destroy event: %s", err))
			}
			return
		}

		return
	}
}

func (i *incrementalStrategy) eventContainerStart(containerId string) error {
	container, err := i.dockerClient.FetchContainer(containerId)
	if err != nil {
		return err
	}

	batchBuilder, err := i.nftDocker.NFT().NewBuildBatch()
	if err != nil {
		return err
	}
	defer func() {
		if err := batchBuilder.Close(); err != nil {
			i.logger.Warn(err.Error())
		}
	}()

	for _, ipInfo := range container.Networks.IPAddresses {
		bridge, err := i.dockerClient.FetchBridge(ipInfo.NetworkID)
		if err != nil {
			i.logger.Error(fmt.Sprintf("failed to fetch bridge for container %s: %s", containerId, err))
			continue
		}
		i.generator.GenerateContainer(container, bridge.Name, batchBuilder, i.nftDocker.Chains(), true)
	}

	return i.nftDocker.NFT().RunBatch(batchBuilder)
}

func (i *incrementalStrategy) eventContainerStop(containerId string) error {
	batchBuilder, err := i.nftDocker.NFT().NewBuildBatch()
	if err != nil {
		return err
	}
	defer func() {
		if err := batchBuilder.Close(); err != nil {
			i.logger.Warn(err.Error())
		}
	}()

	if err := i.nftRuleDeleteContainer(containerId, batchBuilder, i.nftDocker.Chains().PreroutingFilter()); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete container %s rules: %s", containerId, err))
	}

	if err := i.nftRuleDeleteContainer(containerId, batchBuilder, i.nftDocker.Chains().DockerNat()); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete container %s rules: %s", containerId, err))
	}

	if err := i.nftRuleDeleteContainer(containerId, batchBuilder, i.nftDocker.Chains().DockerFilterFirst()); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete container %s rules: %s", containerId, err))
	}

	return i.nftDocker.NFT().RunBatch(batchBuilder)
}

func (i *incrementalStrategy) nftRuleDeleteContainer(containerId string, builder nft.BatchBuilder, chain chain.Docker) error {
	rules, err := chain.ListRules()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if rule.Comment != "container_id:"+containerId {
			continue
		}
		if err := chain.RemoveRuleByHandle(builder, rule.Handle); err != nil {
			i.logger.Error(fmt.Sprintf("failed to delete container %s rule: %s", containerId, err))
		}
	}

	return nil
}

func (i *incrementalStrategy) eventNetworkCreate(bridgeId string) error {
	batchBuilder, err := i.nftDocker.NFT().NewBuildBatch()
	if err != nil {
		return err
	}
	defer func() {
		if err := batchBuilder.Close(); err != nil {
			i.logger.Warn(err.Error())
		}
	}()

	bridge, err := i.dockerClient.FetchBridge(bridgeId)
	if err != nil {
		return err
	}

	i.generator.GenerateBridge(bridge, batchBuilder, i.nftDocker.Chains(), true)
	return i.nftDocker.NFT().RunBatch(batchBuilder)
}

func (i *incrementalStrategy) eventNetworkDestroy(bridgeId string) error {
	batchBuilder, err := i.nftDocker.NFT().NewBuildBatch()
	if err != nil {
		return err
	}
	defer func() {
		if err := batchBuilder.Close(); err != nil {
			i.logger.Warn(err.Error())
		}
	}()

	if err := i.nftRuleDeleteBridge(bridgeId, batchBuilder, i.nftDocker.Chains().DockerFilterSecond()); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete bridge %s rules: %s", bridgeId, err))
	}

	if err := i.nftRuleDeleteBridge(bridgeId, batchBuilder, i.nftDocker.Chains().ForwardFilter()); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete bridge %s rules: %s", bridgeId, err))
	}

	if err := i.nftRuleDeleteBridge(bridgeId, batchBuilder, i.nftDocker.Chains().ForwardBridge()); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete bridge %s rules: %s", bridgeId, err))
	}

	if err := i.nftRuleDeleteBridge(bridgeId, batchBuilder, i.nftDocker.Chains().ForwardCT()); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete bridge %s rules: %s", bridgeId, err))
	}

	if err := i.nftRuleDeleteBridge(bridgeId, batchBuilder, i.nftDocker.Chains().PostroutingNat()); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete bridge %s rules: %s", bridgeId, err))
	}

	return i.nftDocker.NFT().RunBatch(batchBuilder)
}

func (i *incrementalStrategy) nftRuleDeleteBridge(bridgeId string, builder nft.BatchBuilder, chain chain.Docker) error {
	rules, err := chain.ListRules()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if rule.Comment != "bridge_id:"+bridgeId {
			continue
		}
		if err := chain.RemoveRuleByHandle(builder, rule.Handle); err != nil {
			i.logger.Error(fmt.Sprintf("failed to delete bridge %s rule: %s", bridgeId, err))
		}
	}

	return nil
}
