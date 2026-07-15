// notifications - пакет, содержащий модели
// для отправки уведомлений пользователям
package notifications

// Notification - структура для
// отправки Push-уведомлений пользователям
type Notification struct {
	Title string
	Body  string
	Data  map[string]string
}
