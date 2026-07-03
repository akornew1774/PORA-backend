// requests - пакет, содержащий структуры запросов по Api
package requests

// PhoneAndEmailRequest - структура для запроса,
// содержащего только номер телефона и Email в Params или Body
type PhoneAndEmailRequest struct {
	Phone string `json:"phone" form:"phone"`
	Email string `json:"email" form:"email"`
}

// UserIDRequest - структура для запросов, содержащих
// только ID пользователя в Params или Body
type UserIDRequest struct {
	UserID string `json:"user-id" form:"user-id" binding:"required"`
}

// NameRequest - структура для запросов, содержащих
// только название/имя нужной структуры в Params или Body
type NameRequest struct {
	Name string `json:"name" form:"name" binding:"required"`
}
