package blocklist

import (
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/blocklist/sources"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
)

type Config struct {
	BlocklistRepository repository.BlocklistRepository
	Sources             []*SourceConfig
	PathDir             string
}

type SourceConfig struct {
	Name     string
	Interval time.Duration
	Source   sources.BlocklistSource
}
