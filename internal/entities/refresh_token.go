// entities - пакет с сущностями для БД
package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshToken - структура сущности refresh токена.
// Содержит ID, ID привязанного к нему пользователя,
// сам токен, время истечения срока действия
type RefreshToken struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;uniqueIndex"`

	Token     string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

// BeforeCreate создает необходимые отсутствющие поля при создании сущности
func (r *RefreshToken) BeforeCreate(db *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	if r.ExpiresAt.IsZero() {
		r.ExpiresAt = time.Now().UTC().Add(15 * 24 * time.Hour)
	}
	return nil
}
