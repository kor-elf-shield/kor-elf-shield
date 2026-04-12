package config

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/regular_expression"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

var (
	reSystemdUnitValue = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._@-]{0,255}\.(service|socket|target|mount|timer|path|scope|slice|device)$`)
	reSyslogIDValue    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._@-]{0,127}$`)
)

type SourceType string

const (
	SourceTypeJournal SourceType = "journalctl"
	SourceTypeFile    SourceType = "file"
)

type JournalField string

const (
	JournalFieldSystemdUnit      JournalField = "_SYSTEMD_UNIT"
	JournalFieldSyslogIdentifier JournalField = "SYSLOG_IDENTIFIER"
)

type Config struct {
	BinPath BinPath
	Sources []*Source
}

type SourceJournal struct {
	Field JournalField
	Match string
}

func NewSourceJournal(field JournalField, match string) (*SourceJournal, error) {
	v := strings.TrimSpace(match)
	if v == "" {
		return nil, fmt.Errorf("journal match must not be empty")
	}
	if len(v) > 512 {
		return nil, fmt.Errorf("journal match is too long: %d", len(v))
	}
	for _, r := range v {
		if r == 0 || r == '\n' || r == '\r' || unicode.IsControl(r) {
			return nil, fmt.Errorf("journal match contains control characters")
		}
	}
	// to avoid breaking the FIELD=VALUE format and concatenation with '+'
	if strings.ContainsAny(v, "=+") {
		return nil, fmt.Errorf("journal match must not contain '=' or '+'")
	}

	if strings.ContainsAny(v, " \t") {
		return nil, fmt.Errorf("journal match must not contain spaces or tabs")
	}

	switch field {
	case JournalFieldSystemdUnit:
		if !reSystemdUnitValue.MatchString(v) {
			return nil, fmt.Errorf("invalid _SYSTEMD_UNIT value: %q", v)
		}
	case JournalFieldSyslogIdentifier:
		if !reSyslogIDValue.MatchString(v) {
			return nil, fmt.Errorf("invalid SYSLOG_IDENTIFIER value: %q", v)
		}
	default:
		return nil, fmt.Errorf("invalid journal field: %q", field)
	}

	return &SourceJournal{Field: field, Match: v}, nil
}

type SourceFile struct {
	Path string
}

func NewSourceFile(path string) (*SourceFile, error) {
	if err := validate.PathFile(path, "logAlert.rules.source.path"); err != nil {
		return nil, err
	}

	return &SourceFile{Path: path}, nil
}

func (s *SourceJournal) JournalctlMatch() string {
	return string(s.Field) + "=" + s.Match
}

type Source struct {
	Type    SourceType
	Journal *SourceJournal
	File    *SourceFile

	AlertRule                *AlertRule
	BruteForceProtectionRule *brute_force_protection.Rule
}

type AlertRule struct {
	Name           string
	Message        string
	IsNotification bool
	Patterns       []AlertRegexPattern
	Group          *AlertGroup
}

type AlertRegexPattern struct {
	Regexp *regular_expression.LazyRegexp
	Values []PatternValue
}

type PatternValue struct {
	Name  string
	Value uint8
	Type  PatternTypeValue
}

type PatternTypeValue string

const (
	PatternValueIP PatternTypeValue = "ip"
)
