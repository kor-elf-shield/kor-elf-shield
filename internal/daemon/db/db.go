package db

import (
	"errors"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/filesystem"
	"go.etcd.io/bbolt"
)

const (
	app = "app.db"
)

type Repositories interface {
	NotificationsQueue() repository.NotificationsQueueRepository
	AlertGroup() repository.AlertGroupRepository

	Close() error
}

type repositories struct {
	notificationsQueue repository.NotificationsQueueRepository
	alertGroup         repository.AlertGroupRepository

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

	appDB, err := bbolt.Open(dataDir+app, 0600, &bbolt.Options{Timeout: 3 * time.Second})

	return &repositories{
		notificationsQueue: repository.NewNotificationsQueueRepository(appDB),
		alertGroup:         repository.NewAlertGroupRepository(appDB),

		db: []*bbolt.DB{appDB},
	}, nil
}

func (r *repositories) NotificationsQueue() repository.NotificationsQueueRepository {
	return r.notificationsQueue
}

func (r *repositories) AlertGroup() repository.AlertGroupRepository {
	return r.alertGroup
}

func (r *repositories) Close() error {
	for _, db := range r.db {
		_ = db.Close()
	}

	return nil
}
