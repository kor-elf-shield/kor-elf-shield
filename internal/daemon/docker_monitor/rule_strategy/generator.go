package rule_strategy

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Generator interface {
	GenerateAll(chains chain.Chains, isComment bool)
	GenerateBridge(bridge client.Bridge, chain chain.Chains, isComment bool)
	GenerateContainer(container client.Container, bridgeName string, chain chain.Chains, isComment bool)
	ClearChains(chains chain.Chains)
	AddRule(chainData chain.Data, rule string)
}

type generator struct {
	dockerClient client.Docker
	logger       log.Logger
}

func NewGenerator(dockerClient client.Docker, logger log.Logger) Generator {
	return &generator{
		dockerClient: dockerClient,
		logger:       logger,
	}
}

func (g *generator) GenerateAll(chains chain.Chains, isComment bool) {
	listChains := chains.List()

	if err := listChains.ForwardCT.JumpTo(&listChains.ForwardFilter, "", ""); err != nil {
		g.logger.Error(err.Error())
	}
	if err := listChains.ForwardBridge.JumpTo(&listChains.ForwardFilter, "", ""); err != nil {
		g.logger.Error(err.Error())
	}
	if err := listChains.DockerFilterFirst.JumpTo(&listChains.DockerFilter, "", ""); err != nil {
		g.logger.Error(err.Error())
	}
	if err := listChains.DockerFilterSecond.JumpTo(&listChains.DockerFilter, "", ""); err != nil {
		g.logger.Error(err.Error())
	}

	bridges, err := g.dockerClient.FetchBridges()
	if err != nil {
		g.logger.Error(err.Error())
		return
	}

	for _, bridge := range bridges {
		g.GenerateBridge(bridge, chains, isComment)

		if bridge.Containers == nil {
			continue
		}
		for _, container := range bridge.Containers {
			g.GenerateContainer(container, bridge.Name, chains, isComment)
		}
	}
}

func (g *generator) GenerateBridge(bridge client.Bridge, chain chain.Chains, isComment bool) {
	listChains := chain.List()

	var rule string
	comment := ""
	if isComment {
		comment = fmt.Sprintf("comment \"bridge_id:%s\"", bridge.ID)
	}

	rule = fmt.Sprintf("iifname != \"%s\" oifname \"%s\" counter drop %s", bridge.Name, bridge.Name, comment)
	g.AddRule(listChains.DockerFilterSecond, rule)

	rule = fmt.Sprintf("iifname \"%s\" counter accept %s", bridge.Name, comment)
	g.AddRule(listChains.ForwardFilter, rule)

	rule = fmt.Sprintf("oifname \"%s\" counter", bridge.Name)
	if err := listChains.DockerFilter.JumpTo(&listChains.ForwardBridge, rule, comment); err != nil {
		g.logger.Error(err.Error())
	}

	rule = fmt.Sprintf("oifname \"%s\" ct state related,established counter accept %s", bridge.Name, comment)
	g.AddRule(listChains.ForwardCT, rule)

	for _, subnet := range bridge.Subnets {
		rule = fmt.Sprintf("ip saddr %s oifname != \"%s\" counter masquerade %s", subnet, bridge.Name, comment)
		g.AddRule(listChains.PostroutingNat, rule)
	}
}

func (g *generator) GenerateContainer(container client.Container, bridgeName string, chain chain.Chains, isComment bool) {
	listChains := chain.List()
	var rule string
	comment := ""
	if isComment {
		comment = fmt.Sprintf("comment \"container_id:%s\"", container.ID)
	}

	for _, ipInfo := range container.Networks.IPAddresses {
		rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" counter drop %s", ipInfo.NftPrefix(), ipInfo.Address, bridgeName, comment)
		g.AddRule(listChains.PreroutingFilter, rule)

		for _, port := range container.Networks.Ports {
			isZeroAddress := false
			for _, hostInfo := range port.HostPort {
				if hostInfo.IP.Address != "0.0.0.0" && hostInfo.IP.Address != "::" && (hostInfo.IP.Address == "127.0.0.1" || hostInfo.IP.Address == "::1") {
					rule = fmt.Sprintf("%s daddr %s iifname != \"lo\" %s dport %s counter drop %s", hostInfo.IP.NftPrefix(), hostInfo.IP.Address, port.Protocol, hostInfo.Port, comment)
					g.AddRule(listChains.PreroutingFilter, rule)
				}

				if hostInfo.IP.Address == "0.0.0.0" || hostInfo.IP.Address == "::" {
					if isZeroAddress {
						continue
					}
					isZeroAddress = true
					rule = fmt.Sprintf("iifname != \"%s\" %s dport %s counter dnat %s to %s:%s %s", bridgeName, port.Protocol, hostInfo.Port, ipInfo.NftPrefix(), ipInfo.Address, port.Port, comment)
					g.AddRule(listChains.DockerNat, rule)

					rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" oifname \"%s\" %s dport %s counter accept %s", ipInfo.NftPrefix(), ipInfo.Address, bridgeName, bridgeName, port.Protocol, port.Port, comment)
					g.AddRule(listChains.DockerFilterFirst, rule)
					continue
				}
				rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" oifname \"%s\" %s dport %s counter accept %s", ipInfo.NftPrefix(), ipInfo.Address, bridgeName, bridgeName, port.Protocol, port.Port, comment)
				g.AddRule(listChains.DockerFilterFirst, rule)

				rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" %s dport %s counter dnat to %s:%s %s", hostInfo.IP.NftPrefix(), hostInfo.IP.Address, bridgeName, port.Protocol, hostInfo.Port, ipInfo.Address, port.Port, comment)
				g.AddRule(listChains.DockerNat, rule)
			}
		}
	}
}

func (g *generator) ClearChains(chains chain.Chains) {
	listChains := chains.List()

	if err := listChains.DockerNat.Clear(); err != nil {
		g.logger.Error(err.Error())
	}
	if err := listChains.PostroutingNat.Clear(); err != nil {
		g.logger.Error(err.Error())
	}
	if err := listChains.PreroutingFilter.Clear(); err != nil {
		g.logger.Error(err.Error())
	}
	if err := listChains.DockerFilter.Clear(); err != nil {
		g.logger.Error(err.Error())
	}
	if err := listChains.DockerFilterFirst.Clear(); err != nil {
		g.logger.Error(err.Error())
	}
	if err := listChains.DockerFilterSecond.Clear(); err != nil {
		g.logger.Error(err.Error())
	}
	if err := listChains.ForwardFilter.Clear(); err != nil {
		g.logger.Error(err.Error())
	}

	if err := listChains.ForwardBridge.Clear(); err != nil {
		g.logger.Error(err.Error())
	}

	if err := listChains.ForwardCT.Clear(); err != nil {
		g.logger.Error(err.Error())
	}
}

func (g *generator) AddRule(chainData chain.Data, rule string) {
	if err := chainData.AddRule(rule); err != nil {
		g.logger.Error(err.Error())
	}
}
