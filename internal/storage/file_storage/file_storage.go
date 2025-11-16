package filestorage

import (
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

type FileStorage struct {
	Path      string
	watermark image.Image
}

// New - конструктор файлового хранилища
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

// Save сохраняет файл в локальное хранилище без обработки и возвращает путь к нему
func (f *FileStorage) Save(ctx context.Context, origPath string) (string, error) {

	filename := filepath.Base(origPath)
	destPath := filepath.Join(f.Path, "originals", filename)

	err := os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		return "", fmt.Errorf("[fileStorage] failed to create subdir: %w", err)
	}

	input, err := os.Open(origPath)
	if err != nil {
		return "", fmt.Errorf("[fileStorage] failed to open source: %w", err)
	}
	defer input.Close()

	output, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("[fileStorage] failed to create file: %w", err)
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	if err != nil {
		return "", fmt.Errorf("[fileStorage] failed to save file: %w", err)
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
		return nil
	}
	fullPath := destPath
	if !filepath.IsAbs(destPath) {
		fullPath = filepath.Join(f.Path, destPath)
	}

	err := os.Remove(fullPath)

	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("[fileStorage] failed to delete file: %w", err)
	}
	return nil
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
	result := imaging.Clone(base)

	offset := image.Pt(
		base.Bounds().Dx()-watermark.Bounds().Dx()-10,
		base.Bounds().Dy()-watermark.Bounds().Dy()-10,
	)

	draw.Draw(
		result,
		watermark.Bounds().Add(offset),
		watermark,
		image.Point{},
		draw.Over,
	)

	return result
}
