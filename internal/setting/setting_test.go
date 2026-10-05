package setting

import (
	"path/filepath"
	"strings"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/i18nmock"
)

func TestSetting_ValidateBeforeStart_Error(t *testing.T) {
	i18nmock.Use(t)

	var s setting

	s = setting{}
	if err := s.ValidateBeforeStart(); err == nil {
		t.Errorf("ValidateBeforeStart() should return an error")
	}

	s = defaultSettingsForValidate()
	s.PidFile = ""
	if err := s.ValidateBeforeStart(); err == nil {
		t.Errorf("ValidateBeforeStart() should return an error")
	} else if !strings.Contains(err.Error(), "pid_file") {
		t.Errorf("ValidateBeforeStart() should return an error, got %v", err)
	}
}

func TestSetting_ValidateBeforeStart_Success(t *testing.T) {
	i18nmock.Use(t)

	s := defaultSettingsForValidate()
	if err := s.ValidateBeforeStart(); err != nil {
		t.Errorf("ValidateBeforeStart() should not return an error")
	}
}

func TestSetting_ListPathConfigFiles_Success(t *testing.T) {
	s := setting{
		ConfigPath: "/tmp/test_kor-elf-shield.toml",
		OtherSettingsPath: &otherSettingsPath{
			Firewall:      "/tmp/test_firewall.toml",
			Notifications: "/tmp/test_notifications.toml",
			Analyzer:      "/tmp/test_analyzer.toml",
			Docker:        "/tmp/test_docker.toml",
			Blocklists:    "/tmp/test_blocklists.toml",
			GeoIP:         "/tmp/test_geoip.toml",
		},
	}

	list := s.ListPathConfigFiles()
	if list == nil || len(list) == 0 {
		t.Errorf("ListPathConfigFiles() should return at least one file")
	}
}

func TestSetting_ListPathConfigFiles_Error(t *testing.T) {
	s := setting{}
	if s.ListPathConfigFiles() != nil {
		t.Errorf("ListPathConfigFiles() should return at least one file")
	}

	s = setting{
		ConfigPath: "/tmp/test_kor-elf-shield.toml",
	}
	if s.ListPathConfigFiles() != nil {
		t.Errorf("ListPathConfigFiles() should return at least one file")
	}
}

func TestSetting_Validate_Success(t *testing.T) {
	i18nmock.Use(t)

	s := defaultSettingsForValidate()
	if err := s.Validate(); err != nil {
		t.Errorf("Validate() should not return an error. %v", err)
	}
}

func TestSetting_Validate_TestingInterval_Success(t *testing.T) {
	i18nmock.Use(t)

	s := defaultSettingsForValidate()
	s.TestingInterval = int16(1)
	if err := s.Validate(); err != nil {
		t.Errorf("Validate() should not return an error. %v", err)
	}

	s.TestingInterval = int16(30000)
	if err := s.Validate(); err != nil {
		t.Errorf("Validate() should not return an error. %v", err)
	}
}

func TestSetting_Validate_PidFile_Success(t *testing.T) {
	i18nmock.Use(t)

	s := defaultSettingsForValidate()
	s.PidFile = "/tmp/test.pid"
	if err := s.Validate(); err != nil {
		t.Errorf("Validate() should not return an error. %v", err)
	}

	s.PidFile = "/tmp/test.PID"
	if err := s.Validate(); err != nil {
		t.Errorf("Validate() should not return an error. %v", err)
	}
}

func TestSetting_Validate_SocketFile_Success(t *testing.T) {
	i18nmock.Use(t)

	s := defaultSettingsForValidate()
	s.SocketFile = "/tmp/test.sock"
	if err := s.Validate(); err != nil {
		t.Errorf("Validate() should not return an error. %v", err)
	}

	s.SocketFile = "/tmp/test.SOCK"
	if err := s.Validate(); err != nil {
		t.Errorf("Validate() should not return an error. %v", err)
	}
}

func TestSetting_Validate_TestingInterval_Error(t *testing.T) {
	i18nmock.Use(t)

	s := defaultSettingsForValidate()
	s.TestingInterval = int16(0)
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}

	s.TestingInterval = int16(-1)
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}

	s.TestingInterval = int16(30001)
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}
}

