// responses - пакет, содержащий структуры ответов на Api запросы
package responses

import (
	"time"

	"github.com/google/uuid"
)

// FamilyInfo - структура для получения полной
// информации о конкретной семье
type FamilyInfo struct {
	ID uuid.UUID `json:"id"`

	Owner   FamilyMemberInfo   `json:"owner"`
	Members []FamilyMemberInfo `json:"members"`

	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
}

// GetFamiliesResponse - структура для ответа на запрос
// для получения информации о всех семьях текущего пользователя
type GetFamiliesResponse struct {
	Families []FamilyInfo `json:"families"`
}

// GetListsResponse - структура для ответа на
// запрос для получения информации о всех списках
// продуктов определенной семьи
type GetListsResponse struct {
	Lists []ListInfo `json:"lists"`
}

// GetFamiliesResponse - структура для ответа на запрос
// для получения коды и ссылки на определенную семью
type GetFamilyLinkResponse struct {
	LinkCode string `json:"link-code"`
	LinkURL  string `json:"link-url"`
}
