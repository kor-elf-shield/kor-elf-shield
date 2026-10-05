package setting

import (
	"path/filepath"
	"testing"

	log2 "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/i18nmock"
)

func TestLog_ToLoggerOptions_Success(t *testing.T) {
	i18nmock.Use(t)
	dir := t.TempDir()

	s := defaultLogForTest(dir)
	if options, err := s.ToLoggerOptions(); err != nil {
		t.Errorf("ToLoggerOptions() should not return an error. %v", err)
	} else if options.Enabled != s.Enabled {
		t.Errorf("ToLoggerOptions() should return options with Enabled set to %t", s.Enabled)
	} else if options.Development != s.Development {
		t.Errorf("ToLoggerOptions() should return options with Development set to %t", s.Development)
	} else if options.Level.String() != s.Level {
		t.Errorf("ToLoggerOptions() should return options with Level set to %s", s.Level)
	} else if !testEncoderOptions(options.Encoding, s.Encoding) {
		t.Errorf("ToLoggerOptions() should return options with Encoding set to %s", s.Encoding)
	}
}

func TestLog_ToLoggerOptions_Level_Error(t *testing.T) {
	i18nmock.Use(t)
	dir := t.TempDir()

	s := defaultLogForTest(dir)
	s.Level = "invalid"
	if _, err := s.ToLoggerOptions(); err == nil {
		t.Errorf("ToLoggerOptions() should return an error for invalid log level")
	}
}

func TestLog_ToLoggerOptions_Encoding_Error(t *testing.T) {
	i18nmock.Use(t)
	dir := t.TempDir()

	s := defaultLogForTest(dir)
	s.Encoding = "invalid"
	if _, err := s.ToLoggerOptions(); err == nil {
		t.Errorf("ToLoggerOptions() should return an error for invalid log encoding")
	}
}

func TestLog_ToLoggerOptions_Paths_Error(t *testing.T) {
	i18nmock.Use(t)
	dir := t.TempDir()

	s := defaultLogForTest(dir)
	s.Paths = []string{}
	if _, err := s.ToLoggerOptions(); err == nil {
		t.Errorf("ToLoggerOptions() should return an error for empty log paths")
	}

	s.Paths = []string{"/invalid"}
	if _, err := s.ToLoggerOptions(); err == nil {
		t.Errorf("ToLoggerOptions() should return an error for invalid log paths")
	}
}

func TestLog_ToLoggerOptions_LogErrorPaths_Error(t *testing.T) {
	i18nmock.Use(t)
	dir := t.TempDir()

	s := defaultLogForTest(dir)
	s.LogErrorPaths = []string{"/invalid"}
	if _, err := s.ToLoggerOptions(); err == nil {
		t.Errorf("ToLoggerOptions() should return an error for invalid log error paths")
	}
}

func defaultLogForTest(dir string) *log {
	return &log{
		Enabled:       true,
		Level:         fakeRandLogLevel(),
		Encoding:      fakeRandLogEncoding(),
		Development:   fakers.RandBool(),
		Paths:         fakeRandPaths(dir),
		LogErrorPaths: fakeRandPaths(dir),
	}
}

func fakeRandLogLevel() string {
	logLevels := []string{"debug", "info", "warn", "error", "fatal"}
	return fakers.RandItem(logLevels)
}

func fakeRandLogEncoding() string {
	encodings := []string{"json", "text"}
	return fakers.RandItem(encodings)
}

func fakeRandPaths(dir string) []string {
	paths := [][]string{
		{filepath.Join(dir, "kor-elf-shield.log"), filepath.Join(dir, "kor-elf-2.log")},
		{"stdout", filepath.Join(dir, "kor-elf-shield-error.log")},
		{"stdout"},
		{"stderr"},
	}
	return fakers.RandItem(paths)
}

func testEncoderOptions(encoding log2.Encoding, encodingOption string) bool {
	if encoding == log2.ConsoleEncoding {
		return encodingOption == "text"
	}

	return encoding.String() == encodingOption
}
