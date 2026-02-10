package analyzer

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
)

type LogAlert struct {
	Enabled bool `mapstructure:"enabled"`
	Notify  bool `mapstructure:"notify"`
	Rules   []LogAlertRule
}

func defaultLogAlert() LogAlert {
	return LogAlert{
		Enabled: true,
		Notify:  true,
		Rules:   []LogAlertRule{},
	}
}

func (l *LogAlert) Validate() error {
	return nil
}

func (l *LogAlert) ToSources() ([]*config.Source, error) {
	var sources []*config.Source

	if !l.Enabled || len(l.Rules) == 0 {
		return sources, nil
	}

	for _, rule := range l.Rules {
		if !rule.Enabled {
			continue
		}

		source, err := rule.ToSource(l.Notify)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return sources, nil
}
