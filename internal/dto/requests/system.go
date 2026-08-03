// requests - пакет, содержащий структуры запросов по Api
package requests

// NotifyEveryone - структура для запроса на
// уведомление всех пользователей приложения
type NotifyEveryoneRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}
