// storage - пакет, отвечающий за сохранение и передачу изображений
package storage

import (
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"

	"io"
	"os"
	"path/filepath"
)

// LocalStorage - структура для сохранения и удаления
// изображений локально на сервере
type LocalStorage struct {
	BasePath string
}

// NewLocalStorage создает и возвращает новый объект LocalStorage
func NewLocalStorage(basePath string) ports.Storage {
	return &LocalStorage{BasePath: basePath}
}

// Save сохраняет переданное изображение
// по заданному пути в хранилище
func (s *LocalStorage) Save(file io.Reader, path string) (string, error) {
	fullPath := filepath.Join(s.BasePath, path)

	err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
	if err != nil {
		logger.Log.Error("Ошибка при создании директории")
		return "", err
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		logger.Log.Error("Ошибка при создании файла для изображения")
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		logger.Log.Error("Ошибка при копировании изображения в файл")
		return "", err
	}

	return "/" + s.BasePath + "/" + path, nil
}

// Delete удаляет изображение в хранилище по заданному пути
func (s *LocalStorage) Delete(path string) error {
	fullPath := filepath.Join(s.BasePath, path)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		logger.Log.Warn("Файл по указанному пути не найден")
		return nil
	}

	return os.Remove(fullPath)
}
