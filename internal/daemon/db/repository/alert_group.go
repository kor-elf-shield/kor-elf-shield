package repository

import (
	"encoding/json"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"go.etcd.io/bbolt"
)

type AlertGroupRepository interface {
	Update(name string, f func(*entity.AlertGroup) (*entity.AlertGroup, error)) error
}

type alertGroupRepository struct {
	db     *bbolt.DB
	bucket string
}

func NewAlertGroupRepository(appDB *bbolt.DB) AlertGroupRepository {
	return &alertGroupRepository{
		db:     appDB,
		bucket: alertGroupBucket,
	}
}

func (r *alertGroupRepository) Update(name string, f func(*entity.AlertGroup) (*entity.AlertGroup, error)) error {
	entityAlertGroup := &entity.AlertGroup{}
	entityAlertGroup.Reset()

	return r.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(r.bucket))
		if err != nil {
			return err
		}
		key := []byte(name)

		group := b.Get(key)
		if group != nil {
			err = json.Unmarshal(group, entityAlertGroup)
			if err != nil {
				return fmt.Errorf("failed to unmarshal alert group: %w", err)
			}
		}

		entityAlertGroup, err = f(entityAlertGroup)
		if err != nil {
			return err
		}

		data, err := json.Marshal(entityAlertGroup)
		if err != nil {
			return err
		}
		return b.Put(key, data)
	})
}
