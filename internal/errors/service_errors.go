// errors - пакет с возможными ошибками приложения
package errors

// AppError - структура, общая для всех видов ошибок
type AppError struct {
	Code    string
	Message string
}

// Error реализует интерфейс ошибки.
// Необходим, чтобы AppError классифицировалось как ошибка
func (e AppError) Error() string {
	return e.Message
}
