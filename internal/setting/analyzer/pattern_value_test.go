package analyzer

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestPatternValue_ToPatternValue_Success(t *testing.T) {
	s := defaultPatternValueForTest()
	if patternValue, err := s.ToPatternValue(); err != nil {
		t.Errorf("PatternValue.ToPatternValue() error = %v, wantErr nil", err)
	} else if patternValue.Name != s.Name {
		t.Errorf("PatternValue.ToPatternValue() Name = %v, want %v", patternValue.Name, s.Name)
	} else if patternValue.Value != uint8(s.Value) {
		t.Errorf("PatternValue.ToPatternValue() Value = %v, want %v", patternValue.Value, uint8(s.Value))
	} else if string(patternValue.Type) != s.Type {
		t.Errorf("PatternValue.ToPatternValue() Type = %v, want %v", patternValue.Type, s.Type)
	}
}

func TestPatternValue_ToPatternValue_Error(t *testing.T) {
	t.Run("empty Name", func(t *testing.T) {
		s := defaultPatternValueForTest()
		s.Name = ""
		if _, err := s.ToPatternValue(); err == nil {
			t.Errorf("PatternValue.ToPatternValue() error = nil, wantErr not nil")
		}
	})

	t.Run("zero Value", func(t *testing.T) {
		s := defaultPatternValueForTest()
		s.Value = 0
		if _, err := s.ToPatternValue(); err == nil {
			t.Errorf("PatternValue.ToPatternValue() error = nil, wantErr not nil")
		}
	})

	t.Run("negative Value", func(t *testing.T) {
		s := defaultPatternValueForTest()
		s.Value = -1
		if _, err := s.ToPatternValue(); err == nil {
			t.Errorf("PatternValue.ToPatternValue() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid Type", func(t *testing.T) {
		s := defaultPatternValueForTest()
		s.Type = "invalid"
		if _, err := s.ToPatternValue(); err == nil {
			t.Errorf("PatternValue.ToPatternValue() error = nil, wantErr not nil")
		}
	})
}

func TestPatternValue_ToPatternValueForBruteForceProtection_Success(t *testing.T) {
	s := defaultPatternValueForTest()
	if patternValue, err := s.ToPatternValueForBruteForceProtection(); err != nil {
		t.Errorf("PatternValue.ToPatternValueForBruteForceProtection() error = %v, wantErr nil", err)
	} else if patternValue.Name != s.Name {
		t.Errorf("PatternValue.ToPatternValueForBruteForceProtection() Name = %v, want %v", patternValue.Name, s.Name)
	} else if patternValue.Value != uint8(s.Value) {
		t.Errorf("PatternValue.ToPatternValueForBruteForceProtection() Value = %v, want %v", patternValue.Value, uint8(s.Value))
	}
}

func TestPatternValue_ToPatternValueForBruteForceProtection_Error(t *testing.T) {
	t.Run("empty Name", func(t *testing.T) {
		s := defaultPatternValueForTest()
		s.Name = ""
		if _, err := s.ToPatternValueForBruteForceProtection(); err == nil {
			t.Errorf("PatternValue.ToPatternValueForBruteForceProtection() error = nil, wantErr not nil")
		}
	})

	t.Run("zero Value", func(t *testing.T) {
		s := defaultPatternValueForTest()
		s.Value = 0
		if _, err := s.ToPatternValueForBruteForceProtection(); err == nil {
			t.Errorf("PatternValue.ToPatternValueForBruteForceProtection() error = nil, wantErr not nil")
		}
	})

	t.Run("negative Value", func(t *testing.T) {
		s := defaultPatternValueForTest()
		s.Value = -1
		if _, err := s.ToPatternValueForBruteForceProtection(); err == nil {
			t.Errorf("PatternValue.ToPatternValueForBruteForceProtection() error = nil, wantErr not nil")
		}
	})
}

func defaultPatternValueForTest() *PatternValue {
	return &PatternValue{
		Name:  "test",
		Value: fakers.RandInt(1, 100),
		Type:  fakeRandPatternValueType(),
	}
}

func fakeRandPatternValueType() string {
	items := []string{"", "ip"}
	return fakers.RandItem(items)
}
