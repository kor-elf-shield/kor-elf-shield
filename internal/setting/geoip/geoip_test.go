package geoip

import (
	"path/filepath"
	"testing"

	"git.kor-elf.net/kor-elf-shield/geoip2/service/maxmind/mmdb"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/logger_mock"
)

func TestGeoIP_InitSetting_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "geoip.toml")

	if err := testutil.SaveTomlEmpty(path); err != nil {
		t.Fatalf("failed to create toml file: %v", err)
	}

	setting, err := InitSetting(path)
	if err != nil {
		t.Errorf("InitSetting() should not return an error, got %v", err)
		return
	}

	setDefault := settingDefault()
	if setting.Enabled != setDefault.Enabled {
		t.Errorf("InitSetting() should return the default setting, got %v", setting)
	}

	if setting.Service != setDefault.Service {
		t.Errorf("InitSetting() should return the default setting, got %v", setting)
	}

	if setting.Maxmind != setDefault.Maxmind {
		t.Errorf("InitSetting() should return the default setting, got %v", setting)
	}
}

func TestGeoIP_InitSetting_Error(t *testing.T) {
	t.Run("invalid path", func(t *testing.T) {
		if _, err := InitSetting("geoip.toml"); err == nil {
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

func TestGeoIP_ToConfig_Success(t *testing.T) {
	dataDir := t.TempDir()
	logger := logger_mock.New()

	s := defaultSettingForTest()
	if _, err := s.ToConfig(dataDir, logger); err != nil {
		t.Errorf("ToConfig() should not return an error")
	}
}

func TestGeoIP_ToConfig_Error(t *testing.T) {
	t.Run("empty Service", func(t *testing.T) {
		s := defaultSettingForTest()
		s.Service = ""
		if _, err := s.ToConfig("", logger_mock.New()); err == nil {
			t.Errorf("ToConfig() should return an error for empty Service")
		}
	})

	t.Run("invalid Service", func(t *testing.T) {
		s := defaultSettingForTest()
		s.Service = "invalid"
		if _, err := s.ToConfig("", logger_mock.New()); err == nil {
			t.Errorf("ToConfig() should return an error for invalid Service")
		}
	})
}

func defaultSettingForTest() *Setting {
	return &Setting{
		Enabled: true,
		Service: "maxmind",
		Maxmind: Maxmind{
			Username: "test",
			Password: "test",
			Interval: fakers.RandInt(60, 86400),
			URL:      mmdb.DownloadURLCityLite,
			Language: "Russian",
		},
	}
}
