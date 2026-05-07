package rule_strategy

import (
	"fmt"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Generator interface {
	GenerateAll(builder nft.BatchBuilder, chains firewall.NFTDockerChains, isComment bool)
	GenerateBridge(bridge client.Bridge, builder nft.BatchBuilder, chains firewall.NFTDockerChains, isComment bool)
	GenerateContainer(container client.Container, bridgeName string, builder nft.BatchBuilder, chains firewall.NFTDockerChains, isComment bool)
	ClearChains(builder nft.BatchBuilder, chains firewall.NFTDockerChains)
	AddRule(builder nft.BatchBuilder, chainDocker chain.Docker, rule string)
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

func (g *generator) GenerateAll(builder nft.BatchBuilder, chains firewall.NFTDockerChains, isComment bool) {
	if err := chains.ForwardCT().JumpTo(builder, chains.ForwardFilter(), "", ""); err != nil {
		g.logger.Error(err.Error())
	}
	if err := chains.ForwardBridge().JumpTo(builder, chains.ForwardFilter(), "", ""); err != nil {
		g.logger.Error(err.Error())
	}
	if err := chains.DockerFilterFirst().JumpTo(builder, chains.DockerFilter(), "", ""); err != nil {
		g.logger.Error(err.Error())
	}
	if err := chains.DockerFilterSecond().JumpTo(builder, chains.DockerFilter(), "", ""); err != nil {
		g.logger.Error(err.Error())
	}

	bridges, err := g.dockerClient.FetchBridges()
	if err != nil {
		g.logger.Error(err.Error())
		return
	}

	for _, bridge := range bridges {
		g.GenerateBridge(bridge, builder, chains, isComment)

		if bridge.Containers == nil {
			continue
		}
		for _, container := range bridge.Containers {
			g.GenerateContainer(container, bridge.Name, builder, chains, isComment)
		}
	}
}

func (g *generator) GenerateBridge(bridge client.Bridge, builder nft.BatchBuilder, chains firewall.NFTDockerChains, isComment bool) {
	var rule string
	comment := ""
	if isComment {
		comment = fmt.Sprintf("comment \"bridge_id:%s\"", bridge.ID)
	}

	rule = fmt.Sprintf("iifname != \"%s\" oifname \"%s\" counter drop %s", bridge.Name, bridge.Name, comment)
	g.AddRule(builder, chains.DockerFilterSecond(), rule)

	rule = fmt.Sprintf("iifname \"%s\" counter accept %s", bridge.Name, comment)
	g.AddRule(builder, chains.ForwardFilter(), rule)

	rule = fmt.Sprintf("oifname \"%s\" counter", bridge.Name)
	if err := chains.DockerFilter().JumpTo(builder, chains.ForwardBridge(), rule, comment); err != nil {
		g.logger.Error(err.Error())
	}

	rule = fmt.Sprintf("oifname \"%s\" ct state related,established counter accept %s", bridge.Name, comment)
	g.AddRule(builder, chains.ForwardCT(), rule)

	for _, subnet := range bridge.Subnets {
		rule = fmt.Sprintf("ip saddr %s oifname != \"%s\" counter masquerade %s", subnet, bridge.Name, comment)
		g.AddRule(builder, chains.PostroutingNat(), rule)
	}
}

func (g *generator) GenerateContainer(container client.Container, bridgeName string, builder nft.BatchBuilder, chains firewall.NFTDockerChains, isComment bool) {
	var rule string
	comment := ""
	if isComment {
		comment = fmt.Sprintf("comment \"container_id:%s\"", container.ID)
	}

	for _, ipInfo := range container.Networks.IPAddresses {
		rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" counter drop %s", ipInfo.NftPrefix(), ipInfo.Address, bridgeName, comment)
		g.AddRule(builder, chains.PreroutingFilter(), rule)

		for _, port := range container.Networks.Ports {
			isZeroAddress := false
			for _, hostInfo := range port.HostPort {
				if hostInfo.IP.Address != "0.0.0.0" && hostInfo.IP.Address != "::" && (hostInfo.IP.Address == "127.0.0.1" || hostInfo.IP.Address == "::1") {
					rule = fmt.Sprintf("%s daddr %s iifname != \"lo\" %s dport %s counter drop %s", hostInfo.IP.NftPrefix(), hostInfo.IP.Address, port.Protocol, hostInfo.Port, comment)
					g.AddRule(builder, chains.PreroutingFilter(), rule)
				}

				if hostInfo.IP.Address == "0.0.0.0" || hostInfo.IP.Address == "::" {
					if isZeroAddress {
						continue
					}
					isZeroAddress = true
					rule = fmt.Sprintf("iifname != \"%s\" %s dport %s counter dnat %s to %s:%s %s", bridgeName, port.Protocol, hostInfo.Port, ipInfo.NftPrefix(), ipInfo.Address, port.Port, comment)
					g.AddRule(builder, chains.DockerNat(), rule)

					rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" oifname \"%s\" %s dport %s counter accept %s", ipInfo.NftPrefix(), ipInfo.Address, bridgeName, bridgeName, port.Protocol, port.Port, comment)
					g.AddRule(builder, chains.DockerFilterFirst(), rule)
					continue
				}
				rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" oifname \"%s\" %s dport %s counter accept %s", ipInfo.NftPrefix(), ipInfo.Address, bridgeName, bridgeName, port.Protocol, port.Port, comment)
				g.AddRule(builder, chains.DockerFilterFirst(), rule)

				rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" %s dport %s counter dnat to %s:%s %s", hostInfo.IP.NftPrefix(), hostInfo.IP.Address, bridgeName, port.Protocol, hostInfo.Port, ipInfo.Address, port.Port, comment)
				g.AddRule(builder, chains.DockerNat(), rule)
			}
		}
	}
}

func (g *generator) ClearChains(builder nft.BatchBuilder, chains firewall.NFTDockerChains) {
	for _, chain := range chains.List() {
		if err := chain.Clear(builder); err != nil {
			g.logger.Error(err.Error())
		}
	}
}

func (g *generator) AddRule(builder nft.BatchBuilder, chainDocker chain.Docker, rule string) {
	if err := chainDocker.AddRule(builder, rule); err != nil {
		g.logger.Error(err.Error())
	}
}