func TestSetting_Validate_Language_Error(t *testing.T) {
	i18nmock.Use(t)

	s := defaultSettingsForValidate()
	s.Language = ""
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}
}

func TestSetting_Validate_FallbackLanguage_Error(t *testing.T) {
	i18nmock.Use(t)

	s := defaultSettingsForValidate()
	s.FallbackLanguage = ""
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}
}

func TestSetting_Validate_PidFile_Error(t *testing.T) {
	i18nmock.Use(t)

	s := defaultSettingsForValidate()
	s.PidFile = ""
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}

	s.PidFile = "/"
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}

	s.PidFile = "/test/file"
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}

	s.PidFile = "/test/file/"
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}

	s.PidFile = "/tmp/test."
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}
}

func TestSetting_Validate_SocketFile_Error(t *testing.T) {
	i18nmock.Use(t)

	s := defaultSettingsForValidate()
	s.SocketFile = ""
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}

	s.SocketFile = "/"
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}

	s.SocketFile = "/test/file"
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}

	s.SocketFile = "/test/file/"
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}

	s.SocketFile = "/tmp/test."
	if err := s.Validate(); err == nil {
		t.Errorf("Validate() should return an error")
	}
}

func TestSetting_ToDaemonOptions_Success(t *testing.T) {
	i18nmock.Use(t)
	dir := t.TempDir()

	s := defaultSettingsForValidate()
	s.OtherSettingsPath = defaultOtherSettingsPath(dir)

	if err := testutil.SaveTomlEmpty(s.OtherSettingsPath.Firewall); err != nil {
		t.Errorf("SaveToml() should not return an error")
	}

	if err := testutil.SaveTomlEmpty(s.OtherSettingsPath.Analyzer); err != nil {
		t.Errorf("SaveToml() should not return an error")
	}

	if options, err := s.ToDaemonOptions(fakers.RandBool()); err != nil {
		t.Errorf("ToDaemonOptions() should not return an error")
	} else if options.DataDir != s.DataDir {
		t.Errorf("ToDaemonOptions() should return options with DataDir set to %v", s.DataDir)
	} else if options.PathSocketFile != s.SocketFile {
		t.Errorf("ToDaemonOptions() should return options with SocketFile set to %v", s.SocketFile)
	} else if options.PathPidFile != s.PidFile {
		t.Errorf("ToDaemonOptions() should return options with PidFile set to %v", s.PidFile)
	} else if options.PathNftables != s.BinaryLocations.Nftables {
		t.Errorf("ToDaemonOptions() should return options with Nftables set to %v", s.BinaryLocations.Nftables)
	}
}

func TestSetting_ToDaemonOptions_Error(t *testing.T) {
	i18nmock.Use(t)

	s := defaultSettingsForValidate()
	_, err := s.ToDaemonOptions(fakers.RandBool())
	if err == nil {
		t.Errorf("ToDaemonOptions() should return an error")
	}
}

func defaultSettingsForValidate() setting {
	testingInterval := fakers.RandInt(1, 30000)

	return setting{
		TestingInterval:  int16(testingInterval),
		Language:         fakeRandLanguage(),
		FallbackLanguage: fakeRandLanguage(),
		SocketFile:       "/tmp/test_kor-elf-shield.sock",
		PidFile:          "/tmp/test_kor-elf-shield.pid",
		DataDir:          "/tmp/test_kor-elf-shield",
		BinaryLocations: &binaryLocations{
			Nftables:   "/usr/sbin/nft",
			Journalctl: "/bin/journalctl",
			Docker:     "/usr/bin/docker",
		},
	}
}

func defaultOtherSettingsPath(dir string) *otherSettingsPath {
	return &otherSettingsPath{
		Firewall:      filepath.Join(dir, "firewall.toml"),
		Notifications: filepath.Join(dir, "notifications.toml"),
		Analyzer:      filepath.Join(dir, "analyzer.toml"),
		Docker:        filepath.Join(dir, "docker.toml"),
		Blocklists:    filepath.Join(dir, "blocklists.toml"),
		GeoIP:         filepath.Join(dir, "geoip.toml"),
	}
}

func fakeRandLanguage() string {
	languages := []string{"ru", "en", "kk"}
	return fakers.RandItem(languages)
}
