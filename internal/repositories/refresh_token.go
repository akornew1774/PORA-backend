// repositories - пакет с методами для работы с БД
package repositories

import (
	"errors"
	"pora/internal/entities"
	"pora/internal/ports"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshTokenRepo - объект, содержащий методы для работы с сущностью refresh-токена
type RefreshTokenRepo struct {
	db *gorm.DB
}

// NewRefreshTokenRepository создает и возвращает новый объект RefreshTokenRepo
func NewRefreshTokenRepository(db *gorm.DB) ports.RefreshTokenRepository {
	return &RefreshTokenRepo{db: db}
}

// FindByToken находит сущность refresh-токена по его значению
func (r *RefreshTokenRepo) FindByToken(
	token string) (*entities.RefreshToken, error) {

	var refreshToken entities.RefreshToken

	err := r.db.Where("token = ?", token).
		First(&refreshToken).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// CreateRefreshToken создает новый объект токена в БД
func (r *RefreshTokenRepo) CreateRefreshToken(
	refreshToken *entities.RefreshToken) error {
	return r.db.Create(refreshToken).Error
}

// UpdateRefreshToken изменяет данные определенного refresh-токена
func (r *RefreshTokenRepo) UpdateRefreshToken(
	refreshToken *entities.RefreshToken) error {
	return r.db.Save(refreshToken).Debug().Error
}

// DeleteByUserID удаляет все refresh-токены c определенным userID
func (r *RefreshTokenRepo) DeleteByUserID(
	userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).
		Delete(&entities.RefreshToken{}).Error
}
