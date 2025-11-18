package service

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/Vladimirmoscow84/Image_processor/internal/model"
	"github.com/disintegration/imaging"
)

type ImageProcessorService interface {
	ProcessAndSaveImage(ctx context.Context, origPath string) (*model.Image, error)
	DeleteImage(ctx context.Context, image *model.Image) error
	EnqueueImage(ctx context.Context, origPath string) error
	StartKafkaConsumer(ctx context.Context)
}

// ProcessAndSaveImage обрабатывает изображение и сохраняет все версии (original, processed, thumbnail)
// возвращает объект модели с заполненными путями и статусом
func (s *Service) ProcessAndSaveImage(ctx context.Context, origPath string) (*model.Image, error) {

	path, err := s.fs.Save(ctx, origPath)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to save original: %w", err)
	}

	processedPathSaved, thumbPathSaved, err := s.createProcessedVersions(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to create processed/thumb: %w", err)
	}

	var img *model.Image
	images, _ := s.db.GetAllImages(ctx)
	for _, im := range images {
		if im.OriginalPath == path {
			img = im
			break
		}
	}

	if img == nil {
		img = &model.Image{
			OriginalPath:  path,
			ProcessedPath: processedPathSaved,
			ThumbnailPath: thumbPathSaved,
			Status:        "processed",
		}
		id, err := s.db.AddImage(ctx, img)
		if err != nil {
			return nil, fmt.Errorf("[imageprocessor] failed to add image record: %w", err)
		}
		img.ID = id
	} else {
		img.ProcessedPath = processedPathSaved
		img.ThumbnailPath = thumbPathSaved
		img.Status = "processed"
		err = s.db.UpdateImage(ctx, img)
		if err != nil {
			return nil, fmt.Errorf("[imageprocessor] failed to update image record: %w", err)
		}
	}

	return img, nil
}

// createProcessedVersions создаёт processed и thumbnail версии изображения
func (s *Service) createProcessedVersions(ctx context.Context, origPath string) (string, string, error) {
	img, err := imaging.Open(origPath)
	if err != nil {
		return "", "", fmt.Errorf("[imageprocessor] failed to open image: %w", err)
	}

	processedImg := imaging.Resize(img, 1280, 0, imaging.Lanczos)
	processedPath := filepath.Join("processed", filepath.Base(origPath))
	processedSaved, err := s.fs.SaveImage(ctx, processedImg, processedPath)
	if err != nil {
		return "", "", fmt.Errorf("[imageprocessor] failed to save processed image: %w", err)
	}

	thumbImg := imaging.Resize(img, 300, 0, imaging.Lanczos)
	thumbPath := filepath.Join("thumbs", filepath.Base(origPath))
	thumbSaved, err := s.fs.SaveImage(ctx, thumbImg, thumbPath)
	if err != nil {
		return "", "", fmt.Errorf("[imageprocessor] failed to save thumbnail: %w", err)
	}

	return processedSaved, thumbSaved, nil
}

// DeleteImage удаляет все версии изображения и запись из БД
func (s *Service) DeleteImage(ctx context.Context, image *model.Image) error {
	if err := s.fs.Delete(ctx, image.OriginalPath); err != nil {
		return fmt.Errorf("[imageprocessor] failed to delete original: %w", err)
	}
	if err := s.fs.Delete(ctx, image.ProcessedPath); err != nil {
		return fmt.Errorf("[imageprocessor] failed to delete processed: %w", err)
	}
	if err := s.fs.Delete(ctx, image.ThumbnailPath); err != nil {
		return fmt.Errorf("[imageprocessor] failed to delete thumbnail: %w", err)
	}
	if err := s.db.DeleteImage(ctx, image.ID); err != nil {
		return fmt.Errorf("[imageprocessor] failed to delete DB record: %w", err)
	}
	return nil
}

// EnqueueImage отправляет ID изображения в Kafka
func (s *Service) EnqueueImage(ctx context.Context, imageID int) error {
	if s.kafka == nil {
		return fmt.Errorf("[imageprocessor] kafka client is nil")
	}
	return s.kafka.Produce(ctx, strconv.Itoa(imageID))
}

// StartKafkaConsumer запускает фоновый воркер для обработки очереди
func (s *Service) StartKafkaConsumer(ctx context.Context) {
	if s.kafka == nil {
		log.Println("[imageprocessor] kafka is nil, consumer not ready")
		return
	}

	go func() {
		err := s.kafka.Consume(ctx, func(msg string) error {
			id, err := strconv.Atoi(msg)
			if err != nil {
				log.Printf("[worker] invalid image ID: %s", msg)
				return nil
			}

			img, err := s.db.GetImage(ctx, id)
			if err != nil {
				log.Printf("[worker] image with id=%d not found", id)
				return nil
			}

			_, err = s.ProcessAndSaveImage(ctx, img.OriginalPath)
			if err != nil {
				log.Printf("[worker] failed to process %d: %v", id, err)
			} else {
				log.Printf("[worker] successfully processed %d", id)
			}
			return nil
		})
		if err != nil && ctx.Err() == nil {
			log.Printf("[worker] consumer error: %v", err)
		}
	}()
}

// GetImage возвращает изображение по ID
func (s *Service) GetImage(ctx context.Context, id int) (*model.Image, error) {
	return s.db.GetImage(ctx, id)
}

// AddImage добавляет новую запись
func (s *Service) AddImage(ctx context.Context, img *model.Image) (int, error) {
	return s.db.AddImage(ctx, img)
}

// UpdateImage обновляет запись
func (s *Service) UpdateImage(ctx context.Context, img *model.Image) error {
	return s.db.UpdateImage(ctx, img)
}

// GetAllImages возвращает все изображения
func (s *Service) GetAllImages(ctx context.Context) ([]*model.Image, error) {
	return s.db.GetAllImages(ctx)
}
