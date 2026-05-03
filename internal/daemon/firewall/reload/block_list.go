package reload

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/block"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
	nftTable "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/table"
)

func (r *reload) blockList(builder nft.BatchBuilder, beforeLocalInput chain.Chain) (nftTable.BlockList, error) {
	listBlockedIP, err := block.NewListIP(r.nft, builder, r.table.family, r.table.name, "blocked_ip")
	if err != nil {
		return nil, err
	}

	if err := listBlockedIP.AddRuleToChain(beforeLocalInput.AddRule, "drop"); err != nil {
		return nil, err
	}

	listBlockedIPWithPort, err := block.NewListIPWithPort(r.nft, builder, r.table.family, r.table.name, "blocked_ip_port")
	if err != nil {
		return nil, err
	}

	if err := listBlockedIPWithPort.AddRuleToChain(beforeLocalInput.AddRule, "drop"); err != nil {
		return nil, err
	}

	return nftTable.NewBlockList(listBlockedIP, listBlockedIPWithPort), nil
}
