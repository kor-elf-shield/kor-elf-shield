package analyzer

import (
	"fmt"
	"regexp"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
)

var (
	reName = regexp.MustCompile(`^[A-Za-z0-9-_]{0,255}$`)
)

type LogAlertRule struct {
	Enabled  bool   `mapstructure:"enabled"`
	Notify   bool   `mapstructure:"notify"`
	Name     string `mapstructure:"name"`
	Message  string `mapstructure:"message"`
	Source   Source
	Patterns []LogAlertPattern
}

func (l *LogAlertRule) ToSource(isNotify bool) (*config.Source, error) {
	if err := l.validate(); err != nil {
		return nil, err
	}

	source, err := l.Source.ToSource()
	if err != nil {
		return nil, err
	}

	var patterns []config.AlertRegexPattern

	for _, pattern := range l.Patterns {
		p, err := pattern.ToPattern()
		if err != nil {
			return nil, err
		}
		patterns = append(patterns, p)
	}

	if len(patterns) == 0 {
		return nil, fmt.Errorf("patterns is empty")
	}

	source.AlertRule = &config.AlertRule{
		Name:           l.Name,
		Message:        l.Message,
		IsNotification: isNotify && l.Notify,
		Patterns:       patterns,
	}

	return source, nil
}

func (l *LogAlertRule) validate() error {
	if l.Name == "" {
		return fmt.Errorf("name is empty")
	}

	if !reName.MatchString(l.Name) {
		return fmt.Errorf("invalid name")
	}

	return nil
}
