package time_operation

import "time"

func IsRateLimited(lastTriggeredAtUnix int64, eventTime time.Time, rateLimit int64) bool {
	if lastTriggeredAtUnix == 0 {
		return true
	}

	return eventTime.Unix()-lastTriggeredAtUnix < rateLimit
}

func IsReset(lastTriggeredAtUnix int64, eventTime time.Time, resetPeriod int64) bool {
	if resetPeriod == 0 || lastTriggeredAtUnix == 0 {
		return false
	}

	if eventTime.Unix()-lastTriggeredAtUnix > resetPeriod {
		return true
	}

	return false
}
