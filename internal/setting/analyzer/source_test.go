package analyzer

import (
	"fmt"
	"path/filepath"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestSource_ToSource_Success(t *testing.T) {
	s := defaultSourceForTest(t.TempDir())
	if _, err := s.ToSource(); err != nil {
		t.Errorf("Source.ToSource() error = %v, wantErr nil", err)
	}
}

func TestSource_ToSource_Error(t *testing.T) {
	t.Run("empty Type", func(t *testing.T) {
		s := defaultSourceForTest(t.TempDir())
		s.Type = ""
		if _, err := s.ToSource(); err == nil {
			t.Errorf("Source.ToSource() error = %v, wantErr not nil", err)
		}
	})

	t.Run("invalid Type", func(t *testing.T) {
		s := defaultSourceForTest(t.TempDir())
		s.Type = "invalid"
		if _, err := s.ToSource(); err == nil {
			t.Errorf("Source.ToSource() error = %v, wantErr not nil", err)
		}
	})

	invalidPath := []string{
		"invalid.log",
		"",
		" ",
		"./test.log",
		"../test.log",
		"test",
	}
	for _, path := range invalidPath {
		t.Run(fmt.Sprintf("invalid Path: %s", path), func(t *testing.T) {
			s := &Source{
				Type: "file",
				Path: path,
			}
			if _, err := s.ToSource(); err == nil {
				t.Errorf("Source.ToSource() error = %v, wantErr not nil", err)
			}
		})
	}

	invalidField := []string{
		"", " ", "invalid",
	}
	for _, field := range invalidField {
		t.Run(fmt.Sprintf("invalid Field: %s", field), func(t *testing.T) {
			s := &Source{
				Type:  "file",
				Field: field,
				Match: "test.service",
			}
			if _, err := s.ToSource(); err == nil {
				t.Errorf("Source.ToSource() error = %v, wantErr not nil", err)
			}
		})
	}

	invalidMatchFieldSystemdUnit := []string{
		"", " ", "invalid", "=.service", "+.service", "\n.service", "\r.service",
	}
	for _, match := range invalidMatchFieldSystemdUnit {
		t.Run(fmt.Sprintf("invalid Match: %s", match), func(t *testing.T) {
			s := &Source{
				Type:  "journalctl",
				Field: "systemd_unit",
				Match: match,
			}
			if _, err := s.ToSource(); err == nil {
				t.Errorf("Source.ToSource() error = %v, wantErr not nil", err)
			}
		})
	}

	invalidMatchFieldSyslogIdentifier := []string{
		"", " ", "=", "+", "\n", "\r",
	}

	for _, match := range invalidMatchFieldSyslogIdentifier {
		t.Run(fmt.Sprintf("invalid Match: %s", match), func(t *testing.T) {
			s := &Source{
				Type:  "journalctl",
				Field: "syslog_identifier",
				Match: match,
			}
			if _, err := s.ToSource(); err == nil {
				t.Errorf("Source.ToSource() error = %v, wantErr not nil", err)
			}
		})
	}
}

func defaultSourceForTest(dir string) *Source {
	source := &Source{
		Type: fakeRandSourceType(),
	}

	if source.Type == "file" {
		source.Path = filepath.Join(dir, "test.log")
	}

	if source.Type == "journalctl" {
		source.Field = fakeRandSourceField()
		switch source.Field {
		case "systemd_unit":
			source.Match = fakeRandSourceMatchFieldSystemdUnit()
		case "syslog_identifier":
			source.Match = fakeRandSourceMatchFieldSyslogIdentifier()
		}
	}

	return source
}

func fakeRandSourceType() string {
	items := []string{"journalctl", "file"}
	return fakers.RandItem(items)
}

func fakeRandSourceField() string {
	items := []string{"systemd_unit", "syslog_identifier"}
	return fakers.RandItem(items)
}

func fakeRandSourceMatchFieldSystemdUnit() string {
	items := []string{
		"test.service", "test.socket", "test.target",
		"test.mount", "test.timer", "test.path",
		"test.scope", "test.slice", "test.device",
	}
	return fakers.RandItem(items)
}

func fakeRandSourceMatchFieldSyslogIdentifier() string {
	items := []string{"test", "syslog_identifier"}
	return fakers.RandItem(items)
}
