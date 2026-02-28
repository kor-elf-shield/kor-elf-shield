package analyzer

import (
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/regular_expression"
)

type BruteForceProtectionPattern struct {
	Regexp string `mapstructure:"regexp"`
	IP     int    `mapstructure:"ip"`
	Values []PatternValue
}

func (p *BruteForceProtectionPattern) ToPattern() (brute_force_protection.RegexPattern, error) {
	if err := p.validate(); err != nil {
		return brute_force_protection.RegexPattern{}, err
	}

	pattern := brute_force_protection.RegexPattern{
		Regexp: regular_expression.NewLazyRegexp(p.Regexp),
		IP:     uint8(p.IP),
	}

	for _, value := range p.Values {
		v, err := value.ToPatternValueForBruteForceProtection()
		if err != nil {
			return brute_force_protection.RegexPattern{}, err
		}

		pattern.Values = append(pattern.Values, v)
	}

	return pattern, nil
}

func (p *BruteForceProtectionPattern) validate() error {
	if p.IP <= 0 || p.IP > 255 {
		return errors.New("ip must be between 1 and 255")
	}

	if p.Regexp == "" {
		return errors.New("regexp is empty")
	}

	return nil
}
