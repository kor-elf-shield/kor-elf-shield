package checksum

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/secret"
)

// CheckerFile — file integrity check interface.
type CheckerFile interface {
	// ChecksumFile returns the checksum in the format v1:<hex>.
	ChecksumFile(path string) (string, error)

	// ValidateFile checks the expected checksum for the file.
	ValidateFile(path string, expected string) error
}

type checkerFile struct {
	secretProvider secret.Provider
}

func New(secretProvider secret.Provider) (CheckerFile, error) {
	if secretProvider == nil {
		return nil, fmt.Errorf("secret provider is nil")
	}
	return &checkerFile{
		secretProvider: secretProvider,
	}, nil
}

func (c *checkerFile) ChecksumFile(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path is empty")
	}

	secretBytes, err := c.secretProvider.Secret()
	if err != nil {
		return "", fmt.Errorf("get secret: %w", err)
	}

	key := deriveFileKey(secretBytes, path)
	defer zeroBytes(key)

	macHex, err := hmacFileHex(path, key)
	if err != nil {
		return "", err
	}

	return macHex, nil
}

func (c *checkerFile) ValidateFile(path string, expected string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("path is empty")
	}

	expectedBytes, err := decodeChecksum(expected)
	if err != nil {
		return err
	}

	secretBytes, err := c.secretProvider.Secret()
	if err != nil {
		return fmt.Errorf("get secret: %w", err)
	}

	key := deriveFileKey(secretBytes, path)
	defer zeroBytes(key)

	actualHex, err := hmacFileHex(path, key)
	if err != nil {
		return err
	}

	actualBytes, err := hex.DecodeString(actualHex)
	if err != nil {
		return fmt.Errorf("internal checksum decode error: %w", err)
	}

	if subtle.ConstantTimeCompare(expectedBytes, actualBytes) != 1 {
		return fmt.Errorf("checksum mismatch")
	}

	return nil
}

// deriveFileKey creates a separate key for each file:
// K_file = HMAC(masterSecret, "file:"+cleanPath)
func deriveFileKey(masterSecret []byte, path string) []byte {
	clean := filepath.Clean(path)

	mac := hmac.New(sha256.New, masterSecret)
	_, _ = mac.Write([]byte("file:"))
	_, _ = mac.Write([]byte(clean))
	return mac.Sum(nil)
}

func hmacFileHex(path string, key []byte) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	mac := hmac.New(sha256.New, key)
	if _, err := io.Copy(mac, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(mac.Sum(nil)), nil
}

func decodeChecksum(s string) ([]byte, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, fmt.Errorf("expected checksum is empty")
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("expected checksum is not valid hex: %w", err)
	}
	if len(b) != sha256.Size {
		return nil, fmt.Errorf("expected checksum has invalid size")
	}
	return b, nil
}

func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
