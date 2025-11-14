package filestorage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type FileStorage struct {
	Path string
}

// New - конструктор соединения с файловым хранилищем
func New(path string) (*FileStorage, error) {
	if path == "" {
		return nil, fmt.Errorf("[fileStorage] base path is empty")
	}

	err := os.Mkdir(path, 0755)
	if err != nil {
		return nil, fmt.Errorf("[fileStorage] failed to create base dir: %w", err)
	}
	return &FileStorage{
		Path: path,
	}, nil
}

// Save сохраняет файл в локальное хранилище и возвращает путь к нему
func (f *FileStorage) Save(ctx context.Context, origPath string) (string, error) {

	//создание уникального имени файла
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(origPath))
	destPath := filepath.Join(f.Path, fileName)

	//копироание файла
	in, err := os.Open(origPath)
	if err != nil {
		return "", fmt.Errorf("[fileStorage] failed to open file: %w", err)
	}
	defer in.Close()

	out, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("[fileStorage] failed to create file %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return "", fmt.Errorf("[fileStorage] failed to copy file: %w", err)
	}
	return destPath, nil
}

// Delete удаляет файл из локального хранилища
func (f *FileStorage) Delete(ctx context.Context, destPath string) error {
	return os.Remove(destPath)
}

// Update обновляет (перезаписывает) файл старый на новый
func (f *FileStorage) Update(ctx context.Context, destPath, newOrigPath string) (string, error) {
	_ = os.Remove(destPath)

	return f.Save(ctx, newOrigPath)
}
