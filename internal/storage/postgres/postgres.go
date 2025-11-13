package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Vladimirmoscow84/Image_processor/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	DB *sqlx.DB
}

// New - контсруткор соединения к БД
func New(databaseURI string) (*Postgres, error) {
	db, err := sqlx.Connect("pgx", databaseURI)
	if err != nil {
		return nil, fmt.Errorf("[postgres] failed to connect to DB: %w ", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(20 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("[postgres] ping failed: %w ", err)
	}

	log.Println("[postgres] successfull connect to DB")
	return &Postgres{
		DB: db,
	}, nil
}

// Close закрывает соединение с БД
func (p *Postgres) Close() error {
	if p.DB != nil {
		log.Println("[postgres] closing connection on DB")
		return p.DB.Close()
	}
	return nil
}

// AddImage добавляет новую запись в таблицу и возвращает id фронту
func (p *Postgres) AddImage(ctx context.Context, image *model.Image) (int, error) {
	row := p.DB.QueryRowContext(ctx, `
	INSERT INTO images
		(original_path, processed_path, thumbnail_path, status)
	VALUES
		($1,$2,$3,$4)
		RETURNING id;
	`, image.OriginalPath, image.ProcessedPath, image.ThumbnailPath, image.Status)

	var id int
	err := row.Scan(&id)
	if err != nil {
		log.Printf("[postgres] error adding to base: %v", err)
		return 0, fmt.Errorf("[postgres] error adding to base: %w", err)
	}
	image.ID = id
	return id, nil
}
