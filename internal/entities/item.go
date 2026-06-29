// entities - пакет с сущностями для БД
package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// HighestItemPriority задает максимальный
// возможный приоритет товара
const HighestItemPriority = 3
const DefaultItemPriority = 1

// OthersSectionName задает название секции, куда
// помещаются продукты без явно указанной секции
const OthersSectionName = "Прочее"

// Item - структура сущности продукта из списка покупок.
// Содержит ссылку на список, название продукта, приоритет,
// кем был добавлен, количество, время создания и изменения.
type Item struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	ListID uuid.UUID `gorm:"type:uuid;not null"`

	Name    string `gorm:"not null"`
	Section string

	Quantity float64 `gorm:"not null"`
	Unit     string

	Priority int `gorm:"default:1"`
	Urgent   bool

	Checked         bool `gorm:"default:false"`
	RemindEveryDays *int

	AddedByID *uuid.UUID `gorm:"type:uuid;not null"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// BeforeCreate создает необходимые отсутствующие поля при создании сущности
func (i *Item) BeforeCreate(db *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	if i.CreatedAt.IsZero() {
		i.CreatedAt = time.Now().UTC()
	}
	if i.UpdatedAt.IsZero() {
		i.UpdatedAt = time.Now().UTC()
	}
	return nil
}
