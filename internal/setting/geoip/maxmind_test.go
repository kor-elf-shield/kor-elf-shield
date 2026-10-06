package geoip

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/geoip2/service/maxmind/mmdb"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/logger_mock"
)

func TestMaxMind_ToConfig_Success(t *testing.T) {
	dataDir := t.TempDir()
	logger := logger_mock.New()

	s := defaultMaxMindForTest()
	if geoip, err := s.ToConfig(dataDir, logger); err != nil {
		t.Errorf("ToConfig() error = %v", err)
	} else if geoip == nil {
		t.Errorf("ToConfig() geoip = nil")
	}
}

func TestMaxMind_ToConfig_Error(t *testing.T) {
	t.Run("empty Username", func(t *testing.T) {
		dataDir := t.TempDir()
		logger := logger_mock.New()

		s := defaultMaxMindForTest()
		s.Username = ""
		if _, err := s.ToConfig(dataDir, logger); err == nil {
			t.Errorf("ToConfig() error = nil, want error")
		}
	})

	t.Run("empty Password", func(t *testing.T) {
		dataDir := t.TempDir()
		logger := logger_mock.New()

		s := defaultMaxMindForTest()
		s.Password = ""
		if _, err := s.ToConfig(dataDir, logger); err == nil {
			t.Errorf("ToConfig() error = nil, want error")
		}
	})

	t.Run("zero Interval", func(t *testing.T) {
		dataDir := t.TempDir()
		logger := logger_mock.New()

		s := defaultMaxMindForTest()
		s.Interval = 0
		if _, err := s.ToConfig(dataDir, logger); err == nil {
			t.Errorf("ToConfig() error = nil, want error")
		}
	})

	t.Run("negative Interval", func(t *testing.T) {
		dataDir := t.TempDir()
		logger := logger_mock.New()

		s := defaultMaxMindForTest()
		s.Interval = -1
		if _, err := s.ToConfig(dataDir, logger); err == nil {
			t.Errorf("ToConfig() error = nil, want error")
		}
	})

	t.Run("empty Language", func(t *testing.T) {
		dataDir := t.TempDir()
		logger := logger_mock.New()

		s := defaultMaxMindForTest()
		s.Language = ""
		if _, err := s.ToConfig(dataDir, logger); err == nil {
			t.Errorf("ToConfig() error = nil, want error")
		}
	})

	t.Run("invalide Language", func(t *testing.T) {
		dataDir := t.TempDir()
		logger := logger_mock.New()

		s := defaultMaxMindForTest()
		s.Language = "invalid"
		if _, err := s.ToConfig(dataDir, logger); err == nil {
			t.Errorf("ToConfig() error = nil, want error")
		}
	})

	t.Run("invalide URL", func(t *testing.T) {
		dataDir := t.TempDir()
		logger := logger_mock.New()

		s := defaultMaxMindForTest()
		s.URL = "invalid"
		if _, err := s.ToConfig(dataDir, logger); err == nil {
			t.Errorf("ToConfig() error = nil, want error")
		}
	})
}

func defaultMaxMindForTest() *Maxmind {
	return &Maxmind{
		Username: "test",
		Password: "test",
		Interval: fakers.RandInt(60, 86400),
		URL:      mmdb.DownloadURLCityLite,
		Language: fakeLanguageForTest(),
	}
}

func fakeLanguageForTest() string {
	language := []string{
		"Russian", "English",
		"Spanish", "French",
		"German", "Japanese",
		"Brazilian Portuguese",
		"Simplified Chinese",
	}
	return fakers.RandItem(language)
}
