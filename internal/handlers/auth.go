// handlers - пакет, содержащий в себе обработчики Api запросов
package handlers

import (
	"net/http"
	"pora/internal/dto/requests"
	"pora/internal/dto/responses"
	"pora/internal/errors"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"

	"github.com/gin-gonic/gin"
)

// AuthHandler - объект, содержащий методы для обработки
// Api запросов, связанных с авторизацией
type AuthHandler struct {
	authService  ports.AuthService
	tokenService ports.TokenService
}

// NewAuthHandler создает и возвращает новый объект AuthHandler
func NewAuthHandler(authService ports.AuthService,
	tokenService ports.TokenService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		tokenService: tokenService,
	}
}

// CheckUser обрабатывает запрос проверки статуса пользователя.
// Статусы: notFound, notRegistered, registered
func (h *AuthHandler) CheckUser(c *gin.Context) {

	var req requests.PhoneAndEmailRequest

	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Log.Warn("Ошибка при десериализации запроса: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	response, err := h.authService.CheckUser(req.Phone, req.Email)

	if err != nil {
		logger.Log.Warn("Возникла ошибка при работе сервиса")
		c.Error(err)
		return
	}

	logger.Log.Info("Проверка статуса пользователя успешно выполнена")
	c.JSON(http.StatusOK, response)
}

// RefreshTokens обрабатывает запрос обновления токенов
func (h *AuthHandler) RefreshTokens(c *gin.Context) {

	var req requests.RefreshRequest

	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Log.Warn("Ошибка при десериализации запроса: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	response, err := h.authService.GetNewTokens(req.RefreshToken)

	if err != nil {
		logger.Log.Warn("Возникла ошибка при работе сервиса")
		c.Error(err)
		return
	}

	logger.Log.Info("Обновление токенов успешно выполнено")
	c.JSON(http.StatusOK, response)
}

// Logout обрабатывает запрос выхода пользователя из аккаунта
func (h *AuthHandler) Logout(c *gin.Context) {

	accessToken := c.GetHeader("Authorization")

	userID, err := h.tokenService.DecodeAccessToken(accessToken)
	if err != nil {
		logger.Log.Warn("Возникла ошибка при декодировании Access-токена")
		c.Error(err)
		return
	}

	err = h.authService.RemoveRefreshToken(userID)

	if err != nil {
		logger.Log.Warn("Возникла ошибка при работе сервиса")
		c.Error(err)
		return
	}

	logger.Log.Info("Логаут успешно выполнен")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}
