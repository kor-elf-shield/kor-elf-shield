package analyzer

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/i18nmock"
)

func TestBruteForceProtection_Validate_Success(t *testing.T) {
	s := defaultBruteForceProtectionForTest(t.TempDir())
	if err := s.Validate(); err != nil {
		t.Errorf("BruteForceProtection.Validate() error = %v, wantErr nil", err)
	}
}

func TestBruteForceProtection_Validate_Error(t *testing.T) {
	invalidRateLimitCount := []int{-1, -1000, 0}
	for _, rateLimitCount := range invalidRateLimitCount {
		t.Run(fmt.Sprintf("invalid RateLimitCount: %d", rateLimitCount), func(t *testing.T) {
			s := defaultBruteForceProtectionForTest(t.TempDir())
			s.RateLimitCount = rateLimitCount
			if err := s.Validate(); err == nil {
				t.Errorf("BruteForceProtection.Validate() error = nil, wantErr not nil")
			}
		})
	}

	invalidRateLimitPeriod := []int{-1, -1000, 0}
	for _, rateLimitPeriod := range invalidRateLimitPeriod {
		t.Run(fmt.Sprintf("invalid RateLimitPeriod: %d", rateLimitPeriod), func(t *testing.T) {
			s := defaultBruteForceProtectionForTest(t.TempDir())
			s.RateLimitPeriod = rateLimitPeriod
			if err := s.Validate(); err == nil {
				t.Errorf("BruteForceProtection.Validate() error = nil, wantErr not nil")
			}
		})
	}

	invalidRateLimitResetPeriod := []int{-1, -1000}
	for _, rateLimitResetPeriod := range invalidRateLimitResetPeriod {
		t.Run(fmt.Sprintf("invalid RateLimitResetPeriod: %d", rateLimitResetPeriod), func(t *testing.T) {
			s := defaultBruteForceProtectionForTest(t.TempDir())
			s.RateLimitResetPeriod = rateLimitResetPeriod
			if err := s.Validate(); err == nil {
				t.Errorf("BruteForceProtection.Validate() error = nil, wantErr not nil")
			}
		})
	}

	invalidBlockingTime := []int{-1, -1000}
	for _, blockingTime := range invalidBlockingTime {
		t.Run(fmt.Sprintf("invalid BlockingTime: %d", blockingTime), func(t *testing.T) {
			s := defaultBruteForceProtectionForTest(t.TempDir())
			s.BlockingTime = blockingTime
			if err := s.Validate(); err == nil {
				t.Errorf("BruteForceProtection.Validate() error = nil, wantErr not nil")
			}
		})
	}

	invalidSSHNotifyCooldown := []int{-1, -1000}
	for _, sshNotifyCooldown := range invalidSSHNotifyCooldown {
		t.Run(fmt.Sprintf("invalid SSHNotifyCooldown: %d", sshNotifyCooldown), func(t *testing.T) {
			s := defaultBruteForceProtectionForTest(t.TempDir())
			s.SSHNotifyCooldown = sshNotifyCooldown
			if err := s.Validate(); err == nil {
				t.Errorf("BruteForceProtection.Validate() error = nil, wantErr not nil")
			}
		})
	}

	invalidSSHNotifyEvery := []int{-1, -1000}
	for _, sshNotifyEvery := range invalidSSHNotifyEvery {
		t.Run(fmt.Sprintf("invalid SSHNotifyEvery: %d", sshNotifyEvery), func(t *testing.T) {
			s := defaultBruteForceProtectionForTest(t.TempDir())
			s.SSHNotifyEvery = sshNotifyEvery
			if err := s.Validate(); err == nil {
				t.Errorf("BruteForceProtection.Validate() error = nil, wantErr not nil")
			}
		})
	}
}

