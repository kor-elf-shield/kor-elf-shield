package config

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"
)

var (
	reSystemdUnitValue = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._@-]{0,255}\.(service|socket|target|mount|timer|path|scope|slice|device)$`)
	reSyslogIDValue    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._@-]{0,127}$`)
)

type SourceType string

const (
	SourceTypeJournal SourceType = "journalctl"
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

type SourceJournal struct {
	Field JournalField
	Match string
}

func (s *SourceJournal) JournalctlMatch() string {
	return string(s.Field) + "=" + s.Match
}

type Source struct {
	Type SourceType

	Journal   *SourceJournal
	AlertRule *AlertRule
}

type AlertRule struct {
	Name           string
	Message        string
	IsNotification bool
	Patterns       []AlertRegexPattern
	Group          *AlertGroup
}

type AlertRegexPattern struct {
	Regexp *LazyRegexp
	Values []PatternValue
}

type LazyRegexp struct {
	pattern string

	once sync.Once
	re   *regexp.Regexp
	err  error
}

func NewLazyRegexp(pattern string) *LazyRegexp {
	return &LazyRegexp{pattern: pattern}
}

func (lr *LazyRegexp) Get() (*regexp.Regexp, error) {
	lr.once.Do(func() {
		lr.re, lr.err = regexp.Compile(lr.pattern)
	})
	return lr.re, lr.err
}

type PatternValue struct {
	Name  string
	Value uint8
}

type RateLimit struct {
	Count  uint32
	Period time.Duration
}

type AlertGroup struct {
	Name                 string
	Message              string
	RateLimits           []RateLimit
	RateLimitResetPeriod time.Duration
}
