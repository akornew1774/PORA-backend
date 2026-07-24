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

// StatisticsHandler - объект, содержащий методы для обработки
// запросов, связанных с предоставлением пользователям статистики
type StatisticsHandler struct {
	statisticsService ports.StatisticsService
	tokenService      ports.TokenService
}

// NewStatisticsHandler создает и возвращает новый объект StatisticsHandler
func NewStatisticsHandler(statisticsService ports.StatisticsService,
	tokenService ports.TokenService) *StatisticsHandler {
	return &StatisticsHandler{
		statisticsService: statisticsService,
		tokenService:      tokenService,
	}
}

// GetUserProducts обрабатывает запрос на получение
// информации о всех товарах, созданных пользователем
func (h *StatisticsHandler) GetUserProducts(c *gin.Context) {

	accessToken := c.GetHeader("Authorization")

	userID, err := h.tokenService.DecodeAccessToken(accessToken)
	if err != nil {
		logger.Log.Warn("Возникла ошибка при декодировании Access-токена")
		c.Error(err)
		return
	}

	response, err := h.statisticsService.GetUserProducts(userID)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Получение информации о товарах пользователя прошло успешно")
	c.JSON(http.StatusOK, response)
}

// GetBrief обрабатывает запрос на получение информации
//
//	о быстро кончающихся продуктах пользователя
func (h *StatisticsHandler) GetBrief(c *gin.Context) {

	accessToken := c.GetHeader("Authorization")

	userID, err := h.tokenService.DecodeAccessToken(accessToken)
	if err != nil {
		logger.Log.Warn("Возникла ошибка при декодировании Access-токена")
		c.Error(err)
		return
	}

	response, err := h.statisticsService.GetBrief(userID)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Получение краткой информации о товарах пользователя прошло успешно")
	c.JSON(http.StatusOK, response)
}

// SaveBrief обрабатывает запрос на сохранение
// быстро кончающихся товаров пользователя
func (h *StatisticsHandler) SaveBrief(c *gin.Context) {

	var req requests.SaveBriefRequest

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

	err = h.statisticsService.SaveBrief(userID, req)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Сохранение краткой информации о товарах пользователя прошло успешно")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}

// GetLoginTimes обрабатывает запрос на получение
// информации о всех входах пользователя в приложение
func (h *StatisticsHandler) GetLoginTimes(c *gin.Context) {

	accessToken := c.GetHeader("Authorization")

	userID, err := h.tokenService.DecodeAccessToken(accessToken)
	if err != nil {
		logger.Log.Warn("Возникла ошибка при декодировании Access-токена")
		c.Error(err)
		return
	}

	response, err := h.statisticsService.GetLoginTimes(userID)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Получение сведений о входах пользователя прощло успешно")
	c.JSON(http.StatusOK, response)
}
