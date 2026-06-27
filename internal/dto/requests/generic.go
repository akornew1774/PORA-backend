// requests - пакет, содержащий структуры запросов по Api
package requests

// PhoneNumberRequests - структура для запроса,
// содержащего только номер телефона в Params или Body
type PhoneNumberRequest struct {
	Phone string `json:"phone" form:"phone" binding:"required"`
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
