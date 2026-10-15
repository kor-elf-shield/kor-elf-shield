package analyzer

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestLogAlertGroup_ToGroup_Success(t *testing.T) {
	s := defaultLogAlertGroupForTest()
	if group, err := s.ToGroup(); err != nil {
		t.Errorf("LogAlertGroup.ToGroup() error = %v, wantErr nil", err)
	} else if group.Name != s.Name {
		t.Errorf("LogAlertGroup.ToGroup() group.Name = %v, want %v", group.Name, s.Name)
	} else if group.Message != s.Message {
		t.Errorf("LogAlertGroup.ToGroup() group.Message = %v, want %v", group.Message, s.Message)
	} else if group.RateLimitResetPeriod != uint32(s.RateLimitResetPeriod) {
		t.Errorf("LogAlertGroup.ToGroup() group.RateLimitResetPeriod = %v, want %v", group.RateLimitResetPeriod, uint32(s.RateLimitResetPeriod))
	} else if len(group.RateLimits) != len(s.RateLimits) {
		t.Errorf("LogAlertGroup.ToGroup() group.RateLimits = %v, want %v", group.RateLimits, s.RateLimits)
	}
}

func TestLogAlertGroup_ToGroup_Error(t *testing.T) {
	t.Run("negative RateLimitResetPeriod", func(t *testing.T) {
		s := defaultLogAlertGroupForTest()
		s.RateLimitResetPeriod = -1
		if _, err := s.ToGroup(); err == nil {
			t.Errorf("LogAlertGroup.ToGroup() error = nil, wantErr not nil")
		}
	})

	t.Run("empty RateLimits", func(t *testing.T) {
		s := defaultLogAlertGroupForTest()
		s.RateLimits = nil
		if _, err := s.ToGroup(); err == nil {
			t.Errorf("LogAlertGroup.ToGroup() error = nil, wantErr not nil")
		}
	})

	invalidNames := []string{
		"",
		" ",
		"1231",
		"12312sdfsdfsd",
		"#sdf54sdf",
		"dfdfgDff5@fdfs",
		fakers.RandString(256),
		fakers.RandString(500),
	}
	for _, name := range invalidNames {
		t.Run(fmt.Sprintf("invalid Name: %s", name), func(t *testing.T) {
			s := defaultLogAlertGroupForTest()
			s.Name = name

			if _, err := s.ToGroup(); err == nil {
				t.Errorf("LogAlertGroup.ToGroup() error = nil, wantErr not nil")
			}
		})
	}
}

func defaultLogAlertGroupForTest() *LogAlertGroup {
	rateLimit := LogAlertGroupRateLimit{
		Count:  fakers.RandInt(1, 10000),
		Period: fakers.RandInt(0, 86400),
	}

	return &LogAlertGroup{
		Name:                 "test",
		Message:              "test",
		RateLimitResetPeriod: fakers.RandInt(0, 86400),
		RateLimits: []LogAlertGroupRateLimit{
			rateLimit,
		},
	}
}
