package firewall

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
)

type Metadata interface {
	Metadata() (*entity.Metadata, error)
	UpdateMetadata(version string, checksum string) error
}

type metadata struct {
	metadataRepo repository.MetadataRepository
	key          string
}

func NewMetadata(metadataRepo repository.MetadataRepository) Metadata {
	return &metadata{
		metadataRepo: metadataRepo,
		key:          "firewall-file-nft",
	}
}

func (m *metadata) Metadata() (*entity.Metadata, error) {
	return m.metadataRepo.Get(m.key)
}

func (m *metadata) UpdateMetadata(version string, checksum string) error {
	metadataEntity := &entity.Metadata{
		Checksum: checksum,
		Version:  version,
	}

	return m.metadataRepo.Update(m.key, metadataEntity)
}
