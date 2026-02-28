package analyzer

import (
	"errors"
	"fmt"
	"strings"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
)

type Source struct {
	Type  string `mapstructure:"type"`
	Field string `mapstructure:"field"`
	Match string `mapstructure:"match"`
	Path  string `mapstructure:"path"`
}

func (s *Source) ToSource() (*config.Source, error) {
	switch s.Type {
	case "journalctl":
		field, err := s.journalField()
		if err != nil {
			return nil, err
		}

		journal, err := config.NewSourceJournal(field, s.Match)
		if err != nil {
			return nil, err
		}

		return &config.Source{
			Type:    config.SourceTypeJournal,
			Journal: journal,
		}, nil
	case "file":
		file, err := config.NewSourceFile(s.Path)
		if err != nil {
			return nil, err
		}

		return &config.Source{
			Type: config.SourceTypeFile,
			File: file,
		}, nil
	}

	return nil, errors.New(fmt.Sprintf("unknown source type: %s. journalctl or file are allowed.", s.Type))
}

func (s *Source) journalField() (config.JournalField, error) {
	switch strings.ToLower(s.Field) {
	case "systemd_unit":
		return config.JournalFieldSystemdUnit, nil
	case "syslog_identifier":
		return config.JournalFieldSyslogIdentifier, nil
	}

	return "", errors.New(fmt.Sprintf("unknown journal field: %s. systemd_unit or syslog_identifier are allowed.", s.Field))
}
