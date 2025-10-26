package firewall

import (
	nftableCchain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	nftablesFamily "git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

func (f *firewall) reloadInput() error {
	family := nftablesFamily.INET
	tableName := f.config.MetadataNaming.TableName
	chainName := f.config.MetadataNaming.ChainInputName

	baseChain := nftableCchain.BaseChainOptions{
		Type:     nftableCchain.TypeFilter,
		Hook:     nftableCchain.HookInput,
		Priority: 0,
		Policy:   f.config.Policy.Input.ChainDefaultPolicy(),
		Device:   "",
	}
	if err := f.nft.Chain().Add(family, tableName, chainName, baseChain); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "iifname lo counter accept"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "iifname != \"lo\" ct state related,established counter accept"); err != nil {
		return err
	}

	if f.config.Policy.Input == PolicyReject {
		if err := f.nft.Rule().Add(family, tableName, chainName, "reject"); err != nil {
			return err
		}
	}

	return nil
}
