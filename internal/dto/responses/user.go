// responses - пакет, содержащий структуры ответов на Api запросы
package responses

import (
	"time"

	"github.com/google/uuid"
)

// UserInfo - структура, содержащая краткую
// информацию о пользователе
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Surname  string    `json:"surname"`
	ImageURL string    `json:"image-url,omitempty"`
}

// FamilyMemberInfo - структура, содержащая
// полную информаци. об одном члене семьи
type FamilyMemberInfo struct {
	UserInfo
	JoinedAt time.Time `json:"joined-at"`
	Color    string    `json:"color"`
}

// GetMyInfoResponse - структура для ответа на
// запрос получения информации текущего пользователя
type GetMyInfoResponse struct {
	UserInfo
	Lists []ListInfo `json:"lists"`
}
