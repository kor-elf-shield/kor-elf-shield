package repository

import (
	"encoding/json"
	"errors"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"go.etcd.io/bbolt"
	bboltErrors "go.etcd.io/bbolt/errors"
)

type BlockingRepository interface {
	Add(blockedIP entity.Blocking) error
	List(callback func(entity.Blocking) error) error
	DeleteExpired(limit int) (int, error)
	Clear() error
}

type blocking struct {
	db     *bbolt.DB
	bucket string
}

func NewBlockingRepository(appDB *bbolt.DB) BlockingRepository {
	return &blocking{
		db:     appDB,
		bucket: blockingBucket,
	}
}

func (r *blocking) Add(blockedIP entity.Blocking) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(r.bucket))
		if err != nil {
			return err
		}

		data, err := json.Marshal(blockedIP)
		if err != nil {
			return err
		}

		id, err := nextKeyByExpire(bucket, uint64(blockedIP.ExpireAtUnix))
		if err != nil {
			return err
		}

		return bucket.Put(id, data)
	})
}

func (r *blocking) List(callback func(entity.Blocking) error) error {
	return r.db.View(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(r.bucket))
		if err != nil {
			return err
		}

		return bucket.ForEach(func(_, v []byte) error {
			blockedIP := entity.Blocking{}
			err := json.Unmarshal(v, &blockedIP)
			if err != nil {
				return err
			}

			if err := callback(blockedIP); err != nil {
				return err
			}

			return nil
		})
	})
}

func (r *blocking) DeleteExpired(limit int) (int, error) {
	if limit <= 0 {
		return 0, nil
	}

	var deleted int
	err := r.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(r.bucket))
		if err != nil {
			return err
		}

		now := time.Now().Unix()
		c := bucket.Cursor()
		deleted = 0

		for k, v := c.First(); k != nil && deleted < limit; {
			blockedIP := entity.Blocking{}
			if err := json.Unmarshal(v, &blockedIP); err != nil {
				return err
			}

			if blockedIP.ExpireAtUnix <= 0 {
				k, v = c.Next()
				continue
			}

			if blockedIP.ExpireAtUnix > now {
				// Not expired yet
				break
			}

			nextK, nextV := c.Next()
			if err := bucket.Delete(k); err != nil {
				return err
			}
			deleted++
			k = nextK
			v = nextV
		}

		return nil
	})

	return deleted, err
}

func (r *blocking) Clear() error {
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
