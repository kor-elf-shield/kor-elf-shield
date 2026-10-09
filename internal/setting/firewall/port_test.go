package firewall

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestPort_ToPorts_Success(t *testing.T) {
	s := defaultPortForTest()
	if _, _, err := s.ToPorts(); err != nil {
		t.Errorf("Port.ToPorts() error = %v", err)
	}
}

func TestPort_ToPorts_Error(t *testing.T) {
	t.Run("empty Number", func(t *testing.T) {
		s := defaultPortForTest()
		s.Numbers = nil
		if _, _, err := s.ToPorts(); err == nil {
			t.Errorf("Port.ToPorts() error = nil, wantErr not nil")
		}
	})

	t.Run("negative Number", func(t *testing.T) {
		s := defaultPortForTest()
		s.Numbers[0] = -1
		if _, _, err := s.ToPorts(); err == nil {
			t.Errorf("Port.ToPorts() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Direction", func(t *testing.T) {
		s := defaultPortForTest()
		s.Directions = nil
		if _, _, err := s.ToPorts(); err == nil {
			t.Errorf("Port.ToPorts() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Direction", func(t *testing.T) {
		s := defaultPortForTest()
		s.Directions[0] = "invalid"
		if _, _, err := s.ToPorts(); err == nil {
			t.Errorf("Port.ToPorts() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Protocol", func(t *testing.T) {
		s := defaultPortForTest()
		s.Protocols = nil
		if _, _, err := s.ToPorts(); err == nil {
			t.Errorf("Port.ToPorts() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Protocol", func(t *testing.T) {
		s := defaultPortForTest()
		s.Protocols[0] = "invalid"
		if _, _, err := s.ToPorts(); err == nil {
			t.Errorf("Port.ToPorts() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Action", func(t *testing.T) {
		s := defaultPortForTest()
		s.Action = ""
		if _, _, err := s.ToPorts(); err == nil {
			t.Errorf("Port.ToPorts() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Action", func(t *testing.T) {
		s := defaultPortForTest()
		s.Action = "invalid"
		if _, _, err := s.ToPorts(); err == nil {
			t.Errorf("Port.ToPorts() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid LimitRate", func(t *testing.T) {
		s := defaultPortForTest()
		s.LimitRate = "one/s"
		if _, _, err := s.ToPorts(); err == nil {
			t.Errorf("Port.ToPorts() error = nil, wantErr not nil")
		}
	})
}

func defaultPortForTest() *Port {
	return &Port{
		Numbers:    fakers.RandFirewallPorts(),
		Directions: fakers.RandFirewallDirections(),
		Protocols:  fakers.RandFirewallProtocols(),
		Action:     fakers.RandFirewallAction(),
		LimitRate:  fakers.RandFirewallLimitRate(),
	}
}
