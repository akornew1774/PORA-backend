// ports - пакет, содержащий все порты (интерфейсы)
package ports

import (
	"pora/internal/websocket"

	"github.com/google/uuid"
)

// Hub содержит порты для работы менеджера
// всех активных Websocket-соединений
type Hub interface {
	// Register регистрирует новое websocket-соединение с клиентом
	Register(userID uuid.UUID, client *websocket.Client)

	// Unregister удаляет websocket-соединение с клиентом
	Unregister(client *websocket.Client)

	// SendToUser отправляет сообщение всем активным
	// соединениям указанного пользователя
	SendToUser(userID uuid.UUID, message []byte)
}
