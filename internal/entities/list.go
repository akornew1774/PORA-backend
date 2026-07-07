// entities - пакет с сущностями для БД
package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// List - структура сущности списка покупок. Содержит
// название, ссылку на пользователя или группу,
// к которой он относится, а также массив продуктов.
type List struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserID   *uuid.UUID `gorm:"type:uuid"`
	FamilyID *uuid.UUID `gorm:"type:uuid"`

	Name string `gorm:"not null"`

	Items []Item `gorm:"foreignKey:ListID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// BeforeCreate создает необходимые отсутствующие поля при создании сущности
func (l *List) BeforeCreate(db *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	if l.CreatedAt.IsZero() {
		l.CreatedAt = time.Now().UTC()
	}
	if l.UpdatedAt.IsZero() {
		l.UpdatedAt = time.Now().UTC()
	}
	return nil
}
