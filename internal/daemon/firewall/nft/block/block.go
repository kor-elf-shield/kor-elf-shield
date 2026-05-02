package block

import (
	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type Sets interface {
	Add(name string, params string) error
}

type setBatch struct {
	builder nft.BatchBuilder
	family  family.Type
	table   string
}

func NewBatchSet(builder nft.BatchBuilder, family family.Type, table string) Sets {
	return &setBatch{
		builder: builder,
		family:  family,
		table:   table,
	}
}

func (b *setBatch) Add(name string, params string) error {
	command := []string{
		"add set", b.family.String(), b.table, name, "{ " + params + " }",
	}
	return b.builder.Command().Run(command...)
}
