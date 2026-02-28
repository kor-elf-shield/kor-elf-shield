package analyzer

import (
	"fmt"
	"regexp"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
)

var (
	reName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]{0,254}$`)
)

type LogAlert struct {
	Enabled bool `mapstructure:"enabled"`
	Notify  bool `mapstructure:"notify"`
	Groups  []LogAlertGroup
	Rules   []LogAlertRule
}

func defaultLogAlert() LogAlert {
	return LogAlert{
		Enabled: true,
		Notify:  true,
		Groups:  []LogAlertGroup{},
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

	groups, err := l.groups()
	if err != nil {
		return nil, fmt.Errorf("groups: %w", err)
	}

	for _, rule := range l.Rules {
		if !rule.Enabled {
			continue
		}

		var group *config.AlertGroup
		if rule.Group != "" {
			if _, ok := groups[rule.Group]; !ok {
				return nil, fmt.Errorf("group %q not found", rule.Group)
			}
			group = groups[rule.Group]
		}

		source, err := rule.ToSource(l.Notify, group)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (l *LogAlert) groups() (map[string]*config.AlertGroup, error) {
	groups := make(map[string]*config.AlertGroup)
	for _, group := range l.Groups {
		g, err := group.ToGroup()
		if err != nil {
			return nil, err
		}
		groups[g.Name] = g
	}

	return groups, nil
}
