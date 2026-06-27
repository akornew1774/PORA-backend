// entities - пакет с сущностями для БД
package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Family - структура сущности семьи. Содержит
// название, ссылку на владельца, массив участников
// семьи и массив списков покупок.
type Family struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerID uuid.UUID `gorm:"type:uuid;not null"`

	Name string `gorm:"not null"`

	Members []FamilyMember `gorm:"foreignKey:FamilyID;constraint:OnDelete:CASCADE"`
	Lists   []List         `gorm:"foreignKey:FamilyID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"not null"`
}

// BeforeCreate создает необходимые отсутствующие поля при создании сущности
func (f *Family) BeforeCreate(db *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now().UTC()
	}
	return nil
}
