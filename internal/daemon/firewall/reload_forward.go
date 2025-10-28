package firewall

import (
	nftableCchain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	nftablesFamily "git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

func (f *firewall) reloadForward() error {
	family := nftablesFamily.INET
	tableName := f.config.MetadataNaming.TableName
	chainName := f.config.MetadataNaming.ChainForwardName

	baseChain := nftableCchain.BaseChainOptions{
		Type:     nftableCchain.TypeFilter,
		Hook:     nftableCchain.HookForward,
		Priority: 0,
		Policy:   f.config.Policy.Forward.ChainDefaultPolicy(),
		Device:   "",
	}
	if err := f.nft.Chain().Add(family, tableName, chainName, baseChain); err != nil {
		return err
	}

	if f.config.Policy.Forward == PolicyReject {
		if err := f.nft.Rule().Add(family, tableName, chainName, f.config.Policy.Forward.String()); err != nil {
			return err
		}
	}

	return nil
}
