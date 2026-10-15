package analyzer

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestBruteForceProtectionPattern_ToPattern_Success(t *testing.T) {
	s := defaultBruteForceProtectionPatternForTest()
	if _, err := s.ToPattern(); err != nil {
		t.Errorf("BruteForceProtectionPattern.ToPattern() error = %v, wantErr nil", err)
	}
}

func TestBruteForceProtectionPattern_ToPattern_Error(t *testing.T) {
	invalidRegexp := []string{
		"",
		"[",
		"[s-1]",
	}
	for _, regexp := range invalidRegexp {
		t.Run(fmt.Sprintf("invalid Regexp: %s", regexp), func(t *testing.T) {
			s := defaultBruteForceProtectionPatternForTest()
			s.Regexp = regexp
			if alertRegexp, err := s.ToPattern(); err == nil {
				if alertRegexp.Regexp == nil {
					t.Errorf("BruteForceProtectionPattern.ToPattern() error = nil, wantErr nil")
				} else if _, err := alertRegexp.Regexp.Get(); err == nil {
					t.Errorf("regexp.Compile error = nil, wantErr not nil")
				}
			}
		})
	}

	invalidIP := []int{0, -1, -255, -256, 500, 256, 500}
	for _, ip := range invalidIP {
		t.Run(fmt.Sprintf("invalid IP: %d", ip), func(t *testing.T) {
			s := defaultBruteForceProtectionPatternForTest()
			s.IP = ip
			if _, err := s.ToPattern(); err == nil {
				t.Errorf("BruteForceProtectionPattern.ToPattern() error = nil, wantErr not nil")
			}
		})
	}
}

func defaultBruteForceProtectionPatternForTest() *BruteForceProtectionPattern {
	return &BruteForceProtectionPattern{
		Regexp:    "test",
		IP:        fakers.RandInt(1, 255),
		Values:    nil,
		Partition: nil,
	}
}
