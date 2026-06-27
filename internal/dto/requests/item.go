// requests - пакет, содержащий структуры запросов по Api
package requests

import (
	"github.com/google/uuid"
)

// AddItemRequest - структура для запроса на
// добавление товара в список покупок
type AddItemRequest struct {
	Name    string `json:"name" binding:"required"`
	Section string `json:"section"`

	Quantity float64 `json:"quantity" binding:"required"`
	Unit     string  `json:"unit" binding:"required"`

	Priority int  `json:"priority" binding:"required"`
	Urgent   bool `json:"urgent" binding:"required"`

	Checked         bool `json:"checked"`
	RemindEveryDays *int `json:"remind-every-days"`
}

// ChangeItemRequest - структура дял запроса на
// изменение товара в списке покупок
type ChangeItemRequest struct {
	Name    string `json:"name"`
	Section string `json:"section"`

	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`

	Priority int  `json:"priority"`
	Urgent   bool `json:"urgent"`

	Checked         bool `json:"checked"`
	RemindEveryDays *int `json:"remind-every-days"`
}

// NotifyMembersRequest - структура для уведомления
// членов семьи о конкретном продукте из списка
type NotifyMembersRequest struct {
	To      *[]uuid.UUID `json:"to"`
	Message string       `json:"message"`
}
