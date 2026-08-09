package repository

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"go.etcd.io/bbolt"
	bboltErrors "go.etcd.io/bbolt/errors"
)

type BruteForceProtectionGroupRepository interface {
	Update(name string, ip net.IP, partition *string, f func(*entity.BruteForceProtectionGroup) (*entity.BruteForceProtectionGroup, error)) error
	Clear() error
}

type bruteForceProtectionGroupRepository struct {
	db     *bbolt.DB
	bucket string
}

func NewBruteForceProtectionGroupRepository(appDB *bbolt.DB) BruteForceProtectionGroupRepository {
	return &bruteForceProtectionGroupRepository{
		db:     appDB,
		bucket: bruteForceProtectionGroupBucket,
	}
}

func (r *bruteForceProtectionGroupRepository) Update(name string, ip net.IP, partition *string, f func(*entity.BruteForceProtectionGroup) (*entity.BruteForceProtectionGroup, error)) error {
	entityGroup := &entity.BruteForceProtectionGroup{}
	entityGroup.Reset()

	return r.db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(r.bucket))
		if err != nil {
			return err
		}
		key, err := keyGroupIP(name, ip, partition)
		if err != nil {
			return err
		}

		group := b.Get(key)
		if group != nil {
			err = json.Unmarshal(group, entityGroup)
			if err != nil {
				return fmt.Errorf("failed to unmarshal brute force protection group: %w", err)
			}
		}

		entityGroup, err = f(entityGroup)
		if err != nil {
			return err
		}

		data, err := json.Marshal(entityGroup)
		if err != nil {
			return err
		}
		return b.Put(key, data)
	})
}

func (r *bruteForceProtectionGroupRepository) Clear() error {
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

func keyGroupIP(groupID string, ip net.IP, partition *string) ([]byte, error) {
	if ip == nil {
		return nil, fmt.Errorf("ip cannot be nil")
	}

	if len(groupID) == 0 {
		return nil, fmt.Errorf("group id cannot be empty")
	}

	if ip.To4() == nil && ip.To16() == nil {
		return nil, fmt.Errorf("ip is neither IPv4 nor IPv6")
	}

	var ipAddr net.IP
	if ip.To4() != nil {
		ipAddr = ip.To4()
	} else {
		ipAddr = ip.To16()
	}

	if partition != nil {
		partitionHash := sha256.Sum256([]byte(*partition))
		k := make([]byte, 0, len(groupID)+1+len(ipAddr)+1+len(partitionHash))
		k = append(k, groupID...)
		k = append(k, 0x00)
		k = append(k, ipAddr...)
		k = append(k, 0x00)
		k = append(k, partitionHash[:]...)
		return k, nil
	}

	k := make([]byte, 0, len(groupID)+1+len(ipAddr))
	k = append(k, groupID...)
	k = append(k, 0x00)
	k = append(k, ipAddr...)
	return k, nil
}
