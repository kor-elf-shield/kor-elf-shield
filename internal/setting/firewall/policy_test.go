package firewall

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func TestPolicy_ToConfigPolicy_Success(t *testing.T) {
	s := defaultPolicyForTest()
	if config, err := s.ToConfigPolicy(); err != nil {
		t.Errorf("policy.ToConfigPolicy() error = %v", err)
	} else if config.DefaultAllowInput != s.DefaultAllowInput {
		t.Errorf("policy.ToConfigPolicy() default allow input = %v, want %v", config.DefaultAllowInput, s.DefaultAllowInput)
	} else if config.DefaultAllowOutput != s.DefaultAllowOutput {
		t.Errorf("policy.ToConfigPolicy() default allow output = %v, want %v", config.DefaultAllowOutput, s.DefaultAllowOutput)
	} else if config.DefaultAllowForward != s.DefaultAllowForward {
		t.Errorf("policy.ToConfigPolicy() default allow forward = %v, want %v", config.DefaultAllowForward, s.DefaultAllowForward)
	} else if config.InputPriority != s.InputPriority {
		t.Errorf("policy.ToConfigPolicy() input priority = %v, want %v", config.InputPriority, s.InputPriority)
	} else if config.OutputPriority != s.OutputPriority {
		t.Errorf("policy.ToConfigPolicy() output priority = %v, want %v", config.OutputPriority, s.OutputPriority)
	} else if config.ForwardPriority != s.ForwardPriority {
		t.Errorf("policy.ToConfigPolicy() forward priority = %v, want %v", config.ForwardPriority, s.ForwardPriority)
	} else if config.InputDrop.String() != s.InputDrop {
		t.Errorf("policy.ToConfigPolicy() input drop = %v, want %v", config.InputDrop.String(), s.InputDrop)
	} else if config.OutputDrop.String() != s.OutputDrop {
		t.Errorf("policy.ToConfigPolicy() output drop = %v, want %v", config.OutputDrop.String(), s.OutputDrop)
	} else if config.ForwardDrop.String() != s.ForwardDrop {
		t.Errorf("policy.ToConfigPolicy() forward drop = %v, want %v", config.ForwardDrop.String(), s.ForwardDrop)
	}
}

func TestPolicy_ToConfigPolicy_Error(t *testing.T) {
	t.Run("empty InputDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.InputDrop = ""
		if _, err := s.ToConfigPolicy(); err == nil {
			t.Errorf("policy.ToConfigPolicy() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid InputDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.InputDrop = "invalid"
		if _, err := s.ToConfigPolicy(); err == nil {
			t.Errorf("policy.ToConfigPolicy() error = nil, wantErr not nil")
		}
	})

	t.Run("empty OutputDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.OutputDrop = ""
		if _, err := s.ToConfigPolicy(); err == nil {
			t.Errorf("policy.ToConfigPolicy() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid OutputDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.OutputDrop = "invalid"
		if _, err := s.ToConfigPolicy(); err == nil {
			t.Errorf("policy.ToConfigPolicy() error = nil, wantErr not nil")
		}
	})

	t.Run("empty ForwardDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.ForwardDrop = ""
		if _, err := s.ToConfigPolicy(); err == nil {
			t.Errorf("policy.ToConfigPolicy() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid ForwardDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.ForwardDrop = "invalid"
		if _, err := s.ToConfigPolicy(); err == nil {
			t.Errorf("policy.ToConfigPolicy() error = nil, wantErr not nil")
		}
	})
}

func TestPolicy_Validate_Success(t *testing.T) {
	s := defaultPolicyForTest()
	if err := s.Validate(); err != nil {
		t.Errorf("policy.Validate() error = %v", err)
	}
}

func TestPolicy_Validate_Error(t *testing.T) {
	t.Run("empty InputDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.InputDrop = ""
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid InputDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.InputDrop = "invalid"
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("empty OutputDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.OutputDrop = ""
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid OutputDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.OutputDrop = "invalid"
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("empty ForwardDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.ForwardDrop = ""
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid ForwardDrop", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.ForwardDrop = "invalid"
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid negative InputPriority", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.InputPriority = fakers.RandInt(-100, -51)
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid positive InputPriority", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.InputPriority = fakers.RandInt(51, 100)
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid negative OutputPriority", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.OutputPriority = fakers.RandInt(-100, -51)
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid positive OutputPriority", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.OutputPriority = fakers.RandInt(51, 100)
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid negative ForwardPriority", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.ForwardPriority = fakers.RandInt(-100, -51)
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})

	t.Run("invalid positive ForwardPriority", func(t *testing.T) {
		s := defaultPolicyForTest()
		s.ForwardPriority = fakers.RandInt(51, 100)
		if err := s.Validate(); err == nil {
			t.Errorf("policy.Validate() error = nil, wantErr not nil")
		}
	})
}

func defaultPolicyForTest() *policy {
	return &policy{
		DefaultAllowInput:   fakers.RandBool(),
		DefaultAllowOutput:  fakers.RandBool(),
		DefaultAllowForward: fakers.RandBool(),
		InputDrop:           fakers.RandFirewallDrop(),
		InputPriority:       fakers.RandInt(-50, 50),
		OutputDrop:          fakers.RandFirewallDrop(),
		OutputPriority:      fakers.RandInt(-50, 50),
		ForwardDrop:         fakers.RandFirewallDrop(),
		ForwardPriority:     fakers.RandInt(-50, 50),
	}
}
