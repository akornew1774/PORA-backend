// notifications - пакет, содержащий модели
// для отправки уведомлений пользователям
package notifications

import (
	"github.com/google/uuid"
)

// Notification - структура для
// отправки Push-уведомлений пользователям
type Notification struct {
	Title string
	Body  string
	Data  map[string]string
}

// ItemNotification - структура, содержащая необходимые
// для отправки уведомления от пользователя данные
type ItemNotification struct {
	UserID uuid.UUID

	AuthorID   uuid.UUID
	AuthorName string

	ItemID       uuid.UUID
	ItemName     string
	ItemQuantity string
	ItemUnit     string

	ListID   uuid.UUID
	ListName string

	FamilyID   uuid.UUID
	FamilyName string

	Message string
}

// ItemReminder - структура, содержащая необходимые
// для отправки напоминания о продукте данные
type ItemReminder struct {
	UserID uuid.UUID

	ItemID       uuid.UUID
	ItemName     string
	ItemQuantity string
	ItemUnit     string

	ListID   uuid.UUID
	ListName string

	FamilyID   *uuid.UUID
	FamilyName *string
}
