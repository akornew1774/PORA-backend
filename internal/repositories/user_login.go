// repositories - пакет с методами для работы с БД
package repositories

import (
	"pora/internal/entities"
	"pora/internal/ports"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserLoginRepo - объект, содержащий методы для работы
// с сущностью входа пользователя в приложения в БД
type UserLoginRepo struct {
	db *gorm.DB
}

// NewUserLoginRepository создает и возвращает новый объект UserLoginRepo
func NewUserLoginRepository(db *gorm.DB) ports.UserLoginRepository {
	return &UserLoginRepo{db: db}
}

// FindRecentByUserID находит все недавние входы пользователя в приложение
func (r *UserLoginRepo) FindRecentByUserID(userID uuid.UUID,
	since time.Time) ([]entities.UserLogin, error) {

	var logins []entities.UserLogin

	err := r.db.
		Where("user_id = ?", userID).
		Where("created_at >= ?", since).
		Order("created_at DESC").
		Find(&logins).
		Error

	return logins, err
}

// CreateUserLogin создает новую сущность входа пользователя в приложение
func (r *UserLoginRepo) CreateUserLogin(
	login *entities.UserLogin) error {

	return r.db.Create(login).Error
}

// DeleteOldLogins удаляет все UserLogins,
// которые были созданы раньше граничного времени
func (r *UserLoginRepo) DeleteOldLogins(
	cutoffTime time.Time) error {

	err := r.db.
		Where("created_at <= ?", cutoffTime).
		Delete(&entities.UserLogin{}).
		Error

	return err
}
