// responses - пакет, содержащий структуры ответов на Api запросы
package responses

import (
	"github.com/google/uuid"
)

// SectionInfo - структура, содержащая информацию
// о секции в списке продуктов
type SectionInfo struct {
	Name  string     `json:"name"`
	Items []ItemInfo `json:"items"`
}

// ItemInfo - структура, содержащая информацию
// о конкретном товаре в списке продуктов
type ItemInfo struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Section string    `json:"section,omitempty"`

	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`

	Priority int  `json:"priority"`
	Urgent   bool `json:"urgent"`

	Checked         bool `json:"checked"`
	RemindEveryDays *int `json:"remind-every-days"`

	AddedBy *FamilyMemberInfo `json:"added-by,omitempty"`
}
