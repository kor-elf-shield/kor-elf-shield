package firewall

import (
	"fmt"
	"path/filepath"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestOptions_Validate_Success(t *testing.T) {
	dir := t.TempDir()
	s := defaultOptionsForTest(dir)
	if err := s.Validate(); err != nil {
		t.Errorf("options.Validate() error = %v, wantErr nil", err)
	}
}

func TestOptions_Validate_Error(t *testing.T) {
	dir := t.TempDir()
	invalidPaths := []string{
		"",
		" ",
		"test",
		filepath.Join(dir, "nftables"),
		filepath.Join(dir, "test.conf"),
		"./nftables.conf",
		"../nftables.conf",
		"@/nftables.conf",
	}

	for _, path := range invalidPaths {
		t.Run(fmt.Sprintf("SavesRulesPath: %q", path), func(t *testing.T) {
			s := defaultOptionsForTest(dir)
			s.SavesRulesPath = path
			if err := s.Validate(); err == nil {
				t.Errorf("options.Validate() error = nil, wantErr not nil")
			}
		})
	}
}

func TestOptions_ToClearMode_Success(t *testing.T) {
	dir := t.TempDir()
	s := defaultOptionsForTest(dir)
	if clearMode, err := s.ToClearMode(); err != nil {
		t.Errorf("options.ToClearMode() error = %v, wantErr nil", err)
	} else if !testClearMode(clearMode, s.ClearMode) {
		t.Errorf("options.ToClearMode() = %v, want %v", clearMode, s.ClearMode)
	}
}

func TestOptions_ToClearMode_Error(t *testing.T) {
	t.Run("empty ClearMode", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOptionsForTest(dir)
		s.ClearMode = ""
		if _, err := s.ToClearMode(); err == nil {
			t.Errorf("options.ToClearMode() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid ClearMode", func(t *testing.T) {
		dir := t.TempDir()
		s := defaultOptionsForTest(dir)
		s.ClearMode = "invalid"
		if _, err := s.ToClearMode(); err == nil {
			t.Errorf("options.ToClearMode() error = nil, wantErr not nil")
		}
	})
}

func defaultOptionsForTest(dir string) *options {
	return &options{
		Cache:          fakers.RandBool(),
		ClearMode:      fakeRandClearMode(),
		SavesRules:     fakers.RandBool(),
		SavesRulesPath: filepath.Join(dir, "nftables.conf"),
		DnsStrict:      fakers.RandBool(),
		DnsStrictNs:    fakers.RandBool(),
		PacketFilter:   fakers.RandBool(),
	}
}

func fakeRandClearMode() string {
	modes := []string{"global", "own"}
	return fakers.RandItem(modes)
}

func testClearMode(clearMode config.ClearMode, configClearMode string) bool {
	switch clearMode {
	case config.ClearModeOwn:
		return configClearMode == "own"
	case config.ClearModeGlobal:
		return configClearMode == "global"
	}

	return false
}
