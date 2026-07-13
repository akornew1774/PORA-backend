// websocket - пакет, содержащий реализацию постоянного
// соединения между сервером и мобильными устройствами
package websocket

import (
	"github.com/google/uuid"
)

// Message - структура для отправки данных клиенту по WS
type SendChangesMessage struct {
	FamilyID *uuid.UUID `json:"family-id,omitempty"`
	ListID   *uuid.UUID `json:"list-id,omitempty"`
	ItemID   *uuid.UUID `json:"item-id,omitempty"`
}
