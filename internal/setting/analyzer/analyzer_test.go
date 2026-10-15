package analyzer

import (
	"path/filepath"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/i18nmock"
)

func TestAnalyzer_InitSetting_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "analyzer.toml")

	if err := testutil.SaveTomlEmpty(path); err != nil {
		t.Fatalf("failed to create toml file: %v", err)
	}

	_, err := InitSetting(path)
	if err != nil {
		t.Errorf("InitSetting() should not return an error, got %v", err)
		return
	}
}

func TestAnalyzer_InitSetting_Error(t *testing.T) {
	t.Run("invalid path", func(t *testing.T) {
		if _, err := InitSetting("analyzer.toml"); err == nil {
			t.Errorf("InitSetting() should return an error for non-absolute path")
		}
	})

	t.Run("file does not exist", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "not_exists.toml")

		if _, err := InitSetting(path); err == nil {
			t.Errorf("InitSetting() should return an error for missing file")
		}
	})

	t.Run("file not TOML", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "not_toml.txt")
		if err := testutil.SaveTomlEmpty(path); err != nil {
			t.Errorf("testutil.SaveTomlEmpty() failed: %v", err)
		}

		if _, err := InitSetting(path); err == nil {
			t.Errorf("InitSetting() should return an error for non-TOML file")
		}
	})
}

func TestAnalyzer_ToSources_Success(t *testing.T) {
	i18nmock.Use(t)

	s := defaultAnalyzerForTest(t.TempDir())
	if _, err := s.ToSources(); err != nil {
		t.Errorf("Setting.ToSources() error = %v", err)
	}
}

func TestAnalyzer_Validate_Success(t *testing.T) {
	s := defaultAnalyzerForTest(t.TempDir())
	if err := s.Validate(); err != nil {
		t.Errorf("Setting.Validate() error = %v", err)
	}
}

func defaultAnalyzerForTest(dir string) *Setting {
	return &Setting{
		Login:                *defaultLoginForTest(),
		LogAlert:             *defaultLogAlertForTest(dir),
		BruteForceProtection: *defaultBruteForceProtectionForTest(dir),
	}
}
