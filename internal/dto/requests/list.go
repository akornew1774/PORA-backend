// requests - пакет, содержащий структуры запросов по Api
package requests

import (
	"github.com/google/uuid"
)

// CreateList - структура для запроса
// для создания нового списка покупок
type CreateListRequest struct {
	FamilyID *uuid.UUID `json:"family-id"`
	Name     string     `json:"name" binding:"required"`
}
