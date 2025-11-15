package service

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

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

// DeleteImage удаляет все версии изображения (original, processed, thumbnail) и запись в БД
func (s *Service) DeleteImage(ctx context.Context, image *model.Image) error {
	//удаление файлов
	err := s.fs.Delete(ctx, image.OriginalPath)
	if err != nil {
		return fmt.Errorf("[imageprocessor] failed to delete original: %w", err)
	}

	err = s.fs.Delete(ctx, image.ProcessedPath)
	if err != nil {
		return fmt.Errorf("[imageprocessor] failed to delete processed: %w", err)
	}

	err = s.fs.Delete(ctx, image.ThumbnailPath)
	if err != nil {
		return fmt.Errorf("[imageprocessor] failed to delete thumbnail: %w", err)
	}

	// удаление записи из БД
	err = s.db.DeleteImage(ctx, image.ID)
	if err != nil {
		return fmt.Errorf("[imageprocessor] failed to delete DB record: %w", err)
	}

	return nil
}

// UpdateImage заменяет оригинальный файл и пересоздаёт processed и thumbnail
func (s *Service) UpdateImage(ctx context.Context, image *model.Image, newOrigPath string) (*model.Image, error) {
	// удаление старых processed и thumbnail
	err := s.fs.Delete(ctx, image.ProcessedPath)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to delete old processed: %w", err)
	}
	err = s.fs.Delete(ctx, image.ThumbnailPath)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to delete old thumbnail: %w", err)
	}

	// сохранение нового оригинального файла
	newPath, err := s.fs.Update(ctx, image.OriginalPath, newOrigPath)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to update original: %w", err)
	}

	// загрузка нового файла
	img, err := imaging.Open(newPath)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to open new file: %w", err)
	}

	// создание новых processed и thumbnail
	processedImg := imaging.Resize(img, 1280, 0, imaging.Lanczos)
	processedPath := filepath.Join("processed", filepath.Base(newOrigPath))
	processedSaved, err := s.fs.SaveImage(ctx, processedImg, processedPath)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to save processed image: %w", err)
	}

	thumbImg := imaging.Thumbnail(img, 300, 300, imaging.Lanczos)
	thumbPath := filepath.Join("thumbs", filepath.Base(newOrigPath))
	thumbSaved, err := s.fs.SaveImage(ctx, thumbImg, thumbPath)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to save thumbnail: %w", err)
	}

	// оновление модели на основании новых файлов
	image.OriginalPath = newPath
	image.ProcessedPath = processedSaved
	image.ThumbnailPath = thumbSaved

	// сохранение новой модели в БД
	err = s.db.UpdateImage(ctx, image)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to update DB record: %w", err)
	}

	return image, nil
}

// ProcessBatch обрабатывает массив изображений параллельно и возвращает результаты
func (s *Service) ProcessBatch(ctx context.Context, origPaths []string) ([]*model.Image, error) {
	var wg sync.WaitGroup
	results := make([]*model.Image, len(origPaths))

	for i, path := range origPaths {
		wg.Add(1)
		go func(idx int, p string) {
			defer wg.Done()
			imgModel, _ := s.ProcessAndSaveImage(ctx, p)
			results[idx] = imgModel
		}(i, path)
	}

	wg.Wait()
	return results, nil
}
