package analyzer

import (
	"fmt"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestPatternPartition_ToPatternPartition_Success(t *testing.T) {
	s := defaultPatternPartitionForTest()
	if patternPartition, err := s.ToPatternPartition(); err != nil {
		t.Errorf("PatternPartition.ToPatternPartition() error = %v, wantErr nil", err)
	} else if patternPartition.Value != uint8(s.Value) {
		t.Errorf("PatternPartition.ToPatternPartition() = %v, want %v", patternPartition.Value, uint8(s.Value))
	} else if keyword := "test"; patternPartition.Normalize(keyword) != keyword {
		t.Errorf("PatternPartition.Normalize() = %v, want %v", patternPartition.Normalize(keyword), keyword)
	}
}

func TestPatternPartition_ToPatternPartition_Error(t *testing.T) {
	invalidValue := []int{0, -1, -500, 256, 500}
	for _, value := range invalidValue {
		t.Run(fmt.Sprintf("invalid Value=%d", value), func(t *testing.T) {
			s := defaultPatternPartitionForTest()
			s.Value = value
			if _, err := s.ToPatternPartition(); err == nil {
				t.Errorf("PatternPartition.ToPatternPartition() error = nil, wantErr not nil")
			}
		})
	}
}

func defaultPatternPartitionForTest() *PatternPartition {
	return &PatternPartition{
		Value:     fakers.RandInt(1, 255),
		Trim:      fakers.RandBool(),
		LowerCase: fakers.RandBool(),
		Type:      nil,
	}
}
