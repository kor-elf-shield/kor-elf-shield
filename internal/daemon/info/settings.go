package info

import (
	"fmt"
	"os"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

func IsSettingsChanged(repo repository.MetadataRepository, listPathFiles map[string]string, logger log.Logger) bool {
	isChanged := false

	for k, f := range listPathFiles {
		hash, err := fileHash(f)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to get hash of %s: %s", f, err))
			continue
		}

		settingMetadata := NewMetadata(repo, entity.KeySetting(k))
		metadataValue, err := settingMetadata.Get()
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to get metadata value: %s", err))
		}

		if metadataValue != hash {
			isChanged = true
			if err := settingMetadata.Update(hash); err != nil {
				logger.Error(fmt.Sprintf("Failed to update metadata: %s", err))
			}
		}
	}

	return isChanged
}

func fileHash(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	size := info.Size()
	modTime := info.ModTime()

	return fmt.Sprintf("%d-%d", size, modTime.Unix()), nil
}
