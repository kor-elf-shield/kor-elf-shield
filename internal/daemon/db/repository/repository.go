package repository

import (
	"encoding/binary"
	"math"

	"go.etcd.io/bbolt"
)

const (
	notificationsQueueBucket        = "notifications_queue"
	alertGroupBucket                = "alert_group"
	bruteForceProtectionGroupBucket = "brute_force_protection_group"
	blockingBucket                  = "blocking"
	blocklistBucket                 = "blocklist"
	metadataBucket                  = "metadata"
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

func nextKeyByExpire(b *bbolt.Bucket, expireUnixAt uint64) ([]byte, error) {
	seq, err := b.NextSequence()
	if err != nil {
		return nil, err
	}

	// 0 = "forever" -> sort after any finite timestamp
	if expireUnixAt == 0 {
		expireUnixAt = math.MaxUint64
	}

	// 8 bytes expire + 8 bytes seq
	key := make([]byte, 16)

	// Important: BigEndian, so that sorting by bytes matches sorting by number.
	binary.BigEndian.PutUint64(key[0:8], expireUnixAt)
	binary.BigEndian.PutUint64(key[8:16], seq)

	return key, nil
}
