package service

import (
	"context"
	"errors"
	"image"
	"log"

	"github.com/Vladimirmoscow84/Image_processor/internal/model"
)

type imageProcessorRepo interface {
	AddImage(ctx context.Context, image *model.Image) (int, error)
	GetImage(ctx context.Context, id int) (*model.Image, error)
	DeleteImage(ctx context.Context, id int) error
	UpdateImage(ctx context.Context, image *model.Image) error
}

type fileStorageRepo interface {
	Save(ctx context.Context, origPath string) (string, error)
	SaveImage(ctx context.Context, img image.Image, destPath string) (string, error)
	Delete(ctx context.Context, destPath string) error
}

type kafkaProducerConsumer interface {
	Produce(ctx context.Context, msg string) error
	Consume(ctx context.Context, handler func(string) error) error
	Close() error
}

type Service struct {
	db    imageProcessorRepo
	fs    fileStorageRepo
	kafka kafkaProducerConsumer
}

func New(db imageProcessorRepo, fs fileStorageRepo, kafka kafkaProducerConsumer) (*Service, error) {
	if db == nil {
		return nil, errors.New("[service] db client is nil")
	}
	if fs == nil {
		return nil, errors.New("[service] file storage client is nil")
	}
	if kafka == nil {
		log.Println("[service] kafka client is nil, service will be work without queue")
	}
	return &Service{
		db:    db,
		fs:    fs,
		kafka: kafka,
	}, nil
}
