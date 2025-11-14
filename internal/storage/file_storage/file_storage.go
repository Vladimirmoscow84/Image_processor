package filestorage

import (
	"fmt"
	"os"
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
}
