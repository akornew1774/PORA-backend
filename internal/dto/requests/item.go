// requests - пакет, содержащий структуры запросов по Api
package requests

import (
	"github.com/google/uuid"
)

// ChangeItemRequest - структура дял запроса на
// изменение товара в списке покупок
type ChangeItemRequest struct {
	Name    string `json:"name" binding:"required"`
	Section string `json:"section"`

	Quantity float64 `json:"quantity" binding:"required"`
	Unit     string  `json:"unit" binding:"required"`

	Priority int  `json:"priority" binding:"required"`
	Urgent   bool `json:"urgent" binding:"required"`

	Checked         bool `json:"checked"`
	RemindEveryDays *int `json:"remind-every-days"`
}

// NotifyMembersRequest - структура для уведомления
// членов семьи о конкретном продукте из списка
type NotifyMembersRequest struct {
	FamilyID uuid.UUID   `json:"family-id" binding:"required"`
	To       []uuid.UUID `json:"to"`
	Message  string      `json:"message"`
}
