package analyzer

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestBruteForceProtectionGroupRateLimit_ToRateLimit_Success(t *testing.T) {
	s := defaultBruteForceProtectionGroupRateLimitForTest()
	block := newBlockMock(false, nil)
	if _, err := s.ToRateLimit(block); err != nil {
		t.Errorf("BruteForceProtectionGroupRateLimit.ToRateLimit() error = %v", err)
	}
}

func TestBruteForceProtectionGroupRateLimit_ToRateLimit_Error(t *testing.T) {
	invalidCount := []int{0, -1, -100}
	for _, count := range invalidCount {
		t.Run(fmt.Sprintf("invalid Count: %d", count), func(t *testing.T) {
			s := defaultBruteForceProtectionGroupRateLimitForTest()
			s.Count = count
			block := newBlockMock(false, nil)
			if _, err := s.ToRateLimit(block); err == nil {
				t.Errorf("BruteForceProtectionGroupRateLimit.ToRateLimit() error = nil, wantErr nil")
			}
		})
	}

	invalidPeriod := []int{0, -1, -100}
	for _, period := range invalidPeriod {
		t.Run(fmt.Sprintf("invalid Period: %d", period), func(t *testing.T) {
			s := defaultBruteForceProtectionGroupRateLimitForTest()
			s.Period = period
			block := newBlockMock(false, nil)
			if _, err := s.ToRateLimit(block); err == nil {
				t.Errorf("BruteForceProtectionGroupRateLimit.ToRateLimit() error = nil, wantErr nil")
			}
		})
	}

	invalidBlockingTime := []int{-1, -100}
	for _, blockingTime := range invalidBlockingTime {
		t.Run(fmt.Sprintf("invalid BlockingTime: %d", blockingTime), func(t *testing.T) {
			s := defaultBruteForceProtectionGroupRateLimitForTest()
			s.BlockingTime = blockingTime
			block := newBlockMock(false, nil)
			if _, err := s.ToRateLimit(block); err == nil {
				t.Errorf("BruteForceProtectionGroupRateLimit.ToRateLimit() error = nil, wantErr nil")
			}
		})
	}

	t.Run("invalid BlockType", func(t *testing.T) {
		s := defaultBruteForceProtectionGroupRateLimitForTest()
		s.BlockType = "invalid"
		block := newBlockMock(false, nil)
		if _, err := s.ToRateLimit(block); err == nil {
			t.Errorf("BruteForceProtectionGroupRateLimit.ToRateLimit() error = nil, wantErr nil")
		}
	})

	invalidPorts := []string{"", "invalid", "tcp", "80/", "/tcp", "80 tcp", "80 / tcp"}
	for _, port := range invalidPorts {
		t.Run(fmt.Sprintf("invalid Port: %s", port), func(t *testing.T) {
			s := defaultBruteForceProtectionGroupRateLimitForTest()
			s.BlockType = "ip_port"
			s.Ports = []string{port}
			block := newBlockMock(false, nil)
			if _, err := s.ToRateLimit(block); err == nil {
				t.Errorf("BruteForceProtectionGroupRateLimit.ToRateLimit() error = nil, wantErr nil")
			}
		})
	}
}

func defaultBruteForceProtectionGroupRateLimitForTest() *BruteForceProtectionGroupRateLimit {
	blockType, ports := fakeRandBruteForceProtectionBlockTypeAndPorts()

	return &BruteForceProtectionGroupRateLimit{
		Count:        fakers.RandInt(1, 10000),
		Period:       fakers.RandInt(1, 86400),
		BlockingTime: fakers.RandInt(0, 86400),
		BlockType:    blockType,
		Ports:        ports,
	}
}

type blockMock struct {
	shouldPortsBlocked bool
	ports              []types.L4Port
	calls              int
}

func (m *blockMock) PortsBlocked() (bool, []types.L4Port) {
	m.calls++

	if !m.shouldPortsBlocked {
		return false, nil
	}

	// a protective copy to prevent tests from catching side effects
	portsCopy := append([]types.L4Port(nil), m.ports...)
	return true, portsCopy
}

func newBlockMock(shouldPortsBlocked bool, ports []types.L4Port) *blockMock {
	return &blockMock{
		shouldPortsBlocked: shouldPortsBlocked,
		ports:              append([]types.L4Port(nil), ports...),
	}
}
