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

// ItemHandler - объект, содержащий методы для обработки
// Api запросов, связанных с товарами в списках продуктов
type ItemHandler struct {
	itemService  ports.ItemService
	tokenService ports.TokenService
}

// NewItemHandler создает и возвращает новый объект ItemHandler
func NewItemHandler(itemService ports.ItemService,
	tokenService ports.TokenService) *ItemHandler {
	return &ItemHandler{
		itemService:  itemService,
		tokenService: tokenService,
	}
}

// GetItemInfo обрабатывает запрос на получение
// информации об одном конкретном товаре
func (h *ItemHandler) GetItemInfo(c *gin.Context) {

	rawItemID := c.Param("item_id")

	itemID, err := uuid.Parse(rawItemID)
	if err != nil {
		logger.Log.Warn("Некорректный uuid в запросе: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	response, err := h.itemService.GetItemInfo(itemID)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Получение информации о товаре успешно выполнено")
	c.JSON(http.StatusOK, response)
}

// ChangeItem обрабатывает запрос на изменение
// определенного товара в списке покупок
func (h *ItemHandler) ChangeItem(c *gin.Context) {

	var req requests.ChangeItemRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Log.Warn("Ошибка при десериализации запроса: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	rawItemID := c.Param("item_id")

	itemID, err := uuid.Parse(rawItemID)
	if err != nil {
		logger.Log.Warn("Некорректный uuid в запросе: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	err = h.itemService.ChangeItem(itemID, req)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Изменение товара прошло успешно")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}

// DeleteItem обрабатывает запрос на удаленение
// определенного товара из списка покупок
func (h *ItemHandler) DeleteItem(c *gin.Context) {

	rawItemID := c.Param("item_id")

	itemID, err := uuid.Parse(rawItemID)
	if err != nil {
		logger.Log.Warn("Некорректный uuid в запросе: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	err = h.itemService.DeleteItems([]uuid.UUID{itemID})

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Удаление товара прошло успешно")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}

// DeleteItems обрабатывает запрос на удаленение
// нескольких товаров из списка покупок
func (h *ItemHandler) DeleteItems(c *gin.Context) {

	var req requests.DeleteItemsRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Log.Warn("Ошибка при десериализации запроса: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	err = h.itemService.DeleteItems(req.ItemIDs)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Удаление товаров прошло успешно")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}

// MarkAsBought обрабатывает запрос для отметки
// определенного товара как купленного
func (h *ItemHandler) MarkAsBought(c *gin.Context) {
	var req requests.MarkAsBoughtRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		logger.Log.Warn("Ошибка при десериализации запроса: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	rawItemID := c.Param("item_id")

	itemID, err := uuid.Parse(rawItemID)
	if err != nil {
		logger.Log.Warn("Некорректный uuid в запросе: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	err = h.itemService.MarkAsBought(itemID, req.Checked)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Товар успешно помечен как купленный")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}

// NotifyMembers обрабатывает запрос на
// уведомление членов семьи об определенном товаре
func (h *ItemHandler) NotifyMembers(c *gin.Context) {

	var req requests.NotifyMembersRequest

	ctx := c.Request.Context()

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

	rawItemID := c.Param("item_id")

	itemID, err := uuid.Parse(rawItemID)
	if err != nil {
		logger.Log.Warn("Некорректный uuid в запросе: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	err = h.itemService.NotifyMembers(ctx, userID, itemID, req)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Уведомление о товаре прошло успешно")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}
