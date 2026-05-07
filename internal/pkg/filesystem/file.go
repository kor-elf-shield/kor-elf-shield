package filesystem

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
)

func FileHasWritePermissions(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return errors.New(i18n.Lang.T("failed to open file for writing", map[string]any{
			"File": path,
		}))
	}
	_ = file.Close()

	return nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func FileChecksum(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	hash := sha256.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
