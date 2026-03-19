package entity

import (
	"time"
)

type Blocklist struct {
	UpdatedAtUnix int64    `json:"UpdateAtUnix"`
	IPsV4         []string `json:"IPsV4"`
	IPsV6         []string `json:"IPsV6"`
}

// IsFresh returns true if the blocklist is fresh.
func (b *Blocklist) IsFresh(interval time.Duration) bool {
	lastUpdate := time.Unix(b.UpdatedAtUnix, 0)
	return b.UpdatedAtUnix > 0 && time.Since(lastUpdate) <= interval
}
