package entity

import (
	"time"
)

type Blocklist struct {
	UpdatedAtUnix int64  `json:"UpdateAtUnix"`
	Checksum      string `json:"checksum"`
}

// IsFresh returns true if the blocklist is fresh.
func (b *Blocklist) IsFresh(interval time.Duration) bool {
	if b.Checksum == "" {
		return false
	}

	lastUpdate := time.Unix(b.UpdatedAtUnix, 0)
	return b.UpdatedAtUnix > 0 && time.Since(lastUpdate) <= interval
}
