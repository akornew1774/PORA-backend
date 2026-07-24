// requests - пакет, содержащий структуры запросов по Api
package requests

// HelpMessageRequest - структура для запроса для
// отправки сообщения с просьбой помощи
type HelpMessageRequest struct {
	Title   string `json:"title" binding:"required"`
	Message string `json:"message" binding:"required"`
}
