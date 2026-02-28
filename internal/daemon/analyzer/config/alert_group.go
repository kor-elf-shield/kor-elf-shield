package config

import "fmt"

type RateLimit struct {
	Count  uint32
	Period uint32
}

type AlertGroup struct {
	Name                 string
	Message              string
	RateLimits           []RateLimit
	RateLimitResetPeriod uint32
}

func (g *AlertGroup) RateLimit(level uint64) (rateLimit RateLimit, err error) {
	lenRateLimits := len(g.RateLimits) - 1

	if lenRateLimits < 0 {
		return RateLimit{}, fmt.Errorf("rate limits is empty")
	}

	if level <= uint64(lenRateLimits) {
		rateLimit = g.RateLimits[level]
	} else {
		rateLimit = g.RateLimits[lenRateLimits]
	}

	return rateLimit, nil
}
