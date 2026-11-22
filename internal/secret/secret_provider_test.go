package secret

import (
	"fmt"
	"path/filepath"
	"testing"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/logger_mock"
)

func TestNewSecret_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".secret")
	logger := logger_mock.New()
	if secretProvider, err := NewSecret(path, fakers.RandInt(0, 100), logger); err != nil {
		t.Errorf("NewSecret() error = %v", err)
	} else {
		secretProvider.Close()
	}
}

func TestNewSecret_Error(t *testing.T) {
	t.Run("Empty path", func(t *testing.T) {
		if secretProvider, err := NewSecret("", fakers.RandInt(0, 100), logger_mock.New()); err == nil {
			t.Errorf("NewSecret() error = %v", err)
			secretProvider.Close()
		} else if secretProvider != nil {
			secretProvider.Close()
		}
	})

	t.Run("Negative Cache", func(t *testing.T) {
		for _, count := range []int{-1, -100} {
			t.Run(fmt.Sprintf("invalid Count: %d", count), func(t *testing.T) {
				if secretProvider, err := NewSecret(filepath.Join(t.TempDir(), ".secret"), count, logger_mock.New()); err == nil {
					t.Errorf("NewSecret() error = %v", err)
					secretProvider.Close()
				} else if secretProvider != nil {
					secretProvider.Close()
				}
			})
		}
	})

	t.Run("Nil logger", func(t *testing.T) {
		if secretProvider, err := NewSecret(filepath.Join(t.TempDir(), ".secret"), fakers.RandInt(0, 100), nil); err == nil {
			t.Errorf("NewSecret() error = %v", err)
			secretProvider.Close()
		} else if secretProvider != nil {
			secretProvider.Close()
		}
	})
}

func TestGetSecret_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".secret")
	logger := logger_mock.New()

	secretProvider, err := NewSecret(path, fakers.RandInt(0, 100), logger)
	if err != nil {
		t.Errorf("NewSecret() error = %v", err)
	}
	if secret, err := secretProvider.Secret(); err != nil {
		t.Errorf("GetSecret() error = %v", err)
	} else if secret == nil {
		t.Errorf("GetSecret() secret = %v, want non-nil string", secret)
	}

	secretProvider.Close()
}
