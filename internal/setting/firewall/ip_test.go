package firewall

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestIP_ToIPs_Success(t *testing.T) {
	s := defaultIPForTest()
	if _, err := s.ToIPs(); err != nil {
		t.Errorf("IP.ToIPs() error = %v, wantErr nil", err)
	}
}

func TestIP_ToIPs_Error(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		s := &IP{}
		if _, err := s.ToIPs(); err == nil {
			t.Errorf("IP.ToIPs() error = nil, wantErr not nil")
		}
	})

	invalidIPs := []string{
		"",
		" ",
		"invalid",
		"127.0.0.1/33",
		"500.0.0.1",
	}

	for _, ip := range invalidIPs {
		t.Run(fmt.Sprintf("invalid IPS: %s", ip), func(t *testing.T) {
			s := defaultIPForTest()
			s.IPs = []string{ip}
			if _, err := s.ToIPs(); err == nil {
				t.Errorf("IP.ToIPs() error = nil, wantErr not nil")
			}
		})
	}

	t.Run("empty Action", func(t *testing.T) {
		s := defaultIPForTest()
		s.Action = ""
		if _, err := s.ToIPs(); err == nil {
			t.Errorf("IP.ToIPs() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Action", func(t *testing.T) {
		s := defaultIPForTest()
		s.Action = "invalid"
		if _, err := s.ToIPs(); err == nil {
			t.Errorf("IP.ToIPs() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Directions", func(t *testing.T) {
		s := defaultIPForTest()
		s.Directions = []string{}
		if _, err := s.ToIPs(); err == nil {
			t.Errorf("IP.ToIPs() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Directions", func(t *testing.T) {
		s := defaultIPForTest()
		s.Directions = []string{"invalid"}
		if _, err := s.ToIPs(); err == nil {
			t.Errorf("IP.ToIPs() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Protocols", func(t *testing.T) {
		s := defaultIPForTest()
		s.Protocols = []string{"invalid"}
		if _, err := s.ToIPs(); err == nil {
			t.Errorf("IP.ToIPs() error = nil, wantErr not nil")
		}
	})

	invalidPorts := []int{
		-1, -1000,
		65536, 70000,
	}
	for _, port := range invalidPorts {
		t.Run(fmt.Sprintf("invalid Port %d", port), func(t *testing.T) {
			s := defaultIPForTest()
			s.Ports = []int{port}
			if _, err := s.ToIPs(); err == nil {
				t.Errorf("IP.ToIPs() error = nil, wantErr not nil")
			}
		})
	}

	t.Run("invalid LimitRate", func(t *testing.T) {
		s := defaultIPForTest()
		s.LimitRate = "invalid"
		if _, err := s.ToIPs(); err == nil {
			t.Errorf("IP.ToIPs() error = nil, wantErr not nil")
		}
	})
}

func defaultIPForTest() *IP {
	return &IP{
		IPs:        fakeIPsForTest(),
		Action:     fakers.RandFirewallAction(),
		Directions: fakers.RandFirewallDirections(),
		Protocols:  fakers.RandFirewallProtocols(),
		Ports:      fakers.RandFirewallPorts(),
		LimitRate:  fakers.RandFirewallLimitRate(),
	}
}

func fakeIPsForTest() []string {
	items := [][]string{
		{"127.0.0.1", "fe80::260:8ff:fe52:f9d8"},
		{"192.168.0.1", "192.168.1.0/24"},
		{"10.0.0.1"},
		{"172.16.0.1"},
		{"192.168.1.1"},
	}

	return fakers.RandItem(items)
}
