package repository

import (
	"encoding/json"
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"go.etcd.io/bbolt"
)

type BruteForceProtectionNotifyPolicyRepository interface {
	Update(ruleName string, f func(*entity.BruteForceProtectionNotifyPolicy) (*entity.BruteForceProtectionNotifyPolicy, error)) error
}

type bruteForceProtectionNotifyPolicyRepository struct {
	db     *bbolt.DB
	bucket string
}

func NewBruteForceProtectionNotifyPolicyRepository(appDB *bbolt.DB) BruteForceProtectionNotifyPolicyRepository {
	return &bruteForceProtectionNotifyPolicyRepository{
		db:     appDB,
		bucket: bruteForceProtectionNotifyPolicyBucket,
	}
}

func (r *bruteForceProtectionNotifyPolicyRepository) Update(ruleName string, f func(*entity.BruteForceProtectionNotifyPolicy) (*entity.BruteForceProtectionNotifyPolicy, error)) error {
	entityNotify := &entity.BruteForceProtectionNotifyPolicy{}

	return r.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(r.bucket))
		if err != nil {
			return err
		}
		key := []byte(ruleName)

		notify := b.Get(key)
		if notify != nil {
			err = json.Unmarshal(notify, entityNotify)
			if err != nil {
				return fmt.Errorf("failed to unmarshal brute force protection notify policy: %w", err)
			}
		}

		entityNotify, err = f(entityNotify)
		if err != nil {
			return err
		}

		data, err := json.Marshal(entityNotify)
		if err != nil {
			return err
		}
		return b.Put(key, data)
	})
}
