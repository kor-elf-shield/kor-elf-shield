package reload

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/rule"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
)

func (r *reload) forward(builder nft.BatchBuilder, dockerChains firewall.NFTDockerChains) error {
	r.logger.Debug("Reloading forward chain")

	batchForward, err := chain.NewBatchForward(
		builder,
		r.table.family,
		r.table.name,
		r.config.MetadataNaming.ChainForwardName,
		r.config.Policy.DefaultAllowForward,
		r.config.Policy.ForwardPriority,
	)
	if err != nil {
		return err
	}

	localForward, err := chain.NewBatchChain(builder, r.table.family, r.table.name, "local-forward")
	if err != nil {
		return err
	}
	if err := r.forwardAddIPs(batchForward, localForward); err != nil {
		return err
	}

	if r.config.Options.DockerSupport && dockerChains != nil {
		if err := r.docker(builder, batchForward, dockerChains); err != nil {
			return err
		}
	}

	if r.config.Policy.DefaultAllowForward == false {
		drop := r.config.Policy.ForwardDrop.String()
		if err := batchForward.AddRule(drop); err != nil {
			return err
		}
	}

	return nil
}

func (r *reload) forwardAddIPs(batchForward chain.Chain, localForward chain.Chain) error {
	if err := localForward.AddRuleIn(batchForward.AddRule); err != nil {
		return err
	}

	for _, ipConfig := range r.config.IP4.InIPs {
		if ipConfig.Action != types.ActionDrop && ipConfig.Action != types.ActionReject {
			continue
		}
		if err := rule.ForwardAddIP(localForward.AddRule, ipConfig, "ip"); err != nil {
			return err
		}
	}

	if !r.config.IP6.Enable {
		return nil
	}

	for _, ipConfig := range r.config.IP6.InIPs {
		if ipConfig.Action != types.ActionDrop && ipConfig.Action != types.ActionReject {
			continue
		}
		if err := rule.ForwardAddIP(localForward.AddRule, ipConfig, "ip6"); err != nil {
			return err
		}
	}

	return nil
}
