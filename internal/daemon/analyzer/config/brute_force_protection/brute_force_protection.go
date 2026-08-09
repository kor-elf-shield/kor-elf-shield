package brute_force_protection

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/partition"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/regular_expression"
)

type Rule struct {
	Name    string
	Message string

	IsNotification       bool
	NotificationCooldown uint32
	NotificationEvery    uint32

	Patterns []RegexPattern
	Group    *Group
}

type RegexPattern struct {
	Regexp    *regular_expression.LazyRegexp
	Values    []PatternValue
	IP        uint8
	Partition *partition.PatternPartition
}

type RateLimit struct {
	Count               uint32
	Period              uint32
	BlockingTimeSeconds uint32
	BlockConfig         Block
}

type PatternValue struct {
	Name  string
	Value uint8
}
