package analyzer

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestBruteForceProtectionGroup_ToGroup_Success(t *testing.T) {
	s := defaultBruteForceProtectionGroupForTest()
	if group, err := s.ToGroup(); err != nil {
		t.Errorf("BruteForceProtectionGroup.ToGroup() error = %v", err)
	} else if group == nil {
		t.Errorf("BruteForceProtectionGroup.ToGroup() group = nil, error = nil, wantErr nil")
	} else if group.Name != s.Name {
		t.Errorf("BruteForceProtectionGroup.ToGroup() group.Name = %v, want %v", group.Name, s.Name)
	} else if group.Message != s.Message {
		t.Errorf("BruteForceProtectionGroup.ToGroup() group.Message = %v, want %v", group.Message, s.Message)
	} else if group.RateLimitResetPeriod != uint32(s.RateLimitResetPeriod) {
		t.Errorf("BruteForceProtectionGroup.ToGroup() group.RateLimitResetPeriod = %v, want %v", group.RateLimitResetPeriod, uint32(s.RateLimitResetPeriod))
	}
}

func TestBruteForceProtectionGroup_ToGroup_Error(t *testing.T) {
	t.Run("negative RateLimitResetPeriod", func(t *testing.T) {
		s := defaultBruteForceProtectionGroupForTest()
		s.RateLimitResetPeriod = -1
		if _, err := s.ToGroup(); err == nil {
			t.Errorf("BruteForceProtectionGroup.ToGroup() error = nil, wantErr not nil")
		}
	})

	t.Run("empty RateLimits", func(t *testing.T) {
		s := defaultBruteForceProtectionGroupForTest()
		s.RateLimits = nil
		if _, err := s.ToGroup(); err == nil {
			t.Errorf("BruteForceProtectionGroup.ToGroup() error = nil, wantErr not nil")
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
			s := defaultBruteForceProtectionGroupForTest()
			s.Name = name

			if _, err := s.ToGroup(); err == nil {
				t.Errorf("BruteForceProtectionGroup.ToGroup() error = nil, wantErr not nil")
			}
		})
	}

	t.Run("invalid BlockType", func(t *testing.T) {
		s := defaultBruteForceProtectionGroupForTest()
		s.BlockType = "invalid"

		if _, err := s.ToGroup(); err == nil {
			t.Errorf("BruteForceProtectionGroup.ToGroup() error = nil, wantErr not nil")
		}
	})

	invalidPorts := []string{"", "invalid", "tcp", "80/", "/tcp", "80 tcp", "80 / tcp"}
	for _, port := range invalidPorts {
		t.Run(fmt.Sprintf("invalid Port: %s", port), func(t *testing.T) {
			s := defaultBruteForceProtectionGroupForTest()
			s.BlockType = "ip_port"
			s.Ports = []string{port}
			if _, err := s.ToGroup(); err == nil {
				t.Errorf("BruteForceProtectionGroup.ToGroup() error = nil, wantErr nil")
			}
		})
	}

	t.Run("empty RateLimits", func(t *testing.T) {
		s := defaultBruteForceProtectionGroupForTest()
		s.RateLimits = nil
		if _, err := s.ToGroup(); err == nil {
			t.Errorf("BruteForceProtectionGroup.ToGroup() error = nil, wantErr nil")
		}
	})
}

func defaultBruteForceProtectionGroupForTest() *BruteForceProtectionGroup {
	blockType, ports := fakeRandBruteForceProtectionBlockTypeAndPorts()
	rateLimits := fakeRandBruteForceProtectionGroupRateLimits()

	return &BruteForceProtectionGroup{
		Name:                 "test",
		Message:              "test",
		RateLimitResetPeriod: fakers.RandInt(0, 1000),
		BlockType:            blockType,
		Ports:                ports,
		RateLimits:           rateLimits,
	}
}

func fakeRandBruteForceProtectionGroupRateLimits() []BruteForceProtectionGroupRateLimit {
	var rateLimits []BruteForceProtectionGroupRateLimit

	for i := fakers.RandInt(1, 5); i > 0; i-- {
		rateLimits = append(rateLimits, *defaultBruteForceProtectionGroupRateLimitForTest())
	}

	return rateLimits
}
