package analyzer

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/partition"
)

type PatternPartitionType struct {
	Type        string   `mapstructure:"type"`
	Keywords    []string `mapstructure:"keywords"`
	Partitioned bool     `mapstructure:"partitioned"`
}

func (p *PatternPartitionType) ToType(normalize partition.Normalize) (partition.PatternPartitionType, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}

	var keywords []string
	for _, keyword := range p.Keywords {
		keywords = append(keywords, normalize.Normalize(keyword))
	}

	if p.Type == "except" {
		return partition.NewExceptType(keywords, p.Partitioned), nil
	}

	return partition.NewOnlyType(keywords, p.Partitioned), nil
}

func (p *PatternPartitionType) validate() error {
	if p.Type != "only" && p.Type != "except" {
		return errors.New("invalid partition type. only 'only' and 'except' are supported")
	}

	if len(p.Keywords) == 0 {
		return errors.New("invalid partition type. keywords are required")
	}

	return nil
}
