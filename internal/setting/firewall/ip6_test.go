package firewall

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestIp6_Validate(t *testing.T) {
	s := defaultIp6ForTest()
	if err := s.Validate(); err != nil {
		t.Errorf("ip6.Validate() error = %v, wantErr nil", err)
	}
}

func defaultIp6ForTest() *ip6 {
	return &ip6{
		Enable:     true,
		IcmpStrict: fakers.RandBool(),
	}
}
