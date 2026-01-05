package docker_monitor

import (
	"context"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/client"
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
	client client.Docker
	logger log.Logger
	ctx    context.Context

	chains chain.Chains
}

func New(path string, ctx context.Context, logger log.Logger) Docker {
	return &docker{
		client: client.NewDocker(path, ctx, logger),
		logger: logger,
		ctx:    ctx,
	}
}

func (d *docker) NftReload(newNoneChain func(chain string) (nftChain.Chain, error)) error {
	chains, err := chain.NewChains(newNoneChain)
	if err != nil {
		return err
	}
	d.chains = chains

	d.nftRuleReload()

	return nil
}

func (d *docker) NftChains() chain.Chains {
	return d.chains
}

func (d *docker) Run() {
	events := d.client.Events()
	for {
		select {
		case <-d.ctx.Done():
			return
		case msg := <-events:
			if msg == "" {
				continue
			}
			d.logger.Debug("Docker event received: " + msg)
			// TODO: A temporary solution to test how it will interact with nftables in a production environment
			listChains := d.NftChains().List()

			if err := listChains.DockerNat.Clear(); err != nil {
				d.logger.Error(err.Error())
			}
			if err := listChains.PostroutingNat.Clear(); err != nil {
				d.logger.Error(err.Error())
			}
			if err := listChains.PreroutingFilter.Clear(); err != nil {
				d.logger.Error(err.Error())
			}
			if err := listChains.DockerFilter.Clear(); err != nil {
				d.logger.Error(err.Error())
			}
			if err := listChains.DockerFilterFirst.Clear(); err != nil {
				d.logger.Error(err.Error())
			}
			if err := listChains.DockerFilterSecond.Clear(); err != nil {
				d.logger.Error(err.Error())
			}
			if err := listChains.ForwardFilter.Clear(); err != nil {
				d.logger.Error(err.Error())
			}

			if err := listChains.ForwardBridge.Clear(); err != nil {
				d.logger.Error(err.Error())
			}

			if err := listChains.ForwardCT.Clear(); err != nil {
				d.logger.Error(err.Error())
			}

			d.nftRuleReload()
		}
	}
}

func (d *docker) Close() error {
	return d.client.EventsClose()
}

func (d *docker) chainCommand(chainData chain.Data, rule string) {
	if err := chainData.AddRule(rule); err != nil {
		d.logger.Error(err.Error())
	}
}

func (d *docker) nftRuleReload() {
	listChains := d.NftChains().List()

	if err := listChains.ForwardCT.JumpTo(&listChains.ForwardFilter, ""); err != nil {
		d.logger.Error(err.Error())
	}
	if err := listChains.ForwardBridge.JumpTo(&listChains.ForwardFilter, ""); err != nil {
		d.logger.Error(err.Error())
	}
	if err := listChains.DockerFilterFirst.JumpTo(&listChains.DockerFilter, ""); err != nil {
		d.logger.Error(err.Error())
	}
	if err := listChains.DockerFilterSecond.JumpTo(&listChains.DockerFilter, ""); err != nil {
		d.logger.Error(err.Error())
	}

	bridges, err := d.client.FetchBridges()
	if err != nil {
		d.logger.Error(err.Error())
		return
	}
	var rule string
	for _, bridge := range bridges {
		rule = fmt.Sprintf("iifname != \"%s\" oifname \"%s\" counter drop", bridge.Name, bridge.Name)
		d.chainCommand(listChains.DockerFilterSecond, rule)

		rule = fmt.Sprintf("iifname \"%s\" counter accept", bridge.Name)
		d.chainCommand(listChains.ForwardFilter, rule)

		rule = fmt.Sprintf("oifname \"%s\" counter", bridge.Name)
		if err := listChains.DockerFilter.JumpTo(&listChains.ForwardBridge, rule); err != nil {
			d.logger.Error(err.Error())
		}

		rule = fmt.Sprintf("oifname \"%s\" ct state related,established counter accept", bridge.Name)
		d.chainCommand(listChains.ForwardCT, rule)

		rule = fmt.Sprintf("ip saddr %s oifname != \"%s\" counter masquerade", bridge.Subnet, bridge.Name)
		d.chainCommand(listChains.PostroutingNat, rule)

		if bridge.Containers == nil {
			continue
		}
		for _, container := range bridge.Containers {
			for _, ipInfo := range container.Networks.IPAddresses {
				rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" counter drop", ipInfo.NftPrefix(), ipInfo.Address, bridge.Name)
				d.chainCommand(listChains.PreroutingFilter, rule)

				for _, port := range container.Networks.Ports {
					isZeroAddress := false
					for _, hostInfo := range port.HostPort {
						if hostInfo.IP.Address != "0.0.0.0" && hostInfo.IP.Address != "::" && (hostInfo.IP.Address == "127.0.0.1" || hostInfo.IP.Address == "::1") {
							rule = fmt.Sprintf("%s daddr %s iifname != \"lo\" %s dport %s counter drop", hostInfo.IP.NftPrefix(), hostInfo.IP.Address, port.Protocol, hostInfo.Port)
							d.chainCommand(listChains.PreroutingFilter, rule)
						}

						if hostInfo.IP.Address == "0.0.0.0" || hostInfo.IP.Address == "::" {
							if isZeroAddress {
								continue
							}
							isZeroAddress = true
							rule = fmt.Sprintf("iifname != \"%s\" %s dport %s counter dnat %s to %s:%s", bridge.Name, port.Protocol, hostInfo.Port, ipInfo.NftPrefix(), ipInfo.Address, port.Port)
							d.chainCommand(listChains.DockerNat, rule)

							rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" oifname \"%s\" %s dport %s counter accept", ipInfo.NftPrefix(), ipInfo.Address, bridge.Name, bridge.Name, port.Protocol, port.Port)
							d.chainCommand(listChains.DockerFilterFirst, rule)
							continue
						}
						rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" oifname \"%s\" %s dport %s counter accept", ipInfo.NftPrefix(), ipInfo.Address, bridge.Name, bridge.Name, port.Protocol, port.Port)
						d.chainCommand(listChains.DockerFilterFirst, rule)

						rule = fmt.Sprintf("%s daddr %s iifname != \"%s\" %s dport %s counter dnat to %s:%s", hostInfo.IP.NftPrefix(), hostInfo.IP.Address, bridge.Name, port.Protocol, hostInfo.Port, ipInfo.Address, port.Port)
						d.chainCommand(listChains.DockerNat, rule)
					}
				}
			}
		}
	}
}
