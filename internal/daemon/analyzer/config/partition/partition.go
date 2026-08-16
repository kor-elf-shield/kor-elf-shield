package partition

import "strings"

type PatternPartitionType interface {
	Accepts(text string) (bool, string)
	IsPartitioned() bool
}

type PatternPartition struct {
	Value     uint8
	Type      PatternPartitionType
	Normalize func(keyword string) string
}

type Normalize interface {
	Normalize(text string) string
}

type normalize struct {
	Trim      bool
	LowerCase bool
}

func NewNormalize(trim, lowerCase bool) Normalize {
	return &normalize{
		Trim:      trim,
		LowerCase: lowerCase,
	}
}

func (p *normalize) Normalize(text string) string {
	if p.LowerCase {
		text = strings.ToLower(text)
	}
	if p.Trim {
		text = strings.TrimSpace(text)
	}
	return text
}
