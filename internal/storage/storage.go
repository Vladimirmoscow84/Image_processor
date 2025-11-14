package storage

import (
	"fmt"
	"log"

	filestorage "github.com/Vladimirmoscow84/Image_processor/internal/storage/file_storage"
	"github.com/Vladimirmoscow84/Image_processor/internal/storage/postgres"
)

type Storage struct {
	*postgres.Postgres
	*filestorage.FileStorage
}

// New - конструктор storage
func New(pg *postgres.Postgres, fs *filestorage.FileStorage) (*Storage, error) {
	if pg == nil {
		log.Println("[storage] postgres client is nil")
		return nil, fmt.Errorf("[storage] postgres client is nil")
	}
	if fs == nil {
		log.Println("[storage] fileStorage client is nil")
		return nil, fmt.Errorf("[storage] fileStorage client is nil")
	}

	return &Storage{
		Postgres:    pg,
		FileStorage: fs,
	}, nil
}
