// entities - пакет с сущностями для БД
package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BriefItem - структура краткого описания товара
type BriefItem struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`

	Title    string `gorm:"not null"`
	Subtitle *string

	CreatedAt time.Time `gorm:"not null"`
}

// BeforeCreate создает необходимые отсутствующие поля при создании сущности
func (b *BriefItem) BeforeCreate(db *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now().UTC()
	}
	return nil
}
