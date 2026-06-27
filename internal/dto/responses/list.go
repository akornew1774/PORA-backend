// responses - пакет, содержащий структуры ответов на Api запросы
package responses

import (
	"time"

	"github.com/google/uuid"
)

// ListInfo - структура, содержащая полную
// информацию о списке продуктов
type ListInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`

	Sections []SectionInfo `json:"sections"`

	CreatedAt time.Time `json:"created-at"`
}
