// services - пакет, содержащий внутреннюю логику приложения
package services

import (
	"encoding/hex"
	"pora/internal/entities"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"
	"strings"

	"crypto/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenService - объект, содержащий методы для создания и обновления токенов
type TokenService struct {
	userRepo         ports.UserRepository
	refreshTokenRepo ports.RefreshTokenRepository
}

// NewTokenService создает и возвращает новый объект TokenService
func NewTokenService(userRepo ports.UserRepository,
	refreshTokenRepo ports.RefreshTokenRepository) ports.TokenService {
	return &TokenService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

// CreateTokens создает новый access-токен,
// обновляет (или создает новый) refresh-токен,
// после чего возвращает access и refresh токены.
func (s *TokenService) CreateTokens(
	user *entities.User) (string, string, error) {

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		logger.Log.Error("Ошибка при генерации refresh-токена: ", err)
		return "", "", err
	}

	err = s.refreshTokenRepo.DeleteByUserID(user.ID)
	if err != nil {
		logger.Log.Error("Ошибка при обращении к БД: ", err)
		return "", "", err
	}

	newRefreshToken := &entities.RefreshToken{
		UserID: user.ID,
		Token:  refreshToken,
	}

	err = s.refreshTokenRepo.CreateRefreshToken(newRefreshToken)
	if err != nil {
		logger.Log.Error("Ошибка при обращении к БД: ", err)
		return "", "", err
	}

	user.RefreshToken = newRefreshToken

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		logger.Log.Error("Ошибка при обновлении сущности пользователя: ", err)
		return "", "", err
	}

	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		logger.Log.Error("Ошибка при генерации access-токена: ", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// GenerateRefreshToken генерирует новый refresh-токен для пользователя
func (s *TokenService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(bytes)
	return token, nil
}

// GenerateAccessToken создает новый access-токен для пользователя
func (s *TokenService) generateAccessToken(
	userID uuid.UUID) (string, error) {

	jwtTokenKey := []byte(os.Getenv("JWT_TOKEN_KEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id":    userID.String(),
			"expires_at": time.Now().UTC().Add(time.Hour).Unix(),
		})

	accessToken, err := token.SignedString(jwtTokenKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// DecodeAccessToken декодирует полученный access-токен,
// проверяет его валидность
func (s *TokenService) DecodeAccessToken(
	rawToken string) (uuid.UUID, error) {

	prefix := "Bearer "

	if after, ok := strings.CutPrefix(rawToken, prefix); ok {
		rawToken = after
	}

	jwtTokenKey := []byte(os.Getenv("JWT_TOKEN_KEY"))

	token, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.ErrorInvalidInput
		}

		return jwtTokenKey, nil
	})
	if err != nil {
		logger.Log.Warn("Некорректный access-токен: ", err)
		return uuid.Nil, errors.ErrorInvalidInput
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		logger.Log.Warn("Некорректный access-токен")
		return uuid.Nil, errors.ErrorInvalidInput
	}

	rawUserID, ok := claims["user_id"].(string)
	if !ok {
		logger.Log.Warn("Access-токен не содержит ID пользователя")
		return uuid.Nil, errors.ErrorInvalidInput
	}
	userID, err := uuid.Parse(rawUserID)
	if err != nil || userID == uuid.Nil {
		logger.Log.Warn("Некорректный ID пользователя в токене: ", err)
		return uuid.Nil, errors.ErrorInvalidInput
	}

	rawExpiresAt, ok := claims["expires_at"].(float64)
	if !ok {
		logger.Log.Warn("Access-токен не содержит срок действия")
		return uuid.Nil, errors.ErrorInvalidInput
	}
	expiresAt := time.Unix(int64(rawExpiresAt), 0)

	if time.Now().UTC().After(expiresAt) {
		logger.Log.Warn("Время действия токена истекло")
		return uuid.Nil, errors.ErrorAccessExpired
	}

	return userID, nil
}
