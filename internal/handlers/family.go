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

// FamilyHandler - объект, содержащий методы для обработки
// Api запросов, связанных с семьями
type FamilyHandler struct {
	familyService ports.FamilyService
	tokenService  ports.TokenService
}

// NewFamilyHandler создает и возвращает новый объект FamilyHandler
func NewFamilyHandler(familyService ports.FamilyService,
	tokenService ports.TokenService) *FamilyHandler {
	return &FamilyHandler{
		familyService: familyService,
		tokenService:  tokenService,
	}
}

// GetFamilies обрабатывает запрос на
// получение всех семей текущего пользователя
func (h *FamilyHandler) GetFamilies(c *gin.Context) {

	accessToken := c.GetHeader("Authorization")

	userID, err := h.tokenService.DecodeAccessToken(accessToken)
	if err != nil {
		logger.Log.Warn("Возникла ошибка при декодировании Access-токена")
		c.Error(err)
		return
	}

	response, err := h.familyService.GetFamilies(userID)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Получение семей пользователя прошло успешно")
	c.JSON(http.StatusOK, response)
}

// GetLists обрабатывает запрос на получение
// всех списков покупок определенной семьи
func (h *FamilyHandler) GetLists(c *gin.Context) {

	rawFamilyID := c.Param("family_id")

	familyID, err := uuid.Parse(rawFamilyID)
	if err != nil {
		logger.Log.Warn("Некорректный uuid в запросе: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	response, err := h.familyService.GetLists(familyID)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Получение всех списков продуктов семьи прошло успешно")
	c.JSON(http.StatusOK, response)
}

// CreateFamily обрабатывает запрос на
// создание новой семьи текущим пользователем
func (h *FamilyHandler) CreateFamily(c *gin.Context) {

	var req requests.NameRequest

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

	response, err := h.familyService.CreateFamily(userID, req.Name)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Создание новой семьи успешно выполнено")
	c.JSON(http.StatusOK, response)
}

// GetFamilyLink обрабатывает запрос на получение
// уникального кода семьи и ссылки для вступления в нее
func (h *FamilyHandler) GetFamilyLink(c *gin.Context) {

	var req requests.FamilyLinkRequest

	err := c.ShouldBindQuery(&req)
	if err != nil {
		logger.Log.Warn("Ошибка при десериализации запроса: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	familyID, err := uuid.Parse(req.FamilyID)
	if err != nil {
		logger.Log.Warn("Некорректный uuid в запросе: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	response, err := h.familyService.GetFamilyLink(familyID)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Получение ссылки на семью прошло успешно")
	c.JSON(http.StatusOK, response)
}

// JoinFamily обрабатывает запрос на присоединение к семье
func (h *FamilyHandler) JoinFamily(c *gin.Context) {

	familyCode := c.Param("link_code")

	accessToken := c.GetHeader("Authorization")

	userID, err := h.tokenService.DecodeAccessToken(accessToken)
	if err != nil {
		logger.Log.Warn("Возникла ошибка при декодировании Access-токена")
		c.Error(err)
		return
	}

	err = h.familyService.AddMember(userID, familyCode)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Присоединение пользователя к семье прошло успешно")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}

// DeleteFamily обрабатывает запрос на удаление семьи
func (h *FamilyHandler) DeleteFamily(c *gin.Context) {

	rawFamilyID := c.Param("family_id")

	familyID, err := uuid.Parse(rawFamilyID)
	if err != nil {
		logger.Log.Warn("Некорректный uuid в запросе: ", err)
		c.Error(errors.ErrorInvalidInput)
		return
	}

	err = h.familyService.DeleteFamily(familyID)

	if err != nil {
		c.Error(err)
		return
	}

	logger.Log.Info("Удаление семьи успешно выполнено")
	c.JSON(http.StatusOK, responses.GenericResponse{})
}
