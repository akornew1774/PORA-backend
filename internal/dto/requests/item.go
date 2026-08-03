// requests - пакет, содержащий структуры запросов по Api
package requests

import (
	"github.com/google/uuid"
)

// AddItemsRequest - структура запроса на
// добавление нескольких товаров в список покупок
type AddItemsRequest struct {
	Items []ChangeItemRequest `json:"items" binding:"required"`
}

// ChangeItemRequest - структура для запроса на
// изменение или удаление товара в списке покупок
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

// DeleteItemsRequest - структура для запроса на
// удаление нескольких товаров из списка продуктов
type DeleteItemsRequest struct {
	ItemIDs []uuid.UUID `json:"item-ids" binding:"required"`
}

// NotifyMembersRequest - структура для уведомления
// членов семьи о конкретном продукте из списка
type NotifyMembersRequest struct {
	To      []uuid.UUID `json:"to"`
	Message string      `json:"message"`
}
