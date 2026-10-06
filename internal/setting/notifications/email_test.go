package notifications

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
	"github.com/wneessen/go-mail"
)

func TestEmail_Validate_Success(t *testing.T) {
	s := defaultEmailForTest()
	if err := s.Validate(); err != nil {
		t.Errorf("Email.Validate() error = %v, want nil", err)
	}
}

func TestEmail_Validate_Error(t *testing.T) {
	t.Run("empty host", func(t *testing.T) {
		s := defaultEmailForTest()
		s.Host = ""
		if err := s.Validate(); err == nil {
			t.Errorf("Email.Validate() error = nil, want non-nil")
		}
	})

	t.Run("empty username", func(t *testing.T) {
		s := defaultEmailForTest()
		s.Username = ""
		if err := s.Validate(); err == nil {
			t.Errorf("Email.Validate() error = nil, want non-nil")
		}
	})

	t.Run("empty password", func(t *testing.T) {
		s := defaultEmailForTest()
		s.Password = ""
		if err := s.Validate(); err == nil {
			t.Errorf("Email.Validate() error = nil, want non-nil")
		}
	})

	t.Run("empty from", func(t *testing.T) {
		s := defaultEmailForTest()
		s.From = ""
		if err := s.Validate(); err == nil {
			t.Errorf("Email.Validate() error = nil, want non-nil")
		}
	})

	t.Run("empty to", func(t *testing.T) {
		s := defaultEmailForTest()
		s.To = ""
		if err := s.Validate(); err == nil {
			t.Errorf("Email.Validate() error = nil, want non-nil")
		}
	})

	t.Run("zero Port", func(t *testing.T) {
		s := defaultEmailForTest()
		s.Port = 0
		if err := s.Validate(); err == nil {
			t.Errorf("Email.Validate() error = nil, want non-nil")
		}
	})

	t.Run("negative Port", func(t *testing.T) {
		s := defaultEmailForTest()
		s.Port = -1
		if err := s.Validate(); err == nil {
			t.Errorf("Email.Validate() error = nil, want non-nil")
		}
	})
}

func TestEmail_ToTLSConfig_Success(t *testing.T) {
	s := defaultEmailForTest()
	if _, err := s.ToTLSConfig(); err != nil {
		t.Errorf("Email.Validate() error = %v, want nil", err)
	}
}

func TestEmail_ToTLSConfig_Error(t *testing.T) {
	t.Run("empty TLSMode", func(t *testing.T) {
		s := defaultEmailForTest()
		s.TLSMode = ""
		if _, err := s.ToTLSConfig(); err == nil {
			t.Errorf("Email.ToTLSConfig() error = nil, want non-nil")
		}
	})

	t.Run("invalid TLSMode", func(t *testing.T) {
		s := defaultEmailForTest()
		s.TLSMode = "invalid"
		if _, err := s.ToTLSConfig(); err == nil {
			t.Errorf("Email.ToTLSConfig() error = nil, want non-nil")
		}
	})

	t.Run("empty TLSPolicy", func(t *testing.T) {
		s := defaultEmailForTest()
		s.TLSPolicy = ""
		if _, err := s.ToTLSConfig(); err == nil {
			t.Errorf("Email.ToTLSConfig() error = nil, want non-nil")
		}
	})

	t.Run("invalid TLSPolicy", func(t *testing.T) {
		s := defaultEmailForTest()
		s.TLSPolicy = "invalid"
		if _, err := s.ToTLSConfig(); err == nil {
			t.Errorf("Email.ToTLSConfig() error = nil, want non-nil")
		}
	})
}

func TestEmail_ParseAuthType_Success(t *testing.T) {
	authType := fakeRandAuthType()
	if authTypeMail, err := ParseAuthType(authType); err != nil {
		t.Errorf("Email.ParseAuthType() error = %v, want nil", err)
	} else if !testAuthType(authTypeMail, authType) {
		t.Errorf("Email.ParseAuthType() = %v, want %v", string(authTypeMail), authType)
	}
}

func TestEmail_ParseAuthType_Error(t *testing.T) {
	t.Run("empty auth type", func(t *testing.T) {
		if _, err := ParseAuthType(""); err == nil {
			t.Errorf("Email.ParseAuthType() error = nil, want non-nil")
		}
	})
	t.Run("invalid auth type", func(t *testing.T) {
		if _, err := ParseAuthType("INVALID"); err == nil {
			t.Errorf("Email.ParseAuthType() error = nil, want non-nil")
		}
	})
}

func defaultEmailForTest() *Email {
	return &Email{
		Host:      "test",
		Port:      fakers.RandInt(1, 65535),
		Username:  "test",
		Password:  "test",
		AuthType:  fakeRandAuthType(),
		TLSMode:   fakeRandTLSMode(),
		TLSPolicy: fakeRandTLSPolicy(),
		TLSVerify: fakers.RandBool(),
		From:      "test@localhost",
		To:        "test@localhost",
	}
}

func fakeRandAuthType() string {
	items := []string{"PLAIN", "LOGIN", "CRAM-MD5", "NONE"}
	return fakers.RandItem(items)
}

func fakeRandTLSMode() string {
	items := []string{"NONE", "STARTTLS", "IMPLICIT"}
	return fakers.RandItem(items)
}

func fakeRandTLSPolicy() string {
	items := []string{"MANDATORY", "OPPORTUNISTIC"}
	return fakers.RandItem(items)
}

func testAuthType(smtpAuthType mail.SMTPAuthType, authType string) bool {
	if smtpAuthType == mail.SMTPAuthNoAuth {
		return authType == "NONE"
	}

	return string(smtpAuthType) == authType
}
