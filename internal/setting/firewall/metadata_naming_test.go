package firewall

import (
	"fmt"
	"testing"
)

func TestMetadataNaming_Validate_Success(t *testing.T) {
	s := defaultMetadataNamingForTest()
	if err := s.Validate(); err != nil {
		t.Errorf("metadataNaming.Validate() error = %v, wantErr nil", err)
	}
}

func TestMetadataNaming_Validate_Error(t *testing.T) {
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

	t.Run("empty TableName", func(t *testing.T) {
		s := defaultMetadataNamingForTest()
		s.TableName = ""
		if err := s.Validate(); err == nil {
			t.Errorf("metadataNaming.Validate() error = nil, wantErr not nil")
		}
	})

	for _, name := range invalidNames {
		t.Run(fmt.Sprintf("invalid TableName: %s", name), func(t *testing.T) {
			s := defaultMetadataNamingForTest()
			s.TableName = name
			if err := s.Validate(); err == nil {
				t.Errorf("metadataNaming.Validate() error = nil, wantErr not nil")
			}
		})
	}

	t.Run("empty ChainInputName", func(t *testing.T) {
		s := defaultMetadataNamingForTest()
		s.ChainInputName = ""
		if err := s.Validate(); err == nil {
			t.Errorf("metadataNaming.Validate() error = nil, wantErr not nil")
		}
	})

	for _, name := range invalidNames {
		t.Run(fmt.Sprintf("invalid ChainInputName: %s", name), func(t *testing.T) {
			s := defaultMetadataNamingForTest()
			s.ChainInputName = name
			if err := s.Validate(); err == nil {
				t.Errorf("metadataNaming.Validate() error = nil, wantErr not nil")
			}
		})
	}

	t.Run("empty ChainOutputName", func(t *testing.T) {
		s := defaultMetadataNamingForTest()
		s.ChainOutputName = ""
		if err := s.Validate(); err == nil {
			t.Errorf("metadataNaming.Validate() error = nil, wantErr not nil")
		}
	})

	for _, name := range invalidNames {
		t.Run(fmt.Sprintf("invalid ChainOutputName: %s", name), func(t *testing.T) {
			s := defaultMetadataNamingForTest()
			s.ChainOutputName = name
			if err := s.Validate(); err == nil {
				t.Errorf("metadataNaming.Validate() error = nil, wantErr not nil")
			}
		})
	}

	t.Run("empty ChainForwardName", func(t *testing.T) {
		s := defaultMetadataNamingForTest()
		s.ChainForwardName = ""
		if err := s.Validate(); err == nil {
			t.Errorf("metadataNaming.Validate() error = nil, wantErr not nil")
		}
	})

	for _, name := range invalidNames {
		t.Run(fmt.Sprintf("invalid ChainForwardName: %s", name), func(t *testing.T) {
			s := defaultMetadataNamingForTest()
			s.ChainForwardName = name
			if err := s.Validate(); err == nil {
				t.Errorf("metadataNaming.Validate() error = nil, wantErr not nil")
			}
		})
	}
}

func defaultMetadataNamingForTest() *metadataNaming {
	return &metadataNaming{
		TableName:        "table_test",
		ChainInputName:   "chain_input_test",
		ChainOutputName:  "chain_output_test",
		ChainForwardName: "chain_forward_test",
	}
}
