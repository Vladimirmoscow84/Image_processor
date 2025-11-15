package filestorage

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type FileStorage struct {
	Path      string
	watermark image.Image
}

// New - конструктор соединения с файловым хранилищем
func New(path, watermarkPath string) (*FileStorage, error) {

	if path == "" {
		return nil, fmt.Errorf("[fileStorage] base path is empty")
	}

	fs := &FileStorage{
		Path: path,
	}

	if watermarkPath != "" {
		wm, err := imaging.Open(watermarkPath)
		if err != nil {
			return nil, fmt.Errorf("[fileStorage] failed to load watermark: %w", err)
		}
		fs.watermark = wm
	}

	err := os.MkdirAll(path, 0755)
	if err != nil {
		return nil, fmt.Errorf("[fileStorage] failed to create base dir: %w", err)
	}

	return fs, nil
}

// Save сохраняет файл в локальное хранилище и возвращает путь к нему
func (f *FileStorage) Save(ctx context.Context, origPath string) (string, error) {

	filename := filepath.Base(origPath)
	destPath := filepath.Join(f.Path, "originals", filename)

	err := os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		return "", fmt.Errorf("[filestorage] failed to create dir: %w", err)
	}

	data, err := os.ReadFile(origPath)
	if err != nil {
		return "", fmt.Errorf("[filestorage] failed to read source file: %w", err)
	}

	err = os.WriteFile(destPath, data, 0644)
	if err != nil {
		return "", fmt.Errorf("[filestorage]failed to write file: %w", err)
	}

	return destPath, nil
}

// SaveImage сохраняет image.Image в локальное хранилище с водяным знаком и нужным форматом
func (f *FileStorage) SaveImage(ctx context.Context, img image.Image, destPath string) (string, error) {
	fullPath := filepath.Join(f.Path, destPath)

	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		return "", fmt.Errorf("[filestorage] failed to create directories: %w", err)
	}

	// Добавление водянго знака если он задан
	if f.watermark != nil {
		img = applyWatermark(img, f.watermark)
	}

	outFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("[filestorage] failed to create file: %w", err)
	}
	defer outFile.Close()

	ext := strings.ToLower(filepath.Ext(fullPath))
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	case ".png":
		err = png.Encode(outFile, img)
	case ".gif":
		err = gif.Encode(outFile, img, nil)
	default:
		// по умолчанию JPG
		err = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 90})
	}
	if err != nil {
		return "", fmt.Errorf("[filestorage] failed to encode image: %w", err)
	}

	return fullPath, nil
}

// Delete удаляет файл из локального хранилища
func (f *FileStorage) Delete(ctx context.Context, destPath string) error {
	if destPath == "" {
		log.Println("[fileStorage] base path is empty")
		return nil
	}
	err := os.Remove(destPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("[fileStorage] failed to delete file: %w", err)
	}
	return nil
}

// Update обновляет (перезаписывает) старый файл на новый
func (f *FileStorage) Update(ctx context.Context, oldPath, newOrigPath string) (string, error) {
	err := f.Delete(ctx, oldPath)
	if err != nil {
		return "", err
	}

	return f.Save(ctx, newOrigPath)
}

// ResizeImage изменяет размер изображения
func ResizeImage(img image.Image, width, height int) image.Image {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

// CreateThumbnail создает миниатюру
func CreateThumbnail(img image.Image, width, height int) image.Image {
	return imaging.Thumbnail(img, width, height, imaging.Lanczos)
}

// applyWatermark накладывает водяной знак на изображение
func applyWatermark(base, watermark image.Image) image.Image {
	offset := image.Pt(base.Bounds().Dx()-watermark.Bounds().Dx()-10, base.Bounds().Dy()-watermark.Bounds().Dy()-10)
	result := imaging.Clone(base)
	draw.Draw(result, watermark.Bounds().Add(offset), watermark, image.Point{}, draw.Over)
	return result
}
