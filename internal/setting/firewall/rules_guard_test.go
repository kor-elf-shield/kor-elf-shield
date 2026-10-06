package firewall

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestRulesGuard_Validate_Success(t *testing.T) {
	s := defaultRulesGuardForTest()

	if err := s.Validate(); err != nil {
		t.Errorf("RulesGuard.Validate() error = %v, wantErr nil", err)
	}
}

func TestRulesGuard_Validate_Error(t *testing.T) {
	t.Run("interval is less than 60", func(t *testing.T) {
		s := defaultRulesGuardForTest()
		interval := fakers.RandInt(0, 59)
		s.Interval = int32(interval)

		if err := s.Validate(); err == nil {
			t.Errorf("RulesGuard.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("interval is greater than 2147483647", func(t *testing.T) {
		s := defaultRulesGuardForTest()
		interval := fakers.RandInt(2147483648, 4294967295)
		s.Interval = int32(interval)

		if err := s.Validate(); err == nil {
			t.Errorf("RulesGuard.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("negative interval", func(t *testing.T) {
		s := defaultRulesGuardForTest()
		s.Interval = -1

		if err := s.Validate(); err == nil {
			t.Errorf("RulesGuard.Validate() error = nil, wantErr not nil")
		}
	})
}

func TestRulesGuard_ToGuardConfig_Success(t *testing.T) {
	s := defaultRulesGuardForTest()
	config := s.ToGuardConfig()

	if config.Enable != s.Enabled {
		t.Errorf("RulesGuard.ToGuardConfig().Enable = %v, want %v", config.Enable, s.Enabled)
	}

	if config.Notifications != s.Notifications {
		t.Errorf("RulesGuard.ToGuardConfig().Notifications = %v, want %v", config.Notifications, s.Notifications)
	}

	if config.Recovery != s.Recovery {
		t.Errorf("RulesGuard.ToGuardConfig().Recovery = %v, want %v", config.Recovery, s.Recovery)
	}

	if config.Interval != uint32(s.Interval) {
		t.Errorf("RulesGuard.ToGuardConfig().Interval = %v, want %v", config.Interval, uint32(s.Interval))
	}
}

func defaultRulesGuardForTest() *RulesGuard {
	interval := fakers.RandInt(60, 2147483647)

	return &RulesGuard{
		Enabled:       true,
		Notifications: fakers.RandBool(),
		Recovery:      fakers.RandBool(),
		Interval:      int32(interval),
	}
}
