// entities - пакет с сущностями для БД
package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Otp - структура сущности OTP кода.
// Содержит ID, код, телефон пользователя и время истечения срока действия
type Otp struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Code  string `gorm:"not null"`
	Phone string `gorm:"not null"`

	ExpiresAt time.Time `gorm:"not null"`
	IsUsed    bool      `gorm:"default:false"`
}

// BeforeCreate создает необходимые отсутствющие поля при создании сущности
func (o *Otp) BeforeCreate(db *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	if o.ExpiresAt.IsZero() {
		o.ExpiresAt = time.Now().UTC().Add(2 * time.Minute)
	}
	return nil
}
