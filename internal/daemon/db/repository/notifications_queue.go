package repository

import (
	"encoding/json"
	"errors"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"go.etcd.io/bbolt"
	bboltErrors "go.etcd.io/bbolt/errors"
)

type NotificationsQueueRepository interface {
	Add(q entity.NotificationsQueue) error
	Get(limit int) (map[string]entity.NotificationsQueue, error)
	Delete(id string) error

	// Count - return size of notifications queue in db
	Count() (int, error)
	Clear() error
}

type notificationsQueueRepository struct {
	db     *bbolt.DB
	bucket string
}

func NewNotificationsQueueRepository(appDB *bbolt.DB) NotificationsQueueRepository {
	return &notificationsQueueRepository{
		db:     appDB,
		bucket: notificationsQueueBucket,
	}
}

func (r *notificationsQueueRepository) Add(q entity.NotificationsQueue) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(r.bucket))
		if err != nil {
			return err
		}

		data, err := json.Marshal(q)
		if err != nil {
			return err
		}

		id, err := nextID(bucket)
		if err != nil {
			return err
		}

		return bucket.Put(id, data)
	})
}

func (r *notificationsQueueRepository) Get(limit int) (map[string]entity.NotificationsQueue, error) {
	notifications := make(map[string]entity.NotificationsQueue)

	if limit <= 0 {
		return notifications, nil
	}

	err := r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(r.bucket))
		if bucket == nil {
			return nil
		}

		c := bucket.Cursor()
		for k, v := c.First(); k != nil && len(notifications) < limit; k, v = c.Next() {
			var q entity.NotificationsQueue
			if err := json.Unmarshal(v, &q); err != nil {
				return err
			}
			notifications[string(k)] = q
		}

		return nil
	})

	return notifications, err
}

func (r *notificationsQueueRepository) Delete(id string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(r.bucket))
		if bucket == nil {
			return nil
		}

		return bucket.Delete([]byte(id))
	})
}

func (r *notificationsQueueRepository) Count() (int, error) {
	count := 0

	err := r.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(r.bucket))
		if bucket == nil {
			return nil
		}

		count = bucket.Stats().KeyN

		return nil
	})

	return count, err
}

func (r *notificationsQueueRepository) Clear() error {
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
