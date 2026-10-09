package firewall

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestIp4_Validate_Success(t *testing.T) {
	s := defaultIp4ForTest()
	if err := s.Validate(); err != nil {
		t.Errorf("ip4.Validate() error = %v, wantErr nil", err)
	}
}

func TestIp4_Validate_Error(t *testing.T) {
	t.Run("empty IcmpInRate", func(t *testing.T) {
		s := defaultIp4ForTest()
		s.IcmpInRate = ""
		if err := s.Validate(); err == nil {
			t.Errorf("ip4.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid IcmpInRate", func(t *testing.T) {
		s := defaultIp4ForTest()
		s.IcmpInRate = "invalid"
		if err := s.Validate(); err == nil {
			t.Errorf("ip4.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("empty IcmpOutRate", func(t *testing.T) {
		s := defaultIp4ForTest()
		s.IcmpOutRate = ""
		if err := s.Validate(); err == nil {
			t.Errorf("ip4.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid IcmpOutRate", func(t *testing.T) {
		s := defaultIp4ForTest()
		s.IcmpOutRate = "invalid"
		if err := s.Validate(); err == nil {
			t.Errorf("ip4.Validate() error = nil, wantErr not nil")
		}
	})
}

func defaultIp4ForTest() *ip4 {
	return &ip4{
		IcmpIn:            fakers.RandBool(),
		IcmpInRate:        fakeRandLimitRate(),
		IcmpOut:           fakers.RandBool(),
		IcmpOutRate:       fakeRandLimitRate(),
		IcmpTimestampDrop: fakers.RandBool(),
	}
}

func fakeRandLimitRate() string {
	rates := []string{
		"0",
		"1/s",
		"10/minute",
		"100/hour",
		"1000/day",
		"5/week",
		"0.5/second",
		"over 1/s",
		"10 packets/second",
		"1 kbytes/min",
		"2 mbytes/h",
		"1/s burst 5",
		"1/s burst 10 packets",
		"2/min burst 1.5 mbytes",
	}
	return fakers.RandItem(rates)
}
