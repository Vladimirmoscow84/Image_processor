package service

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

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
		return nil, fmt.Errorf("[imageprocessor] failed to create processed and thumb versions: %w", err)
	}

	imageModel := &model.Image{
		OriginalPath:  path,
		ProcessedPath: processedPathSaved,
		ThumbnailPath: thumbPathSaved,
		Status:        "processed",
	}

	_, err = s.db.AddImage(ctx, imageModel)
	if err != nil {
		return nil, fmt.Errorf("[imageprocessor] failed to save image record in DB: %w", err)
	}

	return imageModel, nil
}

// createProcessedVersions создаёт processed и thumbnail версии изображения
func (s *Service) createProcessedVersions(ctx context.Context, origPath string) (string, string, error) {
	img, err := imaging.Open(origPath)
	if err != nil {
		return "", "", fmt.Errorf("[imageprocessor] failed to open image: %w", err)
	}
	//resize
	processedImg := imaging.Resize(img, 1280, 0, imaging.Lanczos)
	processedPath := filepath.Join("processed", filepath.Base(origPath))
	processedSaved, err := s.fs.SaveImage(ctx, processedImg, processedPath)
	if err != nil {
		return "", "", fmt.Errorf("[imageprocessor] failed to save processed image: %w", err)
	}
	//thumbnail
	thumbImg := imaging.Thumbnail(img, 300, 300, imaging.Lanczos)
	thumbPath := filepath.Join("thumbs", filepath.Base(origPath))
	thumbSaved, err := s.fs.SaveImage(ctx, thumbImg, thumbPath)
	if err != nil {
		return "", "", fmt.Errorf("[imageprocessor] failed to save thumbnail: %w", err)
	}

	return processedSaved, thumbSaved, nil
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

// EnqueueImage отправляет путь изображения в Kafka
func (s *Service) EnqueueImage(ctx context.Context, origPath string) error {
	if s.kafka == nil {
		return fmt.Errorf("[imageprocessor] kafka client is nil")
	}
	return s.kafka.Produce(ctx, origPath)
}

// StartKafkaConsumer запускает обработку очереди Kafka
func (s *Service) StartKafkaConsumer(ctx context.Context) {
	if s.kafka == nil {
		log.Println("[imageprocessor] kafka is nil, consumer not ready to work")
		return
	}
	go func() {
		err := s.kafka.Consume(ctx, func(msg string) error {
			_, err := s.ProcessAndSaveImage(ctx, msg)
			if err != nil {
				log.Printf("[worker] failed to process %s: %v", msg, err)
			} else {
				log.Printf("[worker] successfully processed %s", msg)
			}
			return nil
		})
		if err != nil && ctx.Err() == nil {
			log.Printf("[worker] consumer error: %v", err)
		}
	}()
}

func (s *Service) GetImage(ctx context.Context, id int) (*model.Image, error) {
	return s.db.GetImage(ctx, id)
}

func (s *Service) AddImage(ctx context.Context, img *model.Image) (int, error) {
	return s.db.AddImage(ctx, img)
}
