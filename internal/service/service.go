package service

import (
	"context"

	"github.com/Vladimirmoscow84/Image_processor/internal/model"
)

type imageProcessorRepo interface {
	AddImage(ctx context.Context, image *model.Image) (int, error)
	GetImage(ctx context.Context, id int) (*model.Image, error)
	UpdateImage(ctx context.Context, image *model.Image) error
	DeleteImage(ctx context.Context, id int) error
	GetAllImages(ctx context.Context) ([]model.Image, error)
}

type Service struct {
	storage imageProcessorRepo
}

func New(storage imageProcessorRepo) *Service {
	return &Service{
		storage: storage,
	}
}
