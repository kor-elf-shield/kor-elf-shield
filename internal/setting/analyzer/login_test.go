package analyzer

import (
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/i18nmock"
)

func TestLogin_Validate_Success(t *testing.T) {
	s := defaultLoginForTest()
	if err := s.Validate(); err != nil {
		t.Errorf("Login.Validate() error = %v, wantErr nil", err)
	}
}

func TestLogin_ToSources_Success(t *testing.T) {
	i18nmock.Use(t)

	t.Run("Enabled = false", func(t *testing.T) {
		s := defaultLoginForTest()
		s.Enabled = false
		s.SSHEnable = true
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("Login.ToSources() error = %v, wantErr nil", err)
		} else if len(sources) > 0 {
			t.Errorf("Login.ToSources() error = %d, wantErr nil", len(sources))
		}
	})

	t.Run("SSHEnable = true", func(t *testing.T) {
		s := defaultLoginForTest()
		s.SSHEnable = true
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("Login.ToSources() error = %v, wantErr nil", err)
		} else if !testSourceSSH(sources) {
			t.Errorf("Login.ToSources() not source SSH")
		}
	})

	t.Run("SSHEnable = false", func(t *testing.T) {
		s := defaultLoginForTest()
		s.SSHEnable = false
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("Login.ToSources() error = %v, wantErr nil", err)
		} else if testSourceSSH(sources) {
			t.Errorf("Login.ToSources() source SSH")
		}
	})

	t.Run("LocalEnable = true", func(t *testing.T) {
		s := defaultLoginForTest()
		s.LocalEnable = true
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("Login.ToSources() error = %v, wantErr nil", err)
		} else if !testSourceLocal(sources) {
			t.Errorf("Login.ToSources() not source Local")
		}
	})

	t.Run("LocalEnable = false", func(t *testing.T) {
		s := defaultLoginForTest()
		s.LocalEnable = false
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("Login.ToSources() error = %v, wantErr nil", err)
		} else if testSourceLocal(sources) {
			t.Errorf("Login.ToSources() source Local")
		}
	})

	t.Run("SuEnable = true", func(t *testing.T) {
		s := defaultLoginForTest()
		s.SuEnable = true
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("Login.ToSources() error = %v, wantErr nil", err)
		} else if !testSourceSu(sources) {
			t.Errorf("Login.ToSources() not source Su")
		}
	})

	t.Run("SuEnable = false", func(t *testing.T) {
		s := defaultLoginForTest()
		s.SuEnable = false
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("Login.ToSources() error = %v, wantErr nil", err)
		} else if testSourceSu(sources) {
			t.Errorf("Login.ToSources() source Su")
		}
	})

	t.Run("SudoEnable = true", func(t *testing.T) {
		s := defaultLoginForTest()
		s.SudoEnable = true
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("Login.ToSources() error = %v, wantErr nil", err)
		} else if !testSourceSudo(sources) {
			t.Errorf("Login.ToSources() not source Sudo")
		}
	})

	t.Run("SudoEnable = false", func(t *testing.T) {
		s := defaultLoginForTest()
		s.SudoEnable = false
		if sources, err := s.ToSources(); err != nil {
			t.Errorf("Login.ToSources() error = %v, wantErr nil", err)
		} else if testSourceSudo(sources) {
			t.Errorf("Login.ToSources() source Sudo")
		}
	})
}

func defaultLoginForTest() *Login {
	return &Login{
		Enabled:     true,
		Notify:      fakers.RandBool(),
		SSHEnable:   fakers.RandBool(),
		SSHNotify:   fakers.RandBool(),
		LocalEnable: fakers.RandBool(),
		LocalNotify: fakers.RandBool(),
		SuEnable:    fakers.RandBool(),
		SuNotify:    fakers.RandBool(),
		SudoEnable:  fakers.RandBool(),
		SudoNotify:  fakers.RandBool(),
	}
}

func testSourceSSH(sources []*config.Source) bool {
	return testAlertSource(sources, "_login-ssh")
}

func testSourceLocal(sources []*config.Source) bool {
	return testAlertSource(sources, "_login-local")
}

func testSourceSu(sources []*config.Source) bool {
	return testAlertSource(sources, "_login-su")
}

func testSourceSudo(sources []*config.Source) bool {
	return testAlertSource(sources, "_login-sudo")
}

func testAlertSource(sources []*config.Source, name string) bool {
	if len(sources) == 0 {
		return false
	}
	for _, source := range sources {
		if source.AlertRule == nil {
			continue
		}
		if source.AlertRule.Name == name {
			return true
		}
	}
	return false
}
