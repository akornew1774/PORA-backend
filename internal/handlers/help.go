// handlers - пакет, содержащий в себе обработчики Api запросов
package handlers

import (
	"net/http"
	"pora/internal/dto/requests"
	"pora/internal/dto/responses"
	"pora/internal/errors"
	"pora/internal/infrastructure/email"
	"pora/internal/infrastructure/logger"
	"pora/internal/ports"

	"github.com/gin-gonic/gin"
)

// HelpHandler - объект, содержащий методы для обработки
// Api запросов, связанных с помощью пользователям
type HelpHandler struct {
	tokenService ports.TokenService
}

// NewHelpHandler создает и возвращает новый объект HelpHandler
func NewHelpHandler(tokenService ports.TokenService) *HelpHandler {
	return &HelpHandler{tokenService: tokenService}
}

// HandleHelpMessage обрабатывает запрос на передачу
// админам сообщения с просьбой помощи
func (h *HelpHandler) HandleHelpMessage(c *gin.Context) {

	var req requests.HelpMessageRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Log.Warn("Ошибка при десериализации запроса: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	accessToken := c.GetHeader("Authorization")

	userID, err := h.tokenService.DecodeAccessToken(accessToken)
	if err != nil {
		logger.Log.Warn("Возникла ошибка при декодировании Access-токена")
		c.Error(err)
		return
	}

	err = email.SendHelpRequest(req.Title, req.Message, userID)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Обработка просьбы помощи от пользователя прошла успешно")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}
