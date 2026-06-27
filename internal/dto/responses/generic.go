// responses - пакет, содержащий структуры ответов на Api запросы
package responses

import (
	"github.com/google/uuid"
)

// GenericResponse - Общая структура для ответа на запросы.
// Используется для пустых ответов / ответов с message
type GenericResponse struct {
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// SaveImageResponse - структура для ответа на запросы
// для сохранения изображения в хранилище
type SaveImageResponse struct {
	ImageURL string `json:"image-url"`
}

// SaveImageResponse - структура для ответа на запросы,
// в которых требуется возвратить только ID сущности
type IDResponse struct {
	ID uuid.UUID `json:"id"`
}
