package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Vladimirmoscow84/Image_processor/internal/model"
)

type ImageProcessorService interface {
	Upload(ctx context.Context, srcPath string) (*model.Image, error)
	Get(ctx context.Context, id int) (*model.Image, error)
	Delete(ctx context.Context, id int) error
	UpdateImageFile(ctx context.Context, id int, newSrcPath string) (*model.Image, error)
	GetAll(ctx context.Context) ([]model.Image, error)
}

func (s *Service) Upload(ctx context.Context, srcPath string) (*model.Image, error) {

	origSavedPath, err := s.fs.Save(ctx, srcPath)
	if err != nil {
		return nil, fmt.Errorf("[service] failed to save original: %w", err)
	}

	// 2. Генерация thumbnail (заглушка — потом можно добавить реальную обработку)
	thumbPath := origSavedPath // временно используем тот же путь или генерируем рядом

	// 3. Сохраняем запись в БД
	img := &model.Image{
		OriginalPath:  origSavedPath,
		ProcessedPath: "",        // можно позже добавить обработку
		ThumbnailPath: thumbPath, // временно
		Status:        "uploaded",
	}

	id, err := s.db.AddImage(ctx, img)
	if err != nil {
		// если запись в БД провалилась — удаляем файл из локального хранилища
		_ = s.fs.Delete(ctx, origSavedPath)
		return nil, fmt.Errorf("[service] failed to insert image to DB: %w", err)
	}

	img.ID = id
	return img, nil
}

func (s *Service) Get(ctx context.Context, id int) (*model.Image, error) {
	return s.db.GetImage(ctx, id)
}

func (s *Service) GetAll(ctx context.Context) ([]model.Image, error) {
	return s.db.GetAllImages(ctx)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	// 1. получаем запись, чтобы знать путь к файлам
	img, err := s.db.GetImage(ctx, id)
	if err != nil {
		return err
	}

	// 2. удаляем файл с диска
	if img.OriginalPath != "" {
		if err := s.fs.Delete(ctx, img.OriginalPath); err != nil {
			log.Printf("[service] failed to delete file %s: %v", img.OriginalPath, err)
		}
	}
	if img.ThumbnailPath != "" {
		if err := s.fs.Delete(ctx, img.ThumbnailPath); err != nil {
			log.Printf("[service] failed to delete thumb %s: %v", img.ThumbnailPath, err)
		}
	}

	// 3. удаляем запись из БД
	return s.db.DeleteImage(ctx, id)
}

func (s *Service) UpdateImageFile(ctx context.Context, id int, newSrcPath string) (*model.Image, error) {
	// 1. получаем текущую запись
	img, err := s.db.GetImage(ctx, id)
	if err != nil {
		return nil, err
	}

	// 2. заменяем файл в хранилище
	newStoredPath, err := s.fs.Update(ctx, img.OriginalPath, newSrcPath)
	if err != nil {
		return nil, fmt.Errorf("[service] failed to replace image: %w", err)
	}

	// 3. можно пересоздать thumbnail
	img.ThumbnailPath = newStoredPath // заглушка

	// 4. обновляем запись в БД
	img.OriginalPath = newStoredPath
	img.Status = "updated"

	if err := s.db.UpdateImage(ctx, img); err != nil {
		return nil, err
	}

	return img, nil
}
