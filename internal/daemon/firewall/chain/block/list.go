package block

import (
	"fmt"
	"strings"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type List interface {
	Name() string
	AddElement(element string) error
	DeleteElement(element string) error
	ReplaceElements(elements []string) error
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
		"add element",
		l.family.String(), l.table, l.name,
		fmt.Sprintf("{ %s }", element),
	}
	return l.nft.Command().Run(command...)
}

func (l *list) DeleteElement(element string) error {
	command := []string{
		"delete element",
		l.family.String(), l.table, l.name,
		fmt.Sprintf("{ %s }", element),
	}
	return l.nft.Command().Run(command...)
}

func (l *list) ReplaceElements(elements []string) error {
	if len(elements) == 0 {
		return nil
	}

	if err := l.nft.Command().Run("flush set", l.family.String(), l.table, l.name); err != nil {
		return err
	}

	const batchSize = 200
	for _, batch := range chunkStrings(elements, batchSize) {
		command := []string{
			"add element",
			l.family.String(), l.table, l.name,
			fmt.Sprintf("{ %s }", strings.Join(batch, ",")),
		}
		if err := l.nft.Command().Run(command...); err != nil {
			return err
		}
	}

	return nil
}

func chunkStrings(items []string, size int) [][]string {
	if size <= 0 {
		size = 100
	}

	chunks := make([][]string, 0, (len(items)+size-1)/size)
	for start := 0; start < len(items); start += size {
		end := start + size
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[start:end])
	}
	return chunks
}
