package block

import (
	"fmt"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type List interface {
	Name() string
	AddElement(element string) error
}

type list struct {
	nft    nft.NFT
	family family.Type
	table  string
	name   string
}

func newList(nft nft.NFT, family family.Type, table string, name string, params string) (List, error) {
	command := []string{
		"add set", family.String(), table, name, "{ " + params + " }",
	}
	if err := nft.Command().Run(command...); err != nil {
		return nil, err
	}

	return &list{
		nft:    nft,
		family: family,
		table:  table,
		name:   name,
	}, nil
}

func (l *list) Name() string {
	return l.name
}

func (l *list) AddElement(element string) error {
	command := []string{
		"add element inet",
		l.table, l.name,
		fmt.Sprintf("{ %s }", element),
	}
	return l.nft.Command().Run(command...)
}