func TestBruteForceProtection_ToSources_Success(t *testing.T) {
	t.Run("Enabled = true", func(t *testing.T) {
		i18nmock.Use(t)

		s := defaultBruteForceProtectionForTest(t.TempDir())
		if _, err := s.ToSources(); err != nil {
			t.Errorf("BruteForceProtection.ToSources() error = %v, wantErr nil", err)
		}
	})

	t.Run("Enabled = false", func(t *testing.T) {
		i18nmock.Use(t)

		s := defaultBruteForceProtectionForTest(t.TempDir())
		s.Enabled = false
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("BruteForceProtection.ToSources() error = %v, wantErr nil", err)
		} else if len(sources) != 0 {
			t.Errorf("BruteForceProtection.ToSources() sources = %v, want len(sources) = 0", len(sources))
		}
	})

	t.Run("No rules", func(t *testing.T) {
		i18nmock.Use(t)

		s := defaultBruteForceProtectionForTest(t.TempDir())
		s.Rules = nil
		if _, err := s.ToSources(); err != nil {
			t.Errorf("BruteForceProtection.ToSources() error = %v, wantErr nil", err)
		}
	})

	t.Run("SSHEnable = true", func(t *testing.T) {
		i18nmock.Use(t)

		s := defaultBruteForceProtectionForTest(t.TempDir())
		s.SSHEnable = true
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("BruteForceProtection.ToSources() error = %v, wantErr nil", err)
		} else if !testBruteForceProtectionSourceSSH(sources) {
			t.Errorf("BruteForceProtection.ToSources() not source SSH")
		}
	})

	t.Run("SSHEnable = false", func(t *testing.T) {
		i18nmock.Use(t)

		s := defaultBruteForceProtectionForTest(t.TempDir())
		s.SSHEnable = false
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("BruteForceProtection.ToSources() error = %v, wantErr nil", err)
		} else if testBruteForceProtectionSourceSSH(sources) {
			t.Errorf("BruteForceProtection.ToSources() source SSH")
		}
	})
}

func TestBruteForceProtection_ToSources_Error(t *testing.T) {
	t.Run("invalid Rule Group", func(t *testing.T) {
		i18nmock.Use(t)

		s := defaultBruteForceProtectionForTest(t.TempDir())
		s.Rules[0].Group = "invalid"
		if sources, err := s.ToSources(); err == nil {
			t.Errorf("BruteForceProtection.ToSources() error = nil, wantErr not nil")
		} else if len(sources) != 0 {
			t.Errorf("BruteForceProtection.ToSources() sources = %v, want empty", sources)
		}
	})

	t.Run("invalid SSH Group", func(t *testing.T) {
		i18nmock.Use(t)

		s := defaultBruteForceProtectionForTest(t.TempDir())
		s.SSHEnable = true
		s.SSHGroup = "invalid"
		if sources, err := s.ToSources(); err == nil {
			t.Errorf("BruteForceProtection.ToSources() error = nil, wantErr not nil")
		} else if testBruteForceProtectionSourceSSH(sources) {
			t.Errorf("BruteForceProtection.ToSources() source SSH")
		}
	})
}

func defaultBruteForceProtectionForTest(dir string) *BruteForceProtection {
	rule := defaultBruteForceProtectionRuleForTest(dir)

	return &BruteForceProtection{
		Enabled:              true,
		Notify:               fakers.RandBool(),
		RateLimitCount:       fakers.RandInt(1, 1000),
		RateLimitPeriod:      fakers.RandInt(1, 86400),
		RateLimitResetPeriod: fakers.RandInt(0, 86400),
		BlockingTime:         fakers.RandInt(0, 86400),
		SSHEnable:            fakers.RandBool(),
		SSHNotify:            fakers.RandBool(),
		SSHNotifyCooldown:    fakers.RandInt(0, 86400),
		SSHNotifyEvery:       fakers.RandInt(0, 10000),
		SSHGroup:             "",
		Groups:               []BruteForceProtectionGroup{},
		Rules:                []BruteForceProtectionRule{*rule},
	}
}

func testBruteForceProtectionSourceSSH(sources []*config.Source) bool {
	return testBruteForceProtectionSource(sources, "_ssh")
}

func testBruteForceProtectionSource(sources []*config.Source, name string) bool {
	if len(sources) == 0 {
		return false
	}
	for _, source := range sources {
		if source.BruteForceProtectionRule == nil {
			continue
		}
		if source.BruteForceProtectionRule.Name == name {
			return true
		}
	}
	return false
}
