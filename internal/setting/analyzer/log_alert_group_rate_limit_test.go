package analyzer

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestLogAlertGroupRateLimit_ToRateLimit_Success(t *testing.T) {
	s := defaultLogAlertGroupRateLimitForTest()
	if rateLimit, err := s.ToRateLimit(); err != nil {
		t.Errorf("LogAlertGroupRateLimit.ToRateLimit() error = %v, wantErr nil", err)
	} else if rateLimit.Count != uint32(s.Count) {
		t.Errorf("LogAlertGroupRateLimit.ToRateLimit() rateLimit.Count = %v, want %v", rateLimit.Count, uint32(s.Count))
	} else if rateLimit.Period != uint32(s.Period) {
		t.Errorf("LogAlertGroupRateLimit.ToRateLimit() rateLimit.Period = %v, want %v", rateLimit.Period, uint32(s.Period))
	}
}

func TestLogAlertGroupRateLimit_ToRateLimit_Error(t *testing.T) {
	invalidCount := []int{-1, 0, -1000}
	for _, count := range invalidCount {
		t.Run(fmt.Sprintf("invalid Count: %d", count), func(t *testing.T) {
			s := defaultLogAlertGroupRateLimitForTest()
			s.Count = count
			if _, err := s.ToRateLimit(); err == nil {
				t.Errorf("LogAlertGroupRateLimit.ToRateLimit() error = nil, wantErr not nil")
			}
		})
	}

	invalidPeriod := []int{-1, 0, -1000}
	for _, period := range invalidPeriod {
		t.Run(fmt.Sprintf("invalid Period: %d", period), func(t *testing.T) {
			s := defaultLogAlertGroupRateLimitForTest()
			s.Period = period
			if _, err := s.ToRateLimit(); err == nil {
				t.Errorf("LogAlertGroupRateLimit.ToRateLimit() error = nil, wantErr not nil")
			}
		})
	}
}

func defaultLogAlertGroupRateLimitForTest() *LogAlertGroupRateLimit {
	return &LogAlertGroupRateLimit{
		Count:  fakers.RandInt(1, 1000),
		Period: fakers.RandInt(1, 86400),
	}
}
