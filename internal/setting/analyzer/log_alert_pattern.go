package analyzer

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/regular_expression"
)

type LogAlertPattern struct {
	Regexp string `mapstructure:"regexp"`
	Values []PatternValue
}

func (p *LogAlertPattern) ToPattern() (config.AlertRegexPattern, error) {
	pattern := config.AlertRegexPattern{
		Regexp: regular_expression.NewLazyRegexp(p.Regexp),
	}

	for _, value := range p.Values {
		v, err := value.ToPatternValue()
		if err != nil {
			return config.AlertRegexPattern{}, err
		}

		pattern.Values = append(pattern.Values, v)
	}

	return pattern, nil
}
