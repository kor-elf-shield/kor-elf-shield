package info

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
)

type MetadataContainer interface {
	FirewallFileNft() Metadata
}

type Metadata interface {
	Get() (string, error)
	Update(data string) error
}

func NewMetadataContainer(
	firewallFileNft Metadata,
) MetadataContainer {
	return &metadataContainer{
		firewallFileNft: firewallFileNft,
	}
}

func NewMetadata(repo repository.MetadataRepository, key string) Metadata {
	return &metadata{
		key:  key,
		repo: repo,
	}
}

type metadataContainer struct {
	firewallFileNft Metadata
}

func (m *metadataContainer) FirewallFileNft() Metadata {
	return m.firewallFileNft
}

type metadata struct {
	key  string
	repo repository.MetadataRepository
}

func (m *metadata) Get() (string, error) {
	meta, err := m.repo.Get(m.key)
	if err != nil {
		return "", err
	}
	return meta.Value, nil
}

func (m *metadata) Update(value string) error {
	meta := &entity.Metadata{Value: value}
	return m.repo.Update(m.key, meta)
}
