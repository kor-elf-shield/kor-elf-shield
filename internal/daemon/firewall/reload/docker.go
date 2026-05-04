package reload

import (
	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
)

func (r *reload) docker(builder nft.BatchBuilder, batchForward chain.Chain, dockerChains firewall.NFTDockerChains) error {
	r.logger.Debug("Reload docker rules")

	for _, dockerChain := range dockerChains.List() {
		if err := r.addChain(builder, dockerChain.Name(), nftChain.TypeNone); err != nil {
			return err
		}
	}

	if err := batchForward.AddRule("jump docker_forward_filter"); err != nil {
		return err
	}

	return r.dockerPrerouting(builder, dockerChains)
}

func (r *reload) dockerPrerouting(builder nft.BatchBuilder, dockerChains firewall.NFTDockerChains) error {
	var rule []string

	preroutingNat, err := r.addChainWithReturn(builder, "prerouting_nat", nftChain.BaseChainOptions{
		Type:     nftChain.TypeNat,
		Hook:     nftChain.HookPrerouting,
		Priority: -100,
		Policy:   nftChain.PolicyAccept,
		Device:   "",
	})
	if err != nil {
		return err
	}
	rule = []string{"fib daddr type local counter jump ", dockerChains.DockerNat().Name()}
	if err := preroutingNat.AddRule(rule...); err != nil {
		return err
	}

	preroutingFilter, err := r.addChainWithReturn(builder, "prerouting_filter", nftChain.BaseChainOptions{
		Type:     nftChain.TypeFilter,
		Hook:     nftChain.HookPrerouting,
		Priority: -300,
		Policy:   nftChain.PolicyAccept,
		Device:   "",
	})
	if err != nil {
		return err
	}
	rule = []string{"jump ", dockerChains.PreroutingFilter().Name()}
	if err := preroutingFilter.AddRule(rule...); err != nil {
		return err
	}

	outputNat, err := r.addChainWithReturn(builder, "output_nat", nftChain.BaseChainOptions{
		Type:     nftChain.TypeNat,
		Hook:     nftChain.HookOutput,
		Priority: -100,
		Policy:   nftChain.PolicyAccept,
		Device:   "",
	})
	if err != nil {
		return err
	}
	rule = []string{"ip daddr != 127.0.0.0/8 fib daddr type local counter jump ", dockerChains.DockerNat().Name()}
	if err := outputNat.AddRule(rule...); err != nil {
		return err
	}
	rule = []string{"ip6 daddr != ::1 fib daddr type local counter jump ", dockerChains.DockerNat().Name()}
	if err := outputNat.AddRule(rule...); err != nil {
		return err
	}

	postroutingNat, err := r.addChainWithReturn(builder, "postrouting_nat", nftChain.BaseChainOptions{
		Type:     nftChain.TypeNat,
		Hook:     nftChain.HookPostrouting,
		Priority: 300,
		Policy:   nftChain.PolicyAccept,
		Device:   "",
	})
	if err != nil {
		return err
	}
	rule = []string{"jump ", dockerChains.PostroutingNat().Name()}
	if err := postroutingNat.AddRule(rule...); err != nil {
		return err
	}

	return nil
}
