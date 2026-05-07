package info

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

func IsVersionChanged(repo repository.MetadataRepository, version string, logger log.Logger) bool {
	metadata := NewMetadata(repo, entity.MetadataKeyVersion)
	if v, err := metadata.Get(); err != nil {
		logger.Error(err.Error())
		return false
	} else if v != version {
		if err := metadata.Update(version); err != nil {
			logger.Error(err.Error())
		}
		return true
	}

	return false
}
