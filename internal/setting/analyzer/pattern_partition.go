package analyzer

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/partition"
)

type PatternPartition struct {
	Value     int  `mapstructure:"value"`
	Trim      bool `mapstructure:"trim"`
	LowerCase bool `mapstructure:"lower_case"`

	Type *PatternPartitionType `mapstructure:"type"`
}

func (p *PatternPartition) ToPatternPartition() (*partition.PatternPartition, error) {
	if err := p.validate(); err != nil {
		return nil, err
	}

	normalize := partition.NewNormalize(p.Trim, p.LowerCase)

	var patternPartitionType partition.PatternPartitionType
	if p.Type != nil {
		if partitionType, err := p.Type.ToType(normalize); err != nil {
			return nil, err
		} else {
			patternPartitionType = partitionType
		}
	}

	return &partition.PatternPartition{
		Value:     uint8(p.Value),
		Type:      patternPartitionType,
		Normalize: normalize.Normalize,
	}, nil
}

func (p *PatternPartition) validate() error {
	if p.Value <= 0 || p.Value > 255 {
		return errors.New("invalid partition value. min: 1, max: 255")
	}

	return nil
}
