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

// SystemHandler - объект, содержащий методы для
// обработки Api запросов на системных эндпойнтах
type SystemHandler struct {
	pushService ports.PushService
}

// NewSystemHandler создает и возвращает новый объект SystemHandler
func NewSystemHandler(pushService ports.PushService) *SystemHandler {
	return &SystemHandler{pushService: pushService}
}

// NotifyEveryone обрабатывает Api запрос на
// уведомление всех пользователей приложения
func (h *SystemHandler) NotifyEveryone(c *gin.Context) {

	var req requests.NotifyEveryoneRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Log.Warn("Ошибка при десериализации запроса: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	err = h.pushService.NotifyEveryone(
		c.Request.Context(),
		req.Title,
		req.Body,
	)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Уведомление пользователей прошло успешно")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}
