package analyzer

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestPatternPartitionType_ToType_Success(t *testing.T) {
	s := defaultPatternPartitionTypeForTest()
	normalize := newNormalizeMock()
	patternPartitionType, err := s.ToType(normalize)
	if err != nil {
		t.Errorf("PatternPartitionType.ToType() error = %v, wantErr nil", err)
	}

	if patternPartitionType.IsPartitioned() != s.Partitioned {
		t.Errorf("PatternPartitionType.IsPartitioned() = %v, want %v", patternPartitionType.IsPartitioned(), s.Partitioned)
	}

	switch s.Type {
	case "only":
		keywords := s.Keywords[0]
		if isMatch, _ := patternPartitionType.Accepts(keywords); !isMatch {
			t.Errorf("PatternPartitionType.Accepts(%s) = %v, want true", keywords, isMatch)
		}
	case "except":
		keywords := s.Keywords[0] + "test"
		if isMatch, _ := patternPartitionType.Accepts(keywords); !isMatch {
			t.Errorf("PatternPartitionType.Accepts(%s) = %v, want true", keywords, isMatch)
		}
	}
}

func TestPatternPartitionType_ToType_Error(t *testing.T) {
	t.Run("empty Type", func(t *testing.T) {
		s := defaultPatternPartitionTypeForTest()
		s.Type = ""
		normalize := newNormalizeMock()
		if _, err := s.ToType(normalize); err == nil {
			t.Errorf("PatternPartitionType.ToType() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Type", func(t *testing.T) {
		s := defaultPatternPartitionTypeForTest()
		s.Type = "invalid"
		normalize := newNormalizeMock()
		if _, err := s.ToType(normalize); err == nil {
			t.Errorf("PatternPartitionType.ToType() error = nil, wantErr not nil")
		}
	})

	t.Run("empty Keywords", func(t *testing.T) {
		s := defaultPatternPartitionTypeForTest()
		s.Keywords = nil
		normalize := newNormalizeMock()
		if _, err := s.ToType(normalize); err == nil {
			t.Errorf("PatternPartitionType.ToType() error = nil, wantErr not nil")
		}
	})
}

func defaultPatternPartitionTypeForTest() *PatternPartitionType {
	return &PatternPartitionType{
		Type:        fakeRandPatternPartitionTypeType(),
		Keywords:    []string{"test"},
		Partitioned: fakers.RandBool(),
	}
}

func fakeRandPatternPartitionTypeType() string {
	items := []string{"only", "except"}
	return fakers.RandItem(items)
}

type normalizeMock struct {
}

func (m *normalizeMock) Normalize(s string) string {
	return s
}

func newNormalizeMock() *normalizeMock {
	return &normalizeMock{}
}
