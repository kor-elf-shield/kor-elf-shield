package filesystem

import (
	"errors"
	"os"
	"path/filepath"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/i18n"
)

// EnsureDir ensures the existence of a directory with the required rights.
func EnsureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := mkdirAllWithChmod(dir, 0700); err != nil {
			return errors.New(i18n.Lang.T("failed to create directory", map[string]any{
				"Directory": dir,
			}))
		}
	}

	return nil
}

func mkdirAllWithChmod(dir string, mode os.FileMode) error {
	var paths []string

	for {
		paths = append([]string{dir}, paths...)
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	for _, p := range paths {
		info, err := os.Stat(p)
		if err == nil && info.IsDir() {
			// Skipping existing ones
			continue
		}

		if err := os.Mkdir(p, mode); err != nil {
			if os.IsExist(err) {
				// Skipping existing ones
				continue
			}
			return err
		} else {
			// in main.go unix.Umask(0o177)
			// The directory was created with permissions 0600
			if err := os.Chmod(p, mode); err != nil {
				return err
			}
		}
	}

	return nil
}
