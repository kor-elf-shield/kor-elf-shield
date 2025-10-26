package filesystem

import (
	"errors"
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
