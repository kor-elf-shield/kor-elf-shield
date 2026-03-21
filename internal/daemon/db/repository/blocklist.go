package repository

import (
	"encoding/json"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"go.etcd.io/bbolt"
)

type BlocklistRepository interface {
	Get(name string) (*entity.Blocklist, error)
	Update(name string, entity *entity.Blocklist) error
}

type blocklistRepository struct {
	db     *bbolt.DB
	bucket string
}

func NewBlocklistRepository(appDB *bbolt.DB) BlocklistRepository {
	return &blocklistRepository{
		db:     appDB,
		bucket: blocklistBucket,
	}
}

func (r *blocklistRepository) Get(name string) (*entity.Blocklist, error) {
	blocklistEntity := &entity.Blocklist{}

	err := r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(r.bucket))
		if bucket == nil {
			return nil
		}

		data := bucket.Get([]byte(name))
		if data == nil {
			return nil
		}

		return json.Unmarshal(data, blocklistEntity)
	})

	if err != nil {
		return nil, err
	}

	return blocklistEntity, err
}

func (r *blocklistRepository) Update(name string, blocklistEntity *entity.Blocklist) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(r.bucket))
		if err != nil {
			return err
		}
		key := []byte(name)

		data, err := json.Marshal(blocklistEntity)
		if err != nil {
			return err
		}
		return b.Put(key, data)
	})
}
