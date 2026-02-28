package db

import (
	"errors"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/filesystem"
	"go.etcd.io/bbolt"
)

const (
	appDB      = "app.db"
	securityDB = "security.db"
)

type Repositories interface {
	NotificationsQueue() repository.NotificationsQueueRepository
	AlertGroup() repository.AlertGroupRepository
	BruteForceProtectionGroup() repository.BruteForceProtectionGroupRepository
	Blocking() repository.BlockingRepository

	Close() error
}

type repositories struct {
	notificationsQueue        repository.NotificationsQueueRepository
	alertGroup                repository.AlertGroupRepository
	bruteForceProtectionGroup repository.BruteForceProtectionGroupRepository
	blocking                  repository.BlockingRepository

	db []*bbolt.DB
}

func New(dataDir string) (Repositories, error) {
	if dataDir == "" {
		return &repositories{}, errors.New("data directory is empty")
	}
	if dataDir[len(dataDir)-1:] != "/" {
		dataDir += "/"
	}

	err := filesystem.EnsureDir(dataDir)
	if err != nil {
		return &repositories{}, err
	}

	appDB, err := bbolt.Open(dataDir+appDB, 0600, &bbolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return &repositories{}, err
	}

	securityDB, err := bbolt.Open(dataDir+securityDB, 0600, &bbolt.Options{Timeout: 3 * time.Second})

	return &repositories{
		notificationsQueue:        repository.NewNotificationsQueueRepository(appDB),
		alertGroup:                repository.NewAlertGroupRepository(appDB),
		bruteForceProtectionGroup: repository.NewBruteForceProtectionGroupRepository(securityDB),
		blocking:                  repository.NewBlockingRepository(securityDB),

		db: []*bbolt.DB{appDB, securityDB},
	}, nil
}

func (r *repositories) NotificationsQueue() repository.NotificationsQueueRepository {
	return r.notificationsQueue
}

func (r *repositories) AlertGroup() repository.AlertGroupRepository {
	return r.alertGroup
}

func (r *repositories) BruteForceProtectionGroup() repository.BruteForceProtectionGroupRepository {
	return r.bruteForceProtectionGroup
}

func (r *repositories) Blocking() repository.BlockingRepository {
	return r.blocking
}

func (r *repositories) Close() error {
	for _, db := range r.db {
		_ = db.Close()
	}

	return nil
}
