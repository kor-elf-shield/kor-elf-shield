package block

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	nftFirewall "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/rule"
)

type Blocklist interface {
	// ReplaceElements Replace the elements of the list.
	ReplaceElements(ipV4 []string, ipV6 []string, pathSaveNft string) error

	ReplaceElementsWithFile(pathNft string) error

	// AddRuleToChain Add a rule to the parent chain.
	AddRuleToChain(chainAddRuleFunc rule.AddFunc, action string) error
}

type blocklist struct {
	nft nftFirewall.NFT

	listIPv4 List
	listIPv6 List
}

func NewBlocklist(nft nftFirewall.NFT, builder nft.BatchBuilder, family family.Type, table string, name string) (Blocklist, error) {
	listNameV4, listNameV6 := getNamesIP(name)

	params := "type ipv4_addr; flags interval; auto-merge;"
	listIPv4, err := newList(nft, builder, family, table, listNameV4, params)
	if err != nil {
		return nil, err
	}

	params = "type ipv6_addr; flags interval; auto-merge;"
	listIPv6, err := newList(nft, builder, family, table, listNameV6, params)
	if err != nil {
		return nil, err
	}

	return &blocklist{
		nft: nft,

		listIPv4: listIPv4,
		listIPv6: listIPv6,
	}, nil
}

func NewBlocklistWithoutCommand(nft nftFirewall.NFT, family family.Type, table string, name string) Blocklist {
	listNameV4, listNameV6 := getNamesIP(name)

	listIPv4 := newListWithoutCommand(nft, family, table, listNameV4)
	listIPv6 := newListWithoutCommand(nft, family, table, listNameV6)

	return &blocklist{
		nft: nft,

		listIPv4: listIPv4,
		listIPv6: listIPv6,
	}
}

func (l *blocklist) ReplaceElements(ipV4 []string, ipV6 []string, pathSaveNft string) error {
	batchBuilder, err := l.nft.NewBuildBatch()
	if err != nil {
		return err
	}
	defer func() {
		_ = batchBuilder.Close()
	}()

	if err := l.listIPv4.ReplaceBatchElements(batchBuilder, ipV4); err != nil {
		return err
	}
	if err := l.listIPv6.ReplaceBatchElements(batchBuilder, ipV6); err != nil {
		return err
	}

	return l.nft.RunBatchAndMoveFile(batchBuilder, pathSaveNft)
}

func (l *blocklist) ReplaceElementsWithFile(pathNft string) error {
	args := []string{"-f", pathNft}
	return l.nft.NFT().Command().Run(args...)
}

func (l *blocklist) AddRuleToChain(chainAddRuleFunc rule.AddFunc, action string) error {
	addRule := "ip saddr @" + l.listIPv4.Name() + " " + action
	if err := chainAddRuleFunc(addRule); err != nil {
		return err
	}

	addRule = "ip6 saddr @" + l.listIPv6.Name() + " " + action
	if err := chainAddRuleFunc(addRule); err != nil {
		return err
	}

	return nil
}
