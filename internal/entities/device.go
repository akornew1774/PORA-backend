// entities - пакет с сущностями для БД
package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceType задает тип устройства пользователя
type DeviceType string

const (
	AndroidDevice DeviceType = "android"
	IOSDevice     DeviceType = "ios"
)

// Device - структура сущности устройства пользователя.
// Содержит токен и тип устройства пользователя
type Device struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	DeviceToken string     `gorm:"not null;unique"`
	DeviceType  DeviceType `gorm:"not null"`

	CreatedAt time.Time `gorm:"not null"`
}

// BeforeCreate создает необходимые отсутствующие поля при создании сущности
func (d *Device) BeforeCreate(db *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now().UTC()
	}
	return nil
}
