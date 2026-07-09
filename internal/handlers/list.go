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
	"github.com/google/uuid"
)

// ListHandler - объект, содержащий методы для обработки
// Api запросов, связанных с списками покупок
type ListHandler struct {
	listService  ports.ListService
	tokenService ports.TokenService
}

// NewListHandler создает и возвращает новый объект ListHandler
func NewListHandler(listService ports.ListService,
	tokenService ports.TokenService) *ListHandler {
	return &ListHandler{
		listService:  listService,
		tokenService: tokenService,
	}
}

// CreateList обрабатывает запрос на создание
// нового списка продуктов у пользователя/семьи
func (h *ListHandler) CreateList(c *gin.Context) {

	var req requests.CreateListRequest

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

	response, err := h.listService.CreateList(userID, req.FamilyID, req.Name)

	if err != nil {
		logger.Log.Warn("Возникла ошибка при работе сервиса")
		c.Error(err)
		return
	}

	logger.Log.Info("Создание списка продуктов успешно выполнено")
	c.JSON(http.StatusOK, response)
}

// GetListInfo обрабатывает запрос на получение
// информации об определенном списке продуктов
func (h *ListHandler) GetListInfo(c *gin.Context) {

	rawListID := c.Param("list_id")

	listID, err := uuid.Parse(rawListID)
	if err != nil {
		logger.Log.Warn("Некорректный uuid в запросе: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	response, err := h.listService.GetListInfo(listID)

	if err != nil {
		logger.Log.Warn("Возникла ошибка при работе сервиса")
		c.Error(err)
		return
	}

	logger.Log.Info("Получение информации о списке прошло успешно")
	c.JSON(http.StatusOK, response)
}

// AddItem обрабатывает запрос на добавление
// товара в определенный списко продуктов
func (h *ListHandler) AddItem(c *gin.Context) {

	var req requests.ChangeItemRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Log.Warn("Ошибка при десериализации запроса: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	rawListID := c.Param("list_id")

	listID, err := uuid.Parse(rawListID)
	if err != nil {
		logger.Log.Warn("Некорректный uuid в запросе: ", err)
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

	response, err := h.listService.AddItem(userID, listID, req)

	if err != nil {
		logger.Log.Warn("Возникла ошибка при работе сервиса")
		c.Error(err)
		return
	}

	logger.Log.Info("Добавление товара в список прошло успешно")
	c.JSON(http.StatusOK, response)
}

// DeleteList обрабатывает запрос на удаление списка продуктов
func (h *ListHandler) DeleteList(c *gin.Context) {

	rawListID := c.Param("list_id")

	listID, err := uuid.Parse(rawListID)
	if err != nil {
		logger.Log.Warn("Некорректный uuid в запросе: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	err = h.listService.DeleteList(listID)

	if err != nil {
		logger.Log.Warn("Возникла ошибка при работе сервиса")
		c.Error(err)
		return
	}

	logger.Log.Info("Удаление списка продуктов прошло успешно")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}
