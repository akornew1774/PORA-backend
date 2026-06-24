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

// UserHandler - объект, содержащий методы для обработки
// Api запросов, связанных с пользователями
type UserHandler struct {
	userService  ports.UserService
	tokenService ports.TokenService
}

// NewUserHandler создает и возвращает новый объект UserHandler
func NewUserHandler(userService ports.UserService,
	tokenService ports.TokenService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		tokenService: tokenService,
	}
}

// UpdateUser обрабатывает запрос на обноление
// данных какого-то конкретного пользователя
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	var req requests.UpdateUserRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.Error(errors.ErrorInvalidInput)
		return
	}

	accessToken := c.GetHeader("Authorization")

	userID, err := h.tokenService.DecodeAccessToken(accessToken)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.userService.UpdateUserInfo(userID, &req)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Обновление данных пользователя прошло успешно")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}

// SaveImage обрабатывает запрос на сохранение
// изображения для профиля пользователя
func (h *UserHandler) SaveImage(c *gin.Context) {

	fileHeader, err := c.FormFile("image")
	if err != nil {
		logger.Log.Warn("Ошибка при получении файла: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		logger.Log.Warn("Ошибка при открытии полученного файла: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}
	defer file.Close()

	accessToken := c.GetHeader("Authorization")

	userID, err := h.tokenService.DecodeAccessToken(accessToken)
	if err != nil {
		c.Error(err)
		return
	}

	url, err := h.userService.SaveImage(userID, file, fileHeader.Filename)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Сохранение изображения для профиля успешно выполнено")
	c.JSON(http.StatusOK, responses.SaveImageResponse{ImageURL: url})
}
