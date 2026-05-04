package repository

import (
	"encoding/json"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"go.etcd.io/bbolt"
)

type MetadataRepository interface {
	Get(name string) (*entity.Metadata, error)
	Update(name string, entity *entity.Metadata) error
}

type metadataRepository struct {
	db     *bbolt.DB
	bucket string
}

func NewMetadataRepository(appDB *bbolt.DB) MetadataRepository {
	return &metadataRepository{
		db:     appDB,
		bucket: metadataBucket,
	}
}

func (r *metadataRepository) Get(name string) (*entity.Metadata, error) {
	metadataEntity := &entity.Metadata{}

	err := r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(r.bucket))
		if bucket == nil {
			return nil
		}

		data := bucket.Get([]byte(name))
		if data == nil {
			return nil
		}

		return json.Unmarshal(data, metadataEntity)
	})

	if err != nil {
		return nil, err
	}

	return metadataEntity, err
}

func (r *metadataRepository) Update(name string, entity *entity.Metadata) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(r.bucket))
		if err != nil {
			return err
		}
		key := []byte(name)

		data, err := json.Marshal(entity)
		if err != nil {
			return err
		}
		return b.Put(key, data)
	})
}
