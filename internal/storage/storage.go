package storage

import (
	"fmt"
	"github.com/IKostarev/yandex-go-dev/internal/config"
	"github.com/IKostarev/yandex-go-dev/internal/model"
	"github.com/IKostarev/yandex-go-dev/internal/storage/database/postgres"
	"github.com/IKostarev/yandex-go-dev/internal/storage/fs"
	"github.com/IKostarev/yandex-go-dev/internal/storage/mem"
	"github.com/google/uuid"
)

type Storage interface {
	Save(string, string, uuid.UUID) (string, error)
	Get(string, string, uuid.UUID) (string, string)
	GetUserLinks(uuid.UUID) ([]model.UserLink, error)
	CheckIsURLExists(string) (string, error)
	Ping() bool
	Close() error
}

func NewStorage(cfg config.Config) (Storage, error) {
	var s Storage
	var err error

	if cfg.DatabaseDSN != "" {
		if s, err = postgres.NewPostgresDB(cfg.DatabaseDSN); err != nil {
			return nil, fmt.Errorf("cannot database storage: %w", err)
		}
	} else if cfg.FileStoragePath != "" {
		if s, err = fs.NewFsFromFile(cfg.FileStoragePath); err != nil {
			return nil, fmt.Errorf("error NewFs file: %w", err)
		}
	} else {
		if s, err = mem.NewMem(); err != nil {
			return nil, fmt.Errorf("error NewMem: %w", err)
		}
	}

	return s, nil
}
