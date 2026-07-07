// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"io"
	"os"
	"pora/internal/dto/requests"
	"pora/internal/dto/responses"
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
	listService ports.ListService
}

// NewUserService создает и возвращает новый объект UserService
func NewUserService(userRepo ports.UserRepository,
	fileService ports.FileService,
	listService ports.ListService) ports.UserService {
	return &UserService{
		userRepo:    userRepo,
		fileService: fileService,
		listService: listService,
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

	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Email != nil {
		user.Email = *req.Email
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

// GetMyInfo получает информацию о текущем пользователе
func (s *UserService) GetMyInfo(userID uuid.UUID) (
	responses.GetMyInfoResponse, error) {

	var response responses.GetMyInfoResponse

	user, err := s.userRepo.FindWithLists(userID)
	if err != nil {
		logger.Log.Error("Ошибка при поиске пользователя: ", err)
		return response, err
	}

	if user == nil {
		logger.Log.Warn("Пользователь не был найден")
		return response, errors.ErrorUserNotFound
	}

	baseURL := os.Getenv("BASE_URL")

	response.UserInfo = responses.UserInfo{
		ID:       userID,
		Name:     user.Name,
		Surname:  user.Surname,
		ImageURL: baseURL + user.ImageURL,
	}

	for _, list := range user.Lists {

		sections, err := s.listService.
			GetHighestPrioritySections(&list)

		if err != nil {
			logger.Log.Warn("Ошибка при получении секций: ", err)
			continue
		}

		listInfo := responses.ListInfo{
			ID:        list.ID,
			Name:      list.Name,
			Sections:  sections,
			CreatedAt: list.CreatedAt,
			UpdatedAt: list.UpdatedAt,
		}

		response.Lists = append(response.Lists, listInfo)
	}

	return response, nil
}
