package blocklist

import "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/block"

type FalseBlocklist struct {
}

func NewFalseBlocklist() Blocklist {
	return &FalseBlocklist{}
}

func (b *FalseBlocklist) Names() []string {
	return []string{}
}

func (b *FalseBlocklist) NftReload(_ map[string]block.Blocklist) error {
	return nil
}

func (b *FalseBlocklist) Run() {}

func (b *FalseBlocklist) Close() error {
	return nil
}
