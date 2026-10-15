package analyzer

import (
	"fmt"
	"testing"
)

func TestLogAlertPattern_ToPattern_Success(t *testing.T) {
	s := defaultLogAlertPatternForTest()
	if _, err := s.ToPattern(); err != nil {
		t.Errorf("LogAlertPattern.ToPattern() error = %v, wantErr nil", err)
	}
}

func TestLogAlertPattern_ToPattern_Error(t *testing.T) {
	invalidRegexp := []string{
		"",
		"[",
		"[s-1]",
	}
	for _, regexp := range invalidRegexp {
		t.Run(fmt.Sprintf("invalid Regexp: %s", regexp), func(t *testing.T) {
			s := defaultLogAlertPatternForTest()
			s.Regexp = regexp
			if alertRegexp, err := s.ToPattern(); err == nil {
				if alertRegexp.Regexp == nil {
					t.Errorf("LogAlertPattern.ToPattern() error = nil, wantErr nil")
				} else if _, err := alertRegexp.Regexp.Get(); err == nil {
					t.Errorf("regexp.Compile error = nil, wantErr not nil")
				}
			}
		})
	}
}

func defaultLogAlertPatternForTest() *LogAlertPattern {
	return &LogAlertPattern{
		Regexp:    "test",
		Values:    nil,
		Partition: nil,
	}
}
