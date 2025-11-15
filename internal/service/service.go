package service

import (
	"context"
	"image"

	"github.com/Vladimirmoscow84/Image_processor/internal/model"
)

type imageProcessorRepo interface {
	AddImage(ctx context.Context, image *model.Image) (int, error)
	GetImage(ctx context.Context, id int) (*model.Image, error)
	UpdateImage(ctx context.Context, image *model.Image) error
	DeleteImage(ctx context.Context, id int) error
	GetAllImages(ctx context.Context) ([]model.Image, error)
}

type fileStorageRepo interface {
	Save(ctx context.Context, origPath string) (string, error)
	SaveImage(ctx context.Context, img image.Image, destPath string) (string, error)
	Delete(ctx context.Context, destPath string) error
	Update(ctx context.Context, oldPath, newOrigPath string) (string, error)
}

type Service struct {
	db imageProcessorRepo
	fs fileStorageRepo
}

func New(db imageProcessorRepo, fs fileStorageRepo) *Service {
	return &Service{
		db: db,
		fs: fs,
	}
}
