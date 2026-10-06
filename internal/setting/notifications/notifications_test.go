package notifications

import (
	"fmt"
	"path/filepath"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestNotifications_InitSetting_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "notifications.toml")

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
		t.Errorf("InitSetting() default Enabled should be %v", setDefault.Enabled)
	}
	if setting.EnableRetries != setDefault.EnableRetries {
		t.Errorf("InitSetting() default EnableRetries should be %v", setDefault.EnableRetries)
	}
	if setting.RetryInterval != setDefault.RetryInterval {
		t.Errorf("InitSetting() default RetryInterval should be %d, got %d", setDefault.RetryInterval, setting.RetryInterval)
	}
	if setting.ServerName != setDefault.ServerName {
		t.Errorf("InitSetting() default ServerName should be '%s', got %q", setDefault.ServerName, setting.ServerName)
	}
	if setting.Email.AuthType != setDefault.Email.AuthType {
		t.Errorf("InitSetting() default Email.AuthType should be '%s', got %q", setDefault.Email.AuthType, setting.Email.AuthType)
	}
}

func TestNotifications_InitSetting_Error(t *testing.T) {
	t.Run("invalid path", func(t *testing.T) {
		if _, err := InitSetting("notifications.toml"); err == nil {
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

func TestNotifications_Setting_Validate_Success(t *testing.T) {
	s := defaultSettingForTest()
	if err := s.Validate(); err != nil {
		t.Errorf("Setting.Validate() should not return an error for default settings: %v", err)
	}
}

func TestNotifications_Setting_Validate_Error(t *testing.T) {
	t.Run("invalid retry interval", func(t *testing.T) {
		s := defaultSettingForTest()
		s.RetryInterval = -1
		if err := s.Validate(); err == nil {
			t.Errorf("Setting.Validate() should return an error for invalid retry interval")
		}
	})

	t.Run("empty ServerName", func(t *testing.T) {
		s := defaultSettingForTest()
		s.ServerName = ""
		if err := s.Validate(); err == nil {
			t.Errorf("Setting.Validate() should return an error for empty ServerName")
		}
	})

	invalidServerName := []string{
		"",
		" ",
		"name/with/slash",
		"name with space",
		"name@with@at",
		"привет", // Cyrillic
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", // 65 symbols
	}
	for _, name := range invalidServerName {
		t.Run(fmt.Sprintf("invalid ServerName: %s", name), func(t *testing.T) {
			s := defaultSettingForTest()
			s.ServerName = name
			if err := s.Validate(); err == nil {
				t.Errorf("Setting.Validate() should return an error for invalid ServerName: %s", name)
			}
		})
	}
}

func defaultSettingForTest() *Setting {
	retryInterval := fakers.RandInt(1, 32767)

	return &Setting{
		Enabled:       true,
		EnableRetries: fakers.RandBool(),
		RetryInterval: int16(retryInterval),
		ServerName:    "test",
		Email: Email{
			Host:      "test",
			Port:      465,
			Username:  "test",
			Password:  "test",
			AuthType:  "PLAIN",
			TLSMode:   "STARTTLS",
			TLSPolicy: "MANDATORY",
			TLSVerify: true,
			From:      "test@localhost",
			To:        "test@localhost",
		},
	}
}
