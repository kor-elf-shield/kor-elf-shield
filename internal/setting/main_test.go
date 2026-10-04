package setting

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitSetting_InvalidPath(t *testing.T) {
	prev := Config
	t.Cleanup(func() { Config = prev })

	err := InitSetting("relative.toml")
	if err == nil {
		t.Fatalf("expected error for non-absolute path")
	}
}

func TestInitSetting_FileNotFound(t *testing.T) {
	prev := Config
	t.Cleanup(func() { Config = prev })

	missing := filepath.Join(t.TempDir(), "missing.toml")
	err := InitSetting(missing)
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
}

func TestInitSetting_ValidateError(t *testing.T) {
	prev := Config
	t.Cleanup(func() { Config = prev })

	cfgPath := filepath.Join(t.TempDir(), "bad.toml")
	// testing_interval < 1 => Validate() must fail
	writeFile(t, cfgPath, `
		testing_interval = 0
	`)

	err := InitSetting(cfgPath)
	if err == nil {
		t.Fatalf("expected validation error")
	}
	if !strings.Contains(err.Error(), "testing_interval must be greater than 0") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInitSetting_Success(t *testing.T) {
	prev := Config
	t.Cleanup(func() { Config = prev })

	cfgPath := filepath.Join(t.TempDir(), "ok.toml")
	writeFile(t, cfgPath, `
		testing = true
		testing_interval = 10
		language = "en"
		fallback_language = "ru"
		pid_file = "/tmp/kor-elf-shield.pid"
		socket_file = "/tmp/kor-elf-shield.sock"
		data_dir = "/tmp/kor-elf-shield"
	`)

	err := InitSetting(cfgPath)
	if err != nil {
		t.Fatalf("InitSetting returned error: %v", err)
	}

	if Config == nil {
		t.Fatalf("Config is nil after InitSetting")
	}
	if Config.ConfigPath != cfgPath {
		t.Fatalf("unexpected ConfigPath: got=%q want=%q", Config.ConfigPath, cfgPath)
	}
	if Config.TestingInterval != 10 {
		t.Fatalf("unexpected TestingInterval: got=%d want=%d", Config.TestingInterval, 10)
	}
	if Config.Language != "en" {
		t.Fatalf("unexpected Language: got=%q want=%q", Config.Language, "en")
	}
	if Config.FallbackLanguage != "ru" {
		t.Fatalf("unexpected FallbackLanguage: got=%q want=%q", Config.FallbackLanguage, "ru")
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	err := os.WriteFile(path, []byte(strings.TrimSpace(content)+"\n"), 0o644)
	if err != nil {
		t.Fatalf("write file %s: %v", path, err)
	}
}
