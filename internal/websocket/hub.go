// websocket - пакет, содержащий реализацию постоянного
// соединения между сервером и мобильными устройствами
package websocket

import (
	"pora/internal/infrastructure/logger"
	"sync"

	"github.com/google/uuid"
)

// Hub управляет активными Websocket-соединениями пользователей
type Hub struct {
	mu sync.RWMutex

	// clients хранит соединения каждого пользователя
	clients map[uuid.UUID]map[*Client]struct{}

	// owners хранит владельца каждого соединения
	owners map[*Client]uuid.UUID
}

// NewHub создает и возвращает новый Hub
func NewHub() *Hub {
	return &Hub{
		clients: make(map[uuid.UUID]map[*Client]struct{}),
		owners:  make(map[*Client]uuid.UUID),
	}
}

// Register регистрирует новое websocket-соединение с клиентом
func (h *Hub) Register(userID uuid.UUID, client *Client) {

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[userID]; !ok {
		h.clients[userID] = make(map[*Client]struct{})
	}

	h.clients[userID][client] = struct{}{}
	h.owners[client] = userID
}

// Unregister удаляет websocket-соединение с клиентом
func (h *Hub) Unregister(client *Client) {

	h.mu.Lock()
	defer h.mu.Unlock()

	userID, ok := h.owners[client]
	if !ok {
		return
	}

	delete(h.owners, client)

	userClients := h.clients[userID]
	delete(userClients, client)

	if len(userClients) == 0 {
		delete(h.clients, userID)
	}

	if err := client.Close(); err != nil {
		logger.Log.Error("Ошибка при закрытии WS-соединения")
	}
}

// SendToUser отправляет сообщение всем активным
// соединениям указанного пользователя
func (h *Hub) SendToUser(userID uuid.UUID, message []byte) {

	h.mu.RLock()
	defer h.mu.RUnlock()

	userClients, ok := h.clients[userID]
	if !ok {
		return
	}

	for client := range userClients {
		err := client.Send(message)
		if err != nil {
			logger.Log.Error("Ошибка при отправке сообщения через WS")
		}
	}
}
