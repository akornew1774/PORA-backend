// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"io"
	"os"
	"path/filepath"
	"pora/internal/errors"
	"pora/internal/ports"
	"slices"
	"strings"

	"github.com/google/uuid"
)

// FileService - объект, содержащий методы для работы с хранилищем
type FileService struct {
	storage ports.Storage
}

// NewFileService создает и возвращает новый объект FileService
func NewFileService(storage ports.Storage) ports.FileService {
	return &FileService{storage: storage}
}

// SaveProfileImage сохраняет изображение профиля пользователя
// в хранилище и возвращает ссылку на него
func (s *FileService) SaveProfileImage(file io.Reader,
	fileName string, userID uuid.UUID) (string, error) {

	ext := filepath.Ext(fileName)
	if !s.isValidExtension(ext) {
		return "", errors.ErrorInvalidInput
	}

	name := userID.String() + ext
	dirName := os.Getenv("PROFILE_IMAGE_DIR")
	basePath := os.Getenv("BASE_PATH")

	path := filepath.Join(dirName, name)

	err := s.storage.Delete(path)
	if err != nil {
		return "", err
	}

	fullPath, err := s.storage.Save(file, path)
	if err != nil {
		return "", err
	}
	if fullPath == "" {
		return "", errors.ErrorInternal
	}

	url := "/" + basePath + "/" + dirName + "/" + name
	return url, nil
}

// isValidExtension проверяет, является ли расширение валидным
func (s *FileService) isValidExtension(ext string) bool {
	validExtensions := []string{".pdf", ".png", ".jpg", ".jpeg", ".gif", ".webp"}
	ext = strings.ToLower(ext)

	return slices.Contains(validExtensions, ext)
}
