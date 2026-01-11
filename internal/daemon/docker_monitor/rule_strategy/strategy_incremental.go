package rule_strategy

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
	nftChain "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type incrementalStrategy struct {
	dockerClient client.Docker
	chains       chain.Chains
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

func (i *incrementalStrategy) Reload(newNoneChain func(chain string) (nftChain.Chain, error)) error {
	chains, err := chain.NewChains(newNoneChain)
	if err != nil {
		return err
	}
	i.chains = chains

	i.generator.GenerateAll(i.chains, true)

	return nil
}

func (i *incrementalStrategy) Chains() chain.Chains {
	return i.chains
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
			i.eventContainerStop(event.ID)
			return
		}

		return
	}

	if event.Type == "network" {
		if event.Action == "create" {
			if err := i.eventNetworkCreate(event.ID); err != nil {
				i.logger.Error(fmt.Sprintf("failed to handle network create event: %s", err))
			}
		}

		if event.Action == "destroy" {
			i.eventNetworkDestroy(event.ID)
		}

		return
	}
}

func (i *incrementalStrategy) eventContainerStart(containerId string) error {
	container, err := i.dockerClient.FetchContainer(containerId)
	if err != nil {
		return err
	}

	for _, ipInfo := range container.Networks.IPAddresses {
		bridge, err := i.dockerClient.FetchBridge(ipInfo.NetworkID)
		if err != nil {
			i.logger.Error(fmt.Sprintf("failed to fetch bridge for container %s: %s", containerId, err))
			continue
		}
		i.generator.GenerateContainer(container, bridge.Name, i.chains, true)
	}

	return nil
}

func (i *incrementalStrategy) eventContainerStop(containerId string) {
	listChains := i.chains.List()

	if err := i.nftRuleDeleteContainer(containerId, &listChains.PreroutingFilter); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete container %s rules: %s", containerId, err))
	}

	if err := i.nftRuleDeleteContainer(containerId, &listChains.DockerNat); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete container %s rules: %s", containerId, err))
	}

	if err := i.nftRuleDeleteContainer(containerId, &listChains.DockerFilterFirst); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete container %s rules: %s", containerId, err))
	}
}

func (i *incrementalStrategy) nftRuleDeleteContainer(containerId string, chain *chain.Data) error {
	rules, err := chain.ListRules()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if rule.Comment != "container_id:"+containerId {
			continue
		}
		if err := chain.RemoveRuleByHandle(rule.Handle); err != nil {
			i.logger.Error(fmt.Sprintf("failed to delete container %s rule: %s", containerId, err))
		}
	}

	return nil
}

func (i *incrementalStrategy) eventNetworkCreate(bridgeId string) error {
	bridge, err := i.dockerClient.FetchBridge(bridgeId)
	if err != nil {
		return err
	}

	i.generator.GenerateBridge(bridge, i.chains, true)
	return nil
}

func (i *incrementalStrategy) eventNetworkDestroy(bridgeId string) {
	listChains := i.chains.List()

	if err := i.nftRuleDeleteBridge(bridgeId, &listChains.DockerFilterSecond); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete bridge %s rules: %s", bridgeId, err))
	}

	if err := i.nftRuleDeleteBridge(bridgeId, &listChains.ForwardFilter); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete bridge %s rules: %s", bridgeId, err))
	}

	if err := i.nftRuleDeleteBridge(bridgeId, &listChains.ForwardBridge); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete bridge %s rules: %s", bridgeId, err))
	}

	if err := i.nftRuleDeleteBridge(bridgeId, &listChains.ForwardCT); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete bridge %s rules: %s", bridgeId, err))
	}

	if err := i.nftRuleDeleteBridge(bridgeId, &listChains.PostroutingNat); err != nil {
		i.logger.Error(fmt.Sprintf("failed to delete bridge %s rules: %s", bridgeId, err))
	}
}

func (i *incrementalStrategy) nftRuleDeleteBridge(bridgeId string, chain *chain.Data) error {
	rules, err := chain.ListRules()
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if rule.Comment != "bridge_id:"+bridgeId {
			continue
		}
		if err := chain.RemoveRuleByHandle(rule.Handle); err != nil {
			i.logger.Error(fmt.Sprintf("failed to delete bridge %s rule: %s", bridgeId, err))
		}
	}

	return nil
}
