package analyzer

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config/brute_force_protection"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestBruteForceProtectionRule_ToSource_Success(t *testing.T) {
	group := fakeRandGroup()
	s := defaultBruteForceProtectionRuleForTest(t.TempDir())
	isNotify := fakers.RandBool()
	if _, err := s.ToSource(isNotify, group); err != nil {
		t.Errorf("BruteForceProtectionRule.ToSource() error = %v, wantErr nil", err)
	}
}

func TestBruteForceProtectionRule_ToSource_Error(t *testing.T) {
	t.Run("negative NotifyCooldown", func(t *testing.T) {
		s := defaultBruteForceProtectionRuleForTest(t.TempDir())
		s.NotifyCooldown = -1
		if _, err := s.ToSource(false, nil); err == nil {
			t.Errorf("BruteForceProtectionRule.ToSource() error = nil, wantErr not nil")
		}
	})

	t.Run("negative NotifyEvery", func(t *testing.T) {
		s := defaultBruteForceProtectionRuleForTest(t.TempDir())
		s.NotifyEvery = -1
		if _, err := s.ToSource(false, nil); err == nil {
			t.Errorf("BruteForceProtectionRule.ToSource() error = nil, wantErr not nil")
		}
	})

	invalidNames := []string{
		"",
		" ",
		"1231",
		"12312sdfsdfsd",
		"#sdf54sdf",
		"dfdfgDff5@fdfs",
		fakers.RandString(256),
		fakers.RandString(500),
	}
	for _, name := range invalidNames {
		t.Run(fmt.Sprintf("invalid Name: %s", name), func(t *testing.T) {
			s := defaultBruteForceProtectionRuleForTest(t.TempDir())
			s.Name = name

			isNotify := fakers.RandBool()
			if _, err := s.ToSource(isNotify, nil); err == nil {
				t.Errorf("BruteForceProtectionRule.ToSource() error = nil, wantErr not nil")
			}
		})
	}

	t.Run("empty Source", func(t *testing.T) {
		s := defaultBruteForceProtectionRuleForTest(t.TempDir())
		s.Source = Source{}

		isNotify := fakers.RandBool()
		if _, err := s.ToSource(isNotify, nil); err == nil {
			t.Errorf("BruteForceProtectionRule.ToSource() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Patterns", func(t *testing.T) {
		s := defaultBruteForceProtectionRuleForTest(t.TempDir())
		s.Patterns = []BruteForceProtectionPattern{}

		isNotify := fakers.RandBool()
		if _, err := s.ToSource(isNotify, nil); err == nil {
			t.Errorf("BruteForceProtectionRule.ToSource() error = nil, wantErr not nil")
		}
	})
}

func defaultBruteForceProtectionRuleForTest(dir string) *BruteForceProtectionRule {
	pattern := BruteForceProtectionPattern{
		Regexp:    "test",
		IP:        fakers.RandInt(1, 255),
		Values:    nil,
		Partition: nil,
	}

	return &BruteForceProtectionRule{
		Enabled:        true,
		Notify:         fakers.RandBool(),
		NotifyCooldown: fakers.RandInt(0, 86400),
		NotifyEvery:    fakers.RandInt(0, 10000),
		Name:           "test",
		Message:        "test",
		Group:          "",
		Source:         *defaultSourceForTest(dir),
		Patterns:       []BruteForceProtectionPattern{pattern},
	}
}

func fakeRandGroup() *brute_force_protection.Group {
	rateLimitResetPeriod := fakers.RandInt(0, 86400)

	return &brute_force_protection.Group{
		Name:                 "test",
		Message:              "test",
		RateLimits:           nil,
		RateLimitResetPeriod: uint32(rateLimitResetPeriod),
	}
}
