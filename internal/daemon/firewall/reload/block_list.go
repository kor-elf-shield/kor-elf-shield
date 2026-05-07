package reload

import (
	"fmt"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/block"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
	nftTable "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/table"
)

func (r *reload) blockList(builder nft.BatchBuilder, beforeLocalInput chain.Chain, blocks map[string]block.Blocklist) (nftTable.BlockList, error) {
	listBlockedIP, err := block.NewListIP(r.nft, builder, r.table.family, r.table.name, blockedIP)
	if err != nil {
		return nil, err
	}

	if err := listBlockedIP.AddRuleToChain(beforeLocalInput.AddRule, "drop"); err != nil {
		return nil, err
	}

	listBlockedIPWithPort, err := block.NewListIPWithPort(r.nft, builder, r.table.family, r.table.name, blockedIPWithPort)
	if err != nil {
		return nil, err
	}

	if err := listBlockedIPWithPort.AddRuleToChain(beforeLocalInput.AddRule, "drop"); err != nil {
		return nil, err
	}

	return nftTable.NewBlockList(listBlockedIP, listBlockedIPWithPort, blocks), nil
}

func (r *reload) moduleBlockList(builder nft.BatchBuilder, afterLocalInput chain.Chain, blockListNames []string) (blocks map[string]block.Blocklist, err error) {
	r.logger.Debug("Reload blocklist")
	blocks = make(map[string]block.Blocklist)
	for _, blockListName := range blockListNames {
		if blockListName == "" {
			continue
		}
		r.logger.Debug(fmt.Sprintf("Reload blocklist from %s", blockListName))
		blockList, err := block.NewBlocklist(r.nft, builder, r.table.family, r.table.name, getBlocklistName(blockListName))
		if err != nil {
			r.logger.Error(fmt.Sprintf("Failed to create blocklist: %s", err))
			continue
		}
		if err := blockList.AddRuleToChain(afterLocalInput.AddRule, "drop"); err != nil {
			r.logger.Error(fmt.Sprintf("Failed to add rule to chain: %s", err))
			continue
		}
		blocks[blockListName] = blockList
	}

	return blocks, nil
}
