// entities - пакет с сущностями для БД
package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserStatus обозначает текущий статус пользователя
type UserStatus string

const (
	UserStatusNotFound      UserStatus = "notFound"
	UserStatusNotRegistered UserStatus = "notRegistered"
	UserStatusActive        UserStatus = "user"
	UserStatusAdmin         UserStatus = "admin"
	UserStatusDeleted       UserStatus = "deleted"
)

// User - структура сущности пользователя в БД
// Содержит ID, телефон, имя и фамилию пользователя,
// его статус, ссылку на refresh токен
type User struct {
	ID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Phone string    `gorm:"unique"`
	Email string    `gorm:"unique"`

	Name    string
	Surname string
	Status  UserStatus `gorm:"not null"`

	ImageURL string

	RefreshToken *RefreshToken `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	Memberships []FamilyMember `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Lists       []List         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// BeforeCreate создает необходимые отсутствующие поля при создании сущности
func (u *User) BeforeCreate(db *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if len(u.Status) == 0 {

		if len(u.Name) == 0 && len(u.Surname) == 0 {
			u.Status = UserStatusNotRegistered
		} else {
			u.Status = UserStatusActive
		}
	}

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now().UTC()
	}
	return nil
}
