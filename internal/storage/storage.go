package storage

import (
	"errors"
	"log"

	"github.com/Vladimirmoscow84/Image_processor/internal/storage/postgres"
)

type Storage struct {
	*postgres.Postgres
}

// New - конструктор storage
func New(pg *postgres.Postgres) (*Storage, error) {
	if pg == nil {
		log.Println("[storage] postgres client is nil")
		return nil, errors.New("[storage] postgres client is nil")
	}
	return &Storage{
		Postgres: pg,
	}, nil
}
