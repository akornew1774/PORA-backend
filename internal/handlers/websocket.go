// handlers - пакет, содержащий в себе обработчики Api запросов
package handlers

import (
	"pora/internal/infrastructure/logger"
	wsInfrastructure "pora/internal/infrastructure/websocket"
	"pora/internal/ports"
	"pora/internal/websocket"

	"github.com/gin-gonic/gin"
)

// WSHandler - объект, содержащий методы для создания
// Websocket-соединения между сервером и мобильным устройством
type WSHandler struct {
	hub          ports.Hub
	tokenService ports.TokenService
}

// NewWSHandler создает и возвращает новый объект WSHandler
func NewWSHandler(hub ports.Hub,
	tokenService ports.TokenService) *WSHandler {
	return &WSHandler{
		hub:          hub,
		tokenService: tokenService,
	}
}

// Connect устанавливает WebSocket-соединение
func (h *WSHandler) Connect(c *gin.Context) {

	accessToken := c.GetHeader("Authorization")

	userID, err := h.tokenService.DecodeAccessToken(accessToken)
	if err != nil {
		logger.Log.Warn("Возникла ошибка при декодировании Access-токена")
		c.Error(err)
		return
	}

	conn, err := wsInfrastructure.Upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)
	if err != nil {
		logger.Log.Error("Ошибка при подключении WS-соединения: ", err)
		c.Error(err)
		return
	}

	client := websocket.NewClient(conn)

	h.hub.Register(userID, client)

	go func() {
		defer h.hub.Unregister(client)
		client.ReadPump()
	}()

	go func() {
		defer h.hub.Unregister(client)
		client.WritePump()
	}()
}
