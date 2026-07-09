// entities - пакет с сущностями для БД
package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FamilyRole обозначает роль участника в семье
type FamilyRole string

const (
	MemberRole FamilyRole = "member"
	OwnerRole  FamilyRole = "owner"
)

// FamilyMember - структура сущности членства в семье.
// Представляет собой связь между пользователем и семьей,
// содержит ссылки на пользователя и семью, цвет пользователя в семье.
type FamilyMember struct {
	UserID   uuid.UUID `gorm:"type:uuid;primaryKey"`
	FamilyID uuid.UUID `gorm:"type:uuid;primaryKey"`

	Role  FamilyRole  `gorm:"not null"`
	Color MemberColor `gorm:"not null"`

	User   *User   `gorm:"constraint:OnDelete:CASCADE"`
	Family *Family `gorm:"constraint:OnDelete:CASCADE"`

	JoinedAt time.Time `gorm:"not null"`
}

// BeforeCreate создает необходимые отсутствующие поля при создании сущности
func (f *FamilyMember) BeforeCreate(db *gorm.DB) error {
	if f.JoinedAt.IsZero() {
		f.JoinedAt = time.Now().UTC()
	}
	return nil
}
