package config

import (
	"regexp"
	"sync"
	"time"
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
