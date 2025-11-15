package service

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Vladimirmoscow84/Image_processor/internal/model"
	"github.com/disintegration/imaging"
)

type ImageProcessorService interface {
	ProcessAndSaveImage(ctx context.Context, origPath string) (*model.Image, error)
	DeleteImage(ctx context.Context, image *model.Image) error
	UpdateImage(ctx context.Context, image *model.Image, newOrigPath string) (*model.Image, error)
	ProcessBatch(ctx context.Context, origPaths []string) ([]*model.Image, error)
}

// ProcessAndSaveImage обрабатывает изображение и сохраняет все версии (original, processed, thumbnail)
// возвращает объект модели с заполненными путями и статусом
func (s *Service) ProcessAndSaveImage(ctx context.Context, origPath string) (*model.Image, error) {

	path, err := s.fs.Save(ctx, origPath)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to save original: %w", err)
	}

	img, err := imaging.Open(path)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to open image: %w", err)
	}

	processedImg := imaging.Resize(img, 1280, 0, imaging.Lanczos)
	processedPath := filepath.Join("processed", filepath.Base(origPath))
	processedPathSaved, err := s.fs.SaveImage(ctx, processedImg, processedPath)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to save processed image: %w", err)
	}

	thumbImg := imaging.Thumbnail(img, 300, 300, imaging.Lanczos)
	thumbPath := filepath.Join("thumbs", filepath.Base(origPath))
	thumbPathSaved, err := s.fs.SaveImage(ctx, thumbImg, thumbPath)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to save thumbnail: %w", err)
	}

	imageModel := &model.Image{
		OriginalPath:  origPath,
		ProcessedPath: processedPathSaved,
		ThumbnailPath: thumbPathSaved,
		Status:        "processed",
	}

	id, err := s.db.AddImage(ctx, imageModel)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to save image record in DB: %w", err)
	}
	imageModel.ID = id

	return imageModel, nil
}

//
