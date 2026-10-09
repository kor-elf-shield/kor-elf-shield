package blocklists

import (
	"fmt"
	"path/filepath"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/logger_mock"
)

func TestBlocklist_InitSetting_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "blocklists.toml")

	if err := testutil.SaveTomlEmpty(path); err != nil {
		t.Fatalf("failed to create toml file: %v", err)
	}

	_, err := InitSetting(path)
	if err != nil {
		t.Errorf("InitSetting() should not return an error, got %v", err)
		return
	}
}

func TestBlocklist_InitSetting_Error(t *testing.T) {
	t.Run("invalid path", func(t *testing.T) {
		if _, err := InitSetting("blocklists.toml"); err == nil {
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

func TestBlocklist_ToSources_Success(t *testing.T) {
	s := defaultBlocklistForTest()
	logger := logger_mock.New()
	s.ToSources(logger)
	if len(logger.WarnMessages) > 0 {
		t.Errorf("Setting.ToSources() should not return any warning, got %v", logger.WarnMessages)
	}
	if len(logger.ErrorMessages) > 0 {
		t.Errorf("Setting.ToSources() should not return any error, got %v", logger.ErrorMessages)
	}
	if len(logger.FatalMessages) > 0 {
		t.Errorf("Setting.ToSources() should not return any fatal error, got %v", logger.FatalMessages)
	}
}

func TestBlocklist_ToSources_Error(t *testing.T) {
	invalidExcludeIPs := []string{"", "192.168.1.500", "invalid"}
	for _, ExcludeIPs := range invalidExcludeIPs {
		t.Run(fmt.Sprintf("invalid ExcludeIPs :%s", ExcludeIPs), func(t *testing.T) {
			s := defaultBlocklistForTest()
			s.ExcludeIPs = []string{ExcludeIPs}
			logger := logger_mock.New()
			s.ToSources(logger)
			if len(logger.WarnMessages) == 0 {
				t.Errorf("Setting.ToSources() should return warning for invalid ExcludeIPs, got %v", logger.WarnMessages)
			}
		})
	}
}

func defaultBlocklistForTest() *Setting {
	return &Setting{
		Enabled:    true,
		ExcludeIPs: fakeRandExcludeIPs(),
		Sources:    nil,
	}
}

func fakeRandExcludeIPs() []string {
	items := [][]string{
		{},
		{"1.2.3.4", "5.6.7.8"},
		{"127.0.0.1/8"},
		{"10.0.0.0/8"},
		{"172.16.0.0/12"},
		{"192.168.0.0/16"},
		{"::1/128"},
		{"fc00::/7"},
	}

	return fakers.RandItem(items)
}
