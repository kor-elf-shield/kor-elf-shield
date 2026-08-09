package repository

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"go.etcd.io/bbolt"
	bboltErrors "go.etcd.io/bbolt/errors"
)

type AlertGroupRepository interface {
	Update(name string, partition *string, f func(*entity.AlertGroup) (*entity.AlertGroup, error)) error
	Clear() error
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

func (r *alertGroupRepository) Update(name string, partition *string, f func(*entity.AlertGroup) (*entity.AlertGroup, error)) error {
	entityAlertGroup := &entity.AlertGroup{}
	entityAlertGroup.Reset()

	return r.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(r.bucket))
		if err != nil {
			return err
		}
		key, err := keyGroup(name, partition)
		if err != nil {
			return err
		}

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

func (r *alertGroupRepository) Clear() error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		err := tx.DeleteBucket([]byte(r.bucket))
		if errors.Is(err, bboltErrors.ErrBucketNotFound) {
			// If the bucket may not exist, ignore ErrBucketNotFound
			return nil
		}
		_, err = tx.CreateBucketIfNotExists([]byte(r.bucket))
		return err
	})
}

func keyGroup(groupID string, partition *string) ([]byte, error) {
	if len(groupID) == 0 {
		return nil, fmt.Errorf("group id cannot be empty")
	}

	if partition == nil {
		return []byte(groupID), nil
	}

	partitionHash := sha256.Sum256([]byte(*partition))

	k := make([]byte, 0, len(groupID)+1+len(partitionHash))
	k = append(k, groupID...)
	k = append(k, 0x00)
	k = append(k, partitionHash[:]...)
	return k, nil
}
