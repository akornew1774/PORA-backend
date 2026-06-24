// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"pora/internal/dto/responses"
	"pora/internal/entities"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"

	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
)

// AuthService - объект, содержащий методы для авторизации пользователей
type AuthService struct {
	userRepo         ports.UserRepository
	refreshTokenRepo ports.RefreshTokenRepository
	TokenService     ports.TokenService
}

// NewAuthService создает и возвращает новый объект AuthService
func NewAuthService(userRepo ports.UserRepository,
	refreshTokenRepo ports.RefreshTokenRepository,
	TokenService ports.TokenService) ports.AuthService {
	return &AuthService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		TokenService:     TokenService,
	}
}

// CheckUser находит пользователя по телефону и проверяет его
// статус (notFound / notRegistered / registered)
func (s *AuthService) CheckUser(rawPhone string) (
	responses.IsUserResponse, error) {

	num, err := phonenumbers.Parse(rawPhone, "RU")
	if err != nil || !phonenumbers.IsValidNumber(num) {
		return responses.IsUserResponse{},
			errors.ErrorInvalidPhone
	}

	phone := phonenumbers.Format(num, phonenumbers.E164)

	status := string(entities.UserStatusNotFound)

	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return responses.IsUserResponse{}, err
	}

	if user == nil {
		return responses.IsUserResponse{Status: status},
			nil
	}

	response := responses.IsUserResponse{
		UserID: user.ID,
		Status: string(user.Status),
	}

	return response, nil
}

// GetNewTokens получает новые access и refresh токены для пользователя
func (s *AuthService) GetNewTokens(
	refreshTokenStr string) (responses.RefreshResponse, error) {

	refreshToken, err := s.refreshTokenRepo.FindByToken(refreshTokenStr)
	if err != nil {
		logger.Log.Error("Ошибка при обращении к БД: ", err)
		return responses.RefreshResponse{}, err
	}
	if refreshToken == nil {
		logger.Log.Error("Указанный refresh-токен не найден")
		return responses.RefreshResponse{}, errors.ErrorUserNotFound
	}

	user, err := s.userRepo.FindByID(refreshToken.UserID)
	if err != nil {
		logger.Log.Error("Ошибка при обращении к БД: ", err)
		return responses.RefreshResponse{}, err
	}
	if user == nil {
		logger.Log.Error("Пользователь не был найден")
		return responses.RefreshResponse{}, errors.ErrorUserNotFound
	}

	accessToken, refreshTokenStr, err := s.TokenService.CreateTokens(user)
	if err != nil {
		logger.Log.Error("Ошибка при обновлении токенов пользователя: ", err)
		return responses.RefreshResponse{}, err
	}

	response := responses.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
	}
	return response, err
}

// RemoveRefreshToken удаляет refresh-токен по ID пользователя
func (s *AuthService) RemoveRefreshToken(userID uuid.UUID) error {

	err := s.refreshTokenRepo.DeleteByUserID(userID)
	if err != nil {
		logger.Log.Error("Ошибка при удалении Refresh-токена: ", err)
		return err
	}

	return nil
}
