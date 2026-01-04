package chain

import nftChain "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain"

func NewChains(newNoneChain func(chain string) (nftChain.Chain, error)) (Chains, error) {
	chainsData := &chains{}

	if data, err := newChainData("docker_nat", newNoneChain); err != nil {
		return nil, err
	} else {
		chainsData.DockerNat = data
	}

	if data, err := newChainData("docker_postrouting_nat", newNoneChain); err != nil {
		return nil, err
	} else {
		chainsData.PostroutingNat = data
	}

	if data, err := newChainData("docker_prerouting_filter", newNoneChain); err != nil {
		return nil, err
	} else {
		chainsData.PreroutingFilter = data
	}

	if data, err := newChainData("docker_filter", newNoneChain); err != nil {
		return nil, err
	} else {
		chainsData.DockerFilter = data
	}

	if data, err := newChainData("docker_filter_first", newNoneChain); err != nil {
		return nil, err
	} else {
		chainsData.DockerFilterFirst = data
	}

	if data, err := newChainData("docker_filter_second", newNoneChain); err != nil {
		return nil, err
	} else {
		chainsData.DockerFilterSecond = data
	}

	if data, err := newChainData("docker_forward_filter", newNoneChain); err != nil {
		return nil, err
	} else {
		chainsData.ForwardFilter = data
	}

	if data, err := newChainData("docker_forward_bridge", newNoneChain); err != nil {
		return nil, err
	} else {
		chainsData.ForwardBridge = data
	}

	if data, err := newChainData("docker_forward_ct", newNoneChain); err != nil {
		return nil, err
	} else {
		chainsData.ForwardCT = data
	}

	return chainsData, nil
}

func newChainData(chainName string, newNoneChain func(chain string) (nftChain.Chain, error)) (Data, error) {
	data := Data{
		name: chainName,
	}

	newChain, err := newNoneChain(data.name)
	if err != nil {
		return data, err
	}

	data.chain = newChain
	return data, nil
}
