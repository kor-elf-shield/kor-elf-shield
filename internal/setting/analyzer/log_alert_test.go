package analyzer

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestLogAlert_Validate_Success(t *testing.T) {
	s := defaultLogAlertForTest(t.TempDir())
	if err := s.Validate(); err != nil {
		t.Errorf("LogAlert.Validate() error = %v, wantErr nil", err)
	}
}

func TestLogAlert_ToSources_Success(t *testing.T) {
	t.Run("Enabled = true", func(t *testing.T) {
		s := defaultLogAlertForTest(t.TempDir())
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("LogAlert.ToSources() error = %v, wantErr nil", err)
		} else if len(sources) != len(s.Rules) {
			t.Errorf("LogAlert.ToSources() sources = %v, s.Rules = %v, want len(sources) = len(s.Rules)", len(sources), len(s.Rules))
		}
	})

	t.Run("Enabled = false", func(t *testing.T) {
		s := defaultLogAlertForTest(t.TempDir())
		s.Enabled = false
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("LogAlert.ToSources() error = %v, wantErr nil", err)
		} else if len(sources) != 0 {
			t.Errorf("LogAlert.ToSources() sources = %v, want len(sources) = 0", len(sources))
		}
	})

	t.Run("No rules", func(t *testing.T) {
		s := defaultLogAlertForTest(t.TempDir())
		s.Rules = nil
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("LogAlert.ToSources() error = %v, wantErr nil", err)
		} else if len(sources) != 0 {
			t.Errorf("LogAlert.ToSources() sources = %v, want len(sources) = 0", len(sources))
		}
	})
}

func TestLogAlert_ToSources_Error(t *testing.T) {
	t.Run("invalid Rule Group", func(t *testing.T) {
		s := defaultLogAlertForTest(t.TempDir())
		s.Rules[0].Group = "invalid"
		if _, err := s.ToSources(); err == nil {
			t.Errorf("LogAlert.ToSources() error = nil, wantErr not nil")
		}
	})
}

func defaultLogAlertForTest(dir string) *LogAlert {
	rule := defaultLogAlertRuleForTest(dir)

	return &LogAlert{
		Enabled: true,
		Notify:  fakers.RandBool(),
		Groups:  nil,
		Rules:   []LogAlertRule{*rule},
	}
}
