package analyzer

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestLogAlertRule_ToSource_Success(t *testing.T) {
	s := defaultLogAlertRuleForTest(t.TempDir())

	isNotify := fakers.RandBool()
	if _, err := s.ToSource(isNotify, nil); err != nil {
		t.Errorf("LogAlertRule.ToSource() error = %v, wantErr nil", err)
	}
}

func TestLogAlertRule_ToSource_Error(t *testing.T) {
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
			s := defaultLogAlertRuleForTest(t.TempDir())
			s.Name = name

			isNotify := fakers.RandBool()
			if _, err := s.ToSource(isNotify, nil); err == nil {
				t.Errorf("LogAlertRule.ToSource() error = nil, wantErr not nil")
			}
		})
	}

	t.Run("empty Source", func(t *testing.T) {
		s := defaultLogAlertRuleForTest(t.TempDir())
		s.Source = Source{}

		isNotify := fakers.RandBool()
		if _, err := s.ToSource(isNotify, nil); err == nil {
			t.Errorf("LogAlertRule.ToSource() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Patterns", func(t *testing.T) {
		s := defaultLogAlertRuleForTest(t.TempDir())
		s.Patterns = []LogAlertPattern{}

		isNotify := fakers.RandBool()
		if _, err := s.ToSource(isNotify, nil); err == nil {
			t.Errorf("LogAlertRule.ToSource() error = nil, wantErr not nil")
		}
	})
}

func defaultLogAlertRuleForTest(dir string) *LogAlertRule {
	logAlertPattern := LogAlertPattern{
		Regexp:    "Test",
		Values:    []PatternValue{},
		Partition: nil,
	}

	return &LogAlertRule{
		Enabled: true,
		Notify:  fakers.RandBool(),
		Name:    "test",
		Message: "test",
		Group:   "",
		Source:  *defaultSourceForTest(dir),
		Patterns: []LogAlertPattern{
			logAlertPattern,
		},
	}
}
