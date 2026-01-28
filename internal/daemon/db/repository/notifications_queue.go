package repository

import (
	"encoding/json"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"go.etcd.io/bbolt"
)

type NotificationsQueueRepository interface {
	Add(q entity.NotificationsQueue) error
	Get(limit int) (map[string]entity.NotificationsQueue, error)
	Delete(id string) error
}

type notificationsQueueRepository struct {
	db     *bbolt.DB
	bucket string
}

func NewNotificationsQueueRepository(appDB *bbolt.DB) NotificationsQueueRepository {
	return &notificationsQueueRepository{
		db:     appDB,
		bucket: notificationsQueue,
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
