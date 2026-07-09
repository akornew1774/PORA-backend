// requests - пакет, содержащий структуры запросов по Api
package requests

// UpdateRequest - структура для запроса
// на обновление профиля пользователя
type UpdateUserRequest struct {
	Phone *string `json:"phone,omitempty"`
	Email *string `json:"email,omitempty"`

	Name    *string `json:"name,omitempty"`
	Surname *string `json:"surname,omitempty"`
}
