package filesystem

import (
	"errors"
	"kor-elf-shield/internal/i18n"
	"os"
)

// EnsureDir ensures the existence of a directory with the required rights.
func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return errors.New(i18n.Lang.T("failed to create directory", map[string]any{
				"Directory": dir,
			}))
		}
		// The directory was created with permissions 0600
		_ = os.Chmod(dir, 0700)
	}

	return nil
}
