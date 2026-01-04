package firewall

import nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"

func (f *firewall) reloadDocker() error {
	f.logger.Debug("Reload docker rules")
	if err := f.reloadDockerPrerouting(); err != nil {
		return err
	}

	return nil
}

func (f *firewall) reloadDockerPrerouting() error {
	preroutingNat, err := f.chains.NewChain("prerouting_nat", nftChain.BaseChainOptions{
		Type:     nftChain.TypeNat,
		Hook:     nftChain.HookPrerouting,
		Priority: -100,
		Policy:   nftChain.PolicyAccept,
		Device:   "",
	})
	if err != nil {
		return err
	}
	if err := f.docker.NftChains().PreroutingNatJump(preroutingNat.AddRule); err != nil {
		return err
	}

	preroutingFilter, err := f.chains.NewChain("prerouting_filter", nftChain.BaseChainOptions{
		Type:     nftChain.TypeFilter,
		Hook:     nftChain.HookPrerouting,
		Priority: -300,
		Policy:   nftChain.PolicyAccept,
		Device:   "",
	})
	if err != nil {
		return err
	}
	if err := f.docker.NftChains().PreroutingFilterJump(preroutingFilter.AddRule); err != nil {
		return err
	}

	outputNat, err := f.chains.NewChain("output_nat", nftChain.BaseChainOptions{
		Type:     nftChain.TypeNat,
		Hook:     nftChain.HookOutput,
		Priority: -100,
		Policy:   nftChain.PolicyAccept,
		Device:   "",
	})
	if err != nil {
		return err
	}
	if err := f.docker.NftChains().OutputNatJump(outputNat.AddRule); err != nil {
		return err
	}

	postroutingNat, err := f.chains.NewChain("postrouting_nat", nftChain.BaseChainOptions{
		Type:     nftChain.TypeNat,
		Hook:     nftChain.HookPostrouting,
		Priority: 300,
		Policy:   nftChain.PolicyAccept,
		Device:   "",
	})
	if err != nil {
		return err
	}
	if err := f.docker.NftChains().PostroutingNatJump(postroutingNat.AddRule); err != nil {
		return err
	}

	return nil
}
