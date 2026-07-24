// entities - пакет с сущностями для БД
package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserLogin - структура, обозначающая
// время входа пользователя в проложение
type UserLogin struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

// BeforeCreate создает необходимые отсутствующие поля при создании сущности
func (u *UserLogin) BeforeCreate(db *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	return nil
}
