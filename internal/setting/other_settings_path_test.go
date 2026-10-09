package setting

import (
	"path/filepath"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/i18nmock"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/logger_mock"
)

func TestOtherSettingsPath_ToFirewallConfig_Success(t *testing.T) {
	s := defaultOtherSettingsPathForTest(t.TempDir())
	if err := testutil.SaveTomlEmpty(s.Firewall); err != nil {
		t.Errorf("SaveTomlEmpty() error = %v", err)
	}

	if _, _, err := s.ToFirewallConfig(fakers.RandBool()); err != nil {
		t.Errorf("otherSettingsPath.ToFirewallConfig() error = %v", err)
	}
}

func TestOtherSettingsPath_ToFirewallConfig_Error(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		s := otherSettingsPath{
			Firewall:      "",
			Notifications: "",
			Analyzer:      "",
			Docker:        "",
			Blocklists:    "",
			GeoIP:         "",
		}
		if _, _, err := s.ToFirewallConfig(fakers.RandBool()); err == nil {
			t.Errorf("otherSettingsPath.ToFirewallConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		s := defaultOtherSettingsPathForTest(t.TempDir())
		s.Firewall = "invalid"
		if _, _, err := s.ToFirewallConfig(fakers.RandBool()); err == nil {
			t.Errorf("otherSettingsPath.ToFirewallConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("path not toml", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOtherSettingsPathForTest(dir)
		s.Firewall = filepath.Join(dir, "firewall.nottoml")
		if err := testutil.SaveTomlEmpty(s.Firewall); err != nil {
			t.Errorf("SaveTomlEmpty() error = %v", err)
		}
		if _, _, err := s.ToFirewallConfig(fakers.RandBool()); err == nil {
			t.Errorf("otherSettingsPath.ToFirewallConfig() error = nil, wantErr not nil")
		}
	})
}

func TestOtherSettingsPath_ToNotificationsConfig_Success(t *testing.T) {
	s := defaultOtherSettingsPathForTest(t.TempDir())
	if err := testutil.SaveTomlEmpty(s.Notifications); err != nil {
		t.Errorf("SaveTomlEmpty() error = %v", err)
	}

	if _, err := s.ToNotificationsConfig(); err != nil {
		t.Errorf("otherSettingsPath.ToNotificationsConfig() error = %v", err)
	}
}

func TestOtherSettingsPath_ToNotificationsConfig_Error(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		s := otherSettingsPath{
			Firewall:      "",
			Notifications: "",
			Analyzer:      "",
			Docker:        "",
			Blocklists:    "",
			GeoIP:         "",
		}
		if _, err := s.ToNotificationsConfig(); err == nil {
			t.Errorf("otherSettingsPath.ToNotificationsConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		s := defaultOtherSettingsPathForTest(t.TempDir())
		s.Notifications = "invalid"
		if _, err := s.ToNotificationsConfig(); err == nil {
			t.Errorf("otherSettingsPath.ToNotificationsConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("path not toml", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOtherSettingsPathForTest(dir)
		s.Notifications = filepath.Join(dir, "notifications.nottoml")
		if err := testutil.SaveTomlEmpty(s.Firewall); err != nil {
			t.Errorf("SaveTomlEmpty() error = %v", err)
		}
		if _, err := s.ToNotificationsConfig(); err == nil {
			t.Errorf("otherSettingsPath.ToNotificationsConfig() error = nil, wantErr not nil")
		}
	})
}

func TestOtherSettingsPath_ToAnalyzerConfig_Success(t *testing.T) {
	i18nmock.Use(t)

	s := defaultOtherSettingsPathForTest(t.TempDir())
	if err := testutil.SaveTomlEmpty(s.Analyzer); err != nil {
		t.Errorf("SaveTomlEmpty() error = %v", err)
	}
	binary := defaultBinaryLocationsForTest()

	if _, err := s.ToAnalyzerConfig(binary); err != nil {
		t.Errorf("otherSettingsPath.ToAnalyzerConfig() error = %v", err)
	}
}

func TestOtherSettingsPath_ToAnalyzerConfig_Error(t *testing.T) {
	i18nmock.Use(t)

	t.Run("empty path", func(t *testing.T) {
		s := otherSettingsPath{
			Firewall:      "",
			Notifications: "",
			Analyzer:      "",
			Docker:        "",
			Blocklists:    "",
			GeoIP:         "",
		}
		binary := defaultBinaryLocationsForTest()
		if _, err := s.ToAnalyzerConfig(binary); err == nil {
			t.Errorf("otherSettingsPath.ToAnalyzerConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		s := defaultOtherSettingsPathForTest(t.TempDir())
		s.Analyzer = "invalid"
		binary := defaultBinaryLocationsForTest()
		if _, err := s.ToAnalyzerConfig(binary); err == nil {
			t.Errorf("otherSettingsPath.ToAnalyzerConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("path not toml", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOtherSettingsPathForTest(dir)
		s.Analyzer = filepath.Join(dir, "analyzer.nottoml")
		if err := testutil.SaveTomlEmpty(s.Firewall); err != nil {
			t.Errorf("SaveTomlEmpty() error = %v", err)
		}
		binary := defaultBinaryLocationsForTest()
		if _, err := s.ToAnalyzerConfig(binary); err == nil {
			t.Errorf("otherSettingsPath.ToAnalyzerConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("empty path nftables", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOtherSettingsPathForTest(dir)
		if err := testutil.SaveTomlEmpty(s.Firewall); err != nil {
			t.Errorf("SaveTomlEmpty() error = %v", err)
		}
		binary := defaultBinaryLocationsForTest()
		binary.Nftables = ""
		if _, err := s.ToAnalyzerConfig(binary); err == nil {
			t.Errorf("otherSettingsPath.ToAnalyzerConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid path nftables", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOtherSettingsPathForTest(dir)
		if err := testutil.SaveTomlEmpty(s.Firewall); err != nil {
			t.Errorf("SaveTomlEmpty() error = %v", err)
		}
		binary := defaultBinaryLocationsForTest()
		binary.Nftables = "invalid"
		if _, err := s.ToAnalyzerConfig(binary); err == nil {
			t.Errorf("otherSettingsPath.ToAnalyzerConfig() error = nil, wantErr not nil")
		}
	})
}

func TestOtherSettingsPath_ToDockerConfig_Success(t *testing.T) {
	i18nmock.Use(t)

	s := defaultOtherSettingsPathForTest(t.TempDir())
	if err := testutil.SaveTomlEmpty(s.Docker); err != nil {
		t.Errorf("SaveTomlEmpty() error = %v", err)
	}
	binary := defaultBinaryLocationsForTest()

	if _, _, err := s.ToDockerConfig(binary); err != nil {
		t.Errorf("otherSettingsPath.ToDockerConfig() error = %v", err)
	}
}

func TestOtherSettingsPath_ToDockerConfig_Error(t *testing.T) {
	i18nmock.Use(t)

	t.Run("empty path", func(t *testing.T) {
		s := otherSettingsPath{
			Firewall:      "",
			Notifications: "",
			Analyzer:      "",
			Docker:        "",
			Blocklists:    "",
			GeoIP:         "",
		}
		binary := defaultBinaryLocationsForTest()
		if _, _, err := s.ToDockerConfig(binary); err == nil {
			t.Errorf("otherSettingsPath.ToDockerConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		s := defaultOtherSettingsPathForTest(t.TempDir())
		s.Docker = "invalid"
		binary := defaultBinaryLocationsForTest()
		if _, _, err := s.ToDockerConfig(binary); err == nil {
			t.Errorf("otherSettingsPath.ToDockerConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("path not toml", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOtherSettingsPathForTest(dir)
		s.Docker = filepath.Join(dir, "docker.nottoml")
		if err := testutil.SaveTomlEmpty(s.Firewall); err != nil {
			t.Errorf("SaveTomlEmpty() error = %v", err)
		}
		binary := defaultBinaryLocationsForTest()
		if _, _, err := s.ToDockerConfig(binary); err == nil {
			t.Errorf("otherSettingsPath.ToDockerConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("empty path docker", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOtherSettingsPathForTest(dir)
		if err := testutil.SaveTomlEmpty(s.Firewall); err != nil {
			t.Errorf("SaveTomlEmpty() error = %v", err)
		}
		binary := defaultBinaryLocationsForTest()
		binary.Docker = ""
		if _, _, err := s.ToDockerConfig(binary); err == nil {
			t.Errorf("otherSettingsPath.ToDockerConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid path docker", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOtherSettingsPathForTest(dir)
		if err := testutil.SaveTomlEmpty(s.Firewall); err != nil {
			t.Errorf("SaveTomlEmpty() error = %v", err)
		}
		binary := defaultBinaryLocationsForTest()
		binary.Docker = "invalid"
		if _, _, err := s.ToDockerConfig(binary); err == nil {
			t.Errorf("otherSettingsPath.ToDockerConfig() error = nil, wantErr not nil")
		}
	})
}

func TestOtherSettingsPath_ToBlocklistConfig_Success(t *testing.T) {
	i18nmock.Use(t)

	s := defaultOtherSettingsPathForTest(t.TempDir())
	if err := testutil.SaveTomlEmpty(s.Blocklists); err != nil {
		t.Errorf("SaveTomlEmpty() error = %v", err)
	}

	logger := logger_mock.New()
	if _, _, err := s.ToBlocklistConfig(logger); err != nil {
		t.Errorf("otherSettingsPath.ToBlocklistConfig() error = %v", err)
	}
}

func TestOtherSettingsPath_ToBlocklistConfig_Error(t *testing.T) {
	i18nmock.Use(t)

	t.Run("empty path", func(t *testing.T) {
		s := otherSettingsPath{
			Firewall:      "",
			Notifications: "",
			Analyzer:      "",
			Docker:        "",
			Blocklists:    "",
			GeoIP:         "",
		}
		logger := logger_mock.New()
		if _, _, err := s.ToBlocklistConfig(logger); err == nil {
			t.Errorf("otherSettingsPath.ToBlocklistConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		s := defaultOtherSettingsPathForTest(t.TempDir())
		s.Blocklists = "invalid"
		logger := logger_mock.New()
		if _, _, err := s.ToBlocklistConfig(logger); err == nil {
			t.Errorf("otherSettingsPath.ToBlocklistConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("path not toml", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOtherSettingsPathForTest(dir)
		s.Blocklists = filepath.Join(dir, "blocklists.nottoml")
		if err := testutil.SaveTomlEmpty(s.Firewall); err != nil {
			t.Errorf("SaveTomlEmpty() error = %v", err)
		}
		logger := logger_mock.New()
		if _, _, err := s.ToBlocklistConfig(logger); err == nil {
			t.Errorf("otherSettingsPath.ToBlocklistConfig() error = nil, wantErr not nil")
		}
	})
}

func TestOtherSettingsPath_ToGeoIPConfig_Success(t *testing.T) {
	dir := t.TempDir()

	s := defaultOtherSettingsPathForTest(dir)
	if err := testutil.SaveTomlEmpty(s.GeoIP); err != nil {
		t.Errorf("SaveTomlEmpty() error = %v", err)
	}

	logger := logger_mock.New()
	if _, _, err := s.ToGeoIPConfig(dir, logger); err != nil {
		t.Errorf("otherSettingsPath.ToGeoIPConfig() error = %v", err)
	}
}

func TestOtherSettingsPath_ToGeoIPConfig_Error(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		dir := t.TempDir()
		s := otherSettingsPath{
			Firewall:      "",
			Notifications: "",
			Analyzer:      "",
			Docker:        "",
			Blocklists:    "",
			GeoIP:         "",
		}
		logger := logger_mock.New()
		if _, _, err := s.ToGeoIPConfig(dir, logger); err == nil {
			t.Errorf("otherSettingsPath.ToBlocklistConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOtherSettingsPathForTest(dir)
		s.GeoIP = "invalid"
		logger := logger_mock.New()
		if _, _, err := s.ToGeoIPConfig(dir, logger); err == nil {
			t.Errorf("otherSettingsPath.ToBlocklistConfig() error = nil, wantErr not nil")
		}
	})

	t.Run("path not toml", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOtherSettingsPathForTest(dir)
		s.GeoIP = filepath.Join(dir, "geoip.nottoml")
		if err := testutil.SaveTomlEmpty(s.GeoIP); err != nil {
			t.Errorf("SaveTomlEmpty() error = %v", err)
		}
		logger := logger_mock.New()
		if _, _, err := s.ToGeoIPConfig(dir, logger); err == nil {
			t.Errorf("otherSettingsPath.ToGeoIPConfig() error = nil, wantErr not nil")
		}
	})
}

func TestOtherSettingsPath_ListPathFiles_Success(t *testing.T) {
	s := defaultOtherSettingsPathForTest(t.TempDir())
	files := s.ListPathFiles()
	if len(files) != 6 {
		t.Errorf("otherSettingsPath.ListPathFiles() len = %d, want 6", len(files))
	}

	if files["firewall"] != s.Firewall {
		t.Errorf("otherSettingsPath.ListPathFiles() files[\"firewall\"] = %s, want %s", files["firewall"], s.Firewall)
	}

	if files["notifications"] != s.Notifications {
		t.Errorf("otherSettingsPath.ListPathFiles() files[\"notifications\"] = %s, want %s", files["notifications"], s.Notifications)
	}

	if files["analyzer"] != s.Analyzer {
		t.Errorf("otherSettingsPath.ListPathFiles() files[\"analyzer\"] = %s, want %s", files["analyzer"], s.Analyzer)
	}

	if files["docker"] != s.Docker {
		t.Errorf("otherSettingsPath.ListPathFiles() files[\"docker\"] = %s, want %s", files["docker"], s.Docker)
	}

	if files["blocklists"] != s.Blocklists {
		t.Errorf("otherSettingsPath.ListPathFiles() files[\"blocklists\"] = %s, want %s", files["blocklists"], s.Blocklists)
	}

	if files["geoip"] != s.GeoIP {
		t.Errorf("otherSettingsPath.ListPathFiles() files[\"geoip\"] = %s, want %s", files["geoip"], s.GeoIP)
	}
}

func defaultOtherSettingsPathForTest(dir string) *otherSettingsPath {
	return &otherSettingsPath{
		Firewall:      filepath.Join(dir, "firewall.toml"),
		Notifications: filepath.Join(dir, "notifications.toml"),
		Analyzer:      filepath.Join(dir, "analyzer.toml"),
		Docker:        filepath.Join(dir, "docker.toml"),
		Blocklists:    filepath.Join(dir, "blocklists.toml"),
		GeoIP:         filepath.Join(dir, "geoip.toml"),
	}
}

func defaultBinaryLocationsForTest() *binaryLocations {
	return &binaryLocations{
		Nftables:   "/usr/sbin/nft",
		Journalctl: "/bin/journalctl",
		Docker:     "/usr/bin/docker",
	}
}
