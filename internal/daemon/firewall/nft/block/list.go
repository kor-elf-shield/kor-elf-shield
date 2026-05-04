package block

import (
	"fmt"
	"strings"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	nftFirewall "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft"
)

type List interface {
	Name() string
	AddElement(element string) error
	AddBatchElement(builder nft.BatchBuilder, element string) error
	DeleteElement(element string) error
	ReplaceElements(elements []string) error
	ReplaceBatchElements(builder nft.BatchBuilder, elements []string) error
}

type list struct {
	nft    nftFirewall.NFT
	family family.Type
	table  string
	name   string
}

func newList(nft nftFirewall.NFT, builder nft.BatchBuilder, family family.Type, table string, name string, params string) (List, error) {
	command := []string{
		"add set", family.String(), table, name, "{ " + params + " }",
	}
	if err := builder.Command().Run(command...); err != nil {
		return nil, err
	}

	return &list{
		nft:    nft,
		family: family,
		table:  table,
		name:   name,
	}, nil
}

func newListWithoutCommand(nft nftFirewall.NFT, family family.Type, table string, name string) List {
	return &list{
		nft:    nft,
		family: family,
		table:  table,
		name:   name,
	}
}

func (l *list) Name() string {
	return l.name
}

func (l *list) AddElement(element string) error {
	command := []string{
		"add element",
		l.family.String(), l.table, l.name,
		fmt.Sprintf("{ %s }", element),
	}
	return l.nft.NFT().Command().Run(command...)
}

func (l *list) AddBatchElement(builder nft.BatchBuilder, element string) error {
	command := []string{
		"add element",
		l.family.String(), l.table, l.name,
		fmt.Sprintf("{ %s }", element),
	}
	return builder.Command().Run(command...)
}

func (l *list) DeleteElement(element string) error {
	command := []string{
		"delete element",
		l.family.String(), l.table, l.name,
		fmt.Sprintf("{ %s }", element),
	}
	return l.nft.NFT().Command().Run(command...)
}

func (l *list) ReplaceElements(elements []string) error {
	batchBuilder, err := l.nft.NewBuildBatch()
	if err != nil {
		return err
	}
	defer func() {
		_ = batchBuilder.Close()
	}()

	if err := l.replaceElements(batchBuilder, elements); err != nil {
		return err
	}

	return l.nft.RunBatch(batchBuilder)
}

func (l *list) ReplaceBatchElements(builder nft.BatchBuilder, elements []string) error {
	return l.replaceElements(builder, elements)
}

func (l *list) replaceElements(builder nft.BatchBuilder, elements []string) error {
	if err := builder.Command().Run("flush set", l.family.String(), l.table, l.name); err != nil {
		return err
	}

	if len(elements) == 0 {
		return nil
	}

	command := []string{
		"add element",
		l.family.String(), l.table, l.name,
		fmt.Sprintf("{ %s }", strings.Join(elements, ",")),
	}

	return builder.Command().Run(command...)
}
