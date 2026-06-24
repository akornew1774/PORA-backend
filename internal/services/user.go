// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"io"
	"os"
	"pora/internal/dto/requests"
	"pora/internal/entities"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"
	"time"

	"github.com/google/uuid"
)

// UserService - объект, содержащий методы для работы с пользователями
type UserService struct {
	userRepo    ports.UserRepository
	fileService ports.FileService
}

// NewUserService создает и возвращает новый объект UserService
func NewUserService(userRepo ports.UserRepository,
	fileService ports.FileService) ports.UserService {
	return &UserService{
		userRepo:    userRepo,
		fileService: fileService,
	}
}

// UpdateUserInfo обновляет (или дополняет) информацию
// о каком-то конкретном пользователе
func (s *UserService) UpdateUserInfo(
	userID uuid.UUID, req *requests.UpdateUserRequest) error {

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске пользователя: ", err)
		return err
	}
	if user == nil {
		logger.Log.Warn("Пользователь с указанным ID не найден")
		return errors.ErrorUserNotFound
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Surname != nil {
		user.Surname = *req.Surname
	}

	user.Status = entities.UserStatusActive
	user.UpdatedAt = time.Now().UTC()

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		logger.Log.Error("Ошибка при обновлении пользователя: ", err)
		return err
	}

	return nil
}

// SaveImage вызывает FileService для сохранения изображения в
// хранилище, сохраняет ссылку на него в профиле пользователя в БД
func (s *UserService) SaveImage(userID uuid.UUID,
	file io.Reader, fileName string) (string, error) {

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске пользователя: ", err)
		return "", err
	}

	if user == nil {
		logger.Log.Warn("Указанный пользователь не найден")
		return "", errors.ErrorUserNotFound
	}

	url, err := s.fileService.SaveProfileImage(file, fileName, userID)
	if err != nil {
		logger.Log.Error("Ошибка при сохранении изображения профиля: ", err)
		return "", err
	}

	user.ImageURL = url
	user.UpdatedAt = time.Now().UTC()

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		logger.Log.Error("Ошибка при обновлении информации о пользователе: ", err)
		return "", err
	}

	baseURL := os.Getenv("BASE_URL")

	return baseURL + url, nil
}
