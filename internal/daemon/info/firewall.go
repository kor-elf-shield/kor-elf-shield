package info

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
)

func NewMetadataFirewallFileNft(repo repository.MetadataRepository) Metadata {
	return NewMetadata(repo, entity.MetadataKeyFirewallFileNft)
}
