// requests - пакет, содержащий структуры запросов по Api
package requests

import (
	"github.com/google/uuid"
)

// ChangeItemRequest - структура для запроса на
// изменение товара в списке покупок
type ChangeItemRequest struct {
	Name    string `json:"name" binding:"required"`
	Section string `json:"section"`

	Quantity float64 `json:"quantity" binding:"required"`
	Unit     string  `json:"unit" binding:"required"`

	Priority int  `json:"priority"`
	Urgent   bool `json:"urgent"`

	Checked         bool `json:"checked"`
	RemindEveryDays *int `json:"remind-every-days"`
}

// MarkAsBoughtRequest - структура для запроса на
// отметку товара как купленного
type MarkAsBoughtRequest struct {
	Checked bool `json:"checked"`
}

// NotifyMembersRequest - структура для уведомления
// членов семьи о конкретном продукте из списка
type NotifyMembersRequest struct {
	To      []uuid.UUID `json:"to"`
	Message string      `json:"message"`
}
