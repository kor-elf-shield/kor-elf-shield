package repository

import (
	"encoding/binary"

	"go.etcd.io/bbolt"
)

const (
	notificationsQueueBucket = "notifications_queue"
	alertGroupBucket         = "alert_group"
)

func nextID(b *bbolt.Bucket) ([]byte, error) {
	seq, err := b.NextSequence()
	if err != nil {
		return nil, err
	}

	key := make([]byte, 8)
	binary.BigEndian.PutUint64(key, seq)
	return key, nil
}
