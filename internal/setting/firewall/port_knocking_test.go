package firewall

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/ip"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestPortKnocking_ToPortKnocking_Success(t *testing.T) {
	s := defaultPortKnockingForTest()
	if knocking, err := s.ToPortKnocking(); err != nil {
		t.Errorf("portKnocking.ToPortKnocking() error = %v, wantErr nil", err)
	} else if knocking.Name != s.Name {
		t.Errorf("portKnocking.ToPortKnocking() name = %v, want %v", knocking.Name, s.Name)
	} else if knocking.Port.Number() != uint16(s.Port) {
		t.Errorf("portKnocking.ToPortKnocking() port = %v, want %v", knocking.Port.Number(), uint16(s.Port))
	} else if knocking.Port.ProtocolString() != s.Protocol {
		t.Errorf("portKnocking.ToPortKnocking() protocol = %v, want %v", knocking.Port.ProtocolString(), s.Protocol)
	} else if !testPortKnockingIPVersion(knocking.IPVersion, s.IPVersion) {
		t.Errorf("portKnocking.ToPortKnocking() ipVersion = %v, want %v", knocking.IPVersion, s.IPVersion)
	} else if knocking.Knocks == nil {
		t.Errorf("portKnocking.ToPortKnocking() knocks = %v, want not nil", knocking.Knocks)
	}

}

func TestPortKnocking_ToPortKnocking_Error(t *testing.T) {
	invalidNames := []string{
		"",
		" ",
		"name.with.dot",
		"name/with/slash",
		"name with space",
		"name@with@at",
		"привет",                            // Cyrillic
		"abcdefghijklmnopqrstuvwxyzABCDEFG", // 33 symbols
	}

	for _, name := range invalidNames {
		t.Run(fmt.Sprintf("invalid Name: %s", name), func(t *testing.T) {
			s := defaultPortKnockingForTest()
			s.Name = name
			if _, err := s.ToPortKnocking(); err == nil {
				t.Errorf("portKnocking.ToPortKnocking() error = nil, wantErr not nil")
			}
		})
	}

	t.Run("empty IPVersion", func(t *testing.T) {
		s := defaultPortKnockingForTest()
		s.IPVersion = ""
		if _, err := s.ToPortKnocking(); err == nil {
			t.Errorf("portKnocking.ToPortKnocking() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid IPVersion", func(t *testing.T) {
		s := defaultPortKnockingForTest()
		s.IPVersion = "invalid"
		if _, err := s.ToPortKnocking(); err == nil {
			t.Errorf("portKnocking.ToPortKnocking() error = nil, wantErr not nil")
		}
	})

	t.Run("negative Port", func(t *testing.T) {
		s := defaultPortKnockingForTest()
		s.Port = -1
		if _, err := s.ToPortKnocking(); err == nil {
			t.Errorf("portKnocking.ToPortKnocking() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Protocol", func(t *testing.T) {
		s := defaultPortKnockingForTest()
		s.Protocol = ""
		if _, err := s.ToPortKnocking(); err == nil {
			t.Errorf("portKnocking.ToPortKnocking() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Protocol", func(t *testing.T) {
		s := defaultPortKnockingForTest()
		s.Protocol = "invalid"
		if _, err := s.ToPortKnocking(); err == nil {
			t.Errorf("portKnocking.ToPortKnocking() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Knocks", func(t *testing.T) {
		s := defaultPortKnockingForTest()
		s.Knocks = nil
		if _, err := s.ToPortKnocking(); err == nil {
			t.Errorf("portKnocking.ToPortKnocking() error = nil, wantErr not nil")
		}
	})
}

func defaultPortKnockingForTest() *portKnocking {
	timeout := fakers.RandInt(1, 2147483647)
	knock := portKnockingKnock{
		Port:     fakers.RandInt(0, 65535),
		Protocol: fakers.RandFirewallProtocol(),
		Timeout:  int32(timeout),
		Action:   fakeRandPortKnockingKnockAction(),
	}

	return &portKnocking{
		Name:      "test",
		IPVersion: fakeRandPortKnockingIPVersion(),
		Port:      fakers.RandInt(0, 65535),
		Protocol:  fakers.RandFirewallProtocol(),
		Knocks:    []portKnockingKnock{knock},
	}
}

func fakeRandPortKnockingIPVersion() string {
	versions := []string{"ip4", "ip6"}

	return fakers.RandItem(versions)
}

func testPortKnockingIPVersion(ipVersion ip.Version, configIPVersion string) bool {
	if ipVersion == ip.IPv4 {
		return configIPVersion == "ip4"
	}

	return configIPVersion == "ip6"
}
