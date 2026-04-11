package geoip

import (
	"time"

	"git.kor-elf.net/kor-elf-shield/geoip2"
)

type Config struct {
	GeoIP    geoip2.RefreshableGeoIP2
	Interval time.Duration
}
