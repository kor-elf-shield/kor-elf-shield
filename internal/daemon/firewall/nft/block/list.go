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
	ReplaceElementsWithSaveNFTFile(elements []string, pathFile string) error
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
	return l.replaceElements(elements, func(builder nft.BatchBuilder) error {
		return l.nft.RunBatch(builder)
	})
}

func (l *list) ReplaceElementsWithSaveNFTFile(elements []string, pathFile string) error {
	return l.replaceElements(elements, func(builder nft.BatchBuilder) error {
		return l.nft.RunBatchAndMoveFile(builder, pathFile)
	})
}

func (l *list) replaceElements(elements []string, run func(builder nft.BatchBuilder) error) error {
	batchBuilder, err := l.nft.NewBuildBatch()
	if err != nil {
		return err
	}
	defer func() {
		_ = batchBuilder.Close()
	}()

	if err := batchBuilder.Command().Run("flush set", l.family.String(), l.table, l.name); err != nil {
		return err
	}

	if len(elements) == 0 {
		return run(batchBuilder)
	}

	command := []string{
		"add element",
		l.family.String(), l.table, l.name,
		fmt.Sprintf("{ %s }", strings.Join(elements, ",")),
	}
	if err := batchBuilder.Command().Run(command...); err != nil {
		return err
	}

	return run(batchBuilder)
}
