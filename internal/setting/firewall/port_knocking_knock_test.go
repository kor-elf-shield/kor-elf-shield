package firewall

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestPortKnockingKnock_ToKnock_Success(t *testing.T) {
	s := defaultPortKnockingKnockForTest()
	if knock, err := s.ToKnock(); err != nil {
		t.Errorf("portKnockingKnock.ToKnock() error = %v, wantErr nil", err)
	} else if knock.Timeout != uint32(s.Timeout) {
		t.Errorf("portKnockingKnock.ToKnock() timeout = %v, want %v", knock.Timeout, uint32(s.Timeout))
	} else if knock.Action.String() != s.Action {
		t.Errorf("portKnockingKnock.ToKnock() action = %v, want %v", knock.Action, s.Action)
	} else if knock.Port.Number() != uint16(s.Port) {
		t.Errorf("portKnockingKnock.ToKnock() port = %v, want %v", knock.Port.Number(), uint16(s.Port))
	} else if knock.Port.ProtocolString() != s.Protocol {
		t.Errorf("portKnockingKnock.ToKnock() protocol = %v, want %v", knock.Port.ProtocolString(), s.Protocol)
	}
}

func TestPortKnockingKnock_ToKnock_Error(t *testing.T) {
	t.Run("negative Port", func(t *testing.T) {
		s := defaultPortKnockingKnockForTest()
		s.Port = -1
		if _, err := s.ToKnock(); err == nil {
			t.Errorf("portKnockingKnock.ToKnock() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Timeout", func(t *testing.T) {
		s := defaultPortKnockingKnockForTest()
		s.Timeout = 0
		if _, err := s.ToKnock(); err == nil {
			t.Errorf("portKnockingKnock.ToKnock() error = nil, wantErr not nil")
		}
	})

	t.Run("negative Timeout", func(t *testing.T) {
		s := defaultPortKnockingKnockForTest()
		s.Timeout = -1
		if _, err := s.ToKnock(); err == nil {
			t.Errorf("portKnockingKnock.ToKnock() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Protocol", func(t *testing.T) {
		s := defaultPortKnockingKnockForTest()
		s.Protocol = ""
		if _, err := s.ToKnock(); err == nil {
			t.Errorf("portKnockingKnock.ToKnock() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Protocol", func(t *testing.T) {
		s := defaultPortKnockingKnockForTest()
		s.Protocol = "invalid"
		if _, err := s.ToKnock(); err == nil {
			t.Errorf("portKnockingKnock.ToKnock() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Action", func(t *testing.T) {
		s := defaultPortKnockingKnockForTest()
		s.Action = ""
		if _, err := s.ToKnock(); err == nil {
			t.Errorf("portKnockingKnock.ToKnock() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Action", func(t *testing.T) {
		s := defaultPortKnockingKnockForTest()
		s.Action = "invalid"
		if _, err := s.ToKnock(); err == nil {
			t.Errorf("portKnockingKnock.ToKnock() error = nil, wantErr not nil")
		}
	})
}

func defaultPortKnockingKnockForTest() *portKnockingKnock {
	timeout := fakers.RandInt(1, 2147483647)

	return &portKnockingKnock{
		Port:     fakers.RandInt(0, 65535),
		Protocol: fakers.RandFirewallProtocol(),
		Timeout:  int32(timeout),
		Action:   fakeRandPortKnockingKnockAction(),
	}
}

func fakeRandPortKnockingKnockAction() string {
	actions := []string{"accept", "return", "drop", "reject"}
	return actions[fakers.RandInt(0, len(actions)-1)]
}
